package protocol

// HTTPResponse HTTP响应
//
//	author centonhuang
//	update 2025-10-31 01:38:26
type HTTPResponse[BodyT any] struct {
	Body BodyT `json:"data"`
}
