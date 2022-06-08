package responses

type Profile_Car_Response struct {
	Car_Brand   string  `json:"car_brand" binding:"required"`
	Car_Type    string  `json:"car_type" binding:"required"`
	Car_Seats   int     `json:"car_seats" binding:"required"`
	Car_Miles   float64 `json:"car_miles" binding:"required"`
	Car_Gearbox string  `json:"car_gearbox" binding:"required"`
	Car_No      string  `json:"car_no" binding:"required"`
	Car_Price   float64 `json:"car_price" binding:"required"`
	Car_Rating  float64 `json:"car_rating" binding:"required"`
}
