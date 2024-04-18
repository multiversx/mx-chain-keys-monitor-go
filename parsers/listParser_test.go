package parsers

import (
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"

	"github.com/stretchr/testify/assert"
)

func TestNewListParser(t *testing.T) {
	t.Parallel()

	parser := NewListParser("")
	assert.NotNil(t, parser)
}

func TestListParser_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *listParser
	assert.True(t, instance.IsInterfaceNil())

	instance = &listParser{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestListParser_ParseFile(t *testing.T) {
	t.Parallel()

	t.Run("file not found should error", func(t *testing.T) {
		t.Parallel()

		parser := NewListParser("file-not-found")
		err := parser.ParseFile()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "file-not-found")
		assert.Empty(t, parser.Addresses())
		assert.Empty(t, parser.BlsHexKeys())
	})
	t.Run("invalid BLS key size should error", func(t *testing.T) {
		t.Parallel()

		parser := NewListParser("./testdata/invalidBLSKeySize.list")
		err := parser.ParseFile()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid bech32 string length 191")
		assert.Empty(t, parser.Addresses())
		assert.Empty(t, parser.BlsHexKeys())
	})
	t.Run("not a hexed BLS key should error", func(t *testing.T) {
		t.Parallel()

		parser := NewListParser("./testdata/notAHexedBLSKey.list")
		err := parser.ParseFile()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid byte")
		assert.Empty(t, parser.Addresses())
		assert.Empty(t, parser.BlsHexKeys())
	})
	t.Run("not a valid bech32 address should error", func(t *testing.T) {
		t.Parallel()

		parser := NewListParser("./testdata/notAValidBech32Address.list")
		err := parser.ParseFile()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid checksum")
		assert.Empty(t, parser.Addresses())
		assert.Empty(t, parser.BlsHexKeys())
	})
	t.Run("empty file with comments should work", func(t *testing.T) {
		t.Parallel()

		parser := NewListParser("./testdata/emptyFileWithComments.list")
		err := parser.ParseFile()
		assert.Nil(t, err)
		assert.Empty(t, parser.Addresses())
		assert.Empty(t, parser.BlsHexKeys())
	})
	t.Run("correct file should work", func(t *testing.T) {
		t.Parallel()

		parser := NewListParser("./testdata/okMixedIdentities.list")
		err := parser.ParseFile()
		assert.Nil(t, err)

		expectedHexBLSKeys := []string{
			"015c24a0585c3007e02bb9168c7988cccd183285161b26a0fd908b68f4daf64518517b947f58a3c6cb3caebc4a1c84015470b2b43b05d6d9dbd463c817162b7f6c30f2bcb95fd7bc5dce98e5858200087c1d2b095f097dea57c142e4c0c0e088",
			"02dbca1ecef7a29da845c6ddd7b06254c4e6ef4506268e0117fd0350ab8a2f44b2997a02cf5eed3fd54673696d964301c90e5ff3bebc56d1b03138e77afc9d09bcb3d96b2efd93814c805a24761b2ba994be9d4696702966f6d53d149495378c",
		}
		assert.Equal(t, expectedHexBLSKeys, parser.BlsHexKeys())

		expectedAddresses := []core.Address{
			{
				Hex:    "c6762c7eb6edcb341d3e37f3e662363c98e6237b4245f567179661008d5160b0",
				Bech32: "erd1cemzcl4kah9ng8f7xle7vc3k8jvwvgmmgfzl2echjesspr23vzcqdexyy9",
			},
			{
				Hex:    "102a8ba34fce6f9be3b83d159eaae3a1cb8cabd9e31c6d92bb21d940251a9df8",
				Bech32: "erd1zq4ghg60eehehcac852ea2hr589ce27euvwxmy4my8v5qfg6nhuq99r9ez",
			},
		}
		assert.Equal(t, expectedAddresses, parser.Addresses())
	})
}
