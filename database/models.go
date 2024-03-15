package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sitting struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

type User struct {
	// ID       primitive.ObjectID `json:"id" bson:"_id"`
	FullName string    `json:"name" bson:"name"`
	Age      uint      `json:"age" bson:"age"`
	Sittings []Sitting `json:"sittings" bson:"sittings"`
}

type UserwID struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	FullName string             `json:"name" bson:"name"`
	Age      uint               `json:"age" bson:"age"`
	Sittings []Sitting          `json:"sittings" bson:"sittings"`
}

type AdminCredentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
