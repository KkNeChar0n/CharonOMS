package approval

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// GormApprovalFlowManagementRepository GORM实现的审批流实例仓储
type GormApprovalFlowManagementRepository struct {
	db *gorm.DB
}

// NewApprovalFlowManagementRepository 创建审批流实例仓储实例
func NewApprovalFlowManagementRepository(db *gorm.DB) repository.ApprovalFlowManagementRepository {
	return &GormApprovalFlowManagementRepository{db: db}
}

// GetInitiatedFlows 获取用户发起的审批流
func (r *GormApprovalFlowManagementRepository) GetInitiatedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.Table("approval_flow_management fm").
		Select(`
			fm.id,
			fm.approval_flow_template_id,
			fm.approval_flow_type_id,
			fm.step,
			fm.create_user,
			fm.create_time,
			fm.status,
			fm.complete_time,
			ft.name as template_name,
			ftype.name as flow_type_name
		`).
		Joins("LEFT JOIN approval_flow_template ft ON fm.approval_flow_template_id = ft.id").
		Joins("LEFT JOIN approval_flow_type ftype ON fm.approval_flow_type_id = ftype.id").
		Where("fm.create_user = ?", userID)

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("fm.id = ?", id)
	}
	if flowTypeID, ok := filters["approval_flow_type_id"]; ok && flowTypeID != "" {
		query = query.Where("fm.approval_flow_type_id = ?", flowTypeID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("fm.status = ?", status)
	}
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		query = query.Where("fm.create_time >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		query = query.Where("fm.create_time <= ?", endTime)
	}

	err := query.Order("fm.create_time DESC").Find(&results).Error
	return results, err
}

// GetPendingFlows 获取待用户审批的任务
func (r *GormApprovalFlowManagementRepository) GetPendingFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.Table("approval_node_case_user ncu").
		Select(`
			ncu.id as id,
			fm.id as approval_flow_management_id,
			fm.approval_flow_template_id,
			fm.approval_flow_type_id,
			fm.step,
			fm.create_user,
			fm.create_time,
			fm.status,
			ft.name as template_name,
			ftype.name as approval_flow_type_name,
			nc.id as node_case_id,
			nc.type as node_type,
			nc.sort as node_sort,
			ncu.id as user_case_id,
			ncu.create_time as assigned_time
		`).
		Joins("INNER JOIN approval_node_case nc ON ncu.approval_node_case_id = nc.id").
		Joins("INNER JOIN approval_flow_management fm ON nc.approval_flow_management_id = fm.id").
		Joins("LEFT JOIN approval_flow_template ft ON fm.approval_flow_template_id = ft.id").
		Joins("LEFT JOIN approval_flow_type ftype ON fm.approval_flow_type_id = ftype.id").
		Where("ncu.useraccount_id = ?", userID).
		Where("ncu.result IS NULL") // 待审批

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("ncu.id = ?", id)
	}
	if approvalFlowID, ok := filters["approval_flow_id"]; ok && approvalFlowID != "" {
		query = query.Where("fm.id = ?", approvalFlowID)
	}
	if flowTypeID, ok := filters["approval_flow_type_id"]; ok && flowTypeID != "" {
		query = query.Where("fm.approval_flow_type_id = ?", flowTypeID)
	}
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		query = query.Where("ncu.create_time >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		query = query.Where("ncu.create_time <= ?", endTime)
	}

	err := query.Order("ncu.create_time DESC").Find(&results).Error
	return results, err
}

// GetCompletedFlows 获取用户已处理的审批任务
func (r *GormApprovalFlowManagementRepository) GetCompletedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.Table("approval_node_case_user ncu").
		Select(`
			ncu.id as id,
			fm.id as approval_flow_management_id,
			fm.approval_flow_template_id,
			fm.approval_flow_type_id,
			fm.step,
			fm.create_user,
			fm.create_time,
			fm.status,
			ft.name as template_name,
			ftype.name as approval_flow_type_name,
			nc.id as node_case_id,
			nc.type as node_type,
			nc.sort as node_sort,
			ncu.id as user_case_id,
			ncu.result as result,
			ncu.handle_time
		`).
		Joins("INNER JOIN approval_node_case nc ON ncu.approval_node_case_id = nc.id").
		Joins("INNER JOIN approval_flow_management fm ON nc.approval_flow_management_id = fm.id").
		Joins("LEFT JOIN approval_flow_template ft ON fm.approval_flow_template_id = ft.id").
		Joins("LEFT JOIN approval_flow_type ftype ON fm.approval_flow_type_id = ftype.id").
		Where("ncu.useraccount_id = ?", userID).
		Where("ncu.result IS NOT NULL") // 已处理

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("ncu.id = ?", id)
	}
	if approvalFlowID, ok := filters["approval_flow_id"]; ok && approvalFlowID != "" {
		query = query.Where("fm.id = ?", approvalFlowID)
	}
	if flowTypeID, ok := filters["approval_flow_type_id"]; ok && flowTypeID != "" {
		query = query.Where("fm.approval_flow_type_id = ?", flowTypeID)
	}
	if result, ok := filters["result"]; ok && result != "" {
		query = query.Where("ncu.result = ?", result)
	}
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		query = query.Where("ncu.handle_time >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		query = query.Where("ncu.handle_time <= ?", endTime)
	}

	err := query.Order("ncu.handle_time DESC").Find(&results).Error
	return results, err
}

// GetCopiedFlows 获取抄送给用户的通知
func (r *GormApprovalFlowManagementRepository) GetCopiedFlows(userID int, filters map[string]interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.Table("approval_copy_useraccount_case cuc").
		Select(`
			cuc.id as id,
			fm.id as approval_flow_id,
			fm.approval_flow_template_id,
			fm.approval_flow_type_id,
			fm.step,
			fm.create_user,
			fm.status,
			fm.complete_time,
			ft.name as template_name,
			ftype.name as approval_flow_type_name,
			cuc.id as copy_case_id,
			cuc.copy_info,
			cuc.create_time as create_time
		`).
		Joins("INNER JOIN approval_flow_management fm ON cuc.approval_flow_management_id = fm.id").
		Joins("LEFT JOIN approval_flow_template ft ON fm.approval_flow_template_id = ft.id").
		Joins("LEFT JOIN approval_flow_type ftype ON fm.approval_flow_type_id = ftype.id").
		Where("cuc.useraccount_id = ?", userID)

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("cuc.id = ?", id)
	}
	if approvalFlowID, ok := filters["approval_flow_id"]; ok && approvalFlowID != "" {
		query = query.Where("fm.id = ?", approvalFlowID)
	}
	if flowTypeID, ok := filters["approval_flow_type_id"]; ok && flowTypeID != "" {
		query = query.Where("fm.approval_flow_type_id = ?", flowTypeID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("fm.status = ?", status)
	}
	if startTime, ok := filters["start_time"]; ok && startTime != "" {
		query = query.Where("cuc.create_time >= ?", startTime)
	}
	if endTime, ok := filters["end_time"]; ok && endTime != "" {
		query = query.Where("cuc.create_time <= ?", endTime)
	}

	err := query.Order("cuc.create_time DESC").Find(&results).Error
	return results, err
}

// GetDetailByID 获取审批流详情
func (r *GormApprovalFlowManagementRepository) GetDetailByID(flowID int, userID int) (map[string]interface{}, error) {
	// 1. 获取审批流基本信息
	var flowInfo map[string]interface{}
	err := r.db.Table("approval_flow_management fm").
		Select(`
			fm.*,
			ft.name as template_name,
			ftype.name as approval_flow_type_name,
			ua.username as creator_name
		`).
		Joins("LEFT JOIN approval_flow_template ft ON fm.approval_flow_template_id = ft.id").
		Joins("LEFT JOIN approval_flow_type ftype ON fm.approval_flow_type_id = ftype.id").
		Joins("LEFT JOIN useraccount ua ON fm.create_user = ua.id").
		Where("fm.id = ?", flowID).
		Limit(1).
		Find(&flowInfo).Error
	if err != nil {
		return nil, err
	}

	// 2. 获取所有节点及其审批记录
	var nodes []map[string]interface{}
	err = r.db.Table("approval_node_case nc").
		Select(`
			nc.id as node_case_id,
			nc.node_id as template_node_id,
			nc.approval_flow_management_id,
			nc.type,
			nc.sort,
			nc.result as node_result,
			nc.create_time,
			nc.complete_time,
			tn.name as node_name
		`).
		Joins("LEFT JOIN approval_flow_template_node tn ON nc.node_id = tn.id").
		Where("nc.approval_flow_management_id = ?", flowID).
		Order("nc.sort ASC").
		Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	// 3. 为每个节点获取审批人员记录
	for i := range nodes {
		// 安全地获取节点ID，处理int32和int64两种情况
		var nodeCaseID int
		switch v := nodes[i]["node_case_id"].(type) {
		case int:
			nodeCaseID = v
		case int32:
			nodeCaseID = int(v)
		case int64:
			nodeCaseID = int(v)
		case uint:
			nodeCaseID = int(v)
		case uint32:
			nodeCaseID = int(v)
		case uint64:
			nodeCaseID = int(v)
		}

		var users []map[string]interface{}
		err = r.db.Table("approval_node_case_user ncu").
			Select(`
				ncu.id,
				ncu.approval_node_case_id,
				ncu.useraccount_id,
				ncu.result,
				ncu.create_time,
				ncu.handle_time,
				ua.username
			`).
			Joins("LEFT JOIN useraccount ua ON ncu.useraccount_id = ua.id").
			Where("ncu.approval_node_case_id = ?", nodeCaseID).
			Find(&users).Error
		if err != nil {
			return nil, err
		}
		nodes[i]["users"] = users
	}

	// 4. 获取抄送记录
	var copyRecords []map[string]interface{}
	err = r.db.Table("approval_copy_useraccount_case cuc").
		Select(`
			cuc.id,
			cuc.approval_flow_management_id,
			cuc.useraccount_id,
			cuc.copy_info,
			cuc.create_time,
			ua.username
		`).
		Joins("LEFT JOIN useraccount ua ON cuc.useraccount_id = ua.id").
		Where("cuc.approval_flow_management_id = ?", flowID).
		Find(&copyRecords).Error
	if err != nil {
		return nil, err
	}

	// 5. 获取当前用户的审批记录
	var userApproval map[string]interface{}
	err = r.db.Table("approval_node_case_user ncu").
		Select(`
			ncu.id,
			ncu.approval_node_case_id,
			nc.approval_flow_management_id,
			nc.type as node_type,
			nc.sort as node_sort,
			ncu.result,
			ncu.create_time,
			ncu.handle_time
		`).
		Joins("INNER JOIN approval_node_case nc ON ncu.approval_node_case_id = nc.id").
		Where("nc.approval_flow_management_id = ?", flowID).
		Where("ncu.useraccount_id = ?", userID).
		Order("nc.sort ASC").
		Limit(1).
		Find(&userApproval).Error
	if err != nil {
		userApproval = nil // 如果没有审批记录，设为nil
	}

	// 6. 如果是退费类型，获取退费信息
	var refundOrderInfo map[string]interface{}
	flowTypeName, _ := flowInfo["approval_flow_type_name"].(string)
	if flowTypeName == "退费" {
		// 通过审批流ID关联退费订单
		var refundOrder map[string]interface{}
		err = r.db.Table("refund_order ro").
			Select(`
				ro.id as refund_order_id,
				ro.order_id,
				ro.refund_amount,
				ro.submitter,
				ro.submit_time,
				ro.status,
				s.name as student_name,
				s.phone as student_phone,
				g.name as grade_name
			`).
			Joins("LEFT JOIN student s ON ro.student_id = s.id").
			Joins("LEFT JOIN grade g ON s.grade_id = g.id").
			Where("ro.id = (SELECT MAX(id) FROM refund_order WHERE create_time >= (SELECT create_time FROM approval_flow_management WHERE id = ?))", flowID).
			Limit(1).
			Find(&refundOrder).Error

		if err == nil && len(refundOrder) > 0 {
			var refundOrderID int
			switch v := refundOrder["refund_order_id"].(type) {
			case int:
				refundOrderID = v
			case int32:
				refundOrderID = int(v)
			case int64:
				refundOrderID = int(v)
			}

			// 获取退费明细
			var items []map[string]interface{}
			r.db.Table("refund_order_item").
				Select("goods_name, refund_amount").
				Where("refund_order_id = ?", refundOrderID).
				Find(&items)
			refundOrder["items"] = items

			// 获取退费收款分配
			var payments []map[string]interface{}
			r.db.Raw(`
				SELECT
					rp.payment_id,
					rp.payment_type,
					rp.refund_amount,
					CASE
						WHEN rp.payment_type = 0 THEN pc.payment_amount
						WHEN rp.payment_type = 1 THEN tp.payment_amount
					END as payment_amount,
					CASE
						WHEN rp.payment_type = 0 THEN
							CASE pc.payee_entity
								WHEN 0 THEN '北京'
								WHEN 1 THEN '西安'
								ELSE ''
							END
						ELSE ''
					END as payee_entity
				FROM refund_payment rp
				LEFT JOIN payment_collection pc ON rp.payment_id = pc.id AND rp.payment_type = 0
				LEFT JOIN taobao_payment tp ON rp.payment_id = tp.id AND rp.payment_type = 1
				WHERE rp.refund_order_id = ?
			`, refundOrderID).Scan(&payments)
			refundOrder["payments"] = payments

			// 获取淘宝补充信息
			var taobaoSupplement map[string]interface{}
			r.db.Table("refund_taobao_supplement").
				Where("refund_order_id = ?", refundOrderID).
				Limit(1).
				Find(&taobaoSupplement)
			refundOrder["taobao_supplement"] = taobaoSupplement

			// 获取常规补充信息
			var regularSupplements []map[string]interface{}
			r.db.Table("refund_regular_supplement").
				Where("refund_order_id = ?", refundOrderID).
				Find(&regularSupplements)
			refundOrder["regular_supplements"] = regularSupplements

			refundOrderInfo = refundOrder
		}
	}

	// 7. 检查当前用户的权限信息
	canApprove := false
	canCancel := false

	// 检查是否可以审批（是否在待审批人员中）
	var pendingCount int64
	err = r.db.Table("approval_node_case_user ncu").
		Joins("INNER JOIN approval_node_case nc ON ncu.approval_node_case_id = nc.id").
		Where("nc.approval_flow_management_id = ?", flowID).
		Where("ncu.useraccount_id = ?", userID).
		Where("ncu.result IS NULL").
		Count(&pendingCount).Error
	if err == nil && pendingCount > 0 {
		canApprove = true
	}

	// 检查是否可以撤销（是发起人且状态为待审批）
	if status, ok := flowInfo["status"].(int8); ok && status == 0 {
		if createUser, ok := flowInfo["create_user"].(int64); ok && int(createUser) == userID {
			canCancel = true
		} else if createUser, ok := flowInfo["create_user"].(int32); ok && int(createUser) == userID {
			canCancel = true
		}
	}

	// 组装结果
	result := map[string]interface{}{
		"flow_info":         flowInfo,
		"user_approval":     userApproval,
		"all_nodes":         nodes,
		"refund_order_info": refundOrderInfo,
		"copy_records":      copyRecords,
		"can_approve":       canApprove,
		"can_cancel":        canCancel,
	}

	return result, nil
}

// GetByID 根据ID查询审批流
func (r *GormApprovalFlowManagementRepository) GetByID(id int) (*entity.ApprovalFlowManagement, error) {
	var flow entity.ApprovalFlowManagement
	err := r.db.First(&flow, id).Error
	if err != nil {
		return nil, err
	}
	return &flow, nil
}

// CreateFromTemplate 从模板创建审批流实例（事务：审批流、第一个节点、审批人员）
func (r *GormApprovalFlowManagementRepository) CreateFromTemplate(templateID int, userID int) (int, error) {
	var flowID int

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取模板信息
		var template entity.ApprovalFlowTemplate
		if err := tx.First(&template, templateID).Error; err != nil {
			return fmt.Errorf("获取模板信息失败: %w", err)
		}

		// 检查模板是否启用
		if template.Status != 0 {
			return fmt.Errorf("模板已禁用，无法创建审批流")
		}

		// 2. 创建审批流实例
		flow := entity.ApprovalFlowManagement{
			ApprovalFlowTemplateID: templateID,
			ApprovalFlowTypeID:     template.ApprovalFlowTypeID,
			Step:                   0,
			CreateUser:             userID,
			Status:                 0, // 待审批
		}
		if err := tx.Create(&flow).Error; err != nil {
			return fmt.Errorf("创建审批流实例失败: %w", err)
		}
		flowID = flow.ID

		// 3. 获取第一个节点（sort最小的节点）
		var firstNode entity.ApprovalFlowTemplateNode
		err := tx.Where("template_id = ?", templateID).
			Order("sort ASC").
			First(&firstNode).Error
		if err != nil {
			return fmt.Errorf("获取第一个节点失败: %w", err)
		}

		// 4. 创建第一个节点实例
		nodeCase := entity.ApprovalNodeCase{
			NodeID:                   firstNode.ID,
			ApprovalFlowManagementID: flowID,
			Type:                     firstNode.Type,
			Sort:                     firstNode.Sort,
		}
		if err := tx.Create(&nodeCase).Error; err != nil {
			return fmt.Errorf("创建节点实例失败: %w", err)
		}

		// 5. 获取第一个节点的审批人员
		var approvers []entity.ApprovalNodeUserAccount
		err = tx.Where("node_id = ?", firstNode.ID).Find(&approvers).Error
		if err != nil {
			return fmt.Errorf("获取节点审批人员失败: %w", err)
		}

		if len(approvers) == 0 {
			return fmt.Errorf("节点未配置审批人员")
		}

		// 6. 创建审批人员记录
		for _, approver := range approvers {
			userCase := entity.ApprovalNodeCaseUser{
				ApprovalNodeCaseID: nodeCase.ID,
				UserAccountID:      approver.UserAccountID,
			}
			if err := tx.Create(&userCase).Error; err != nil {
				return fmt.Errorf("创建审批人员记录失败: %w", err)
			}
		}

		// 7. 获取并创建抄送记录
		var copyUsers []entity.ApprovalCopyUserAccount
		err = tx.Where("approval_flow_template_id = ?", templateID).Find(&copyUsers).Error
		if err != nil {
			return fmt.Errorf("获取抄送人员失败: %w", err)
		}

		for _, copyUser := range copyUsers {
			copyCase := entity.ApprovalCopyUserAccountCase{
				ApprovalFlowManagementID: flowID,
				UserAccountID:            copyUser.UserAccountID,
				CopyInfo:                 "审批流已创建",
			}
			if err := tx.Create(&copyCase).Error; err != nil {
				return fmt.Errorf("创建抄送记录失败: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return flowID, nil
}

// UpdateStatus 更新审批流状态
func (r *GormApprovalFlowManagementRepository) UpdateStatus(flowID int, status int8) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// 如果状态是完成状态（已通过、已驳回、已撤销），更新完成时间
	if status == 10 || status == 20 || status == 99 {
		now := time.Now()
		updates["complete_time"] = &now
	}

	return r.db.Model(&entity.ApprovalFlowManagement{}).
		Where("id = ?", flowID).
		Updates(updates).Error
}

// Cancel 撤销审批流
func (r *GormApprovalFlowManagementRepository) Cancel(flowID int, userID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 验证审批流是否存在且是发起人
		var flow entity.ApprovalFlowManagement
		if err := tx.First(&flow, flowID).Error; err != nil {
			return fmt.Errorf("审批流不存在: %w", err)
		}

		if flow.CreateUser != userID {
			return fmt.Errorf("只有发起人才能撤销审批流")
		}

		// 2. 检查状态是否允许撤销
		if flow.Status != 0 {
			return fmt.Errorf("只能撤销待审批状态的审批流")
		}

		// 3. 更新审批流状态为已撤销
		now := time.Now()
		err := tx.Model(&entity.ApprovalFlowManagement{}).
			Where("id = ?", flowID).
			Updates(map[string]interface{}{
				"status":        99,
				"complete_time": &now,
			}).Error
		if err != nil {
			return fmt.Errorf("更新审批流状态失败: %w", err)
		}

		return nil
	})
}

// IncrementStep 增加step步骤
func (r *GormApprovalFlowManagementRepository) IncrementStep(flowID int) error {
	return r.db.Model(&entity.ApprovalFlowManagement{}).
		Where("id = ?", flowID).
		UpdateColumn("step", gorm.Expr("step + ?", 1)).Error
}
