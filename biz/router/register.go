// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	demo_student "demo_student_client/biz/router/demo/student"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	demo_student.Register(r)
}