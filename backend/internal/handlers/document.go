package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neathhatake/the_Scan/internal/middleware"
	"github.com/neathhatake/the_Scan/internal/services"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
	"github.com/neathhatake/the_Scan/pkg/respone"
)




type DocumentHandler struct {
	docservice *services.DocumentService
}


func NewDocumentHandler(docservice *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		docservice: docservice,
	}
}


// POST/scan

func(h *DocumentHandler) Scan(c *gin.Context){
	file , header , err := c.Request.FormFile("file")
	if err != nil {
		respone.Error(c, apperrors.BadRequest("file field is Required"))

		return 
	}
	defer file.Close()

	imageData , err := io.ReadAll(file)
	if err != nil {
		respone.Error(c, apperrors.Internal("failed to read file", err))
		return 
	}

	title := c.Query("title")
	if title == ""{
		title = "Untitled Scan"
	}
	doc , svcErr := h.docservice.ScanandSave(&services.Scan{
		UserID: middleware.UserID(c),
		FileName: header.Filename,
		ImageData: imageData,
		Title: title,

	})
	if svcErr != nil {
		respone.Error(c, svcErr)
		return
	}

	c.JSON(http.StatusCreated, respone.Envelop{
		Success: true,
		Data: doc,
	})
}




// GET /documents?page=1&limit=20

func (h *DocumentHandler) ListDocuments(c *gin.Context){
	page , _ := strconv.Atoi(c.DefaultQuery("page","1"))
	limit , _ := strconv.Atoi(c.DefaultQuery("limit","20"))

	doc , total , err := h.docservice.ListAllDoc(middleware.UserID(c), page , limit)
	if err != nil {
		respone.Error(c, err)
		return 
	}

	respone.OKPaginated(c,doc , total , page , limit)

}




/// GET/documents/:id/download

func (h *DocumentHandler) DownloadDocuments(c *gin.Context){
	id , err := parseID(c)
	if err != nil {
		respone.Error(c, err)
		return
	}

	path , title , svcErr := h.docservice.GetDocPath(uint(id), middleware.UserID(c))
	if svcErr != nil {
		respone.Error(c,err)
		return
	}
	c.FileAttachment(path,title+".pdf")

}


// PUT/documents/:id
func (h *DocumentHandler) RenameDocument(c *gin.Context){
	id , err := parseID(c)
	if err != nil {
		respone.Error(c, err)
		return
	}
	var req struct {
		Title string `json:"title" binding:"required,min=1,max=255"`
	} 
	if err := c.ShouldBindJSON(&req); err != nil {
		respone.Error(c, respone.ValidationError(err))
		return
	}

	doc , svcErr := h.docservice.RenameDoc(uint(id),middleware.UserID(c),req.Title)
	if svcErr != nil {
		respone.Error(c, svcErr)
		return
	}
	respone.OK(c, doc)

	
}


func parseID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return 0, apperrors.BadRequest("invalid id parameter")
	}
	return id, nil
}




// DELETE/documents/:id

func (h *DocumentHandler) DeleteDocuments(c *gin.Context){
	id, err := parseID(c)
	if err != nil {
		respone.Error(c, err)
		return
	}
	if svcErr := h.docservice.DeleteDoc(uint(id), middleware.UserID(c)); svcErr != nil {
		respone.Error(c, svcErr)
		return
	}
	respone.NoContent(c)
}