package hiorder

import "github.com/hootuu/hyle/hfsm"

type Status = hfsm.State

const (
	_ Status = iota
	Draft
	Initial   //以存储
	Consensus //达成一致.已完成支付等环节
	Executing //执行中
	Completed //已完成
	Canceled  //已取消
	Timeout   //已超时
)

const (
	_ hfsm.Event = iota
	SubmitEvent
	ConsenseEvent
	ExecuteEvent
	CompleteEvent
	TimeoutEvent
	CancelEvent
)

type ExStatus = hfsm.State
