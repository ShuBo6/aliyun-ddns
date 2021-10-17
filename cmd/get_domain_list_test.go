package main

import (
	"fmt"
	alidns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"testing"
)

func TestGetDomainList(t *testing.T) {

	c,err:=LoadConfig("../conf/conf.yaml")
	client, err := alidns.NewClientWithAccessKey(c.RegionID, c.AccessKey, c.SecretKey)
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = "shubo6.cn"
	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	for _, record := range response.DomainRecords.Record {
		if record.Type=="A"&&record.RR=="zoo" {
			fmt.Printf("response is %#v\n", record)
		}
	}


}
