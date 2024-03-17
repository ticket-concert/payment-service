package queries_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"payment-service/internal/modules/payment/repositories/queries"
	"payment-service/internal/pkg/errors"
	"testing"

	mocklog "payment-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
)

func TestGetTransactionStatus(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	repo := queries.NewQueryMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.GetTransactionStatus(context.Background(), "transactionId")

	assert.NoError(t, err)
}

func TestGetTransactionStatusErrReadAll(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	queries.ReadAll = func(r io.Reader) ([]byte, error) {
		return nil, errors.BadRequest("error ReadAll")
	}

	repo := queries.NewQueryMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.GetTransactionStatus(context.Background(), "transactionId")

	assert.Error(t, err)
}

func TestGetTransactionStatusErrDo(t *testing.T) {
	repo := queries.NewQueryMidtransRepository("URL", &mocklog.Logger{})

	_, err := repo.GetTransactionStatus(context.Background(), "transactionId")

	assert.Error(t, err)
}

func TestGetTransactionStatusErrStatusCode(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"failed"}`))
	}))
	defer mockServer.Close()

	repo := queries.NewQueryMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.GetTransactionStatus(context.Background(), "transactionId")

	assert.NoError(t, err)
}

func TestGetTransactionStatusErrResp(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	repo := queries.NewQueryMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.GetTransactionStatus(context.Background(), "transactionId")

	assert.Error(t, err)
}

func TestGetTransactionStatusErrHttp(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	queries.NewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return &http.Request{}, errors.BadRequest("error")
	}

	repo := queries.NewQueryMidtransRepository(mockServer.URL, &mocklog.Logger{})

	_, err := repo.GetTransactionStatus(context.Background(), "transactionId")

	assert.Error(t, err)
}
