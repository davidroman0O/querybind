# QueryBind 

The `querybind` package is a utility library for the [Fiber](https://gofiber.io/) web framework, designed to simplify query parameter binding and state management in HTMX-driven applications. It provides a way to bind query parameters to Go structs and manage URL state seamlessly.

## Features

- Bind URL query parameters to Go structs dynamically using struct tags.
- Maintain and manipulate browser URL state with HTMX requests without full page reloads.
- Provide functional options to customize the behavior of the response binding.

## Installation

To use the `querybind` package in your project, run:

```bash
go get github.com/davidroman0O/querybind
```

## Usage

### QueryBind

Bind query parameters from the request to a struct:

```go
type FilterParams struct {
    Genres []string `querybind:"genres"`
    Years  []int    `querybind:"years"`
}

// ...

// return refreshed html of a component
app.Get("/page/component", func(c *fiber.Ctx) error {
    var params FilterParams
    if err := querybind.Bind(&params, c); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }
    // Use params here...
    // return component html
})
```

### ResponseBind

Update the URL state using `ResponseBind` by setting a `HX-Push-Url` response header:

```go
// return refreshed html of a component
app.Get("/page/component", func(c *fiber.Ctx) error {
    // Assume params is populated... and you did some processing on some data, whatever
    querybind.ResponseBind(c, params, querybind.WithPath("/page")) // for the component, you might want to keep the path of the page
    // Continue with response... that return the html
})
```

## Options

The `ResponseBind` function can be customized with the following option:

- `WithPath(path string)`: Customizes the path used in the `HX-Push-Url` header.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgements

- The [Fiber](https://gofiber.io/) team for creating a fantastic web framework.
- The creators and contributors of [HTMX](https://htmx.org/) for enabling modern, dynamic web interactions.

