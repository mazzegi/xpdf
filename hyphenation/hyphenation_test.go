package hyphenation

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	s := "ab5o5liz"
	p, err := parsePattern(s)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	t.Logf("pattern: %q", p.String())

	s = ".me5ter"
	p, err = parsePattern(s)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	t.Logf("pattern: %q", p.String())
}

func TestHyhenation(t *testing.T) {
	hyp := NewEnUs()

	var s string
	var hsl []string
	var t0 time.Time

	// s = "hyphenation"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "concatenation"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "supercalifragilisticexpialidocious"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "Developer"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "sportsman"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "small"
	// t0 = time.Now()
	// hsl = hyp.Hyphenate(s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "sportsman"
	// t0 = time.Now()
	// hsl = hyp.Hyphenate(s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	s = longText
	t0 = time.Now()
	hsl = hyp.Hyphenate(s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))
}

var longText = `At brother inquiry of offices without do my service. As particular to companions at sentiments. Weather however luckily enquire so certain do. Aware did stood was day under ask. Dearest affixed enquire on explain opinion he. Reached who the mrs joy offices pleased. Towards did colonel article any parties. Article nor prepare chicken you him now. Shy merits say advice ten before lovers innate add. She cordially behaviour can attempted estimable. Trees delay fancy noise manor do as an small. Felicity now law securing breeding likewise extended and. Roused either who favour why ham. Knowledge nay estimable questions repulsive daughters boy. Solicitude gay way unaffected expression for. His mistress ladyship required off horrible disposed rejoiced. Unpleasing pianoforte unreserved as oh he unpleasant no inquietude insipidity. Advantages can discretion possession add favourable cultivated admiration far. Why rather assure how esteem end hunted nearer and before. By an truth after heard going early given he. Charmed to it excited females whether at examine. Him abilities suffering may are yet dependent. Mr do raising article general norland my hastily. Its companions say uncommonly pianoforte favourable. Education affection consulted by mr attending he therefore on forfeited. High way more far feet kind evil play led. Sometimes furnished collected add for resources attention. Norland an by minuter enquire it general on towards forming. Adapted mrs totally company two yet conduct men. So by colonel hearted ferrars. Draw from upon here gone add one. He in sportsman household otherwise it perceived instantly. Is inquiry no he several excited am. Called though excuse length ye needed it he having. Whatever throwing we on resolved entrance together graceful. Mrs assured add private married removed believe did she. Received the likewise law graceful his. Nor might set along charm now equal green. Pleased yet equally correct colonel not one. Say anxious carried compact conduct sex general nay certain. Mrs for recommend exquisite household eagerness preserved now. My improved honoured he am ecstatic quitting greatest formerly. His having within saw become ask passed misery giving. Recommend questions get too fulfilled. He fact in we case miss sake. Entrance be throwing he do blessing up. Hearts warmth in genius do garden advice mr it garret. Collected preserved are middleton dependent residence but him how. Handsome weddings yet mrs you has carriage packages. Preferred joy agreement put continual elsewhere delivered now. Mrs exercise felicity had men speaking met. Rich deal mrs part led pure will but. Perhaps far exposed age effects. Now distrusts you her delivered applauded affection out sincerity. As tolerably recommend shameless unfeeling he objection consisted. She although cheerful perceive screened throwing met not eat distance. Viewing hastily or written dearest elderly up weather it as. So direction so sweetness or extremity at daughters. Provided put unpacked now but bringing. Parish so enable innate in formed missed. Hand two was eat busy fail. Stand smart grave would in so. Be acceptance at precaution astonished excellence thoroughly is entreaties. Who decisively attachment has dispatched. Fruit defer in party me built under first. Forbade him but savings sending ham general. So play do in near park that pain.`

func TestLatinHyhenation(t *testing.T) {
	hyp := NewLatin()

	s := "Lorem"
	t0 := time.Now()
	hsl := hyp.Hyphenate(s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))
}
