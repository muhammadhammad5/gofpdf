package gofpdf_test

import (
	"testing"

	"github.com/muhammadhammad5/gofpdf"
)

// Reproduces: panic: runtime error: index out of range [65533] with length 256
// when SplitText is called with UTF-8 characters not representable by 8-bit core fonts.
func TestSplitTextWithTranslatorDoesNotPanic(t *testing.T) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 12)

	// Guard against forks that don't have UnicodeTranslatorFromDescriptor.
	tr := pdf.UnicodeTranslatorFromDescriptor
	if tr == nil {
		t.Skip("fork does not expose UnicodeTranslatorFromDescriptor; skipping")
	}

	text := "this was — emoji 😊 will be replaced"
	safe := tr("")(text) // map UTF-8 → cp1252, replacing unsupported runes

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("did not expect panic with translator: %v", r)
		}
	}()

	_ = pdf.SplitText(safe, 180)
}
