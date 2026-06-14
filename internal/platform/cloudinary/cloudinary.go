package cloudinary

import (
	"os"

	cloudinary "github.com/cloudinary/cloudinary-go/v2"
)

var Client *cloudinary.Cloudinary

func Init() {
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	client, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		panic("Cloudinary failed to initialize: " + err.Error())
	}
	Client = client
}
