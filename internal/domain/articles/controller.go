package articles

const articles_images_bucket = "articles_images"

type Controller struct {
	service *Service
}

func CreateController(service *Service) *Controller {
	return &Controller{service}
}
