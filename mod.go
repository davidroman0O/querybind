package querybind

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Bind binds query parameters to a struct of type T based on `query` tags.
func Bind[T any](c *fiber.Ctx) (*T, error) {
	var t T
	val := reflect.ValueOf(&t).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := typ.Field(i).Tag.Get("querybind")
		if tag == "" {
			continue
		}

		paramValue := c.Query(tag)
		if paramValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.Slice:
			separator := ","
			sliceOf := field.Type().Elem()
			values := strings.Split(paramValue, separator)
			slice := reflect.MakeSlice(field.Type(), 0, len(values))
			for _, value := range values {
				elem := reflect.New(sliceOf).Elem()
				if err := setValueFromString(elem, value); err != nil {
					return nil, err
				}
				slice = reflect.Append(slice, elem)
			}
			field.Set(slice)
		default:
			if err := setValueFromString(field, paramValue); err != nil {
				return nil, err
			}
		}
	}
	return &t, nil
}

// setValueFromString sets a value from a string based on the value's type and returns an error if any.
func setValueFromString(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			v.SetInt(intValue)
		} else {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			v.SetUint(uintValue)
		} else {
			return err
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			v.SetBool(boolValue)
		} else {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			v.SetFloat(floatValue)
		} else {
			return err
		}
	default:
		return fmt.Errorf("unsupported kind %s", v.Kind())
	}
	return nil
}

// ResponseBind sets the HX-Push-Url header in the response based on the struct's `query` tags.
func ResponseBind[T any](c *fiber.Ctx, value T) {
	val := reflect.ValueOf(value)
	typ := reflect.TypeOf(value)

	params := url.Values{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("querybind")
		if tag == "" {
			continue
		}

		var stringValue string
		switch field.Kind() {
		case reflect.Slice:
			s := []string{}
			for j := 0; j < field.Len(); j++ {
				s = append(s, fmt.Sprint(field.Index(j)))
			}
			stringValue = strings.Join(s, ",")
		default:
			stringValue = fmt.Sprint(field.Interface())
		}

		if stringValue != "" {
			params.Add(tag, stringValue)
		}
	}

	fullURL := c.BaseURL() + c.Path() + "?" + params.Encode()
	c.Set("HX-Push-Url", fullURL)
}
