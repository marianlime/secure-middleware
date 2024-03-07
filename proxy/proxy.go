package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Proxy struct {
	targetURL *url.URL
	client    *http.Client
}

func NewProxy(target string) (*Proxy, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		targetURL: parsedURL,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	forwardReq, err := http.NewRequest(req.Method, p.targetURL.String() + req.URL.Path, req.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	forwardReq.RequestURI = ""

	forwardReq.Header = make(http.Header)
	for k, vv := range req.Header{
		for _, v := range vv {
			forwardReq.Header.Add(k, v)
		}
	}
	forwardReq.Host = p.targetURL.Host

	client := &http.Client{}
	resp, err := client.Do(forwardReq)
	if err != nil {
		log.Printf("Error forwarding request: %v", err)
		http.Error(w, "Failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Failed to copy response body: %v", err)
	}
}

func (p *Proxy) CloneRequest(req *http.Request) *http.Request {
	outReq := req.Clone(req.Context())
	outReq.URL = p.targetURL
	outReq.URL.Scheme = p.targetURL.Scheme
	outReq.URL.Host = p.targetURL.Host
	outReq.Header = make(http.Header)

	for k, vv := range req.Header {
		for _, v := range vv {
			outReq.Header.Add(k, v)
		}
	}

	removeHopByHopHeaders(outReq.Header)
	return outReq

}

func copyHeaders(dest, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dest.Add(k, v)
		}
	}
	removeHopByHopHeaders(dest)
}

func removeHopByHopHeaders(header http.Header) {
	hopByHopHeader := []string{"Connection", "Keep-Alive", "Proxy-Authen", "Proxy-Authorize", "Te", "Trailers", "Transfer-Encoding", "Upgrade"}
	for _, h := range hopByHopHeader {
		header.Del(h)
	}
}