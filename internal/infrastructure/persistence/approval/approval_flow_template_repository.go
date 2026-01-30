package approval

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
	"gorm.io/gorm"
)

// GormApprovalFlowTemplateRepository GORM实现的审批流模板仓储
type GormApprovalFlowTemplateRepository struct {
	db *gorm.DB
}

// NewApprovalFlowTemplateRepository 创建审批流模板仓储实例
func NewApprovalFlowTemplateRepository(db *gorm.DB) repository.ApprovalFlowTemplateRepository {
	return &GormApprovalFlowTemplateRepository{db: db}
}

// GetList 获取审批流模板列表（含关联类型名称）
func (r *GormApprovalFlowTemplateRepository) GetList(filters map[string]interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.Table("approval_flow_template t").
		Select("t.*, ft.name as flow_type_name").
		Joins("LEFT JOIN approval_flow_type ft ON t.approval_flow_type_id = ft.id")

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("t.id = ?", id)
	}
	if typeID, ok := filters["approval_flow_type_id"]; ok && typeID != "" {
		query = query.Where("t.approval_flow_type_id = ?", typeID)
	}
	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("t.name LIKE ?", "%"+name.(string)+"%")
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("t.status = ?", status)
	}

	err := query.Find(&results).Error
	return results, err
}

// GetByID 获取模板基本信息
func (r *GormApprovalFlowTemplateRepository) GetByID(id int) (*entity.ApprovalFlowTemplate, error) {
	var template entity.ApprovalFlowTemplate
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetDetailByID 获取模板完整详情（含节点、审批人员、抄送人员）
func (r *GormApprovalFlowTemplateRepository) GetDetailByID(id int) (map[string]interface{}, error) {
	// 获取模板基本信息和类型名称
	var templateInfo map[string]interface{}
	err := r.db.Table("approval_flow_template t").
		Select("t.*, ft.name as flow_type_name").
		Joins("LEFT JOIN approval_flow_type ft ON t.approval_flow_type_id = ft.id").
		Where("t.id = ?", id).
		Take(&templateInfo).Error
	if err != nil {
		return nil, err
	}

	// 获取节点列表（按sort排序）
	var nodes []map[string]interface{}
	err = r.db.Table("approval_flow_template_node").
		Where("template_id = ?", id).
		Order("sort ASC").
		Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	// 为每个节点获取审批人员
	for i := range nodes {
		// 安全地获取nodeID（处理int32和int64）
		var nodeID int
		switch v := nodes[i]["id"].(type) {
		case int32:
			nodeID = int(v)
		case int64:
			nodeID = int(v)
		case int:
			nodeID = v
		}

		var approvers []map[string]interface{}
		err = r.db.Table("approval_node_useraccount anu").
			Select("ua.id, ua.username").
			Joins("INNER JOIN useraccount ua ON anu.useraccount_id = ua.id").
			Where("anu.node_id = ?", nodeID).
			Find(&approvers).Error
		if err != nil {
			return nil, err
		}
		nodes[i]["approvers"] = approvers
	}

	// 获取抄送人员
	var copyUsers []map[string]interface{}
	err = r.db.Table("approval_copy_useraccount acu").
		Select("ua.id, ua.username").
		Joins("INNER JOIN useraccount ua ON acu.useraccount_id = ua.id").
		Where("acu.approval_flow_template_id = ?", id).
		Find(&copyUsers).Error
	if err != nil {
		return nil, err
	}

	// 组装结果
	result := map[string]interface{}{
		"template":   templateInfo,
		"nodes":      nodes,
		"copy_users": copyUsers,
	}

	return result, nil
}

// Create 创建审批流模板（事务：模板、节点、人员、抄送）
func (r *GormApprovalFlowTemplateRepository) Create(
	template *entity.ApprovalFlowTemplate,
	nodes []entity.ApprovalFlowTemplateNode,
	nodeApprovers map[int][]int,
	copyUsers []int,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建模板
		if err := tx.Create(template).Error; err != nil {
			return err
		}

		// 2. 创建节点并设置template_id
		nodeIDMap := make(map[int]int) // 索引 -> 数据库ID
		for i := range nodes {
			nodes[i].TemplateID = template.ID
			if err := tx.Create(&nodes[i]).Error; err != nil {
				return err
			}
			nodeIDMap[i] = nodes[i].ID
		}

		// 3. 创建节点审批人员
		for nodeIndex, approverIDs := range nodeApprovers {
			nodeID := nodeIDMap[nodeIndex]
			for _, approverID := range approverIDs {
				approver := entity.ApprovalNodeUserAccount{
					NodeID:        nodeID,
					UserAccountID: approverID,
				}
				if err := tx.Create(&approver).Error; err != nil {
					return err
				}
			}
		}

		// 4. 创建抄送人员
		for _, userID := range copyUsers {
			copyUser := entity.ApprovalCopyUserAccount{
				ApprovalFlowTemplateID: template.ID,
				UserAccountID:          userID,
			}
			if err := tx.Create(&copyUser).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateStatus 更新模板状态
func (r *GormApprovalFlowTemplateRepository) UpdateStatus(id int, status int8) error {
	return r.db.Model(&entity.ApprovalFlowTemplate{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// DisableSameTypeTemplates 禁用同类型的其他模板
func (r *GormApprovalFlowTemplateRepository) DisableSameTypeTemplates(typeID int, excludeID int) error {
	return r.db.Model(&entity.ApprovalFlowTemplate{}).
		Where("approval_flow_type_id = ? AND id != ?", typeID, excludeID).
		Update("status", 1).Error
}

// GetNodesByTemplateID 获取模板的所有节点（按sort排序）
func (r *GormApprovalFlowTemplateRepository) GetNodesByTemplateID(templateID int) ([]entity.ApprovalFlowTemplateNode, error) {
	var nodes []entity.ApprovalFlowTemplateNode
	err := r.db.Where("template_id = ?", templateID).Order("sort ASC").Find(&nodes).Error
	return nodes, err
}

// GetNodeApprovers 获取节点的审批人员
func (r *GormApprovalFlowTemplateRepository) GetNodeApprovers(nodeID int) ([]int, error) {
	var approvers []entity.ApprovalNodeUserAccount
	err := r.db.Where("node_id = ?", nodeID).Find(&approvers).Error
	if err != nil {
		return nil, err
	}

	approverIDs := make([]int, len(approvers))
	for i, a := range approvers {
		approverIDs[i] = a.UserAccountID
	}
	return approverIDs, nil
}

// GetCopyUsers 获取模板的抄送人员
func (r *GormApprovalFlowTemplateRepository) GetCopyUsers(templateID int) ([]int, error) {
	var copyUsers []entity.ApprovalCopyUserAccount
	err := r.db.Where("approval_flow_template_id = ?", templateID).Find(&copyUsers).Error
	if err != nil {
		return nil, err
	}

	userIDs := make([]int, len(copyUsers))
	for i, cu := range copyUsers {
		userIDs[i] = cu.UserAccountID
	}
	return userIDs, nil
}
