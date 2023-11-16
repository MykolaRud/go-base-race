package services

import (
	"DB_race/models"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type FileExportService struct {
	catalog string
}

func (service *FileExportService) SetCatalog(catalog string) {
	service.catalog = catalog
}

func (service *FileExportService) Export(article models.Article) (path string, err error) {
	if service.catalog == "" {
		return "", errors.New("output catalog is not set")
	}

	filename := fmt.Sprintf("%s/%s.txt", service.catalog, strconv.Itoa(article.Id))
	data := fmt.Sprintf("Title: %s\nUrl: %s\nCreated at: %s", article.Title, article.URL, article.CreatedAt)
	byteData := []byte(data)

	fileErr := os.WriteFile(filename, byteData, 0666)
	if fileErr != nil {
		return "", fileErr
	}

	return filename, nil
}
