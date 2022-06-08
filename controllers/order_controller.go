package controllers

import (
	"context"
	"fmt"
	"net/http"
	"server/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, exists := ctx.Get("email")
		if !exists {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "You're not logged in."})
			return
		}

		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input models.OrderInput
		defer cancel()

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(&input); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		newOrder := models.Order{
			Id:         primitive.NewObjectID(),
			UserEmail:  fmt.Sprint(email),
			CarNo:      input.CarNo,
			From_Time:  input.From_Time,
			To_Time:    input.To_Time,
			From_Where: input.From_Where,
			To_Where:   input.To_Where,
		}

		_, err := newOrder.SaveOrder(c)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": "Created Order Succesfully"})
	}
}
