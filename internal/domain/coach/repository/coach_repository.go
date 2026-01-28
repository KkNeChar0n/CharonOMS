package repository

// CoachRepository 教练仓储接口
type CoachRepository interface {
	// GetCoachList 获取教练列表（含关联信息）
	GetCoachList() ([]map[string]interface{}, error)

	// GetActiveCoaches 获取启用状态的教练列表
	GetActiveCoaches() ([]map[string]interface{}, error)

	// GetCoachByID 根据ID获取教练详情
	GetCoachByID(id int) (map[string]interface{}, error)

	// CreateCoach 创建教练
	CreateCoach(name string, sexID, subjectID int, phone string, status int) (int, error)

	// UpdateCoach 更新教练信息
	UpdateCoach(id int, name string, sexID, subjectID int, phone string) error

	// UpdateCoachStatus 更新教练状态
	UpdateCoachStatus(id int, status int) error

	// DeleteCoach 删除教练（级联删除关联）
	DeleteCoach(id int) error

	// AddCoachStudents 添加教练学生关联
	AddCoachStudents(coachID int, studentIDs []int) error

	// RemoveCoachStudents 删除教练学生关联
	RemoveCoachStudents(coachID int, studentIDs []int) error

	// CoachExists 检查教练是否存在
	CoachExists(id int) (bool, error)
}
