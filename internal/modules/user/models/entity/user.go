package entity

import (
	"time"
)

type User struct {
	UserId       string      `json:"userId" bson:"userId"`
	FullName     string      `json:"fullName" bson:"fullName"`
	Email        string      `json:"email" bson:"email"`
	Password     string      `json:"password" bson:"password"`
	NIK          string      `json:"nik" bson:"nik"`
	MobileNumber string      `json:"mobileNumber" bson:"mobileNumber"`
	Address      string      `json:"address" bson:"address"`
	Subdistrict  Subdistrict `json:"subdistrict" bson:"subdistrict"`
	Country      Country     `json:"country" bson:"country"`
	RtRw         string      `json:"rtrw" bson:"rtrw"`
	Role         string      `json:"role" bson:"role"`
	KKNumber     string      `json:"kkNumber" bson:"kkNumber"`
	Status       string      `json:"status" bson:"status"`
	LoginAt      time.Time   `json:"loginAt" bson:"loginAt"`
	CreatedAt    time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt" bson:"updatedAt"`
}

// Move to domain address
type Subdistrict struct {
	Id           string `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
	DistrictId   string `json:"districtId" bson:"districtId"`
	DistrictName string `json:"districtName" bson:"districtName"`
	CityId       string `json:"cityId" bson:"cityId"`
	CityName     string `json:"cityName" bson:"cityName"`
	ProvinceId   string `json:"provinceId" bson:"provinceId"`
	ProvinceName string `json:"provinceName" bson:"provinceName"`
}

type Country struct {
	Id            int    `json:"id" bson:"id"`
	Code          string `json:"code" bson:"code"`
	Name          string `json:"name" bson:"name"`
	FullName      string `json:"fullName" bson:"fullName"`
	ContinentId   string `json:"continentId" bson:"continentId"`
	ContinentName string `json:"continentName" bson:"continentName"`
	Latitude      string `json:"latitude" bson:"latitude"`
	Longitude     string `json:"longitude" bson:"longitude"`
}
