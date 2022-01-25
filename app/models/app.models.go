package models

type SkuCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SkuContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SkuRetailLink struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	BrandLogo File   `json:"brandLogo"`
}

type File struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Subtask struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UserSubtask struct {
	Subtask
	Status     string `json:"status"`
	IsComplete bool   `json:"isComplete"`
}
