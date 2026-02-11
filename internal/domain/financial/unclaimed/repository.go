package unclaimed

import "time"

// UnclaimedRepository 常规待认领款项仓储接口
type UnclaimedRepository interface {
	// Create 创建待认领记录
	Create(unclaimed *Unclaimed) error

	// GetByID 根据ID获取待认领记录
	GetByID(id int) (*Unclaimed, error)

	// List 获取待认领列表
	List(filters map[string]interface{}) ([]*Unclaimed, error)

	// Update 更新待认领记录
	Update(unclaimed *Unclaimed) error

	// Delete 删除待认领记录
	Delete(id int) error

	// FindMatchingPayment 查找匹配的未核验收款记录
	// 用于导入Excel时自动匹配
	FindMatchingPayment(merchantOrder string, paymentMethod int, paymentAmount float64, payeeEntity int) (int, error)

	// UpdatePaymentStatus 更新收款记录状态为已支付
	UpdatePaymentStatus(paymentID int, arrivalTime *time.Time) error
}
