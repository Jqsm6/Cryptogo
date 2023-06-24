Get information about the invoice.

#### Request

- URL: `GET /v1/invoice/info`
- Headers:
  - Content-Type: `application/json`
- Body:
  - `id` (string, required): ID assigned at the time of creation.

Example:

```json
{
  "id": "abcd1234"
}
```

#### Response

- Status Code: `200 OK`
- Body:
  - `id` (string): ID of the created invoice
  - `state` (string): payment status. Can be 'paid' and 'notpayed'
  - `currency` (string): available only ETH, BTC, BNB
  - `amount` (string): the amount to be paid
  - `from_address` (string): the address from which the payment will be made
  - `to_address` (string): address to pay the invoice

Example:

```json
{
{
    "id": "abcd1234",
    "state": "notpayed",
    "currency": "ETH",
    "amount": "0.05",
    "to_address": "0xA1B2C3D4E5F6G7H8I9J0K1L2M3N4O5P6Q7R8S9T0",
    "from_address": "0x1234567890123456789012345678901234567890"
}
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
        "message": "Invoice with this ID was not found."
    }
}
```