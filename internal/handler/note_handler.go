package handler

import (
	"net/http"
	"strconv"

	"github.com/Napat/mcpserver-demo/internal/service"
	"github.com/Napat/mcpserver-demo/models"
	"github.com/Napat/mcpserver-demo/pkg/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateNoteRequest is a data structure for creating a note
type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// UpdateNoteRequest is a data structure for updating a note
type UpdateNoteRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// NoteHandler handles note operations
type NoteHandler struct {
	noteService service.INoteService
	logger      *zap.Logger
}

// NewNoteHandler creates a new instance of NoteHandler
func NewNoteHandler(noteService service.INoteService, logger *zap.Logger) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
		logger:      logger,
	}
}

// GetAllNotes retrieves all notes for a user
func (h *NoteHandler) GetAllNotes(c echo.Context) error {
	userID := middleware.GetUserIDFromToken(c)

	notes, err := h.noteService.GetAllByUserID(userID)
	if err != nil {
		h.logger.Error("Failed to get notes", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get notes")
	}

	return c.JSON(http.StatusOK, notes)
}

// GetNote retrieves a note by ID
func (h *NoteHandler) GetNote(c echo.Context) error {
	userID := middleware.GetUserIDFromToken(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid note ID")
	}

	note, err := h.noteService.GetByID(uint(noteID), userID)
	if err != nil {
		if err.Error() == "unauthorized access to note" {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Note not found")
		}
		h.logger.Error("Failed to get note", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get note")
	}

	return c.JSON(http.StatusOK, note)
}

// CreateNote creates a new note
func (h *NoteHandler) CreateNote(c echo.Context) error {
	userID := middleware.GetUserIDFromToken(c)

	req := new(CreateNoteRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	note := models.Note{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := h.noteService.Create(&note); err != nil {
		h.logger.Error("Failed to create note", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create note")
	}

	return c.JSON(http.StatusCreated, note)
}

// UpdateNote updates a note
func (h *NoteHandler) UpdateNote(c echo.Context) error {
	userID := middleware.GetUserIDFromToken(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid note ID")
	}

	req := new(UpdateNoteRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	note := &models.Note{
		ID:      uint(noteID),
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	err = h.noteService.Update(note, userID)
	if err != nil {
		if err.Error() == "unauthorized access to note" {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Note not found")
		}
		h.logger.Error("Failed to update note", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update note")
	}

	return c.JSON(http.StatusOK, note)
}

// DeleteNote deletes a note
func (h *NoteHandler) DeleteNote(c echo.Context) error {
	userID := middleware.GetUserIDFromToken(c)
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid note ID")
	}

	err = h.noteService.Delete(uint(noteID), userID)
	if err != nil {
		if err.Error() == "unauthorized access to note" {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Note not found")
		}
		h.logger.Error("Failed to delete note", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete note")
	}

	return c.NoContent(http.StatusNoContent)
}
