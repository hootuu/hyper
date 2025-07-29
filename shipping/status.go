package shipping

import "github.com/hootuu/hyle/hfsm"

// 物流状态常量定义（兼容快递100/京东等常见状态）
const (
	_ Status = iota
	Initialized
	Submitted
	Completed
	Failed
	Canceled
)

const (
	_ hfsm.Event = iota
	SubmitEvent
	CancelEvent
	FailEvent
	CompleteEvent
)
