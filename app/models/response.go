package models

type ResponseCreate struct {
	OrderID   int64  `json:"orderID"`
	Status    string `json:"status"`
	UpdatedOn string `json:"updatedOn"`
}

type ResponseUpdate struct {
	OrderID        int64  `json:"orderID"`
	PreviousStatus string `json:"previousStatus"`
	NewStatus      string `json:"newStatus"`
	UpdatedOn      string `json:"updatedOn"`
}

type ResponseGet struct {
	OrderID             int64     `json:"orderID"`
	ExternalReferenceID string    `json:"externalReferenceID"`
	Channel             string    `json:"channel"`
	ChannelTranslate    string    `json:"channelTranslate"`
	PurchaseDate        string    `json:"purchaseDate"`
	TotalValue          float64   `json:"totalValue"`
	Buyer               Buyer     `json:"buyer"`
	Products            []Product `json:"product"`
	Status              string    `json:"status"`
	StatusTranslate     string    `json:"statusTranslate"`
	Events              []Event   `json:"events"`
}
