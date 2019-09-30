package web

import "net/http"

type IFilter interface {
	RegisterInfo() (string, func(http.ResponseWriter, *http.Request))
}
