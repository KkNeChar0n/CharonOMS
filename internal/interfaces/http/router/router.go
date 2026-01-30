package router

import (
	accountService "charonoms/internal/application/service/account"
	attributeService "charonoms/internal/application/service/attribute"
	authService "charonoms/internal/application/service/auth"
	basicService "charonoms/internal/application/service/basic"
	brandService "charonoms/internal/application/service/brand"
	classifyService "charonoms/internal/application/service/classify"
	coachService "charonoms/internal/application/service/coach"
	contractService "charonoms/internal/application/service/contract"
	goodsService "charonoms/internal/application/service/goods"
	rbacService "charonoms/internal/application/service/rbac"
	studentService "charonoms/internal/application/service/student"
	activityService "charonoms/internal/application/activity"
	activityTemplateService "charonoms/internal/application/activity_template"
	approvalService "charonoms/internal/application/service/approval"
	"charonoms/internal/infrastructure/config"
	"charonoms/internal/infrastructure/persistence"
	approvalImpl "charonoms/internal/infrastructure/persistence/approval"
	accountImpl "charonoms/internal/infrastructure/persistence/mysql/account"
	attributeImpl "charonoms/internal/infrastructure/persistence/mysql/attribute"
	authImpl "charonoms/internal/infrastructure/persistence/mysql/auth"
	brandImpl "charonoms/internal/infrastructure/persistence/mysql/brand"
	classifyImpl "charonoms/internal/infrastructure/persistence/mysql/classify"
	coachImpl "charonoms/internal/infrastructure/persistence/mysql/coach"
	contractImpl "charonoms/internal/infrastructure/persistence/mysql/contract"
	goodsImpl "charonoms/internal/infrastructure/persistence/mysql/goods"
	rbacImpl "charonoms/internal/infrastructure/persistence/mysql/rbac"
	studentImpl "charonoms/internal/infrastructure/persistence/mysql/student"
	"charonoms/internal/infrastructure/persistence/mysql"
	"charonoms/internal/interfaces/http/handler/account"
	"charonoms/internal/interfaces/http/handler/approval"
	"charonoms/internal/interfaces/http/handler/attribute"
	"charonoms/internal/interfaces/http/handler/auth"
	"charonoms/internal/interfaces/http/handler/basic"
	"charonoms/internal/interfaces/http/handler/brand"
	"charonoms/internal/interfaces/http/handler/classify"
	"charonoms/internal/interfaces/http/handler/coach"
	"charonoms/internal/interfaces/http/handler/contract"
	"charonoms/internal/interfaces/http/handler/goods"
	"charonoms/internal/interfaces/http/handler/placeholder"
	"charonoms/internal/interfaces/http/handler/rbac"
	"charonoms/internal/interfaces/http/handler/student"
	"charonoms/internal/interfaces/http/handler"
	"charonoms/internal/interfaces/http/middleware"
	approvalDomainService "charonoms/internal/domain/approval/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter configure router
func SetupRouter(cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg.CORS))

	// Initialize dependencies
	setupDependencies(r, cfg)

	// Static files (frontend)
	r.Static("/frontend", "./frontend")
	r.StaticFile("/", "./frontend/index.html")

	return r
}

// setupDependencies setup dependency injection
func setupDependencies(r *gin.Engine, cfg *config.Config) {
	// RBAC module (初始化在前，因为 AuthService 需要依赖 roleRepo)
	roleRepo := rbacImpl.NewRoleRepository(mysql.DB)
	permissionRepo := rbacImpl.NewPermissionRepository(mysql.DB)
	menuRepo := rbacImpl.NewMenuRepository(mysql.DB)
	rbacSvc := rbacService.NewRBACService(roleRepo, permissionRepo, menuRepo)
	rbacHdl := rbac.NewRBACHandler(rbacSvc)

	// Auth module
	authRepo := authImpl.NewAuthRepository(mysql.DB)
	authSvc := authService.NewAuthService(authRepo, roleRepo, cfg.JWT)
	authHdl := auth.NewAuthHandler(authSvc)

	// Basic module (sex, grade, subject)
	basicRepo := persistence.NewBasicRepository(mysql.DB)
	basicSvc := basicService.NewBasicService(basicRepo)
	basicHdl := basic.NewBasicHandler(basicSvc)

	// Account module
	accountRepo := accountImpl.NewAccountRepository(mysql.DB)
	accountSvc := accountService.NewAccountService(accountRepo)
	accountHdl := account.NewAccountHandler(accountSvc)

	// Student module
	studentRepo := studentImpl.NewStudentRepository(mysql.DB)
	studentSvc := studentService.NewStudentService(studentRepo)
	studentHdl := student.NewStudentHandler(studentSvc)

	// Coach module
	coachRepo := coachImpl.NewCoachRepository(mysql.DB)
	coachSvc := coachService.NewCoachService(coachRepo)
	coachHdl := coach.NewCoachHandler(coachSvc)

	// Contract module
	contractRepo := contractImpl.NewContractRepository(mysql.DB)
	contractSvc := contractService.NewContractService(contractRepo)
	contractHdl := contract.NewContractHandler(contractSvc)

	// Brand module
	brandRepo := brandImpl.NewBrandRepository(mysql.DB)
	brandSvc := brandService.NewBrandService(brandRepo)
	brandHdl := brand.NewBrandHandler(brandSvc)

	// Classify module
	classifyRepo := classifyImpl.NewClassifyRepository(mysql.DB)
	classifySvc := classifyService.NewClassifyService(classifyRepo)
	classifyHdl := classify.NewClassifyHandler(classifySvc)

	// Attribute module
	attributeRepo := attributeImpl.NewAttributeRepository(mysql.DB)
	attributeSvc := attributeService.NewAttributeService(attributeRepo)
	attributeHdl := attribute.NewAttributeHandler(attributeSvc)

	// Goods module
	goodsRepo := goodsImpl.NewGoodsRepository(mysql.DB)
	goodsSvc := goodsService.NewGoodsService(goodsRepo)
	goodsHdl := goods.NewGoodsHandler(goodsSvc)

	// Activity Template module
	activityTemplateRepo := persistence.NewActivityTemplateRepository(mysql.DB)
	activityTemplateSvc := activityTemplateService.NewService(activityTemplateRepo, mysql.DB)
	activityTemplateHdl := handler.NewActivityTemplateHandler(activityTemplateSvc)

	// Activity module
	activityRepo := persistence.NewActivityRepository(mysql.DB)
	activitySvc := activityService.NewService(activityRepo, activityTemplateRepo, mysql.DB)
	activityHdl := handler.NewActivityHandler(activitySvc)

	// Approval module
	approvalTypeRepo := approvalImpl.NewApprovalFlowTypeRepository(mysql.DB)
	approvalTemplateRepo := approvalImpl.NewApprovalFlowTemplateRepository(mysql.DB)
	approvalMgmtRepo := approvalImpl.NewApprovalFlowManagementRepository(mysql.DB)
	approvalNodeRepo := approvalImpl.NewApprovalNodeCaseRepository(mysql.DB)
	approvalDomainSvc := approvalDomainService.NewApprovalFlowService(approvalMgmtRepo, approvalNodeRepo, approvalTemplateRepo)
	approvalTypeSvc := approvalService.NewApprovalFlowTypeService(approvalTypeRepo)
	approvalTemplateSvc := approvalService.NewApprovalFlowTemplateService(approvalTemplateRepo, approvalTypeRepo)
	approvalMgmtSvc := approvalService.NewApprovalFlowManagementService(approvalMgmtRepo, approvalNodeRepo, approvalDomainSvc)
	approvalHdl := approval.NewApprovalHandler(approvalTypeSvc, approvalTemplateSvc, approvalMgmtSvc)

	// Placeholder handler for unimplemented features
	placeholderHdl := placeholder.NewPlaceholderHandler()

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (no JWT required)
		api.POST("/login", authHdl.Login)
		api.POST("/logout", authHdl.Logout)

		// Routes that require authentication
		authorized := api.Group("/")
		authorized.Use(middleware.JWTAuth())
		{
			// User info
			authorized.GET("/profile", authHdl.GetProfile)
			authorized.GET("/sync-role", authHdl.SyncRole)
			authorized.GET("/user/permissions", authHdl.GetUserPermissions)

			// Menu (for frontend navigation)
			authorized.GET("/menu", rbacHdl.GetMenu)
			authorized.GET("/menus", rbacHdl.GetMenu) // Compatible with frontend call

			// Role management
			roles := authorized.Group("/roles")
			{
				roles.GET("", rbacHdl.GetRoles)
				roles.POST("", rbacHdl.CreateRole)
				roles.PUT("/:id", rbacHdl.UpdateRole)
				roles.PUT("/:id/status", rbacHdl.UpdateRoleStatus)
				roles.GET("/:id/permissions", rbacHdl.GetRolePermissions)
				roles.PUT("/:id/permissions", rbacHdl.UpdateRolePermissions)
			}

			// Permission management
			permissions := authorized.Group("/permissions")
			{
				permissions.GET("", rbacHdl.GetPermissions)
				permissions.PUT("/:id/status", rbacHdl.UpdatePermissionStatus)
				permissions.GET("/tree", rbacHdl.GetPermissionTree)
			}

			// Menu management (updated route: menu-management -> menu_management)
			menus := authorized.Group("/menu_management")
			{
				menus.GET("", rbacHdl.GetMenus)
				menus.PUT("/:id", rbacHdl.UpdateMenu)
				menus.PUT("/:id/status", rbacHdl.UpdateMenuStatus)
			}

			// Menu management alias (前端使用连字符)
			menusAlias := authorized.Group("/menu-management")
			{
				menusAlias.GET("", rbacHdl.GetMenus)
				menusAlias.PUT("/:id", rbacHdl.UpdateMenu)
				menusAlias.PUT("/:id/status", rbacHdl.UpdateMenuStatus)
			}

			// Basic data (sex, grade, subject)
			authorized.GET("/sexes", basicHdl.GetAllSexes)
			authorized.GET("/grades/active", basicHdl.GetActiveGrades)
			authorized.GET("/subjects/active", basicHdl.GetActiveSubjects)

			// Account management
			accounts := authorized.Group("/accounts")
			{
				accounts.GET("", accountHdl.GetAccounts)
				accounts.POST("", accountHdl.CreateAccount)
				accounts.PUT("/:id", accountHdl.UpdateAccount)
				accounts.PUT("/:id/status", accountHdl.UpdateAccountStatus)
			}

			// Student Management
			students := authorized.Group("/students")
			{
				students.GET("/active", studentHdl.GetActiveStudents) // Must be before /:id
				students.GET("", studentHdl.GetStudents)
				students.POST("", studentHdl.CreateStudent)
				students.PUT("/:id", studentHdl.UpdateStudent)
				students.PUT("/:id/status", studentHdl.UpdateStudentStatus)
				students.DELETE("/:id", studentHdl.DeleteStudent)
			}

			// Coach Management
			coaches := authorized.Group("/coaches")
			{
				coaches.GET("/active", coachHdl.GetActiveCoaches) // Must be before /:id
				coaches.GET("", coachHdl.GetCoaches)
				coaches.POST("", coachHdl.CreateCoach)
				coaches.PUT("/:id", coachHdl.UpdateCoach)
				coaches.PUT("/:id/status", coachHdl.UpdateCoachStatus)
				coaches.DELETE("/:id", coachHdl.DeleteCoach)
			}

			// Order Management - Placeholder routes
			authorized.GET("/orders", placeholderHdl.HandlePlaceholder)
			authorized.GET("/childorders", placeholderHdl.HandlePlaceholder)
			authorized.GET("/refund_orders", placeholderHdl.HandlePlaceholder)
			authorized.GET("/refund_childorders", placeholderHdl.HandlePlaceholder)

			// Brand Management
			brands := authorized.Group("/brands")
			{
				brands.GET("", brandHdl.GetBrands)
				brands.GET("/active", brandHdl.GetActiveBrands)
				brands.POST("", brandHdl.CreateBrand)
				brands.PUT("/:id", brandHdl.UpdateBrand)
				brands.PUT("/:id/status", brandHdl.UpdateBrandStatus)
			}

			// Classify Management
			classifies := authorized.Group("/classifies")
			{
				classifies.GET("", classifyHdl.GetClassifies)
				classifies.GET("/parents", classifyHdl.GetParents)
				classifies.GET("/active", classifyHdl.GetActiveClassifies)
				classifies.POST("", classifyHdl.CreateClassify)
				classifies.PUT("/:id", classifyHdl.UpdateClassify)
				classifies.PUT("/:id/status", classifyHdl.UpdateClassifyStatus)
			}

			// Attribute Management
			attributes := authorized.Group("/attributes")
			{
				attributes.GET("", attributeHdl.GetAttributes)
				attributes.GET("/active", attributeHdl.GetActiveAttributes)
				attributes.POST("", attributeHdl.CreateAttribute)
				attributes.PUT("/:id", attributeHdl.UpdateAttribute)
				attributes.PUT("/:id/status", attributeHdl.UpdateAttributeStatus)
				attributes.GET("/:id/values", attributeHdl.GetAttributeValues)
				attributes.POST("/:id/values", attributeHdl.SaveAttributeValues)
			}

			// Goods Management (静态路径必须在参数路径之前)
			goods := authorized.Group("/goods")
			{
				goods.GET("", goodsHdl.GetGoods)
				goods.GET("/active-for-order", goodsHdl.GetActiveForOrder)
				goods.GET("/available-for-combo", goodsHdl.GetAvailableForCombo)
				goods.POST("", goodsHdl.CreateGoods)
				goods.GET("/:id", goodsHdl.GetGoodsByID)
				goods.GET("/:id/included-goods", goodsHdl.GetIncludedGoods)
				goods.GET("/:id/total-price", goodsHdl.GetTotalPrice)
				goods.PUT("/:id", goodsHdl.UpdateGoods)
				goods.PUT("/:id/status", goodsHdl.UpdateStatus)
			}

			// Approval Flow Management
			// 审批流类型管理
			approvalTypes := authorized.Group("/approval-flow-types")
			{
				approvalTypes.GET("", approvalHdl.GetApprovalFlowTypes)
				approvalTypes.POST("", approvalHdl.CreateApprovalFlowType)
				approvalTypes.PUT("/:id/status", approvalHdl.UpdateApprovalFlowTypeStatus)
			}

			// 审批流模板管理
			approvalTemplates := authorized.Group("/approval-flow-templates")
			{
				approvalTemplates.GET("", approvalHdl.GetApprovalFlowTemplates)
				approvalTemplates.GET("/:id", approvalHdl.GetApprovalFlowTemplateDetail)
				approvalTemplates.POST("", approvalHdl.CreateApprovalFlowTemplate)
				approvalTemplates.PUT("/:id/status", approvalHdl.UpdateApprovalFlowTemplateStatus)
			}

			// 审批流实例管理
			approvalFlows := authorized.Group("/approval-flows")
			{
				approvalFlows.GET("/initiated", approvalHdl.GetInitiatedFlows)
				approvalFlows.GET("/pending", approvalHdl.GetPendingFlows)
				approvalFlows.GET("/completed", approvalHdl.GetCompletedFlows)
				approvalFlows.GET("/copied", approvalHdl.GetCopiedFlows)
				approvalFlows.GET("/:id/detail", approvalHdl.GetApprovalFlowDetail)
				approvalFlows.POST("/create-from-template", approvalHdl.CreateFromTemplate)
				approvalFlows.PUT("/:id/cancel", approvalHdl.CancelApprovalFlow)
				approvalFlows.POST("/approve", approvalHdl.ApproveFlow)
				approvalFlows.POST("/reject", approvalHdl.RejectFlow)
			}

			// Approval Flow Management - menu placeholder
			authorized.GET("/approval_flow_type", placeholderHdl.HandlePlaceholder)
			authorized.GET("/approval_flow_template", placeholderHdl.HandlePlaceholder)
			authorized.GET("/approval_flow_management", placeholderHdl.HandlePlaceholder)

			// Activity Template Management (API)
			activityTemplates := authorized.Group("/activity-templates")
			{
				activityTemplates.GET("", activityTemplateHdl.ListTemplates)
				activityTemplates.GET("/active", activityTemplateHdl.ListActiveTemplates)
				activityTemplates.POST("", activityTemplateHdl.CreateTemplate)
				activityTemplates.GET("/:id", activityTemplateHdl.GetTemplate)
				activityTemplates.PUT("/:id", activityTemplateHdl.UpdateTemplate)
				activityTemplates.DELETE("/:id", activityTemplateHdl.DeleteTemplate)
				activityTemplates.PUT("/:id/status", activityTemplateHdl.UpdateTemplateStatus)
			}

			// Activity Management (API)
			activities := authorized.Group("/activities")
			{
				activities.GET("", activityHdl.ListActivities)
				activities.POST("", activityHdl.CreateActivity)
				activities.GET("/by-date-range", activityHdl.GetActivitiesByDateRange)
				activities.GET("/:id", activityHdl.GetActivity)
				activities.PUT("/:id", activityHdl.UpdateActivity)
				activities.DELETE("/:id", activityHdl.DeleteActivity)
				activities.PUT("/:id/status", activityHdl.UpdateActivityStatus)
			}

			// Marketing Management - Menu placeholder (前端菜单导航使用)
			authorized.GET("/activity_template", placeholderHdl.HandlePlaceholder)
			authorized.GET("/activity_management", placeholderHdl.HandlePlaceholder)

			// Contract Management
			contracts := authorized.Group("/contracts")
			{
				contracts.GET("", contractHdl.GetContracts)
				contracts.POST("", contractHdl.CreateContract)
				contracts.PUT("/:id/revoke", contractHdl.RevokeContract)
				contracts.PUT("/:id/terminate", contractHdl.TerminateContract)
				contracts.GET("/:id", contractHdl.GetContractByID)
			}

			// Contract Management - menu placeholder
			authorized.GET("/contract_management", placeholderHdl.HandlePlaceholder)

			// Finance Management - Placeholder routes
			authorized.GET("/payment_collection", placeholderHdl.HandlePlaceholder)
			authorized.GET("/separate_account", placeholderHdl.HandlePlaceholder)
			authorized.GET("/refund_management", placeholderHdl.HandlePlaceholder)
			authorized.GET("/refund_payment_detail", placeholderHdl.HandlePlaceholder)
		}
	}
}