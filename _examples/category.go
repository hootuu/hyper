package main

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyper/category"
)

func main() {
	helix.AfterStartup(func() {
		_, _ = category.Backend().Add(category.Root, "TEST", "ticon")
		//cat, err := category.NewCategory("background")
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//for i := 0; i < 0; i++ {
		//	if id, err := cat.Add(category.Root, "NAME-"+cast.ToString(i)+"-"+cast.ToString(time.Now().Unix()), "A"); err != nil {
		//		fmt.Println(err)
		//		return
		//	} else {
		//		//fmt.Println(id)
		//		for j := 0; j < 3; j++ {
		//			if jID, err := cat.Add(id, "NAME-"+cast.ToString(i)+"-"+cast.ToString(j)+"-"+cast.ToString(time.Now().Unix()), "B"); err != nil {
		//				fmt.Println(err)
		//				return
		//			} else {
		//				for n := 0; n < 3; n++ {
		//					if _, err := cat.Add(jID, "NAME-"+cast.ToString(i)+"-"+cast.ToString(j)+"-"+cast.ToString(n)+"-"+cast.ToString(time.Now().Unix()), "B"); err != nil {
		//						fmt.Println(err)
		//						return
		//					}
		//				}
		//			}
		//		}
		//	}
		//}
		//arr, err := cat.Get(category.Root, 1)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//fmt.Println(hjson.MustToString(arr))
	})
	helix.Startup()
}
