package xdoc

import (
	"encoding/json"

	"github.com/mazzegi/xpdf/style"
	"github.com/pkg/errors"
)

type Diff struct {
	Path string
	Org  interface{}
	Mod  interface{}
}

func stylesDiff(org style.Styles, mod style.Styles) ([]Diff, error) {
	orgMSI := toMapStringInterface(org)
	modMSI := toMapStringInterface(mod)
	return interfaceDiff("styles", orgMSI, modMSI)
}

func toMapStringInterface(sty style.Styles) map[string]interface{} {
	bs, _ := json.Marshal(sty)
	msi := map[string]interface{}{}
	json.Unmarshal(bs, &msi)
	return msi
}

func interfaceDiff(path string, org interface{}, mod interface{}) ([]Diff, error) {
	diffs := []Diff{}
	switch org := org.(type) {
	case map[string]interface{}:
		for orgKey, orgVal := range org {
			modMSI, ok := mod.(map[string]interface{})
			if !ok {
				return nil, errors.Errorf("%s: org value is map-string-iface. mod is %T", path, mod)
			}
			modVal, ok := modMSI[orgKey]
			if !ok {
				return nil, errors.Errorf("%s (MSI): key %s does not exist in modified MSI", path, orgKey)
			}
			msiDiffs, err := interfaceDiff(path+":"+orgKey, orgVal, modVal)
			if err != nil {
				return nil, errors.Wrap(err, "interface diff")
			}
			diffs = append(diffs, msiDiffs...)
		}
	default:
		if org != mod {
			diffs = append(diffs, Diff{
				Path: path,
				Org:  org,
				Mod:  mod,
			})
		}
	}
	return diffs, nil
}
