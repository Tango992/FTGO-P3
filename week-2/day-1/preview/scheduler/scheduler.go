package scheduler

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Scheduler struct {
	Cron *cron.Cron
	Db   *mongo.Database
}

func NewScheduler(db *mongo.Database) Scheduler {
	return Scheduler{
		Cron: cron.New(),
		Db:   db,
	}
}

func (s Scheduler) SettlePayment() {
	s.Cron.AddFunc("0 13 * * *", func() {
		result, err := s.Db.Collection("payments").UpdateMany(context.TODO(), bson.M{"status": "payment"}, bson.M{"$set": bson.M{"status": "settlement"}})
		if err != nil {
			fmt.Println("Cron job failed:", err.Error())
		}
		fmt.Printf("Cron scheduler updated %v entries\n", result.ModifiedCount)
	})
}

func (s Scheduler) StartCronJob() {
	s.SettlePayment()
	s.Cron.Start()
}
