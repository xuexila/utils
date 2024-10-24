package httpServer

import "sync"

var (
	LoginSessionMap sync.Map
	err             error
)
