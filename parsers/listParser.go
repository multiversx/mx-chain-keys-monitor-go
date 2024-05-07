package parsers

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

const commentMarker = "#"

var bech32PubKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(core.AddressLen, core.AddressHRP)

type listParser struct {
	filename   string
	mut        sync.RWMutex
	blsHexKeys []string
	addresses  []core.Address
}

// NewListParser creates a new file list parser
func NewListParser(filename string) *listParser {
	return &listParser{
		filename: filename,
	}
}

// ParseFile will try to parse to file and split the identities in 2. Errors if something is wrong with the file
func (parser *listParser) ParseFile() error {
	dataFile, err := os.ReadFile(parser.filename)
	if err != nil {
		return err
	}

	parser.mut.Lock()
	defer parser.mut.Unlock()

	parser.emptyLists()

	spitLines := strings.Split(string(dataFile), "\n")
	for index, line := range spitLines {
		line = strings.Trim(line, " \t\r\n")
		if len(line) == 0 {
			continue
		}
		if strings.Index(line, commentMarker) == 0 {
			continue
		}
		err = parser.processLine(line)
		if err != nil {
			parser.emptyLists()
			return fmt.Errorf("%w on line %d", err, index)
		}
	}

	return nil
}

func (parser *listParser) emptyLists() {
	parser.blsHexKeys = make([]string, 0)
	parser.addresses = make([]core.Address, 0)
}

func (parser *listParser) processLine(line string) error {
	if strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\"") {
		line = line[1 : len(line)-1]
	}

	if len(line) == core.BLSHexKeyLen {
		// try to unhex it
		_, err := hex.DecodeString(line)
		if err != nil {
			return err
		}

		parser.blsHexKeys = append(parser.blsHexKeys, line)

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
	parser.addresses = append(parser.addresses, address)

	return nil
}

// BlsHexKeys returns the BLS keys addresses in hex format
func (parser *listParser) BlsHexKeys() []string {
	parser.mut.RLock()
	defer parser.mut.RUnlock()

	return parser.blsHexKeys
}

// Addresses returns identities addresses
func (parser *listParser) Addresses() []core.Address {
	parser.mut.RLock()
	defer parser.mut.RUnlock()

	return parser.addresses
}

// IsInterfaceNil returns true if there is no value under the interface
func (parser *listParser) IsInterfaceNil() bool {
	return parser == nil
}
