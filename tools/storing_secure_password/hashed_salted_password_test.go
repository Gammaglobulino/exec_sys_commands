package storing_secure_password

import (
	"../storing_secure_password"
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	saltkey, err := storing_secure_password.GenerateSalt()
	assert.Nil(t, err)
	assert.NotEmpty(t, saltkey)
	log.Println(saltkey)
}

func TestHashPassword(t *testing.T) {
	hashedPassword := storing_secure_password.HashPassword("pan")
	assert.NotEmpty(t, hashedPassword)
	log.Println(hashedPassword)
}

func TestEnteringPasswordWholeProcess(t *testing.T) {
	//copy to a main function
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for text != "fff4f0d581e3bbcf3f8c944ba24a6932a23a0619314c63fed9ab4d482ef81411" { // exit only if hash-verified passoword
		fmt.Print("Enter your password: ")
		scanner.Scan()
		text = scanner.Text()
		text = storing_secure_password.HashPassword(text)
		if text != "fff4f0d581e3bbcf3f8c944ba24a6932a23a0619314c63fed9ab4d482ef81411" {
			fmt.Println("Your hashed password was: ", text)
		}
	}

}
