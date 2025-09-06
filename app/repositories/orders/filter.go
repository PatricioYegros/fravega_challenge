package orders

import (
	"challenge_pyegros/app/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ApplyOrderIdFilter(filters *models.Filters, query []bson.M) []bson.M {
	if filters.OrderId > 0 {
		query = append(query, bson.M{"id": filters.OrderId})
	}
	return query
}

func ApplyDocumentNumberFilter(filters *models.Filters, query []bson.M) []bson.M {
	if len(filters.DocumentNumber) > 0 {
		query = append(query, bson.M{"buyer.documentNumber": filters.DocumentNumber})
	}
	return query
}

func ApplyStatusFilter(filters *models.Filters, query []bson.M) []bson.M {
	if len(filters.Status) > 0 {
		query = append(query, bson.M{"status": filters.Status})
	}
	return query
}

func ApplyCreatedOnFilter(filters *models.Filters, query []bson.M) []bson.M {
	if len(filters.CreatedOnFrom) > 0 && len(filters.CreatedOnTo) > 0 {
		query = append(query, bson.M{
			"purchaseDate": bson.M{
				"$gte": filters.CreatedOnFrom,
				"$lte": filters.CreatedOnTo,
			},
		})
	}
	return query
}

func ApplyFilters(filters models.Filters) []bson.M {
	var query []bson.M
	query = ApplyOrderIdFilter(&filters, query)
	query = ApplyDocumentNumberFilter(&filters, query)
	query = ApplyStatusFilter(&filters, query)
	query = ApplyCreatedOnFilter(&filters, query)
	return query
}
