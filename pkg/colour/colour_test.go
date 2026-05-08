package colour

import (
	"slices"
	"testing"
)

func TestColour_NamedColours(t *testing.T) {
	c := New()
	got := c.Colour("{red}error{reset}")
	want := Red + "error" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColour_AllDefaults(t *testing.T) {
	c := New()
	for placeholder, code := range defaultColours {
		got := c.Colour(placeholder)
		if got != code {
			t.Errorf("Colour(%q) = %q, want %q", placeholder, got, code)
		}
	}
}

func TestColour_NoPlaceholders(t *testing.T) {
	c := New()
	got := c.Colour("plain text")
	if got != "plain text" {
		t.Errorf("got %q, want %q", got, "plain text")
	}
}

func TestColour_MultiplePlaceholders(t *testing.T) {
	c := New()
	got := c.Colour("{bold}{red}error:{reset} something broke")
	want := Bold + Red + "error:" + Reset + " something broke"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColour_RGB(t *testing.T) {
	c := New()
	got := c.Colour("{rgb:255,100,0}orange{reset}")
	want := "\033[38;2;255;100;0m" + "orange" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColour_MultipleRGB(t *testing.T) {
	c := New()
	got := c.Colour("{rgb:255,0,0}red{rgb:0,255,0}green{reset}")
	want := "\033[38;2;255;0;0m" + "red" + "\033[38;2;0;255;0m" + "green" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColour_RGBZeroValues(t *testing.T) {
	c := New()
	got := c.Colour("{rgb:0,0,0}black{reset}")
	want := "\033[38;2;0;0;0m" + "black" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColour_Disabled(t *testing.T) {
	c := New(Disable())
	got := c.Colour("{red}error{reset}")
	if got != "error" {
		t.Errorf("got %q, want %q", got, "error")
	}
}

func TestColour_DisabledRGB(t *testing.T) {
	c := New(Disable())
	got := c.Colour("{rgb:255,100,0}orange{reset}")
	if got != "orange" {
		t.Errorf("got %q, want %q", got, "orange")
	}
}

func TestColour_DisabledPlainText(t *testing.T) {
	c := New(Disable())
	got := c.Colour("plain text")
	if got != "plain text" {
		t.Errorf("got %q, want %q", got, "plain text")
	}
}

func TestRGB(t *testing.T) {
	got := RGB(10, 20, 30)
	want := "\033[38;2;10;20;30m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWithColour_Custom(t *testing.T) {
	c := New(WithColour("{highlight}", Yellow+Bold))
	got := c.Colour("{highlight}text{reset}")
	want := Yellow + Bold + "text" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWithColours_Custom(t *testing.T) {
	c := New(WithColours(map[string]string{
		"{error}":   Red,
		"{success}": Green,
	}))
	got := c.Colour("{error}fail{reset} {success}pass{reset}")
	want := Red + "fail" + Reset + " " + Green + "pass" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWithColour_OverridesDefault(t *testing.T) {
	c := New(WithColour("{red}", Blue))
	got := c.Colour("{red}actually blue{reset}")
	want := Blue + "actually blue" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColourf(t *testing.T) {
	c := New()
	got := c.Colourf("{red}error: %s{reset}", "something broke")
	want := Red + "error: something broke" + Reset
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColourf_DoesNotColourArgs(t *testing.T) {
	c := New()
	got := c.Colourf("prefix %s suffix", "{red}not coloured{reset}")
	want := "prefix {red}not coloured{reset} suffix"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColourf_Disabled(t *testing.T) {
	c := New(Disable())
	got := c.Colourf("{red}count: %d{reset}", 42)
	want := "count: 42"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColour_InvalidRGBPanics(t *testing.T) {
	c := New()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic but did not get one")
		}
	}()
	c.Colour("{rgb:256,0,0}text{reset}")
}

func TestColour_InvalidBgRGBPanics(t *testing.T) {
	c := New()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic but did not get one")
		}
	}()
	c.Colour("{bg-rgb:0,999,0}text{reset}")
}

func TestNew_DefaultsPresent(t *testing.T) {
	c := New()
	colours := c.Colours
	for placeholder := range defaultColours {
		if !slices.Contains(mapKeys(colours), placeholder) {
			t.Errorf("missing default colour %q", placeholder)
		}
	}
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
