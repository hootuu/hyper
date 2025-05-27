package main

import (
	"fmt"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyper/address/maps"
	"github.com/hootuu/hyper/address/maps/amap"
	"time"
)

func main() {
	p := maps.NewAmapProvider("213c33139661d94b70f97c4646744fab")
	s := time.Now()
	err := p.RegionSync(func(r *maps.Region) error {
		fmt.Println(hjson.MustToString(r))
		//if r.Name == "太湖县" {
		//	os.Exit(-1)
		//}
		return nil
	})
	fmt.Println(err)
	fmt.Println("elapse ms: ", time.Now().Sub(s).Milliseconds()/1000)
}

func main2() {
	amapCli := amap.NewClient("213c33139661d94b70f97c4646744fab")
	data, err := amapCli.District("100000", 1, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hjson.MustToString(data))
}
