package format

import (
	"errors"
	"fmt"
	"strings"

	"github.com/princjef/gomarkdoc/lang"
)

// Asciidoc provides a Format which is compatible asciidoc format specification.
type Asciidoc struct{}

// Bold converts the provided text to bold
func (f *Asciidoc) Bold(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	return fmt.Sprintf("*%s*", text), nil
}

// CodeBlock wraps the provided code as a code block and tags it with the
// provided language (or no language if the empty string is provided).
func (f *Asciidoc) CodeBlock(language, code string) (string, error) {
	if code == "" {
		return "", nil
	}

	if "" == language {
		language = "go"
	}
	return fmt.Sprintf(`[source,%s]\n----\n%s\n----`, language, code), nil
}

// Header converts the provided text into a header of the provided level. The
// level is expected to be at least 1.
func (f *Asciidoc) Header(level int, text string) (string, error) {
	return f.header(level, escape(text))
}

// RawHeader converts the provided text into a header of the provided level
// without escaping the header text. The level is expected to be at least 1.
func (f *Asciidoc) RawHeader(level int, text string) (string, error) {
	return f.header(level, text)
}

// LocalHref generates an href for navigating to a header with the given
// headerText located within the same document as the href itself.
func (f *Asciidoc) LocalHref(headerText string) (string, error) {
	return fmt.Sprintf("xref:%s[%s]", f.genref(headerText), headerText), nil
}

// CodeHref always returns the empty string.
func (f *Asciidoc) CodeHref(loc lang.Location) (string, error) {
	return "", nil
}

// Link generates a link with the given text and href values.
func (f *Asciidoc) Link(text, href string) (string, error) {
	if text == "" {
		return "", nil
	}

	if href == "" {
		return text, nil
	}

	return fmt.Sprintf("%s[%s]", href, text), nil
}

// ListEntry generates an unordered list entry with the provided text at the
// provided zero-indexed depth. A depth of 0 is considered the topmost level of
// list.
func (f *Asciidoc) ListEntry(depth int, text string) (string, error) {
	if text == "" {
		return "", nil
	}

	prefix := strings.Repeat("**", depth)
	return fmt.Sprintf("%s %s\n", prefix, text), nil
}

// Accordion generates a collapsible content. Asciidoc handles collapsable
// with a titiel. The body is not escaped since asciidoc can handle all types
// of elements within a collapsable.
func (f *Asciidoc) Accordion(title, body string) (string, error) {
	if "" == title {
		title = "Description"
	}

	return fmt.Sprintf(".%s\n[%%collapsible]\n====\n%s\n====\n", title, body), nil

}

// AccordionHeader generates the header visible when an accordion is collapsed.
//
// The AccordionHeader is expected to be used in conjunction with
// AccordionTerminator() when the demands of the body's rendering requires it to
// be generated independently. The result looks conceptually like the following:
//
//	accordion := format.AccordionHeader("Accordion Title") + "Accordion Body" + format.AccordionTerminator()
func (f *Asciidoc) AccordionHeader(title string) (string, error) {
	return fmt.Sprintf(".%s\nn[%%collapsible]\n====\n", title), nil
}

// AccordionTerminator generates the code necessary to terminate an accordion
// after the body. It is expected to be used in conjunction with
// AccordionHeader(). See AccordionHeader for a full description.
func (f *Asciidoc) AccordionTerminator() (string, error) {
	return "\n====\n", nil
}

// Paragraph formats a paragraph with the provided text as the contents.
func (f *Asciidoc) Paragraph(text string) (string, error) {
	return fmt.Sprintf("%s\n\n", text), nil
}

// Escape escapes special markdown characters from the provided text.
func (f *Asciidoc) Escape(text string) string {
	return escape(text)
}

func (f *Asciidoc) header(level int, text string) (string, error) {
	if level < 1 {
		return "", errors.New("format: header level cannot be less than 1")
	}

	switch level {
	case 1:
		return fmt.Sprintf("[[%s]]\n= %s\n\n", f.genref(text), text), nil
	case 2:
		return fmt.Sprintf("[[%s]]\n== %s\n\n", f.genref(text), text), nil
	case 3:
		return fmt.Sprintf("[[%s]]\n=== %s\n\n", f.genref(text), text), nil
	case 4:
		return fmt.Sprintf("[[%s]]\n==== %s\n\n", f.genref(text), text), nil
	case 5:
		return fmt.Sprintf("[[%s]]\n===== %s\n\n", f.genref(text), text), nil
	default:
		// Only go up to 6 levels. Anything higher is also level 6
		return fmt.Sprintf("[[%s]]\n====== %s\n\n", f.genref(text), text), nil
	}
}

func (f *Asciidoc) genref(text string) string {
	result := plainText(text)
	result = strings.ToLower(result)
	result = strings.TrimSpace(result)
	result = gfmWhitespaceRegex.ReplaceAllString(result, "-")
	result = gfmRemoveRegex.ReplaceAllString(result, "")

	return result
}
