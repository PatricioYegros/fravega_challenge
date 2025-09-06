package models

type Filters struct {
	OrderId        int64  `json:"orderId"`
	DocumentNumber string `json:"documentNumber"`
	Status         string `json:"status"`
	CreatedOnFrom  string `json:"createdOnFrom"`
	CreatedOnTo    string `json:"createdOnTo"`
}
