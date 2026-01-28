package goods

// GoodsListResponse 商品列表响应
type GoodsListResponse struct {
	Goods []map[string]interface{} `json:"goods"`
}

// CreateGoodsRequest 创建商品请求
// 前端通过Vue select/input发送的数字字段为字符串，使用 interface{} 接收后在service层转换
type CreateGoodsRequest struct {
	Name              string      `json:"name"`
	BrandID           interface{} `json:"brandid"`
	ClassifyID        interface{} `json:"classifyid"`
	IsGroup           interface{} `json:"isgroup"`
	Price             interface{} `json:"price"`
	AttributeValueIDs []int       `json:"attributevalue_ids"`
	IncludedGoodsIDs  []int       `json:"included_goods_ids"`
}

// UpdateGoodsRequest 更新商品请求
type UpdateGoodsRequest struct {
	Name              string      `json:"name"`
	BrandID           interface{} `json:"brandid"`
	ClassifyID        interface{} `json:"classifyid"`
	IsGroup           interface{} `json:"isgroup"`
	Price             interface{} `json:"price"`
	AttributeValueIDs []int       `json:"attributevalue_ids"`
	IncludedGoodsIDs  []int       `json:"included_goods_ids"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status interface{} `json:"status"`
}

// CreateGoodsResponse 创建商品响应
type CreateGoodsResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// MessageResponse 通用消息响应
type MessageResponse struct {
	Message string `json:"message"`
}
