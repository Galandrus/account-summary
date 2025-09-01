package models

import "time"

type Account struct {
	ID        string         `json:"id" bson:"id"`
	Email     string         `json:"email" bson:"email"`
	Summary   AccountSummary `json:"summary" bson:"summary"`
	CreatedAt time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" bson:"updatedAt"`
}
