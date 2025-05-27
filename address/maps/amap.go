package maps

import (
	"fmt"
	"github.com/hootuu/hyle/hlog"
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

func (p *AmapProvider) RegionSync(call func(r *Region) error) error {
	lstCallAmap := time.Now()
	page := 1
	for true {
		callAmapInterval := time.Now().Sub(lstCallAmap)
		if callAmapInterval < 1500*time.Millisecond {
			time.Sleep(1500*time.Millisecond - callAmapInterval)
		}
		lstCallAmap = time.Now()
		arr, err := p.cli.District(amap.China, page, 1)
		if err != nil {
			return err
		}
		if len(arr) == 0 {
			break
		}
		for _, country := range arr {
			if err := call(p.convertDistrict(country)); err != nil {
				return err
			}
			if len(country.Districts) == 0 {
				hlog.Err("hyper.addr.amap.RegionSync: len(country.Districts) == 0")
				return fmt.Errorf("no any districts with country")
			}
			for _, province := range country.Districts {
				if err := call(p.convertDistrict(province)); err != nil {
					return err
				}
				callAmapInterval := time.Now().Sub(lstCallAmap)
				if callAmapInterval < 1500*time.Millisecond {
					time.Sleep(1500*time.Millisecond - callAmapInterval)
				}
				lstCallAmap = time.Now()
				err := p.doLoadRegion(province.Adcode, call)
				if err != nil {
					return err
				}
			}
		}
		page += 1
	}
	return nil
}

func (p *AmapProvider) doLoadRegion(adCode string, call func(r *Region) error) error {
	page := 1
	for true {
		arr, err := p.cli.District(adCode, page, 3)
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
				if err := call(p.convertDistrict(one)); err != nil {
					return err
				}
				if len(one.Districts) == 0 {
					continue
				}
				for _, two := range one.Districts {
					if err := call(p.convertDistrict(two)); err != nil {
						return err
					}
					if len(two.Districts) == 0 {
						continue
					}
					for _, three := range two.Districts {
						if err := call(p.convertDistrict(three)); err != nil {
							return err
						}
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
		Code:  data.Adcode,
		Name:  data.Name,
		Level: 0,
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
