package entity

import "errors"

var (
	// ErrTemplateNameRequired 模板名称不能为空
	ErrTemplateNameRequired = errors.New("模板名称不能为空")

	// ErrInvalidTemplateType 无效的活动类型
	ErrInvalidTemplateType = errors.New("活动类型必须为 1/2/3")

	// ErrInvalidSelectType 无效的选择方式
	ErrInvalidSelectType = errors.New("选择方式必须为 1/2")

	// ErrTemplateNotFound 活动模板不存在
	ErrTemplateNotFound = errors.New("活动模板不存在")

	// ErrTemplateHasActivities 模板有关联活动，无法删除
	ErrTemplateHasActivities = errors.New("该模板有关联活动，无法删除")

	// ErrTemplateIsEnabled 活动模板启用中，无法编辑
	ErrTemplateIsEnabled = errors.New("活动模板启用中，无法编辑")

	// ErrTemplateIsDisabled 活动模板已禁用，无法使用
	ErrTemplateIsDisabled = errors.New("活动模板已禁用，无法使用")

	// ErrMissingRelationConfig 缺少关联配置
	ErrMissingRelationConfig = errors.New("模板名称、活动类型和选择方式不能为空")
)
