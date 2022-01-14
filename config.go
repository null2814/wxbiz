package qywx

import "sync"

type BaseInfo struct {
	CorpId         string   `yaml:"corpid" json:"corpid"`
	CorpSecret     string   `yaml:"corpsecret" json:"corpsecret"`
	AgentId        int32    `yaml:"agentId" json:"agentId"`
	EncodingAeskey string   `yaml:"encodingAeskey" json:"encodingAeskey"`
	ToUserList     []string `yaml:"userList" json:"userList"`
	Token          string   `yaml:"token" json:"token"`
	ReceiverId     string   `yaml:"receiverId" json:"receiverId"`
	l              *sync.RWMutex
}

func (b *BaseInfo) SetBaseInfo(new *BaseInfo) {
	b.l.Lock()
	defer b.l.Unlock()
	b.CorpId = new.CorpId
	b.CorpSecret = new.CorpSecret
	b.AgentId = new.AgentId
	b.EncodingAeskey = new.EncodingAeskey
	b.ToUserList = new.ToUserList
	b.Token = new.Token
	b.ReceiverId = new.ReceiverId
}

func (b *BaseInfo) GetAgentId() int32 {
	b.l.RLock()
	defer b.l.Unlock()
	id := b.AgentId
	return id
}

func (b *BaseInfo) GetCorpID() string {
	b.l.RLock()
	defer b.l.Unlock()
	id := b.CorpId
	return id
}

type MessageInfo struct {
	// ToUserList []string
	// Token      string
	// ReceiverId string
	// l          *sync.RWMutex
}
