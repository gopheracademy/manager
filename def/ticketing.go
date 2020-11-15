package def

import (
	"time"

	"github.com/shopspring/decimal" //this is an arbitrary pick
)

// EventSlot holds information for any sellable/giftable slot we have in the event for
// a Talk or any other activity that requires admission.
type EventSlot struct {
	ID          uint
	Name        string
	Description string
	Cost        decimal.Decimal
	Capacity    int // int should be enough even if we organize glastonbury
	StartDate   time.Time
	EndDate     time.Time
	// DependsOn means that these two Slots need to be acquired together, user must either buy
	// both Slots or pre-own one of the one it depends on.
	DependsOn *EventSlot
	// PurchaseableFrom indicates when this item is on sale, for instance early bird tickets are the first
	// ones to go on sale.
	PurchaseableFrom time.Time
	// PuchaseableUntil indicates when this item stops being on sale, for instance early bird tickets can
	// no loger be purchased N months before event.
	PurchaseableUntil time.Time
	// AvailableToPublic indicates is this is something that will appear on the tickets purchase page (ie, we can
	// issue sponsor tickets and those cannot be bought individually)
	AvailableToPublic bool
}

// ClaimPayment represents a payment for N claims
type ClaimPayment struct {
	ID string //uuid
	// ClaimsPayed would be what in a bill one see as detail.
	ClaimsPayed []*SlotClaim
	Payment     []FinancialInstrument
	Invoice     string // let us fill this once we know how to invoice
}

// TotalDue returns the total cost to cover by this payment.
func (c *ClaimPayment) TotalDue() decimal.Decimal {
	totalDue := decimal.Zero
	for _, sc := range c.ClaimsPayed {
		totalDue = totalDue.Add(sc.EventSlot.Cost)
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
	Amount     decimal.Decimal
}

// Total implements FinancialInstrument
func (p *PaymentMethodMoney) Total() decimal.Decimal {
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
	Amount decimal.Decimal
}

// Total implements FinancialInstrument
func (p *PaymentMethodConferenceDiscount) Total() decimal.Decimal {
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
	Amount decimal.Decimal
}

// Total implements FinancialInstrument
func (p *PaymentMethodCreditNote) Total() decimal.Decimal {
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
type FinancialInstrument interface {
	// Total is the total amount fulfilled by this instrument
	Total() decimal.Decimal
	// Type is the type of asset represented
	Type() AssetType
}

// PaymentBalanced returns true or false depending on balancing status and missing
// payment amount if any.
func PaymentBalanced(amount decimal.Decimal, payments ...FinancialInstrument) (bool, decimal.Decimal) {
	receivables := decimal.Zero
	received := decimal.Zero
	for _, p := range payments {
		switch p.Type() {
		case ATCash, ATDiscount:
			received.Add(p.Total())
		case ATReceivable:
			receivables.Add(p.Total())
		}
	}
	missing := amount.Sub(received).Sub(receivables)
	return missing.LessThanOrEqual(decimal.Zero), missing
}

// DebtBalanced returns true if all cretid notes or similar instruments have been covered or an
// amount if not.
func DebtBalanced(payments ...FinancialInstrument) (bool, decimal.Decimal) {
	receivables := decimal.Zero
	received := decimal.Zero
	for _, p := range payments {
		switch p.Type() {
		case ATCash, ATDiscount:
			received.Add(p.Total())
		case ATReceivable:
			receivables.Add(p.Total())
		}
	}
	missing := receivables.Sub(received)
	return missing.LessThanOrEqual(decimal.Zero), missing
}
