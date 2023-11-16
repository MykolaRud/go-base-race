package repositories

import (
	"DB_race/interfaces"
	"DB_race/models"
	"sync"
)

type ArticleRepository struct {
	DBHandler interfaces.IDbHandler
	mu        sync.Mutex
}

func (repo *ArticleRepository) LockNextArticle() (article models.Article, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	article, err = repo.GetNextArticle()

	if err != nil {
		return article, err
	}

	repo.SetArticleProcessed(article.Id)

	return article, err
}

func (repo *ArticleRepository) GetNextArticle() (article models.Article, err error) {
	row := repo.DBHandler.QueryRow("SELECT id, title, url, created_at FROM articles WHERE is_processed = 0 ORDER BY id")
	errScan := row.Scan(&article.Id, &article.Title, &article.URL, &article.CreatedAt)
	if errScan != nil {
		return models.Article{}, errScan
	}

	return article, nil
}

func (repo *ArticleRepository) SetArticleProcessed(id int) {
	repo.DBHandler.Execute("UPDATE articles SET is_processed=is_processed+1 WHERE id = ?", id)
}

func (repo *ArticleRepository) ResetArticles() {
	repo.DBHandler.Execute("UPDATE articles SET is_processed=0 WHERE 1")
}
