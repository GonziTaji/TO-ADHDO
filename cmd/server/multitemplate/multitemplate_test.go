package multitemplate

import (
	"fmt"
	"strings"
	"testing"
)

func validate(input_path string, expected_name string) error {
	if out := getTemplateName(input_path); !strings.EqualFold(out, expected_name) {
		return fmt.Errorf("Invalid template name. Expected: %s; Got: %s ", expected_name, out)
	}

	return nil
}

func TestGetTemplateNamePage(t *testing.T) {
	err := validate(
		"domain/articles/pages/catalog/page.html",
		"articles/pages/catalog",
	)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTemplateNameComponentAsDir(t *testing.T) {
	err := validate(
		"domain/wishlist/components/card/index.html",
		"domain/wishlist/card",
	)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTemplateNameComponentAsFile(t *testing.T) {
	err := validate(
		"domain/shared/partials/nav.html",
		"shared/partials/nav",
	)

	if err != nil {
		t.Error(err.Error())
	}
}
