package main

import (
	"article/config"
	docs "article/docs" // docs is generated by Swag CLI, you have to import it.
	"article/handlers"
	"article/storage"
	"article/storage/postgres"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	conf := config.Load()
	AUTH := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.PostgresHost,
		conf.PostgresPort,
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresDatabase,
	)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = conf.App
	docs.SwaggerInfo.Version = conf.AppVersion

	var DB storage.StorageI
	var err error
	DB, err = postgres.InitDB(AUTH)
	if err != nil {
		panic(err)
	}
	h := handlers.Handler{
		IM:   DB,
		Conf: conf,
	}
	if conf.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	if conf.Environment != "production" {
		router.Use(gin.Logger(), gin.Recovery()) // Later they will be replaced by custom Logger and Recovery
	}
	v1 := router.Group("/v2")
	{
		v1.POST("/article", h.CreateArticle)
		v1.GET("/article/:id", h.GetArticleByID)
		v1.GET("/article", h.GetArticleList)
		v1.PUT("/article", h.UpdateArticle)
		v1.DELETE("/article/:id", h.DeleteArticle)

		v1.POST("/author", h.CreateAuthor)
		v1.GET("/author/:id", h.GetAuthorByID)
		v1.GET("/author", h.GetAuthorList)
		v1.PUT("/author", h.UpdateAuthor)
		v1.DELETE("/author/:id", h.DeleteAuthor)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
