package model

import "time"

type BaseModel struct {
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy *string    `db:"created_by" json:"created_by,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	DeletedBy *string    `db:"deleted_by" json:"deleted_by,omitempty"`
}
