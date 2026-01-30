package enum

// ArticleStatus 文章状态
//
//	@author centonhuang
//	@update 2026-01-28 21:18:34
type ArticleStatus string

const (
	// ArticleStatusPublished 已发布
	//	@author centonhuang
	//	@update 2026-01-28 21:18:34
	ArticleStatusPublished ArticleStatus = "published"

	// ArticleStatusHidden 隐藏
	//	@author centonhuang
	//	@update 2026-01-28 21:18:34
	ArticleStatusHidden ArticleStatus = "hidden"
)
