package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url""
)

type Proxt struct {
	TargetURL *url.URL
}
