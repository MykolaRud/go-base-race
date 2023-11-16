package repositories

import (
	"DB_race/interfaces"
	"DB_race/models"
)

type ArticleRepository struct {
	interfaces.IDbHandler
}

func (repo *ArticleRepository) GetNextArticle() (article models.Article, err error) {

	row := repo.QueryRow("SELECT id, title, url, created_at FROM articles WHERE is_processed = 0 ORDER BY id")
	errScan := row.Scan(&article.Id, &article.Title, &article.URL, &article.CreatedAt)
	if errScan != nil {
		return models.Article{}, errScan
	}

	return article, nil
}

func (repo *ArticleRepository) SetArticleProcessed(id int) {
	repo.Execute("UPDATE articles SET is_processed=is_processed+1 WHERE id = ?", id)
}

func (repo *ArticleRepository) ResetArticles() {
	repo.Execute("UPDATE articles SET is_processed=0 WHERE 1")
}
