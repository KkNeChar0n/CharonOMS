package entity

import "time"

// Student 学生实体
type Student struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	SexID     int       `gorm:"column:sex_id;not null" json:"sex_id"`
	GradeID   int       `gorm:"column:grade_id;not null" json:"grade_id"`
	Phone     string    `gorm:"column:phone;type:varchar(20);not null" json:"phone"`
	Status    int       `gorm:"column:status;default:0" json:"status"` // 0=启用，1=禁用
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Student) TableName() string {
	return "student"
}

// StudentCoach 学生与教练的关联实体
type StudentCoach struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentID int       `gorm:"column:student_id;not null" json:"student_id"`
	CoachID   int       `gorm:"column:coach_id;not null" json:"coach_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (StudentCoach) TableName() string {
	return "student_coach"
}
