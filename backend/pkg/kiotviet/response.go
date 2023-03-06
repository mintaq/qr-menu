package kiotviet

type ConnectTokenResponse struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn     int `json:"expires_in"`
	TokenType     string `json:"token_type"`
	Scope     string `json:"scope"`
}

type ProductsResponse struct {
	Total int `json:"total"`
	PageSize int `json:"pageSize"`
	Data []ProductResponse `json:"data"`
}

type ProductResponse struct {
	Content     string       `json:"description" gorm:"default:null"`
	Summary     string       `json:"summary" gorm:"default:null"`
	CreatedOn   string    `json:"createdDate" gorm:"default:null"`
	Alias       string       `json:"alias"`
	ProductId   uint64       `json:"id"`
	Images     	[]string     `json:"images" gorm:"default:null"`
	Options     OptionsResponse  `json:"attributes" gorm:"default:null"`
	ProductType string       `json:"categoryName" gorm:"default:null"`
	PublishedOn string    `json:"startDate" gorm:"default:null"`
	Tags        string       `json:"tags" gorm:"default:null"`
	ProductName string       `json:"name" gorm:"column:product_name" validate:"required"`
	ModifiedOn  string    `json:"modifiedDate" gorm:"default:null"`
	Variants    string    `json:"variants" gorm:"default:null"`
	Vendor      string       `json:"vendor" gorm:"default:null"`
}

type OptionsResponse []OptionResponse

type OptionResponse struct {
	ProductId uint64   `json:"productId"`
	Name      string   `json:"attributeName"`
	Value    string `json:"attributeValue"`
}

type CollectionsResponse struct {
	Total int `json:"total"`
	PageSize int `json:"pageSize"`
	Data []CollectionResponse `json:"data"`
}

type CollectionResponse struct {
	CollectionId uint64          `json:"categoryId"`
	Description  string          `json:"description" gorm:"default:null"`
	Alias        string          `json:"alias" gorm:"default:null"`
	Name         string          `json:"categoryName" validate:"required"`
}
