package student

import (
	"charonoms/internal/domain/student/entity"
	"charonoms/internal/domain/student/repository"
	"errors"

	"gorm.io/gorm"
)

// StudentRepositoryImpl 学生仓储实现
type StudentRepositoryImpl struct {
	db *gorm.DB
}

// NewStudentRepository 创建学生仓储实例
func NewStudentRepository(db *gorm.DB) repository.StudentRepository {
	return &StudentRepositoryImpl{db: db}
}

// GetStudentList 获取学生列表（含关联信息）
func (r *StudentRepositoryImpl) GetStudentList() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Table("student s").
		Select(`s.id,
				s.name as student_name,
				s.sex_id,
				sex.name as sex,
				s.grade_id,
				grade.name as grade,
				s.phone,
				s.status,
				GROUP_CONCAT(DISTINCT coach.name ORDER BY coach.name SEPARATOR ', ') as coach_names`).
		Joins("LEFT JOIN sex ON s.sex_id = sex.id").
		Joins("LEFT JOIN grade ON s.grade_id = grade.id").
		Joins("LEFT JOIN student_coach sc ON s.id = sc.student_id").
		Joins("LEFT JOIN coach ON sc.coach_id = coach.id").
		Group("s.id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 确保返回空数组而不是nil
	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetActiveStudents 获取启用状态的学生列表
func (r *StudentRepositoryImpl) GetActiveStudents() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Table("student").
		Select("id, name as student_name").
		Where("status = ?", 0).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 确保返回空数组而不是nil
	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetStudentByID 根据ID查询学生
func (r *StudentRepositoryImpl) GetStudentByID(id int) (*entity.Student, error) {
	var student entity.Student
	err := r.db.First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// CreateStudent 创建学生
func (r *StudentRepositoryImpl) CreateStudent(student *entity.Student) error {
	return r.db.Create(student).Error
}

// UpdateStudent 更新学生信息
func (r *StudentRepositoryImpl) UpdateStudent(student *entity.Student) error {
	return r.db.Model(&entity.Student{}).
		Where("id = ?", student.ID).
		Updates(map[string]interface{}{
			"name":     student.Name,
			"sex_id":   student.SexID,
			"grade_id": student.GradeID,
			"phone":    student.Phone,
		}).Error
}

// UpdateStudentStatus 更新学生状态
func (r *StudentRepositoryImpl) UpdateStudentStatus(id int, status int) error {
	return r.db.Model(&entity.Student{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// DeleteStudent 删除学生（级联删除关联）
func (r *StudentRepositoryImpl) DeleteStudent(id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除学生教练关联
		if err := tx.Where("student_id = ?", id).Delete(&entity.StudentCoach{}).Error; err != nil {
			return err
		}

		// 删除学生记录
		if err := tx.Delete(&entity.Student{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckStudentHasOrders 检查学生是否有订单
func (r *StudentRepositoryImpl) CheckStudentHasOrders(studentID int) (bool, error) {
	var count int64

	// 检查orders表是否存在
	var tableExists int
	err := r.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'orders'").Scan(&tableExists).Error
	if err != nil {
		return false, err
	}

	// 如果订单表不存在，返回false
	if tableExists == 0 {
		return false, nil
	}

	// 查询学生是否有订单
	err = r.db.Table("orders").Where("student_id = ?", studentID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// AddStudentCoaches 添加学生教练关联
func (r *StudentRepositoryImpl) AddStudentCoaches(studentID int, coachIDs []int) error {
	if len(coachIDs) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, coachID := range coachIDs {
			sc := &entity.StudentCoach{
				StudentID: studentID,
				CoachID:   coachID,
			}
			// 使用FirstOrCreate避免重复插入
			if err := tx.Where("student_id = ? AND coach_id = ?", studentID, coachID).
				FirstOrCreate(sc).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveStudentCoaches 删除学生的所有教练关联
func (r *StudentRepositoryImpl) RemoveStudentCoaches(studentID int) error {
	return r.db.Where("student_id = ?", studentID).Delete(&entity.StudentCoach{}).Error
}

// ValidateStudentExists 验证学生是否存在
func (r *StudentRepositoryImpl) ValidateStudentExists(id int) error {
	var count int64
	err := r.db.Model(&entity.Student{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("学生不存在")
	}
	return nil
}
