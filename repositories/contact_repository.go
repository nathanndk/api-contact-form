package repositories

import (
	"api-contact-form/models"

	"gorm.io/gorm"
)

/*
Package repositories implements data persistence and retrieval for contact entities.

This file provides a GORM-backed ContactRepository implementation that performs
CRUD operations for the Contact model. The implementation intentionally relies
on GORM's soft-delete behavior (gorm.DeletedAt in the model) so that:
  - calling Delete(...) performs a soft delete (sets deleted_at timestamp);
  - normal queries (Find, First) automatically exclude soft-deleted rows.

See GORM documentation for Delete / Soft Delete behavior.
*/

// ContactRepository defines the interface for contact data operations.
type ContactRepository interface {
	// Create inserts a new contact record into the database.
	Create(contact *models.Contact) error

	// FindAll retrieves all non-deleted contacts.
	// Note: GORM automatically excludes soft-deleted rows when the model
	// uses gorm.DeletedAt.
	FindAll() ([]models.Contact, error)

	// FindByID retrieves a contact by primary key (ID). Soft-deleted records
	// are excluded by default.
	FindByID(id uint) (*models.Contact, error)

	// Update persists changes to an existing contact.
	Update(contact *models.Contact) error

	// Delete performs a soft-delete for the provided contact (sets deleted_at).
	// For a hard delete, callers can use db.Unscoped().Delete(...) directly.
	Delete(contact *models.Contact) error
}

// contactRepository is a GORM-based implementation of ContactRepository.
type contactRepository struct {
	db *gorm.DB
}

// NewContactRepository constructs a new ContactRepository backed by the provided GORM DB.
func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

// Create inserts a new contact into the database using GORM.
//
// On success, the contact struct will have its ID and timestamps populated by GORM.
func (r *contactRepository) Create(contact *models.Contact) error {
	return r.db.Create(contact).Error
}

// FindAll returns all contacts that are not soft-deleted.
//
// This relies on GORM's global soft-delete scope (models with gorm.DeletedAt
// are excluded automatically from normal queries).
func (r *contactRepository) FindAll() ([]models.Contact, error) {
	var contacts []models.Contact
	if err := r.db.Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

// FindByID looks up a contact by primary key and returns it.
//
// If no record is found, GORM will return an error (e.g., gorm.ErrRecordNotFound).
// Soft-deleted records are excluded by default; use r.db.Unscoped().First(...) if you
// intentionally need deleted records.
func (r *contactRepository) FindByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update persists changes to an existing contact record.
//
// This uses Save(...) which performs an update based on the primary key.
func (r *contactRepository) Update(contact *models.Contact) error {
	return r.db.Save(contact).Error
}

// Delete performs a soft delete using GORM's Delete(...) method.
//
// GORM will set the model's DeletedAt timestamp rather than physically removing
// the row. To permanently remove rows, use Unscoped().Delete(...).
func (r *contactRepository) Delete(contact *models.Contact) error {
	return r.db.Delete(contact).Error
}
