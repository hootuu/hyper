package vwh

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
	uniVwhID = "HYPER_VWH_UNI"
)

var gUniVwhID ID

func UniVwhID() ID {
	return gUniVwhID
}

func initUni() error {
	uniLink := collar.Build("VWH", uniVwhID).ToID()
	vwhM, err := hdb.Get[VirtualWhM](hyperplt.DB(), "link = ?", uniLink)
	if err != nil {
		return err
	}
	if vwhM == nil {
		vwhM = &VirtualWhM{
			Template: hdb.Template{
				Ctrl: ctrl.MustNewCtrl().MustSet(0, true),
				Tag:  hjson.MustToBytes(tag.NewTag("UNI")),
				Meta: hjson.MustToBytes(dict.NewDict().Set("id", uniVwhID)),
			},
			ID:   nextID(),
			Link: uniLink,
			Memo: "the uni vwh",
		}
		err := hdb.GetOrCreate[VirtualWhM](hyperplt.DB(), vwhM, "link = ?", uniLink)
		if err != nil {
			return err
		}
	}
	gUniVwhID = vwhM.ID
	return nil
}
