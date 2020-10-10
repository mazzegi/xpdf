package main

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/mazzegi/xpdf/font"
)

func main() {
	def := font.Definitions{
		MonoFont: "dejavu-mono",
		Fonts: []font.Definition{
			{
				Name: "dejavu",
				StyleDefinitions: []font.StyleDefinition{
					{
						Style:    font.Regular,
						FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
					},
					{
						Style:    font.Bold,
						FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
					},
				},
			},
			{
				Name: "dejavu-mono",
				StyleDefinitions: []font.StyleDefinition{
					{
						Style:    font.Regular,
						FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuMono.ttf",
					},
					{
						Style:    font.Bold,
						FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuMono-Bold.ttf",
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	toml.NewEncoder(&buf).Encode(def)
	fmt.Printf("%s\n", buf.String())
}
