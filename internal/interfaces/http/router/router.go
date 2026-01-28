package router

import (
	accountService "charonoms/internal/application/service/account"
	authService "charonoms/internal/application/service/auth"
	basicService "charonoms/internal/application/service/basic"
	coachService "charonoms/internal/application/service/coach"
	contractService "charonoms/internal/application/service/contract"
	rbacService "charonoms/internal/application/service/rbac"
	studentService "charonoms/internal/application/service/student"
	"charonoms/internal/infrastructure/config"
	"charonoms/internal/infrastructure/persistence"
	accountImpl "charonoms/internal/infrastructure/persistence/mysql/account"
	authImpl "charonoms/internal/infrastructure/persistence/mysql/auth"
	coachImpl "charonoms/internal/infrastructure/persistence/mysql/coach"
	contractImpl "charonoms/internal/infrastructure/persistence/mysql/contract"
	rbacImpl "charonoms/internal/infrastructure/persistence/mysql/rbac"
	studentImpl "charonoms/internal/infrastructure/persistence/mysql/student"
	"charonoms/internal/infrastructure/persistence/mysql"
	"charonoms/internal/interfaces/http/handler/account"
	"charonoms/internal/interfaces/http/handler/auth"
	"charonoms/internal/interfaces/http/handler/basic"
	"charonoms/internal/interfaces/http/handler/coach"
	"charonoms/internal/interfaces/http/handler/contract"
	"charonoms/internal/interfaces/http/handler/placeholder"
	"charonoms/internal/interfaces/http/handler/rbac"
	"charonoms/internal/interfaces/http/handler/student"
	"charonoms/internal/interfaces/http/middleware"

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
	// Auth module
	authRepo := authImpl.NewAuthRepository(mysql.DB)
	authSvc := authService.NewAuthService(authRepo, cfg.JWT)
	authHdl := auth.NewAuthHandler(authSvc)

	// RBAC module
	roleRepo := rbacImpl.NewRoleRepository(mysql.DB)
	permissionRepo := rbacImpl.NewPermissionRepository(mysql.DB)
	menuRepo := rbacImpl.NewMenuRepository(mysql.DB)
	rbacSvc := rbacService.NewRBACService(roleRepo, permissionRepo, menuRepo)
	rbacHdl := rbac.NewRBACHandler(rbacSvc)

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

			// Goods Management - Placeholder routes
			authorized.GET("/brands", placeholderHdl.HandlePlaceholder)
			authorized.GET("/attributes", placeholderHdl.HandlePlaceholder)
			authorized.GET("/classifies", placeholderHdl.HandlePlaceholder)
			authorized.GET("/goods", placeholderHdl.HandlePlaceholder)

			// Approval Flow Management - Placeholder routes
			authorized.GET("/approval_flow_type", placeholderHdl.HandlePlaceholder)
			authorized.GET("/approval_flow_template", placeholderHdl.HandlePlaceholder)
			authorized.GET("/approval_flow_management", placeholderHdl.HandlePlaceholder)

			// Marketing Management - Placeholder routes
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