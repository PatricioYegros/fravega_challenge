package orders

import (
	"challenge_pyegros/app/models"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/alicebob/miniredis/v2"
)

var (
	order = models.Order{
		ExternalReferenceID: "abc-123",
		Channel:             "Ecommerce",
		PurchaseDate:        "2024-05-01T14:00:00Z",
		TotalValue:          2000,
		Buyer: models.Buyer{
			FirstName:      "Patricio",
			LastName:       "Yegros",
			DocumentNumber: "87654321",
			Phone:          "+541112345678",
		},
		Products: []models.Product{
			{
				Sku:         "P001",
				Name:        "Producto A",
				Description: "Descripción",
				Price:       1000,
				Quantity:    2,
			},
		},
		Events: []models.Event{},
	}

	counter = models.Counter{
		ID:            "orders",
		SequenceValue: 1,
	}

	event = models.Event{
		Id:   "event-001",
		Type: "PaymentReceived",
		Date: "2024-05-01T15:00:00Z",
		User: "adminUser123",
	}

	filters = models.Filters{
		OrderId:        1,
		DocumentNumber: order.Buyer.DocumentNumber,
		Status:         "Created",
		CreatedOnFrom:  order.PurchaseDate,
		CreatedOnTo:    order.PurchaseDate,
	}
)

func CreateCacheForTesting(t *testing.T) *redis.Client {
	s := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})

	return rdb
}
func TestCreateOrderSuccess(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {

		ordersRepo := NewRepository(mt.Client, rdb)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ordersRepo.obtainID = func() (int64, error) {
			return counter.SequenceValue, nil // Devuelve el valor de counter.SequenceValue
		}

		model, err := ordersRepo.CreateOrder(order)

		response := &models.ResponseCreate{
			OrderID:   1,
			Status:    "Created",
			UpdatedOn: "2024-05-01T14:00:00Z",
		}

		assert.Nil(t, err)
		assert.Equal(t, model, response)

	})
}

func TestCreateOrderTotalMismatch(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("total mismatch", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		localOrder := order
		localOrder.TotalValue = 3000
		model, err := ordersRepo.CreateOrder(localOrder)
		assert.Nil(t, model)
		assert.Equal(t, err, ErrTotalMismatch)
	})
}

func TestCreateOrderInvalidExternalReferenceID(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("invalid external reference id", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		localOrder := order
		localOrder.ExternalReferenceID = "invalid_id"
		model, err := ordersRepo.CreateOrder(localOrder)
		assert.Nil(t, model)
		assert.Equal(t, err, ErrMismatchExternalReference)
	})
}

func TestCreateOrderObtainIDError(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("obtain id error", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		ordersRepo.obtainID = func() (int64, error) {
			return 0, ErrCreatingAutoIncrementalId
		}
		model, err := ordersRepo.CreateOrder(order)
		assert.Nil(t, model)
		assert.Equal(t, err, ErrCreatingAutoIncrementalId)
	})
}

func TestCreateOrderFailsInsertOne(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("fails insert one", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		ordersRepo.obtainID = func() (int64, error) {
			return counter.SequenceValue, nil
		}
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    11000,
			Message: "duplicate key error",
			Index:   0,
		}))
		model, err := ordersRepo.CreateOrder(order)
		assert.Nil(t, model)
		assert.NotNil(t, err)
	})
}

func TestGetOrderByIDSuccess(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		firstResponse := mtest.CreateCursorResponse(1, "orders.orders", mtest.FirstBatch, bson.D{
			{Key: "id", Value: 1},
			{Key: "externalReferenceID", Value: order.ExternalReferenceID},
			{Key: "channel", Value: order.Channel},
			{Key: "purchaseDate", Value: order.PurchaseDate},
			{Key: "totalValue", Value: order.TotalValue},
			{Key: "buyer", Value: order.Buyer},
			{Key: "products", Value: order.Products},
			{Key: "status", Value: "Created"},
			{Key: "events", Value: order.Events},
		})

		mt.AddMockResponses(firstResponse)

		model, err := ordersRepo.GetOrderByID(1)
		response := &models.ResponseGet{
			OrderID:             1,
			ExternalReferenceID: order.ExternalReferenceID,
			Channel:             order.Channel,
			ChannelTranslate:    "Comercio Electrónico",
			PurchaseDate:        order.PurchaseDate,
			TotalValue:          order.TotalValue,
			Buyer:               order.Buyer,
			Products:            order.Products,
			Status:              "Created",
			StatusTranslate:     "Creado",
			Events:              order.Events,
		}

		assert.Nil(t, err)
		assert.Equal(t, model, response)
	})
}

func TestGetOrderByIDFailsFindOne(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("fails find one", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    66,
			Message: "some error",
			Name:    "SomeError",
			Labels:  []string{},
		}))

		model, err := ordersRepo.GetOrderByID(1)
		assert.Nil(t, model)
		assert.NotNil(t, err)
	})
}

func TestGetOrderByFiltersSuccess(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		firstResponse := mtest.CreateCursorResponse(1, "orders.orders", mtest.FirstBatch, bson.D{
			{Key: "id", Value: 1},
			{Key: "externalReferenceID", Value: order.ExternalReferenceID},
			{Key: "channel", Value: order.Channel},
			{Key: "purchaseDate", Value: order.PurchaseDate},
			{Key: "totalValue", Value: order.TotalValue},
			{Key: "buyer", Value: order.Buyer},
			{Key: "products", Value: order.Products},
			{Key: "status", Value: "Created"},
			{Key: "events", Value: order.Events},
		})

		endOfCursor := mtest.CreateCursorResponse(
			0,               // Cursor ID of 0 indicates the end of the cursor
			"orders.orders", // Namespace
			mtest.NextBatch, // Flag indicating a subsequent batch (and its end)
		)

		mt.AddMockResponses(firstResponse, endOfCursor)

		model, err := ordersRepo.GetOrderByFilters(filters)
		response := order
		response.OrderID = 1
		response.Status = "Created"
		responseAll := []models.Order{response}
		assert.Nil(t, err)
		assert.Equal(t, model, responseAll)
	})
}

func TestGetOrderByFilterErrorFindOne(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("error find one", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    66,
			Message: "some error",
			Name:    "SomeError",
			Labels:  []string{},
		}))
		model, err := ordersRepo.GetOrderByFilters(filters)
		assert.Nil(t, model)
		assert.NotNil(t, err)
	})
}

func TestGetOrderByFilterErrorCursorAll(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("error cursor all", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		firstResponse := mtest.CreateCursorResponse(1, "orders.orders", mtest.FirstBatch, bson.D{
			{Key: "id", Value: 1},
			{Key: "externalReferenceID", Value: order.ExternalReferenceID},
			{Key: "channel", Value: order.Channel},
			{Key: "purchaseDate", Value: order.PurchaseDate},
			{Key: "totalValue", Value: order.TotalValue},
			{Key: "buyer", Value: order.Buyer},
			{Key: "products", Value: order.Products},
			{Key: "status", Value: "Created"},
			{Key: "events", Value: order.Events},
		})
		mt.AddMockResponses(firstResponse, mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    66,
			Message: "some error",
			Name:    "SomeError",
			Labels:  []string{},
		}))
		model, err := ordersRepo.GetOrderByFilters(filters)
		assert.Nil(t, model)
		assert.NotNil(t, err)
	})
}

func TestValidateStateTransition(t *testing.T) {

	status, err := validateStateTransition("Created", "PaymentReceived")
	assert.NoError(t, err)
	assert.Equal(t, "PaymentReceived", status)

	status, err = validateStateTransition("Created", "Canceled")
	assert.NoError(t, err)
	assert.Equal(t, "Canceled", status)

	status, err = validateStateTransition("PaymentReceived", "Invoiced")
	assert.NoError(t, err)
	assert.Equal(t, "Invoiced", status)

	status, err = validateStateTransition("Invoiced", "Returned")
	assert.NoError(t, err)
	assert.Equal(t, "Returned", status)

	status, err = validateStateTransition("Invalid", "Invalid")
	assert.Error(t, err)
}

func TestCheckUniqueEventID(t *testing.T) {
	events := []models.Event{
		{Id: "1", Date: "2022-01-01", Type: "Created"},
	}
	newEvent := models.Event{Id: "2", Date: "2022-01-02", Type: "PaymentReceived"}
	unique, err := checkUniqueEventID(events, newEvent)
	assert.True(t, unique)
	assert.NoError(t, err)

	duplicateEvent := models.Event{Id: "1", Date: "2022-01-01", Type: "Created"}
	unique, err = checkUniqueEventID(events, duplicateEvent)
	assert.False(t, unique)
	assert.NoError(t, err)

	otherEventWithSameID := models.Event{Id: "1", Date: "2022-01-02", Type: "PaymentReceived"}
	unique, err = checkUniqueEventID(events, otherEventWithSameID)
	assert.False(t, unique)
	assert.Error(t, err)
}

func TestValidateExternalReferenceIdError(t *testing.T) {
	ok, err := validateExternalReferenceId("abc-123", "Unknown")
	assert.False(t, ok)
	assert.Error(t, err)
}

func TestTranslate(t *testing.T) {
	trans := translate("CallCenter", "PaymentReceived")
	assert.Equal(t, []string{"Centro de Llamadas", "Pago Recibido"}, trans)

	trans = translate("Store", "Canceled")
	assert.Equal(t, []string{"Tienda", "Cancelado"}, trans)

	trans = translate("Affiliate", "Invoiced")
	assert.Equal(t, []string{"Afiliado", "Facturado"}, trans)

	trans = translate("Affiliate", "Returned")
	assert.Equal(t, []string{"Afiliado", "Devuelto"}, trans)

	trans = translate("Invalid", "Invalid")
	assert.Equal(t, []string{"Invalid", "Invalid"}, trans)
}

func TestUpdateEventOrderSuccess(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		firstResponse := mtest.CreateCursorResponse(1, "orders.orders", mtest.FirstBatch, bson.D{
			{Key: "id", Value: 1},
			{Key: "externalReferenceID", Value: order.ExternalReferenceID},
			{Key: "channel", Value: order.Channel},
			{Key: "purchaseDate", Value: order.PurchaseDate},
			{Key: "totalValue", Value: order.TotalValue},
			{Key: "buyer", Value: order.Buyer},
			{Key: "products", Value: order.Products},
			{Key: "status", Value: "Created"},
			{Key: "events", Value: order.Events},
		})

		mt.AddMockResponses(firstResponse)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    11000,
			Message: "duplicate key error",
			Index:   0,
		}))

		model, err := ordersRepo.UpdateEventOrder(1, event)

		assert.NotNil(t, err)
		assert.Nil(t, model)
	})
}

func TestUpdateEventOrderFailFindOne(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("fails find one", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    66,
			Message: "some error",
			Name:    "SomeError",
			Labels:  []string{},
		}))

		model, err := ordersRepo.UpdateEventOrder(1, event)
		assert.Nil(t, model)
		assert.NotNil(t, err)
	})
}

func TestUpdateEventOrderFailsUniqueEventID(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("fails unique event id", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)

		firstResponse := mtest.CreateCursorResponse(1, "orders.orders", mtest.FirstBatch, bson.D{
			{Key: "id", Value: 1},
			{Key: "externalReferenceID", Value: order.ExternalReferenceID},
			{Key: "channel", Value: order.Channel},
			{Key: "purchaseDate", Value: order.PurchaseDate},
			{Key: "totalValue", Value: order.TotalValue},
			{Key: "buyer", Value: order.Buyer},
			{Key: "products", Value: order.Products},
			{Key: "status", Value: "Created"},
			{Key: "events", Value: []models.Event{event}},
		})

		localEvent := event
		localEvent.Type = "Invoiced"

		mt.AddMockResponses(firstResponse)

		model, err := ordersRepo.UpdateEventOrder(1, localEvent)
		assert.Nil(t, model)
		assert.NotNil(t, err)
	})
}

func TestUpdateEventOrderFailsRepeatEvent(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("fails unique event id", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)

		firstResponse := mtest.CreateCursorResponse(1, "orders.orders", mtest.FirstBatch, bson.D{
			{Key: "id", Value: 1},
			{Key: "externalReferenceID", Value: order.ExternalReferenceID},
			{Key: "channel", Value: order.Channel},
			{Key: "purchaseDate", Value: order.PurchaseDate},
			{Key: "totalValue", Value: order.TotalValue},
			{Key: "buyer", Value: order.Buyer},
			{Key: "products", Value: order.Products},
			{Key: "status", Value: "Created"},
			{Key: "events", Value: []models.Event{event}},
		})

		mt.AddMockResponses(firstResponse)

		model, err := ordersRepo.UpdateEventOrder(1, event)
		assert.NotNil(t, model)
		assert.Nil(t, err)
	})
}

func TestUpdateEventOrderFailsEventTransition(t *testing.T) {
	rdb := CreateCacheForTesting(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("fails unique event id", func(mt *mtest.T) {
		ordersRepo := NewRepository(mt.Client, rdb)

		firstResponse := mtest.CreateCursorResponse(1, "orders.orders", mtest.FirstBatch, bson.D{
			{Key: "id", Value: 1},
			{Key: "externalReferenceID", Value: order.ExternalReferenceID},
			{Key: "channel", Value: order.Channel},
			{Key: "purchaseDate", Value: order.PurchaseDate},
			{Key: "totalValue", Value: order.TotalValue},
			{Key: "buyer", Value: order.Buyer},
			{Key: "products", Value: order.Products},
			{Key: "status", Value: "Created"},
			{Key: "events", Value: []models.Event{event}},
		})

		localEvent := models.Event{
			Id:   "event-002",
			Type: "Invoiced",
			Date: "2025-05-01T15:00:00Z",
			User: "admin002",
		}

		mt.AddMockResponses(firstResponse)

		model, err := ordersRepo.UpdateEventOrder(1, localEvent)
		assert.Nil(t, model)
		assert.NotNil(t, err)
	})
}
