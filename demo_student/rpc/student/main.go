package main

import (
	demo "demo_student/kitex_gen/demo/studentservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	etcd "github.com/kitex-contrib/registry-etcd"
	"time"

	//"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	//etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	svr := demo.NewServer(new(StudentServiceImpl),
		server.WithServiceAddr(addr),
		server.WithReadWriteTimeout(1000*time.Second),
		server.WithExitWaitTime(1000*time.Second),

		//指定 Registry 与服务基本信息
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "demo",
				Tags: map[string]string{
					"Cluster": "student",
				},
			},
		),
	)
	//svr := demo.NewServer(new(StudentServiceImpl))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
