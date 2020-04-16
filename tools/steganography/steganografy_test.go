package steganography

import (
	"github.com/stretchr/testify/assert"
	image2 "image"
	"image/jpeg"
	"math/rand"
	"os"
	"testing"
)

func TestGenerateImageWithRandomPixels(t *testing.T) {
	image := image2.NewRGBA(image2.Rect(0, 0, 100, 200))
	for p := 0; p < (100 * 200); p++ {
		pixelOffset := 4 * p
		image.Pix[0+pixelOffset] = uint8(rand.Intn(256)) //R
		image.Pix[1+pixelOffset] = uint8(rand.Intn(256)) //G
		image.Pix[2+pixelOffset] = uint8(rand.Intn(256)) //B
		image.Pix[3+pixelOffset] = 255                   //Alpha
	}

	outFile, err := os.Create("randomizedImage.jpg")
	assert.Nil(t, err)
	defer outFile.Close()
	jpeg.Encode(outFile, image, nil)
	err = outFile.Close()
	assert.Nil(t, err)

}
