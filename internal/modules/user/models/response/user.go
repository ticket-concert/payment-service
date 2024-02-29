package response

type RegisterUser struct {
	Email string `json:"email"`
}

type VerifyRegister struct {
	AuthToken    string `json:"authToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiredAt    string `json:"expiredAt"`
}

type LoginUserResp struct {
	AuthToken    string `json:"authToken" bson:"authToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
	ExpiredAt    string `json:"expiredAt" bson:"password"`
}

type GetProfile struct {
	UserId       string `json:"user_id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	NIK          string `json:"nik"`
	MobileNumber string `json:"mobile_number"`
	Address      string `json:"address"`
	RtRw         string `json:"rt_rw"`
	Role         string `json:"role"`
	KKNumber     string `json:"kk_number"`
}
