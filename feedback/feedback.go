package feedback

import (
	"errors"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func addFeedback(fbM *FbM) (*FbM, error) {
	if fbM.Person == "" {
		return nil, errors.New("require Person")
	}
	if fbM.Title == "" {
		return nil, errors.New("require Title")
	}
	fbM.ID = idx.New()
	err := hpg.Create[FbM](zplt.HelixPgDB().PG(), fbM)
	if err != nil {
		hlog.Err("hyper.feedback.addFeedback: Create", zap.Error(err))
		return nil, err
	}
	return fbM, nil
}
