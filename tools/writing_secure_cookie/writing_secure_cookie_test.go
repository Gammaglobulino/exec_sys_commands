package writing_secure_cookie

import (
	"../writing_secure_cookie"
	"log"
	"net/http"
	"testing"
)

func TestWriteSecureCookie(t *testing.T) {
	http.HandleFunc("/", writing_secure_cookie.IndexHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
