package model

func BuildIndexArticleResponse(maxCount int, subject string) []Article {
	resp := make([]Article, maxCount)
	DB.Table("articles").Limit(8).Order("create_time asc").Where("status='active' and subject=?", subject).Scan(&resp)
	return resp
}
