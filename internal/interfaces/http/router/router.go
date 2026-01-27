package router

import (
	authService "charonoms/internal/application/service/auth"
	basicService "charonoms/internal/application/service/basic"
	rbacService "charonoms/internal/application/service/rbac"
	"charonoms/internal/infrastructure/config"
	"charonoms/internal/infrastructure/persistence"
	authImpl "charonoms/internal/infrastructure/persistence/mysql/auth"
	rbacImpl "charonoms/internal/infrastructure/persistence/mysql/rbac"
	"charonoms/internal/infrastructure/persistence/mysql"
	"charonoms/internal/interfaces/http/handler/auth"
	"charonoms/internal/interfaces/http/handler/basic"
	"charonoms/internal/interfaces/http/handler/rbac"
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

			// Menu management
			menus := authorized.Group("/menu-management")
			{
				menus.GET("", rbacHdl.GetMenus)
				menus.PUT("/:id", rbacHdl.UpdateMenu)
				menus.PUT("/:id/status", rbacHdl.UpdateMenuStatus)
			}

			// Basic data (sex, grade, subject)
			authorized.GET("/sexes", basicHdl.GetAllSexes)
			authorized.GET("/grades/active", basicHdl.GetActiveGrades)
			authorized.GET("/subjects/active", basicHdl.GetActiveSubjects)

			// TODO: Add other business module routes here
			// Student management
			// students := authorized.Group("/students")
			// students.Use(middleware.Permission("student"))
			// {
			//     students.GET("", studentHdl.List)
			//     students.POST("", studentHdl.Create)
			//     students.GET("/:id", studentHdl.Get)
			//     students.PUT("/:id", studentHdl.Update)
			//     students.DELETE("/:id", studentHdl.Delete)
			// }

			// ... other business module routes
		}
	}
}