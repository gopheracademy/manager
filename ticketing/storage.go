package ticketing

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ShiftLeftSecurity/gaum/db/chain"
	"github.com/ShiftLeftSecurity/gaum/db/connection"
	"github.com/ShiftLeftSecurity/gaum/db/logging"
	"github.com/ShiftLeftSecurity/gaum/db/postgres"
	"github.com/gopheracademy/manager/def"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	uuid "github.com/satori/go.uuid"
)

// NewSQLStorage returns a new Storage connected to the postgres-like db indicated by the connectionString.
func NewSQLStorage(connectionString string, logger *log.Logger) (*SQLStorage, error) {
	logLevel := connection.Error

	connector := postgres.Connector{
		ConnectionString: connectionString,
	}
	maxConnLifetime := 1 * time.Minute

	// you could open this without the config info but I put it here so other people looking at it
	// know where to tweak if necessary.
	db, err := connector.Open(&connection.Information{
		Logger:          logging.NewGoLogger(logger),
		LogLevel:        logLevel,
		ConnMaxLifetime: &maxConnLifetime,
		CustomDial: func(network, addr string) (net.Conn, error) {
			d := &net.Dialer{
				KeepAlive: time.Minute,
			}
			return d.Dial(network, addr)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("initializing db connection: %w", err)
	}
	return &SQLStorage{conn: db}, nil
}

// NewSQLStorageFromConnection returns a new SQLStorage using the passed connection
func NewSQLStorageFromConnection(conn connection.DB) *SQLStorage {
	return &SQLStorage{conn: conn}
}

// SQLStorage provides a Postgres Flavored storage backend to store ticketing information.
type SQLStorage struct {
	conn connection.DB
}

var _ PurchaseStore = &SQLStorage{}

// AtomicOperation begins a transaction and returns commit and rollback functions along with a new
// SQLStorage wrapping the tx
func (s *SQLStorage) AtomicOperation() (func() error, func() error, PurchaseStore, error) {
	tx, err := s.conn.BeginTransaction()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("beginning transation: %w", err)
	}
	return tx.CommitTransaction, tx.RollbackTransaction, &SQLStorage{conn: tx}, nil
}

const (
	ticketIDUniqueConstraint = "ticket_id_is_unique"
	tableSlotClaims          = "slot_claim"
)

// CreateAttendee creates a new attendee in the database and returns it.
func (s *SQLStorage) CreateAttendee(a *Attendee) (*Attendee, error) {
	results := []Attendee{}
	err := chain.New(s.conn).Insert(map[string]interface{}{
		"email":        a.Email,
		"coc_accepted": a.CoCAccepted,
	}).Table(tableAttendee).Returning("*").
		Fetch(&results)
	if err != nil {
		return nil, fmt.Errorf("creating new attendee: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("attendee was not created")
	}
	newClaims := make([]SlotClaim, len(a.Claims))
	for i := range a.Claims {
		c := a.Claims[i]
		returnedClaims := []SlotClaim{}
		err := chain.New(s.conn).Insert(map[string]interface{}{
			"attendee_id":   results[0].ID,
			"slot_claim_id": c.Redeemed,
		}).Table(tableAttendeeSlotClaims).
			OnConflict(func(c *chain.OnConflict) {
				// This claim was someone else's, this might be the result of transfering.
				c.OnConstraint(slotClaimIDUniqueConstraint).
					DoUpdate().
					Set("attendee_id", results[0].ID)
			}).Returning("*").Fetch(&returnedClaims)
		if err != nil {
			return nil, fmt.Errorf("inserting attendee claims: %w", err)
		}
		if len(returnedClaims) == 0 {
			return nil, fmt.Errorf("could not create claim")
		}
		newClaims[i] = returnedClaims[0]
	}
	newAttendee := results[0]
	newAttendee.Claims = newClaims
	return &newAttendee, nil
}

// ReadAttendeeByEmail returns an attendee for that email if one exists.
func (s *SQLStorage) ReadAttendeeByEmail(email string) (*Attendee, error) {
	if email == "" {
		return nil, fmt.Errorf("email is empty")
	}
	return s.readAttendee(email, 0)
}

// ReadAttendeeByID returns an attendee for the given ID if one exists.
func (s *SQLStorage) ReadAttendeeByID(id uint64) (*Attendee, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is not valid")
	}
	return s.readAttendee("", id)
}

func (s *SQLStorage) readAttendee(email string, id uint64) (*Attendee, error) {
	results := []Attendee{}
	q := chain.New(s.conn).Select("*").From(tableAttendee)
	if email != "" {
		q.AndWhere("email = ?", email)
	}
	if id != 0 {
		q.AndWhere("id = ?", id)
	}
	err := q.Fetch(&results)
	if err != nil {
		return nil, fmt.Errorf("reading attendee by email: %w", err)
	}
	if len(results) == 0 {
		return nil, nil
	}
	claims := []SlotClaim{}
	ats := chain.TablePrefix(attendeesToSlotTable)
	tsc := chain.TablePrefix(tableSlotClaims)
	err = chain.New(s.conn).Select("*").From(tableSlotClaims).
		Join(attendeesToSlotTable,
			chain.CompareExpressions(chain.Eq, tsc("id"), ats("slot_claim_id"))).
		AndWhere(ats("attendee_id = ?"), results[0].ID).
		Fetch(&claims)
	if err != nil {
		return nil, fmt.Errorf("reading claims for attendee: %w", err)
	}
	newAttendee := results[0]
	newAttendee.Claims = claims
	return &newAttendee, nil
}

const eventSlotTable = "event_slot"

// CreateEventSlot saves a slot in the database.
func (s *SQLStorage) CreateEventSlot(e *EventSlot) (*EventSlot, error) {
	results := []EventSlot{}
	insertMap := map[string]interface{}{
		"event_id":            e.Event.ID,
		"name":                e.Name,
		"description":         e.Description,
		"cost":                e.Cost,
		"capacity":            e.Capacity,
		"start_date":          e.StartDate,
		"end_date":            e.EndDate,
		"purchaseable_from":   e.PurchaseableFrom,
		"purchaseable_until":  e.PurchaseableUntil,
		"available_to_public": e.AvailableToPublic,
	}
	if e.DependsOn != nil {
		insertMap["depends_on_id"] = e.DependsOn.ID
	}
	if e.Event != nil {
		insertMap["event_id"] = e.Event.ID
	}
	err := chain.New(s.conn).Insert(insertMap).
		Table(eventSlotTable).Returning("*").
		Fetch(&results)
	if err != nil {
		return nil, fmt.Errorf("creating new attendee: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("attendee was not created")
	}
	return &results[0], nil
}

type wrapEventSlot struct {
	EventSlot
	DependsOnID uint64 `gaum:"field_name:depends_on_id"`
	EventID     uint64 `gaum:"field_name:event_id"`
}

// ReadEventSlotByID returns an event slot identified by the passed ID.
func (s *SQLStorage) ReadEventSlotByID(id uint64) (*EventSlot, error) {
	return s.readEventSlotByID(id, true)
}

func (s *SQLStorage) readEventSlotByID(id uint64, loadDeps bool) (*EventSlot, error) {
	results := []wrapEventSlot{}
	err := chain.New(s.conn).Select("*").
		From(eventSlotTable).
		AndWhere("id = ?", id).Fetch(&results)
	if err != nil {
		return nil, fmt.Errorf("reading event slots by id: %w", err)
	}
	if len(results) == 0 {
		return nil, nil
	}
	slot := results[0].EventSlot
	if results[0].DependsOnID != 0 && loadDeps {
		slot.DependsOn, err = s.readEventSlotByID(results[0].DependsOnID, false)
		if err != nil {
			return nil, fmt.Errorf("loading dependency: %w", err)
		}
	}

	events := []def.Event{}
	err = chain.New(s.conn).Select("*").
		From("event").
		AndWhere("id = ?", results[0].EventID).Fetch(&events)
	if err != nil {
		return nil, fmt.Errorf("reading event by id: %w", err)
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("could not find event for slot")
	}
	slot.Event = &events[0]
	return &slot, nil
}

// UpdateEventSlot updates event slot fields from the passed instance
func (s *SQLStorage) UpdateEventSlot(e *EventSlot) error {
	updateMap := map[string]interface{}{
		"event_id":            e.Event.ID,
		"name":                e.Name,
		"description":         e.Description,
		"cost":                e.Cost,
		"capacity":            e.Capacity,
		"start_date":          e.StartDate,
		"end_date":            e.EndDate,
		"purchaseable_from":   e.PurchaseableFrom,
		"purchaseable_until":  e.PurchaseableUntil,
		"available_to_public": e.AvailableToPublic,
	}
	if e.DependsOn != nil {
		updateMap["depends_on_id"] = e.DependsOn.ID
	}
	if e.Event != nil {
		updateMap["event_id"] = e.Event.ID
	}
	affected, err := chain.New(s.conn).UpdateMap(updateMap).
		Table(eventSlotTable).AndWhere("id = ?", e.ID).
		ExecResult()
	if err != nil {
		return fmt.Errorf("creating new attendee: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("attendee was not updated")
	}
	return nil
}

// CreateSlotClaim saves a slot claim and returns it with the populated ID
func (s *SQLStorage) CreateSlotClaim(slotClaim *SlotClaim) (*SlotClaim, error) {
	var err error
	// FIXME: Add a check for capacity not exceeded on event.
	for i := 0; i < 3; i++ {
		q := chain.New(s.conn)
		results := []SlotClaim{}
		err = q.Insert(map[string]interface{}{
			"ticket_id":     slotClaim.TicketID,
			"redeemed":      slotClaim.Redeemed,
			"event_slot_id": slotClaim.EventSlot.ID,
		}).
			Table(tableSlotClaims).
			Returning("id, ticket_id, redeemed").Fetch(&results)

		if err != nil {
			// there is no SQL for "on error change the inserting statement", only to change
			// the existing one.
			if err, ok := err.(pgx.PgError); ok {
				if pgerrcode.IsIntegrityConstraintViolation(err.Code) {
					// if this clashes again entropy might be broken, check if not 2020
					slotClaim.TicketID = uuid.NewV4().String()
					continue
				}
			}
			return nil, fmt.Errorf("saving slot claim: %w", err)
		}
		if len(results) == 0 {
			return nil, nil
		}
		results[0].EventSlot = slotClaim.EventSlot
		return &results[0], nil
	}
	return nil, fmt.Errorf("failed to insert: %w", err)
}

const (
	tableAttendee               = "attendee"
	tableAttendeeSlotClaims     = "attendee_to_slot_claims"
	slotClaimIDUniqueConstraint = "slot_claim_id_is_unique"
)

// UpdateAttendee saves the passed attendee attributes on top of the existing one.
func (s *SQLStorage) UpdateAttendee(attendee *Attendee) (*Attendee, error) {
	rows, err := chain.New(s.conn).UpdateMap(map[string]interface{}{
		"email":        attendee.Email,
		"coc_accepted": attendee.CoCAccepted,
	}).Table(tableAttendee).
		AndWhere("id = ?", attendee.ID).ExecResult()
	if err != nil {
		return nil, fmt.Errorf("updating attendee: %w", err)
	}
	if rows == 0 { // attendee does not exist
		return nil, nil
	}

	for i := range attendee.Claims {
		c := attendee.Claims[i]
		err := chain.New(s.conn).Insert(map[string]interface{}{
			"attendee_id":   attendee.ID,
			"slot_claim_id": c.Redeemed,
		}).Table(tableAttendeeSlotClaims).
			OnConflict(func(c *chain.OnConflict) {
				// This claim was someone else's, this might be the result of transfering.
				c.OnConstraint(slotClaimIDUniqueConstraint).
					DoUpdate().
					Set("attendee_id", attendee.ID)
			}).Exec()
		if err != nil {
			return nil, fmt.Errorf("updating attendee claims: %w", err)
		}
	}
	return attendee, nil
}

const (
	tableClaimPayment                = "claim_payment"
	tableFinancialInstrumentMoney    = "payment_method_money"
	tableMoneyToPayment              = "payment_method_money_to_claim_payment"
	tableFinancialInstrumentDiscount = "payment_method_event_discount"
	tableDiscountToPayment           = "payment_method_event_discount_to_claim_payment"
	tableFinancialInstrumentCredit   = "payment_method_credit_note"
	tableCreditToPayment             = "payment_method_credit_note_to_claim_payment"
)

func insertMoneyPayment(conn connection.DB, claimPaymentID uint64, payment *PaymentMethodMoney) (*PaymentMethodMoney, error) {
	if payment.ID != 0 { // BIGSERIAL starts in 1
		return payment, nil
	}
	money := []PaymentMethodMoney{}
	err := chain.New(conn).Insert(map[string]interface{}{
		"amount": payment.Amount,
		"ref":    payment.PaymentRef,
	}).Table(tableFinancialInstrumentMoney).
		Returning("*").Fetch(&money)
	if err != nil {
		return nil, fmt.Errorf("inserting money payment: %w", err)
	}
	if len(money) == 0 {
		return nil, fmt.Errorf("failed to insert money payment")
	}
	err = chain.New(conn).Insert(map[string]interface{}{
		"payment_method_money_id": money[0].ID,
		"claim_payment_id":        claimPaymentID,
	}).Table(tableMoneyToPayment).Exec()
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument money to payment: %w", err)
	}
	return &money[0], nil
}

func insertDiscountPayment(conn connection.DB, claimPaymentID uint64, payment *PaymentMethodConferenceDiscount) (*PaymentMethodConferenceDiscount, error) {
	if payment.ID != 0 { // BIGSERIAL starts in 1
		return payment, nil
	}
	discount := []PaymentMethodConferenceDiscount{}
	err := chain.New(conn).Insert(map[string]interface{}{
		"amount": payment.Amount,
		"detail": payment.Detail,
	}).Table(tableFinancialInstrumentDiscount).
		Returning("*").Fetch(&discount)
	if err != nil {
		return nil, fmt.Errorf("inserting discount payment: %w", err)
	}
	if len(discount) == 0 {
		return nil, fmt.Errorf("failed to insert discount payment")
	}

	err = chain.New(conn).Insert(map[string]interface{}{
		"payment_method_event_discount_id": discount[0].ID,
		"claim_payment_id":                 claimPaymentID,
	}).Table(tableDiscountToPayment).Exec()
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument discount to payment: %w", err)
	}
	return &discount[0], nil

}

func insertCreditPayment(conn connection.DB, claimPaymentID uint64, payment *PaymentMethodCreditNote) (*PaymentMethodCreditNote, error) {
	if payment.ID != 0 { // BIGSERIAL starts in 1
		return payment, nil
	}
	credit := []PaymentMethodCreditNote{}
	err := chain.New(conn).Insert(map[string]interface{}{
		"amount": payment.Amount,
		"detail": payment.Detail,
	}).Table(tableFinancialInstrumentCredit).
		Returning("*").Fetch(&credit)
	if err != nil {
		return nil, fmt.Errorf("inserting credit payment: %w", err)
	}
	if len(credit) == 0 {
		return nil, fmt.Errorf("failed to insert credit payment")
	}

	err = chain.New(conn).Insert(map[string]interface{}{
		"payment_method_credit_note_id": credit[0].ID,
		"claim_payment_id":              claimPaymentID,
	}).Table(tableCreditToPayment).Exec()
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument credit to payment: %w", err)
	}

	return &credit[0], nil
}

// CreateClaimPayment Creates c ClaimPayment record and asociates it with all the relevant payments.
func (s *SQLStorage) CreateClaimPayment(c *ClaimPayment) (*ClaimPayment, error) {
	claimPayments := []ClaimPayment{}
	err := chain.New(s.conn).Insert(map[string]interface{}{
		"invoice": c.Invoice,
	}).Table(tableClaimPayment).Returning("*").Fetch(&claimPayments)
	if err != nil {
		return nil, fmt.Errorf("inserting payment for claims: %w", err)
	}
	if len(claimPayments) == 0 {
		return nil, fmt.Errorf("claim payment was not created")
	}

	processedPayments := make([]FinancialInstrument, len(c.Payment), len(c.Payment))

	for i, cp := range c.Payment {
		switch payment := cp.(type) {
		case *PaymentMethodMoney:
			processedPayments[i], err = insertMoneyPayment(s.conn, claimPayments[0].ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting money payment: %w", err)
			}
		case *PaymentMethodConferenceDiscount:
			processedPayments[i], err = insertDiscountPayment(s.conn, claimPayments[0].ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting discount payment: %w", err)
			}
		case *PaymentMethodCreditNote:
			processedPayments[i], err = insertCreditPayment(s.conn, claimPayments[0].ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting credit payment: %w", err)
			}
		default:
			return nil, fmt.Errorf("not sure how to process payments of type %T", cp)
		}
	}
	newClaim := claimPayments[0]
	newClaim.ClaimsPayed = c.ClaimsPayed
	newClaim.Payment = processedPayments
	return &newClaim, nil
}

// UpdateClaimPayment saves the invoice and payments of this claim payment assuming it exists
func (s *SQLStorage) UpdateClaimPayment(c *ClaimPayment) (*ClaimPayment, error) {
	updated, err := chain.New(s.conn).UpdateMap(map[string]interface{}{
		"invoice": c.Invoice,
	}).Table(tableClaimPayment).
		AndWhere("id = ?", c.ID).
		ExecResult()
	if err != nil {
		return nil, fmt.Errorf("inserting payment for claims: %w", err)
	}
	if updated == 0 {
		return nil, fmt.Errorf("claim payment was not found")
	}

	processedPayments := make([]FinancialInstrument, len(c.Payment), len(c.Payment))

	for i, cp := range c.Payment {
		switch payment := cp.(type) {
		case *PaymentMethodMoney:
			processedPayments[i], err = insertMoneyPayment(s.conn, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("processing financial instrument money to payment: %w", err)
			}
		case *PaymentMethodConferenceDiscount:
			processedPayments[i], err = insertDiscountPayment(s.conn, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("processing financial instrument discount to payment: %w", err)
			}
		case *PaymentMethodCreditNote:
			processedPayments[i], err = insertCreditPayment(s.conn, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("processing financial instrument credit to payment: %w", err)
			}
		default:
			return nil, fmt.Errorf("not sure how to process payments of type %T", cp)
		}
	}
	newClaim := ClaimPayment{
		ID:          c.ID,
		ClaimsPayed: c.ClaimsPayed,
		Payment:     processedPayments,
		Invoice:     c.Invoice,
	}
	return &newClaim, nil
}

const attendeesToSlotTable = "attendees_to_slot_claims"

// ChangeSlotClaimOwner changes the passed claims owner from source to target
func (s *SQLStorage) ChangeSlotClaimOwner(slots []SlotClaim, source *Attendee, target *Attendee) (*Attendee, *Attendee, error) {
	if source == nil || target == nil {
		return nil, nil, fmt.Errorf("either source or target is undefined")
	}

	if len(slots) == 0 {
		return nil, nil, fmt.Errorf("no slots to transfer")
	}

	if len(slots) > len(source.Claims) {
		return nil, nil, fmt.Errorf("the passed source lacks those claims")
	}

	claimIDs := make([]uint64, 0, len(slots))
	claimIDsIndex := map[uint64]bool{}
	for _, slot := range slots {
		if slot.ID == 0 {
			return nil, nil, fmt.Errorf("some slot claims lack IDs, perhaps the have not been saved yet")
		}
		claimIDs = append(claimIDs, slot.ID)
		claimIDsIndex[slot.ID] = true
	}
	affected, err := chain.New(s.conn).UpdateMap(map[string]interface{}{
		"attendee_id": target.ID,
	}).AndWhere("attendee_id = ?", source.ID).
		AndWhere("slot_claim_id IN (?)", claimIDs).ExecResult()
	if err != nil {
		return nil, nil, fmt.Errorf("chaingin slot claims ownershio: %w", err)
	}
	if int64(len(slots)) != affected {
		return nil, nil, fmt.Errorf("got %d claims to change but only changed %d", len(slots), affected)
	}

	newClaims := make([]SlotClaim, 0, len(source.Claims)-len(claimIDs))
	for i := range source.Claims {
		if claimIDsIndex[source.Claims[i].ID] {
			target.Claims = append(target.Claims, source.Claims[i])
			continue
		}
		newClaims = append(newClaims, source.Claims[i])
	}
	source.Claims = newClaims
	return source, target, nil
}
