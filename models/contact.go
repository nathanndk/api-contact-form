// Package models defines the data models for the API Contact Form application.
//
// Contact is the record representation of a contact message stored in the database.
// This model uses idiomatic, DB-agnostic GORM types/tags so it works with both
// PostgreSQL and MySQL. Use gorm.DeletedAt to enable GORM's soft-delete behavior.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Contact represents a contact message submitted through the API.
//
// Notes:
//  - Keep tags DB-agnostic: let GORM choose the exact SQL type for time.Time
//    so the same code works across MySQL/MariaDB and PostgreSQL.
//  - Use gorm.DeletedAt for soft deletes; GORM will set deleted_at and exclude
//    soft-deleted rows from normal queries automatically.
type Contact struct {
	// ID is the primary key. GORM treats a field named ID as the primary key
	// by convention; explicit `primaryKey` keeps intent clear.
	ID uint `gorm:"primaryKey;column:id" json:"id"`

	// FullName is the name of the person submitting the contact message.
	// Keep length constraints here so migrations create appropriate columns.
	FullName string `gorm:"column:full_name;type:VARCHAR(100);not null" json:"full_name"`

	// Email is the email address of the submitter.
	// Consider adding a unique index at the DB level if you want to enforce uniqueness.
	Email string `gorm:"column:email_address;type:VARCHAR(100);not null" json:"email"`

	// Phone is the phone number.
	Phone string `gorm:"column:phone_number;type:VARCHAR(20);not null" json:"phone"`

	// Message stores the contact message content.
	Message string `gorm:"column:message_text;type:TEXT;not null" json:"message"`

	// CreatedAt / UpdatedAt are automatically maintained by GORM.
	// Do NOT hardcode a DB-specific type like DATETIME — let GORM map time.Time
	// to the appropriate type (TIMESTAMP/TIMESTAMPTZ for Postgres, DATETIME for MySQL).
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// DeletedAt enables GORM soft deletes. Use gorm.DeletedAt instead of time.Time
	// so GORM can handle NULL vs timestamp semantics correctly across DB engines.
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName overrides the default table name that GORM derives from the struct.
func (Contact) TableName() string {
	return "contact_messages"
}
