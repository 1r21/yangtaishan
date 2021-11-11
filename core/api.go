package core

type Api struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type BaseArticle struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Title    string `json:"title"`
	ImageUrl string `json:"cover"`
}

type Article struct {
	BaseArticle
	Source     string `json:"source"`
	AudioUrl   string `json:"audio_url"`
	Transcript string `json:"transcript"`
}
