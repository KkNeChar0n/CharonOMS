package student

// CreateStudentRequest 创建学生请求
type CreateStudentRequest struct {
	Name     string `json:"student_name" binding:"required"` // 前端使用student_name
	SexID    int    `json:"sex_id" binding:"required"`
	GradeID  int    `json:"grade_id" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	CoachIDs []int  `json:"coach_ids"` // 可选，关联的教练ID列表
}

// UpdateStudentRequest 更新学生请求
type UpdateStudentRequest struct {
	Name    string `json:"student_name" binding:"required"` // 前端使用student_name
	SexID   int    `json:"sex_id" binding:"required"`
	GradeID int    `json:"grade_id" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
}

// UpdateStudentStatusRequest 更新学生状态请求
type UpdateStudentStatusRequest struct {
	Status *int `json:"status" binding:"required,oneof=0 1"` // 使用指针类型，可以区分0和未提供的值
}

// StudentListResponse 学生列表响应
type StudentListResponse struct {
	Students []map[string]interface{} `json:"students"`
}

// ActiveStudentResponse 启用学生列表响应
type ActiveStudentResponse struct {
	Students []map[string]interface{} `json:"students"`
}

// StudentDetailDTO 学生详情DTO
type StudentDetailDTO struct {
	ID          int    `json:"id"`
	StudentName string `json:"student_name"`
	SexID       int    `json:"sex_id"`
	Sex         string `json:"sex"`
	GradeID     int    `json:"grade_id"`
	Grade       string `json:"grade"`
	Phone       string `json:"phone"`
	Status      int    `json:"status"`
	CoachNames  string `json:"coach_names"`
}
