package models

type Menu struct {
	BasicModel
	StoreId   uint64 `json:"store_id" validate:"required"`
	Name      string `json:"name"`
	QrCodeSrc string `json:"qr_code_src"`
	Url       string `json:"url"`
}
