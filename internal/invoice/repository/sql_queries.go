package repository

const (
	createInvoice = `INSERT INTO payments (invoice_id, state, currency, amount, to_address, private_key) 
	VALUES ($1, $2, $3, $4, $5, $6)`
	infoInvoice                = `SELECT state, to_address, amount, currency FROM payments WHERE invoice_id = $1`
	changeInvoiceState = `UPDATE payments SET state = 'paid' WHERE invoice_id = $1`
)
