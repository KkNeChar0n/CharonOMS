package goods

import (
	"fmt"

	"charonoms/internal/domain/goods/entity"
	"charonoms/internal/domain/goods/repository"

	"gorm.io/gorm"
)

// GoodsRepositoryImpl 商品仓储实现
type GoodsRepositoryImpl struct {
	db *gorm.DB
}

// NewGoodsRepository 创建商品仓储实例
func NewGoodsRepository(db *gorm.DB) repository.GoodsRepository {
	return &GoodsRepositoryImpl{db: db}
}

// GetList 获取商品列表（含品牌、分类、属性信息）
func (r *GoodsRepositoryImpl) GetList(classifyID *int, status *int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			g.id,
			g.name,
			g.brandid,
			b.name as brand_name,
			g.classifyid,
			c.name as classify_name,
			g.isgroup,
			g.price,
			g.status,
			g.create_time,
			g.update_time,
			GROUP_CONCAT(
				CONCAT(a.name, ':', av.name)
				ORDER BY a.name SEPARATOR ','
			) as attributes
		FROM goods g
		LEFT JOIN brand b ON g.brandid = b.id
		LEFT JOIN classify c ON g.classifyid = c.id
		LEFT JOIN goods_attributevalue gav ON g.id = gav.goodsid
		LEFT JOIN attribute_value av ON gav.attributevalueid = av.id
		LEFT JOIN attribute a ON av.attributeid = a.id
		WHERE 1=1
	`

	var args []interface{}

	// 添加分类过滤
	if classifyID != nil {
		query += " AND g.classifyid = ?"
		args = append(args, *classifyID)
	}

	// 添加状态过滤
	if status != nil {
		query += " AND g.status = ?"
		args = append(args, *status)
	}

	query += `
		GROUP BY g.id, g.name, g.brandid, b.name, g.classifyid, c.name, g.isgroup, g.price, g.status, g.create_time, g.update_time
		ORDER BY g.create_time DESC
	`

	if err := r.db.Raw(query, args...).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get goods list: %w", err)
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetByID 根据ID获取商品详情（含属性值ID数组和包含的商品ID数组）
func (r *GoodsRepositoryImpl) GetByID(id int) (map[string]interface{}, error) {
	var result map[string]interface{}

	// 获取商品基本信息
	query := `
		SELECT
			g.id,
			g.name,
			g.brandid,
			b.name as brand_name,
			g.classifyid,
			c.name as classify_name,
			g.isgroup,
			g.price,
			g.status,
			g.create_time,
			g.update_time
		FROM goods g
		LEFT JOIN brand b ON g.brandid = b.id
		LEFT JOIN classify c ON g.classifyid = c.id
		WHERE g.id = ?
	`

	if err := r.db.Raw(query, id).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get goods by id: %w", err)
	}

	if result == nil || len(result) == 0 {
		return nil, nil
	}

	// 获取属性值ID数组
	var attributeValueIDs []int
	attrQuery := `
		SELECT attributevalueid
		FROM goods_attributevalue
		WHERE goodsid = ?
		ORDER BY attributevalueid
	`
	if err := r.db.Raw(attrQuery, id).Scan(&attributeValueIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get attribute value ids: %w", err)
	}
	if attributeValueIDs == nil {
		attributeValueIDs = []int{}
	}
	result["attributevalue_ids"] = attributeValueIDs

	// 获取包含的商品ID数组（如果是组合商品）
	var includedGoodsIDs []int
	goodsQuery := `
		SELECT goodsid
		FROM goods_goods
		WHERE parentsid = ?
		ORDER BY goodsid
	`
	if err := r.db.Raw(goodsQuery, id).Scan(&includedGoodsIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get included goods ids: %w", err)
	}
	if includedGoodsIDs == nil {
		includedGoodsIDs = []int{}
	}
	result["included_goods_ids"] = includedGoodsIDs

	return result, nil
}

// GetActiveForOrder 获取可用于下单的商品列表（status=0，含total_price）
func (r *GoodsRepositoryImpl) GetActiveForOrder() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			g.id,
			g.name,
			g.brandid,
			b.name as brand_name,
			g.classifyid,
			c.name as classify_name,
			g.isgroup,
			g.price,
			CASE
				WHEN g.isgroup = 1 THEN g.price
				ELSE COALESCE((
					SELECT SUM(child.price)
					FROM goods_goods gg
					INNER JOIN goods child ON gg.goodsid = child.id
					WHERE gg.parentsid = g.id
				), 0)
			END as total_price,
			g.status,
			g.create_time
		FROM goods g
		LEFT JOIN brand b ON g.brandid = b.id
		LEFT JOIN classify c ON g.classifyid = c.id
		WHERE g.status = 0
		ORDER BY g.create_time DESC
	`

	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get active goods for order: %w", err)
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetAvailableForCombo 获取可用于组合的单品商品列表（isgroup=1, status=0），excludeID>0时排除该商品
func (r *GoodsRepositoryImpl) GetAvailableForCombo(excludeID int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			g.id,
			g.name,
			g.brandid,
			b.name as brand_name,
			g.classifyid,
			c.name as classify_name,
			g.price,
			g.status
		FROM goods g
		LEFT JOIN brand b ON g.brandid = b.id
		LEFT JOIN classify c ON g.classifyid = c.id
		WHERE g.isgroup = 1 AND g.status = 0
	`

	if excludeID > 0 {
		query += fmt.Sprintf(" AND g.id != %d", excludeID)
	}

	query += " ORDER BY g.name ASC"

	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get available goods for combo: %w", err)
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetIncludedGoods 根据父商品ID获取包含的子商品列表
func (r *GoodsRepositoryImpl) GetIncludedGoods(parentsID int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			g.id,
			g.name,
			g.brandid,
			b.name as brand_name,
			g.classifyid,
			c.name as classify_name,
			g.price,
			g.status
		FROM goods_goods gg
		INNER JOIN goods g ON gg.goodsid = g.id
		LEFT JOIN brand b ON g.brandid = b.id
		LEFT JOIN classify c ON g.classifyid = c.id
		WHERE gg.parentsid = ?
		ORDER BY g.name ASC
	`

	if err := r.db.Raw(query, parentsID).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get included goods: %w", err)
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	return results, nil
}

// GetTotalPrice 计算商品总价（单商品返回price，组合商品返回子商品价格之和）
func (r *GoodsRepositoryImpl) GetTotalPrice(id int) (float64, error) {
	var totalPrice float64

	query := `
		SELECT
			CASE
				WHEN g.isgroup = 1 THEN g.price
				ELSE COALESCE((
					SELECT SUM(child.price)
					FROM goods_goods gg
					INNER JOIN goods child ON gg.goodsid = child.id
					WHERE gg.parentsid = g.id
				), 0)
			END as total_price
		FROM goods g
		WHERE g.id = ?
	`

	if err := r.db.Raw(query, id).Scan(&totalPrice).Error; err != nil {
		return 0, fmt.Errorf("failed to get total price: %w", err)
	}

	return totalPrice, nil
}

// Create 创建商品
func (r *GoodsRepositoryImpl) Create(name string, brandID int, classifyID int, isGroup int, price float64, attributeValueIDs []int, includedGoodsIDs []int) (int, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// 创建商品
	goods := &entity.Goods{
		Name:       name,
		BrandID:    brandID,
		ClassifyID: classifyID,
		IsGroup:    &isGroup,
		Price:      price,
		Status:     0,
	}

	if err := tx.Create(goods).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create goods: %w", err)
	}

	goodsID := goods.ID

	// 创建属性值关联
	if len(attributeValueIDs) > 0 {
		for _, attrValueID := range attributeValueIDs {
			goodsAttrValue := &entity.GoodsAttributeValue{
				GoodsID:          goodsID,
				AttributeValueID: attrValueID,
			}
			if err := tx.Create(goodsAttrValue).Error; err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to create goods attribute value: %w", err)
			}
		}
	}

	// 创建商品组合关联（如果是组合商品）
	if isGroup == 0 && len(includedGoodsIDs) > 0 {
		for _, includedGoodsID := range includedGoodsIDs {
			goodsGoods := &entity.GoodsGoods{
				GoodsID:   includedGoodsID,
				ParentsID: goodsID,
			}
			if err := tx.Create(goodsGoods).Error; err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to create goods goods: %w", err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return goodsID, nil
}

// Update 更新商品
func (r *GoodsRepositoryImpl) Update(id int, name string, brandID int, classifyID int, isGroup int, price float64, attributeValueIDs []int, includedGoodsIDs []int) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// 更新商品基本信息
	updates := map[string]interface{}{
		"name":       name,
		"brandid":    brandID,
		"classifyid": classifyID,
		"isgroup":    isGroup,
		"price":      price,
	}

	result := tx.Model(&entity.Goods{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update goods: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("goods not found")
	}

	// 删除旧的属性值关联
	if err := tx.Where("goodsid = ?", id).Delete(&entity.GoodsAttributeValue{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old attribute values: %w", err)
	}

	// 创建新的属性值关联
	if len(attributeValueIDs) > 0 {
		for _, attrValueID := range attributeValueIDs {
			goodsAttrValue := &entity.GoodsAttributeValue{
				GoodsID:          id,
				AttributeValueID: attrValueID,
			}
			if err := tx.Create(goodsAttrValue).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create goods attribute value: %w", err)
			}
		}
	}

	// 删除旧的商品组合关联
	if err := tx.Where("parentsid = ?", id).Delete(&entity.GoodsGoods{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old goods goods: %w", err)
	}

	// 创建新的商品组合关联（如果是组合商品）
	if isGroup == 0 && len(includedGoodsIDs) > 0 {
		for _, includedGoodsID := range includedGoodsIDs {
			goodsGoods := &entity.GoodsGoods{
				GoodsID:   includedGoodsID,
				ParentsID: id,
			}
			if err := tx.Create(goodsGoods).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create goods goods: %w", err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateStatus 更新商品状态
func (r *GoodsRepositoryImpl) UpdateStatus(id int, status int) error {
	result := r.db.Model(&entity.Goods{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update goods status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("goods not found")
	}

	return nil
}
