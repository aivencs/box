package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aivencs/box/pkg/config"
)

var Conf = BindConf{}

type BindConf struct {
	Header ConfHeader `json:"header"`
	Body   ConfBody   `json:"body"`
}

type ConfHeader struct {
	Application string `json:"application"`
	Shutdown    int    `json:"shutdown"`
}

type ConfBody struct {
	Request Request `json:"request"`
}

type Request struct {
	Request BindConfRequest `json:"request"`
	Render  BindConfRender  `json:"render"`
}

type BindConfRequest struct {
	Timeout int           `json:"timeout"`
	Proxy   BindConfProxy `json:"proxy"`
}

type BindConfRender struct {
	Timeout int           `json:"timeout"`
	Proxy   BindConfProxy `json:"proxy"`
	PoolCap int           `json:"pool_cap"`
	Monitor bool          `json:"monitor"`
	Host    string        `json:"host"`
}

type BindConfProxy struct {
	Auth     bool   `json:"auth"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func main() {
	ctx := context.Background()
	// 初始化配置对象
	err := config.InitConf(ctx, config.Consul, config.Option{
		Application: "delnic-service-request",
		Env:         "dev",
		Auth:        true,
		Username:    "piker019",
		Password:    "cro00Too01",
		Host:        "http://consul.wqzcir.com",
		Type:        "yaml",
		Bind:        &Conf,
		Update:      true,
		Interval:    10,
	})
	if err != nil {
		log.Fatal(err)
	}
	// 使用方法
	for i := 0; i < 1000; i++ {
		fmt.Println("bind-", i, ": ", Conf) // 直接访问
		// 期间可以修改配置中的内容，以观察自动定时更新是否生效
		time.Sleep(time.Second * 3)
	}
}
