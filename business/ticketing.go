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

func CoverCredit(store PurchaseStore) error {}

func ReleaseClaims() error
