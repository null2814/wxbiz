package qywx

const (
	TextType SMType = "SendMessageTextType"
	CardType SMType = "SendMessageTemplateCardType"
)

type GetAuthorizationApiResponse struct {
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

type GetTicketApiResponse struct {
	ErrCode   int32  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int32  `json:"expires_in"`
}

type SendMessageBaseType struct {
	ToUser                 string `json:"touser"`                   // 指定接收消息的成员
	ToParty                string `json:"toparty"`                  // 指定接收消息的部门
	ToTag                  string `json:"totag"`                    // 指定接收消息的标签
	MsgType                string `json:"msgtype"`                  // 消息类型
	AgentId                int32  `json:"agentid"`                  // 企业应用的id
	EnableIdTrans          int32  `json:"enable_id_trans"`          // 表示是否开启id转译
	EnableDuplicateCheck   int32  `json:"enable_duplicate_check"`   // 表示是否开启重复消息检查
	DuplicateCheckInterval int32  `json:"duplicate_check_interval"` // 表示是否重复消息检查的时间间隔
}

/*
{
    "template_card" : {
        "card_type" : "button_interaction",
        "main_title" : {
            "title" : "ES新建应用审批",
            "desc" : "ES新建应用审批相关信息"
        },
        "task_id": "task_id",

		// 暂不需要:
        "button_selection": {
            "question_key": "btn_question_key1",
            "title": "企业微信评分",
            "option_list": [
                {
                    "id": "btn_selection_id1",
                    "text": "100分"
                },
                {
                    "id": "btn_selection_id2",
                    "text": "101分"
                }
            ],
            "selected_id": "btn_selection_id1"
        },

        "button_list": [
            {
                "text": "按钮1",
                "style": 1,
                "key": "button_key_1"
            },
            {
                "text": "按钮2",
                "style": 2,
                "key": "button_key_2"
            }
        ]
    },
}
*/
type (
	SendMessageTemplateCardType struct {
		SendMessageBaseType
		TemplateCard SendMessageTemplateCardContent `json:"template_card"`
	}
	SendMessageTemplateCardContent struct {
		CardType  string `json:"card_type"` // 模板卡片类型,button_interaction
		MainTitle struct {
			Title string `json:"title"`
			// Desc  string `json:"desc"` // 更多操作界面的描述
		} `json:"main_title"` // 一级标题
		SubTitleText string                                 `json:"sub_title_text"`
		TaskId       string                                 `json:"task_id"`     // 任务id，不能重复，最长128字节
		ButtonList   []SendMessageTemplateCardContentButton `json:"button_list"` // 按钮列表，列表长度不超过6
	}
	SendMessageTemplateCardContentButton struct {
		Text  string `json:"text"`  // 按钮文案
		Style int32  `json:"style"` // 按钮样式
		Key   string `json:"key"`   // 按钮key值，最长支持1024字节，不可重复
	}
)

type SendMessageTextType struct {
	SendMessageBaseType
	Text struct {
		Content string `json:"content"` // 消息内容，最长不超过2048个字节
	} `json:"text"`
	Safe int32 `json:"safe"` //表示是否是保密消息
}

type SendMessageResponse struct {
	ErrCode      int32  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
	Msgid        string `json:"msgid"`
	ResponseCode string `json:"response_code"`
}

type CallbackCheckRequest struct {
	MsgSignature string `json:"msg_signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Echostr      string `json:"echostr"`
}

type CallbackRequest struct {
	MsgSignature string
	Timestamp    string
	Nonce        string
	Body         []byte
}

type CallbackRequestBody struct {
	ToUser     string `xml:"ToUserName"`
	ToAgentID  string `xml:"AgentID"`
	MsgEncrypt string `xml:"Encrypt"`
}
