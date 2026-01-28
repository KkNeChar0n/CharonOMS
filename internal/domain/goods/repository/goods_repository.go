package repository

// GoodsRepository 商品仓储接口
type GoodsRepository interface {
	// GetList 获取商品列表（含品牌、分类、属性信息）
	GetList() ([]map[string]interface{}, error)

	// GetByID 根据ID获取商品详情（含属性值ID数组和包含的商品ID数组）
	GetByID(id int) (map[string]interface{}, error)

	// GetActiveForOrder 获取可用于下单的商品列表（status=0，含total_price）
	GetActiveForOrder() ([]map[string]interface{}, error)

	// GetAvailableForCombo 获取可用于组合的单品商品列表（isgroup=1, status=0），excludeID>0时排除该商品
	GetAvailableForCombo(excludeID int) ([]map[string]interface{}, error)

	// GetIncludedGoods 根据父商品ID获取包含的子商品列表
	GetIncludedGoods(parentsID int) ([]map[string]interface{}, error)

	// GetTotalPrice 计算商品总价（单商品返回price，组合商品返回子商品价格之和）
	GetTotalPrice(id int) (float64, error)

	// Create 创建商品
	Create(name string, brandID int, classifyID int, isGroup int, price float64, attributeValueIDs []int, includedGoodsIDs []int) (int, error)

	// Update 更新商品
	Update(id int, name string, brandID int, classifyID int, isGroup int, price float64, attributeValueIDs []int, includedGoodsIDs []int) error

	// UpdateStatus 更新商品状态
	UpdateStatus(id int, status int) error
}
