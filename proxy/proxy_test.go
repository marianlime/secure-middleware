package proxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startMockServer() (*httptest.Server, string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(r.URL.Path))
	})
	server := httptest.NewServer(handler)

	return server, server.URL
}

func TestRequestForwarding(t *testing.T) {
	mockServer, mockServerURL := startMockServer()
	defer mockServer.Close()

	proxyInstance, err := NewProxy(mockServerURL)
	if err != nil {
		t.Fatalf("Failed to create new proxy : %v", err)
	}

	proxyServer := httptest.NewServer(proxyInstance)
	defer proxyServer.Close()

	testPath := "/testpath"
	resp, err := http.Get(proxyServer.URL + testPath)
	if err != nil {
		t.Fatalf("Failed to make request to proxy : %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if got := string(body); got != testPath {
		t.Errorf("Expected path '%s', got '%s'", testPath, got)
	}

}

func TestRequestForwardingWithDifferentMehtods(t *testing.T) {
	methods := []string{"POST", "PUT", "DELETE"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					t.Errorf("Expected '%s', got '%s'", method, r.Method)
				}
			}))
			defer mockServer.Close()
			proxyInstance, err := NewProxy(mockServer.URL)
			if err != nil {
				t.Fatalf("Failed to create a new proxy : %v", err)
			}
			proxyServer := httptest.NewServer(proxyInstance)
			defer proxyServer.Close()

			req, _ := http.NewRequest(method, proxyServer.URL, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to make %s request to proxy : %v", method, err)
			}
			resp.Body.Close()
		})

	}
}

func TestHeaderForwarding(t *testing.T) {
	headerKey := "XYZ-Test-Header"
	headerValue := "HeaderValue"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerKey) != headerValue {
			t.Errorf("Expected header '%s' to be forwarded", headerKey)
		}
	}))
	defer mockServer.Close()

	proxyInstance, err := NewProxy(mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to creata a new proxy : %v", err)
	}
	proxyServer := httptest.NewServer(proxyInstance)
	defer proxyServer.Close()

	req, err := http.NewRequest("GET", proxyServer.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request : %v", err)
	}
	req.Header.Add(headerKey, headerValue)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request with header : %v", err)
	}
	defer resp.Body.Close()
}

func TestErrorHandling(t *testing.T) {
	invalidURL := "http://127.0.0.1:9999"
	proxyInstance, err := NewProxy(invalidURL)
	if err != nil {
		t.Fatalf("Failed to create proxy instance: %v", err)
	}
	proxyServer := httptest.NewServer(proxyInstance)
	defer proxyServer.Close()

	resp, err := http.Get(proxyServer.URL)
	if err != nil {
		t.Fatalf("Failed to make request to proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status code %d, got %d", http.StatusBadGateway, resp.StatusCode)
	}
}