package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewInput struct {
	To      string   `json:"to_car_no" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Rating  float64  `json:"rating" binding:"required"`
	Avatars []string `json:"avatars"`
}

type Review struct {
	Id        primitive.ObjectID `json:"id" binding:"required"`
	From      string             `json:"from_user_email" binding:"required"` // id of user who reviews
	To        string             `json:"to_car_no" binding:"required"`       // id of car what user reviews
	Content   string             `json:"content" binding:"required"`         // content of review
	Rating    float64            `json:"rating" binding:"required"`          // rating
	Avatars   []string           `json:"avatars"`                            // review pictures of car
	ReviewdAt time.Time          `json:"reviewd_at" binding:"required"`      // reviewd time
}

func (newReview *Review) SaveReview(c context.Context) (*Review, error) {
	_, err := reviewCollection.InsertOne(c, newReview)
	return nil, err
}

func GetCarRating(c context.Context, car_no string) (float64, error) {
	cur, err := reviewCollection.Find(c, bson.M{"to": car_no})
	if err != nil {
		return 0, err
	}

	defer cur.Close(c)
	sum, cnt := 0.0, 0.0

	for cur.Next(c) {
		var res Review
		err := cur.Decode(&res)
		if err != nil {
			continue
		}
		sum = sum + res.Rating
		cnt = cnt + 1
	}

	if cnt < 1 {
		return 0, errors.New("zero divisoin error")
	}

	return sum / cnt, nil
}
