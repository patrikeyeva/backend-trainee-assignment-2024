package db

type ResponseUserBanner struct {
	Title    string `db:"title"`
	Text     string `db:"text"`
	Url      string `db:"url"`
	IsActive bool   `db:"is_active"`
}
