package request

const (
	RoleUser        = `user`
	RoleAdmin       = `admin`
	RoleStackHolder = `stackholder`
)

var MapOfRole = map[string]string{
	RoleUser:        RoleUser,
	RoleAdmin:       RoleAdmin,
	RoleStackHolder: RoleStackHolder,
}

type RegisterUser struct {
	FullName      string `json:"full_name" validate:"required"`
	Email         string `json:"email" validate:"required,min=1,max=50"`
	Password      string `json:"password" validate:"required,min=8,max=20"`
	NIK           string `json:"nik" validate:"required"`
	MobileNumber  string `json:"mobile_number" validate:"required"`
	Address       string `json:"address"`
	ProvinceId    string `json:"province_id" validate:"required_if=CountryId 100"`
	CityId        string `json:"city_id" validate:"required_if=CountryId 100"`
	DistrictId    string `json:"district_id" validate:"required_if=CountryId 100"`
	SubdictrictId string `json:"subdictrict_id" validate:"required_if=CountryId 100"`
	CountryId     string `json:"country_id" validate:"required"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	RtRw          string `json:"rt_rw"`
	Role          string `json:"role" validate:"required"`
	KKNumber      string `json:"kk_number"`
}

type UpdateUser struct {
	FullName      string `json:"full_name" validate:"required"`
	MobileNumber  string `json:"mobile_number" validate:"required"`
	Address       string `json:"address"`
	ProvinceId    string `json:"province_id" validate:"required"`
	CityId        string `json:"city_id" validate:"required"`
	DistrictId    string `json:"district_id" validate:"required"`
	SubdictrictId string `json:"subdictrict_id" validate:"required"`
	RtRw          string `json:"rt_rw"`
	Role          string `json:"role"`
}

type VerifyRegisterUser struct {
	Email string `json:"email" validate:"required,min=1,max=50"`
	Otp   string `json:"otp_number" validate:"required"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required"`
}

type GetProfile struct {
	UserId string
}
