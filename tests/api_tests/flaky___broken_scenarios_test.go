//nolint:gocritic
//nolint:gocyclo
package api_tests

import (
	"encoding/hex"
	"github.com/0chain/gosdk/core/encryption"
	"github.com/0chain/system_test/internal/api/util/endpoint"
	"testing"

	"github.com/0chain/system_test/internal/api/model"
	"github.com/0chain/system_test/internal/api/util/crypto"
	"github.com/stretchr/testify/require"
)

/*
Tests in here are skipped until the feature has been fixed
*/
func Test___BrokenScenariosRegisterWallet(t *testing.T) {
	t.Skip()
	t.Parallel()

	t.Run("Register wallet API call should be successful, ignoring invalid creation date", func(t *testing.T) {
		t.Parallel()

		mnemonic := crypto.GenerateMnemonic(t)
		expectedKeyPair := crypto.GenerateKeys(t, mnemonic)
		publicKeyBytes, _ := hex.DecodeString(expectedKeyPair.PublicKey.SerializeToHexStr())
		expectedClientId := encryption.Hash(publicKeyBytes)
		invalidCreationDate := -1

		walletRequest := model.ClientPutWalletRequest{Id: expectedClientId, PublicKey: expectedKeyPair.PublicKey.SerializeToHexStr(), CreationDate: &invalidCreationDate}

		registeredWallet, httpResponse, err := v1ClientPut(t, walletRequest, endpoint.ConsensusByHttpStatus(endpoint.HttpOkStatus))

		require.Nil(t, err, "Unexpected error [%s] occurred registering wallet with http response [%s]", err, httpResponse)
		require.NotNil(t, registeredWallet, "Registered wallet was unexpectedly nil! with http response [%s]", httpResponse)
		require.Equal(t, endpoint.HttpOkStatus, httpResponse.Status())
		require.Equal(t, registeredWallet.Id, expectedClientId)
		require.Equal(t, registeredWallet.PublicKey, expectedKeyPair.PublicKey.SerializeToHexStr())
		require.Greater(t, *registeredWallet.CreationDate, 0, "Creation date is an invalid value!")
		require.NotNil(t, registeredWallet.Version)
	})

	t.Run("Register wallet API call should be unsuccessful given an invalid request - client id invalid", func(t *testing.T) {
		t.Parallel()

		mnemonic := crypto.GenerateMnemonic(t)
		expectedKeyPair := crypto.GenerateKeys(t, mnemonic)
		walletRequest := model.ClientPutWalletRequest{Id: "invalid", PublicKey: expectedKeyPair.PublicKey.SerializeToHexStr()}

		walletResponse, httpResponse, err := v1ClientPut(t, walletRequest, endpoint.ConsensusByHttpStatus("400 Bad Request"))

		require.Nil(t, walletResponse, "Expected returned wallet to be nil but was [%s] with http response [%s]", walletResponse, httpResponse)
		require.NotNil(t, err, "Expected error when registering wallet but was nil.")
		require.Equal(t, "400 Bad Request", httpResponse.Status())
	})

	t.Run("Register wallet API call should be unsuccessful given an invalid request - public key invalid", func(t *testing.T) {
		t.Parallel()

		mnemonic := crypto.GenerateMnemonic(t)
		expectedKeyPair := crypto.GenerateKeys(t, mnemonic)
		publicKeyBytes, _ := hex.DecodeString(expectedKeyPair.PublicKey.SerializeToHexStr())
		clientId := encryption.Hash(publicKeyBytes)
		walletRequest := model.ClientPutWalletRequest{Id: clientId, PublicKey: "invalid"}

		walletResponse, httpResponse, err := v1ClientPut(t, walletRequest, endpoint.ConsensusByHttpStatus("400 Bad Request"))

		require.Nil(t, walletResponse, "Expected returned wallet to be nil but was [%s] with http response [%s]", walletResponse, httpResponse)
		require.NotNil(t, err, "Expected error when registering wallet but was nil.")
		require.Equal(t, "400 Bad Request", httpResponse.Status())
	})

	t.Run("Register wallet API call should be unsuccessful given an invalid request - empty json body", func(t *testing.T) {
		t.Parallel()

		walletRequest := model.ClientPutWalletRequest{}
		walletResponse, httpResponse, err := v1ClientPut(t, walletRequest, endpoint.ConsensusByHttpStatus("400 Bad Request"))

		require.Nil(t, walletResponse, "Expected returned wallet to be nil but was [%s] with http response [%s]", walletResponse, httpResponse)
		require.NotNil(t, err, "Expected error when registering wallet but was nil.")
		require.Equal(t, "400 Bad Request", httpResponse.Status())
	})
}
