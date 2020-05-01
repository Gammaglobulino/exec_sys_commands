package storing_secure_password

import (
	"../storing_secure_password"
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	hashedPassword := storing_secure_password.HashPassword("---")
	assert.NotEmpty(t, hashedPassword)
	log.Println(hashedPassword)
	ioutil.WriteFile("amkey.key", []byte(hashedPassword), 0666)
}

func TestEnteringPasswordWholeProcess(t *testing.T) {
	//copy to a main function
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for text != "0111a24456fbeb7a17468dc6afc704b12111d573c25629ab1ab04de8efbbc222" { // exit only if hash-verified passoword
		fmt.Print("Enter your password: ")
		scanner.Scan()
		text = scanner.Text()
		text = storing_secure_password.HashPassword(text)
		if text != "0111a24456fbeb7a17468dc6afc704b12111d573c25629ab1ab04de8efbbc222" {
			fmt.Println("Your hashed password was: ", text)
		}
	}

}

func TestAMPassword(t *testing.T) {
	text := storing_secure_password.HashPassword("--")
	assert.EqualValues(t, text, "0111a24456fbeb7a17468dc6afc704b12111d573c25629ab1ab04de8efbbc222")
}
