#### Request

- URL: `POST /v1/invoice/create`
- Headers:
  - Content-Type: `application/json`
- Body:
  - `currency` (string, required): available only ETH, BTC, BNB
  - `amount` (float64, required): the amount to be paid
  - `from_address` (string, required): the address from which the payment will be made

Example:

```json
{
  "currency": "ETH",
  "amount": 0.05,
  "from_address": "0x1234567890123456789012345678901234567890"
}
```

#### Response

- Status Code: `201 Created`
- Body:
  - `id` (string): ID of the created invoice
  - `to_address` (string): address to pay the invoice

Example:

```json
{
  "id": "abcd1234",
  "to_address": "0xA1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P6Q7R8S9T0"
}
```

#### Errors

- Status Code: `400 Bad Request`
- Body:
  - `error` (object): object containing the 'code' and 'message' fields
  - `code` (int): code that describes the error
  - `message` (string): description of the error

Example:

```json
{
  "error": {
    "code": 400,
    "message": "Invalid request body. Use the documentation."
  }
}
```
