package imageio

import (
	"log"
	"testing"
)

func TestGetRomoteFile(t *testing.T) {
	plat := GetPlatform()
	err := GetRomoteFile(FNAME_PER_PLATFORM[plat])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Done.")
}
