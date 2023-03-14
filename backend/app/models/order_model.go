package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"golang.org/x/exp/slices"
)

type Order struct {
	BasicModel
	StoreId               uint64             `json:"store_id"`
	BillingAddress        BillingAddress     `json:"billing_address" gorm:"default:null"`
	BrowserIP             string             `json:"browser_ip" gorm:"default:null"`
	BuyerAcceptsMarketing bool               `json:"buyer_accepts_marketing" gorm:"default:null"`
	CancelReason          string             `json:"cancel_reason" gorm:"default:null"`
	CancelledOn           time.Time          `json:"cancelled_on" gorm:"default:null"`
	CartToken             string             `json:"cart_token"`
	ClientDetails         ClientDetails      `json:"client_details" gorm:"default:null"`
	ClosedOn              time.Time          `json:"closed_on" gorm:"default:null"`
	Currency              string             `json:"currency" gorm:"default:null"`
	Customer              Customer           `json:"customer" gorm:"default:null"`
	DiscountCodes         DiscountCodeArray  `json:"discount_codes" gorm:"default:null"`
	Email                 string             `json:"email" gorm:"default:null"`
	FinancialStatus       string             `json:"financial_status" gorm:"default:null"`
	Status                string             `json:"status" gorm:"default:null"`
	Fulfillments          Fulfillments       `json:"fulfillments" gorm:"default:null"`
	FulfillmentStatus     string             `json:"fulfillment_status" gorm:"default:null"`
	Tags                  string             `json:"tags" gorm:"default:null"`
	LandingSite           string             `json:"landing_site" gorm:"default:null"`
	LineItems             LineItemArray      `json:"line_items" gorm:"default:null"`
	Name                  string             `json:"name" gorm:"default:null"` // name of order. Eg: "#0001"
	Note                  string             `json:"note" gorm:"default:null"`
	NoteAttributes        NoteAttributeArray `json:"note_attributes" gorm:"default:null"`
	Number                uint64             `json:"number" gorm:"default:null"`       // The unique number that identifies the Order for the Shop. This number is a self-incrementing number and starts at 1000. Eg:
	OrderNumber           uint64             `json:"order_number" gorm:"default:null"` // The unique number that identifies the Order. This number is used by Shop owners and customers.
	PaymentGatewayNames   StringArray        `json:"payment_gateway_names" gorm:"default:null"`
	ProcessedOn           time.Time          `json:"processed_on" gorm:"default:null"`
	ProcessingMethod      string             `json:"processing_method" gorm:"default:null"`
	ReferringSite         string             `json:"referring_site" gorm:"default:null"`
	Refunds               string             `json:"refunds" gorm:"default:null"`
	ShippingAddress       ShippingAddress    `json:"shipping_address" gorm:"default:null"`
	ShippingLines         ShippingLineArray  `json:"shipping_lines" gorm:"default:null"`
	SourceName            string             `json:"source_name" gorm:"default:null"`
	Token                 string             `json:"token" gorm:"default:null"`
	TotalDiscount         float64            `json:"total_discount" gorm:"default:null"`
	TotalLineItemsPrice   float64            `json:"total_line_items_price" gorm:"default:null"`
	TotalPrice            float64            `json:"total_price" gorm:"default:null"`
	TotalWeight           int                `json:"total_weight" gorm:"default:null"`
	ModifiedOn            time.Time          `json:"modified_on" gorm:"default:null"`
	Gateway               string             `json:"gateway" gorm:"default:null"`
}

func (*Order) GetColumnsUpdateOnConflict() []string {
	return []string{
		"billing_address", "browser_ip", "buyer_accepts_marketing", "cancel_reason", "note_attributes", "line_items", "total_line_items_price", "total_price", "total_weight",
	}
}

type StringArray []string

func (sla *StringArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla StringArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

func (o *Order) UpdatePaymentGatewayNames(paymentMethod string) {
	if !slices.Contains(o.PaymentGatewayNames, paymentMethod) {
		o.PaymentGatewayNames = append(o.PaymentGatewayNames, paymentMethod)
	}
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

func (sla *BillingAddress) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla BillingAddress) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ClientDetails struct {
	AcceptLanguage string `json:"accept_language"`
	BrowserHeight  string `json:"browser_height"`
	BrowserWidth   string `json:"browser_width"`
	BrowserIP      string `json:"browser_ip"`
	SessionHash    string `json:"session_hash"`
	UserAgent      string `json:"user_agent"`
}

func (sla *ClientDetails) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ClientDetails) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
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

func (sla *Customer) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla Customer) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type DiscountCode struct {
	Amount int    `json:"amount"`
	Code   string `json:"code"`
	Type   string `json:"type"`
}

type DiscountCodeArray []DiscountCode

func (sla *DiscountCodeArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla DiscountCodeArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
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

func (sla *Fulfillments) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla Fulfillments) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
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
	Taxable             bool    `json:"taxable"`
	TaxLines            string  `json:"tax_lines"`
	TotalDiscount       float64 `json:"total_discount"`
}

type LineItemArray []LineItem

func (sla *LineItemArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla LineItemArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type NoteAttribute struct {
	Name      string `json:"name"`
	Attribute string `json:"attribute"`
}

type NoteAttributeArray []NoteAttribute

func (sla *NoteAttributeArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla NoteAttributeArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
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

func (sla *ShippingAddress) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ShippingAddress) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ShippingLine struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Source   string  `json:"source"`
	Title    string  `json:"title"`
	TaxLines string  `json:"tax_lines"`
}

type ShippingLineArray []ShippingLine

func (sla *ShippingLineArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ShippingLineArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type PayOrderReqBody struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=cash bank_transfer vnpay card"`
}
