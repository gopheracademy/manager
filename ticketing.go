package main

//this is an arbitrary pick

// ClaimPayment represents a payment for N claims
type ClaimPayment struct {
	ID string //uuid
	// ClaimsPayed would be what in a bill one see as detail.
	ClaimsPayed []*SlotClaim
	Payment     []FinancialInstrument
	Invoice     string // let us fill this once we know how to invoice
}

// TotalDue returns the total cost to cover by this payment.
func (c *ClaimPayment) TotalDue() int {
	totalDue := 0
	for _, sc := range c.ClaimsPayed {
		totalDue = totalDue + sc.EventSlot.Cost
	}
	return totalDue
}

// Fulfilled returns true if the payment of this invoice has been fulfilled
func (c *ClaimPayment) Fulfilled() bool {
	totalDue := c.TotalDue()
	f, _ := PaymentBalanced(totalDue, c.Payment...)
	b, _ := DebtBalanced(c.Payment...)
	return f && b
}

// SlotClaim represents one occupancy of one slot.
type SlotClaim struct {
	ID        string // uuid
	EventSlot *EventSlot
	// TicketID should only be valid when combined with the correct Attendee ID/Email
	TicketID string // uuid
	// Redeemed represents whether this has been used (ie the Attendee enrolled in front desk
	// or into the online conf system) until this is not true, transfer/refund might be possible.
	Redeemed bool
}

// Attendee is a person attending one or more Slots of the Conference.
type Attendee struct {
	ID    uint
	Email string
	// CoCAccepted, claims cannot be used without this.
	CoCAccepted bool
	Claims      []SlotClaim
}

// Finance Section

// PaymentMethodMoney represents a payment in cash.
type PaymentMethodMoney struct {
	PaymentRef string // stripe payment ID/Log?
	Amount     int
}

// Total implements FinancialInstrument
func (p *PaymentMethodMoney) Total() int {
	return p.Amount
}

// Type implements FinancialInstrument
func (p *PaymentMethodMoney) Type() AssetType {
	return ATCash
}

var _ FinancialInstrument = &PaymentMethodMoney{}

// PaymentMethodConferenceDiscount represents a discount issued by the event.
type PaymentMethodConferenceDiscount struct {
	// Detail describes what kind of discount was issued (ie 100% sponsor, 30% grant)
	Detail string
	Amount int
}

// Total implements FinancialInstrument
func (p *PaymentMethodConferenceDiscount) Total() int {
	return p.Amount
}

// Type implements FinancialInstrument
func (p *PaymentMethodConferenceDiscount) Type() AssetType {
	return ATDiscount
}

var _ FinancialInstrument = &PaymentMethodConferenceDiscount{}

// PaymentMethodCreditNote represents credit extended to defer payment.
type PaymentMethodCreditNote struct {
	Detail string
	Amount int
}

// Total implements FinancialInstrument
func (p *PaymentMethodCreditNote) Total() int {
	return p.Amount
}

// Type implements FinancialInstrument
func (p *PaymentMethodCreditNote) Type() AssetType {
	return ATReceivable
}

var _ FinancialInstrument = &PaymentMethodCreditNote{}

// AssetType is a type of accounting asset.
type AssetType string

const (
	// ATCash in this context means it is money, like a stripe payment
	ATCash AssetType = "cash"
	// ATReceivable in this context means it is a promise of payment
	ATReceivable AssetType = "receivable"
	// ATDiscount in this context means an issued discount (represented as a fixed amount for
	// accounting's sake)
	ATDiscount AssetType = "discount"
)

// FinancialInstrument represents any kind of instrument used to cover a debt.
// oto: "skip"
type FinancialInstrument interface {
	// Total is the total amount fulfilled by this instrument
	Total() int
	// Type is the type of asset represented
	Type() AssetType
}

// PaymentBalanced returns true or false depending on balancing status and missing
// payment amount if any.
func PaymentBalanced(amount int, payments ...FinancialInstrument) (bool, int) {
	receivables := 0
	received := 0
	for _, p := range payments {
		switch p.Type() {
		case ATCash, ATDiscount:
			received = received + p.Total()
		case ATReceivable:
			receivables = receivables + p.Total()
		}
	}
	missing := amount - received - receivables
	return missing <= 0, missing
}

// DebtBalanced returns true if all cretid notes or similar instruments have been covered or an
// amount if not.
func DebtBalanced(payments ...FinancialInstrument) (bool, int) {
	receivables := 0
	received := 0
	for _, p := range payments {
		switch p.Type() {
		case ATCash, ATDiscount:
			received = received + p.Total()
		case ATReceivable:
			receivables = receivables + p.Total()

		}
	}
	missing := receivables - received
	return missing <= 0, missing
}
