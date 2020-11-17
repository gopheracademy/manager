package business

import (
	"fmt"

	"github.com/gopheracademy/manager/build"
	uuid "github.com/satori/go.uuid"
)

type PurchaseStore interface {
	CreateSlotClaim(*build.SlotClaim) (*build.SlotClaim, error)

	UpdateAttendee(*build.Attendee) (*build.Attendee, error)

	CreateClaimPayment(*build.ClaimPayment) (*build.ClaimPayment, error)

	UpdateClaimPayment(*build.ClaimPayment) (*build.ClaimPayment, error)
	ChangeSlotClaimOwner([]build.SlotClaim, *build.Attendee, *build.Attendee) (*build.Attendee, *build.Attendee, error)
}

// ClaimSlots claims N slots for an attendee.
func ClaimSlots(storer PurchaseStore,
	attendee *build.Attendee, slots ...build.EventSlot) ([]build.SlotClaim, error) {
	var err error
	var claims = make([]build.SlotClaim, 0, len(slots))
	for i := range slots {
		slot := slots[i]
		sc := &build.SlotClaim{
			EventSlot: &slot,
			TicketID:  uuid.NewV4().String(),
		}
		sc, err = storer.CreateSlotClaim(sc)
		if err != nil {
			return nil, fmt.Errorf("Claiming a slot: %w", err)
		}
		claims[i] = *sc
	}
	attendee.Claims = append(attendee.Claims, claims...)
	_, err = storer.UpdateAttendee(attendee)
	if err != nil {
		return nil, fmt.Errorf("Updating claimed slots for attendee: %w", err)
	}
	return claims, nil
}

// PayClaims assigns payments and/or credits to a set of claims.
func PayClaims(store PurchaseStore,
	attendee *build.Attendee, claims []build.SlotClaim,
	payments []build.FinancialInstrument) (*build.ClaimPayment, error) {
	ptrClaims := make([]*build.SlotClaim, len(claims))
	for i := range claims {
		ptrClaims[i] = &claims[i]
	}
	claimPayment := &build.ClaimPayment{
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
	currencyType build.AssetType
}

func (e *ErrInvalidCurrency) Error() string {
	return fmt.Sprintf("the debt cannot be covered with %s", e.currencyType)
}

// CoverCredit adds funds to a payment to cover for receivables.
func CoverCredit(store PurchaseStore,
	existingPayment *build.ClaimPayment,
	payments []build.FinancialInstrument) error {
	for i := range payments {
		if payments[i].Type() == build.ATReceivable {
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
	source, target *build.Attendee, claims []build.SlotClaim) (*build.Attendee, *build.Attendee, error) {
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
