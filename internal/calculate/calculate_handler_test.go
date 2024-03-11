package calculate

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mayu13/gymshark-assignment/api"
	"github.com/mayu13/gymshark-assignment/internal/packs"
	"github.com/stretchr/testify/assert"
)

// Mock for Packs Manager
type mockPackManager struct {
	SetPackSizesFunc   func(sizes []int) error
	CalculatePacksFunc func(itemOrder int) ([]packs.Pack, error)
}

func (m *mockPackManager) SetPackSizes(sizes []int) error {
	return m.SetPackSizesFunc(sizes)
}

func (m *mockPackManager) CalculatePacks(itemOrder int) ([]packs.Pack, error) {
	return m.CalculatePacksFunc(itemOrder)
}

func TestSetPackSizesHandler(t *testing.T) {

	t.Run("successful request", func(t *testing.T) {
		pm := &mockPackManager{
			SetPackSizesFunc: func(sizes []int) error {
				return nil
			},
		}
		handler := SetPackSizesHandler(pm)

		reqBody, _ := json.Marshal(api.SetPackSizes{Sizes: []int{1, 2, 3}})
		req := httptest.NewRequest(http.MethodPost, "/packs", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("request with invalid body", func(t *testing.T) {
		pm := &mockPackManager{}
		handler := SetPackSizesHandler(pm)

		req := httptest.NewRequest(http.MethodPost, "/packs", io.NopCloser(bytes.NewBufferString("{invalid json")))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid JSON in request", func(t *testing.T) {
		pm := &mockPackManager{}
		handler := SetPackSizesHandler(pm)

		reqBody := []byte(`{"sizes": "not-an-array"}`)
		req := httptest.NewRequest(http.MethodPost, "/packs", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("failure in SetPackSizes method", func(t *testing.T) {
		pm := &mockPackManager{
			SetPackSizesFunc: func(sizes []int) error {
				return errors.New("internal error")
			},
		}
		handler := SetPackSizesHandler(pm)

		reqBody, _ := json.Marshal(api.SetPackSizes{Sizes: []int{1, 2, 3}})
		req := httptest.NewRequest(http.MethodPost, "/packs", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestCalculatePacksHandler(t *testing.T) {
	t.Run("successful calculation", func(t *testing.T) {
		pm := &mockPackManager{
			CalculatePacksFunc: func(itemsCount int) ([]packs.Pack, error) {
				return []packs.Pack{{Size: 5, Quantity: 2}}, nil
			},
		}
		handler := CalculatePacksHandler(pm)

		reqBody, _ := json.Marshal(api.CalculatePacksRequest{ItemsCount: 10})
		req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		var response api.CalculatePacksResponse
		json.Unmarshal(rr.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Len(t, response.Packs, 1)
		assert.Equal(t, 5, response.Packs[0].Size)
		assert.Equal(t, 2, response.Packs[0].Count)
	})

	t.Run("invalid request body", func(t *testing.T) {
		pm := &mockPackManager{}
		handler := CalculatePacksHandler(pm)

		req := httptest.NewRequest(http.MethodPost, "/calculate", ioutil.NopCloser(bytes.NewBufferString("{invalid json")))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("failure in CalculatePacks method", func(t *testing.T) {
		pm := &mockPackManager{
			CalculatePacksFunc: func(itemsCount int) ([]packs.Pack, error) {
				return nil, errors.New("calculation error")
			},
		}
		handler := CalculatePacksHandler(pm)

		reqBody, _ := json.Marshal(api.CalculatePacksRequest{ItemsCount: 10})
		req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
