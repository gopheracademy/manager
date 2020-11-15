CREATE TABLE event_slot (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT, 
    cost DECIMAL, 
    capacity INT, 
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    depends_on BIGINT, 
    purchaseable_from TIMESTAMP WITH TIME ZONE,
    purchaseable_until TIMESTAMP WITH TIME ZONE,
    available_to_public BOOLEAN,
    FOREIGN KEY(depends_on) REFERENCES event_slot(id)
);

CREATE TABLE slot_claim (
    id UUID PRIMARY KEY,
    event_slot_id BIGINT, 
    ticket_id UUID,
    redeemed BOOLEAN,
    FOREIGN KEY(event_slot_id) REFERENCES event_slot(id)
);

CREATE TABLE attendee (
    id BIGSERIAL PRIMARY KEY,
    email BIGINT, 
    coc_accepted BOOLEAN
);

CREATE TABLE attendee_to_slot_claims (
    attendee_id BIGINT,
    slot_claim_id UUID UNIQUE,
    FOREIGN KEY (attendee_id) REFERENCES attendee(id),
    FOREIGN KEY (slot_claim_id) REFERENCES slot_claim(id)
);

CREATE TABLE claim_payment (
    id UUID PRIMARY KEY,
    invoice TEXT -- just in case we need to store the whole thing.
);

CREATE TABLE payment_method_money (
    id BIGSERIAL PRIMARY KEY,
    amount DECIMAL,
    ref VARCHAR(250)
);

CREATE TABLE payment_method_money_to_claim_payment (
    payment_method_money_id BIGINT,
    claim_payment_id UUID,
    FOREIGN KEY (payment_method_money_id) REFERENCES payment_method_money(id),
    FOREIGN KEY (claim_payment_id) REFERENCES claim_payment(id)
);

CREATE TABLE payment_method_credit_note (
    id BIGSERIAL PRIMARY KEY,
    amount DECIMAL,
    detail VARCHAR(250)
);

CREATE TABLE payment_method_credit_note_to_claim_payment (
    payment_method_credit_note_id BIGINT,
    claim_payment_id UUID,
    FOREIGN KEY (payment_method_credit_note_id) REFERENCES payment_method_credit_note(id),
    FOREIGN KEY (claim_payment_id) REFERENCES claim_payment(id)
);

CREATE TABLE payment_method_event_discount (
    id BIGSERIAL PRIMARY KEY,
    amount DECIMAL,
    detail VARCHAR(250)
);

CREATE TABLE payment_method_event_discount_to_claim_payment (
    payment_method_event_discount_id BIGINT,
    claim_payment_id UUID,
    FOREIGN KEY (payment_method_event_discount_id) REFERENCES payment_method_event_discount(id),
    FOREIGN KEY (claim_payment_id) REFERENCES claim_payment(id)
);

