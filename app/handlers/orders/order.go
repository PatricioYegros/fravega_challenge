package orders

import (
	"challenge_pyegros/app/models"
	"challenge_pyegros/app/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

// OrderHandler holds dependencies like the MongoDB client.
type OrderHandler struct {
	mongoClient *mongo.Client
}

// NewOrderHandler creates a new OrderHandler instance.
func NewOrderHandler(client *mongo.Client) *OrderHandler {
	return &OrderHandler{mongoClient: client}
}

// CreateOrder creates a new order in MongoDB.
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create order")
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var order models.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}
	err = utils.CheckFormatDate(order.PurchaseDate)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
		return
	}
	var response *models.ResponseCreate
	response, err = h.u.CreateOrder(order)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	// Logic to parse request body and insert a new order into MongoDB
	w.Write(json)
}

// UpdateEventOrder updates the state of the order considering the Event
func (h *Handler) UpdateEventOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var event models.Event
	orderID := chi.URLParam(r, "orderId")
	orderIDInt, err := strconv.Atoi(orderID)
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	err = utils.CheckFormatDate(event.Date)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
		return
	}

	var response *models.ResponseUpdate
	response, err = h.u.UpdateEventOrder(int64(orderIDInt), event)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	// Logic to parse request body and insert a new order into MongoDB
	w.Write(json)
}

// GetOrderByID retrieves a order by ID from MongoDB.
func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderID := chi.URLParam(r, "orderId")
	orderIDInt, err := strconv.Atoi(orderID)

	response, err := h.u.GetOrderByID(int64(orderIDInt))
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	// Logic to parse request body and insert a new order into MongoDB
	w.Write(json)
}

// GetOrderByFilters retrieves a order by Filters from MongoDB
func (h *Handler) GetOrderByFilters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	filters := utils.GetFilters(r)

	err := utils.CheckFormatDate(filters.CreatedOnFrom)
	if err != nil {
		http.Error(w, "The From date is not in the correct format", http.StatusInternalServerError)
		return
	}
	err = utils.CheckFormatDate(filters.CreatedOnTo)
	if err != nil {
		http.Error(w, "The To date is not in the correct format", http.StatusInternalServerError)
		return
	}

	response, err := h.u.GetOrderByFilters(filters)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	w.Write(json)
}
