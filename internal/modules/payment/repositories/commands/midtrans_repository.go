package commands

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-service/configs"
	"payment-service/internal/modules/payment"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/modules/payment/models/response"
	"payment-service/internal/pkg/log"
	"time"
)

type midtransRepository struct {
	logger log.Logger
}

func NewCommandMidtransRepository(log log.Logger) payment.MidtransRepositoryCommand {
	return &midtransRepository{
		logger: log,
	}
}

func (m midtransRepository) TransferBank(ctx context.Context, payload request.BankTransferRequest) (*response.BankTransferResponse, error) {
	result := &response.BankTransferResponse{}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	bankTransferUrl := fmt.Sprintf("%s/v2/charge", configs.GetConfig().Midtrans.BaseUrl)

	request, err := http.NewRequest(http.MethodPost, bankTransferUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", configs.GetConfig().Midtrans.BasicAuth)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   180 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	errBody := json.Unmarshal(respBody, &result)
	if errBody != nil {
		return nil, errBody
	}

	return result, nil
}
