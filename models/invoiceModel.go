package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID               primitive.ObjectID `bson:"_id"`
	Invoice_id       string             `json:"invoice_id"`
	Order_id         string             `json:"order_id" validate:"required"`
	Payment_method   *string            `json:"payment_method"`
	Payment_status   *string            `json:"payment_status"`
	Payment_due_date time.Time          `json:"payment_due_date"`
	UpdatedAt        time.Time          `json:"updated_at"`
	CreatedAt        time.Time          `json:"created_at"`
}
