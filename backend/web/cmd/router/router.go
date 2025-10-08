package router

import (
	"github/mahirjain_10/sse-backend/backend/internal/middleware"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/websocket"
	"github/mahirjain_10/sse-backend/backend/web/cmd/handlers/alert"
	"github/mahirjain_10/sse-backend/backend/web/cmd/handlers/auth"

	"github.com/gin-gonic/gin"
	// "github/mahirjain_10/sse-backend/backend/internal/test"
)

// registerRoutes handles the grouping and organization of routes
func RegisterRoutes(r *gin.Engine, hub *websocket.Hub, app *types.App) {
	// WebSocket endpoint
	// r.GET("/ws/get-stock-price-socket", func(c *gin.Context) {
	// 	websocket.ServeWs(c, hub, c.Writer, c.Request)
	// })

	// Auth group
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", func(c *gin.Context) {
			auth.RegisterUser(c, r, app)
		})
		authRoutes.POST("/login", func(c *gin.Context) {
			auth.LoginUser(c, r, app)
		})
	}

	// Alert group
	alertRoutes := r.Group("/api/alert")
	{
		// in making
		// alertRoutes.GET("/load-monitor-active-stocks", func(c *gin.Context) {
		// 	alert.LoadMonitorActiveStocks(c, r, app)
		// })
		alertRoutes.POST("/get-current-price", func(c *gin.Context) {
			alert.GetCurrentStockPriceAndTime(c, r, app)
		})
		alertRoutes.POST("/create-stock-alert", middleware.AuthMiddleware(app), func(c *gin.Context) {
			alert.CreateStockAlert(c, r, app)
		})
		alertRoutes.PUT("/update-stock-alert", func(c *gin.Context) {
			alert.UpdateStockAlert(c, r, app)
		})
		alertRoutes.PUT("/update-stock-alert-status", func(c *gin.Context) {
			alert.UpdateActiveStatus(c, r, app)
		})
		alertRoutes.DELETE("/delete-stock-alert", func(c *gin.Context) {
			alert.DeleteStockAlert(c, r, app)
		})
		alertRoutes.POST("/alert-notification", func(c *gin.Context) {
			alert.SendAlertNotification(c, r, app, hub)
		})
		alertRoutes.POST("/start-monitoring", func(c *gin.Context) {
			alert.StartStockAlertMonitoring(c, r, app)
		})
		alertRoutes.POST("/stop-monitoring", func(c *gin.Context) {
			alert.StopStockAlertMonitoring(c, r, app, hub)
		})

		alertRoutes.POST("/send-fcm-notification", func(c *gin.Context) {
			alert.SendFCMNotification(c, r, app)
		})

	}

	monitorStockRoutes := r.Group("api/monitor-stock")
	{
		monitorStockRoutes.POST("/start-monitoring", middleware.AuthMiddleware(app), func(c *gin.Context) {
			alert.ManualStartStockMonitoring(c, r, app)
		})
		monitorStockRoutes.POST("/stop-monitoring", middleware.AuthMiddleware(app), func(c *gin.Context) {
			alert.ManualStopStockMonitoring(c, r, app)
		})
	}

	// Load/Unload group
	loadUnloadRoutes := r.Group("/api/load-unload")
	{
		loadUnloadRoutes.POST("/load-stock-alerts", middleware.AuthMiddleware(app), func(c *gin.Context) {
			alert.LoadUserAlerts(c, app) // Assuming LoadUserAlerts is in the alert package
		})
		loadUnloadRoutes.POST("/unload-stock-alerts", middleware.AuthMiddleware(app), func(c *gin.Context) {
			alert.UnloadUserAlerts(c, app) // Assuming UnloadUserAlerts is in the alert package
		})
	}
}

