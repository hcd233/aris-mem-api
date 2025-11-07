package protocol

// HTTPResponse HTTP响应
//
//	author centonhuang
//	update 2025-10-31 01:38:26
type HTTPResponse[BodyT any] struct {
	Body BodyT `json:"data"`
}

// RedirectResponse 重定向响应
//
//	@author centonhuang
//	@update 2025-11-02 04:01:39
type RedirectResponse struct {
	Status int    `json:"status" doc:"Status code"`
	Url    string `json:"url" doc:"URL for redirect"`
}
