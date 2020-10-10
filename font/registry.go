package font

type Style string

const (
	Regular    Style = "regular"
	Bold       Style = "bold"
	Italic     Style = "italic"
	BoldItalic Style = "bold+italic"
)

type Descriptor struct {
	Name     string
	Style    Style
	FilePath string
}

type Registry struct {
	monoFont string
	fonts    []Descriptor
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (d *Registry) Register(fd Descriptor) {
	d.fonts = append(d.fonts, fd)
}

func (d *Registry) MonoFont() string {
	return d.monoFont
}

func (d *Registry) Each(do func(fd Descriptor) error) error {
	for _, fd := range d.fonts {
		err := do(fd)
		if err != nil {
			return err
		}
	}
	return nil
}
