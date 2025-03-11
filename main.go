package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type item struct {
	name   string
	status bool
}

func main() {

	var todos [1024]item
	counter := 0

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.GET("/fetch-tasks", func(c echo.Context) error {
		fmt.Println("Fetching tasks...")
		if counter == 0 {
			return c.String(http.StatusOK, "")
		}

		var result string
		for i := 0; i < counter; i++ {
			if todos[i].name == "" {
				result += ""
			} else {
				result += fmt.Sprintf(`
				<tr id="task-%d">
					<td>%d</td>
					<td>%s</td>
					<td><input type="checkbox" %s hx-put="/todo/%d" hx-trigger="change" hx-target="#task-%d" hx-swap="outerHTML"></td>
					<td><button hx-delete="/todo/%d" hx-target="#task-%d" hx-swap="outerHTML">Delete</button></td>
				</tr>`, i, i+1, todos[i].name, checked(todos[i].status), i, i, i, i)
			}
		}
		return c.String(http.StatusOK, result)
	})

	e.POST("/todo/", func(c echo.Context) error {
		name := c.FormValue("name")
		if strings.Trim(name, " ") == "" {
			return c.String(http.StatusOK, "")
		}
		todos[counter] = item{name, false}
		result := fmt.Sprintf(`
			<tr id="task-%d">
				<td>%d</td>
				<td>%s</td>
				<td><input type="checkbox" %s hx-put="/todo/%d" hx-trigger="change" hx-target="#task-%d" hx-swap="outerHTML"></td>
				<td><button hx-delete="/todo/%d" hx-target="#task-%d" hx-swap="outerHTML">Delete</button></td>
			</tr>`, counter, (counter + 1), todos[counter].name, checked(todos[counter].status), counter, counter, counter, counter)
		counter++

		return c.String(http.StatusOK, result)
	})

	e.PUT("/todo/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(err)
		}
		if (id >= counter) || (todos[id].name == "") {
			return c.String(http.StatusBadRequest, "Does not exist")
		}

		todos[id].status = !todos[id].status

		result := fmt.Sprintf(`
			<tr id="task-%d">
				<td>%d</td>
				<td>%s</td>
				<td><input type="checkbox" id="status-%d" %s hx-put="/todo/%d" hx-trigger="change" hx-target="#task-%d" hx-swap="outerHTML"></td>
				<td><button hx-delete="/todo/%d" hx-target="#task-%d" hx-swap="outerHTML">Delete</button></td>
			</tr>`, id, (id + 1), todos[id].name, id, checked(todos[id].status), id, id, id, id)

		return c.String(http.StatusOK, result)
	})

	e.DELETE("/todo/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(err)
		}
		if id >= counter {
			return c.String(http.StatusBadRequest, "Does not exist")
		}
		todos[id].name = ""
		return c.String(http.StatusOK, "")
	})

	fmt.Println("Starting HTTP Echo Server on :8080")
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}

	/*e.GET("/todo/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(err)
		}
		if (id >= counter) || (todos[id].name == "") {
			return c.String(http.StatusBadRequest, "Does not exist")
		}
		result := fmt.Sprintf(`
		<tr id="task-%d">
			<td>%d</td>
			<td>%s</td>
			<td><input type="checkbox" %s hx-put="/todo/%d" hx-trigger="change" hx-target="#task-%d" hx-swap="outerHTML"></td>
			<td><button hx-delete="/todo/%d" hx-target="#task-%d" hx-swap="outerHTML">Delete</button></td>
		</tr>`, counter, (counter + 1), todos[counter].name, checked(todos[counter].status), counter, counter, counter, counter)
		return c.String(http.StatusOK, result)
	})*/
}

func checked(status bool) string {
	if status {
		return "checked"
	}
	return ""
}
