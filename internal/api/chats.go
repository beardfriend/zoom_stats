package api

import (
	"context"
	"fmt"
	"time"

	"zoom_stats/internal/models"

	z "github.com/beardfriend/zoom_chatting_parser"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatAPI struct {
	db         *mongo.Client
	chatParser *z.Parser
}

func NewChatAPI(db *mongo.Client, chatParser *z.Parser) *ChatAPI {
	return &ChatAPI{db, chatParser}
}

type ListRequest struct {
	StartDateTime time.Time `form:"start_date_time" time_format:"2006-01-02T15:04:05"`
	EndDateTime   time.Time `form:"end_date_time" time_format:"2006-01-02T15:04:05"`
	Count         int       `form:"count"`
	Category      string    `form:"category"`
}

func (h *ChatAPI) GetList(c *gin.Context) {
	var query ListRequest

	if err := c.ShouldBindQuery(&query); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"msg": "error"})
		return
	}

	matchStage := bson.D{
		{"$match", bson.D{
			{"chatted_at", bson.D{
				{"$gte", query.StartDateTime},
				{"$lt", query.EndDateTime},
			}},
		}},
	}

	sortStage := bson.D{
		{"$sort", bson.D{{"react_ids", -1}}},
	}

	limitStage := bson.D{
		{"$limit", query.Count},
	}

	pipeline := mongo.Pipeline{matchStage, sortStage, limitStage}

	cursor, err := h.db.Database("zoom_chat").Collection("chats").Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		panic(err)
	}

	c.JSON(200, results)
}

func (h *ChatAPI) Upload(c *gin.Context) {
	// Form Data

	fileHeader, _ := c.FormFile("file")
	date := c.PostForm("date")

	// date Validation

	// File Access

	multipartFile, err := fileHeader.Open()
	if err != nil {
		panic(err)
	}

	result, err := h.chatParser.Parse(multipartFile)
	if err != nil {
		c.JSON(400, gin.H{"msg": "file check"})
		return
	}

	idMap := make(map[uint]primitive.ObjectID)
	for _, v := range result.ZoomChatHistory {
		idMap[v.Id] = primitive.NewObjectID()
	}

	DAO := make([]interface{}, 0)
	for _, v := range result.ZoomChatHistory {

		// make reactIds replyIds
		reactIds := make([]primitive.ObjectID, 0)
		replyIds := make([]primitive.ObjectID, 0)

		for id := range v.ReactIds {
			reactIds = append(reactIds, idMap[uint(id)])
		}

		for id := range v.ReplyIds {
			replyIds = append(replyIds, idMap[uint(id)])
		}

		// make chattedAt
		fmt.Println(v.Text)
		chattedAt, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", date, v.ChatedAt))
		DAO = append(DAO, &models.Chat{
			ID:        idMap[v.Id],
			ReactIds:  reactIds,
			ReplyIds:  replyIds,
			TextType:  v.TextType,
			Sender:    v.SenderName,
			Text:      v.Text,
			Removed:   v.Removed,
			ChattedAt: chattedAt,
		})
	}

	collection := h.db.Database("zoom_chat").Collection("chats")

	res, err := collection.InsertMany(context.TODO(), DAO)
	if err != nil {

		c.JSON(500, gin.H{"msg": "database error"})
		return
	}
	fmt.Println(res)
	c.JSON(201, gin.H{"msg": "created"})
}
