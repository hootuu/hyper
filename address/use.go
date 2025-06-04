package address

import (
	"context"
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func RegionChildren(id RegionID, deep ...int) ([]*Region, error) {
	d := 1
	if len(deep) > 0 {
		d = deep[0]
	}
	return getRegion(id, d)
}

func AddAddress(addr *Address) (*Address, error) {
	addrM, err := addAddr(&AddrM{
		Owner:   addr.Owner,
		CName:   addr.Contact.Name,
		CMobi:   addr.Contact.Mobi,
		Default: addr.Default,
		Region:  addr.Region,
		Addr:    addr.Address,
		LocX:    addr.Location.Lon,
		LocY:    addr.Location.Lat,
	})
	if err != nil {
		return nil, err
	}
	return addrM.ToAddress(), nil
}

func MutAddress(addr *Address) error {
	return mutAddr(&AddrM{
		ID:      addr.ID,
		Owner:   addr.Owner,
		CName:   addr.Contact.Name,
		CMobi:   addr.Contact.Mobi,
		Default: addr.Default,
		Region:  addr.Region,
		Addr:    addr.Address,
		LocX:    addr.Location.Lon,
		LocY:    addr.Location.Lat,
	})
}

func UseAddress(id string) error {
	return useAddr(id)
}

func MyAddress(owner sattva.Identification) ([]*Address, error) {
	arrM, err := listAddrByOwner(owner)
	if err != nil {
		return nil, err
	}
	var arr []*Address
	for _, m := range arrM {
		arr = append(arr, m.ToAddress())
	}
	return arr, nil
}

func DelAddress(id string) error {
	return delAddr(id)
}

func init() {
	helix.Use(helix.BuildHelix("hyper_address", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&RegionM{},
			&AddrM{},
		)
		if err != nil {
			return nil, err
		}
		err = regionInit()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
