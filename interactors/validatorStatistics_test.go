package interactors

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidatorStatisticsInteractor(t *testing.T) {
	t.Parallel()

	t.Run("nil HTTP client wrapper should error", func(t *testing.T) {
		instance, err := NewValidatorStatisticsInteractor(nil)
		assert.Nil(t, instance)
		assert.Equal(t, errNilHTTPClientWrapper, err)
	})
	t.Run("should work", func(t *testing.T) {
		instance, err := NewValidatorStatisticsInteractor(&mock.HTTPClientWrapperStub{})
		assert.NotNil(t, instance)
		assert.Nil(t, err)
	})
}

func TestValidatorStatisticsInteractor_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *validatorStatisticsInteractor
	assert.True(t, instance.IsInterfaceNil())

	instance = &validatorStatisticsInteractor{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestValidatorStatisticsInteractor_Query(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("http client wrapper errors should return error", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			GetHTTPHandler: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				assert.Equal(t, "validator/statistics", endpoint)
				return nil, http.StatusBadRequest, expectedErr
			},
		}
		instance, _ := NewValidatorStatisticsInteractor(wrapper)
		response, err := instance.Query(context.Background())
		assert.Nil(t, response)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("http client wrapper returns not-ok status code should return error", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			GetHTTPHandler: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				assert.Equal(t, "validator/statistics", endpoint)
				return nil, http.StatusBadRequest, nil
			},
		}
		instance, _ := NewValidatorStatisticsInteractor(wrapper)
		response, err := instance.Query(context.Background())
		assert.Nil(t, response)
		assert.ErrorIs(t, err, errReturnCodeIsNotOk)
	})
	t.Run("http client wrapper returns a non-json response should return error", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			GetHTTPHandler: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				assert.Equal(t, "validator/statistics", endpoint)
				return []byte(`not-a-json-response`), http.StatusOK, nil
			},
		}
		instance, _ := NewValidatorStatisticsInteractor(wrapper)
		response, err := instance.Query(context.Background())
		assert.Nil(t, response)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})
	t.Run("http client wrapper return an invalid response should return error", func(t *testing.T) {
		t.Parallel()

		wrapper := &mock.HTTPClientWrapperStub{
			GetHTTPHandler: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				assert.Equal(t, "validator/statistics", endpoint)
				return []byte(`
{
  "data": {
  },
    "error": "",
    "code": "successful"
}`), http.StatusOK, nil
			},
		}
		instance, _ := NewValidatorStatisticsInteractor(wrapper)
		response, err := instance.Query(context.Background())
		assert.Nil(t, response)
		assert.Equal(t, errInvalidResponse, err)
	})
	t.Run("should work with real .json file", func(t *testing.T) {
		t.Parallel()

		data, err := os.ReadFile("testdata/response.json")
		require.Nil(t, err)

		wrapper := &mock.HTTPClientWrapperStub{
			GetHTTPHandler: func(ctx context.Context, endpoint string) ([]byte, int, error) {
				assert.Equal(t, "validator/statistics", endpoint)
				return data, http.StatusOK, nil
			},
		}

		expectedMap := map[string]*core.ValidatorStatistics{
			"0012e9c0f84ae4ce0345ab349b9532ac57d7bfea499e2bca4708ff7df70721ed7668bfcb35e0b5a1e9680cf11ed2ed15bef0cadbcb2e8b65ed1f1403e77fdc6a403a46bd6987f251faae51dc6687b3d507a8e5b4b86563be70d96dc0ccd8ce89": {
				TempRating: 90.69571,
				Rating:     90.69571,
			},
			"001e8e6180f3f52f6714b58fbdbd4055a62d2df46f55befa9815047491c53d25592b6efa37e36e6e13e9f9ea6559c80e35e6dc1950a7ff5e2b3dfcc9904df5eea6597bdbf7bb57fa1ce55ad6acb27f3a5a62599030c3c105d0d691c10caf9b0c": {
				TempRating: 100,
				Rating:     100,
			},
			"0026a4b6d8f4b6a2e22141341efb5dddf4db130a7e04d539dfd8c70bf3139d016ed958ccfd0bdcaf6aa866d11a09e21058f9dcb96fab9b863fe832cbed7f1705970ab9ad8c4de9da69e59a890751740064bfd84b7eb9e714e0e03fa8d776d004": {
				TempRating: 90.69571,
				Rating:     100,
			},
			"026c5e9d87c584b787050ffd7bc3484b1fb598b5b9bacab19c97a1da48498b5f5e3dd0a2befb3b89d82814af17bd2112813cb36a7d727f2bb90287075c31d6385beb6264ae761600cb2679368ef5d8cffe5c883a738335ca8771fb901b54488f": {
				TempRating: 99.69571,
				Rating:     100,
			},
			"02816b1072fa705ebcb640b156f7cca2e4b404760911aafeab829c83921156758ba9f4c18e39423ee536c5cdc33cb1182af82ddd4d6c1ef587c0db12b9d4c12df28f4e76c1b24eb2ad190f5217e5a05e099243f7672803975313d2a67c19098c": {
				TempRating: 49.69571,
				Rating:     50,
			},
			"0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80": {
				TempRating: 48.69571,
				Rating:     50,
			},
		}

		instance, _ := NewValidatorStatisticsInteractor(wrapper)
		response, err := instance.Query(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, expectedMap, response)
	})
}
