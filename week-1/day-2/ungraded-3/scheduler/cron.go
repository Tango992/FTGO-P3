package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"ungraded-3/models"

	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Schedule struct {
	Cron       *cron.Cron
	Collection *mongo.Collection
}

func New(collection *mongo.Collection) Schedule {
	return Schedule{
		Cron:       cron.New(),
		Collection: collection,
	}
}

func (s Schedule) PostMessage() {
	s.Cron.AddFunc("@every 3s", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		url := "https://api.api-ninjas.com/v1/jokes"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("X-Api-Key", os.Getenv("API_KEY"))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		var contentTmp []map[string]any
		if err := json.NewDecoder(res.Body).Decode(&contentTmp); err != nil {
			log.Fatal(err)
		}

		content := contentTmp[0]["joke"].(string)

		postData := models.Message{
			Sender:   "john@mail.com",
			Receiver: "jane@mail.com",
			Type:     "text",
			Content:  content,
		}

		mongoRes, err := s.Collection.InsertOne(ctx, postData)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s[%v] %sPosted message (ObjectID %s)\n", "\033[0m", time.Now().Local().Format("2006-01-02 15:04:05"), "\033[36m", mongoRes.InsertedID.(primitive.ObjectID).Hex())
	})
}

func (s Schedule) ClearMessages() {
	s.Cron.AddFunc("@every 15s", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := s.Collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s[%v] %sDeleted %d datas\n", "\033[0m", time.Now().Local().Format("2006-01-02 15:04:05"), "\033[33m", res.DeletedCount)
	})
}

func (s Schedule) StarWorker() {
	s.PostMessage()
	s.ClearMessages()
	s.Cron.Start()
}
