package repository

import (
	"charonoms/internal/domain/student/entity"
)

// StudentRepository 学生仓储接口
type StudentRepository interface {
	// GetStudentList 获取学生列表（含关联信息）
	GetStudentList() ([]map[string]interface{}, error)

	// GetActiveStudents 获取启用状态的学生列表
	GetActiveStudents() ([]map[string]interface{}, error)

	// GetStudentByID 根据ID查询学生
	GetStudentByID(id int) (*entity.Student, error)

	// CreateStudent 创建学生
	CreateStudent(student *entity.Student) error

	// UpdateStudent 更新学生信息
	UpdateStudent(student *entity.Student) error

	// UpdateStudentStatus 更新学生状态
	UpdateStudentStatus(id int, status int) error

	// DeleteStudent 删除学生（级联删除关联）
	DeleteStudent(id int) error

	// CheckStudentHasOrders 检查学生是否有订单
	CheckStudentHasOrders(studentID int) (bool, error)

	// AddStudentCoaches 添加学生教练关联
	AddStudentCoaches(studentID int, coachIDs []int) error

	// RemoveStudentCoaches 删除学生的所有教练关联
	RemoveStudentCoaches(studentID int) error
}
