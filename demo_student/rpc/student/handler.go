package main

import (
	"context"
	demo "demo_student/kitex_gen/demo"
	"errors"
	"log"
	"strconv"
	"sync"
)

// StudentServiceImpl implements the last service interface defined in the IDL.
type StudentServiceImpl struct{}

var students sync.Map

// Register implements the StudentServiceImpl interface.

// Register 实现了 StudentServiceImpl 接口。
func (s *StudentServiceImpl) Register(ctx context.Context, student *demo.Student) (resp *demo.RegisterResp, err error) {
	resp = &demo.RegisterResp{}
	log.Println("姓名" + student.Name)
	// 将学生数据添加到 students 的 sync.Map 中
	students.Store(student.Id, student)
	// 设置响应
	//resp.id = student.Id
	resp.Success = true
	resp.Message = "学生注册成功" + strconv.Itoa(int(student.Id))

	return resp, nil
}

// Query implements the StudentServiceImpl interface.
// Query 实现了 StudentServiceImpl 接口。
func (s *StudentServiceImpl) Query(ctx context.Context, req *demo.QueryReq) (resp *demo.Student, err error) {
	log.Println("查询id" + strconv.Itoa(int(req.Id)))
	resp = &demo.Student{}

	// 从 students 的 sync.Map 中根据 req.Id 查找学生数据
	value, ok := students.Load(req.Id)
	if !ok {
		log.Println("查询id" + strconv.Itoa(int(req.Id)) + "empyty")

		// 如果没查到，则返回空
		return &demo.Student{}, nil
	}

	// 将查到的学生数据转换为 demo.Student 类型并返回
	student, ok := value.(*demo.Student)
	if !ok {
		return nil, errors.New("类型断言失败")
	}

	return student, nil
}
