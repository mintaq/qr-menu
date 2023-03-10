package repository

const (
	FULFILLMENT_STATUS_FULFILLED string = "fulfilled"
	FULFILLMENT_STATUS_NULL      string = "null"
	FULFILLMENT_STATUS_PARTIAL   string = "partial"

	FINANCIAL_STATUS_PENDING            string = "pending"
	FINANCIAL_STATUS_AUTHORIZED         string = "authorized"
	FINANCIAL_STATUS_PARTIALLY_PAID     string = "partially_paid"
	FINANCIAL_STATUS_PAID               string = "paid"
	FINANCIAL_STATUS_PARTIALLY_REFUNDED string = "partially_refunded"
	FINANCIAL_STATUS_REFUNDED           string = "refunded"
	FINANCIAL_STATUS_VOIDED             string = "voided"

	ORDER_STATUS_OPEN      string = "open"
	ORDER_STATUS_CLOSED    string = "closed"
	ORDER_STATUS_CANCELLED string = "cancelled"
)
