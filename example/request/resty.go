package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aivencs/box/pkg/request"
)

func main() {
	ctx := context.WithValue(context.Background(), "trace", "r001")
	ctx = context.WithValue(ctx, "label", "request")
	for i := 0; i < 2; i++ {
		err := request.InitRequest(ctx, "resty", request.Option{})
		if err != nil {
			log.Fatal(err)
		}
		res, err := request.Get(ctx, request.Param{
			Link:             "https://wind.in-en.com/windcorp/",
			Timeout:          2,
			EnableSkipVerify: true,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("status: ", res.StatusCode, "err: ", err)
	}
}
