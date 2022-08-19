package httpclient

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

type Client struct {
	HttpClient *http.Client
}

// DefaultClient 超时时间为10秒的缺省实例指针
var DefaultClient *Client

// New 返回Client实例指针
//
// 传入参数
//
//	timeSeconds 访问超时秒数
func New(timeSeconds int) *Client {
	if timeSeconds <= 0 {
		timeSeconds = 10
	}
	return &Client{
		HttpClient: &http.Client{
			Timeout: time.Duration(timeSeconds) * time.Second,
		},
	}
}

func init() {
	timeSeconds := g.Cfg().MustGet(gctx.New(), "httpclient.timeoutSeconds").Int()
	if timeSeconds <= 0 {
		timeSeconds = 10
	}
	DefaultClient = New(timeSeconds)
}

// Exec 通过给定方式访问URL，获取结果
//
// 传入参数：
//
//	ctx
//	request        HttpRequest 结构体
//
// 返回参数
//
//	result         HttpResponse 结构体
//	err            错误
func (c *Client) Exec(ctx context.Context, request *HttpRequest) (result *HttpResponse, err error) {
	req, _ := http.NewRequest(gstr.ToUpper(request.Method), request.Url, getReader(request.Body))
	if !g.IsEmpty(request.Headers) {
		req.Header = *request.Headers
	}
	if !g.IsEmpty(request.Cookies) {
		for _, value := range request.Cookies {
			req.AddCookie(value)
		}
	}

	retryCount := request.RetryCount
	if retryCount < 0 {
		retryCount = 1000
	}

	for i := 0; i < retryCount; i++ {
		result, err = execRequest(ctx, c.HttpClient, req)
		if err == nil {
			return result, nil
		}
		g.Log().Debugf(ctx, "%d times retry get %s", i+1, request.Url)
	}
	g.Log().Warningf(ctx, "aborting after retried %d times for %s", retryCount, request.Url)
	return nil, nil
}

// DoGet GET方式访问URL，获取结果，是Exec的快捷方法
//
// 传入参数：
//
//	ctx
//	url            要访问的url
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoGet(ctx context.Context, url string, retryCount int) (result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodGet,
		Url:        url,
		RetryCount: retryCount,
	})
}

// DoGetWithHeaders GET方式访问URL，获取结果，是Exec的快捷方法
//
// 传入参数：
//
//	ctx
//	url            要访问的url
//	header         请求的Http header
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoGetWithHeaders(ctx context.Context, url string, headers *http.Header, retryCount int) (result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodGet,
		Url:        url,
		Headers:    headers,
		RetryCount: retryCount,
	})
}

// DoGetWithCookies GET方式访问URL，获取结果，是Exec的快捷方法
//
// 传入参数：
//
//	ctx
//	url            要访问的url
//	cookies        请求的Http cookies
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoGetWithCookies(ctx context.Context, url string, cookies []*http.Cookie, retryCount int) (result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodGet,
		Url:        url,
		Cookies:    cookies,
		RetryCount: retryCount,
	})
}

// DoGetWithHeadersAndCookies GET方式访问URL，获取结果，是Exec的快捷方法
// 传入参数：
//
//	ctx
//	url            要访问的url
//	header         请求的Http header
//	cookies        请求的Http cookies
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoGetWithHeadersAndCookies(ctx context.Context, url string, headers *http.Header, cookies []*http.Cookie, retryCount int) (result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodGet,
		Url:        url,
		Headers:    headers,
		Cookies:    cookies,
		RetryCount: retryCount,
	})
}

// DoPost Post方式访问URL，获取结果，是Exec的快捷方法
//
// 传入参数：
//
//	ctx
//	url            要访问的url
//	body           请求的Http body
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoPost(ctx context.Context, url string, body string, retryCount int) (result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodPost,
		Url:        url,
		Body:       body,
		RetryCount: retryCount,
	})
}

// DoPostWithHeaders Post方式访问URL，获取结果，是Exec的快捷方法
//
// 传入参数：
//
//	ctx
//	url            要访问的url
//	header         请求的Http header
//	body           请求的Http body
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoPostWithHeaders(ctx context.Context, url string, headers *http.Header, body string, retryCount int) (result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodPost,
		Url:        url,
		Headers:    headers,
		Body:       body,
		RetryCount: retryCount,
	})
}

// DoPostWithCookies Post方式访问URL，获取结果，是Exec的快捷方法
//
// 传入参数：
//
//	ctx
//	url            要访问的url
//	cookies        请求的Http cookies
//	body           请求的Http body
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoPostWithCookies(ctx context.Context, url string, cookies []*http.Cookie, body string, retryCount int) (result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodPost,
		Url:        url,
		Cookies:    cookies,
		Body:       body,
		RetryCount: retryCount,
	})
}

// DoPostWithHeadersAndCookies Post方式访问URL，获取结果，是Exec的快捷方法
//
// 传入参数：
//
//	ctx
//	url            要访问的url
//	header         请求的Http header
//	cookies        请求的Http cookies
//	body           请求的Http body
//	retryCount     重试次数，如为-1则无限重试
//
// 返回参数
//
//	result         HttpResponse结构体
//	err            错误
func (c *Client) DoPostWithHeadersAndCookies(ctx context.Context, url string, headers *http.Header, cookies []*http.Cookie, body string, retryCount int) (
	result *HttpResponse, err error) {
	return c.Exec(ctx, &HttpRequest{
		Method:     MethodPost,
		Url:        url,
		Headers:    headers,
		Cookies:    cookies,
		Body:       body,
		RetryCount: retryCount,
	})
}

func execRequest(ctx context.Context, client *http.Client, req *http.Request) (*HttpResponse, error) {
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("http get error", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	return &HttpResponse{
		Body:       string(body),
		StatusCode: resp.StatusCode,
		Header:     &resp.Header,
	}, nil
}

func getReader(body string) io.Reader {
	if g.IsEmpty(body) {
		return nil
	}
	return strings.NewReader(body)
}