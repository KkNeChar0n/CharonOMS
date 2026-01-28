package attribute

import (
	"charonoms/internal/domain/attribute/entity"
	"charonoms/internal/domain/attribute/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// toInt 将interface{}转换为int
func toInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		if val == "" {
			return 0, errors.New("值不能为空")
		}
		return strconv.Atoi(val)
	default:
		return 0, errors.New("无效的数据类型")
	}
}

// AttributeService 属性业务服务
type AttributeService struct {
	repo repository.AttributeRepository
}

// NewAttributeService 创建属性服务实例
func NewAttributeService(repo repository.AttributeRepository) *AttributeService {
	return &AttributeService{repo: repo}
}

// GetAttributeList 获取属性列表
func (s *AttributeService) GetAttributeList() (*AttributeListResponse, error) {
	attributes, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取属性列表失败: %w", err)
	}

	return &AttributeListResponse{Attributes: attributes}, nil
}

// GetActiveAttributes 获取启用属性列表
func (s *AttributeService) GetActiveAttributes() (*ActiveAttributeResponse, error) {
	attributes, err := s.repo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("获取启用属性列表失败: %w", err)
	}

	return &ActiveAttributeResponse{Attributes: attributes}, nil
}

// GetAttributeByID 获取属性详情
func (s *AttributeService) GetAttributeByID(id int) (*entity.Attribute, error) {
	attribute, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("属性不存在")
	}
	return attribute, nil
}

// CreateAttribute 创建属性
func (s *AttributeService) CreateAttribute(req *CreateAttributeRequest) (int, error) {
	// 验证名称和分类不能为空
	if strings.TrimSpace(req.Name) == "" {
		return 0, errors.New("名称和分类不能为空")
	}

	// 转换classify为int
	classify, err := toInt(req.Classify)
	if err != nil {
		return 0, errors.New("分类值无效")
	}

	// 验证分类值必须为0或1
	if classify != 0 && classify != 1 {
		return 0, errors.New("分类值必须为0或1")
	}

	// 创建属性实体
	attribute := &entity.Attribute{
		Name:     req.Name,
		Classify: classify,
		Status:   0, // 默认启用
	}

	// 保存属性
	if err := s.repo.Create(attribute); err != nil {
		return 0, fmt.Errorf("创建属性失败: %w", err)
	}

	return attribute.ID, nil
}

// UpdateAttribute 更新属性信息
func (s *AttributeService) UpdateAttribute(id int, req *UpdateAttributeRequest) error {
	// 验证名称和分类不能为空
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("名称和分类不能为空")
	}

	// 转换classify为int
	classify, err := toInt(req.Classify)
	if err != nil {
		return errors.New("分类值无效")
	}

	// 验证分类值必须为0或1
	if classify != 0 && classify != 1 {
		return errors.New("分类值必须为0或1")
	}

	// 检查属性是否存在
	_, err = s.repo.GetByID(id)
	if err != nil {
		return errors.New("属性不存在")
	}

	// 更新属性信息
	attribute := &entity.Attribute{
		ID:       id,
		Name:     req.Name,
		Classify: classify,
	}

	if err := s.repo.Update(attribute); err != nil {
		return fmt.Errorf("更新属性失败: %w", err)
	}

	return nil
}

// UpdateAttributeStatus 更新属性状态
func (s *AttributeService) UpdateAttributeStatus(id int, req *UpdateAttributeStatusRequest) error {
	// Status是指针类型，binding已经确保了它不为nil且值在0或1之间
	status := *req.Status

	// 检查属性是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("属性不存在")
	}

	// 更新状态
	if err := s.repo.UpdateStatus(id, status); err != nil {
		return fmt.Errorf("更新属性状态失败: %w", err)
	}

	return nil
}

// GetAttributeValues 获取属性值列表
func (s *AttributeService) GetAttributeValues(attributeID int) (*AttributeValuesResponse, error) {
	// 检查属性是否存在
	_, err := s.repo.GetByID(attributeID)
	if err != nil {
		return nil, errors.New("属性不存在")
	}

	// 获取属性值
	values, err := s.repo.GetValues(attributeID)
	if err != nil {
		return nil, fmt.Errorf("获取属性值失败: %w", err)
	}

	// 转换为响应格式
	result := make([]map[string]interface{}, 0, len(values))
	for _, v := range values {
		result = append(result, map[string]interface{}{
			"id":   v.ID,
			"name": v.Name,
		})
	}

	return &AttributeValuesResponse{Values: result}, nil
}

// SaveAttributeValues 保存属性值
func (s *AttributeService) SaveAttributeValues(attributeID int, req *SaveValuesRequest) error {
	// 检查属性是否存在
	_, err := s.repo.GetByID(attributeID)
	if err != nil {
		return errors.New("属性不存在")
	}

	// 验证至少需要填入一条属性值
	if len(req.Values) == 0 {
		return errors.New("至少需要填入一条属性值")
	}

	// 验证属性值不能为空
	var newValues []entity.AttributeValue
	for _, valueName := range req.Values {
		if strings.TrimSpace(valueName) == "" {
			return errors.New("属性值不能为空")
		}
		newValues = append(newValues, entity.AttributeValue{
			Name:        valueName,
			AttributeID: attributeID,
		})
	}

	// 保存属性值（全量替换）
	if err := s.repo.SaveValues(attributeID, newValues); err != nil {
		return fmt.Errorf("保存属性值失败: %w", err)
	}

	return nil
}
