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

type midtransRepository struct {
	logger log.Logger
}

func NewQueryMidtransRepository(log log.Logger) payment.MidtransRepositoryQuery {
	return &midtransRepository{
		logger: log,
	}
}

func (m midtransRepository) GetTransactionStatus(ctx context.Context, transactionId string) (*response.TransactionStatusResponse, error) {
	result := &response.TransactionStatusResponse{}

	transactionStatusUrl := fmt.Sprintf("%s/v2/%s/status", configs.GetConfig().Midtrans.BaseUrl, transactionId)

	request, err := http.NewRequest(http.MethodGet, transactionStatusUrl, nil)
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

	fmt.Println(response)

	if response.StatusCode > 300 {
		return nil, err
	}

	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(respBody))

	errBody := json.Unmarshal(respBody, &result)
	fmt.Println(errBody)
	if errBody != nil {
		return nil, errBody
	}

	return result, nil
}
