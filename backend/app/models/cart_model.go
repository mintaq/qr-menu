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

func (c *Cart) UpdateCart(index, quantity int) *Cart {
	if index >= 0 && index < len(c.Items) {
		if quantity == 0 {
			c.Items[index] = c.Items[len(c.Items)-1]
			c.Items = c.Items[:len(c.Items)-1]
		} else {
			c.Items[index].Quantity = quantity
		}
	}

	return c
}
