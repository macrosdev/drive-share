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

func CreateReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, exists := ctx.Get("email")
		if !exists {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "You're not logged in."})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input models.ReviewInput
		defer cancel()

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(&input); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		newReview := models.Review{
			Id:        primitive.NewObjectID(),
			From:      fmt.Sprint(email),
			To:        input.To,
			Content:   input.Content,
			Rating:    input.Rating,
			Avatars:   input.Avatars,
			ReviewdAt: time.Now(),
		}

		_, err := newReview.SaveReview(c)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": "Reviewd successfully"})
	}
}
