package utils

import (
	"challenge_pyegros/app/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetFilters(r *http.Request) models.Filters {
	orderId, _ := getQueryValue(r, "orderId")
	orderIdInt, _ := strconv.Atoi(orderId)
	documentNumber, _ := getQueryValue(r, "documentNumber")
	status, _ := getQueryValue(r, "status")
	createdOnFrom, _ := getQueryValue(r, "createdOnFrom")
	createdOnTo, _ := getQueryValue(r, "createdOnTo")

	filters := models.Filters{
		OrderId:        int64(orderIdInt),
		DocumentNumber: documentNumber,
		Status:         status,
		CreatedOnFrom:  createdOnFrom,
		CreatedOnTo:    createdOnTo,
	}

	return filters
}

func getQueryValue(r *http.Request, key string) (string, error) {
	if !r.URL.Query().Has(key) {
		return "", fmt.Errorf("missing query parameter: %s", key)
	}

	return r.URL.Query().Get(key), nil
}

func CheckFormatDate(date string) error {
	var err = errors.New("The date is not in the correct format")
	_, errParse := time.Parse(time.RFC3339, date)
	if errParse != nil {
		return err
	}
	return nil
}
