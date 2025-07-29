package payment

import (
	"fmt"
	"github.com/hootuu/hyle/crypto/hmd5"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func BuildJobID(paymentID ID, seq int) JobID {
	return hmd5.MD5(fmt.Sprintf("%d_%d", paymentID, seq))
}

var gJobExecutorMap = map[Channel]JobExecutor{}

func RegisterJobExecutor(jobExecutor JobExecutor) {
	channel := jobExecutor.GetChannel()
	_, ok := gJobExecutorMap[channel]
	if ok {
		hlog.Fix("hpay.RegisterJobExecutor: job executor is already registered",
			zap.String("channel", channel))
	}
	gJobExecutorMap[jobExecutor.GetChannel()] = jobExecutor
}

func MustGetJobExecutor(channel Channel) (JobExecutor, error) {
	executor, ok := gJobExecutorMap[channel]
	if !ok {
		return nil, fmt.Errorf("job executor not registered: %s", channel)
	}
	return executor, nil
}
