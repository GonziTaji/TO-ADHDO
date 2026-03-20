package articles

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/tags"
)

// PageHandler is a function that serves an HTML page. Injected by the server layer.
type PageHandler = gin.HandlerFunc

func RegisterRoutes(router *gin.Engine, db *sql.DB, servePage func(string) gin.HandlerFunc) {
	store := CreateStore(db)
	tagsStore := tags.CreateStore(db)
	service := CreateService(store, &Views{}, tagsStore)
	controller := CreateController(service)

	// -----------------------------------------------------------------------
	// Static assets for old template-based pages (kept for compatibility)
	// -----------------------------------------------------------------------
	router.Group("catalog").Static("/static", "domain/articles/static")
	router.Group("admin/articles").Static("/static", "domain/articles/static")

	// -----------------------------------------------------------------------
	// Page routes – serve built React HTML files
	// -----------------------------------------------------------------------
	router.GET("/catalog", servePage("articles/catalog.html"))
	router.GET("/catalog/:article_id", servePage("articles/view.html"))

	router.GET("/articles", servePage("articles/list.html"))
	router.GET("/articles/new", servePage("articles/form.html"))
	router.GET("/articles/:article_id/edit", servePage("articles/form.html"))

	// -----------------------------------------------------------------------
	// Backward-compatibility redirects
	// -----------------------------------------------------------------------
	router.GET("/catalog/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/catalog")
	})
	router.GET("/admin/articles/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/articles")
	})
	router.GET("/admin/articles/new", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/articles/new")
	})
	router.GET("/admin/articles/:article_id/edit", func(ctx *gin.Context) {
		id := ctx.Param("article_id")
		ctx.Redirect(http.StatusMovedPermanently, "/articles/"+id+"/edit")
	})

	// -----------------------------------------------------------------------
	// JSON API routes
	// -----------------------------------------------------------------------
	api := router.Group("/api")

	api.GET("/catalog", controller.ApiCatalogHandler)

	api.GET("/articles", controller.ApiListHandler)
	api.GET("/articles/:article_id", controller.ApiGetHandler)
	api.POST("/articles", controller.ApiCreateHandler)
	api.PUT("/articles/:article_id", controller.ApiUpdateHandler)
	api.DELETE("/articles/:article_id", controller.ApiDeleteHandler)
	api.POST("/articles/uploads", controller.ApiUploadImageHandler)
}
