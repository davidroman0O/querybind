package querybind

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type TestStruct struct {
	StringValue string  `querybind:"str"`
	IntValue    int     `querybind:"int"`
	BoolValue   bool    `querybind:"bool"`
	FloatValue  float64 `querybind:"float"`
	SliceValue  []int   `querybind:"slice"`
}

func TestBind_Success(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		result, err := Bind[TestStruct](c)
		if err != nil {
			t.Errorf("Bind returned an error: %v", err)
		}

		if result.StringValue != "hello" {
			t.Errorf("Expected StringValue to be 'hello', got '%s'", result.StringValue)
		}
		if result.IntValue != 42 {
			t.Errorf("Expected IntValue to be 42, got '%d'", result.IntValue)
		}
		if result.BoolValue != true {
			t.Errorf("Expected BoolValue to be true, got '%t'", result.BoolValue)
		}
		if result.FloatValue != 3.14 {
			t.Errorf("Expected FloatValue to be 3.14, got '%f'", result.FloatValue)
		}
		if len(result.SliceValue) != 3 || result.SliceValue[0] != 1 || result.SliceValue[1] != 2 || result.SliceValue[2] != 3 {
			t.Errorf("Expected SliceValue to be [1, 2, 3], got '%v'", result.SliceValue)
		}
		return nil
	})

	queryParams := url.Values{}
	queryParams.Set("str", "hello")
	queryParams.Set("int", "42")
	queryParams.Set("bool", "true")
	queryParams.Set("float", "3.14")
	queryParams.Set("slice", "1,2,3")

	req := httptest.NewRequest("GET", "/?"+queryParams.Encode(), nil)
	app.Test(req, -1)
}

func TestBind_MissingQueryParams(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		result, err := Bind[TestStruct](c)
		if err != nil {
			t.Errorf("Bind returned an error: %v", err)
		}

		if result.StringValue != "" {
			t.Errorf("Expected StringValue to be '', got '%s'", result.StringValue)
		}
		// Similar assertions for other fields
		return nil
	})

	req := httptest.NewRequest("GET", "/", nil)
	app.Test(req, -1)
}

func TestBind_InvalidQueryValues(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		_, err := Bind[TestStruct](c)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		return nil
	})

	queryParams := url.Values{}
	queryParams.Set("str", "hello")
	queryParams.Set("int", "invalid_int")
	queryParams.Set("bool", "invalid_bool")
	queryParams.Set("float", "invalid_float")
	queryParams.Set("slice", "1,invalid_int,3")

	req := httptest.NewRequest("GET", "/?"+queryParams.Encode(), nil)
	app.Test(req, -1)
}

// TestBind_EmptyQueryParams tests binding when query parameters are empty.
func TestBind_EmptyQueryParams(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		result, err := Bind[TestStruct](c)
		if err != nil {
			t.Errorf("Bind returned an error: %v", err)
		}

		if result.StringValue != "" || result.IntValue != 0 || result.BoolValue != false || result.FloatValue != 0.0 || result.SliceValue != nil {
			t.Errorf("Expected all fields to be their zero values, got %+v", result)
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/?", nil)
	app.Test(req, -1)
}

// TestBind_UnsupportedTypes tests binding to a struct with unsupported field types.
func TestBind_UnsupportedTypes(t *testing.T) {
	type UnsupportedStruct struct {
		UnsupportedField map[string]string `querybind:"unsupported"`
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		_, err := Bind[UnsupportedStruct](c)
		if err == nil {
			t.Errorf("Expected an error for unsupported field type, got nil")
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/?unsupported=key:value", nil)
	app.Test(req, -1)
}

// TestBind_StructPointer tests binding to a struct pointer.
func TestBind_StructPointer(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		result, err := Bind[*TestStruct](c)
		if err != nil {
			t.Errorf("Bind returned an error: %v", err)
		}

		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		} else {
			// Dereferencing the result pointer to access the struct fields
			dereferencedResult := *result
			if dereferencedResult.StringValue != "" || dereferencedResult.IntValue != 0 || dereferencedResult.BoolValue != false || dereferencedResult.FloatValue != 0.0 || dereferencedResult.SliceValue != nil {
				t.Errorf("Expected all fields to be their zero values, got %+v", dereferencedResult)
			}
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/?", nil)
	app.Test(req, -1)
}

// TestBind_InvalidQuerySyntax tests handling of query parameters with invalid syntax.
func TestBind_InvalidQuerySyntax(t *testing.T) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		_, err := Bind[TestStruct](c)
		if err == nil {
			t.Errorf("Expected an error for invalid query syntax, got nil")
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/?invalid%zz", nil)
	app.Test(req, -1)
}

// Additional test cases for specific scenarios can be added here...
