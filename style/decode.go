package style

import (
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func trimWS(s string) string {
	return strings.Trim(s, " \r\n\t")
}

func parseRaw(s string) (map[string]string, error) {
	raw := map[string]string{}
	styleStrs := strings.Split(s, ";")
	for _, styleStr := range styleStrs {
		styleStr = trimWS(styleStr)
		if len(styleStr) == 0 {
			continue
		}
		styleKV := strings.Split(styleStr, ":")
		if len(styleKV) != 2 {
			return nil, errors.Errorf("invalid style syntax (%s) must be of (key:val)", styleStr)
		}
		raw[trimWS(styleKV[0])] = trimWS(styleKV[1])
	}
	return raw, nil
}

type Unmarshaler interface {
	UnmarshalStyle(v string) error
}

type MutateFnc func(styles *Styles)

func MutateNone(styles *Styles) {}

type Mutator struct {
	fncs []MutateFnc
}

func (m *Mutator) Append(other *Mutator) {
	for _, f := range other.fncs {
		m.fncs = append(m.fncs, f)
	}
}

func (m *Mutator) Mutate(styles *Styles) {
	for _, fnc := range m.fncs {
		fnc(styles)
	}
}

func DecodeMutator(r io.Reader) (*Mutator, error) {
	m := &Mutator{
		fncs: []MutateFnc{},
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read-all")
	}
	raw, err := parseRaw(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "parse-raw")
	}
	protoType := reflect.TypeOf(Styles{})
	for k, v := range raw {
		fnc, found, err := makeMutateFnc(protoType, k, v, []int{})
		if err != nil {
			return nil, errors.Wrapf(err, "make-mutate-fnc (%s, %s)", k, v)
		}
		if !found {
			continue
		}
		m.fncs = append(m.fncs, fnc)
	}
	return m, nil
}

func makeMutateFnc(rt reflect.Type, key, val string, indexPath []int) (MutateFnc, bool, error) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		currIndexPath := appendedCopy(indexPath, i)
		if t := field.Tag.Get("style"); t == key {
			var setValue func(reflect.Value)
			um, impls := reflect.New(field.Type).Interface().(Unmarshaler)
			if impls {
				err := um.UnmarshalStyle(val)
				if err != nil {
					return nil, false, err
				}
				setValue = func(rv reflect.Value) {
					rv.Set(reflect.ValueOf(um).Elem())
				}
			} else {
				var err error
				setValue, err = makeSetValueFnc(field.Type.Kind(), val)
				if err != nil {
					return nil, false, errors.Wrapf(err, "make set value func (%s, %s)", key, val)
				}
			}

			return func(s *Styles) {
				rVal := reflect.ValueOf(s).Elem()
				for _, fIdx := range currIndexPath {
					rVal = rVal.Field(fIdx)
				}
				setValue(rVal)
			}, true, nil
		}

		if field.Type.Kind() == reflect.Struct {
			fnc, found, err := makeMutateFnc(field.Type, key, val, currIndexPath)
			if err != nil {
				return nil, true, err
			} else if found {
				return fnc, true, nil
			}
		}
	}
	return nil, false, nil
}

func makeSetValueFnc(kind reflect.Kind, styleValue string) (func(v reflect.Value), error) {
	switch kind {
	case reflect.String:
		return func(v reflect.Value) {
			v.SetString(styleValue)
		}, nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(styleValue, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-float (%s)", styleValue)
		}
		return func(v reflect.Value) {
			v.SetFloat(f)
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(styleValue, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-int (%s)", styleValue)
		}
		return func(v reflect.Value) {
			v.SetInt(n)
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(styleValue, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-uint (%s)", styleValue)
		}
		return func(v reflect.Value) {
			v.SetUint(n)
		}, nil
	default:
		return nil, errors.Errorf("unsupported style kind (%s)", kind)
	}
}

func appendedCopy(sl []int, a int) []int {
	c := make([]int, len(sl)+1)
	for i, v := range sl {
		c[i] = v
	}
	c[len(c)-1] = a
	return c
}
