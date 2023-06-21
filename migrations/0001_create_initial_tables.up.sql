DROP TABLE IF EXISTS payments CASCADE;

CREATE TABLE payments
(
    invoice_id UNIQUE NOT NULL,
    state TEXT NOT NULL,
    currency TEXT NOT NULL,
    amount TEXT NOT NULL,
    to_address TEXT NOT NULL,
    from_address TEXT NOT NULL,
);

CREATE INDEX IF NOT EXISTS payments_idx
    ON payments (invoice_id);