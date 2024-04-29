package interactors

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

const (
	getBlsKeysStatusFuncName = "getBlsKeysStatus"
	validatorScAddress       = "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l"
	endpoint                 = "vm-values/query"
	stakedStatus             = "staked"
)

type vmQueryRequest struct {
	ScAddress string   `json:"scAddress"`
	FuncName  string   `json:"funcName"`
	Caller    string   `json:"caller"`
	Args      []string `json:"args"`
}

type vmQueryResponse struct {
	Data struct {
		Data struct {
			ReturnData [][]byte `json:"returnData"`
			ReturnCode string   `json:"returnCode"`
		} `json:"data"`
	} `json:"data"`
	Code string `json:"code"`
}

type blsKeysFetcher struct {
	httpClientWrapper  HTTPClientWrapper
	addresses          []core.Address
	timeBetweenQueries time.Duration
}

// NewBLSKeysFetcher creates a new instance of type bls keys fetcher
func NewBLSKeysFetcher(
	httpClientWrapper HTTPClientWrapper,
	addresses []core.Address,
	timeBetweenQueries time.Duration,
) (*blsKeysFetcher, error) {
	if check.IfNil(httpClientWrapper) {
		return nil, errNilHTTPClientWrapper
	}

	return &blsKeysFetcher{
		httpClientWrapper:  httpClientWrapper,
		addresses:          addresses,
		timeBetweenQueries: timeBetweenQueries,
	}, nil
}

// GetAllBLSKeys will fetch all BLS keys for the set addresses
func (fetcher *blsKeysFetcher) GetAllBLSKeys(ctx context.Context, sender string) ([]string, error) {
	allBLSKeys := make([]string, 0)

	timer := time.NewTimer(fetcher.timeBetweenQueries)
	defer timer.Stop()

	for _, address := range fetcher.addresses {
		if len(address.Hex) == 0 {
			continue
		}

		blsKeys, err := fetcher.getBlsKeys(ctx, address, sender)
		if err != nil {
			return nil, err
		}

		timer.Reset(fetcher.timeBetweenQueries)
		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil, errContextClosing
		}

		allBLSKeys = append(allBLSKeys, blsKeys...)
	}

	return allBLSKeys, nil
}

func (fetcher *blsKeysFetcher) getBlsKeys(ctx context.Context, address core.Address, sender string) ([]string, error) {
	log.Debug("blsKeysFetcher.getBlsKeys", "address", address.Bech32)
	request := &vmQueryRequest{
		ScAddress: validatorScAddress,
		FuncName:  getBlsKeysStatusFuncName,
		Caller:    validatorScAddress,
		Args:      []string{address.Hex},
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	responseBytes, statusCode, err := fetcher.httpClientWrapper.PostHTTP(ctx, endpoint, requestBytes)
	if err != nil {
		return nil, err
	}
	if !core.IsHttpStatusCodeSuccess(statusCode) {
		return nil, fmt.Errorf("%w, but %d", errReturnCodeIsNotOk, statusCode)
	}
	response := &vmQueryResponse{}
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for i := 0; i < len(response.Data.Data.ReturnData)-1; i += 2 {
		blsKey := response.Data.Data.ReturnData[i]
		status := response.Data.Data.ReturnData[i+1]
		if len(blsKey) != core.BLSKeyLen {
			log.Warn("invalid response fetched", "returned data", fmt.Sprintf("%+v", response.Data.Data.ReturnData))
			continue
		}
		if string(status) != stakedStatus {
			continue
		}

		keys = append(keys, hex.EncodeToString(blsKey))
	}

	log.Debug("blsKeysFetcher.getBlsKeys", "sender", sender, "address", address.Bech32, "num keys", len(keys))

	return keys, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (fetcher *blsKeysFetcher) IsInterfaceNil() bool {
	return fetcher == nil
}
