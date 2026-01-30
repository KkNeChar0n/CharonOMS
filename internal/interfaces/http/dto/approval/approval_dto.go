package approval

import "time"

// ApprovalFlowTypeListRequest 审批流类型列表请求
type ApprovalFlowTypeListRequest struct {
	ID     string `form:"id"`
	Name   string `form:"name"`
	Status string `form:"status"`
}

// ApprovalFlowTypeCreateRequest 审批流类型创建请求
type ApprovalFlowTypeCreateRequest struct {
	Name   string `json:"name" binding:"required"`
	Status int8   `json:"status"`
}

// ApprovalFlowTypeStatusRequest 审批流类型状态更新请求
type ApprovalFlowTypeStatusRequest struct {
	Status int8 `json:"status"`
}

// ApprovalFlowTemplateListRequest 审批流模板列表请求
type ApprovalFlowTemplateListRequest struct {
	ID                  string `form:"id"`
	ApprovalFlowTypeID  string `form:"approval_flow_type_id"`
	Name                string `form:"name"`
	Status              string `form:"status"`
}

// ApprovalFlowTemplateNode 审批流模板节点
type ApprovalFlowTemplateNode struct {
	Name      string `json:"name"`
	Type      int8   `json:"type"` // 0=会签，1=或签
	Approvers []int  `json:"approvers"`
}

// CreateApprovalFlowTemplateRequest 创建审批流模板请求
type CreateApprovalFlowTemplateRequest struct {
	Name               string                     `json:"name"`
	ApprovalFlowTypeID int                        `json:"approval_flow_type_id"`
	Nodes              []ApprovalFlowTemplateNode `json:"nodes"`
	CopyUsers          []int                      `json:"copy_users"`
}

// ApprovalFlowTemplateStatusRequest 审批流模板状态更新请求
type ApprovalFlowTemplateStatusRequest struct {
	Status int8 `json:"status"`
}

// InitiatedFlowsRequest 我发起的审批流请求
type InitiatedFlowsRequest struct {
	ID                 string `form:"id"`
	ApprovalFlowTypeID string `form:"approval_flow_type_id"`
	Status             string `form:"status"`
}

// PendingFlowsRequest 待我审批的任务请求
type PendingFlowsRequest struct {
	ID                 string `form:"id"`
	ApprovalFlowID     string `form:"approval_flow_id"`
	ApprovalFlowTypeID string `form:"approval_flow_type_id"`
}

// CompletedFlowsRequest 处理完成的审批请求
type CompletedFlowsRequest struct {
	ID                 string `form:"id"`
	ApprovalFlowID     string `form:"approval_flow_id"`
	ApprovalFlowTypeID string `form:"approval_flow_type_id"`
}

// CopiedFlowsRequest 抄送我的通知请求
type CopiedFlowsRequest struct {
	ID                 string `form:"id"`
	ApprovalFlowID     string `form:"approval_flow_id"`
	ApprovalFlowTypeID string `form:"approval_flow_type_id"`
}

// CreateFromTemplateRequest 从模板创建审批流请求
type CreateFromTemplateRequest struct {
	TemplateID int `json:"approval_flow_template_id"`
	Title      string `json:"title"`
	Info       string `json:"info"`
}

// ApproveRequest 审批通过请求
type ApproveRequest struct {
	NodeCaseUserID int `json:"node_case_user_id" binding:"required"`
}

// RejectRequest 审批驳回请求
type RejectRequest struct {
	NodeCaseUserID int `json:"node_case_user_id" binding:"required"`
}

// ApprovalFlowTypeResponse 审批流类型响应
type ApprovalFlowTypeResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Status     int8      `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// ApprovalFlowTemplateResponse 审批流模板响应
type ApprovalFlowTemplateResponse struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	ApprovalFlowTypeID int       `json:"approval_flow_type_id"`
	FlowTypeName       string    `json:"flow_type_name"`
	Creator            string    `json:"creator"`
	Status             int8      `json:"status"`
	CreateTime         time.Time `json:"create_time"`
	UpdateTime         time.Time `json:"update_time"`
}

// NodeResponse 节点响应
type NodeResponse struct {
	ID        int                  `json:"id"`
	Name      string               `json:"name"`
	Sort      int                  `json:"sort"`
	Type      int8                 `json:"type"`
	Approvers []map[string]interface{} `json:"approvers"`
}

// ApprovalFlowTemplateDetailResponse 审批流模板详情响应
type ApprovalFlowTemplateDetailResponse struct {
	ID                 int                      `json:"id"`
	Name               string                   `json:"name"`
	ApprovalFlowTypeID int                      `json:"approval_flow_type_id"`
	FlowTypeName       string                   `json:"flow_type_name"`
	Creator            string                   `json:"creator"`
	Status             int8                     `json:"status"`
	CreateTime         time.Time                `json:"create_time"`
	UpdateTime         time.Time                `json:"update_time"`
	Nodes              []map[string]interface{} `json:"nodes"`
	CopyUsers          []map[string]interface{} `json:"copy_users"`
}

// InitiatedFlowResponse 我发起的审批流响应
type InitiatedFlowResponse struct {
	ID                     int        `json:"id"`
	ApprovalFlowTemplateID int        `json:"approval_flow_template_id"`
	ApprovalFlowTypeID     int        `json:"approval_flow_type_id"`
	FlowTypeName           string     `json:"flow_type_name"`
	Step                   int        `json:"step"`
	CreateUser             int        `json:"create_user"`
	CreateTime             time.Time  `json:"create_time"`
	Status                 int8       `json:"status"`
	CompleteTime           *time.Time `json:"complete_time"`
}

// PendingFlowResponse 待我审批的任务响应
type PendingFlowResponse struct {
	ID                 int       `json:"id"`
	ApprovalFlowID     int       `json:"approval_flow_id"`
	ApprovalFlowTypeID int       `json:"approval_flow_type_id"`
	FlowTypeName       string    `json:"flow_type_name"`
	CreateUserName     string    `json:"create_user_name"`
	CreateTime         time.Time `json:"create_time"`
	NodeType           int8      `json:"node_type"`
	NodeSort           int       `json:"node_sort"`
	NodeCaseUserID     int       `json:"node_case_user_id"`
}

// CompletedFlowResponse 处理完成的审批响应
type CompletedFlowResponse struct {
	ID                 int        `json:"id"`
	ApprovalFlowID     int        `json:"approval_flow_id"`
	ApprovalFlowTypeID int        `json:"approval_flow_type_id"`
	FlowTypeName       string     `json:"flow_type_name"`
	CreateUserName     string     `json:"create_user_name"`
	CreateTime         time.Time  `json:"create_time"`
	NodeType           int8       `json:"node_type"`
	NodeSort           int        `json:"node_sort"`
	Result             int8       `json:"result"`
	HandleTime         *time.Time `json:"handle_time"`
}

// CopiedFlowResponse 抄送我的通知响应
type CopiedFlowResponse struct {
	ID                 int        `json:"id"`
	ApprovalFlowID     int        `json:"approval_flow_id"`
	ApprovalFlowTypeID int        `json:"approval_flow_type_id"`
	FlowTypeName       string     `json:"flow_type_name"`
	CreateUserName     string     `json:"create_user_name"`
	CreateTime         time.Time  `json:"create_time"`
	CompleteTime       *time.Time `json:"complete_time"`
	CopyInfo           string     `json:"copy_info"`
}
