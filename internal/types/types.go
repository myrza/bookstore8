package types

type Author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Biography string `json:"biography"`
	Birthday  string `json:"birthday"`
}

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID string `json:"authorid"`
	ISBN     string `json:"isbn"`
	Year     string `json:"year"`
}

type AuthorAndBook struct {
	AuthorID  int    `json:"author_id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Biography string `json:"biography"`
	Birthday  string `json:"birthday"`
	Title     string `json:"title"`
	ISBN      string `json:"isbn"`
	Year      string `json:"year"`
}
