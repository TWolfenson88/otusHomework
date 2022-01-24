package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrRegexp = errors.New("value doesn't match regular expression")
	ErrLen    = errors.New("length doesn't match expected")
	ErrMax    = errors.New("value is larger than max")
	ErrMin    = errors.New("value is less than min")
	ErrIn     = errors.New("value out of range")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for _, val := range v {
		sb.WriteString(val.Field)
		sb.WriteString(": ")
		sb.WriteString(val.Err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func validateStringByTag(field string, value string, tags []string, vErr ValidationErrors) (ValidationErrors, error) {
	var err error
	for _, tag := range tags {
		tagValue := strings.Split(tag, ":")[1]
		switch {
		case strings.HasPrefix(tag, "len:"):
			vErr, err = lenValidate(field, value, tagValue, vErr)
		case strings.HasPrefix(tag, "in:"):
			vErr = inValidate(field, value, tagValue, vErr)
		case strings.HasPrefix(tag, "regexp:"):
			vErr, err = regexpValidate(field, value, tagValue, vErr)
		}
	}
	return vErr, err
}

func validateIntByTag(field string, value int, tags []string, vErr ValidationErrors) (ValidationErrors, error) {
	var err error
	for _, tag := range tags {
		tagValue := strings.Split(tag, ":")[1]
		switch {
		case strings.HasPrefix(tag, "in:"):
			i := strconv.Itoa(value)
			vErr = inValidate(field, i, tagValue, vErr)
		case strings.HasPrefix(tag, "min:"):
			vErr, err = minValidate(field, value, tagValue, vErr)
		case strings.HasPrefix(tag, "max:"):
			vErr, err = maxValidate(field, value, tagValue, vErr)
		}
	}
	return vErr, err
}

func lenValidate(fieldName string, value string, tagValue string, vErr ValidationErrors) (ValidationErrors, error) {
	var e ValidationError
	i, err := strconv.Atoi(tagValue)
	if err != nil {
		return vErr, fmt.Errorf("atoi error: %w", err)
	}
	if len(value) != i {
		e.Field = fieldName
		e.Err = ErrLen
		return append(vErr, e), nil
	}
	return vErr, nil
}

func inValidate(fieldName string, value string, tagValue string, vErr ValidationErrors) ValidationErrors {
	var e ValidationError
	dict := strings.Split(tagValue, ",")
	var ok bool
	for _, v := range dict {
		if v == value {
			ok = true
		}
	}
	if !ok {
		e.Field = fieldName
		e.Err = ErrIn
		return append(vErr, e)
	}
	return vErr
}

func regexpValidate(fieldName string, value string, tagValue string, vErr ValidationErrors) (ValidationErrors, error) {
	var e ValidationError
	matched, err := regexp.Match(tagValue, []byte(value))
	if err != nil {
		return vErr, fmt.Errorf("match error: %w", err)
	}
	if !matched {
		e.Field = fieldName
		e.Err = ErrRegexp
		return append(vErr, e), nil
	}
	return vErr, nil
}

func minValidate(fieldName string, value int, tagValue string, vErr ValidationErrors) (ValidationErrors, error) {
	var e ValidationError
	i, err := strconv.Atoi(tagValue)
	if err != nil {
		return vErr, fmt.Errorf("atoi error: %w", err)
	}
	if i > value {
		e.Field = fieldName
		e.Err = ErrMin
		return append(vErr, e), nil
	}
	return vErr, nil
}

func maxValidate(fieldName string, value int, tagValue string, vErr ValidationErrors) (ValidationErrors, error) {
	var e ValidationError
	i, err := strconv.Atoi(tagValue)
	if err != nil {
		return vErr, fmt.Errorf("atoi error: %w", err)
	}
	if i < value {
		e.Field = fieldName
		e.Err = ErrMax
		return append(vErr, e), nil
	}
	return vErr, nil
}

func typeSwitch(fieldName string, val interface{}, tags []string, vErr ValidationErrors) (ValidationErrors, error) {
	var err error
	switch h := val.(type) {
	case int:
		vErr, err = validateIntByTag(fieldName, h, tags, vErr)
	case string:
		vErr, err = validateStringByTag(fieldName, h, tags, vErr)
	case []string:
		for _, v := range h {
			vErr, err = validateStringByTag(fieldName, v, tags, vErr)
		}
	case []int:
		for _, v := range h {
			vErr, err = validateIntByTag(fieldName, v, tags, vErr)
		}
	}
	return vErr, err
}

func typeSelector(field string, value reflect.Value, tags []string, vErr ValidationErrors) (*ValidationErrors, error) {
	var err error
	switch {
	case value.Kind() == reflect.String:
		val := value.String()
		vErr, err = typeSwitch(field, val, tags, vErr)
	case value.Kind() == reflect.Int:
		val := int(value.Int())
		vErr, err = typeSwitch(field, val, tags, vErr)
	case value.Kind() == reflect.Int64:
		val := int(value.Int())
		vErr, err = typeSwitch(field, val, tags, vErr)
	case value.Kind() == reflect.Slice:
		vErr, err = typeSwitch(field, value.Interface(), tags, vErr)
	}
	return &vErr, err
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not a pointer to struct", v)
	}
	vType := val.Type()
	vErr := new(ValidationErrors)
	var err error
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		fv := val.Field(i)
		var tags []string
		tag := field.Tag.Get("validate")
		if strings.Contains(tag, "|") {
			tags = strings.Split(tag, "|")
		} else {
			tags = append(tags, tag)
		}
		if len(tag) != 0 {
			vErr, err = typeSelector(field.Name, fv, tags, *vErr)
			if err != nil {
				return err
			}
		}
	}
	return vErr
}
