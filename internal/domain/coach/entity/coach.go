package entity

import "time"

// Coach 教练实体
type Coach struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	SexID     int       `gorm:"column:sex_id;not null" json:"sex_id"`
	SubjectID int       `gorm:"column:subject_id;not null" json:"subject_id"`
	Phone     string    `gorm:"column:phone;type:varchar(20);not null" json:"phone"`
	Status    int       `gorm:"column:status;default:0" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Coach) TableName() string {
	return "coach"
}

// StudentCoach 学生教练关联实体
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
