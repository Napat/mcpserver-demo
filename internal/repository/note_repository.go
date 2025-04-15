package repository

import (
	"errors"

	"github.com/Napat/mcpserver-demo/models"
	"gorm.io/gorm"
)

//go:generate mockgen -source=./note_repository.go -destination=./mocks/mock_note_repository.go -package=mocks

// INoteRepository is an interface for managing note data in the database
type INoteRepository interface {
	Create(note *models.Note) error
	FindByID(id uint) (*models.Note, error)
	FindByUserID(userID uint) ([]models.Note, error)
	Update(note *models.Note) error
	Delete(id uint) error
}

// NoteRepository is a struct that implements INoteRepository
type NoteRepository struct {
	db *gorm.DB
}

// NewNoteRepository creates a new instance of NoteRepository
func NewNoteRepository(db *gorm.DB) INoteRepository {
	return &NoteRepository{
		db: db,
	}
}

// Create adds a new note to the database
func (r *NoteRepository) Create(note *models.Note) error {
	return r.db.Create(note).Error
}

// FindByID finds a note by ID
func (r *NoteRepository) FindByID(id uint) (*models.Note, error) {
	var note models.Note
	result := r.db.First(&note, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("note not found")
		}
		return nil, result.Error
	}
	return &note, nil
}

// FindByUserID finds all notes for a user
func (r *NoteRepository) FindByUserID(userID uint) ([]models.Note, error) {
	var notes []models.Note
	result := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notes)

	if result.Error != nil {
		return nil, result.Error
	}
	return notes, nil
}

// Update updates a note
func (r *NoteRepository) Update(note *models.Note) error {
	return r.db.Model(note).
		Select("title", "content", "updated_at").
		Updates(note).Error
}

// Delete removes a note
func (r *NoteRepository) Delete(id uint) error {
	return r.db.Delete(&models.Note{}, id).Error
}
