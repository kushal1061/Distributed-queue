package models

import "time"

// Task represents a unit of work
type Task struct {
	ID             string                 `bson:"_id,omitempty" json:"id"` //store uuid
	Type           string                 `bson:"type" json:"type"`
	Payload        map[string]interface{} `bson:"payload" json:"payload"` // JSON-like payload - This allows flexibility in handling different payload structures.
	Status         string                 `bson:"status" json:"status"`   // queueud, in progress, success, failed
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
	CompletionTime time.Time              `bson:"completion_time,omitempty" json:"completion_time,omitempty"`
	UpdatedAt      time.Time              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	TimeToComplete int                    `bson:"time_to_complete" json:"timeToComplete"` // Time in seconds
}
