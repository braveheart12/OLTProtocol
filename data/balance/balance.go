/*

   ____             _          _
  / __ \           | |        | |
 | |  | |_ __   ___| | ___  __| | __ _  ___ _ __
 | |  | | '_ \ / _ \ |/ _ \/ _` |/ _` |/ _ \ '__|
 | |__| | | | |  __/ |  __/ (_| | (_| |  __/ |
  \____/|_| |_|\___|_|\___|\__,_|\__, |\___|_|
                                  __/ |
                                 |___/

	Copyright 2017 - 2019 OneLedger

*/

package balance

// Wrap the amount with owner information
type Balance struct {
	Amounts map[string]Coin `json:"amounts"`
}

/*
	Generators
*/
func NewBalance() *Balance {
	amounts := make(map[string]Coin, 0)
	result := &Balance{
		Amounts: amounts,
	}
	return result
}

/*
	methods
*/
func (b *Balance) FindCoin(currency Currency) *Coin {
	if coin, ok := b.Amounts[currency.StringKey()]; ok {
		return &coin
	}
	return nil
}

// Add a new or existing coin
func (b *Balance) AddCoin(coin Coin) *Balance {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		b.Amounts[coin.Currency.StringKey()] = coin
		return b
	}

	amt, err := result.Plus(coin)
	if err != nil {
		return b
	}
	b.Amounts[coin.Currency.StringKey()] = amt
	return b
}

func (b *Balance) MinusCoin(coin Coin) (*Balance, error) {
	result := b.FindCoin(coin.Currency)
	if result == nil {
		return b, ErrInsufficientBalance
	}

	var err error

	b.Amounts[coin.Currency.StringKey()], err = result.Minus(coin)
	if err != nil {
		return b, err
	}

	return b, nil
}

func (b *Balance) GetCoin(currency Currency) Coin {
	result := b.FindCoin(currency)
	if result == nil {
		// NOTE: Missing coins are actually zero value coins.
		return currency.NewCoinFromInt(0)
	}
	return b.Amounts[currency.StringKey()]
}

func (b *Balance) setAmount(coin Coin) *Balance {
	b.Amounts[coin.Currency.StringKey()] = coin
	return b
}

func (b Balance) IsEnoughBalance(balance Balance) bool {
	for i, coin := range balance.Amounts {
		v, ok := b.Amounts[i]
		if !ok {
			v = coin.Currency.NewCoinFromInt(0)
		}

		_, err := v.Minus(coin)
		if err != nil {
			return false
		}
	}
	return true
}

// String method used in fmt and Dump
func (b Balance) String() string {
	buffer := ""
	for _, coin := range b.Amounts {
		if buffer != "" {
			buffer += ", "
		}
		buffer += coin.String()
	}
	return buffer
}
