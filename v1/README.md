# OpenAPI v1
## OpenAPI 签名算法
## Usage
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
	client, _ := openapi.NewClient("192.168.100.162:8109", "c92f21653163ec91d238129c", "66555962666442714948516f4c70504c4c79686d544a795172766f73506c5967")
	c := cmdb.NewClient(openapi.WrapClient("cmdb", client))
	d, err := c.Instance.GetDetail(context.Background(), &instance.GetDetailRequest{
		ObjectId:   "APP",
		InstanceId: "582f3323eef19",
		Fields:     "name",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(d.GetFields())
}
```

