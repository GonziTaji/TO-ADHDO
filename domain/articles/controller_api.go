package articles

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/articles/model"
	"github.com/yogusita/to-adhdo/domain/uploads"
)

// --- DTOs ---

type ArticleListItemDTO struct {
	Id                string               `json:"Id"`
	Name              string               `json:"Name"`
	Price             int                  `json:"Price"`
	ThumbnailUrl      string               `json:"ThumbnailUrl"`
	Tags              []ArticleTagDTO      `json:"Tags"`
	Prices            []ArticlePriceDTO    `json:"Prices"`
	AvailableForTrade bool                 `json:"AvailableForTrade"`
	Condition         *ArticleConditionDTO `json:"Condition,omitempty"`
}

type ArticleDetailDTO struct {
	Id                string               `json:"Id"`
	Name              string               `json:"Name"`
	Description       string               `json:"Description"`
	Price             int                  `json:"Price"`
	ThumbnailUrl      string               `json:"ThumbnailUrl"`
	ImagesUrls        []string             `json:"ImagesUrls"`
	Tags              []ArticleTagDTO      `json:"Tags"`
	Prices            []ArticlePriceDTO    `json:"Prices"`
	IsDeleted         bool                 `json:"IsDeleted"`
	AvailableForTrade bool                 `json:"AvailableForTrade"`
	Condition         *ArticleConditionDTO `json:"Condition,omitempty"`
}

type ArticleTagDTO struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type ArticlePriceDTO struct {
	Price       int    `json:"Price"`
	Description string `json:"Description"`
}

type ArticleConditionDTO struct {
	Slug  string `json:"Slug"`
	Label string `json:"Label"`
}

type UploadResponseDTO struct {
	Filename string `json:"Filename"`
	Url      string `json:"Url"`
}

// --- Catalog API ---

func (c *Controller) ApiCatalogHandler(ctx *gin.Context) {
	var options model.CatalogFilterOptions
	if err := ctx.Bind(&options); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	catalogData, err := c.service.Catalog(options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]ArticleListItemDTO, 0, len(catalogData.Articles))
	for _, a := range catalogData.Articles {
		dto := ArticleListItemDTO{
			Id:                a.Id,
			Name:              a.Name,
			Price:             a.Price,
			ThumbnailUrl:      a.ThumbnailUrl,
			AvailableForTrade: a.AvailableForTrade,
		}
		for _, t := range a.Tags {
			dto.Tags = append(dto.Tags, ArticleTagDTO{Name: t.Name})
		}
		if a.Condition.Slug != "" {
			dto.Condition = &ArticleConditionDTO{
				Slug:  a.Condition.Slug,
				Label: a.Condition.Label,
			}
		}
		items = append(items, dto)
	}

	ctx.JSON(http.StatusOK, items)
}

// --- Articles admin API ---

func (c *Controller) ApiListHandler(ctx *gin.Context) {
	var options ListingArticlesOptions
	if err := ctx.ShouldBindQuery(&options); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articles, err := c.service.List(&options)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]ArticleListItemDTO, 0, len(articles))
	for _, a := range articles {
		dto := ArticleListItemDTO{
			Id:                a.Id,
			Name:              a.Name,
			AvailableForTrade: a.AvailableForTrade,
		}
		for _, t := range a.Tags {
			dto.Tags = append(dto.Tags, ArticleTagDTO{Id: t.Id, Name: t.Name})
		}
		for _, p := range a.Prices {
			dto.Prices = append(dto.Prices, ArticlePriceDTO{Price: p.Price, Description: p.Description})
			if dto.Price == 0 {
				dto.Price = p.Price
			}
		}
		for _, img := range a.Images {
			if dto.ThumbnailUrl == "" {
				dto.ThumbnailUrl = img.Url
			}
		}
		items = append(items, dto)
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *Controller) ApiGetHandler(ctx *gin.Context) {
	articleId := ctx.Param("article_id")
	if articleId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article_id is required"})
		return
	}

	article, err := c.service.GetDetails(articleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if article.Id == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	dto := ArticleDetailDTO{
		Id:                article.Id,
		Name:              article.Name,
		Description:       article.Description,
		Price:             article.Price,
		ImagesUrls:        article.ImagesUrls,
		IsDeleted:         article.IsDeleted,
		AvailableForTrade: article.AvailableForTrade,
	}

	if len(article.ImagesUrls) > 0 {
		dto.ThumbnailUrl = article.ImagesUrls[0]
	}

	for _, t := range article.Tags {
		dto.Tags = append(dto.Tags, ArticleTagDTO{Id: t.Id, Name: t.Name})
	}

	if article.Condition.Slug != "" {
		dto.Condition = &ArticleConditionDTO{
			Slug:  article.Condition.Slug,
			Label: article.Condition.Label,
		}
	}

	ctx.JSON(http.StatusOK, dto)
}

func (c *Controller) ApiCreateHandler(ctx *gin.Context) {
	var form ArticleFormData
	if err := ctx.Bind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articleId, err := c.service.CreateFromForm(form)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Location", "/articles/"+articleId+"/edit")
	ctx.JSON(http.StatusCreated, gin.H{"Id": articleId})
}

func (c *Controller) ApiUpdateHandler(ctx *gin.Context) {
	articleId := ctx.Param("article_id")

	var form ArticleFormData
	if err := ctx.Bind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form.Id = articleId

	if err := c.service.UpdateFromForm(form); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Id": articleId})
}

func (c *Controller) ApiDeleteHandler(ctx *gin.Context) {
	articleId := ctx.Param("article_id")
	if articleId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "article_id is required"})
		return
	}

	if err := c.service.Delete(articleId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) ApiUploadImageHandler(ctx *gin.Context) {
	fileheader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fileheader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileExt := filepath.Ext(fileheader.Filename)

	savedFilename, err := uploads.SaveFile(articles_images_bucket, fileExt, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := file.Close(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, UploadResponseDTO{
		Filename: savedFilename,
		Url:      uploads.GetFilePublicUrl(articles_images_bucket, savedFilename),
	})
}
