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

func CreateCar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email, exists := ctx.Get("email")
		if !exists {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "You're not logged in."})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input models.CarInput
		user_email := fmt.Sprint(email)
		defer cancel()

		check := models.IsEmailRegistered(user_email, c)
		if !check {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You're now registering cars in Not-Registered user"})
			return
		}

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(&input); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		newCar := models.Car{
			Id:          primitive.NewObjectID(),
			Car_Type_Id: primitive.NewObjectID(),
			Car_Price:   input.Car_Price,
			Car_No:      input.Car_No,
			Owner_Email: user_email,
		}

		newCarType := models.CarType{
			Id:          primitive.NewObjectID(),
			Car_Brand:   input.Car_Brand,
			Car_Type:    input.Car_Type,
			Car_Seats:   input.Car_Seats,
			Car_Miles:   input.Car_Miles,
			Car_Gearbox: input.Car_Gearbox,
			Car_Fuel:    input.Car_Fuel,
		}

		_, err := newCar.SaveCar(c, newCarType)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": "Car Registered Successfully"})
	}
}

func GetCarProfileByEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		vemail, exists := ctx.Get("email")
		if !exists {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "You're not logged in."})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		email := fmt.Sprint(vemail)
		defer cancel()

		check := models.IsEmailRegistered(email, c)
		if !check {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You're now trying with Not-Registered user"})
			return
		}

		res, err := models.GetCarProfileByEmail(c, email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}
