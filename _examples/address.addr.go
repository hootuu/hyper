package main

import (
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyper/address"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		for i := 0; i < 110; i++ {
			addr, err := address.AddAddress(&address.Address{
				Owner:   "abcd",
				Region:  8001001001001001,
				Address: "ADDR-" + cast.ToString(time.Now().UnixMilli()),
				Contact: address.Contact{
					Name: "张继国3",
					Mobi: "18988998899",
				},
				Default:  true,
				Location: address.Location{},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(hjson.MustToString(addr))
			err = address.UseAddress(addr.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = address.MutAddress(&address.Address{
				ID:      addr.ID,
				Owner:   addr.Owner,
				Region:  addr.Region,
				Address: "ADDR-MUT-" + cast.ToString(time.Now().UnixMilli()),
				Contact: address.Contact{
					Name: "张继国3-MUT",
					Mobi: "18988998899",
				},
				Default: false,
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			err = address.DelAddress(addr.ID)

			if err != nil {
				fmt.Println(err)
				return
			}
		}

	})
	helix.Startup()
}
