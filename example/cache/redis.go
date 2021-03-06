package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aivencs/box/pkg/cache"
)

func main() {
	ctx := context.WithValue(context.Background(), "trace", "ctx-cache-001")
	opt := cache.Option{
		Host:     "localhost:6379",
		Auth:     true,
		Username: "",
		Password: "password",
		Database: "",
		Table:    "",
		DB:       1,
	}
	err := cache.InitCache(ctx, cache.REDIS, opt)
	if err != nil {
		log.Fatal(err)
	}
	payload := "19619c9e08f0ed4cc147e211efa8c3f0"
	r, err := cache.SetEx(ctx, payload, 1, 20)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r, err) // output: OK nil
	ov := cache.Overdue(ctx, payload)
	fmt.Println(ov)                             // output: true
	fmt.Println(cache.Set(ctx, payload, "105")) // output: OK nil
	val, err := cache.Get(ctx, payload)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(val.([]uint8)), err) // output: 105 <nil>
}
