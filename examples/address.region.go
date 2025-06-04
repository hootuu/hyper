package main

import (
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyper/address"
)

func main() {
	helix.AfterStartup(func() {
		arr, err := address.RegionChildren(0)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(hjson.MustToString(arr))
	})
	helix.Startup()
}
