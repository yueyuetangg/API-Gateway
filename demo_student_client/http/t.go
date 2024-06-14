package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"

	//"demo_student_client/kitex_gen/demo/studentservice"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	//"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	//"github.com/cloudwego/hertz/pkg/protocol/http1/req"
	//"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"net/http"

	//"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
)

var (
	cli genericclient.Client
)

func main() {

	//基于本地文件解析IDL实现idl热更新
	content := `
	 namespace go demo

    struct College {
        1: required string name(go.tag = 'json:"name"'),
        2: string address(go.tag = 'json:"address"'),
    }

    struct Student {
        1: required i32 id(api.body='id'),
        2: required string name(api.body='name'),
        3: required College college(api.body='college'),
        4: optional list<string> email(api.body='email'),
    }

    struct RegisterResp {
        1: bool success(api.body='success'),
        2: string message(api.body='message'),
    }

    struct QueryReq {
        1: required i32 id(api.query='id')
    }

    service StudentService {
        RegisterResp Register(1: Student student)(api.post = '/add-student-info')
        Student Query(1: QueryReq req)(api.get = '/query')
    }
	`
	//path := "./idl/student.thrift"
	////includes := map[string]string{
	////	path: content,
	////}
	p, err := generic.NewThriftContentProviderWithDynamicGo(content, nil)
	if err != nil {
		panic(err)
		fmt.Println("idl pro err")

	}
	go func() {
		// dynamic update
		err = p.UpdateIDL(content, nil)
		if err != nil {
			fmt.Println("up err")
		}
	}()
	//p, err := generic.NewThriftFileProvider("./idl/student.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.HTTPThriftGeneric(p)
	//注册发现
	resolver, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	c, err := genericclient.NewClient("demo", g, client.WithResolver(resolver), client.WithTag("Cluster", "student"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c
	hz := server.New(server.WithHostPorts("localhost:8887"))
	// 创建 HTTP 服务器
	//hz.Handle = Handler
	//update-idl 接口，问题：thrift放到Json不能换行+转义字符
	hz.POST("/update-idl", func(ctx context.Context, c *app.RequestContext) {
		// dynamic update
		type UpdateIDLRequest struct {
			Content string
		}

		var err error
		var req UpdateIDLRequest
		err = c.BindAndValidate(&req)
		content := req.Content
		fmt.Println("更新idl为" + content)
		err = p.UpdateIDL(content, nil)
		if err != nil {
			fmt.Println(err)
			c.String(400, err.Error())
		}
		c.JSON(200, "updateok")
	})
	hz.Any("/*wildcard", Handler)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

func Handler(ctx context.Context, c *app.RequestContext) {
	//fmt.Println("err打印")
	req, err := adaptor.GetCompatRequest(&c.Request)
	if err != nil {
		fmt.Println(err)
		return
	}
	// You may build more logic on req
	fmt.Println(req.URL.String())
	// caution: don't pass in c.GetResponse() as it return a copy of response
	rw := adaptor.GetCompatResponseWriter(&c.Response)
	//fmt.Println("err打印2")

	handler2(rw, req)

}
func handler2(w http.ResponseWriter, req *http.Request) {

	// 将 HTTP 请求转换为泛化请求
	fmt.Println(req)
	customReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 发起泛化调用
	resp, err := cli.GenericCall(context.Background(), "", customReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 将泛化响应写入 HTTP 响应
	genericResp := resp.(*generic.HTTPResponse)
	genericResp.Write(w)
}
