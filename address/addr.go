package address

import (
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func verifyAddr(m *AddrM) error {
	if m.Region == 0 {
		return fmt.Errorf("require region: %d", m.Region)
	}
	if m.Owner == "" {
		return errors.New("require owner")
	}
	if m.Addr == "" {
		return errors.New("require addr")
	}
	if m.CName == "" {
		return errors.New("require contact name")
	}
	if m.CMobi == "" {
		return errors.New("require contact mobi")
	}
	return nil
}

func addAddr(m *AddrM) (*AddrM, error) {
	if err := verifyAddr(m); err != nil {
		return nil, err
	}
	total, err := hdb.Count[AddrM](zplt.HelixPgDB().PG(), "owner = ?", m.Owner)
	if err != nil {
		hlog.Err("hyper.address.addAddr: Count", zap.Error(err))
		return nil, err
	}
	addrMax := hcfg.GetInt64("hyper.address.max", 100)
	if total > addrMax {
		return nil, fmt.Errorf("the maximum number of storage addresses is %d: %d", addrMax, total)
	}

	//regionM, err := hdb.MustGet[RegionM](zplt.HelixPgDB().PG(), "id = ?", m.Region)
	//if err != nil {
	//	return nil, err
	//}
	//m.FullAddr = regionM.Address + m.Addr
	m.FullAddr = m.Addr
	m.ID = idx.New()

	err = hdb.Tx(zplt.HelixPgDB().PG(), func(tx *gorm.DB) error {
		if m.Default {
			err := hdb.Update[AddrM](tx, map[string]any{
				"is_default": false,
			}, "owner = ?", m.Owner)
			if err != nil {
				return err
			}
		}
		err := hdb.Create[AddrM](tx, m)
		if err != nil {
			hlog.Err("hyper.address.addAddr: Create", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func mutAddr(mutM *AddrM) error {
	if err := verifyAddr(mutM); err != nil {
		return err
	}
	dbM, err := hdb.Get[AddrM](zplt.HelixPgDB().PG(), "id = ?", mutM.ID)
	if err != nil {
		hlog.Err("hyper.address.mutAddr: Get", zap.Error(err))
		return err
	}
	if dbM == nil {
		return fmt.Errorf("no such addr: %s", mutM.ID)
	}
	mut := make(map[string]any)
	if dbM.Region != mutM.Region {
		mut["region"] = mutM.Region

		regionM, err := hdb.MustGet[RegionM](zplt.HelixPgDB().PG(), "id = ?", mutM.Region)
		if err != nil {
			return err
		}
		mut["full_addr"] = regionM.Address + mutM.Addr
	}
	if dbM.Addr != mutM.Addr {
		mut["addr"] = mutM.Addr
	}
	if dbM.CName != mutM.CName {
		mut["cname"] = mutM.CName
	}
	if dbM.CMobi != mutM.CMobi {
		mut["cmobi"] = mutM.CMobi
	}
	if dbM.LocX != mutM.LocX {
		mut["loc_x"] = mutM.LocX
	}
	if dbM.LocY != mutM.LocY {
		mut["loc_y"] = mutM.LocY
	}
	if dbM.Default != mutM.Default {
		mut["is_default"] = mutM.Default
	}
	err = hdb.Tx(zplt.HelixPgDB().PG(), func(tx *gorm.DB) error {
		if mutM.Default {
			err := hdb.Update[AddrM](tx, map[string]any{
				"is_default": false,
			}, "owner = ?", mutM.Owner)
			if err != nil {
				return err
			}
		}
		err := hdb.Update[AddrM](tx, mut, "id = ?", mutM.ID)
		if err != nil {
			hlog.Err("hyper.address.addAddr: Update", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func useAddr(id string) error {
	mut := map[string]any{
		"usage": gorm.Expr("`usage` + ?", 1),
	}
	err := hdb.Update[AddrM](zplt.HelixPgDB().PG(), mut, "id = ?", id)
	if err != nil {
		hlog.Err("hyper.address.useAddr", zap.Error(err))
		return err
	}
	return nil
}

func listAddrByOwner(owner sattva.Identification) ([]*AddrM, error) {
	arrM, err := hdb.Find[AddrM](func() *gorm.DB {
		return zplt.HelixPgDB().PG().Where("owner = ?", owner)
	})
	if err != nil {
		return nil, err
	}
	if len(arrM) == 0 {
		return []*AddrM{}, nil
	}
	return arrM, nil
}

func delAddr(id string) error {
	err := hdb.Delete[AddrM](zplt.HelixPgDB().PG(), "id = ?", id)
	if err != nil {
		hlog.Err("hyper.address.delAddr", zap.Error(err))
		return err
	}
	return nil
}
