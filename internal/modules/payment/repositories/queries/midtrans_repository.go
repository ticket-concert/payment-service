package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-service/configs"
	"payment-service/internal/modules/payment"
	"payment-service/internal/modules/payment/models/response"
	"payment-service/internal/pkg/log"
	"time"
)

var (
	NewRequest = http.NewRequest
	ReadAll    = io.ReadAll
)

type midtransRepository struct {
	baseUrl string
	logger  log.Logger
}

func NewQueryMidtransRepository(baseUrl string, log log.Logger) payment.MidtransRepositoryQuery {
	return &midtransRepository{
		baseUrl: baseUrl,
		logger:  log,
	}
}

func (m midtransRepository) GetTransactionStatus(ctx context.Context, transactionId string) (*response.TransactionStatusResponse, error) {
	result := &response.TransactionStatusResponse{}

	transactionStatusUrl := fmt.Sprintf("%s/v2/%s/status", m.baseUrl, transactionId)
	request, err := NewRequest(http.MethodGet, transactionStatusUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", configs.GetConfig().Midtrans.BasicAuth)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 300 {
		return nil, err
	}

	defer response.Body.Close()

	respBody, err := ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(respBody))

	errBody := json.Unmarshal(respBody, &result)
	if errBody != nil {
		return nil, errBody
	}

	return result, nil
}
