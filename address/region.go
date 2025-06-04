package address

import (
	"fmt"
	"github.com/hootuu/helix/components/honce"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyper/address/maps"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegionID = htree.ID

func getRegion(id RegionID, deep int) ([]*Region, error) {
	if deep < 1 || deep > 5 {
		return nil, fmt.Errorf("invalid deep: %d", deep)
	}
	minId, maxId, base, err := gRegionTree.Factory().DirectChildren(id)
	if err != nil {
		return nil, err
	}
	arrM, err := hpg.Find[RegionM](func() *gorm.DB {
		return zplt.HelixPgDB().PG().
			Where("id % ? = 0 AND id BETWEEN ? AND ?", base, minId, maxId)
	})
	if err != nil {
		return nil, err
	}

	var arr []*Region
	for _, rM := range arrM {
		r := &Region{
			ID:       rM.ID,
			Map:      rM.Map,
			Code:     rM.Code,
			Name:     rM.Name,
			Children: nil,
		}
		arr = append(arr, r)
		if deep > 1 {
			r.Children, err = getRegion(r.ID, deep-1)
			if err != nil {
				return nil, err
			}
		}
	}

	return arr, nil
}

var gRegionTree *htree.Tree

func regionSave(parentId htree.ID, r *maps.Region) (htree.ID, error) {
	if parentId == 0 {
		parentId = gRegionTree.Root()
	}
	exist, err := hpg.Exist[RegionM](zplt.HelixPgDB().PG(),
		"map = ? AND code = ? AND name = ?",
		r.Map, r.Code, r.Name,
	)
	if err != nil {
		hlog.Err("hyper.address.regionSave: Exist", zap.Error(err))
		return 0, err
	}
	if exist {
		return 0, nil
	}
	var newId htree.ID
	err = gRegionTree.Next(parentId, func(id htree.ID) error {
		newId = id
		return nil
	})
	if err != nil {
		return 0, err
	}
	err = hpg.Create[RegionM](zplt.HelixPgDB().PG(), &RegionM{
		ID:      newId,
		Map:     r.Map,
		Code:    r.Code,
		Name:    r.Name,
		Address: r.Address,
	})
	if err != nil {
		hlog.Err("hyper.address.regionSave: Create", zap.Error(err))
		return 0, err
	}
	return newId, nil
}

func regionInit() error {
	var err error
	gRegionTree, err = htree.NewTree("hyper_addr_region", 8, []uint{3, 3, 3, 3, 3})
	if err != nil {
		return err
	}
	err = honce.Do("hyper.address.region.init.v1", func() error {
		helix.AfterStartup(func() {
			amapKey, err := hcfg.MustGetString("hyper.address.amap.key")
			if err != nil {
				hsys.Exit(err)
				return
			}
			p := maps.NewAmapProvider(amapKey)
			err = p.RegionSync(regionSave)
			if err != nil {
				hsys.Exit(err)
				return
			}
		})
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
