package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderInput struct {
	CarNo      string    `json:"car_no" binding:"required"`
	From_Time  time.Time `json:"from_time" binding:"required"`
	To_Time    time.Time `json:"to_time" binding:"required"`
	From_Where string    `json:"from_where" binding:"required"`
	To_Where   string    `json:"from_to" binding:"required"`
}

type Order struct {
	Id         primitive.ObjectID `json:"id" binding:"required"`
	UserEmail  string             `json:"user_email" binding:"required"` // Id of user who makes order
	CarNo      string             `json:"car_no" binding:"required"`     // car id number what user orders
	From_Time  time.Time          `json:"from_time" binding:"required"`  // pick up time
	To_Time    time.Time          `json:"to_time" binding:"required"`    // drop off time
	From_Where string             `json:"from_where" binding:"required"` // pick up place
	To_Where   string             `json:"to_where" binding:"required"`   // drop off place
}

func (newOrder *Order) SaveOrder(c context.Context) (*Order, error) {
	_, err := orderCollection.InsertOne(c, newOrder)
	return nil, err
}
