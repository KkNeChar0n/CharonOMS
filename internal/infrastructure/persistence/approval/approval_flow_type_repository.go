package approval

import (
	"charonoms/internal/domain/approval/entity"
	"charonoms/internal/domain/approval/repository"
	"gorm.io/gorm"
)

// GormApprovalFlowTypeRepository GORM实现的审批流类型仓储
type GormApprovalFlowTypeRepository struct {
	db *gorm.DB
}

// NewApprovalFlowTypeRepository 创建审批流类型仓储实例
func NewApprovalFlowTypeRepository(db *gorm.DB) repository.ApprovalFlowTypeRepository {
	return &GormApprovalFlowTypeRepository{db: db}
}

// GetList 获取审批流类型列表
func (r *GormApprovalFlowTypeRepository) GetList(filters map[string]interface{}) ([]entity.ApprovalFlowType, error) {
	var types []entity.ApprovalFlowType
	query := r.db

	// 应用筛选条件
	if id, ok := filters["id"]; ok && id != "" {
		query = query.Where("id = ?", id)
	}
	if name, ok := filters["name"]; ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&types).Error
	return types, err
}

// GetByID 根据ID查询审批流类型
func (r *GormApprovalFlowTypeRepository) GetByID(id int) (*entity.ApprovalFlowType, error) {
	var flowType entity.ApprovalFlowType
	err := r.db.First(&flowType, id).Error
	if err != nil {
		return nil, err
	}
	return &flowType, nil
}

// Create 创建审批流类型
func (r *GormApprovalFlowTypeRepository) Create(flowType *entity.ApprovalFlowType) error {
	return r.db.Create(flowType).Error
}

// UpdateStatus 更新审批流类型状态
func (r *GormApprovalFlowTypeRepository) UpdateStatus(id int, status int8) error {
	return r.db.Model(&entity.ApprovalFlowType{}).Where("id = ?", id).Update("status", status).Error
}
