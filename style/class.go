package style

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Selector struct {
	Name    string
	mutator *Mutator
}

type Class struct {
	Name      string
	mutator   *Mutator
	Selectors map[string]Selector
}

type Classes map[string]Class

func (c Class) Mutate(styles *Styles) {
	c.mutator.Mutate(styles)
}

func (c Class) MutateWithSelector(sel string, styles *Styles) {
	if sel, ok := c.Selectors[sel]; ok {
		sel.mutator.Mutate(styles)
		return
	}
	// if no special selector is found, decode with default
	c.mutator.Mutate(styles)
}

func (cs Classes) Mutate(styles *Styles, useClasses ...string) {
	for _, cn := range useClasses {
		if c, ok := cs[cn]; ok {
			c.Mutate(styles)
		}
	}
}

func (cs Classes) MutateWithSelector(sel string, styles *Styles, useClasses ...string) {
	for _, cn := range useClasses {
		if c, ok := cs[cn]; ok {
			c.MutateWithSelector(sel, styles)
		}
	}
}

func DecodeClasses(r io.Reader) (Classes, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read-all")
	}
	s := string(b)

	s = strings.Replace(s, "\r", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\t", " ", -1)
	cs := Classes{}
	pos := 0
	for {
		curr := s[pos:]
		i := strings.IndexByte(curr, '{')
		if i < 0 {
			return cs, nil
		}
		name := trimWS(curr[:i])
		if len(name) == 0 {
			return nil, errors.Errorf("style class without name")
		}
		in := strings.IndexByte(curr[i:], '}')
		if in < 0 {
			return nil, errors.Errorf("non matching brace")
		}
		in += i
		mut, err := DecodeMutator(bytes.NewBufferString(curr[i+1 : in]))
		if err != nil {
			return nil, errors.Wrap(err, "parse style")
		}
		className := name
		selName := ""
		idxSel := strings.Index(name, ":")
		if idxSel > 0 {
			className = name[:idxSel]
			selName = name[idxSel+1:]
			if cl, ok := cs[className]; ok {
				cl.Selectors[selName] = Selector{
					Name:    selName,
					mutator: mut,
				}
				fmt.Printf("added selector (%s) to class (%s) \n", selName, className)
			} else {
				return nil, errors.Errorf("no base class for (%s:%s)", className, selName)
			}
		} else {
			//no selector
			cl := Class{
				Name:      string(className),
				mutator:   mut,
				Selectors: map[string]Selector{},
			}
			cs[className] = cl
		}
		pos += in + 1
	}
}
