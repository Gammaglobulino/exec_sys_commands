package writing_secure_cookie

import (
	"../secure_crypting_AES"
	"fmt"
	"log"
	"net/http"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pk, err := secure_crypting_AES.GeneratePrimaryKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(pk))
	secureSessionCookie := http.Cookie{
		Name:       "GammaCookie",
		Value:      pk,
		Path:       "/",
		Domain:     "vaultgamma.com",
		Expires:    time.Now().Add(60 * time.Minute),
		RawExpires: "",
		Secure:     true,
		HttpOnly:   true,
	}
	http.SetCookie(w, &secureSessionCookie)
	fmt.Fprintln(w, "Cookie has been set")

}
