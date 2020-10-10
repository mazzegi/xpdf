package font

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type StyleDefinition struct {
	Style    Style
	FilePath string
}

type Definition struct {
	Name             string
	StyleDefinitions []StyleDefinition
}

type Definitions struct {
	MonoFont string
	Fonts    []Definition
}

func LoadRegistryFromFile(file string) (*Registry, error) {
	var defs Definitions
	_, err := toml.DecodeFile(file, &defs)
	if err != nil {
		return nil, errors.Wrapf(err, "toml-decode-file %q", file)
	}

	//resolves filepaths relative to dir of definition file
	absFile, err := filepath.Abs(file)
	if err != nil {
		return nil, errors.Wrapf(err, "abs-file-path %q", file)
	}
	absDir := filepath.Dir(absFile)
	resolve := func(file string) string {
		if filepath.IsAbs(file) {
			return file
		}
		return filepath.Join(absDir, file)
	}

	reg := NewRegistry()
	reg.monoFont = defs.MonoFont
	for _, fnt := range defs.Fonts {
		for _, sty := range fnt.StyleDefinitions {
			reg.Register(Descriptor{
				Name:     fnt.Name,
				Style:    sty.Style,
				FilePath: resolve(sty.FilePath),
			})
		}
	}
	return reg, nil
}
