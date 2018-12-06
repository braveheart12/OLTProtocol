package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/status"
)

func CreatePaymentRequest(app Application, quotient data.Coin, height int64) action.Transaction {
	chainId := app.Admin.Get(chainKey)
	identities := app.Validators.Approved
	payto := make([]action.SendTo, len(identities))

	for i, identity := range identities {
		if identity.Name == "" {
			log.Error("Missing Party argument")
			return nil
		}

		party, err := app.Identities.FindName(identity.Name)
		if err == status.MISSING_DATA {
			log.Debug("CreatePaymentRequest", "PartyMissingData", err)
			return nil
		}
		//todo : here we use index of the approved identity in the
		// Validators.approved to simplify the verification of validators in payment
		// need to makes this more secure.
		payto[i] = action.SendTo{
			AccountKey: party.AccountKey,
			Amount:     quotient,
		}
	}

	payment, err := app.Accounts.FindName(global.Current.PaymentAccount)
	if err != status.SUCCESS {
		log.Fatal("Payment Account not found")
	}

	// Create base transaction
	send := &action.Payment{
		Base: action.Base{
			Type:     action.PAYMENT,
			ChainId:  string(chainId.([]byte)),
			Owner:    payment.AccountKey(),
			Signers:  GetSigners(payment.AccountKey(), app),
			Sequence: height, //global.Current.Sequence,
		},
		PayTo: payto,
	}

	return send
}
