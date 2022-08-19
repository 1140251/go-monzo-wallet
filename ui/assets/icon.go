package assets

import (
	"bytes"
	"embed"
	"image"
	"strings"
)

//go:embed *
var content embed.FS

var Icons map[string]image.Image

func init() {
	icons, err := ReadIcons()
	if err != nil {
		panic("Error loading icons")
	}

	Icons = icons
}

func ReadIcons() (map[string]image.Image, error) {
	entries, err := content.ReadDir("icons")
	if err != nil {
		return nil, err
	}

	decredIcons := make(map[string]image.Image)
	for _, entry := range entries {

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}

		imgBytes, err := content.ReadFile("icons/" + entry.Name())
		if err != nil {
			return nil, err
		}

		img, _, err := image.Decode(bytes.NewBuffer(imgBytes))
		if err != nil {
			return nil, err
		}

		split := strings.Split(entry.Name(), ".")
		decredIcons[split[0]] = img
	}

	return decredIcons, nil
}
