package handlers

import (
	"log"
	"testing"
)

func TestIsNSFW(t *testing.T) {
	url := "https://images-ext-2.discordapp.net/external/0zMEE4p_lWKxQ2ezlESyyiTNeW4AqPvH6ddM2n9WpTk/https/media.discordapp.net/attachments/792622562118074371/820861677369032735/USER_SCOPED_TEMP_DATA_orca-image--259473315.jpeg"
	ok, err := isNSFW(url)
	log.Print(ok)
	if err != nil {
		log.Print(err.Error())
		t.FailNow()
	}
}
