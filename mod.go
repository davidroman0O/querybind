package querybind

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type ResponseBindParams struct {
	Path *string
}

type ResponseBindOption func(*ResponseBindParams)

// Bind binds query parameters to a struct of type T based on `query` tags.
func Bind[T any](ctx HTTPContext) (*T, error) {
	var t T
	val := reflect.ValueOf(&t).Elem()
	typ := val.Type()

	// Parse the Referer URL's query parameters if present.
	referer := ctx.Referer()
	var allQueryParams url.Values
	var err error
	if referer != "" {
		parsedURL, err := url.Parse(referer)
		if err != nil {
			return nil, err
		}
		allQueryParams = parsedURL.Query()
	} else {
		allQueryParams, err = url.ParseQuery(ctx.OriginalURL())
		if err != nil {
			return nil, err
		}
	}

	// Combine with the current HTMX request's query parameters
	ctx.QueryParams(&allQueryParams)

	// Bind the combined query parameters to the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := typ.Field(i).Tag.Get("querybind")
		if tag == "" {
			continue
		}

		// Set the field value from the combined query parameters
		paramValue := allQueryParams.Get(tag)
		if paramValue == "" {
			continue
		}

		if err := setValueFromString(field, paramValue); err != nil {
			return nil, err
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
	case reflect.Slice:
		// Assumes a comma-separated list for slice types
		elements := strings.Split(value, ",")
		slice := reflect.MakeSlice(v.Type(), len(elements), len(elements))
		for i, s := range elements {
			elem := slice.Index(i)
			if err := setValueFromString(elem, s); err != nil {
				return err
			}
		}
		v.Set(slice)
	default:
		return fmt.Errorf("unsupported kind %s for query binding", v.Kind())
	}
	return nil
}

func WithPath(path string) ResponseBindOption {
	return func(params *ResponseBindParams) {
		params.Path = &path
	}
}

// ResponseBind sets the HX-Push-Url header in the response based on the struct's `query` tags.
func ResponseBind[T any](ctx HTTPContext, value T, options ...ResponseBindOption) {
	params := ResponseBindParams{}
	for _, option := range options {
		option(&params)
	}

	// If no custom path is provided, use the request's path
	if params.Path == nil {
		defaultPath := ctx.Path()
		params.Path = &defaultPath
	}

	val := reflect.ValueOf(value)
	typ := reflect.TypeOf(value)

	queryParams := url.Values{}
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
			queryParams.Add(tag, stringValue)
		}
	}

	// Encode the query parameters manually to avoid encoding commas
	encodedQuery := encodeQueryParams(queryParams)

	fullURL := ctx.BaseUrl() + *params.Path + "?" + encodedQuery

	ctx.SetHeader("HX-Push-Url", fullURL)
}

// encodeQueryParams encodes the parameters without encoding commas.
func encodeQueryParams(params url.Values) string {
	query := ""
	for key, values := range params {
		for _, value := range values {
			if query != "" {
				query += "&"
			}
			query += url.QueryEscape(key) + "=" + url.QueryEscape(value)
		}
	}
	// Replace encoded commas back to literal commas
	query = strings.ReplaceAll(query, "%2C", ",")
	return query
}
