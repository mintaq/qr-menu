package models

import "time"

type Order struct {
	BasicModel
	StoreId               uint64          `json:"store_id"`
	BillingAddress        BillingAddress  `json:"billing_address"`
	BrowserIP             string          `json:"browser_ip"`
	BuyerAcceptsMarketing bool            `json:"buyer_accepts_marketing"`
	CancelReason          string          `json:"cancel_reason"`
	CancelledOn           time.Time       `json:"cancelled_on"`
	CartToken             string          `json:"cart_token"`
	ClientDetails         ClientDetails   `json:"client_details"`
	ClosedOn              time.Time       `json:"closed_on"`
	Currency              string          `json:"currency"`
	Customer              Customer        `json:"customer"`
	DiscountCodes         []DiscountCode  `json:"discount_codes"`
	Email                 string          `json:"email"`
	FinancialStatus       string          `json:"financial_status"`
	Status                string          `json:"status"`
	Fulfillments          Fulfillments    `json:"fulfillments"`
	FulfillmentStatus     string          `json:"fulfillment_status"`
	Tags                  string          `json:"tags"`
	LandingSite           string          `json:"landing_site"`
	LineItems             []LineItem      `json:"line_items"`
	Name                  string          `json:"name"` // name of order. Eg: "#0001"
	Note                  string          `json:"note"`
	NoteAttributes        []NoteAttribute `json:"note_attributes"`
	Number                uint64          `json:"number"`
	OrderNumber           uint64          `json:"order_number"`
	PaymentGatewayNames   []string        `json:"payment_gate_way_names"`
	ProcessedOn           time.Time       `json:"processed_on"`
	ProcessingMethod      string          `json:"processing_method"`
	ReferringSite         string          `json:"referring_site"`
	Refunds               string          `json:"refunds"`
	ShippingAddress       ShippingAddress `json:"shipping_address"`
	ShippingLines         []ShippingLine  `json:"shipping_lines"`
	SourceName            string          `json:"source_name"`
	Token                 string          `json:"token"`
	TotalDiscount         float64         `json:"total_discount"`
	TotalLineItemsPrice   float64         `json:"total_line_items_price"`
	TotalPrice            float64         `json:"total_price"`
	TotalWeight           int             `json:"total_weight"`
	ModifiedOn            time.Time       `json:"modified_on"`
	Gateway               string          `json:"gateway"`
}

type BillingAddress struct {
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	City         string `json:"city"`
	Company      string `json:"company"`
	Country      string `json:"country"`
	FirstName    string `json:"first_name"`
	Id           uint64 `json:"id"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Province     string `json:"province"`
	Zip          string `json:"zip"`
	Name         string `json:"name"`
	ProvinceCode string `json:"province_code"`
	CountryCode  string `json:"country_code"`
	Default      bool   `json:"default"`
}

type ClientDetails struct {
	AcceptLanguage string `json:"accept_language"`
	BrowserHeight  string `json:"browser_height"`
	BrowserWidth   string `json:"browser_width"`
	BrowserIP      string `json:"browser_ip"`
	SessionHash    string `json:"session_hash"`
	UserAgent      string `json:"user_agent"`
}

type Customer struct {
	AcceptsMarketing bool      `json:"accepts_marketing"`
	CreatedOn        time.Time `json:"created_on"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	Id               uint64    `json:"id"`
	LastName         string    `json:"last_name"`
	Note             string    `json:"note"`
	OrdersCount      int       `json:"orders_count"`
	State            string    `json:"state"`
	TotalSpent       float64   `json:"total_spent"`
	ModifiedOn       time.Time `json:"modified_on"`
	Tags             string    `json:"tags"`
}

type DiscountCode struct {
	Amount int    `json:"amount"`
	Code   string `json:"code"`
	Type   string `json:"type"`
}

type Fulfillments struct {
	CreatedOn       time.Time `json:"created_on"`
	Id              uint64    `json:"id"`
	OrderId         uint64    `json:"order_id"`
	Status          string    `json:"status"`
	TrackingCompany string    `json:"tracking_company"`
	TrackingNumber  string    `json:"tracking_number"`
	ModifiedOn      time.Time `json:"modified_on"`
}

type LineItem struct {
	Id                  uint64  `json:"id"`
	FulfillableQuantity int     `json:"fulfillable_quantity"`
	FulfillmentService  string  `json:"fulfillment_service"`
	Grams               int     `json:"grams"`
	Price               float64 `json:"price"`
	ProductId           uint64  `json:"product_id"`
	Quantity            int     `json:"quantity"`
	RequiresShipping    bool    `json:"requires_shipping"`
	SKU                 string  `json:"sku"`
	Title               string  `json:"title"` // title of product
	VariantId           uint64  `json:"variant_id"`
	VariantTitle        string  `json:"variant_title"`
	Vendor              string  `json:"vendor"`
	Name                string  `json:"name"` // title of variant
	GiftCard            bool    `json:"gift_card"`
	Taxable             bool    `json:"table"`
	TaxLines            string  `json:"tax_lines"`
	TotalDiscount       float64 `json:"total_discount"`
}

type NoteAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ShippingAddress struct {
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	City         string `json:"city"`
	Company      string `json:"company"`
	Country      string `json:"country"`
	FirstName    string `json:"first_name"`
	Id           uint64 `json:"id"`
	LastName     string `json:"last_name"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Phone        string `json:"phone"`
	Province     string `json:"province"`
	Zip          string `json:"zip"`
	Name         string `json:"name"`
	ProvinceCode string `json:"province_code"`
	CountryCode  string `json:"country_code"`
}

type ShippingLine struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Source   string  `json:"source"`
	Title    string  `json:"title"`
	TaxLines string  `json:"tax_lines"`
}
