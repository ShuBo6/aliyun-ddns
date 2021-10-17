package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"time"
)
type Config struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Domain    string `yaml:"domain"`
	Record    string `yaml:"record"`
	RegionID  string `yaml:"region_id"`
}


func GetIP() (string, error) {
	responseClient, errClient := http.Get("https://ipw.cn/api/ip/myip") // 获取外网 IP
	if errClient != nil {
		return "", errClient
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer responseClient.Body.Close()

	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))
	return clientIP, nil
}
func LoadConfig(filename string) (*Config, error) {
	c := new(Config)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func GetRecord(client *alidns.Client, conf *Config) (*alidns.Record, string) {
	res := new(alidns.Record)
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = conf.Domain
	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return nil, fmt.Sprintf("%s", err.Error())
	}
	for _, record := range response.DomainRecords.Record {
		if record.Type == "A" && record.RR == conf.Record {
			fmt.Printf("old Record is %#v\n", record)
			res = &record
		}
	}
	if res == nil {
		return nil, fmt.Sprintf("can't find a record(%s) in domain(%s) with type(%s)", conf.Record, conf.Domain, "A")
	}
	return res, ""
}
func UpdateRecord(client *alidns.Client, conf *Config, record *alidns.Record, ip string) string {
	if record.Value==ip {
		return fmt.Sprintf("ip(%s) need't update.",record.Value)
	}else {
		req := alidns.CreateUpdateDomainRecordRequest()
		req.RecordId = record.RecordId
		req.RR = conf.Record
		req.Type = "A"
		req.Value = ip
		_, err := client.UpdateDomainRecord(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		return 	fmt.Sprintf("update success!,new record(%s.%s) IP(%s)", conf.Record,conf.Domain,ip)
	}


}
func main() {
	for  {
		ip, err := GetIP()
		if err != nil {
			fmt.Print(err.Error())
		}
		config, err := LoadConfig("conf.yaml")
		if err != nil {
			fmt.Print(err.Error())
		}
		client, err := alidns.NewClientWithAccessKey(config.RegionID, config.AccessKey, config.SecretKey)
		if err != nil {
			fmt.Print(err.Error())
		}
		record, errMassge := GetRecord(client, config)
		if errMassge != "" {
			fmt.Print(errMassge)
		}
		fmt.Println(UpdateRecord(client, config, record, ip))
		time.Sleep(time.Minute*10)
	}


}
