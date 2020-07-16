package model

type Order struct {
	Entity string      `json:"entity"`
	Item   []OrderItem `json:"item"`
	Date   string      `json:"trandate"` // same with Trandate
}

type OrderItem struct {
	Item     string `json:"item"`
	Rate     string `json:"rate"`
	Quantity string `json:"quantity"`
}

type OrderNetsuite struct {
	URI  string `json:"uri"`
	Data Order  `json:"data"`
}
