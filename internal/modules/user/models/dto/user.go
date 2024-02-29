package dto

import "time"

type UserResp struct {
	UserId       string    `json:"user_id" bson:"userId"`
	FullName     string    `json:"full_name" bson:"fullName"`
	Email        string    `json:"email" bson:"email"`
	MobileNumber string    `json:"mobileNumber" bson:"mobileNumber"`
	Role         string    `json:"role" bson:"role"`
	CreatedAt    time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updatedAt"`
}

type UserData struct {
	Data UserResp `json:"data" bson:"data"`
}
