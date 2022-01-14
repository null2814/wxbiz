package qywx

const QWURL = "https://qyapi.weixin.qq.com"

type StatusDescription struct {
	Code int32
	Desc string
}

var (
	StatusReject = StatusDescription{0, "拒绝"}
	StatusAccept = StatusDescription{1, "接收"}
)
