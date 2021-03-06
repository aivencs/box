package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/aivencs/box/pkg/kit"
	"github.com/aivencs/box/pkg/logger"
	"github.com/aivencs/box/pkg/validate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 使用枚举限定选择
type TypeSupport string

type MethodType string

const (
	SERVER_ECHO TypeSupport = "echo"
	// 限定请求方式
	GET    MethodType = "GET"
	POST   MethodType = "POST"
	DELETE MethodType = "DELETE"
	PUT    MethodType = "PUT"
	// 定义默认值
	DEFAULT_HOST = ":"
)

var server Server
var once sync.Once

func init() {
	ctx := context.WithValue(context.Background(), "trace", "init-for-server")
	validate.InitValidate(ctx, validate.VALIDATOR, validate.Option{})
}

type Server interface {
	Work()
	AddRouter(payload RouterPayload, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
}

type EchoServer struct {
	Kernel *echo.Echo
	Port   int
	Host   string
}

type Option struct {
	Host               string `json:"host" label:"服务地址" desc:"默认为开放访问"`
	Port               int    `json:"port" label:"端口号" validate:"required,min=3000,max=10000"`
	DisableMiddCors    bool   `json:"disable_midd_cors" label:"cors中间件开关" desc:"默认开启"`
	DisableMiddRecover bool   `json:"disable_midd_recover" label:"recover中间件开关" desc:"默认开启"`
}

func InitServer(ctx context.Context, support TypeSupport, option Option) error {
	c := server
	message, err := validate.Work(ctx, option)
	if err != nil {
		return errors.New(message)
	}
	once.Do(func() {
		c, err = ServerFactory(ctx, support, option)
		if err != nil {
			return
		}
		if c == nil {
			err = errors.New("初始化失败")
		}
		server = c
	})
	return err
}

func ServerFactory(ctx context.Context, support TypeSupport, option Option) (Server, error) {
	switch support {
	case SERVER_ECHO:
		return NewEchoServer(ctx, option)
	default:
		return NewEchoServer(ctx, option)
	}
}

func NewEchoServer(ctx context.Context, option Option) (Server, error) {
	svr := echo.New()
	if !option.DisableMiddCors {
		svr.Use(middleware.CORS())
	}
	if !option.DisableMiddRecover {
		svr.Use(middleware.Recover())
	}
	if len(option.Host) == 0 {
		option.Host = DEFAULT_HOST
	} else {
		option.Host = kit.JoinString(option.Host, DEFAULT_HOST)
	}
	return &EchoServer{
		Kernel: svr,
		Port:   option.Port,
		Host:   option.Host,
	}, nil
}

func (c *EchoServer) Work() {
	port := strconv.Itoa(c.Port)
	c.Kernel.Logger.Fatal(c.Kernel.Start(kit.JoinString(c.Host, port)))
}

type RouterPayload struct {
	Method MethodType
	Path   string
	Label  string
}

func (c *EchoServer) AddRouter(payload RouterPayload, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	switch payload.Method {
	case GET:
		c.Kernel.GET(payload.Path, h, m...)
	case POST:
		c.Kernel.POST(payload.Path, h, m...)
	case DELETE:
		c.Kernel.DELETE(payload.Path, h, m...)
	case PUT:
		c.Kernel.PUT(payload.Path, h, m...)
	default:
		c.Kernel.GET(payload.Path, h, m...)
	}
}

func Work() {
	server.Work()
}

func AddRouter(payload RouterPayload, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	server.AddRouter(payload, h, m...)
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// 接口响应结果
type ServerResponse struct {
	Code    logger.Code `json:"code"`
	Trace   string      `json:"trace"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func EmptyHandler(c echo.Context) error {
	res := ServerResponse{
		Code:    logger.PVERROR,
		Trace:   c.Get("trace").(string),
		Message: c.Get("message").(string),
		Result:  nil,
	}
	return c.JSONPretty(http.StatusOK, res, "")
}

type Header struct {
	X_REQUEST_ID string `json:"X-REQUEST-ID" label:"追踪编码" validate:"required,min=16,max=100"`
}

type T struct {
	A int64
	B float64
}

// 基础日志中间件
func loggerBase(next echo.HandlerFunc, inp bool, oup bool) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var header Header
		startT := time.Now()
		response := ServerResponse{}
		// 获取追踪编码
		ids := c.Request().Header.Values("X-REQUEST-ID")
		if len(ids) > 0 {
			header.X_REQUEST_ID = ids[0]
		}
		// 创建新的Context
		ctx := context.WithValue(context.Background(), "trace", header.X_REQUEST_ID)
		ctx = context.WithValue(ctx, "label", c.Request().URL.Path)
		// 设置框架的Context
		c.Set("trace", header.X_REQUEST_ID)
		c.Set("label", c.Request().URL.Path)
		c.Set("context", ctx)
		// 校验追踪编码
		message, err := validate.Work(context.Background(), &header)
		// 拦截响应
		responseBuffer := new(bytes.Buffer)
		mw := io.MultiWriter(c.Response().Writer, responseBuffer)
		writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
		c.Response().Writer = writer
		if err != nil {
			// 通过调用空接口避免接口和日志输出不统一
			c.Set("message", message)
			EmptyHandler(c)
		} else {
			// 调用目标接口
			err = next(c)
		}
		// 解析响应内容
		json.Unmarshal(responseBuffer.Bytes(), &response)
		duration := time.Since(startT).Milliseconds()
		// 构建日志信息
		input := map[string]interface{}{
			"host":       c.Request().Host,
			"path":       c.Path(),
			"user-agent": c.Request().UserAgent(),
			"method":     c.Request().Method,
		}
		output := map[string]interface{}{
			"status_code": c.Response().Status,
		}
		if inp {
			input["param"] = c.Get("param")
		}
		if oup {
			output["response"] = response
		}
		logger.Info(ctx, logger.Message{
			Text:  response.Message,
			Label: c.Request().URL.Path,
			Attr: logger.Attr{
				Monitor: logger.Monitor{
					Final:           true,
					ProcessDuration: duration,
					Code:            response.Code,
					Level:           logger.GetLevelBaseCode(response.Code),
				},
				Inp: input,
				Oup: output,
			},
		})
		return
	}
}

// 不带输入输出的日志中间件
func LoggerNormal(next echo.HandlerFunc) echo.HandlerFunc {
	return loggerBase(next, false, false)
}

// 附带输入输出的日志中间件
func LoggerPlus(next echo.HandlerFunc) echo.HandlerFunc {
	return loggerBase(next, true, true)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
