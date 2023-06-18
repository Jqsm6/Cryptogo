package repository

const (
	createInvoice = `INSERT INTO payments (id_, status, currency, amount, from_address, to_address, private_key) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	checkInvoice                = `SELECT status, to_address, amount, currency FROM payments WHERE id_ = $1`
	changeInvoiceStatusEndpoint = `UPDATE payments SET status = 1 WHERE id_ = $1`
)
