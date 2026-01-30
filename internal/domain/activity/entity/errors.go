package entity

import "errors"

var (
	// ErrActivityNameRequired 活动名称不能为空
	ErrActivityNameRequired = errors.New("活动名称、关联模板、开始时间和结束时间不能为空")

	// ErrTemplateIDRequired 模板ID不能为空
	ErrTemplateIDRequired = errors.New("活动名称、关联模板、开始时间和结束时间不能为空")

	// ErrTimeRequired 时间不能为空
	ErrTimeRequired = errors.New("活动名称、关联模板、开始时间和结束时间不能为空")

	// ErrInvalidTimeRange 无效的时间范围
	ErrInvalidTimeRange = errors.New("活动名称、关联模板、开始时间和结束时间不能为空")

	// ErrActivityNotFound 活动不存在
	ErrActivityNotFound = errors.New("活动不存在")

	// ErrActivityIsEnabled 活动启用中，无法编辑
	ErrActivityIsEnabled = errors.New("活动启用中，无法编辑")

	// ErrTemplateNotEnabled 活动模板未启用
	ErrTemplateNotEnabled = errors.New("活动模板未启用")

	// ErrPaymentTimeRequired 预计付款时间不能为空
	ErrPaymentTimeRequired = errors.New("预计付款时间不能为空")

	// ErrInvalidThresholdAmount 无效的门槛金额
	ErrInvalidThresholdAmount = errors.New("门槛金额必须大于0")

	// ErrInvalidDiscountValue 无效的折扣值
	ErrInvalidDiscountValue = errors.New("折扣值无效")
)
