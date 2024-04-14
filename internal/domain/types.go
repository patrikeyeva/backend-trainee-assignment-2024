package domain

import "time"

type ResponseUserBanner struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type ResponseBanner struct {
	BannerID  int                `json:"banner_id"`
	TagIds    []int              `json:"tag_ids"`
	FeatureID int                `json:"feature_id"`
	Content   ResponseUserBanner `json:"content"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type CacheUserBanner struct {
	Title      string
	Text       string
	Url        string
	IsActive   bool
	Expiration time.Time
}

type PostRequestBody struct {
	TagIds    []int              `json:"tag_ids"`
	FeatureID int                `json:"feature_id"`
	Content   ResponseUserBanner `json:"content"`
	IsActive  bool               `json:"is_active"`
}
