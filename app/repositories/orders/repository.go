package orders

import (
	"challenge_pyegros/app/models"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrTotalMismatch             = errors.New("Total value does not match sum of products")
	ErrParsingID                 = errors.New("Error parsing order ID")
	ErrMismatchExternalReference = errors.New("External ReferenceId does not match with channel")
	ErrChannelNotFound           = errors.New("Channel not found")
	ErrCreatingAutoIncrementalId = errors.New("Error creating auto incremental ID")
	ErrGettingAutoIncrementalId  = errors.New("Error getting auto incremental ID")
	ErrUpdatingAutoIncrementalId = errors.New("Error updating auto incremental ID")
	ErrAnotherEventWithSameID    = errors.New("Another event with same ID already exists")
)

type Repository struct {
	db       *mongo.Client
	obtainID func() (int64, error)
}

func NewRepository(client *mongo.Client) *Repository {
	repo := &Repository{db: client}
	repo.obtainID = repo.defaultObtainID
	return repo
}

func (r *Repository) CreateOrder(order models.Order) (*models.ResponseCreate, error) {
	if !validateTotal(order.Products, order.TotalValue) {
		return nil, ErrTotalMismatch
	}

	isValid, err := validateExternalReferenceId(order.ExternalReferenceID, order.Channel)
	if err != nil || !isValid {
		return nil, err
	}

	var id int64
	id, err = r.obtainID()
	if err != nil {
		return nil, err
	}
	order.OrderID = id
	order.Status = "Created"
	order.Events = []models.Event{}

	collection := r.db.Database("orders").Collection("orders")
	_, err = collection.InsertOne(context.TODO(), order)
	if err != nil {
		return nil, err
	}

	response := &models.ResponseCreate{
		OrderID:   id,
		Status:    "Created",
		UpdatedOn: order.PurchaseDate,
	}
	return response, nil
}

func (r *Repository) UpdateEventOrder(orderID int64, event models.Event) (*models.ResponseUpdate, error) {
	collection := r.db.Database("orders").Collection("orders")

	filter := bson.M{"id": orderID}

	var order models.Order
	err := collection.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	response := &models.ResponseUpdate{
		OrderID:        order.OrderID,
		PreviousStatus: order.Status,
		NewStatus:      "newStatus",
		UpdatedOn:      event.Date,
	}

	unique, errEvent := checkUniqueEventID(order.Events, event)
	if errEvent != nil {
		return nil, errEvent
	} else {
		if !unique {
			return response, nil
		}
	}

	newStatus, err := validateStateTransition(order.Status, event.Type)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": bson.M{"status": newStatus}, "$push": bson.M{"events": event}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *Repository) GetOrderByID(orderID int64) (*models.ResponseGet, error) {
	collection := r.db.Database("orders").Collection("orders")

	filter := bson.M{"id": orderID}

	var order *models.Order
	err := collection.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	translations := translate(order.Channel, order.Status)

	response := &models.ResponseGet{
		OrderID:             order.OrderID,
		ExternalReferenceID: order.ExternalReferenceID,
		Channel:             order.Channel,
		ChannelTranslate:    translations[0],
		PurchaseDate:        order.PurchaseDate,
		TotalValue:          order.TotalValue,
		Buyer:               order.Buyer,
		Products:            order.Products,
		Status:              order.Status,
		StatusTranslate:     translations[1],
		Events:              order.Events,
	}

	return response, nil
}

func (r *Repository) GetOrderByFilters(filters models.Filters) ([]models.Order, error) {
	collection := r.db.Database("orders").Collection("orders")

	query := ApplyFilters(filters)

	projection := bson.M{
		"events": bson.M{
			"$slice": -1,
		},
	}

	filtersQuery := bson.M{}
	if len(query) > 0 {
		filtersQuery = bson.M{"$and": query}
	}

	cursor, err := collection.Find(context.TODO(), filtersQuery, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	var orders = []models.Order{}
	if err = cursor.All(context.TODO(), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func validateTotal(products []models.Product, total float64) bool {
	var totalCalculated float64

	for _, product := range products {
		totalCalculated = (product.Price * float64(product.Quantity)) + totalCalculated
	}

	if totalCalculated != total {
		return false
	}

	return true
}

func validateStateTransition(actualStatus string, typeEvent string) (string, error) {
	err := errors.New("Invalid state transition")

	switch typeEvent {
	case "PaymentReceived":
		if actualStatus == "Created" {
			return "PaymentReceived", nil
		}
	case "Canceled":
		if actualStatus == "Created" {
			return "Canceled", nil
		}
	case "Invoiced":
		if actualStatus == "PaymentReceived" {
			return "Invoiced", nil
		}
	case "Returned":
		if actualStatus == "Invoiced" {
			return "Returned", nil
		}
	}

	return "", err
}

func checkUniqueEventID(events []models.Event, newEvent models.Event) (bool, error) {
	for _, event := range events {
		if event.Id == newEvent.Id {
			if event.Date == newEvent.Date && event.Type == newEvent.Type {
				return false, nil
			} else {
				return false, ErrAnotherEventWithSameID
			}
		}
	}
	return true, nil
}

func validateExternalReferenceId(id string, channel string) (bool, error) {
	values := map[string]string{
		"Ecommerce":  "abc-123",
		"CallCenter": "def-456",
		"Store":      "ghi-789",
		"Affiliate":  "jkl-012",
	}

	v, ok := values[channel]
	if ok {
		if v == id {
			return true, nil
		} else {
			return false, ErrMismatchExternalReference
		}
	} else {
		return false, ErrChannelNotFound
	}
}

func (r *Repository) defaultObtainID() (int64, error) {
	filter := bson.M{"_id": "orders"}
	var counter models.Counter
	collection := r.db.Database("orders").Collection("counters")
	err := collection.FindOne(context.TODO(), filter).Decode(&counter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newCounter := models.Counter{
				ID:            "orders",
				SequenceValue: 1,
			}
			_, err := collection.InsertOne(context.TODO(), newCounter)
			if err != nil {
				return 0, ErrCreatingAutoIncrementalId
			}
			return 1, nil
		} else {
			fmt.Println(err.Error())
			return 0, ErrGettingAutoIncrementalId
		}
	} else {
		update := bson.M{"$inc": bson.M{"sequence_value": 1}}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return 0, ErrUpdatingAutoIncrementalId
		}
		return counter.SequenceValue + 1, nil
	}
}

func translate(channel string, status string) []string {
	var translations = []string{}

	switch channel {
	case "Ecommerce":
		translations = append(translations, "Comercio Electrónico")
	case "CallCenter":
		translations = append(translations, "Centro de Llamadas")
	case "Store":
		translations = append(translations, "Tienda")
	case "Affiliate":
		translations = append(translations, "Afiliado")
	default:
		translations = append(translations, channel)
	}

	switch status {
	case "Created":
		translations = append(translations, "Creado")
	case "PaymentReceived":
		translations = append(translations, "Pago Recibido")
	case "Canceled":
		translations = append(translations, "Cancelado")
	case "Invoiced":
		translations = append(translations, "Facturado")
	case "Returned":
		translations = append(translations, "Devuelto")
	default:
		translations = append(translations, status)
	}

	return translations
}
