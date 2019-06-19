package ons

import (
	"encoding/hex"
	"fmt"

	"github.com/Oneledger/protocol/data/ons"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
)

var _ action.Msg = DomainCreate{}

type DomainCreate struct {
	Owner   action.Address
	Account action.Address
	Name    string
	Price   action.Amount
}

func (dc DomainCreate) Signers() []action.Address {
	return []action.Address{dc.Owner}
}

func (dc DomainCreate) Type() action.Type {
	return action.DOMAIN_CREATE
}

func (dc DomainCreate) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(dc.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: dc.Owner.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

var _ action.Tx = domainCreateTx{}

type domainCreateTx struct {
}

func (domainCreateTx) Validate(ctx *action.Context, msg action.Msg, fee action.Fee, memo string, signatures []action.Signature) (bool, error) {
	ok, err := action.ValidateBasic(msg, fee, memo, signatures)
	if err != nil {
		return ok, err
	}

	create, ok := msg.(*DomainCreate)
	if !ok {
		return false, action.ErrWrongTxType
	}

	if create.Owner == nil || len(create.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !create.Price.IsValid(ctx) || create.Price.Currency != "OLT" {
		return false, action.ErrInvalidAmount
	}

	coin := create.Price.ToCoin(ctx)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromInt(CREATE_PRICE)) {
		return false, action.ErrNotEnoughFund
	}

	return true, nil
}

func (domainCreateTx) ProcessCheck(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	create, ok := msg.(*DomainCreate)
	if !ok {
		return false, action.Response{Log: "DomainCreate cast failed"}
	}

	b, err := ctx.Balances.Get(create.Owner.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("failed to get balance for owner: %s", hex.EncodeToString(create.Owner))}
	}
	price := create.Price.ToCoin(ctx)

	b, err = b.MinusCoin(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if ctx.Domains.Exists(create.Name) {
		return false, action.Response{Log: errors.New("Domain already exist").Error()}
	}
	result := action.Response{
		Tags: create.Tags(),
	}

	return true, result
}

func (domainCreateTx) ProcessDeliver(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	create, ok := msg.(*DomainCreate)
	if !ok {
		return false, action.Response{Log: "DomainCreate cast failed"}
	}

	b, err := ctx.Balances.Get(create.Owner.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("failed to get balance for owner: %s", hex.EncodeToString(create.Owner))}
	}
	price := create.Price.ToCoin(ctx)

	b, err = b.MinusCoin(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.Balances.Set(create.Owner.Bytes(), *b)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "set balance of owner").Error()}
	}

	if ctx.Domains.Exists(create.Name) {
		return false, action.Response{Log: errors.New("domain already exist: " + create.Name).Error()}
	}
	domain := ons.NewDomain(
		create.Owner,
		create.Account,
		create.Name,
		ctx.Header.Height,
	)
	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	result := action.Response{
		Tags: create.Tags(),
	}
	return true, result
}

func (domainCreateTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
