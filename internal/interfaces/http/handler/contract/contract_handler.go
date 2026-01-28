package contract

import (
	"net/http"
	"strconv"
	"strings"

	contractService "charonoms/internal/application/service/contract"

	"github.com/gin-gonic/gin"
)

// ContractHandler 合同HTTP处理器
type ContractHandler struct {
	contractService *contractService.ContractService
}

// NewContractHandler 创建合同HTTP处理器实例
func NewContractHandler(contractService *contractService.ContractService) *ContractHandler {
	return &ContractHandler{
		contractService: contractService,
	}
}

// GetContracts 获取合同列表
// GET /api/contracts
func (h *ContractHandler) GetContracts(c *gin.Context) {
	contracts, err := h.contractService.GetContractList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contracts": contracts})
}

// GetContractByID 获取合同详情
// GET /api/contracts/:id
func (h *ContractHandler) GetContractByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract id"})
		return
	}

	contract, err := h.contractService.GetContractByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contract": contract})
}

// CreateContract 创建合同
// POST /api/contracts
func (h *ContractHandler) CreateContract(c *gin.Context) {
	var req contractService.CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token获取当前用户名作为发起人
	initiator, _ := c.Get("username")
	initiatorStr, ok := initiator.(string)
	if !ok {
		initiatorStr = ""
	}

	contractID, err := h.contractService.CreateContract(&req, initiatorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      contractID,
		"message": "合同新增成功",
	})
}

// RevokeContract 撤销合同
// PUT /api/contracts/:id/revoke
func (h *ContractHandler) RevokeContract(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract id"})
		return
	}

	if err := h.contractService.RevokeContract(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "操作成功"})
}

// TerminateContract 中止合作
// PUT /api/contracts/:id/terminate
func (h *ContractHandler) TerminateContract(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract id"})
		return
	}

	var req contractService.TerminateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.contractService.TerminateContract(id, &req); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "操作成功"})
}
