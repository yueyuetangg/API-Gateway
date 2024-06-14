package main

import (
	"context"
	student "demo_student_client/kitex_gen/demo"
)

// StudentServiceImpl implements the last service interface defined in the IDL.
type StudentServiceImpl struct{}

// Register implements the StudentServiceImpl interface.
func (s *StudentServiceImpl) Register(ctx context.Context, student *student.Student) (resp *student.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Query implements the StudentServiceImpl interface.
func (s *StudentServiceImpl) Query(ctx context.Context, req *student.QueryReq) (resp *student.Student, err error) {
	// TODO: Your code here...
	return
}
