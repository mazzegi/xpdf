package markup

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		items []Item
	}{
		{
			name: "basic",
			in:   "lorem ipsum dolor",
			items: []Item{
				&TextItem{Text: "lorem ipsum dolor", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
			},
		},
		{
			name: "style #1",
			in:   "lorem _ipsum_ dolor",
			items: []Item{
				&TextItem{Text: "lorem ", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "ipsum", Style: TextStyle{Italic: true, Bold: false, Mono: false}},
				&TextItem{Text: " dolor", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
			},
		},
		{
			name: "style #2",
			in:   "lorem __ipsum__ dolor",
			items: []Item{
				&TextItem{Text: "lorem ", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "ipsum", Style: TextStyle{Italic: false, Bold: true, Mono: false}},
				&TextItem{Text: " dolor", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
			},
		},
		{
			name: "style #3",
			in:   "lorem __ipsum *dolor*__",
			items: []Item{
				&TextItem{Text: "lorem ", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "ipsum ", Style: TextStyle{Italic: false, Bold: true, Mono: false}},
				&TextItem{Text: "dolor", Style: TextStyle{Italic: true, Bold: true, Mono: false}},
			},
		},
		{
			name: "style #4",
			in:   "_lorem_ ipsum sed __dolor__",
			items: []Item{
				&TextItem{Text: "lorem", Style: TextStyle{Italic: true, Bold: false, Mono: false}},
				&TextItem{Text: " ipsum sed ", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "dolor", Style: TextStyle{Italic: false, Bold: true, Mono: false}},
			},
		},
		{
			name: "style #5",
			in:   "lorem ipsum `sed` dolor",
			items: []Item{
				&TextItem{Text: "lorem ipsum ", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "sed", Style: TextStyle{Italic: false, Bold: false, Mono: true}},
				&TextItem{Text: " dolor", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
			},
		},
		{
			name: "style #6",
			in:   "lorem ipsum\\`sed` dolor",
			items: []Item{
				&TextItem{Text: "lorem ipsum", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&ControlItem{Op: LineFeed},
				&TextItem{Text: "sed", Style: TextStyle{Italic: false, Bold: false, Mono: true}},
				&TextItem{Text: " dolor", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			items := Parse(test.in)
			if err := ensureItemsEqual(items, test.items); err != nil {
				t.Fatalf("%v\nhave: %s\nwant: %s", err, dumpItems(items), dumpItems(test.items))
			}
		})
	}
}

func ensureItemsEqual(itemsHave []Item, itemsWant []Item) error {
	if len(itemsHave) != len(itemsWant) {
		return errors.Errorf("size of items is different (have: %d, want: %d)", len(itemsHave), len(itemsWant))
	}
	for i, haveItem := range itemsHave {
		wantItem := itemsWant[i]
		if reflect.TypeOf(haveItem) != reflect.TypeOf(wantItem) {
			return errors.Errorf("items at %d have different types (have: %T, want: %T)", i, haveItem, wantItem)
		}
		switch wantTypedItem := wantItem.(type) {
		case *ControlItem:
			haveTypedItem := haveItem.(*ControlItem)
			if *haveTypedItem != *wantTypedItem {
				return errors.Errorf("items at %d are not equal", i)
			}
		case *TextItem:
			haveTypedItem := haveItem.(*TextItem)
			if *haveTypedItem != *wantTypedItem {
				return errors.Errorf("items at %d are not equal", i)
			}
		}
	}
	return nil
}

func dumpItems(items []Item) string {
	sl := []string{}
	for _, i := range items {
		sl = append(sl, fmt.Sprintf("[%s]", i.String()))
	}
	return strings.Join(sl, ", ")
}

func TestWords(t *testing.T) {
	tests := []struct {
		in    string
		words []Item
	}{
		{
			in: "Duis autem vel eum iriure",
			words: []Item{
				&TextItem{Text: "Duis", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "autem", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "vel", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "eum", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "iriure", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
			},
		},
		{
			in: "Duis autem __vel__ eum iriure\\sed october esse ",
			words: []Item{
				&TextItem{Text: "Duis", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "autem", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "vel", Style: TextStyle{Italic: false, Bold: true, Mono: false}},
				&TextItem{Text: "eum", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "iriure", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&ControlItem{Op: LineFeed},
				&TextItem{Text: "sed", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "october", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
				&TextItem{Text: "esse", Style: TextStyle{Italic: false, Bold: false, Mono: false}},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-words-#%d", i), func(t *testing.T) {
			items := Parse(test.in).Words()
			if err := ensureItemsEqual(items, test.words); err != nil {
				t.Fatalf("%v\nhave: %s\nwant: %s", err, dumpItems(items), dumpItems(test.words))
			}
		})
	}
}
