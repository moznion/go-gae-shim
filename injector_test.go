package gaeshim

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func Test_BeforeAll(t *testing.T) {
	if os.Getenv("__GAE-SHIM-TEST-FOO") != "" {
		t.Errorf("ENV[__GAE-SHIM-TEST-FOO] is not empty")
	}

	if os.Getenv("__GAE-SHIM-TEST-BUZ") != "" {
		t.Errorf("ENV[__GAE-SHIM-TEST-BUZ] is not empty")
	}

	err := Inject("test/test.yml")
	if err != nil {
		t.Fatalf("unexpected error has occurred: %v", err)
	}
}

func Test_ShouldSetEnvironmentVariables(t *testing.T) {
	if os.Getenv("__GAE-SHIM-TEST-FOO") != "bar" {
		t.Errorf("ENV[__GAE-SHIM-TEST-FOO] gives unexpected error")
	}

	if os.Getenv("__GAE-SHIM-TEST-BUZ") != "qux" {
		t.Errorf("ENV[__GAE-SHIM-TEST-BUZ] gives unexpected error")
	}
}

func TestFoo(t *testing.T) {
	port, err := retrieveEmptyPort()
	if err != nil {
		t.Fatalf("%s", err)
	}
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

	baseURL := fmt.Sprintf("http://localhost:%d", port)

	cnt := 0
	c := time.Tick(100 * time.Millisecond)
	for range c {
		_, err := http.Get(baseURL)
		if err == nil {
			break
		}

		if cnt++; cnt > 20 {
			t.Fatalf("could not connect to test httpd")
		}
	}

	imgBag := map[string]string{
		"img.jpg": "image/jpeg",
		"img.png": "image/png",
	}

	for img, contentType := range imgBag {
		r, _ := http.Get(fmt.Sprintf("%s/img/%s", baseURL, img))
		if r.StatusCode != 200 {
			t.Errorf("could not get %s", img)
		}
		if gotContentType := r.Header.Get("Content-Type"); gotContentType != contentType {
			t.Errorf("unexpected Content-Type has come: expected=%s, got=%s", contentType, gotContentType)
		}
	}

	videos := []string{"video.mp4", "video.mov"}
	for _, video := range videos {
		r, _ := http.Get(fmt.Sprintf("%s/video/%s", baseURL, video))
		if r.StatusCode != 200 {
			t.Errorf("could not get %s", video)
		}
		if gotContentType := r.Header.Get("Content-Type"); gotContentType != "video/mp4" {
			t.Errorf("unexpected Content-Type has come (maybe `mime_type` doesn't work): expected=video/mp4, got=%s", gotContentType)
		}
		if testHeader1 := r.Header.Get("X-TEST-HEADER-1"); testHeader1 != "test1" {
			t.Error("unexpected X-TEST-HEADER-1 has come")
		}
		if testHeader2 := r.Header.Get("X-TEST-HEADER-2"); testHeader2 != "test2" {
			t.Error("unexpected X-TEST-HEADER-2 has come")
		}
	}
}

func retrieveEmptyPort() (int, error) {
	for port := 10000; port < 20000; port++ {
		addr := fmt.Sprintf("localhost:%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			defer l.Close()
			return port, nil
		}
	}
	return 0, errors.New("could not retrieve empty port")
}
