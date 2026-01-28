package coach

import (
	"fmt"

	"charonoms/internal/domain/coach/entity"
	"charonoms/internal/domain/coach/repository"

	"gorm.io/gorm"
)

// CoachRepositoryImpl 教练仓储实现
type CoachRepositoryImpl struct {
	db *gorm.DB
}

// NewCoachRepository 创建教练仓储实例
func NewCoachRepository(db *gorm.DB) repository.CoachRepository {
	return &CoachRepositoryImpl{db: db}
}

// GetCoachList 获取教练列表（含关联信息）
func (r *CoachRepositoryImpl) GetCoachList() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			c.id,
			c.name as coach_name,
			c.sex_id,
			s.name as sex,
			c.subject_id,
			sub.subject as subject,
			c.phone,
			c.status
		FROM coach c
		LEFT JOIN sex s ON c.sex_id = s.id
		LEFT JOIN subject sub ON c.subject_id = sub.id
		ORDER BY c.id DESC
	`

	fmt.Println("DEBUG: Executing coach list query with sub.subject")

	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get coach list: %w", err)
	}

	// 确保返回空数组而不是nil
	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetActiveCoaches 获取启用状态的教练列表
func (r *CoachRepositoryImpl) GetActiveCoaches() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			c.id,
			c.name as coach_name
		FROM coach c
		WHERE c.status = 0
		ORDER BY c.id DESC
	`

	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get active coaches: %w", err)
	}

	// 确保返回空数组而不是nil
	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetCoachByID 根据ID获取教练详情
func (r *CoachRepositoryImpl) GetCoachByID(id int) (map[string]interface{}, error) {
	var result map[string]interface{}

	query := `
		SELECT
			c.id,
			c.name as coach_name,
			c.sex_id,
			s.name as sex,
			c.subject_id,
			sub.subject as subject,
			c.phone,
			c.status
		FROM coach c
		LEFT JOIN sex s ON c.sex_id = s.id
		LEFT JOIN subject sub ON c.subject_id = sub.id
		WHERE c.id = ?
	`

	if err := r.db.Raw(query, id).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get coach by id: %w", err)
	}

	return result, nil
}

// CreateCoach 创建教练
func (r *CoachRepositoryImpl) CreateCoach(name string, sexID, subjectID int, phone string, status int) (int, error) {
	coach := &entity.Coach{
		Name:      name,
		SexID:     sexID,
		SubjectID: subjectID,
		Phone:     phone,
		Status:    status,
	}

	if err := r.db.Create(coach).Error; err != nil {
		return 0, fmt.Errorf("failed to create coach: %w", err)
	}

	return coach.ID, nil
}

// UpdateCoach 更新教练信息
func (r *CoachRepositoryImpl) UpdateCoach(id int, name string, sexID, subjectID int, phone string) error {
	updates := map[string]interface{}{
		"name":       name,
		"sex_id":     sexID,
		"subject_id": subjectID,
		"phone":      phone,
	}

	result := r.db.Model(&entity.Coach{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update coach: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("coach not found")
	}

	return nil
}

// UpdateCoachStatus 更新教练状态
func (r *CoachRepositoryImpl) UpdateCoachStatus(id int, status int) error {
	result := r.db.Model(&entity.Coach{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update coach status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("coach not found")
	}

	return nil
}

// DeleteCoach 删除教练（级联删除关联）
func (r *CoachRepositoryImpl) DeleteCoach(id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除学生教练关联
		if err := tx.Where("coach_id = ?", id).Delete(&entity.StudentCoach{}).Error; err != nil {
			return fmt.Errorf("failed to delete student-coach associations: %w", err)
		}

		// 删除教练
		result := tx.Delete(&entity.Coach{}, id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete coach: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("coach not found")
		}

		return nil
	})
}

// AddCoachStudents 添加教练学生关联
func (r *CoachRepositoryImpl) AddCoachStudents(coachID int, studentIDs []int) error {
	if len(studentIDs) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, studentID := range studentIDs {
			sc := &entity.StudentCoach{
				StudentID: studentID,
				CoachID:   coachID,
			}
			// 使用FirstOrCreate避免重复插入
			if err := tx.Where("student_id = ? AND coach_id = ?", studentID, coachID).
				FirstOrCreate(sc).Error; err != nil {
				return fmt.Errorf("failed to add coach-student association: %w", err)
			}
		}
		return nil
	})
}

// RemoveCoachStudents 删除教练学生关联
func (r *CoachRepositoryImpl) RemoveCoachStudents(coachID int, studentIDs []int) error {
	if len(studentIDs) == 0 {
		return nil
	}

	return r.db.Where("coach_id = ? AND student_id IN ?", coachID, studentIDs).
		Delete(&entity.StudentCoach{}).Error
}

// CoachExists 检查教练是否存在
func (r *CoachRepositoryImpl) CoachExists(id int) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Coach{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check coach existence: %w", err)
	}
	return count > 0, nil
}
