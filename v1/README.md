# OpenAPI v1
## OpenAPI 签名算法
## Usage
```shell script
go get github.com/easyopsapis/openapi-go/v1
```
```go
package main

import (
	"context"
	cmdb "github.com/easyopsapis/easyops-api-go/protorepo-cmdb"
	"github.com/easyopsapis/easyops-api-go/protorepo-cmdb/instance"
	"github.com/easyopsapis/openapi-go/v1"
	"log"
)

func main() {
	// OpenAPI服务地址: 192.168.100.162:8109
	// accessKey: c92f21653163ec91d238129c
	// secretKey: 66555962666442714948516f4c70504c4c79686d544a795172766f73506c5967
	gw, _ := openapi.NewClient("192.168.100.162:8109", "c92f21653163ec91d238129c", "66555962666442714948516f4c70504c4c79686d544a795172766f73506c5967")

	// 根据服务名封装client, 注入cmdb client中
	c := cmdb.NewClient(openapi.WrapClient("cmdbservice", gw)) // 服务名: cmdbservice

	// 利用 cmdb client 查询实例详情
	d, err := c.Instance.GetDetail(context.Background(), &instance.GetDetailRequest{
		ObjectId:   "APP",
		InstanceId: "582f3323eef19",
		Fields:     "name",
	})

	if err != nil {
		// 当请求返回错误时, 利用 gerr 包从 err 中提取 Status 信息 ( 包含错误码, 错误信息等信息 )
		status := gerr.FromError(err)
		fmt.Println(status.Code(), status.Message())
		os.Exit(1)
	}

	fmt.Println(d.Fields["name"].GetStringValue())
}
```

