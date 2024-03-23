package aliyun

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/coscms/sms"
	"github.com/webx-top/com"
)

var _ sms.Sender = &Aliyun{}

func New() *Aliyun {
	return &Aliyun{
		RegionId: `cn-hangzhou`,
	}
}

type Aliyun struct {
	RegionId     string
	AccessKey    string
	AccessSecret string
	SignName     string //默认签名
	TmplCode     string //默认模板代码
	client       *dysmsapi.Client
}

func (a *Aliyun) Send(c *sms.Config) error {
	tmplCode := c.Template
	signName := c.SignName
	//code := fmt.Sprint(c.Extra(`code`))
	if len(tmplCode) == 0 {
		tmplCode = a.TmplCode
	}
	if len(signName) == 0 {
		signName = a.SignName
	}
	if a.client == nil {
		var err error
		a.client, err = dysmsapi.NewClientWithAccessKey(a.RegionId, a.AccessKey, a.AccessSecret)
		if err != nil {
			return err
		}
	}
	b, e := com.JSONEncode(c.ExtraData)
	if e != nil {
		return e
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = c.Mobile   //手机号变量值
	request.SignName = signName       //签名
	request.TemplateCode = tmplCode   //模板编码
	request.TemplateParam = string(b) //"{\"code\":\"" + code + "\"}"
	response, err := a.client.SendSms(request)
	if err != nil {
		return err
	}
	if response.IsSuccess() {
		return nil
	}
	// if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {}
	return errors.New(`AliyunSMS: ` + response.Message)
}
