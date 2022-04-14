package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aivencs/box/pkg/logger"
)

func main() {
	/* example for error code */
	ber := logger.NewErr(logger.ErrOption{Code: logger.PVERROR, Label: "数据长度不足"})
	fmt.Println("ber: ", ber)

	/* example for zap logger */
	ctx := context.WithValue(context.Background(), "trace", "t001")
	// 初始化日志对象
	err := logger.InitLogger(ctx, logger.ZAP, logger.Option{
		// Application: "zap-log",
		Env:    "dev",
		Label:  "detail",
		Encode: logger.JSON})
	if err != nil {
		log.Fatal(err)
	}
	// 1
	ctx = context.WithValue(context.Background(), "trace", "t002")
	logger.Info(ctx, logger.Message{Text: "操作失败", Remark: "标题替代正文", Traceback: "按规则未找到正文",
		Attr: logger.Attr{
			Inp: map[string]interface{}{"link": "http://localhost:9087"},
			Oup: map[string]interface{}{"res": "title"},
			Monitor: logger.Monitor{
				Final:           true,
				Level:           logger.FATAL,
				Code:            logger.CHECK,
				ProcessDuration: 200,
				ProcessDelay:    20930,
			},
		},
	})
	// 2
	ctx = context.WithValue(context.Background(), "trace", "t02")
	logger.Error(ctx, logger.Message{Text: "work", Remark: "说明", Traceback: "调用超时", Label: "render",
		Attr: logger.Attr{
			Inp: map[string]interface{}{"application": "spanic-service-net"},
			Oup: map[string]interface{}{"result": ""},
			Monitor: logger.Monitor{
				Level:           logger.ERROR,
				Code:            logger.PVERROR,
				ProcessDuration: 5001,
				ProcessDelay:    2037,
			},
		},
	})
}
