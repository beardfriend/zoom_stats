package main

import (
	"context"

	"zoom_stats/internal/api"
	"zoom_stats/internal/database"

	z "github.com/beardfriend/zoom_chatting_parser"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

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
	chat.GET("/most-talkers", chatAPI.GetMostTalkers)
	chat.GET("/participation", chatAPI.GetParticipation)
	chat.GET("/most-reacted-people", chatAPI.GetMostReactedPeople)
	chat.GET("/most-reactors", chatAPI.GetMostReactor)

	r.Run()
}
