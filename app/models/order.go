package models

type Order struct {
	OrderID             int64     `bson:"id" json:"orderID"`
	ExternalReferenceID string    `bson:"externalReferenceID" json:"externalReferenceID"`
	Channel             string    `bson:"channel" json:"channel"`
	PurchaseDate        string    `bson:"purchaseDate" json:"purchaseDate"`
	TotalValue          float64   `bson:"totalValue" json:"totalValue"`
	Buyer               Buyer     `bson:"buyer" json:"buyer"`
	Products            []Product `bson:"products" json:"products"`
	Status              string    `bson:"status" json:"status"`
	Events              []Event   `bson:"events" json:"events"`
}

type Buyer struct {
	FirstName      string `bson:"firstName" json:"firstName"`
	LastName       string `bson:"lastName" json:"lastName"`
	DocumentNumber string `bson:"documentNumber" json:"documentNumber"`
	Phone          string `bson:"phone" json:"phone"`
}

type Product struct {
	Sku         string  `bson:"sku" json:"sku"`
	Name        string  `bson:"name" json:"name"`
	Description string  `bson:"description" json:"description"`
	Price       float64 `bson:"price" json:"price"`
	Quantity    int64   `bson:"quantity" json:"quantity"`
}
