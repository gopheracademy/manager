package ticketing

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// PurchaseStore offers functionality for persistence of ticketing models.
type PurchaseStore interface {
	// AtomicOperation returns a store which will act as one single atomic operation.
	// It returns a commit and cancel functions and the Store .
	AtomicOperation() (func() error, func() error, PurchaseStore, error)
	// CreateSlotClaim saves a slot claim and returns it with the populated ID
	CreateSlotClaim(*SlotClaim) (*SlotClaim, error)

	// UpdateAttendee saves the passed attendee attributes on top of the existing one.
	UpdateAttendee(*Attendee) (*Attendee, error)

	CreateClaimPayment(*ClaimPayment) (*ClaimPayment, error)

	UpdateClaimPayment(*ClaimPayment) (*ClaimPayment, error)
	ChangeSlotClaimOwner([]SlotClaim, *Attendee, *Attendee) (*Attendee, *Attendee, error)
}

// ClaimSlots claims N slots for an attendee.
func ClaimSlots(storer PurchaseStore,
	attendee *Attendee, slots ...EventSlot) ([]SlotClaim, error) {
	succed, fail, atomic, err := storer.AtomicOperation()
	if err != nil {
		return nil, fmt.Errorf("beginning atomic operation: %w", err)
	}
	var claims = make([]SlotClaim, 0, len(slots))
	for i := range slots {
		slot := slots[i]
		sc := &SlotClaim{
			EventSlot: &slot,
			TicketID:  uuid.NewV4().String(),
		}
		sc, err = atomic.CreateSlotClaim(sc)
		if err != nil {
			if atomicErr := fail(); atomicErr != nil {
				err = fmt.Errorf("%w (also cancelling atomic operation: %v)", err, atomicErr)
			}
			return nil, fmt.Errorf("Claiming a slot: %w", err)
		}
		claims[i] = *sc
	}
	attendee.Claims = append(attendee.Claims, claims...)
	_, err = atomic.UpdateAttendee(attendee)
	if err != nil {
		if atomicErr := fail(); atomicErr != nil {
			err = fmt.Errorf("%w (also cancelling atomic operation: %v)", err, atomicErr)
		}
		return nil, fmt.Errorf("Updating claimed slots for attendee: %w", err)
	}
	if err := succed(); err != nil {
		return nil, fmt.Errorf("confirming atomic operation: %w", err)
	}
	return claims, nil
}

// PayClaims assigns payments and/or credits to a set of claims.
func PayClaims(store PurchaseStore,
	attendee *Attendee, claims []SlotClaim,
	payments []FinancialInstrument) (*ClaimPayment, error) {
	ptrClaims := make([]*SlotClaim, len(claims))
	for i := range claims {
		ptrClaims[i] = &claims[i]
	}
	claimPayment := &ClaimPayment{
		ClaimsPayed: ptrClaims,
		Payment:     payments,
	}

	claimPayment, err := store.CreateClaimPayment(claimPayment)
	if err != nil {
		return nil, fmt.Errorf("paying for claims: %w", err)
	}
	return claimPayment, nil
}

// ErrInvalidCurrency should be returned when paying with the wrong kind of instrument
// for instance covering credit with credit.
type ErrInvalidCurrency struct {
	currencyType AssetType
}

func (e *ErrInvalidCurrency) Error() string {
	return fmt.Sprintf("the debt cannot be covered with %s", e.currencyType)
}

// CoverCredit adds funds to a payment to cover for receivables.
func CoverCredit(store PurchaseStore,
	existingPayment *ClaimPayment,
	payments []FinancialInstrument) error {
	for i := range payments {
		if payments[i].Type() == ATReceivable {
			return &ErrInvalidCurrency{currencyType: payments[i].Type()}
		}
		existingPayment.Payment = append(existingPayment.Payment, payments[i])
	}
	_, err := store.UpdateClaimPayment(existingPayment)
	if err != nil {
		return fmt.Errorf("saving new payments %w", err)
	}
	return nil
}

// TransferClaims transfer claims from one user the the other, assuming they belong to the first.
func TransferClaims(storer PurchaseStore,
	source, target *Attendee, claims []SlotClaim) (*Attendee, *Attendee, error) {
	var err error
	sourceClaimsMap := map[uint64]bool{}
	for i := range source.Claims {
		sourceClaimsMap[source.Claims[i].ID] = true
	}
	for i := range claims {
		if belongsToSource := sourceClaimsMap[claims[i].ID]; !belongsToSource {
			return nil, nil, fmt.Errorf("%d claim for slot %s does not belong to %s", claims[i].ID, claims[i].EventSlot.Name, source.Email)
		}
	}
	if source, target, err = storer.ChangeSlotClaimOwner(claims, source, target); err != nil {
		return nil, nil, fmt.Errorf("reowning slot claim: %w", err)
	}

	return source, target, nil
}
