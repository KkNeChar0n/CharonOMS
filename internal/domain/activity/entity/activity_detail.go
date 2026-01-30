package entity

// ActivityDetail 活动详情实体（满折规则）
type ActivityDetail struct {
	ID              int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ActivityID      int     `gorm:"column:activity_id;not null" json:"activity_id"`
	ThresholdAmount float64 `gorm:"column:threshold_amount;type:decimal(10,2);not null" json:"threshold_amount"` // 门槛金额
	DiscountValue   float64 `gorm:"column:discount_value;type:decimal(10,2);not null" json:"discount_value"`     // 折扣值
}

// TableName 指定表名
func (ActivityDetail) TableName() string {
	return "activity_detail"
}

// Validate 验证活动详情
func (d *ActivityDetail) Validate(activityType int) error {
	if d.ThresholdAmount <= 0 {
		return ErrInvalidThresholdAmount
	}

	if d.DiscountValue <= 0 {
		return ErrInvalidDiscountValue
	}

	// 满折类型：折扣值为百分比形式，必须在 0-100 之间
	// 例如：90表示9折（顾客付90%），7表示0.7折（顾客付7%）
	if activityType == 2 && d.DiscountValue > 100 {
		return ErrInvalidDiscountValue
	}

	return nil
}
