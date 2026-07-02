package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/neathhatake/the_Scan/internal/config"
	"github.com/neathhatake/the_Scan/internal/models"

	"github.com/neathhatake/the_Scan/internal/repositories"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
)



type DocumentService struct {
	repo repositories.DocumentsRepository	
	user repositories.UserRepository
	cfg *config.Config
}



// ScanAndSave forwards image to the Python scanner, stores the PDF,
// and persists the document record + updates storage stats.

func NewDocumentService(repo repositories.DocumentsRepository, user repositories.UserRepository, cfg *config.Config) *DocumentService {
	return &DocumentService{
		repo: repo,
		user: user,
		cfg: cfg,
	}
}

type  Scan struct {
	UserID uint 
	FileName string
	ImageData []byte
	Title string 
}

func (s *DocumentService) ScanandSave(input *Scan ) (*models.ResponeDoc, error){
	pdfByte , err := s.callScan(input.FileName, input.ImageData)
	if err != nil {
		return nil ,err
	}

	if err := os.MkdirAll(s.cfg.PDFStorageDir, 0o755); err != nil {
		return nil , apperrors.Internal("failed to create storage directory", err)
	}
	pdfName := fmt.Sprintf("user_%d_%d.pdf", input.UserID, time.Now().UnixMilli())
	pdfPath := filepath.Join(s.cfg.PDFStorageDir, pdfName)

	if err := os.WriteFile(pdfPath,pdfByte, 0o644); err != nil {
		return nil , apperrors.Internal("failed to write PDF", err)
	}
	doc := &models.ScanDocument{
		UserID:           input.UserID,
		Title:            input.Title,
		OriginalFilename: input.FileName,
		PDFPath:          pdfPath,
		PDFSizeBytes:     int64(len(pdfByte)),
		PageCount:        1,
		Status:           "done",
	}
	if err := s.repo.CreateDoc(doc); err != nil {
		_ = os.Remove(pdfPath)
		return nil , err
	}
	_ = s.user.IncrementStorage(input.UserID, doc.PDFSizeBytes)
	res := models.ToDocumentResponse(doc)
	return &res , nil 
}

func(s *DocumentService) DeleteDoc(id uint , userID uint) error {
	doc , err := s.repo.FindByIDandUser(id,userID)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteDoc(doc); err != nil {
		return err
	}

	_ = os.Remove(doc.PDFPath)
	_ = s.user.DecrementStorage(userID, doc.PDFSizeBytes)
	return nil

}


func (s *DocumentService) RenameDoc(id  , userID uint , title string) (*models.ResponeDoc,  error) {
	doc , err := s.repo.FindByIDandUser(id,userID)
	if err != nil {
		return nil ,err
	}
	if err := s.repo.UpdateDoc(doc, map[string]any{"title": title}); err != nil {
		return nil, err
	}
	doc.Title = title
	res := models.ToDocumentResponse(doc)
	return &res , nil
}


func (s *DocumentService) GetDocPath(id , userID uint) (string ,string, error){
	doc , err := s.repo.FindByIDandUser(id,userID)
	if err != nil {
		return "" ,"", err
	}
	if _ , err := os.Stat(doc.PDFPath); os.IsNotExist(err){
		return "" ,"", apperrors.NotFound("document not found on server!")
	}
	return doc.PDFPath, doc.Title , nil
}



func (s *DocumentService) callScan(filename string , imageData []byte) ([]byte , error){
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw , err := w.CreateFormFile("file",filename)
	if err != nil {
		return nil , apperrors.Internal("failed to build multipart request", err)
	}

	if _ , err := fw.Write(imageData); err != nil {
		return nil , apperrors.Internal("failed to write Image request", err)
	}
	w.Close()

	resp, err := http.Post(s.cfg.ScannerURL+"/scan", w.FormDataContentType(), &buf)
	if err != nil {
		return nil, apperrors.Internal("scanner service unavailable", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apperrors.Internal(fmt.Sprintf("scanner returned HTTP %d", resp.StatusCode), nil)
	}
	return io.ReadAll(resp.Body)
	
	
	
}

func (s *DocumentService) ListAllDoc(userID uint , page , limit int) ([]models.ResponeDoc,int64, error){
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	docs, total, err := s.repo.ListByUser(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	out := make([]models.ResponeDoc, len(docs))
	for i, d := range docs {
		out[i] = models.ToDocumentResponse(&d)
	}
	return out, total, nil
}


