package main

import (
	"DB_race/infrastructures"
	"DB_race/interfaces"
	"DB_race/repositories"
	"DB_race/services"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"log"
	"sync"
)

var (
	MySQLConfig = mysql.Config{
		ParseTime: true,
	}
	Repo *repositories.ArticleRepository
)

func main() {
	Setup()

	container := dig.New()
	container.Provide(initDBConnection)
	container.Provide(initDBHandler)
	container.Invoke(initRepo)

	Repo.ResetArticles()

	var wgArticleExport sync.WaitGroup

	wgArticleExport.Add(1)
	go func() {
		ExportToFileWorkerLoop()
		wgArticleExport.Done()
	}()

	wgArticleExport.Add(1)
	go func() {
		ExportToConsoleWorkerLoop()
		wgArticleExport.Done()
	}()

	wgArticleExport.Add(1)
	go func() {
		ExportToEmailWorker()
		wgArticleExport.Done()
	}()

	wgArticleExport.Wait()
}

func ExportToEmailWorker() {
	emailExport := services.EmailExportService{}
	emailExport.SetCredentials(viper.GetString("smtp_host"), viper.GetString("smtp_port"), viper.GetString("smtp_username"), viper.GetString("smtp_password"))

	for {
		article, err := Repo.LockNextArticle()
		if err != nil {
			fmt.Printf("get article error: %s", err)
			break
		}

		subject, exportErr := emailExport.Export(article)
		if exportErr != nil {
			fmt.Println("email export article error: ", exportErr.Error())
			break
		}

		fmt.Println("  exported ", article.Id, " to email ", subject)

	}
}

func ExportToFileWorkerLoop() {
	fileExport := services.FileExportService{}
	fileExport.SetCatalog("./exportData")

	for {
		article, err := Repo.LockNextArticle()
		if err != nil {
			fmt.Printf("get article error: %s", err)
			break
		}

		filename, exportErr := fileExport.Export(article)
		if exportErr != nil {
			fmt.Printf("file export article error: %s", exportErr.Error())
			break
		}

		fmt.Println("  exported ", article.Id, " to file ", filename)

	}
}

func ExportToConsoleWorkerLoop() {

	consoleExport := services.ConsoleExportService{}

	for {
		article, err := Repo.LockNextArticle()
		if err != nil {
			fmt.Printf("get article error: %s", err)
			break
		}

		filename, exportErr := consoleExport.Export(article)
		if exportErr != nil {
			fmt.Printf("console export article error: %s", exportErr.Error())
			break
		}

		fmt.Println("  exported ", article.Id, " to console ", filename)
	}
}

func Setup() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	MySQLConfig = mysql.Config{
		User:      viper.GetString("db_user"),
		Passwd:    viper.GetString("db_password"),
		Net:       "tcp",
		Addr:      viper.GetString("db_address"),
		DBName:    viper.GetString("db_name"),
		ParseTime: true,
	}
}

func initDBConnection() *sql.DB {
	var err error

	Conn, err := sql.Open("mysql", MySQLConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	return Conn
}

func initDBHandler() interfaces.IDbHandler {
	return &infrastructures.MySQLHandler{}
}

func initRepo(dbHandler interfaces.IDbHandler, Conn *sql.DB) {
	dbHandler.SetConn(Conn)

	Repo = &repositories.ArticleRepository{DBHandler: dbHandler}
}
