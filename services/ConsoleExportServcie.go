package services

import (
	"DB_race/models"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type ConsoleExportService struct {
}

func (service *ConsoleExportService) Export(article models.Article) (path string, err error) {
	data := fmt.Sprintf("Title: %s\nUrl: %s\nCreated at: %s", article.Title, article.URL, article.CreatedAt)
	byteData := []byte(data)
	md5HexSum := md5.Sum(byteData)

	fmt.Println("article's ", article.Id, " md5 is ", hex.EncodeToString(md5HexSum[:]))

	return "console", nil
}
