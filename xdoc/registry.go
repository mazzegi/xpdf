package xdoc

import (
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

var registry *Registry

func init() {
	registry = NewRegistry()
	registry.RegisterInstruction(&Font{})
	registry.RegisterInstruction(&Box{})
	registry.RegisterInstruction(&Text{})
	registry.RegisterInstruction(&LineFeed{})
	registry.RegisterInstruction(&SetX{})
	registry.RegisterInstruction(&SetY{})
	registry.RegisterInstruction(&Image{})
	registry.RegisterInstruction(&Table{})
	registry.RegisterInstruction(&TableRow{})
	registry.RegisterInstruction(&TableCell{})

	registry.RegisterInstruction(&Grid{})
	registry.RegisterInstruction(&GridRow{})
	registry.RegisterInstruction(&GridPart{})

	registry.RegisterInstruction(&Paragraph{})
	registry.RegisterInstruction(&LineBreak{})
	registry.RegisterInstruction(&PageBreak{})
}

type Registry struct {
	instructions map[string]Instruction
}

func NewRegistry() *Registry {
	return &Registry{
		instructions: map[string]Instruction{},
	}
}

func xmlName(v any) (string, error) {
	ty := reflect.TypeOf(v)
	if ty.Kind() != reflect.Ptr {
		return "", errors.Errorf("register (%T). Instruction must be a ptr type (kind is %s).", ty.Name(), ty.Kind())
	}
	ty = reflect.TypeOf(reflect.ValueOf(v).Elem().Interface())
	fxml, ok := ty.FieldByName("XMLName")
	if !ok {
		return "", errors.Errorf("(%T) contains no XMLName", ty.Name())
	}
	xn := fxml.Tag.Get("xml")
	return xn, nil
}

func (r *Registry) RegisterInstruction(prototype Instruction) error {
	xn, err := xmlName(prototype)
	if err != nil {
		return fmt.Errorf("xml-name: %w", err)
	}
	r.instructions[xn] = prototype
	return nil
}

func (r *Registry) DecodeInstruction(d *xml.Decoder, start xml.StartElement) (Instruction, error) {
	proto, contains := r.instructions[start.Name.Local]
	if !contains {
		return nil, errors.Errorf("registry-decode: (%s) is NOT a registered instruction", start.Name.Local)
	}

	pointerToI := reflect.New(reflect.TypeOf(proto))
	err := d.DecodeElement(pointerToI.Interface(), &start)
	if err != nil {
		return nil, err
	}

	inst := pointerToI.Elem().Interface().(Instruction)
	err = inst.DecodeAttrs(start.Attr)
	if err != nil {
		return nil, err
	}
	return inst, nil
}
