package htxp

import (
	"encoding/json"
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/zeromicro/go-zero/core/logx"
)

type AliyunSmsClient struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
	RegionId        string
}

type SmsRequest struct {
	PhoneNumbers  string `json:"PhoneNumbers"`
	SignName      string `json:"SignName"`
	TemplateCode  string `json:"TemplateCode"`
	TemplateParam string `json:"TemplateParam"`
}

type SmsResponse struct {
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	BizId     string `json:"BizId"`
}

func NewAliyunSmsClient(accessKeyId, accessKeySecret, signName string) *AliyunSmsClient {
	return &AliyunSmsClient{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		SignName:        signName,
		RegionId:        "cn-hangzhou", // 阿里云短信服务默认区域
	}
}

// 发送短信
func (c *AliyunSmsClient) SendSms(mobiles []string, templateCode string, parameters map[string]string) (*SmsResponse, error) {
	if len(mobiles) == 0 {
		logx.Errorf("手机号列表不能为空")
		return nil, errors.New("手机号列表不能为空")
	}

	// 阿里云短信服务一次只能发送给一个手机号
	phoneNumber := mobiles[0]

	// 将参数转换为JSON字符串
	templateParam, err := json.Marshal(parameters)
	if err != nil {
		logx.Errorf("参数序列化失败: %v", err)
		return nil, errors.New("参数序列化失败")
	}

	// 创建阿里云短信客户端
	client, err := dysmsapi.NewClientWithAccessKey(c.RegionId, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		logx.Errorf("创建阿里云客户端失败: %v", err)
		return nil, errors.New("创建阿里云客户端失败")
	}

	// 创建发送短信请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumber
	request.SignName = c.SignName
	request.TemplateCode = templateCode
	request.TemplateParam = string(templateParam)

	// 打印请求信息用于调试
	logx.Infof("发送短信请求: PhoneNumbers=%s, TemplateCode=%s, TemplateParam=%s",
		phoneNumber, templateCode, string(templateParam))

	// 发送短信
	response, err := client.SendSms(request)
	if err != nil {
		logx.Errorf("发送短信失败: %v", err)
		return nil, errors.New("发送短信失败")
	}

	// 打印响应信息用于调试
	logx.Infof("短信发送响应: RequestId=%s, Code=%s, Message=%s, BizId=%s", response.RequestId, response.Code, response.Message, response.BizId)

	// 构建响应对象
	smsResp := &SmsResponse{
		RequestId: response.RequestId,
		Code:      response.Code,
		Message:   response.Message,
		BizId:     response.BizId,
	}

	// 判断是否成功
	if response.Code == "OK" {
		return smsResp, nil
	}
	logx.Errorf("短信发送失败: %s", response.Message)
	return smsResp, errors.New("短信发送失败")
}
