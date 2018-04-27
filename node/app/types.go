/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package app

// Coin is the basic amount, specified in integers, at the smallest increment (i.e. a satoshi, not a bitcoin)
type Coin struct {
	Currency string `json:"denom"`
	Amount   int64  `json:"amount"`
}

type Coins []Coin

// inputs into a send transaction (similar to Bitcoin)
type SendInput struct {
	Address   Address   `json:"address"`
	Coins     Coins     `json:"coins"`
	Sequence  int       `json:"sequence"`
	Signature Signature `json:"signature"`
	PubKey    PublicKey `json:"pub_key"`
}

// outputs for a send transaction (similar to Bitcoin)
type SendOutput struct {
	Address Address `json:"address"`
	Coins   Coins   `json:"coins"`
}