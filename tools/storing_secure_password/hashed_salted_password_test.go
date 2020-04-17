package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	saltkey, err := GenerateSalt()
	assert.Nil(t, err)
	assert.NotEmpty(t, saltkey)
	log.Println(saltkey)
}

func TestHashPassword(t *testing.T) {
	hashedPassword := HashPassword("andreama@microsoft.com")
	assert.NotEmpty(t, hashedPassword)
	log.Println(hashedPassword)
}
