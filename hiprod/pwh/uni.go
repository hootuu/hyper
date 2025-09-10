package pwh

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
)

const (
	uniPwhID = "HYPER_PWH_UNI"
)

var gUniPwhID ID

func UniPwhID() ID {
	return gUniPwhID
}

func initUni() error {
	uniLink := collar.Build("PWH", uniPwhID).ToID()
	pwhM, err := hdb.Get[PhysicalWhM](hyperplt.DB(), "link = ?", uniLink)
	if err != nil {
		return err
	}
	if pwhM == nil {
		pwhM = &PhysicalWhM{
			Template: hdb.Template{
				Ctrl: ctrl.MustNewCtrl().MustSet(0, true),
				Tag:  hjson.MustToBytes(tag.NewTag("UNI")),
				Meta: hjson.MustToBytes(dict.NewDict().Set("id", uniPwhID)),
			},
			ID:   nextID(),
			Link: uniLink,
			Memo: "the uni pwh",
		}
		err := hdb.GetOrCreate[PhysicalWhM](hyperplt.DB(), pwhM, "link = ?", uniLink)
		if err != nil {
			return err
		}
	}
	gUniPwhID = pwhM.ID
	return nil
}
