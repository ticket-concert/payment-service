package commands_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/modules/payment/repositories/commands"
	"testing"

	mocklog "payment-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
)

func TestTransferBank(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	repo := commands.NewCommandMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.TransferBank(context.Background(), request.BankTransferRequest{})

	assert.NoError(t, err)
}

func TestTransferBankErrMarshal(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	commands.Marshal = func(v any) ([]byte, error) {
		return nil, errors.New("error")
	}

	repo := commands.NewCommandMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.TransferBank(context.Background(), request.BankTransferRequest{})

	assert.Error(t, err)
}

func TestTransferBankErrRequest(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	commands.NewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("error NewRequest")
	}

	repo := commands.NewCommandMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.TransferBank(context.Background(), request.BankTransferRequest{})

	assert.Error(t, err)
}

func TestTransferBankErrDo(t *testing.T) {

	repo := commands.NewCommandMidtransRepository("URL", &mocklog.Logger{})

	_, err := repo.TransferBank(context.Background(), request.BankTransferRequest{})
	assert.Error(t, err)
}

func TestTransferBankErrReadAll(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	commands.ReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("error ReadAll")
	}

	repo := commands.NewCommandMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.TransferBank(context.Background(), request.BankTransferRequest{})

	assert.Error(t, err)
}

func TestTransferBankErrUnMarshal(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	repo := commands.NewCommandMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.TransferBank(context.Background(), request.BankTransferRequest{})

	assert.Error(t, err)
}
