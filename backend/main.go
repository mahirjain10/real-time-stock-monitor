package main

import (
	"log"
	_ "net/http/pprof" // Import pprof for profiling

	"github/mahirjain_10/sse-backend/backend/internal/app"
	"github/mahirjain_10/sse-backend/backend/internal/sse"

	"github.com/gin-gonic/gin"

	// "github/mahirjain_10/sse-backend/backend/internal/test"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/websocket"
	"github/mahirjain_10/sse-backend/backend/web/cmd/router"
)

func main() {

	// file, err := app.InitializeLogger()
	// if err != nil {
	// 	slog.Error("message", "Error initalizing logger", "error", err)
	// }
	// defer file.Close()
	// Start pprof server in a separate goroutine
	// go func() {
	// 	// log.Println("Starting pprof server on :6060")
	// 	slog.Info("Starting pprof server on :6060")
	// 	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
	// 		log.Fatalf("pprof server failed: %v", err)
	// 	}
	// }()

	// Initialize Gin router
	r := gin.Default()
	err := app.InitalizeEnv()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
		return
	}

	// ctx := context.Background()
	// Initialize the database and Redis client using the new helper function
	db, redisClient, err := app.InitializeServices()
	if err != nil {
		log.Fatalf("Error initializing services: %v", err)
		return
	}

	defer db.Close()
	sseServer := sse.NewSSEServer()

	// Intitialize FCM client
	fcmClient, err := app.InitializeFCMClient()
	if err != nil {
		log.Fatalf("Error initializing FCM client: %v", err)
	}
	var appInstance = types.App{
		DB:          db,
		RedisClient: redisClient,
		SSEServer:   sseServer,
		FCMClient: fcmClient,
	}

	// Keep the application running
	// Initialize database tables
	if err := app.InitializeDatabaseTables(db); err != nil {
		log.Fatalf("Error initializing database tables: %v", err)
		return
	}

	hub := websocket.NewHub()
	// c := cron.StartCron(&appInstance, hub)
	// defer c.Stop()
	// go hub.Run()

	r.GET("/events", func(c *gin.Context) {
		sse.SSEHandler(&appInstance, c.Writer, c.Request)
	})
	r.GET("/disconnect", func(c *gin.Context) {
		sse.StopMonitoringHandler(&appInstance, c.Writer, c.Request)
	})

	// Register routes
	// go utils.Subscribe(appInstance.RedisClient, ctx)
	// go utils.SubscribeToPubSub(appInstance.RedisClient, ctx, "alert-topic")

	router.RegisterRoutes(r, hub, &appInstance)
// use
	// go func() {
	// 	log.Println("Starting WebSocket Load Test...")
	// 	test.Wstest() // Call load test function
	// }()
	// log.Fatal(r.Run("192.168.245.128:8080"))

	log.Fatal(r.Run("0.0.0.0:8000"))

}
