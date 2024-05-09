package core

import "time"

// BLSKeyLen is the constant for the length of the BLS key
const BLSKeyLen = 96

// BLSHexKeyLen is the constant for the length of the hexed BLS key
const BLSHexKeyLen = 2 * BLSKeyLen

// AddressLen is the constant for the length of the address
const AddressLen = 32

// AddressHRP is the bech32 HRP used in addresses
const AddressHRP = "erd"

// EveryWeekDay is the constant that encodes each week day option
const EveryWeekDay = time.Weekday(-1)
