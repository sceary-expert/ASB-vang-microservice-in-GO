package router

import (
	"vang/function/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Routers
	app.Post("/messages", handlers.SaveMessages)
	// app.Post("/message/query", handlers.QueryMessagesHandle)
	// app.Put("/message", handlers.UpdateMessageHandle)
	// app.Put("/room/deactive/:roomId", handlers.DeactiveUserRoomHandle)
	// app.Delete("/message/:messageId", handlers.DeleteMessageHandle)
	// app.Post("/room/active", handlers.ActivePeerRoom)

	// app.Get("/active-room/:roomId", handlers.GetActiveRoomHandle)
	// app.Post("/rooms", handlers.GetUserRooms)
	// app.Put("/read", handlers.UpdateReadMessageHandle)
}
