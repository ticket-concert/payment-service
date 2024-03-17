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

var (
	NewRequest = http.NewRequest
	ReadAll    = io.ReadAll
	Marshal    = json.Marshal
)

type midtransRepository struct {
	baseUrl string
	logger  log.Logger
}

func NewCommandMidtransRepository(baseUrl string, log log.Logger) payment.MidtransRepositoryCommand {
	return &midtransRepository{
		baseUrl: baseUrl,
		logger:  log,
	}
}

func (m midtransRepository) TransferBank(ctx context.Context, payload request.BankTransferRequest) (*response.BankTransferResponse, error) {
	result := &response.BankTransferResponse{}

	body, err := Marshal(payload)
	if err != nil {
		return nil, err
	}

	bankTransferUrl := fmt.Sprintf("%s/v2/charge", m.baseUrl)

	request, err := NewRequest(http.MethodPost, bankTransferUrl, bytes.NewBuffer(body))
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

	respBody, err := ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	errBody := json.Unmarshal(respBody, &result)
	if errBody != nil {
		return nil, errBody
	}

	return result, nil
}
