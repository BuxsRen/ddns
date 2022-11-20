package utils

import (
	"ddns/conf"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// GetIp 取外网IP
func GetIp() string {
	responseClient, errClient := http.Get("https://www.cip.cc/") // 获取外网 IP
	if errClient != nil {
		Logs.Error("获取外网IP失败")
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(responseClient.Body)
	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))
	reg := regexp.MustCompile(`<pre>IP	: (.*?)\n`)
	res := reg.FindAllStringSubmatch(clientIP, 1)
	if len(res) > 0 {
		clientIP = res[0][1]
	}
	return clientIP
}

// SetIp 同步到域名解析
func SetIp(ip string) {
	// 配置身份验证
	credential := common.NewCredential(
		conf.Get().SecretId,
		conf.Get().SecretKey,
	)

	// 配置请求地址
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(credential, "", cpf)

	request := dnspod.NewModifyRecordRequest()

	// 构建请求参数
	request.Domain = common.StringPtr(conf.Get().Domain)
	if conf.Get().SubDomain != "@" {
		request.SubDomain = common.StringPtr(conf.Get().SubDomain)
	}
	request.RecordType = common.StringPtr("A")
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(ip)
	request.TTL = common.Uint64Ptr(600)
	request.RecordId = common.Uint64Ptr(conf.Get().RecordId)
	// 发起请求
	response, err := client.ModifyRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		Logs.Error(err)
		return
	}
	if err != nil {
		Logs.Error(err)
		return
	}
	Logs.Info(response.ToJsonString())
}

// GetDomainIp 取域名当前解析的ip
func GetDomainIp() string {
	// 配置身份验证
	credential := common.NewCredential(
		conf.Get().SecretId,
		conf.Get().SecretKey,
	)
	// 配置请求地址
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(credential, "", cpf)
	// 构建请求参数
	request := dnspod.NewDescribeRecordRequest()

	request.Domain = common.StringPtr(conf.Get().Domain)
	request.RecordId = common.Uint64Ptr(conf.Get().RecordId)
	// 发起请求
	response, err := client.DescribeRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		Logs.Error(err)
		return ""
	}
	if err != nil {
		Logs.Error(err)
		return ""
	}
	str := response.ToJsonString()
	var data map[string]interface{}

	e := json.Unmarshal([]byte(str), &data)

	if e != nil {
		Logs.Error(err)
		return ""
	}
	return data["Response"].(map[string]interface{})["RecordInfo"].(map[string]interface{})["Value"].(string)
}
