package dto

import (
	"fmt"
	"strings"
	"time"
)

// FlexibleTime 支持多种格式的时间类型
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON 自定义JSON反序列化，支持多种时间格式
func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		ft.Time = time.Time{}
		return nil
	}

	// 支持的时间格式列表
	formats := []string{
		time.RFC3339,                // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05",       // "2026-01-01T15:45:00"
		"2006-01-02T15:04",          // "2026-01-01T15:45"
		"2006-01-02 15:04:05",       // "2026-01-01 15:45:00"
		"2006-01-02 15:04",          // "2026-01-01 15:45"
		"2006-01-02",                // "2026-01-01"
	}

	var err error
	for _, format := range formats {
		ft.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("无法解析时间格式: %s", s)
}

// CreateActivityDTO 创建活动请求
type CreateActivityDTO struct {
	TemplateID int                 `json:"template_id" binding:"required"`
	Name       string              `json:"name" binding:"required"`
	StartTime  FlexibleTime        `json:"start_time" binding:"required"`
	EndTime    FlexibleTime        `json:"end_time" binding:"required"`
	Details    []ActivityDetailDTO `json:"details"`
	Status     int                 `json:"status"`
}

// UpdateActivityDTO 更新活动请求
type UpdateActivityDTO struct {
	TemplateID int                 `json:"template_id" binding:"required"`
	Name       string              `json:"name" binding:"required"`
	StartTime  FlexibleTime        `json:"start_time" binding:"required"`
	EndTime    FlexibleTime        `json:"end_time" binding:"required"`
	Details    []ActivityDetailDTO `json:"details"`
	Status     int                 `json:"status"`
}

// UpdateActivityStatusDTO 更新活动状态请求
type UpdateActivityStatusDTO struct {
	Status *int `json:"status" binding:"required,oneof=0 1"`
}

// ActivityDetailDTO 活动详情
type ActivityDetailDTO struct {
	ID              int     `json:"id,omitempty"`
	ActivityID      int     `json:"activity_id,omitempty"`
	ThresholdAmount float64 `json:"threshold_amount"`
	DiscountValue   float64 `json:"discount_value"`
}

// ActivityDTO 活动响应
type ActivityDTO struct {
	ID                   int                 `json:"id"`
	Name                 string              `json:"name"`
	TemplateID           int                 `json:"template_id"`
	TemplateName         string              `json:"template_name"`
	TemplateType         int                 `json:"template_type"`
	TemplateSelectType   int                 `json:"template_select_type,omitempty"`
	StartTime            time.Time           `json:"start_time"`
	EndTime              time.Time           `json:"end_time"`
	Status               int                 `json:"status"`
	CreateTime           time.Time           `json:"create_time"`
	Details              []ActivityDetailDTO `json:"details,omitempty"`
}

// ActivityDetailResponseDTO 活动详情响应
type ActivityDetailResponseDTO struct {
	ID             int                 `json:"id"`
	Name           string              `json:"name"`
	TemplateID     int                 `json:"template_id"`
	TemplateName   string              `json:"template_name"`
	TemplateType   int                 `json:"template_type"`
	SelectType     int                 `json:"select_type"`
	StartTime      time.Time           `json:"start_time"`
	EndTime        time.Time           `json:"end_time"`
	Status         int                 `json:"status"`
	Details        []ActivityDetailDTO `json:"details"`
}

// ActivitiesByDateRangeDTO 按日期范围查询活动响应
type ActivitiesByDateRangeDTO struct {
	HasDuplicate  bool          `json:"has_duplicate"`
	DuplicateType *int          `json:"duplicate_type"`
	TypeName      string        `json:"type_name,omitempty"`
	Activities    []ActivityDTO `json:"activities"`
}
