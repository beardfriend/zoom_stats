package api

import (
	"context"
	"fmt"
	"strings"
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
	Count         int       `form:"count,default=10"`
}

type ReactorResponse struct {
	ID    string `json:"name" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}

func (h *ChatAPI) GetMostReactor(c *gin.Context) {
	var query ListRequest

	if err := c.ShouldBindQuery(&query); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"msg": "error"})
		return
	}

	splited := strings.Split(query.StartDateTime.String(), "+")[0]
	splited2 := strings.Split(query.EndDateTime.String(), "+")[0]
	startDateTime, _ := time.Parse("2006-01-02 15:04:05", splited[:len(splited)-1])
	endDateTime, _ := time.Parse("2006-01-02 15:04:05", splited2[:len(splited2)-1])

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "chatted_at", Value: bson.D{
				{Key: "$gte", Value: startDateTime},
				{Key: "$lt", Value: endDateTime},
			}},
			{
				Key: "text_type", Value: bson.D{
					{Key: "$eq", Value: 2},
				},
			},
		}},
	}

	groupByStage := bson.D{
		{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$sender"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
		},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}},
	}

	limitStage := bson.D{
		{Key: "$limit", Value: query.Count},
	}

	pipeline := mongo.Pipeline{matchStage, groupByStage, sortStage, limitStage}

	cursor, err := h.db.Database("zoom_chat").Collection("chats").Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}
	var results []ReactorResponse
	if err := cursor.All(context.Background(), &results); err != nil {
		panic(err)
	}

	c.JSON(200, results)
}

func (h *ChatAPI) GetMostTalkers(c *gin.Context) {
	var query ListRequest

	if err := c.ShouldBindQuery(&query); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"msg": "error"})
		return
	}

	matchStage := bson.D{
		bson.E{
			Key: "$match", Value: bson.D{
				bson.E{
					Key: "chatted_at", Value: bson.D{
						bson.E{Key: "$gte", Value: query.StartDateTime},
						bson.E{Key: "$lt", Value: query.EndDateTime},
					},
				},
				bson.E{
					Key: "$or", Value: bson.A{
						bson.D{
							bson.E{Key: "text_type", Value: 1},
						},
						bson.D{
							bson.E{Key: "text_type", Value: 3},
						},
					},
				},
			},
		},
	}

	groupByStage := bson.D{
		{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$sender"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
		},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}},
	}

	limitStage := bson.D{
		{Key: "$limit", Value: query.Count},
	}

	pipeline := mongo.Pipeline{matchStage, groupByStage, sortStage, limitStage}

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

func (h *ChatAPI) GetMostReactedPeople(c *gin.Context) {
	var query ListRequest

	if err := c.ShouldBindQuery(&query); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"msg": "error"})
		return
	}
	splited := strings.Split(query.StartDateTime.String(), "+")[0]
	splited2 := strings.Split(query.EndDateTime.String(), "+")[0]
	startDateTime, _ := time.Parse("2006-01-02 15:04:05", splited[:len(splited)-1])
	endDateTime, _ := time.Parse("2006-01-02 15:04:05", splited2[:len(splited2)-1])

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "chatted_at", Value: bson.D{
				{Key: "$gte", Value: startDateTime},
				{Key: "$lt", Value: endDateTime},
			}},
			bson.E{
				Key: "$or", Value: bson.A{
					bson.D{
						bson.E{Key: "text_type", Value: 1},
					},
					bson.D{
						bson.E{Key: "text_type", Value: 3},
					},
				},
			},
		}},
	}
	groupByStage := bson.D{
		{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$sender"},
				{Key: "count", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$size", Value: "$react_ids"},
					}},
				}},
			},
		},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}},
	}

	limitStage := bson.D{
		{Key: "$limit", Value: query.Count},
	}

	pipeline := mongo.Pipeline{matchStage, groupByStage, sortStage, limitStage}

	cursor, err := h.db.Database("zoom_chat").Collection("chats").Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}
	var results []ReactorResponse
	if err := cursor.All(context.Background(), &results); err != nil {
		panic(err)
	}

	c.JSON(200, results)
}

type ListResponse struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	ReactIds  []primitive.ObjectID `json:"react_ids" bson:"react_ids"`
	ReplyIds  []primitive.ObjectID `json:"reply_ids" bson:"reply_ids"`
	TextType  z.ChatContentType    `json:"text_type" bson:"text_type"`
	Sender    string               `json:"sender" bson:"sender"`
	Text      string               `json:"text" bson:"text"`
	Removed   bool                 `json:"removed" bson:"removed"`
	ChattedAt time.Time            `json:"chatted_at" bson:"chatted_at"`
}

func (h *ChatAPI) GetList(c *gin.Context) {
	var query ListRequest

	if err := c.ShouldBindQuery(&query); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"msg": "error"})
		return
	}
	splited := strings.Split(query.StartDateTime.String(), "+")[0]
	splited2 := strings.Split(query.EndDateTime.String(), "+")[0]
	startDateTime, _ := time.Parse("2006-01-02 15:04:05", splited[:len(splited)-1])
	endDateTime, _ := time.Parse("2006-01-02 15:04:05", splited2[:len(splited2)-1])

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "chatted_at", Value: bson.D{
				{Key: "$gte", Value: startDateTime},
				{Key: "$lt", Value: endDateTime},
			}},
		}},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{{Key: "react_ids", Value: -1}}},
	}

	limitStage := bson.D{
		{Key: "$limit", Value: query.Count},
	}

	pipeline := mongo.Pipeline{matchStage, sortStage, limitStage}

	cursor, err := h.db.Database("zoom_chat").Collection("chats").Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}
	var results []ListResponse

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
		chattedAt, _ := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%sT%s", date, v.ChatedAt))
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
