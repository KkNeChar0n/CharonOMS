package taobao

// TaobaoPaymentRepository 淘宝收款仓储接口
type TaobaoPaymentRepository interface {
	// Create 创建淘宝收款记录
	Create(payment *TaobaoPayment) error

	// GetByID 根据ID获取淘宝收款记录
	GetByID(id int) (*TaobaoPayment, error)

	// List 获取淘宝收款列表
	List(filters map[string]interface{}) ([]*TaobaoPayment, error)

	// Update 更新淘宝收款记录
	Update(payment *TaobaoPayment) error

	// Delete 删除淘宝收款记录
	Delete(id int) error

	// GetTotalPaid 获取订单的淘宝收款总额（仅统计已到账状态）
	GetTotalPaid(orderID int) (float64, error)

	// ListUnclaimed 获取淘宝待认领列表（状态10和20）
	ListUnclaimed(filters map[string]interface{}) ([]*TaobaoPayment, error)

	// FindByMerchantOrderAndAmount 根据商户订单号和金额查找记录（用于导入匹配）
	FindByMerchantOrderAndAmount(merchantOrder string, amount float64) (*TaobaoPayment, error)
}
