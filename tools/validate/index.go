package validate

import (
	"errors"
	"net/url"
	"regexp"
)

type Validator struct {
	Regex   *regexp.Regexp
	Message string
}

func validateForm(validatorMap map[string]*Validator, form url.Values, keys ...string) []error {
	errs := make([]error, 0)
	for _, key := range keys {
		if validator, ok := validatorMap[key]; ok {
			value := form.Get(key)
			if !validator.Regex.MatchString(value) {
				errs = append(errs, errors.New(validator.Message))
			}
		}
	}
	return errs
}
