package payment

import "github.com/hootuu/hyle/hfsm"

type Status hfsm.State

const (
	_ Status = iota
	Initialized
	Prepared
	Completed
	Timeout
	Canceled
)

type JobStatus hfsm.State

const (
	_ JobStatus = iota
	JobInitialized
	JobPrepared
	JobCompleted
	JobTimeout
	JobCanceled
)
