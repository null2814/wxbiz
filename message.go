package qywx

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type (
	SMType string
	CBType string
)

const ButtonClickCBType CBType = "ButtonClick"

type BaseMsg struct {
	SMMsgType SMType
	CBMsgType CBType
	MsgBody   interface{}
}

type (
	ButtonClickCallbackRequestBody struct {
		ToUserName    string                                       `xml:"ToUserName"`   // 企业微信CorpID
		FromUserName  string                                       `xml:"FromUserName"` // 成员UserID
		CreateTime    int                                          `xml:"CreateTime"`   // 消息创建时间（整型）
		MsgType       string                                       `xml:"MsgType"`      // 消息类型 event
		Event         string                                       `xml:"Event"`        // 事件类型 template_card_event
		EventKey      string                                       `xml:"EventKey"`     // 按钮btn:key值
		TaskId        string                                       `xml:"TaskId"`       // task_id
		CardType      string                                       `xml:"CardType"`     // 通用模板卡片的类型
		ResponseCode  string                                       `xml:"ResponseCode"` // 用于调用更新卡片接口的ResponseCode
		AgentID       int                                          `xml:"AgentID"`      // 企业应用的id，整型
		SelectedItems []ButtonClickCallbackRequestBodySelectedItem `xml:"SelectedItems"`
	}
	ButtonClickCallbackRequestBodySelectedItem struct {
		QuestionKey string                                               `xml:"QuestionKey"` // 问题的key值
		OptionIds   []ButtonClickCallbackRequestBodySelectedItemOptionId `xml:"OptionIds"`   // 对应问题的选项列表
	}
	ButtonClickCallbackRequestBodySelectedItemOptionId string
)

type (
	UpdateTemplateCardRequest struct {
		AtAll        int32                    `json:"atall"`
		AgentId      int32                    `json:"agentid"`
		ResponseCode string                   `json:"response_code"`
		Button       UpdateTemplateCardButton `json:"button"`
	}
	UpdateTemplateCardButton struct {
		ReplaceName string `json:"replace_name"`
	}
)

type UpdateTemplateCardResponse struct {
	ErrCode         int32    `json:"errcode"`
	ErrMsg          string   `json:"errmsg"`
	InvalidUserIds  []string `json:"invalid_userids"`
	InvalidPartyIds []string `json:"invalid_partyids"`
	InvalidTagIds   []string `json:"invalid_tagids"`
}

func (r *UpdateTemplateCardRequest) GetUpdateTemplateCardResp(ctx context.Context, c *QYWXClient, access *AccessToken) (*UpdateTemplateCardResponse, error) {
	var url strings.Builder
	url.WriteString(QWURL)
	url.WriteString("/cgi-bin/message/update_template_card")
	url.WriteString(fmt.Sprintf("?access_toke=%s", access.Token))
	err := c.GetAuthorization(ctx, access)
	if nil != err {
		return nil, err
	}
	utxr := new(UpdateTemplateCardResponse)
	req, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	resp, err := NetPOSTJson(ctx, url.String(), req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, utxr)
	return utxr, err
}
