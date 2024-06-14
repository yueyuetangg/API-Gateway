package main

import (
	"bytes"
	"context"
	"demo_student/kitex_gen/demo"
	"demo_student/kitex_gen/demo/studentservice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"net/http"
	"time"

	//"github.com/cloudwego/kitex/server"
	"log"
)

var (
	cli studentservice.Client
)

func main() {
	c, err := studentservice.NewClient("demo", client.WithHostPorts("127.0.0.1:9999"),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithRPCTimeout(500*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	cli = c
	hz := server.New(server.WithHostPorts("localhost:8889"))
	hz.GET("/query", Handler)
	hz.POST("/add-student-info", Addlder)
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}
func Handler(ctx context.Context, c *app.RequestContext) {
	var err error
	var req demo.QueryReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	// 调用服务端的 Query 方法查询学生信息
	resp, err := cli.Query(ctx, &demo.QueryReq{
		Id: req.Id,
	})
	if err != nil {
		c.String(consts.StatusInternalServerError, "调用 Query 函数失败")
		return
	}

	// 返回查询结果给客户端
	c.JSON(consts.StatusOK, resp)

}
func Addlder(ctx context.Context, c *app.RequestContext) {
	var student demo.Student
	if err := json.NewDecoder(bytes.NewReader(c.Request.Body())).Decode(&student); err != nil {
		log.Println("解析 JSON 数据失败:", err)
		c.String(http.StatusBadRequest, "解析 JSON 数据失败")
		return
	}

	// 调用服务端的 Register 函数
	resp, err := cli.Register(ctx, &student)
	if err != nil {
		log.Println("调用 Register 函数失败:", err)
		c.String(http.StatusInternalServerError, "调用 Register 函数失败")
		return
	}

	// 返回注册结果给客户端
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Println("序列化响应数据失败:", err)
		c.String(http.StatusInternalServerError, "序列化响应数据失败")
		return
	}

	c.String(http.StatusOK, string(respJSON))

}
