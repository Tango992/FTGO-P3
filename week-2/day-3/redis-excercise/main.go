package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"redis-excercise/config"
	"redis-excercise/models"
)

func main() {
	cacheClient := config.InitCache()
	ctx := context.TODO()

	order := models.Order{Amount: 10000, Discount: 200}
	orderByte, err :=json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	
	if err := cacheClient.Set(ctx, "hello", orderByte, 0).Err(); err != nil {
		log.Fatal(err)
	}

	val, _ := cacheClient.Get(ctx, "hello").Result()
	
	var cacheOrder models.Order
	if err := json.Unmarshal([]byte(val), &cacheOrder); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(cacheOrder)

	cacheClient.Del(ctx, "hello")
}