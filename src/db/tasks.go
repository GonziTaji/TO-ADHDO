package db

type TaskTemplate struct {
	Id          string
	Name        string
	Description string
	TagsIds     []string
	CreatedAt   int
	UpdatedAt   int
	DeletedAt   int
}

type Tag struct {
	Id        string
	Name      string
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}

func GetAvailableTaskTemplates(limit int8) {
}

func GetAllTaskTags(limit int8) {
}

func SaveTask() {
}
