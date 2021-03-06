package stripe

import "encoding/json"

// BalanceTransactionStatus is the list of allowed values for the balance transaction's status.
type BalanceTransactionStatus string

const (
	BalanceTransactionStatusAvailable BalanceTransactionStatus = "available"
	BalanceTransactionStatusPending   BalanceTransactionStatus = "pending"
)

// BalanceTransactionType is the list of allowed values for the balance transaction's type.
type BalanceTransactionType string

const (
	BalanceTransactionTypeAdjustment               BalanceTransactionType = "adjustment"
	BalanceTransactionTypeApplicationFee           BalanceTransactionType = "application_fee"
	BalanceTransactionTypeApplicationFeeRefund     BalanceTransactionType = "application_fee_refund"
	BalanceTransactionTypeCharge                   BalanceTransactionType = "charge"
	BalanceTransactionTypePayment                  BalanceTransactionType = "payment"
	BalanceTransactionTypePaymentFailureRefund     BalanceTransactionType = "payment_failure_refund"
	BalanceTransactionTypePaymentRefund            BalanceTransactionType = "payment_refund"
	BalanceTransactionTypePayout                   BalanceTransactionType = "payout"
	BalanceTransactionTypePayoutCancel             BalanceTransactionType = "payout_cancel"
	BalanceTransactionTypePayoutFailure            BalanceTransactionType = "payout_failure"
	BalanceTransactionTypeRecipientTransfer        BalanceTransactionType = "recipient_transfer"
	BalanceTransactionTypeRecipientTransferCancel  BalanceTransactionType = "recipient_transfer_cancel"
	BalanceTransactionTypeRecipientTransferFailure BalanceTransactionType = "recipient_transfer_failure"
	BalanceTransactionTypeRefund                   BalanceTransactionType = "refund"
	BalanceTransactionTypeStripeFee                BalanceTransactionType = "stripe_fee"
	BalanceTransactionTypeTransfer                 BalanceTransactionType = "transfer"
	BalanceTransactionTypeTransferRefund           BalanceTransactionType = "transfer_refund"
)

// BalanceTransactionSourceType consts represent valid balance transaction sources.
type BalanceTransactionSourceType string

const (
	BalanceTransactionSourceTypeApplicationFee    BalanceTransactionSourceType = "application_fee"
	BalanceTransactionSourceTypeCharge            BalanceTransactionSourceType = "charge"
	BalanceTransactionSourceTypeDispute           BalanceTransactionSourceType = "dispute"
	BalanceTransactionSourceTypePayout            BalanceTransactionSourceType = "payout"
	BalanceTransactionSourceTypeRecipientTransfer BalanceTransactionSourceType = "recipient_transfer"
	BalanceTransactionSourceTypeRefund            BalanceTransactionSourceType = "refund"
	BalanceTransactionSourceTypeReversal          BalanceTransactionSourceType = "reversal"
	BalanceTransactionSourceTypeTransfer          BalanceTransactionSourceType = "transfer"
)

// BalanceTransactionSource describes the source of a balance Transaction.
// The Type should indicate which object is fleshed out.
// For more details see https://stripe.com/docs/api#retrieve_balance_transaction
type BalanceTransactionSource struct {
	ApplicationFee    *ApplicationFee              `json:"-"`
	Charge            *Charge                      `json:"-"`
	Dispute           *Dispute                     `json:"-"`
	ID                string                       `json:"id"`
	Payout            *Payout                      `json:"-"`
	RecipientTransfer *RecipientTransfer           `json:"-"`
	Refund            *Refund                      `json:"-"`
	Reversal          *Reversal                    `json:"-"`
	Transfer          *Transfer                    `json:"-"`
	Type              BalanceTransactionSourceType `json:"object"`
}

// BalanceParams is the set of parameters that can be used when retrieving a balance.
// For more details see https://stripe.com/docs/api#balance.
type BalanceParams struct {
	Params `form:"*"`
}

// BalanceTransactionParams is the set of parameters that can be used when retrieving a transaction.
// For more details see https://stripe.com/docs/api#retrieve_balance_transaction.
type BalanceTransactionParams struct {
	Params `form:"*"`
}

// BalanceTransactionListParams is the set of parameters that can be used when listing balance transactions.
// For more details see https://stripe.com/docs/api/#balance_history.
type BalanceTransactionListParams struct {
	ListParams       `form:"*"`
	AvailableOn      *int64            `form:"available_on"`
	AvailableOnRange *RangeQueryParams `form:"available_on"`
	Created          *int64            `form:"created"`
	CreatedRange     *RangeQueryParams `form:"created"`
	Currency         *string           `form:"currency"`
	Payout           *string           `form:"payout"`
	Source           *string           `form:"source"`
	Type             *string           `form:"type"`
}

// Balance is the resource representing your Stripe balance.
// For more details see https://stripe.com/docs/api/#balance.
type Balance struct {
	Available []*Amount `json:"available"`
	Livemode  bool      `json:"livemode"`
	Pending   []*Amount `json:"pending"`
}

// BalanceTransaction is the resource representing the balance transaction.
// For more details see https://stripe.com/docs/api/#balance.
type BalanceTransaction struct {
	Amount      int64                    `json:"amount"`
	AvailableOn int64                    `json:"available_on"`
	Created     int64                    `json:"created"`
	Currency    Currency                 `json:"currency"`
	Description string                   `json:"description"`
	ID          string                   `json:"id"`
	Fee         int64                    `json:"fee"`
	FeeDetails  []*BalanceTransactionFee `json:"fee_details"`
	Net         int64                    `json:"net"`
	Recipient   string                   `json:"recipient"`
	Source      string                   `json:"source"`
	Status      BalanceTransactionStatus `json:"status"`
	Type        BalanceTransactionType   `json:"type"`
}

// BalanceTransactionList is a list of transactions as returned from a list endpoint.
type BalanceTransactionList struct {
	ListMeta
	Data []*BalanceTransaction `json:"data"`
}

// Amount is a structure wrapping an amount value and its currency.
type Amount struct {
	Value    int64    `json:"amount"`
	Currency Currency `json:"currency"`
}

// BalanceTransactionFee is a structure that breaks down the fees in a transaction.
type BalanceTransactionFee struct {
	Amount      int64    `json:"amount"`
	Application string   `json:"application"`
	Currency    Currency `json:"currency"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
}

// UnmarshalJSON handles deserialization of a Transaction.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (t *BalanceTransaction) UnmarshalJSON(data []byte) error {
	type bt BalanceTransaction
	var tt bt
	err := json.Unmarshal(data, &tt)
	if err == nil {
		*t = BalanceTransaction(tt)
	} else {
		// the id is surrounded by "\" characters, so strip them
		t.ID = string(data[1 : len(data)-1])
	}

	return nil
}

// UnmarshalJSON handles deserialization of a BalanceTransactionSource.
// This custom unmarshaling is needed because the specific
// type of transaction source it refers to is specified in the JSON
func (s *BalanceTransactionSource) UnmarshalJSON(data []byte) error {
	type source BalanceTransactionSource
	var ss source
	err := json.Unmarshal(data, &ss)
	if err == nil {
		*s = BalanceTransactionSource(ss)

		switch s.Type {
		case BalanceTransactionSourceTypeApplicationFee:
			err = json.Unmarshal(data, &s.ApplicationFee)
		case BalanceTransactionSourceTypeCharge:
			err = json.Unmarshal(data, &s.Charge)
		case BalanceTransactionSourceTypeDispute:
			err = json.Unmarshal(data, &s.Dispute)
		case BalanceTransactionSourceTypePayout:
			err = json.Unmarshal(data, &s.Payout)
		case BalanceTransactionSourceTypeRecipientTransfer:
			err = json.Unmarshal(data, &s.RecipientTransfer)
		case BalanceTransactionSourceTypeRefund:
			err = json.Unmarshal(data, &s.Refund)
		case BalanceTransactionSourceTypeReversal:
			err = json.Unmarshal(data, &s.Reversal)
		case BalanceTransactionSourceTypeTransfer:
			err = json.Unmarshal(data, &s.Transfer)
		}

		if err != nil {
			return err
		}
	} else {
		// the id is surrounded by "\" characters, so strip them
		s.ID = string(data[1 : len(data)-1])
	}

	return nil
}

// MarshalJSON handles serialization of a BalanceTransactionSource.
func (s *BalanceTransactionSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ID)
}
