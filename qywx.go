package qywx

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"sync"

	wxbiz "qywx/wxbizjsonmsgcrypt"

	"strconv"
	"strings"
	"time"
)

type QYWXClient struct {
	base  *BaseInfo
	wxcpt *wxbiz.WXBizMsgCrypt
}

func NewQYWXClient(ctx context.Context, base *BaseInfo) *QYWXClient {
	cli := &QYWXClient{
		base: base,
	}
	if base.l == nil {
		base.l = new(sync.RWMutex)
	}
	cli.InitWXBizMsgCrypt(ctx)
	return cli
}

type (
	AccessToken struct {
		Token     string
		ExpiresIn int64
	}
	CallbackRespFunc func(ctx context.Context, msgTyp CBType, msg []byte) ([]byte, error)
	SendMsgBodyFunc  func(ctx context.Context, msg *BaseMsg) ([]byte, error)
)

func NewSendMsgBodyFunc(c *QYWXClient, typ SMType) SendMsgBodyFunc {
	efunc := func(msg interface{}) ([]byte, error) {
		return json.Marshal(msg)
	}
	switch typ {
	case TextType:
		return func(ctx context.Context, msg *BaseMsg) ([]byte, error) {
			body := msg.MsgBody.(*SendMessageTextType)
			body.SendMessageBaseType = SendMessageBaseType{
				ToUser:                 strings.Join(c.GetToUserList(ctx), "|"),
				MsgType:                "text",
				AgentId:                c.base.GetAgentId(),
				EnableIdTrans:          0,
				EnableDuplicateCheck:   0,
				DuplicateCheckInterval: 1800,
			}
			return efunc(body)
		}
	case CardType:
		return func(ctx context.Context, msg *BaseMsg) ([]byte, error) {
			body := msg.MsgBody.(*SendMessageTemplateCardType)
			body.SendMessageBaseType = SendMessageBaseType{
				ToUser:                 strings.Join(c.GetToUserList(ctx), "|"),
				MsgType:                "template_card",
				AgentId:                c.base.GetAgentId(),
				EnableIdTrans:          0,
				EnableDuplicateCheck:   0,
				DuplicateCheckInterval: 1800,
			}
			return efunc(body)
		}
	default:
		return nil
	}
}

func (c *QYWXClient) InitWXBizMsgCrypt(ctx context.Context) {
	if c.base == nil || c.base.l == nil {
		return
	}
	c.base.l.RLock()
	defer c.base.l.RUnlock()
	c.wxcpt = wxbiz.NewWXBizMsgCrypt(c.base.Token, c.base.EncodingAeskey, c.base.ReceiverId, wxbiz.JsonType)
}

func (c *QYWXClient) GetWXBizMsgCrypt(ctx context.Context) *wxbiz.WXBizMsgCrypt {
	if c.wxcpt == nil {
		c.InitWXBizMsgCrypt(ctx)
	}
	return c.wxcpt
}

func (c *QYWXClient) GetToUserList(ctx context.Context) []string {
	if c.base == nil || c.base.ToUserList == nil {
		return []string{}
	}
	return c.base.ToUserList
}

func (c *QYWXClient) SetBaseInfo(ctx context.Context, b *BaseInfo) {
	if c.base == nil {
		c.base = b
	} else {
		c.base.SetBaseInfo(b)
	}
	c.InitWXBizMsgCrypt(ctx)
}

func (c *QYWXClient) GetAuthorization(ctx context.Context, access *AccessToken) error {
	if access.ExpiresIn > time.Now().Unix() {
		return nil
	}
	if c.base.CorpId == "" || c.base.CorpSecret == "" {
		return errors.New("authorization info of QiyeWeixin gettoken api init failed")
	}
	var url strings.Builder
	url.WriteString(QWURL)
	url.WriteString("/cgi-bin/gettoken?")
	url.WriteString("corpid=")
	url.WriteString(c.base.CorpId)
	url.WriteString("&corpsecret=")
	url.WriteString(c.base.CorpSecret)
	resp, err := NetGET(ctx, url.String())
	if err != nil {
		return err
	}
	r := new(GetAuthorizationApiResponse)
	err = json.Unmarshal(resp, r)
	if err != nil {
		return err
	}
	ExpiresIn, err := time.ParseDuration(fmt.Sprintf("%ds", r.ExpiresIn-60))
	if err != nil {
		return err
	}
	access.ExpiresIn = time.Now().Add(ExpiresIn).Unix()
	access.Token = r.AccessToken
	return err
}

func (c *QYWXClient) SendMessage(ctx context.Context, access *AccessToken, msg *BaseMsg) (*SendMessageResponse, error) {
	err := c.GetAuthorization(ctx, access)
	if nil != err {
		return nil, err
	}
	var url strings.Builder
	url.WriteString(QWURL)
	url.WriteString("/cgi-bin/message/send")
	url.WriteString(fmt.Sprintf("?access_token=%s", access.Token))
	smresp := new(SendMessageResponse)
	reqBody, err := NewSendMsgBodyFunc(c, msg.SMMsgType)(ctx, msg)
	if err != nil {
		return smresp, err
	}
	resp, err := NetPOSTJson(ctx, url.String(), reqBody)
	if err != nil {
		return smresp, err
	}
	err = json.Unmarshal(resp, smresp)
	if err != nil {
		return smresp, err
	}
	fmt.Println(string(resp))
	return smresp, nil
}

func (c *QYWXClient) MakeResponseData(ctx context.Context, msg string, timestamp int64) (string, error) {
	wxcpt := c.GetWXBizMsgCrypt(ctx)
	reqNonce := fmt.Sprintf("Nonce_%v", wxcpt.GetRandString(16))
	timeStamp := strconv.Itoa(int(timestamp))
	encrypt, sig, err := wxcpt.GetEncryptMsg(msg, timeStamp, reqNonce)
	if err != nil {
		return "", errors.New(err.ErrMsg)
	}
	respBody := fmt.Sprintf(`<xml><Encrypt><![CDATA[%v]]></Encrypt><MsgSignature><![CDATA[%v]]></MsgSignature><TimeStamp>%v</TimeStamp><Nonce><![CDATA[%v]]></Nonce></xml>`, encrypt, sig, timeStamp, reqNonce)
	return respBody, nil
}

func (c *QYWXClient) CallbackCheck(ctx context.Context, req *CallbackCheckRequest) ([]byte, error) {
	verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr := req.MsgSignature, req.Timestamp, req.Nonce, req.Echostr
	echoStr, cryptErr := c.GetWXBizMsgCrypt(ctx).VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	if nil != cryptErr {
		return nil, errors.New(cryptErr.ErrMsg)
	}
	return echoStr, nil
}

func (c *QYWXClient) callbackDec(ctx context.Context, req *CallbackRequest) ([]byte, error) {
	reqMsgSign, reqTimestamp, reqNonce := req.MsgSignature, req.Timestamp, req.Nonce
	body := new(CallbackRequestBody)
	err := xml.Unmarshal(req.Body, body)
	if err != nil {
		return nil, err
	}
	reqDataString := fmt.Sprintf(`{"tousername":"%s","encrypt":"%s","agentid":"%s"}`, body.ToUser, body.MsgEncrypt, body.ToAgentID)
	reqData := []byte(reqDataString)

	msg, cryptErr := c.GetWXBizMsgCrypt(ctx).DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, reqData)
	if nil != cryptErr {
		return nil, errors.New(cryptErr.ErrMsg)
	}
	return msg, nil
}

func (c *QYWXClient) CallbackResp(ctx context.Context, typ CBType, req *CallbackRequest, cb CallbackRespFunc) ([]byte, error) {
	m, err := c.callbackDec(ctx, req)
	if err != nil {
		return nil, err
	}
	showinfo, err := cb(ctx, typ, m)
	if err != nil {
		return nil, err
	}
	if showinfo == nil {
		showinfo = []byte("")
	}
	toUser := strings.Join(c.GetToUserList(ctx), "|")
	fromUser := c.base.GetCorpID()
	timeStamp := time.Now().Unix()
	var msg string
	switch typ {
	case ButtonClickCBType:
		respT := `<xml><ToUserName><![CDATA[%v]]></ToUserName><FromUserName><![CDATA[%v]]></FromUserName><CreateTime>%d</CreateTime><MsgType><![CDATA[update_button]]></MsgType><Button><ReplaceName><![CDATA[%v]]></ReplaceName></Button></xml>`
		msg = fmt.Sprintf(respT, toUser, fromUser, timeStamp, string(showinfo))
	}
	respStr, err := c.MakeResponseData(ctx, msg, timeStamp)
	if err != nil {
		return nil, err
	}
	return []byte(respStr), nil
}
