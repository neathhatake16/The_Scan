package models

import "time"

type ScanDocument struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"-"`
	Title            string    `gorm:"size:255;not null;default:'Untitled Scan'" json:"title"`
	OriginalFilename string    `gorm:"size:255" json:"original_filename"`
	PDFPath          string    `gorm:"size:500;not null" json:"-"`
	PDFSizeBytes     int64     `gorm:"default:0" json:"pdf_size_bytes"`
	PageCount        int       `gorm:"default:1" json:"page_count"`
	Status           string    `gorm:"size:20;default:'done'" json:"status"`
	CreatedAt        time.Time `gorm:"index" json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

//DTO

type ResponeDoc struct {
	ID               uint      `json:"id"`
	Title            string    `json:"title"`
	OriginalFilename string    `json:"original_filename"`
	PDFSizeBytes     int64     `json:"pdf_size_bytes"`
	PageCount        int       `json:"page_count"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type ReNameDocRequest struct {
	Title string `json:"title" binding:"required,min=1,max=255"`
}

func ToDocumentResponse(d *ScanDocument) ResponeDoc {
	return ResponeDoc{
		ID:               d.ID,
		Title:            d.Title,
		OriginalFilename: d.OriginalFilename,
		PDFSizeBytes:     d.PDFSizeBytes,
		PageCount:        d.PageCount,
		Status:           d.Status,
		CreatedAt:        d.CreatedAt,
	}
}
