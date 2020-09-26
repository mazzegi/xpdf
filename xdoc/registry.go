package xdoc

import (
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

var instructionRegistry *Registry

func init() {
	instructionRegistry = NewRegistry()
	instructionRegistry.Register(&Font{})
	instructionRegistry.Register(&Box{})
	instructionRegistry.Register(&Text{})
	instructionRegistry.Register(&LineFeed{})
	instructionRegistry.Register(&SetX{})
	instructionRegistry.Register(&SetY{})
	instructionRegistry.Register(&SetXY{})
	instructionRegistry.Register(&Image{})
	instructionRegistry.Register(&Table{})
	instructionRegistry.Register(&TableRow{})
	instructionRegistry.Register(&TableCell{})
}

type Registry struct {
	types map[string]Instruction
}

func NewRegistry() *Registry {
	return &Registry{
		types: map[string]Instruction{},
	}
}

func (r *Registry) Register(prototype Instruction) error {
	ty := reflect.TypeOf(prototype)
	if ty.Kind() != reflect.Ptr {
		return errors.Errorf("register (%T). Instruction must be a ptr type (kind is %s).", ty.Name(), ty.Kind())
	}
	ty = reflect.TypeOf(reflect.ValueOf(prototype).Elem().Interface())
	fxml, ok := ty.FieldByName("XMLName")
	if !ok {
		return errors.Errorf("(%T) contains no XMLName", ty.Name())
	}
	xmlName := fxml.Tag.Get("xml")
	r.types[xmlName] = prototype
	return nil
}

func (r *Registry) Decode(d *xml.Decoder, start xml.StartElement) (Instruction, error) {
	proto, contains := r.types[start.Name.Local]
	if !contains {
		return nil, fmt.Errorf("registry-decode: (%s) is not registered", start.Name.Local)
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
