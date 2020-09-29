package font

type Style string

const (
	Regular    Style = "regular"
	Bold       Style = "bold"
	Italic     Style = "italic"
	BoldItalic Style = "bold/italic"
)

type Descriptor struct {
	Name     string
	Style    Style
	FilePath string
}

type Directory struct {
	fonts []Descriptor
}

func NewDirectory() *Directory {
	return &Directory{}
}

func (d *Directory) Register(fd Descriptor) {
	d.fonts = append(d.fonts, fd)
}

func (d *Directory) Each(do func(fd Descriptor) error) error {
	for _, fd := range d.fonts {
		err := do(fd)
		if err != nil {
			return err
		}
	}
	return nil
}
