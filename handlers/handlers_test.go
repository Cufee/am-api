package handlers

import (
	"log"
	"testing"
)

func TestCDNUpload(t *testing.T) {
	pid := 1013072123
	imageURL := "http://res.cloudinary.com/vkodev/image/upload/v1602915772/Aftermath/202905960405139456.jpg"

	url, err := uploadToCDN(imageURL, pid)
	if err != nil {
		log.Println(err)
		t.FailNow()
		return
	}
	log.Println(url)
}
