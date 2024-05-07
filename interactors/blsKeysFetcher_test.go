package interactors

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewBLSKeysFetcher(t *testing.T) {
	t.Parallel()

	t.Run("nil HTTP client wrapper should error", func(t *testing.T) {
		instance, err := NewBLSKeysFetcher(nil, nil, 0)
		assert.Nil(t, instance)
		assert.Equal(t, errNilHTTPClientWrapper, err)
	})
	t.Run("should work", func(t *testing.T) {
		instance, err := NewBLSKeysFetcher(&mock.HTTPClientWrapperStub{}, nil, 0)
		assert.NotNil(t, instance)
		assert.Nil(t, err)
	})
}

func TestBlsKeysFetcher_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *blsKeysFetcher
	assert.True(t, instance.IsInterfaceNil())

	instance = &blsKeysFetcher{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestBlsKeysFetcher_GetAllBLSKeys(t *testing.T) {
	t.Parallel()

	addr1 := "address1"
	addr2 := "address2"
	addresses := []core.Address{
		{
			Hex:    addr1,
			Bech32: addr1,
		},
		{
			Hex:    addr2,
			Bech32: addr2,
		},
	}

	bls1 := bytes.Repeat([]byte("1"), core.BLSKeyLen)
	bls2 := bytes.Repeat([]byte("2"), core.BLSKeyLen)
	bls3 := bytes.Repeat([]byte("3"), core.BLSKeyLen)
	bls4 := bytes.Repeat([]byte("4"), core.BLSKeyLen)

	resp := map[string][][]byte{
		addr1: {bls1, []byte(stakedStatus), bls2, []byte(stakedStatus)},
		addr2: {bls3, []byte(stakedStatus), bls4, []byte("not-staked")},
	}

	badResp := map[string][][]byte{
		addr1: {[]byte("not a bls key"), []byte(stakedStatus), bls1, []byte(stakedStatus), bls2, []byte(stakedStatus), []byte("extra, non wanted param")},
		addr2: {bls3, []byte(stakedStatus), bls4, []byte("not-staked")},
	}

	expectedErr := errors.New("expected error")
	t.Run("empty provided address should return empty list", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				assert.Fail(t, "should have not called post")

				return nil, http.StatusInternalServerError, nil
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, nil, 0)
		list, err := instance.GetAllBLSKeys(context.Background(), "test")
		assert.Nil(t, err)
		assert.Empty(t, list)
		assert.NotNil(t, list)
	})
	t.Run("not a valid address should not call post", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				assert.Fail(t, "should have not called post")

				return nil, http.StatusInternalServerError, nil
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, []core.Address{{}}, 0)
		list, err := instance.GetAllBLSKeys(context.Background(), "test")
		assert.Nil(t, err)
		assert.Empty(t, list)
		assert.NotNil(t, list)
	})
	t.Run("server responds not-ok should error", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				return nil, http.StatusInternalServerError, nil
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, addresses, 0)
		list, err := instance.GetAllBLSKeys(context.Background(), "test")
		assert.ErrorIs(t, err, errReturnCodeIsNotOk)
		assert.Empty(t, list)
	})
	t.Run("server errors should error", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				return nil, http.StatusInternalServerError, expectedErr
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, addresses, 0)
		list, err := instance.GetAllBLSKeys(context.Background(), "test")
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, list)
	})
	t.Run("server returns invalid response should error", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				return []byte("not-a-json"), http.StatusOK, nil
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, addresses, 0)
		list, err := instance.GetAllBLSKeys(context.Background(), "test")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid character")
		assert.Empty(t, list)
	})
	t.Run("should skip bad data", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				request := &vmQueryRequest{}
				err := json.Unmarshal(data, request)
				assert.Nil(t, err)

				assert.Equal(t, getBlsKeysStatusFuncName, request.FuncName)
				assert.Equal(t, validatorScAddress, request.ScAddress)
				assert.Equal(t, validatorScAddress, request.Caller)

				respSlices := badResp[request.Args[0]]

				response := &vmQueryResponse{}
				response.Data.Data.ReturnData = respSlices

				respBytes, err := json.Marshal(response)
				assert.Nil(t, err)

				return respBytes, http.StatusOK, nil
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, addresses, 0)
		list, err := instance.GetAllBLSKeys(context.Background(), "test")
		assert.Nil(t, err)

		expectedList := []string{
			hex.EncodeToString(bls1),
			hex.EncodeToString(bls2),
			hex.EncodeToString(bls3),
		}
		assert.Equal(t, expectedList, list)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				request := &vmQueryRequest{}
				err := json.Unmarshal(data, request)
				assert.Nil(t, err)

				assert.Equal(t, getBlsKeysStatusFuncName, request.FuncName)
				assert.Equal(t, validatorScAddress, request.ScAddress)
				assert.Equal(t, validatorScAddress, request.Caller)

				respSlices := resp[request.Args[0]]

				response := &vmQueryResponse{}
				response.Data.Data.ReturnData = respSlices

				respBytes, err := json.Marshal(response)
				assert.Nil(t, err)

				return respBytes, http.StatusOK, nil
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, addresses, 0)
		list, err := instance.GetAllBLSKeys(context.Background(), "test")
		assert.Nil(t, err)

		expectedList := []string{
			hex.EncodeToString(bls1),
			hex.EncodeToString(bls2),
			hex.EncodeToString(bls3),
		}
		assert.Equal(t, expectedList, list)
	})
	t.Run("context expires", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			PostHTTPHandler: func(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
				request := &vmQueryRequest{}
				err := json.Unmarshal(data, request)
				assert.Nil(t, err)

				assert.Equal(t, getBlsKeysStatusFuncName, request.FuncName)
				assert.Equal(t, validatorScAddress, request.ScAddress)
				assert.Equal(t, validatorScAddress, request.Caller)

				respSlices := resp[request.Args[0]]

				response := &vmQueryResponse{}
				response.Data.Data.ReturnData = respSlices

				respBytes, err := json.Marshal(response)
				assert.Nil(t, err)

				return respBytes, http.StatusOK, nil
			},
		}

		instance, _ := NewBLSKeysFetcher(wrapper, addresses, time.Hour)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		list, err := instance.GetAllBLSKeys(ctx, "test")
		assert.Equal(t, errContextClosing, err)
		assert.Empty(t, list)
	})
}
