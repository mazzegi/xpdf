package xdoc

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/mazzegi/xpdf/style"
	"github.com/pkg/errors"
)

type Styled struct {
	Mutators []*style.Mutator
	Classes  []string
}

func (i *Styled) DecodeAttrs(attrs []xml.Attr) error {
	for _, a := range attrs {
		if a.Name.Local == "style" {
			mut, err := style.DecodeMutator(bytes.NewBufferString(a.Value))
			if err != nil {
				return errors.Wrapf(err, "decode style applier (%s)", a.Value)
			}
			i.Mutators = append(i.Mutators, mut)
		} else if a.Name.Local == "class" {
			i.Classes = append(i.Classes, strings.Fields(a.Value)...)
		}
	}
	return nil
}

func (i *Styled) MutatedStyles(cs style.Classes, styles style.Styles) style.Styles {
	ms := styles
	cs.Mutate(&ms, i.Classes...)
	for _, mut := range i.Mutators {
		mut.Mutate(&ms)
	}
	return ms
}

func (i *Styled) MutatedStylesWithSelector(sel string, cs style.Classes, styles style.Styles) style.Styles {
	ms := styles
	cs.MutateWithSelector(sel, &ms, i.Classes...)
	for _, mut := range i.Mutators {
		mut.Mutate(&ms)
	}
	return ms
}

type NoStyles struct{}

func (i *NoStyles) DecodeAttrs(attrs []xml.Attr) error {
	return nil
}

func (i *NoStyles) MutatedStyles(cs style.Classes, styles style.Styles) style.Styles {
	return styles
}

func (i *NoStyles) MutatedStylesWithSelector(sel string, cs style.Classes, styles style.Styles) style.Styles {
	return styles
}
