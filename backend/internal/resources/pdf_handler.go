package resources

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
)

const maxPDFSize = 10 << 20

func ImportPDFHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxPDFSize)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Choose a PDF file up to 10 MB."})
		return
	}
	defer file.Close()

	if !strings.EqualFold(filepath.Ext(header.Filename), ".pdf") || header.Size > maxPDFSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Choose a PDF file up to 10 MB."})
		return
	}
	data, err := io.ReadAll(io.LimitReader(file, maxPDFSize+1))
	if err != nil || len(data) > maxPDFSize || !bytes.HasPrefix(data, []byte("%PDF-")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "That file is not a valid PDF or is too large."})
		return
	}

	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't read this PDF. Try an unencrypted PDF with selectable text."})
		return
	}
	textReader, err := reader.GetPlainText()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't extract text from this PDF."})
		return
	}
	content, err := io.ReadAll(io.LimitReader(textReader, 2<<20))
	if err != nil || strings.TrimSpace(string(content)) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No selectable text was found in this PDF. Scanned PDFs are not supported."})
		return
	}

	title := strings.TrimSpace(strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename)))
	if title == "" {
		title = "Untitled PDF"
	}
	err = CreateResource(c.MustGet("userID").(string), CreateResourceRequest{
		Title: title, Type: "pdf", Description: fmt.Sprintf("Imported from %s", header.Filename),
		SourceURL: "upload://" + filepath.Base(header.Filename), Content: string(content),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't save this resource."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "PDF imported successfully"})
}
