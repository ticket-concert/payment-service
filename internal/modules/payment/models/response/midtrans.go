package response

type VaNumber struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}

type BankTransferResponse struct {
	StatusCode        string     `json:"status_code"`
	StatusMessage     string     `json:"status_message"`
	TransactionID     string     `json:"transaction_id"`
	OrderID           string     `json:"order_id"`
	MerchantID        string     `json:"merchant_id"`
	GrossAmount       string     `json:"gross_amount"`
	Currency          string     `json:"currency"`
	PaymentType       string     `json:"payment_type"`
	SignatureKey      string     `json:"signature_key"`
	TransactionTime   string     `json:"transaction_time"`
	TransactionStatus string     `json:"transaction_status"`
	FraudStatus       string     `json:"fraud_status"`
	PermataVaNumber   string     `json:"permata_va_number"`
	VaNumbers         []VaNumber `json:"va_numbers"`
	ExpiryTime        string     `json:"expiry_time"`
}

type PaymentAmount struct {
	Amount string `json:"amount"`
	PaidAt string `json:"paid_at"`
}

type TransactionStatusResponse struct {
	StatusCode        string          `json:"status_code"`
	TransactionID     string          `json:"transaction_id"`
	GrossAmount       string          `json:"gross_amount"`
	Currency          string          `json:"currency"`
	OrderID           string          `json:"order_id"`
	PaymentType       string          `json:"payment_type"`
	SignatureKey      string          `json:"signature_key"`
	TransactionStatus string          `json:"transaction_status"`
	FraudStatus       string          `json:"fraud_status"`
	StatusMessage     string          `json:"status_message"`
	MerchantID        string          `json:"merchant_id"`
	VANumbers         []VaNumber      `json:"va_numbers"`
	PaymentAmounts    []PaymentAmount `json:"payment_amounts"`
	TransactionTime   string          `json:"transaction_time"`
	ExpiryTime        string          `json:"expiry_time"`
}
