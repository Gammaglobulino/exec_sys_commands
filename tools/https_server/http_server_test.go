package https_server

import (
	"../../client_server_connection/core/handle_connections"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func TestHTTPsServer(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":8181"

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "You requested:"+request.URL.Path)
	})
	log.Println("Listening to..", localip)
	err = http.ListenAndServeTLS(localip, "PEMCertificate", "privatePEM", nil)
	if err != nil {
		log.Fatal(err)
	}
}
