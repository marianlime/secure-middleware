package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	
)

type Proxy struct {
	TargetURL *url.URL
}

func NewProxy(target string) (*Proxy, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &Proxy{TargetURL: parsedURL}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	forwardReq, err := http.NewRequest(req.Method, p.TargetURL.String(), req.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	forwardReq.Header = req.Header

	client := &http.Client{}
	resp, err := client.Do(forwardReq)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Failed to show response body: %v", err)
	}
}	
