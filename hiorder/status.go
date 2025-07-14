package hiorder

import "github.com/hootuu/hyle/hfsm"

type Status = hfsm.State

const (
	_         Status = iota
	Draft            //草稿状态,未持久化
	Initial          //以存储
	Consensus        //达成一致.已完成支付等环节
	Executing        //执行中
	Completed        //已完成
	Canceled         //已取消
	Timeout          //已超时
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

// FOR Ex
const (
	ExOnDraft     ExStatus = Draft * 1000
	ExOnInitial   ExStatus = Initial * 1000
	ExOnConsensus ExStatus = Consensus * 1000
	ExOnExecuting ExStatus = Executing * 1000
	ExOnCompleted ExStatus = Completed * 1000
	ExOnCanceled  ExStatus = Canceled * 1000
	ExOnTimeout   ExStatus = Timeout * 1000
)
