package constants

type MetaData struct {
	Page      int64 `json:"page"`
	Count     int64 `json:"count"`
	TotalPage int64 `json:"totalPage"`
	TotalData int64 `json:"totalData"`
}

var PaymentType = []string{"bca", "bni", "permata"}

const (
	SoldevInstitutionCode = "501"
	Pending               = "pending"
	Settlement            = "settlement"
	Expired               = "expired"
)
