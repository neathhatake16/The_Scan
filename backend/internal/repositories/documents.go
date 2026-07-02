package repositories

import "github.com/neathhatake/the_Scan/internal/models"



type DocumentsRepository interface {
	CreateDoc(doc *models.ScanDocument) error
	FindByIDandUser(id ,userID uint) (*models.ScanDocument, error)
	ListByUser(userID uint , offset , limit int) ([]models.ScanDocument,int64, error) 
	FindByID(id uint) (*models.ScanDocument, error)
	UpdateDoc(doc *models.ScanDocument, fields map[string]any) error
	DeleteDoc(doc *models.ScanDocument) error
	
}



