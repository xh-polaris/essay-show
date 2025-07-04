package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
)

var client *HttpClient

// HttpClient 是一个简单的 HTTP 客户端
type HttpClient struct {
	Client *http.Client
	Config *config.Config
}

// NewHttpClient 创建一个新的 HttpClient 实例
func NewHttpClient() *HttpClient {
	return &HttpClient{
		Client: &http.Client{},
	}
}

func GetHttpClient() *HttpClient {
	if client == nil {
		client = NewHttpClient()
	}
	return client
}

type params struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// SendRequest 发送 HTTP 请求
func (c *HttpClient) SendRequest(ctx context.Context, method, url string, headers map[string]string, body interface{}) (map[string]interface{}, error) {
	// 将 body 序列化为 JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("请求体序列化失败: %w", err)
	}

	// 创建新的请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.WithContext(ctx)

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("关闭请求失败: %v", closeErr)
		}
	}()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("unexpected status code: %d, response body: %s", resp.StatusCode, responseBody)
		return nil, fmt.Errorf(errMsg)
	}

	// 反序列化响应体
	var responseMap map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseMap); err != nil {
		return nil, fmt.Errorf("反序列化响应失败: %w", err)
	}

	return responseMap, nil
}

// SignUp 用于用户初始化
func (c *HttpClient) SignUp(ctx context.Context, authType string, authId string, verifyCode *string) (map[string]interface{}, error) {

	body := make(map[string]interface{})
	body["authType"] = authType
	body["authId"] = authId
	body["verifyCode"] = *verifyCode
	body["appId"] = consts.AppId

	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["Charset"] = consts.CharSetUTF8

	// 如果是测试环境则向测试环境的中台发送请求
	if config.GetConfig().State == "test" {
		header["X-Xh-Env"] = "test"
	}

	resp, err := c.SendRequest(ctx, consts.Post, consts.PlatformSignInUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SignIn 用于用户登录
func (c *HttpClient) SignIn(ctx context.Context, authType string, authId string, verifyCode *string, password *string) (map[string]interface{}, error) {

	body := make(map[string]interface{})
	body["authType"] = authType
	body["authId"] = authId
	if verifyCode != nil {
		body["verifyCode"] = *verifyCode
	}
	if password != nil {
		body["password"] = *password
	}
	body["appId"] = consts.AppId

	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["Charset"] = consts.CharSetUTF8

	// 如果是测试环境则向测试环境中台发送请求
	if config.GetConfig().State == "test" {
		header["X-Xh-Env"] = "test"
	}

	resp, err := c.SendRequest(ctx, consts.Post, consts.PlatformSignInUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetPassword 用于用户登录
func (c *HttpClient) SetPassword(ctx context.Context, authorization string, password string) (map[string]interface{}, error) {

	body := make(map[string]interface{})
	body["password"] = password
	body["appId"] = consts.AppId

	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["Charset"] = consts.CharSetUTF8
	header["Authorization"] = authorization

	// 如果是测试环境则向测试环境中台发送请求
	if config.GetConfig().State == "test" {
		header["X-Xh-Env"] = "test"
	}

	resp, err := c.SendRequest(ctx, consts.Post, consts.PlatformSetPasswordUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SendVerifyCode SetPassword 用于用户登录
func (c *HttpClient) SendVerifyCode(ctx context.Context, authType string, authId string) (map[string]interface{}, error) {

	body := make(map[string]interface{})
	body["authType"] = authType
	body["authId"] = authId

	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["Charset"] = consts.CharSetUTF8

	// 如果是测试环境则向测试环境中台发送请求
	if config.GetConfig().State == "test" {
		header["X-Xh-Env"] = "test"
	}

	resp, err := c.SendRequest(ctx, consts.Post, consts.PlatformSendVerifyCodeUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// BetaEvaluate 用Beta接口进行批改
func (c *HttpClient) BetaEvaluate(ctx context.Context, title string, text string, grade *int64, essayType *string) (map[string]interface{}, error) {

	body := make(map[string]interface{})

	// 请求体
	body["title"] = title
	body["content"] = text
	if grade != nil {
		body["grade"] = *grade
	}
	if essayType != nil {
		body["essayType"] = *essayType
	}

	// 请求头
	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["Charset"] = consts.CharSetUTF8

	// 如果是测试环境则向测试环境发送请求
	if config.GetConfig().State == "test" {
		header["X-Xh-Env"] = "test"
	}

	resp, err := c.SendRequest(ctx, consts.Post, consts.BetaEvaluateUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// BeeTitleUrlOCR 蜜蜂ocr - 带标题
func (c *HttpClient) BeeTitleUrlOCR(ctx context.Context, images []string, left string) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	// 图片url列表
	body["images"] = images
	// 保留类型
	if len(left) > 0 {
		body["leftType"] = left
	}

	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	if config.GetConfig().State == "test" {
		header["X-Xh-Env"] = "test"
	}

	resp, err := c.SendRequest(ctx, consts.Post, consts.BeeTitleUrlOcr, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
