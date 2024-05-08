package parsers

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

const commentMarker = "#"

var bech32PubKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(core.AddressLen, core.AddressHRP)

type listParser struct {
}

// NewListParser creates a new file list parser
func NewListParser() *listParser {
	return &listParser{}
}

// ParseFile will try to parse to file and split the identities in 2. Errors if something is wrong with the file
func (parser *listParser) ParseFile(filename string) (*core.IdentitiesHolder, error) {
	dataFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	result := &core.IdentitiesHolder{}

	spitLines := strings.Split(string(dataFile), "\n")
	for index, line := range spitLines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.Index(line, commentMarker) == 0 {
			continue
		}
		err = parser.processLine(line, result)
		if err != nil {
			return nil, fmt.Errorf("%w on line %d", err, index)
		}
	}

	return result, nil
}

func (parser *listParser) processLine(line string, identitiesHolder *core.IdentitiesHolder) error {
	if strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\"") {
		line = line[1 : len(line)-1]
	}

	if len(line) == core.BLSHexKeyLen {
		// try to unhex it
		_, err := hex.DecodeString(line)
		if err != nil {
			return err
		}

		identitiesHolder.BlsHexKeys = append(identitiesHolder.BlsHexKeys, line)

		return nil
	}
	decoded, err := bech32PubKeyConverter.Decode(line)
	if err != nil {
		return err
	}

	address := core.Address{
		Hex:    hex.EncodeToString(decoded),
		Bech32: line,
	}
	identitiesHolder.Addresses = append(identitiesHolder.Addresses, address)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (parser *listParser) IsInterfaceNil() bool {
	return parser == nil
}
