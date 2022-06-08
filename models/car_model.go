package models

import (
	"context"
	"errors"
	"server/responses"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CarInput struct {
	Car_Brand   string  `json:"car_brand" binding:"required"`
	Car_Type    string  `json:"car_type" binding:"required"`
	Car_Seats   int     `json:"car_seats" binding:"required"`
	Car_Miles   float64 `json:"car_miles" binding:"required"`
	Car_Gearbox string  `json:"car_gearbox" binding:"required"`
	Car_Fuel    float64 `json:"car_fuel" binding:"required"`
	Car_Price   float64 `json:"car_price"`
	Car_No      string  `json:"car_no"`
}

type CarType struct {
	Id          primitive.ObjectID `json:"id"`
	Car_Brand   string             `json:"car_brand" binding:"required"`   // car brand
	Car_Type    string             `json:"car_type" binding:"required"`    // car type
	Car_Seats   int                `json:"car_seats" binding:"required"`   // seat count
	Car_Miles   float64            `json:"car_miles" binding:"required"`   // miles per hour : max speed
	Car_Gearbox string             `json:"car_gearbox" binding:"required"` // gearbox type
	Car_Fuel    float64            `json:"car_fuel" binding:"required"`    // capacity
}

type Car struct {
	Id          primitive.ObjectID `json:"id" binding:"required"`
	Car_Type_Id primitive.ObjectID `json:"car_type_id" binding:"required"`
	Car_Price   float64            `json:"car_price" binding:"required"` // hourly
	Car_No      string             `json:"car_no" binding:"required"`    // car id number
	Owner_Email string             `json:"owner_email"`                  // owner email
}

func IsCarNoExist(c context.Context, no string) bool {
	car := Car{}
	err := carCollection.FindOne(c, bson.M{"car_no": no}).Decode(&car)
	return err == nil
}

func (newCar *Car) SaveCar(c context.Context, newType CarType) (*Car, error) {

	if IsCarNoExist(c, newCar.Car_No) {
		return nil, errors.New("car with the same no already exists")
	}

	_, errs := newType.SaveCarType(c)
	if errs != nil && errs.Error() != "exist" {
		return nil, errs
	}

	id, erra := FindCarTypeId(c, newType)
	if erra != nil {
		return nil, erra
	}

	newCar.Car_Type_Id = id

	_, err := carCollection.InsertOne(c, newCar)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("car with the same no already exists")
		}
	}

	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"car_no": 1}, Options: opt}

	if _, err := carCollection.Indexes().CreateOne(c, index); err != nil {
		return nil, errors.New("could not create index of car number")
	}

	return nil, nil
}

func (newCarType *CarType) SaveCarType(c context.Context) (*CarType, error) {
	// validate if cartype is existing.
	var ntp CarType
	res := cartypeCollection.FindOne(c, bson.M{
		"car_brand":   newCarType.Car_Brand,
		"car_type":    newCarType.Car_Type,
		"car_seats":   newCarType.Car_Seats,
		"car_miles":   newCarType.Car_Miles,
		"car_gearbox": newCarType.Car_Gearbox,
		"car_fuel":    newCarType.Car_Fuel,
	}).Decode(&ntp)

	if res != nil {
		_, err := cartypeCollection.InsertOne(c, newCarType)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return nil, errors.New("exist")
}

func FindCarTypeId(c context.Context, newType CarType) (primitive.ObjectID, error) {
	var resType CarType
	res := cartypeCollection.FindOne(c, bson.M{
		"car_brand":   newType.Car_Brand,
		"car_type":    newType.Car_Type,
		"car_seats":   newType.Car_Seats,
		"car_miles":   newType.Car_Miles,
		"car_gearbox": newType.Car_Gearbox,
		"car_fuel":    newType.Car_Fuel,
	}).Decode(&resType)

	if res != nil {
		return primitive.NewObjectID(), res
	}
	return resType.Id, nil
}

func GetCarTypeById(c context.Context, id primitive.ObjectID) (*CarType, error) {
	var restype CarType
	res := cartypeCollection.FindOne(c, bson.M{"id": id}).Decode(&restype)

	if res != nil {
		return nil, res
	}

	return &restype, nil
}

func GetCarProfileByEmail(c context.Context, email string) ([]responses.Profile_Car_Response, error) {
	cur, err := carCollection.Find(c, bson.M{"owner_email": email})
	if err != nil {
		return nil, err
	}

	defer cur.Close(c)

	var result []responses.Profile_Car_Response

	for cur.Next(c) {
		var res Car
		err := cur.Decode(&res)
		if err != nil {
			continue
		}

		// need to aggregate
		curType, merr := GetCarTypeById(c, res.Car_Type_Id)
		if merr != nil {
			continue
		}

		rating, err := GetCarRating(c, res.Car_No)
		if err != nil {
			rating = 0.0
		}

		curCar := responses.Profile_Car_Response{
			Car_Brand:   curType.Car_Brand,
			Car_Type:    curType.Car_Type,
			Car_Seats:   curType.Car_Seats,
			Car_Miles:   curType.Car_Miles,
			Car_Gearbox: curType.Car_Gearbox,
			Car_No:      res.Car_No,
			Car_Price:   res.Car_Price,
			Car_Rating:  rating,
		}
		result = append(result, curCar)
	}

	return result, nil
}
