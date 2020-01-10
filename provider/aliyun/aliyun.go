package aliyun

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/messagebird/sachet"
)

type AliyunConfig struct {
	AccessKey       string `yaml:"access_key"`
	AccessKeySecret string `yaml:"access_key_secret"`

	SignName         string `yaml:"sign_name"`
	TemplateCode     string `yaml:"template_code"`
	TemplateParamKey string `yaml:"template_param_key"`
}

type Aliyun struct {
	client *dysmsapi.Client
	config *AliyunConfig
}

func NewAliyun(config AliyunConfig) (*Aliyun, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", config.AccessKey, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return &Aliyun{
		client: client,
		config: &config,
	}, nil
}

func (aliyun *Aliyun) Send(message sachet.Message) error {
	var err error = nil
	switch message.Type {
	case "", "text":
		request := dysmsapi.CreateSendSmsRequest()
		request.Scheme = "https"
		request.SignName = aliyun.config.SignName
		request.TemplateCode = aliyun.config.TemplateCode
		request.PhoneNumbers = strings.Join(message.To, ",")
		templateParam := make(map[string]string)
		templateParam[aliyun.config.TemplateParamKey] = message.Text
		templateParamByte, err := json.Marshal(templateParam)
		if err == nil {
			request.TemplateParam = string(templateParamByte)
			var response *dysmsapi.SendSmsResponse
			response, err = aliyun.client.SendSms(request)
			if err == nil && (!response.IsSuccess() || response.Code != "") {
				err = fmt.Errorf(response.String())
			}
		}
	default:
		return fmt.Errorf("unknown message type %s", message.Type)
	}
	return err
}
