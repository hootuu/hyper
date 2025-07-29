package payment

import "github.com/hootuu/hyle/hfsm"

type Status = hfsm.State

const (
	_ Status = iota
	Initialized
	Prepared
	Executing
	Completed
	Timeout
	Canceled
)

const (
	_ hfsm.Event = iota
	PrepareEvent
	ExecuteEvent
	TimeoutEvent
	CancelEvent
	CompleteEvent
)

type JobStatus = hfsm.State

const (
	_ JobStatus = iota
	JobInitialized
	JobPrepared
	JobCompleted
	JobTimeout
	JobCanceled
)

const (
	_ hfsm.Event = iota
	JobPrepareEvent
	JobTimeoutEvent
	JobCancelEvent
	JobCompleteEvent
)
