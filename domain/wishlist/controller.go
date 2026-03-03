package wishlist

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/domain/tags"
)

type Controller struct {
	store *Store
}

func CreateController(store *Store) Controller {
	return Controller{store}
}

func (c *Controller) GetAdminListHandler(ctx *gin.Context) {
	list, err := c.store.GetAdminList()

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "wishlist/wishitem/list", gin.H{
		"List": list,
	})
}

func (c *Controller) GetListHandler(ctx *gin.Context) {
	list, err := c.store.GetWishlist()

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "wishlist", gin.H{
		"List": list,
	})
}

func (c *Controller) GetHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	wi, err := c.store.GetWishitem(id)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "wishlist/wishitem/view", wi)
}

func (c *Controller) GetFormHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	var wi Wishitem
	var err error

	log.Printf("id: %s\n", id)

	if len(id) > 0 {
		wi, err = c.store.GetWishitem(id)

		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		wi.Id = ""
		wi.Name = ""
		wi.ExternalUrl = ""
		wi.ObservedPrice = "0"
	}

	// TODO: get tags
	tags := []tags.Tag{}

	ctx.HTML(http.StatusOK, "wishlist/wishitem/form", WishitemFormTemplateData{
		Record: wi,
		Tags:   tags,
	})
}

func (c *Controller) CreateHandler(ctx *gin.Context) {
	var fd WishitemFormData

	if err := ctx.Bind(&fd); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.store.SaveWishitem(fd)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Header("location", "/wishlist/"+id)
	ctx.Status(http.StatusCreated)
}

func (c *Controller) UpdateHandler(ctx *gin.Context) {
	var fd WishitemFormData

	if err := ctx.Bind(&fd); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err := c.store.SaveWishitem(fd)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) DeleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.store.DeleteWishitem(id); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) GetPreview(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Preview is not yet available")

	// TODO: scrapper
	// linkpreview.net could not bring the metadata from mercadolibre. Neither a direct request.
	//
	// url := ctx.Query("url")
	//
	// api_key, exists := os.LookupEnv("LINKPREVIEW_API_KEY")
	//
	// if !exists || len(api_key) == 0 {
	// 	log.Println("Env var \"LINKPREVIEW_API_KEY\" not set or empty")
	// 	ctx.String(http.StatusInternalServerError, "Could not load preview")
	// 	return
	// }
	//
	// req, err := http.NewRequest("GET", fmt.Sprintf("https://api.linkpreview.net?q=%s", url), nil)
	//
	// if err != nil {
	// 	ctx.String(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	//
	// req.Header.Set("X-Linkpreview-Api-Key", api_key)
	// req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WishAppBot/1.0)")
	//
	// res, err := http.DefaultClient.Do(req)
	//
	// if err != nil {
	// 	ctx.String(http.StatusInternalServerError, "Error: "+err.Error())
	// 	return
	// }
	// defer res.Body.Close()
	//
	// if res.StatusCode != 200 {
	// 	ctx.String(http.StatusBadGateway, "Error: "+res.Status)
	// }
	//
	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// if err != nil {
	// 	ctx.String(http.StatusInternalServerError, "Error: "+err.Error())
	// 	return
	// }
	//
	// metadata := WishitemMeta{
	// 	URL: url,
	// }
	//
	// // OpenGraph tags
	// metadata.Title, _ = doc.Find("meta[property='og:title']").Attr("content")
	// metadata.Description, _ = doc.Find("meta[property='og:description']").Attr("content")
	// metadata.Image, _ = doc.Find("meta[property='og:image']").Attr("content")
	// price := doc.Find("meta[itemprop='price']").AttrOr("content", "0")
	// metadata.Price, err = strconv.Atoi(price)
	//
	// if err != nil {
	// 	log.Printf("Error transforming price from page. Value found: %s\n. Error: %s", price, err.Error())
	// }
	//
	// // title fallback
	// if metadata.Title == "" {
	// 	metadata.Title = doc.Find("title").First().Text()
	// }
	//
	// log.Printf("meta: %v\n", metadata)
	//
	// body, err := io.ReadAll(res.Body)
	//
	// if err != nil {
	// 	ctx.String(http.StatusInternalServerError, "Error: "+err.Error())
	// 	return
	// }
	//
	// ctx.String(http.StatusOK, fmt.Sprintf("%s", body))
}
