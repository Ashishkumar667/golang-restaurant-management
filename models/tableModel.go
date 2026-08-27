package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID					primitive.ObjectID         `bson:"_id"`
	Number_of_guest		*int 					   `json:"number_of_guest" validate:"required"`
	Table_Number		*int					   `json:"table_number" validate:"required"`
	Created_At			time.Time                  `json:"created_at"`
	Updated_At			time.Time				   `json:"updated_at"`
	Table_id			string                     `json:"table_id"`
}