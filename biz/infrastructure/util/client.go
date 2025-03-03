package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/xh-polaris/essay-show/biz/infrastructure/config"
	"github.com/xh-polaris/essay-show/biz/infrastructure/consts"
	"io"
	"log"
	"net/http"
	"time"
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

type evaluateRequestBody struct {
	Key       string `json:"key"`
	Method    string `json:"method"`
	Action    string `json:"action"`
	Params    string `json:"params"` // string 类型
	Timestamp int64  `json:"timestamp"`
}

type params struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// SendRequest 发送 HTTP 请求
func (c *HttpClient) SendRequest(method, url string, headers map[string]string, body interface{}) (map[string]interface{}, error) {
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
func (c *HttpClient) SignUp(authType string, authId string, verifyCode *string) (map[string]interface{}, error) {

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

	resp, err := c.SendRequest(consts.Post, consts.PlatformSignInUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SignIn 用于用户登录
func (c *HttpClient) SignIn(authType string, authId string, verifyCode *string, password *string) (map[string]interface{}, error) {

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

	resp, err := c.SendRequest(consts.Post, consts.PlatformSignInUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetPassword 用于用户登录
func (c *HttpClient) SetPassword(authorization string, password string) (map[string]interface{}, error) {

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

	resp, err := c.SendRequest(consts.Post, consts.PlatformSetPasswordUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SendVerifyCode SetPassword 用于用户登录
func (c *HttpClient) SendVerifyCode(authType string, authId string) (map[string]interface{}, error) {

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

	resp, err := c.SendRequest(consts.Post, consts.PlatformSendVerifyCodeUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// BetaEvaluate 用Beta接口进行批改
func (c *HttpClient) BetaEvaluate(title string, text string) (map[string]interface{}, error) {

	//body := make(map[string]interface{})
	jsonParams, err := json.Marshal(params{
		Title:   title,
		Content: text,
	})
	if err != nil {
		return nil, consts.ErrInvalidParams
	}

	body := evaluateRequestBody{
		Key:       config.GetConfig().CallKey,
		Method:    consts.Post,
		Action:    consts.Beta,
		Params:    string(jsonParams),
		Timestamp: time.Now().Unix(),
	}

	//body["key"] = config.GetConfig().CallKey
	//body["method"] = consts.Post
	//body["action"] = consts.Beta
	//body["params"] = string(jsonParams)
	//now := time.Now().Unix()
	//body["timestamp"] = now

	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["Charset"] = consts.CharSetUTF8
	header["Signature"] = signature(body)

	// 如果是测试环境则向测试环境发送请求
	//if config.GetConfig().State == "test" {
	//	header["X-Xh-Env"] = "test"
	//}
	header["X-Xh-Env"] = "test" // 暂时都先用测试环境批改

	resp, err := c.SendRequest(consts.Post, consts.OpenApiCallUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// BeeOCR 蜜蜂ocr
func (c *HttpClient) BeeOCR(url string) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	body["image_url"] = url

	header := make(map[string]string)
	header["Content-Type"] = consts.ContentTypeJson
	header["x-app-secret"] = config.GetConfig().OCRSecret
	header["x-app-key"] = config.GetConfig().OCRKey

	resp, err := c.SendRequest(consts.Post, consts.BeeOCRUrl, header, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func signature(b evaluateRequestBody) string {

	bodyBytes, err := json.Marshal(b)
	if err != nil {
		return ""
	}

	// 签名
	hash := sha256.Sum256(bodyBytes)
	hashedSign := hex.EncodeToString(hash[:])

	return hashedSign
}
