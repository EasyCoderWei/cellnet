package cellnet

import "fmt"

type RegisterMessageContext struct {
	*MessageMeta
}

type CallbackHandler struct {
	userCallback func(*SessionEvent)
}

func (self *CallbackHandler) Call(ev *SessionEvent) {

	self.userCallback(ev)

}

func NewCallbackHandler(userCallback func(*SessionEvent)) EventHandler {
	return &CallbackHandler{
		userCallback: userCallback,
	}
}

// 注册消息处理回调
// DispatcherHandler -> socket.DecodePacketHandler -> socket.CallbackHandler
func RegisterMessage(p Peer, msgName string, userCallback func(*SessionEvent)) *RegisterMessageContext {

	meta := MessageMetaByName(msgName)

	if meta == nil {
		panic(fmt.Sprintf("message register failed, %s", msgName))
	}

	p.AddHandler(int(meta.ID), HandlerLink(NewQueuePostHandler(p.Queue(), HandlerLink(NewCallbackHandler(userCallback)))))

	return &RegisterMessageContext{MessageMeta: meta}
}

// 注册消息处理的一系列Handler
// DispatcherHandler -> socket.DecodePacketHandler -> ...
func RegisterHandler(p Peer, msgName string, handlers ...EventHandler) *RegisterMessageContext {

	meta := MessageMetaByName(msgName)

	if meta == nil {
		panic(fmt.Sprintf("message register failed, %s", msgName))
	}

	p.AddHandler(int(meta.ID), HandlerLink(NewQueuePostHandler(p.Queue(), handlers)))

	return &RegisterMessageContext{MessageMeta: meta}
}
