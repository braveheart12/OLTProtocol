package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/status"
)

func CreatePaymentRequest(app Application, identities []id.Identity, quotient data.Coin) []byte {
	var signers []id.PublicKey
	chainId := app.Admin.Get(chainKey)
	inputs := make([]action.SendInput, 0)
	outputs := make([]action.SendOutput, 0)

	for _, identity := range identities {
		if identity.Name == "" {
			log.Error("Missing Party argument")
			return nil
		}

		party, err := app.Identities.FindName(identity.Name)
		if err == status.MISSING_DATA {
			log.Debug("CreatePaymentRequest", "PartyMissingData", err)
			return nil
		}

		partyBalance := app.Utxo.Get(party.AccountKey)
		if partyBalance == nil {
			interimBalance := data.NewBalance(0, "OLT")
			partyBalance = &interimBalance
		}

		//fee := conv.GetCoin(args.Fee, args.Currency)
		//gas := conv.GetCoin(args.Gas, args.Currency)

		inputs = append(inputs,
			action.NewSendInput(party.AccountKey, partyBalance.Amount))

		outputs = append(outputs,
			action.NewSendOutput(party.AccountKey, partyBalance.Amount.Plus(quotient)))
	}

	payment, err := app.Accounts.FindName("Payment-OneLedger")
	if err != status.SUCCESS {
		log.Fatal("Payment Account not found")
	}
	paymentBalance := app.Utxo.Get(payment.AccountKey())
	log.Debug("CreatePaymentRequest", "paymentBalance", paymentBalance)

	numberValidators := data.NewCoin(int64(len(identities)), "OLT")
	totalPayment := quotient.Multiply(numberValidators)

	inputs = append(inputs,
		action.NewSendInput(payment.AccountKey(), paymentBalance.Amount))

	outputs = append(outputs,
		action.NewSendOutput(payment.AccountKey(), paymentBalance.Amount.Minus(totalPayment)))

	// Create base transaction
	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  string(chainId.([]byte)),
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     data.NewCoin(0, "OLT"),
		Gas:     data.NewCoin(0, "OLT"),
	}

	signed := action.SignTransaction(send)
	packet := action.PackRequest(action.SEND, signed)

	return packet
}