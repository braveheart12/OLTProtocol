package data

import (
	"encoding/hex"
	"math/big"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"golang.org/x/crypto/ripemd160"
)

// CurrencyId is a readable wrapper on ids in int
type CurrencyId int

const (
	OLT CurrencyId = iota
	BTC
	ETH
	VT
)

func init() {

	serial.Register(Currency{})
	serial.Register(CurrencyId(0))

	// The currency configurations are initialized in init
	// so that they can be loaded from a config file in future.

	Currencies = map[string]Currency{
		"OLT": Currency{"OLT", ONELEDGER, 0},
		"BTC": Currency{"BTC", BITCOIN, 1},
		"ETH": Currency{"ETH", ETHEREUM, 2},
		"VT":  Currency{"VT", ONELEDGER, 3},
	}

	CurrencyIdMap = map[int]string{
		0: "OLT",
		1: "BTC",
		2: "ETH",
		3:  "VT",
	}

	CurrenciesExtra = map[string]Extra{
		"OLT": Extra{big.NewFloat(1000000000000000000), 6, 'f'},
		"BTC": Extra{big.NewFloat(1), 0, 'f'}, // TODO: This needs to be set correctly
		"ETH": Extra{big.NewFloat(1), 0, 'f'}, // TODO: This needs to be set correctly
		"VT":  Extra{big.NewFloat(1), 0, 'f'},
	}

}

//Currency datatype holds the type of curency,
// the chain it uses and its id
type Currency struct {
	Name  string     `json:"name"`
	Chain ChainType  `json:"chain"`
	Id    int `json:"id"`
}


// TODO: Separated from Currency to avoid serializing big floats and giving out this info
type Extra struct {
	Units   *big.Float
	Decimal int
	Format  uint8
}


// TODO: These need to be driven from a domain database, also they are many-to-one with chains
var Currencies map[string]Currency

var CurrencyIdMap map[int]string

var CurrenciesExtra map[string]Extra

//
//
//
// Key sets a encodable key for the currency entry,
// we may end up using currencyCodes instead.
func (c Currency) Key() string {
	hasher := ripemd160.New()

	// serialize to a
	b, err := serial.Serialize(c, serial.JSON)
	if err != nil {
		log.Fatal("hash serialize failed", "err", err)
	}

	// hash the binary data
	_, err = hasher.Write(b)
	if err != nil {
		log.Fatal("hasher failed", "err", err)
	}
	b = hasher.Sum(nil)

	// encode to hex representation
	return hex.EncodeToString(b)
}

// Look up the currency by its name
func GetCurrency(currency string) Currency {
	return Currencies[currency]
}

// Get base of a currency
func GetBase(currency string) *big.Float {
	return GetExtra(currency).Units
}

// Get a CurrencyExtra by its name
func GetExtra(currency string) Extra {
	if value, ok := CurrenciesExtra[currency]; ok {
		return value
	}
	return CurrenciesExtra["OLT"]
}
