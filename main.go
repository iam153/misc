package main

import (
	"bloom/bloom"
	"fmt"
)

func main() {
	// n限制128
	f := bloom.NewBloomFilter(12, 2048)

	keys := [][]byte{[]byte("abc"), []byte("bcd"), []byte("124213")}

	for _, k := range keys {
		err := f.Insert(k)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	}

	for _, k := range keys {
		exist, err := f.Exist(k)
		fmt.Printf("exist:%v, err:%v\n", exist, err)
	}

	e, _ := f.Exist([]byte("ccc"))
	fmt.Println(e)
	f.Exist([]byte("cccccc"))
	fmt.Println(e)
	f.Exist([]byte("ccccccccc"))
	fmt.Println(e)

	// fmt.Println(byte(1 << int64(2)))
}
