package request

import (
	"context"
	"crypto/tls"
	"errors"
	"net/url"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/aivencs/box/pkg/logger"
	"github.com/aivencs/box/pkg/validate"
	"github.com/go-resty/resty/v2"
)

// 使用枚举限定使用时的选项
type TypeSupport string
type MethodType string

const (
	RESTY TypeSupport = "resty"
	// 限定请求方式
	GET  MethodType = "GET"
	POST MethodType = "POST"
	// 定义默认值
	DEFAULT_TIEOUT = 10
)

var once sync.Once

// 定义对象
var request Request

func init() {
	ctx := context.WithValue(context.Background(), "trace", "init-for-request")
	validate.InitValidate(ctx, validate.VALIDATOR, validate.Option{})
}

type Request interface {
	Get(ctx context.Context, param Param) (Result, error)
	Post(ctx context.Context, param Param) (Result, error)
}

// 请求结果
type Result struct {
	Text       string
	StatusCode int
}

// 初始化时所用参数
type Option struct {
}

// 结构体
// 基于Resty
type RestyRequest struct{}

// 请求参数
type Param struct {
	Link             string     `json:"link" label:"网址" validate:"required,url"`
	Method           MethodType `json:"method" label:"请求方式"`
	Payload          string     `json:"payload" label:"参数"`
	Timeout          int        `json:"timeout" label:"超时时间"`
	Proxy            string     `json:"proxy" label:"IP代理"`
	EnableSkipVerify bool       `json:"enable_skip_verify" label:"跳过证书" desc:"默认不开启"`
	EnableHeader     bool       `json:"enable_header" label:"根据网址设置请求头基本参数" desc:"默认不开启"`
}

// 初始化对象
func InitRequest(ctx context.Context, support TypeSupport, option Option) error {
	c := request
	var err error
	message, err := validate.Work(ctx, &option)
	if err != nil {
		return logger.NewErr(logger.ErrOption{Code: logger.PVERROR, Label: message, Err: err})
	}
	once.Do(func() {
		c = RequestFactory(ctx, support, option)
		if c == nil {
			err = errors.New("初始化失败")
		}
		request = c
	})
	if err != nil {
		return logger.NewErr(logger.ErrOption{Code: logger.DVERROR, Err: err})
	}
	return nil
}

// 抽象工厂
func RequestFactory(ctx context.Context, support TypeSupport, option Option) Request {
	switch support {
	case RESTY:
		return NewRestyRequest(ctx, option)
	default:
		return NewRestyRequest(ctx, option)
	}
}

// 创建基于Resty的请求对象
func NewRestyRequest(ctx context.Context, option Option) Request {
	return &RestyRequest{}
}

func (c *RestyRequest) Get(ctx context.Context, param Param) (Result, error) {
	return c.work(ctx, param)
}

func (c *RestyRequest) Post(ctx context.Context, param Param) (Result, error) {
	return c.work(ctx, param)
}

func (c *RestyRequest) work(ctx context.Context, param Param) (Result, error) {
	var response *resty.Response
	var err error
	// 参数校验
	message, err := validate.Work(ctx, param)
	if err != nil {
		return Result{}, logger.NewErr(logger.ErrOption{Code: logger.PVERROR, Err: err, Label: message})
	}
	// 前期准备
	serviceSafeString, _ := url.Parse(param.Link)
	client := resty.New()
	// 设置追踪编码
	client.SetHeaders(map[string]string{"X-REQUEST-ID": ctx.Value("trace").(string)})
	// 参数项的应用
	if param.Timeout == 0 {
		param.Timeout = DEFAULT_TIEOUT
	}
	client.SetTimeout(time.Duration(param.Timeout) * time.Second)
	if param.EnableSkipVerify {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if param.EnableHeader {
		client.SetHeaders(map[string]string{
			"X-REQUEST-ID": ctx.Value("trace").(string),
			"Host":         serviceSafeString.Host,
			"Referer":      serviceSafeString.Host,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		})
	}
	if utf8.RuneCountInString(param.Proxy) > 6 {
		client.SetProxy(param.Proxy)
	}
	// 发出请求
	switch param.Method {
	case GET:
		response, err = client.R().SetBody(param.Payload).Get(param.Link)
	case POST:
		client.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
		response, err = client.R().SetBody(param.Payload).Post(param.Link)
	default:
		response, err = client.R().SetBody(param.Payload).Get(param.Link)
	}
	// 请求结果处理
	if err != nil {
		return Result{}, logger.NewErr(logger.ErrOption{Code: logger.DVERROR, Err: err, Label: "请求时发生错误"})
	}
	// 状态码处理
	if response.RawResponse.StatusCode > 201 {
		switch response.RawResponse.StatusCode {
		case 429:
			err = logger.NewErr(logger.ErrOption{Code: logger.LIMITERROR})
		case 404:
			err = logger.NewErr(logger.ErrOption{Code: logger.CHECK, Label: "资源不存在"})
		case 200:
			err = nil
		case 201:
			err = nil
		default:
			err = logger.NewErr(logger.ErrOption{Code: logger.STATUSERROR})
		}
	}
	// 构造结果并返回
	return Result{
		Text:       response.String(),
		StatusCode: response.RawResponse.StatusCode,
	}, err
}

// 暴露给外部调用
func Get(ctx context.Context, param Param) (Result, error) {
	param.Method = GET
	return request.Get(ctx, param)
}

// 暴露给外部调用
func Post(ctx context.Context, param Param) (Result, error) {
	param.Method = POST
	return request.Get(ctx, param)
}
