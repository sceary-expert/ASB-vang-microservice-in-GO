package function

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	// "example.com/core/config"
	"example.com/core/pkg/log"
	"example.com/function/database"
	micros "example.com/function/micros"
	"example.com/function/router"
)

// Cache state
var app *fiber.App

func init() {

	micros.InitConfig()

	// fmt.Println("here")

	// Initialize app
	app = fiber.New()
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(
		logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} - ${header:}\n​",
		},
	))
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     *config.AppConfig.Origin,
	// 	AllowCredentials: true,
	// 	AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Headers, X-Requested-With, X-HTTP-Method-Override, access-control-allow-origin, access-control-allow-headers",
	// }))
	router.SetupRoutes(app)
	// fmt.Println("finally here")
}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {
	// func Handle(){

	fmt.Println("Jai Bajrang Bali")
	ctx := context.Background()
	// fmt.Println("in handle function at line 54")
	if database.Db == nil {
		// fmt.Println("database.db is nil at line 58")
		var startErr error
		// fmt.Println("start error is empty now at line 60")
		startErr = database.Connect(ctx)
		// fmt.Println("starterr has been assigned with some value at line 62")
		if startErr != nil {
			// fmt.Println("start error at line 64")
			log.Error("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
		// fmt.Println(" no start error at line 62")
	}
	fmt.Println("in handle function at line 64")
	adaptor.FiberApp(app)(w, r)

}
