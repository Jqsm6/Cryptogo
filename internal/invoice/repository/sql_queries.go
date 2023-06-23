package repository

const (
	create      = `INSERT INTO payments (invoice_id, state, currency, amount, to_address, from_address) VALUES ($1, $2, $3, $4, $5, $6)`
	info        = `SELECT invoice_id, state, to_address, amount, currency, from_address FROM payments WHERE invoice_id = $1`
	changeState = `UPDATE payments SET state = 'paid' WHERE invoice_id = $1`
	checkID     = `SELECT COUNT(*) FROM payments WHERE invoice_id = $1`
	checkHash   = `SELECT COUNT(*) FROM payments WHERE transaction_hash = $1`
	updateHash  = `UPDATE payments SET transaction_hash = $1 WHERE invoice_id = $2`
)
