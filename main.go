package main

import (
	"ddns/conf"
	"ddns/utils"
	"fmt"
	"time"
)

// 自动获取外网ip并同步到到腾讯云域名解析
func main() {

	fmt.Println("DDNS动态域名IP解析腾讯云专版")
	fmt.Printf("当前解析的域名：%v\n", conf.Get().Domain)
	utils.Logs.Info("[服务运行中]")

	tk := time.NewTicker(time.Duration(conf.Get().Time) * time.Minute)

	var record = utils.GetDomainIp()

	for {
		select {
		case <-tk.C:
			ip := utils.GetIp()
			if ip != "" && record != "" {
				if ip != record {
					utils.Logs.Infof("[IP变动] %v -> %v", record, ip)
					utils.SetIp(ip)
					record = utils.GetDomainIp()
				}
			}
		}
	}

}
