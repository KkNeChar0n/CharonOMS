package entity

// AttributeValue 属性值实体
type AttributeValue struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	AttributeID int    `gorm:"column:attributeid;not null" json:"attributeid"`
}

// TableName 指定表名
func (AttributeValue) TableName() string {
	return "attribute_value"
}
