package repoimpl

import (
	"github.com/neathhatake/the_Scan/internal/models"
	"github.com/neathhatake/the_Scan/internal/repositories"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
	"gorm.io/gorm"
)



type documentsRepo struct {
	db *gorm.DB


}

func NewDocumentRepository(db *gorm.DB) repositories.DocumentsRepository {
	return &documentsRepo{
		db: db,
	}
}


func (r *documentsRepo) CreateDoc(doc *models.ScanDocument) error{
	if err := r.db.Create(doc).Error; err != nil {
		return apperrors.Internal("failed to create document...",err)
	}
	return nil 
}



func (r *documentsRepo) FindByID(id uint) (*models.ScanDocument, error){
	var doc models.ScanDocument	
	if err := r.db.First(&doc, id).Error; err != nil {
		return nil , apperrors.NotFound("failed to find document...")
	}
	return &doc , nil
}



func (r *documentsRepo) FindByIDandUser(id, userID uint) (*models.ScanDocument, error){
	var doc models.ScanDocument
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&doc).Error
	if err != nil {
		return nil, apperrors.NotFound("document not found")
	}
	return &doc, nil

}


func (r *documentsRepo)  ListByUser(userID uint, offset, limit int) ([]models.ScanDocument, int64, error){
	var docs []models.ScanDocument
	var total int64
	
	r.db.Model(&models.ScanDocument{}).Where("user_id = ?", userID).Count(&total)

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&docs).Error
	if err != nil {
		return nil, 0, apperrors.Internal("failed to list documents", err)
	}
	return docs, total, nil
}





func (r *documentsRepo) UpdateDoc(doc *models.ScanDocument, feild map[string]any ) error{
	if err := r.db.Model(doc).Updates(feild).Error; err != nil {
		return apperrors.Internal("faild to update document...",err)
	}
	return nil
}


func (r *documentsRepo) DeleteDoc(doc *models.ScanDocument) error{
	if err := r.db.Delete(doc).Error; err != nil {
		return apperrors.Internal("failed to delete document...",err)
	}
	return nil
}






