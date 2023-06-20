package repository

const (
	createInvoice = `INSERT INTO payments (invoice_id, state, currency, amount, to_address, from_address) 
	VALUES ($1, $2, $3, $4, $5, $6)`
	infoInvoice        = `SELECT invoice_id, state, to_address, amount, currency, from_address FROM payments WHERE invoice_id = $1`
	changeInvoiceState = `UPDATE payments SET state = 'paid' WHERE invoice_id = $1`
	checkID            = `SELECT COUNT(*) FROM payments WHERE invoice_id = $1`
)
