package data

import (
	"errors"
	"fmt"
	"github.com/Oneledger/protocol/node/serial"
	"math/big"

	"github.com/Oneledger/protocol/node/log"
)

var (
	ErrParsingBigInt       = errors.New("error parsing big int from fmt.Sscan")
	ErrWrongBalanceAdapter = errors.New("error in asserting to BalanceAdapter")
)

// BalanceAdapter is an easy to serailize representation of a Balance object. A full Balance object can be recostructed
// from a BalanceAdapter object and vice versa.
// There is a map flattening ofcourse for
type BalanceAdapter struct {
	Data []CoinData `json:"data"`
}

//

// CoinData is a flattening of coin map in a balance data type
type CoinData struct {
	CurId    int `json:"curr_id"`
	CurName  string     `json:"curr_name"`
	CurChain ChainType  `json:"curr_chain"`

	Amount string `json:"amt"`
}

//


func init() {
	serial.Register(BalanceAdapter{})
	serial.Register(CoinData{})
}

// NewBalanceAdapter creates a BalanceAdapter from a given Balance object,
// the coins are flattened to a list in the generator itself
// ideally there should be no change done to a data after this step. This datatype can go straight to serialization.
func NewBalanceAdapter(bal *Balance) *BalanceAdapter {
	//initialize with source pointer
	badap := &BalanceAdapter{}

	if bal == nil {
		return badap
	}

	// this allows to reserve capacity so the process of adding
	// items to the list
	badap.Data = make([]CoinData, 0, len(bal.Amounts))

	for id, coin := range bal.Amounts {
		cd := CoinData{
			CurId:    id,
			CurName:  coin.Currency.Name,
			CurChain: coin.Currency.Chain,
			Amount:   coin.Amount.String(),
		}

		badap.Data = append(badap.Data, cd)
	}

	return badap
}

//

// Extract recreates the Balance object form the Data BalanceAdapter holds after deserialization/
func (ba *BalanceAdapter) Extract() (*Balance, error) {

	b := &Balance{}
	b.Amounts = make(map[int]Coin)

	d := ba.Data
	for i := range d {

		curID := d[i].CurId

		//convert string representation to big int
		amt := new(big.Int)
		_, err := fmt.Sscan(d[i].Amount, amt)
		if err != nil {
			log.Error("error in parsing bigint for balance ", err)
			return nil, ErrParsingBigInt
		}

		coin := Coin{Amount: amt}
		coin.Currency.Id = curID
		coin.Currency.Name = d[i].CurName
		coin.Currency.Chain = d[i].CurChain

		b.Amounts[curID] = coin
	}

	return b, nil
}
