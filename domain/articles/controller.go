package articles

import (
	"github.com/yogusita/to-adhdo/domain/tags"
)

const articles_images_bucket = "articles_images"

type Controller struct {
	store     *Store
	views     *Views
	tagsStore *tags.Store
}

func CreateController(store *Store, views *Views, tagsStore *tags.Store) *Controller {
	return &Controller{store, views, tagsStore}
}
