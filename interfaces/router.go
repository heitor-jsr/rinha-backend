package interfaces

import "net/http"

type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
