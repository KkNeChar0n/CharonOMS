package classify

import (
	"charonoms/internal/domain/classify/entity"
	"charonoms/internal/domain/classify/repository"
	"errors"
	"fmt"
	"strconv"
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

// toIntPtr 将interface{}转换为*int
func toIntPtr(v interface{}) (*int, error) {
	if v == nil {
		return nil, nil
	}

	switch val := v.(type) {
	case int:
		return &val, nil
	case float64:
		intVal := int(val)
		return &intVal, nil
	case string:
		if val == "" {
			return nil, nil
		}
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return &intVal, nil
	default:
		return nil, errors.New("无效的数据类型")
	}
}

// ClassifyService 分类业务服务
type ClassifyService struct {
	repo repository.ClassifyRepository
}

// NewClassifyService 创建分类服务实例
func NewClassifyService(repo repository.ClassifyRepository) *ClassifyService {
	return &ClassifyService{repo: repo}
}

// GetClassifyList 获取分类列表
func (s *ClassifyService) GetClassifyList() (*ClassifyListResponse, error) {
	classifies, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %w", err)
	}

	return &ClassifyListResponse{Classifies: classifies}, nil
}

// GetParents 获取一级分类列表
func (s *ClassifyService) GetParents() (*ParentsListResponse, error) {
	parents, err := s.repo.GetParents()
	if err != nil {
		return nil, fmt.Errorf("获取一级分类列表失败: %w", err)
	}

	// 转换为 map 格式
	result := make([]map[string]interface{}, len(parents))
	for i, parent := range parents {
		result[i] = map[string]interface{}{
			"id":   parent.ID,
			"name": parent.Name,
		}
	}

	return &ParentsListResponse{Parents: result}, nil
}

// GetActiveClassifies 获取启用分类列表
func (s *ClassifyService) GetActiveClassifies() (*ActiveClassifyResponse, error) {
	classifies, err := s.repo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("获取启用分类列表失败: %w", err)
	}

	return &ActiveClassifyResponse{Classifies: classifies}, nil
}

// GetClassifyByID 获取分类详情
func (s *ClassifyService) GetClassifyByID(id int) (*entity.Classify, error) {
	classify, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("类型不存在")
	}
	return classify, nil
}

// CreateClassify 创建分类
func (s *ClassifyService) CreateClassify(req *CreateClassifyRequest) (int, error) {
	// 验证必填字段
	if req.Name == "" {
		return 0, errors.New("名称和级别不能为空")
	}

	// 转换level为int
	level, err := toInt(req.Level)
	if err != nil {
		return 0, errors.New("级别值无效")
	}

	// 验证级别值
	if level != 0 && level != 1 {
		return 0, errors.New("级别值必须为0或1")
	}

	// 转换parentID为*int
	parentID, err := toIntPtr(req.ParentID)
	if err != nil {
		return 0, errors.New("父级ID无效")
	}

	// 验证二级分类必须有父级
	if level == 1 {
		if parentID == nil || *parentID <= 0 {
			return 0, errors.New("二级类型必须选择父级类型")
		}
	}

	// 验证名称唯一性
	if level == 0 {
		// 一级分类：检查名称在所有一级分类中是否唯一
		isUnique, err := s.repo.CheckNameUnique(req.Name, 0, nil, 0)
		if err != nil {
			return 0, fmt.Errorf("检查名称唯一性失败: %w", err)
		}
		if !isUnique {
			return 0, errors.New("该一级类型名称已存在")
		}
	} else {
		// 二级分类：检查名称在同一父级下是否唯一
		isUnique, err := s.repo.CheckNameUnique(req.Name, 1, parentID, 0)
		if err != nil {
			return 0, fmt.Errorf("检查名称唯一性失败: %w", err)
		}
		if !isUnique {
			return 0, errors.New("该父级类型下已存在同名的二级类型")
		}
	}

	// 创建分类实体
	classify := &entity.Classify{
		Name:     req.Name,
		Level:    level,
		ParentID: parentID,
		Status:   0, // 默认启用
	}

	// 保存分类
	if err := s.repo.Create(classify); err != nil {
		return 0, fmt.Errorf("创建分类失败: %w", err)
	}

	return classify.ID, nil
}

// UpdateClassify 更新分类信息
func (s *ClassifyService) UpdateClassify(id int, req *UpdateClassifyRequest) error {
	// 验证必填字段
	if req.Name == "" {
		return errors.New("名称和级别不能为空")
	}

	// 转换level为int
	level, err := toInt(req.Level)
	if err != nil {
		return errors.New("级别值无效")
	}

	// 验证级别值
	if level != 0 && level != 1 {
		return errors.New("级别值必须为0或1")
	}

	// 转换parentID为*int
	parentID, err := toIntPtr(req.ParentID)
	if err != nil {
		return errors.New("父级ID无效")
	}

	// 验证二级分类必须有父级
	if level == 1 {
		if parentID == nil || *parentID <= 0 {
			return errors.New("二级类型必须选择父级类型")
		}
	}

	// 检查分类是否存在
	_, err = s.repo.GetByID(id)
	if err != nil {
		return errors.New("类型不存在")
	}

	// 验证名称唯一性
	if level == 0 {
		// 一级分类：检查名称在所有一级分类中是否唯一
		isUnique, err := s.repo.CheckNameUnique(req.Name, 0, nil, id)
		if err != nil {
			return fmt.Errorf("检查名称唯一性失败: %w", err)
		}
		if !isUnique {
			return errors.New("该一级类型名称已存在")
		}
	} else {
		// 二级分类：检查名称在同一父级下是否唯一
		isUnique, err := s.repo.CheckNameUnique(req.Name, 1, parentID, id)
		if err != nil {
			return fmt.Errorf("检查名称唯一性失败: %w", err)
		}
		if !isUnique {
			return errors.New("该父级类型下已存在同名的二级类型")
		}
	}

	// 更新分类信息
	classify := &entity.Classify{
		ID:       id,
		Name:     req.Name,
		Level:    level,
		ParentID: parentID,
	}

	if err := s.repo.Update(classify); err != nil {
		return fmt.Errorf("更新分类失败: %w", err)
	}

	return nil
}

// UpdateClassifyStatus 更新分类状态
func (s *ClassifyService) UpdateClassifyStatus(id int, req *UpdateClassifyStatusRequest) error {
	// Status是指针类型，binding已经确保了它不为nil且值在0或1之间
	status := *req.Status

	// 验证状态值
	if status != 0 && status != 1 {
		return errors.New("状态值必须为0或1")
	}

	// 检查分类是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("类型不存在")
	}

	// 更新状态
	if err := s.repo.UpdateStatus(id, status); err != nil {
		return fmt.Errorf("更新分类状态失败: %w", err)
	}

	return nil
}
