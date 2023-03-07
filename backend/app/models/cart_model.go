package models

// Cart struct to describe cart object.
type Cart struct {
	UserToken  string     `json:"user_token"`
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
	Note       string     `json:"note"`
	ItemsCount uint       `json:"items_count"`
}

type CartItem struct {
	Product
	Quantity int `json:"quantity"`
}

type AddItemToCartReqBody struct {
	ProductId uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (c *Cart) CountTotalPrice() *Cart {
	var totalPrice float64
	for i := range c.Items {
		totalPrice += c.Items[i].Price
	}

	c.TotalPrice = totalPrice
	return c
}

func (c *Cart) CountTotalItems() *Cart {
	var totalItems uint
	for i := range c.Items {
		totalItems += uint(c.Items[i].Quantity)
	}

	c.ItemsCount = totalItems
	return c
}

func (c *Cart) UpdateCountableFields() *Cart {
	var totalItems uint
	var totalPrice float64
	for i := range c.Items {
		totalItems += uint(c.Items[i].Quantity)
		totalPrice += c.Items[i].Price
	}

	c.ItemsCount = totalItems
	c.TotalPrice = totalPrice
	return c
}

func (c *Cart) UpdateCartByIndex(index, quantity int) *Cart {
	if index < 0 || index >= len(c.Items) {
		return c
	}

	if quantity == 0 {
		c.Items = append(c.Items[:index], c.Items[index+1:]...)
	} else {
		c.Items[index].Quantity = quantity
	}

	return c
}


func (c *Cart) HasProduct(productId uint64) bool {
	for index := range c.Items {
		if c.Items[index].Product.ProductId == productId {
			return true
		}
	}

	return false
}

func (c *Cart) UpdateCartByProductId(productId uint64, quantity int) *Cart {
	for index := range c.Items {
		if c.Items[index].Product.ProductId == productId {
			c.Items[index].Quantity += quantity
			return c
		}
	}

	return c
}
