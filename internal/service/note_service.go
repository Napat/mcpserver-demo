package service

import (
	"errors"

	"github.com/Napat/mcpserver-demo/internal/repository"
	"github.com/Napat/mcpserver-demo/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=./note_service.go -destination=./mocks/mock_note_service.go -package=mocks

// INoteService interface for managing note business logic
type INoteService interface {
	Create(note *models.Note) error
	GetByID(id, userID uint) (*models.Note, error)
	GetAllByUserID(userID uint) ([]models.Note, error)
	Update(note *models.Note, userID uint) error
	Delete(id, userID uint) error
}

// NoteService struct for handling note business logic
type NoteService struct {
	noteRepo repository.INoteRepository
	logger   *zap.Logger
}

// NewNoteService creates a new instance of NoteService
func NewNoteService(noteRepo repository.INoteRepository, logger *zap.Logger) INoteService {
	return &NoteService{
		noteRepo: noteRepo,
		logger:   logger,
	}
}

// Create creates a new note
func (s *NoteService) Create(note *models.Note) error {
	return s.noteRepo.Create(note)
}

// GetByID retrieves a note by ID and checks access permissions
func (s *NoteService) GetByID(id, userID uint) (*models.Note, error) {
	note, err := s.noteRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if note.UserID != userID {
		return nil, errors.New("unauthorized access to note")
	}

	return note, nil
}

// GetAllByUserID retrieves all notes for a user
func (s *NoteService) GetAllByUserID(userID uint) ([]models.Note, error) {
	return s.noteRepo.FindByUserID(userID)
}

// Update updates a note and checks access permissions
func (s *NoteService) Update(note *models.Note, userID uint) error {
	existing, err := s.noteRepo.FindByID(note.ID)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return errors.New("unauthorized access to note")
	}

	return s.noteRepo.Update(note)
}

// Delete removes a note and checks access permissions
func (s *NoteService) Delete(id, userID uint) error {
	existing, err := s.noteRepo.FindByID(id)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return errors.New("unauthorized access to note")
	}

	return s.noteRepo.Delete(id)
}
