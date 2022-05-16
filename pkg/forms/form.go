package forms

import (
	"regexp"
	"time"
	"strconv"
	"fmt"
	"unicode/utf8"
	"strings"
	"net/url"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=? ^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

func (f *Form) IntegerRange(field string, min int, max int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	intValue, err := strconv.ParseInt(value, 0, 32)
	if err != nil {
		f.Errors.Add(field, "This field must be a number")
	} else if  int(intValue) < min || int(intValue) > max {
		f.Errors.Add(field, fmt.Sprintf("This field must be in %d-%d range", min, max))

	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) ValidDate(field string, format string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if _, err := time.Parse(format, value); err != nil {
		f.Errors.Add(field, "This field is invalid")
	}
}

// TODO: add method for validating time

func (f* Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

