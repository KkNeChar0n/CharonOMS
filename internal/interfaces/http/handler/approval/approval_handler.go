package approval

import (
	approvalService "charonoms/internal/application/service/approval"
	"charonoms/internal/domain/approval/entity"
	approvalDTO "charonoms/internal/interfaces/http/dto/approval"
	"charonoms/pkg/response"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ApprovalHandler 审批流HTTP处理器
type ApprovalHandler struct {
	flowTypeService     *approvalService.ApprovalFlowTypeService
	flowTemplateService *approvalService.ApprovalFlowTemplateService
	flowMgmtService     *approvalService.ApprovalFlowManagementService
}

// NewApprovalHandler 创建审批流处理器实例
func NewApprovalHandler(
	flowTypeService *approvalService.ApprovalFlowTypeService,
	flowTemplateService *approvalService.ApprovalFlowTemplateService,
	flowMgmtService *approvalService.ApprovalFlowManagementService,
) *ApprovalHandler {
	return &ApprovalHandler{
		flowTypeService:     flowTypeService,
		flowTemplateService: flowTemplateService,
		flowMgmtService:     flowMgmtService,
	}
}

// GetApprovalFlowTypes 获取审批流类型列表
// @route GET /api/approval_flow_type
func (h *ApprovalHandler) GetApprovalFlowTypes(c *gin.Context) {
	var req approvalDTO.ApprovalFlowTypeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.ID != "" {
		filters["id"] = req.ID
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}

	// 查询列表
	types, err := h.flowTypeService.GetList(filters)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 转换为响应格式
	var responseData []approvalDTO.ApprovalFlowTypeResponse
	for _, t := range types {
		responseData = append(responseData, approvalDTO.ApprovalFlowTypeResponse{
			ID:         t.ID,
			Name:       t.Name,
			Status:     t.Status,
			CreateTime: t.CreateTime,
			UpdateTime: t.UpdateTime,
		})
	}

	response.Success(c, gin.H{
		"approval_flow_types": responseData,
	})
}

// UpdateApprovalFlowTypeStatus 更新审批流类型状态
// @route PUT /api/approval_flow_type/:id/status
func (h *ApprovalHandler) UpdateApprovalFlowTypeStatus(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "审批流类型ID无效")
		return
	}

	// 绑定请求体
	var req approvalDTO.ApprovalFlowTypeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 更新状态
	if err := h.flowTypeService.UpdateStatus(id, req.Status); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

// CreateApprovalFlowType 创建审批流类型
// @route POST /api/approval-flow-types
func (h *ApprovalHandler) CreateApprovalFlowType(c *gin.Context) {
	// 绑定请求体
	var req approvalDTO.ApprovalFlowTypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证必填字段
	if req.Name == "" {
		response.BadRequest(c, "类型名称不能为空")
		return
	}

	// 创建实体
	flowType := &entity.ApprovalFlowType{
		Name:   req.Name,
		Status: req.Status,
	}

	// 调用服务创建
	if err := h.flowTypeService.Create(flowType); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"id": flowType.ID,
	})
}

// GetApprovalFlowTemplates 获取审批流模板列表
// @route GET /api/approval_flow_template
func (h *ApprovalHandler) GetApprovalFlowTemplates(c *gin.Context) {
	var req approvalDTO.ApprovalFlowTemplateListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.ID != "" {
		filters["id"] = req.ID
	}
	if req.ApprovalFlowTypeID != "" {
		filters["approval_flow_type_id"] = req.ApprovalFlowTypeID
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}

	// 查询列表（已包含关联数据）
	templates, err := h.flowTemplateService.GetList(filters)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"approval_flow_templates": templates,
	})
}

// GetApprovalFlowTemplateDetail 获取审批流模板详情
// @route GET /api/approval_flow_template/:id
func (h *ApprovalHandler) GetApprovalFlowTemplateDetail(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "模板ID无效")
		return
	}

	// 查询详情（包含节点、审批人员、抄送人员）
	detail, err := h.flowTemplateService.GetDetail(id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, detail)
}

// CreateApprovalFlowTemplate 创建审批流模板
// @route POST /api/approval_flow_template
func (h *ApprovalHandler) CreateApprovalFlowTemplate(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			response.InternalServerError(c, "Panic: "+fmt.Sprint(r))
		}
	}()

	// 检查Service是否为nil
	if h.flowTemplateService == nil {
		response.InternalServerError(c, "flowTemplateService is nil")
		return
	}

	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	// 绑定请求体
	var req approvalDTO.CreateApprovalFlowTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 创建模板实体
	template := &entity.ApprovalFlowTemplate{
		Name:               req.Name,
		ApprovalFlowTypeID: req.ApprovalFlowTypeID,
		Creator:            strconv.Itoa(userID),
		Status:             0, // 默认启用
	}

	// 创建节点实体列表
	var nodes []entity.ApprovalFlowTemplateNode
	nodeApprovers := make(map[int][]int)
	for i, nodeReq := range req.Nodes {
		nodes = append(nodes, entity.ApprovalFlowTemplateNode{
			Name: nodeReq.Name,
			Sort: i + 1,
			Type: nodeReq.Type,
		})
		nodeApprovers[i] = nodeReq.Approvers
	}

	// 调用仓储创建（事务）
	if err := h.flowTemplateService.Create(template, nodes, nodeApprovers, req.CopyUsers); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "模板创建成功", gin.H{
		"id": template.ID,
	})
}

// UpdateApprovalFlowTemplateStatus 更新审批流模板状态
// @route PUT /api/approval_flow_template/:id/status
func (h *ApprovalHandler) UpdateApprovalFlowTemplateStatus(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "模板ID无效")
		return
	}

	// 绑定请求体
	var req approvalDTO.ApprovalFlowTemplateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 如果启用模板，需要禁用同类型的其他模板
	if req.Status == 0 {
		template, err := h.flowTemplateService.GetByID(id)
		if err != nil {
			response.HandleError(c, err)
			return
		}

		if err := h.flowTemplateService.DisableSameTypeTemplates(template.ApprovalFlowTypeID, id); err != nil {
			response.HandleError(c, err)
			return
		}
	}

	// 更新状态
	if err := h.flowTemplateService.UpdateStatus(id, req.Status); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "状态更新成功", nil)
}

// GetInitiatedFlows 获取我发起的审批流
// @route GET /api/approval_flow_management/initiated
func (h *ApprovalHandler) GetInitiatedFlows(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	var req approvalDTO.InitiatedFlowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.ID != "" {
		filters["id"] = req.ID
	}
	if req.ApprovalFlowTypeID != "" {
		filters["approval_flow_type_id"] = req.ApprovalFlowTypeID
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}

	// 查询列表
	flows, err := h.flowMgmtService.GetInitiatedFlows(userID, filters)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"initiated_flows": flows,
	})
}

// GetPendingFlows 获取待我审批的任务
// @route GET /api/approval_flow_management/pending
func (h *ApprovalHandler) GetPendingFlows(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	var req approvalDTO.PendingFlowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.ID != "" {
		filters["id"] = req.ID
	}
	if req.ApprovalFlowID != "" {
		filters["approval_flow_id"] = req.ApprovalFlowID
	}
	if req.ApprovalFlowTypeID != "" {
		filters["approval_flow_type_id"] = req.ApprovalFlowTypeID
	}

	// 查询列表
	flows, err := h.flowMgmtService.GetPendingFlows(userID, filters)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"pending_flows": flows,
	})
}

// GetCompletedFlows 获取处理完成的审批
// @route GET /api/approval_flow_management/completed
func (h *ApprovalHandler) GetCompletedFlows(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	var req approvalDTO.CompletedFlowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.ID != "" {
		filters["id"] = req.ID
	}
	if req.ApprovalFlowID != "" {
		filters["approval_flow_id"] = req.ApprovalFlowID
	}
	if req.ApprovalFlowTypeID != "" {
		filters["approval_flow_type_id"] = req.ApprovalFlowTypeID
	}

	// 查询列表
	flows, err := h.flowMgmtService.GetCompletedFlows(userID, filters)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"completed_flows": flows,
	})
}

// GetCopiedFlows 获取抄送我的通知
// @route GET /api/approval_flow_management/copied
func (h *ApprovalHandler) GetCopiedFlows(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	var req approvalDTO.CopiedFlowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if req.ID != "" {
		filters["id"] = req.ID
	}
	if req.ApprovalFlowID != "" {
		filters["approval_flow_id"] = req.ApprovalFlowID
	}
	if req.ApprovalFlowTypeID != "" {
		filters["approval_flow_type_id"] = req.ApprovalFlowTypeID
	}

	// 查询列表
	flows, err := h.flowMgmtService.GetCopiedFlows(userID, filters)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"copied_flows": flows,
	})
}

// GetApprovalFlowDetail 获取审批流详情
// @route GET /api/approval_flow_management/:id
func (h *ApprovalHandler) GetApprovalFlowDetail(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	// 解析路径参数
	idStr := c.Param("id")
	flowID, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "审批流ID无效")
		return
	}

	// 查询详情
	detail, err := h.flowMgmtService.GetDetail(flowID, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, detail)
}

// CreateFromTemplate 从模板创建审批流
// @route POST /api/approval_flow_management/create_from_template
func (h *ApprovalHandler) CreateFromTemplate(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	// 绑定请求体
	var req approvalDTO.CreateFromTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 创建审批流实例
	flowID, err := h.flowMgmtService.CreateFromTemplate(req.TemplateID, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "审批流创建成功", gin.H{
		"id": flowID,
	})
}

// CancelApprovalFlow 撤销审批流
// @route POST /api/approval_flow_management/:id/cancel
func (h *ApprovalHandler) CancelApprovalFlow(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	// 解析路径参数
	idStr := c.Param("id")
	flowID, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "审批流ID无效")
		return
	}

	// 撤销审批流
	if err := h.flowMgmtService.Cancel(flowID, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "审批流已撤销", nil)
}

// ApproveFlow 审批通过
// @route POST /api/approval_flow_management/approve
func (h *ApprovalHandler) ApproveFlow(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	// 绑定请求体
	var req approvalDTO.ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 调用领域服务处理审批通过逻辑
	if err := h.flowMgmtService.Approve(req.NodeCaseUserID, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "审批通过", nil)
}

// RejectFlow 审批驳回
// @route POST /api/approval_flow_management/reject
func (h *ApprovalHandler) RejectFlow(c *gin.Context) {
	// 获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}
	userID := int(userIDVal.(uint))

	// 绑定请求体
	var req approvalDTO.RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 调用领域服务处理审批驳回逻辑
	if err := h.flowMgmtService.Reject(req.NodeCaseUserID, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMessage(c, "审批已驳回", nil)
}
