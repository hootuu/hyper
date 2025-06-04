package maps

import (
	"fmt"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyper/address/maps/amap"
	"go.uber.org/zap"
	"time"
)

type AmapResponse struct {
}

type AmapProvider struct {
	cli *amap.Client
}

func NewAmapProvider(key string) *AmapProvider {
	cli := amap.NewClient(key)
	return &AmapProvider{cli: cli}
}

func (p *AmapProvider) RegionSync(call func(parentId htree.ID, r *Region) (htree.ID, error)) error {
	hsys.Info("Region Syncing From AMAP ...")
	defer hlog.Elapse("hyper.address.RegionSync")()
	lstCallAmap := time.Now()
	page := 1
	for true {
		callAmapInterval := time.Now().Sub(lstCallAmap)
		if callAmapInterval < 1500*time.Millisecond {
			time.Sleep(1500*time.Millisecond - callAmapInterval)
		}
		lstCallAmap = time.Now()
		arr, err := p.cli.District(amap.China, page, 2)
		if err != nil {
			return err
		}
		if len(arr) == 0 {
			break
		}
		for _, country := range arr {
			countryId, err := call(0, p.convertDistrict(country))
			if err != nil {
				return err
			}
			if len(country.Districts) == 0 {
				hlog.Err("hyper.addr.amap.RegionSync: len(country.Districts) == 0")
				return fmt.Errorf("no any districts with country")
			}
			for _, province := range country.Districts {
				hsys.Info("Amap Sync Province ", province.Name, " ...")
				province.Address = province.Name
				provinceId, err := call(countryId, p.convertDistrict(province))
				if err != nil {
					return err
				}
				if len(province.Districts) == 0 {
					continue
				}
				for _, city := range province.Districts {
					city.Address = province.Address + city.Name
					cityId, err := call(provinceId, p.convertDistrict(city))
					if err != nil {
						return err
					}
					callAmapInterval := time.Now().Sub(lstCallAmap)
					if callAmapInterval < 1100*time.Millisecond {
						time.Sleep(1100*time.Millisecond - callAmapInterval)
					}
					lstCallAmap = time.Now()
					err = p.doLoadRegion(cityId, city, call)
					if err != nil {
						return err
					}
				}
				hsys.Success("Amap Sync Province ", province.Name, " [OK]")
			}
		}
		page += 1
	}
	hsys.Info("Region Sync From AMAP [OK]")
	return nil
}

func (p *AmapProvider) doLoadRegion(id htree.ID, city *amap.District, call func(parentId htree.ID, r *Region) (htree.ID, error)) error {
	page := 1
	for true {
		arr, err := p.cli.District(city.Adcode, page, 2)
		if err != nil {
			hlog.Err("hyper.addr.amap.doLoadRegion: ", zap.Error(err))
			return err
		}
		if len(arr) == 0 {
			break
		}
		for _, item := range arr {
			if len(item.Districts) == 0 {
				continue
			}
			for _, one := range item.Districts {
				one.Address = city.Address + one.Name
				oneId, err := call(id, p.convertDistrict(one))
				if err != nil {
					return err
				}
				if len(one.Districts) == 0 {
					continue
				}
				for _, two := range one.Districts {
					two.Address = one.Address + two.Name
					_, err := call(oneId, p.convertDistrict(two))
					if err != nil {
						return err
					}
				}
			}
		}
		page += 1
	}
	return nil
}

func (p *AmapProvider) convertDistrict(data *amap.District) *Region {
	r := &Region{
		Map:     AMAP,
		Code:    data.Adcode,
		Name:    data.Name,
		Level:   0,
		Address: data.Address,
	}
	switch data.Level {
	case "country":
		r.Level = Country
	case "province":
		r.Level = Province
	case "city":
		r.Level = City
	case "district":
		r.Level = District
	case "street":
		r.Level = Street
	}

	return r
}
