package models

import (
	"time"

	"gorm.io/gorm"
)

// Note is a model for storing notes
type Note struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null;index:idx_notes_title" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `gorm:"not null;index:idx_notes_user_id" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;index:idx_notes_created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;index:idx_notes_updated_at" json:"updated_at"`
}

// TableName defines the table name
func (Note) TableName() string {
	return "notes"
}

// BeforeCreate runs before creating data
func (n *Note) BeforeCreate(tx *gorm.DB) error {
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate runs before updating data
func (n *Note) BeforeUpdate(tx *gorm.DB) error {
	n.UpdatedAt = time.Now()
	return nil
}
