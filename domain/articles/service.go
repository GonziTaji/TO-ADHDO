package articles

import (
	"fmt"
	"slices"

	"github.com/yogusita/to-adhdo/domain/articles/model"
	"github.com/yogusita/to-adhdo/domain/tags"
)

type Service struct {
	store     *Store
	views     *Views
	tagsStore *tags.Store
}

func CreateService(store *Store, views *Views, tagsStore *tags.Store) *Service {
	return &Service{store, views, tagsStore}
}

func (s *Service) GetDetails(articleID string) (model.ArticleDetails, error) {
	return s.store.GetDetails(articleID)
}

func (s *Service) Catalog(options model.CatalogFilterOptions) (model.CatalogData, error) {
	catalogItems, err := s.store.Catalog(options)

	if err != nil {
		return model.CatalogData{}, err
	}

	allTags, err := s.tagsStore.List(tags.ListingTagsOptions{})

	if err != nil {
		return model.CatalogData{}, err
	}

	tagsOptions := make([]model.TagOption, 0, len(allTags))

	for _, tag := range allTags {
		to := model.TagOption{Id: tag.Id, Name: tag.Name}

		to.Selected = slices.ContainsFunc(options.TagsIdsFilter, func(tagidFilter string) bool {
			return tagidFilter == tag.Id
		})

		tagsOptions = append(tagsOptions, to)
	}

	return model.CatalogData{
		Articles: catalogItems,
		Tags:     tagsOptions,
		Options:  options,
	}, nil
}

func (s *Service) GetFormData(articleID string) (model.ArticleFormTemplateData, error) {
	article := model.Article{}

	if articleID != "" {
		var err error
		article, err = s.store.Get(articleID)

		if err != nil {
			return model.ArticleFormTemplateData{}, err
		}
	}

	allTags, err := s.tagsStore.List(tags.ListingTagsOptions{})

	if err != nil {
		return model.ArticleFormTemplateData{}, err
	}

	tagOptions := make([]model.TagOption, len(allTags))
	tagsIdsInArticle := make(map[string]bool)

	for _, tag := range article.Tags {
		tagsIdsInArticle[tag.Id] = true
	}

	for i, tag := range allTags {
		tagOptions[i] = model.TagOption{
			Name:     tag.Name,
			Id:       tag.Id,
			Disabled: tagsIdsInArticle[tag.Id],
		}
	}

	return model.ArticleFormTemplateData{
		Article:    article,
		TagOptions: tagOptions,
	}, nil
}

func (s *Service) List(options *ListingArticlesOptions) ([]model.Article, error) {
	return s.store.List(options)
}

func (s *Service) CreateFromForm(form ArticleFormData) (string, error) {
	article, err := buildArticleFromForm(form)

	if err != nil {
		return "", err
	}

	return s.store.Create(&article)
}

func (s *Service) UpdateFromForm(form ArticleFormData) error {
	article, err := buildArticleFromForm(form)

	if err != nil {
		return err
	}

	return s.store.Update(&article)
}

func (s *Service) Delete(articleID string) error {
	return s.store.Delete(articleID)
}

func buildArticleFromForm(form ArticleFormData) (model.Article, error) {
	if len(form.TagNames) != len(form.TagIds) {
		return model.Article{}, fmt.Errorf("tags names and ids mismatch")
	}

	if len(form.ArticleImageFilenames) != len(form.ArticleImageIds) {
		return model.Article{}, fmt.Errorf("article image filenames and ids mismatch")
	}

	article := model.Article{
		Id:                form.Id,
		Name:              form.Name,
		Description:       form.Description,
		Tags:              []tags.Tag{},
		AvailableForTrade: form.AvailableForTrade,
		Images:            []model.ArticleImage{},
	}

	for i, name := range form.TagNames {
		tag := tags.Tag{
			Id:   form.TagIds[i],
			Name: name,
		}

		article.Tags = append(article.Tags, tag)
	}

	for i, id := range form.ArticleImageIds {
		image := model.ArticleImage{
			Id:       id,
			Filename: form.ArticleImageFilenames[i],
		}

		article.Images = append(article.Images, image)
	}

	if form.NewPrice != 0 {
		article.Prices = append(article.Prices, model.ArticlePrice{
			Price:       form.NewPrice,
			Description: form.NewPriceDescription,
		})
	}

	return article, nil
}
