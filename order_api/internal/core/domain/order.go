package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrProductNotFound = errors.New("product not found")

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID string  `bson:"product_id" json:"product_id"`
	Quantity  int     `bson:"quantity"   json:"quantity"`
	UnitPrice float64 `bson:"unit_price" json:"unit_price"`
}

type Order struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"    json:"id"`
	CustomerID     string             `bson:"customer_id"      json:"customer_id"`
	Items          []OrderItem        `bson:"items"            json:"items"`
	TotalAmount    float64            `bson:"total_amount"     json:"total_amount"`
	Status         OrderStatus        `bson:"status"           json:"status"`
	IdempotencyKey string             `bson:"idempotency_key"  json:"-"`
	CreatedAt      time.Time          `bson:"created_at"       json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"       json:"updated_at"`
}

type Product struct {
	ID    string  `bson:"_id"`
	Price float64 `bson:"price"`
}
