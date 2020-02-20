package rlfilter

import "net/http"

type Limiter interface {
	Limit(handler WebHandler) WebHandler
}
type WebHandler func(resp http.ResponseWriter, req *http.Request)
