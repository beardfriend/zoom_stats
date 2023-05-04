package main

import (
	"context"

	"zoom_stats/internal/api"
	"zoom_stats/internal/database"

	z "github.com/beardfriend/zoom_chatting_parser"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	client, err := database.NewMongoDB(context.Background(), "mongodb://localhost:27017", "zoom_chat")
	if err != nil {
		panic(err)
	}
	chattingParser := z.NewParser()

	router := r.Group("/api")
	chatAPI := api.NewChatAPI(client, chattingParser)
	chat := router.Group("/chats")
	chat.POST("/upload", chatAPI.Upload)
	chat.GET("/", chatAPI.GetList)

	r.Run()
}
