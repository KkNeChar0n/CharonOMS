package coach

// CoachListResponse 教练列表响应
type CoachListResponse struct {
	Coaches []map[string]interface{} `json:"coaches"`
}

// ActiveCoachResponse 启用教练列表响应
type ActiveCoachResponse struct {
	Coaches []map[string]interface{} `json:"coaches"`
}

// CreateCoachRequest 创建教练请求
type CreateCoachRequest struct {
	CoachName  string `json:"coach_name" binding:"required"`
	SexID      int    `json:"sex_id" binding:"required"`
	SubjectID  int    `json:"subject_id" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	StudentIDs []int  `json:"student_ids"`
}

// UpdateCoachRequest 更新教练请求
type UpdateCoachRequest struct {
	CoachName string `json:"coach_name" binding:"required"`
	SexID     int    `json:"sex_id" binding:"required"`
	SubjectID int    `json:"subject_id" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

// UpdateCoachStatusRequest 更新教练状态请求
type UpdateCoachStatusRequest struct {
	Status *int `json:"status" binding:"required,oneof=0 1"`
}

// CreateCoachResponse 创建教练响应
type CreateCoachResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// MessageResponse 通用消息响应
type MessageResponse struct {
	Message string `json:"message"`
}
