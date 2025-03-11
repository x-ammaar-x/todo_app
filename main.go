package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Todo struct {
	name   string
	status bool
}

func main() {

	var todos [1024]Todo
	counter := 0

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		fmt.Println("Handling index GET for page returning full page reload")
		//Need to handle the availability of hx-boost header.
		//basically if the page is already loaded and you are reloading or something calls hx-boost in the
		//hx-get attribute it should only return what is necessary.
		//need to hx-boost needs to included in the body to use htmx
		//so make sure the header is set, if not you always return the full page
		return c.File("index.html")
	})

	e.GET("/todo/all", func(c echo.Context) error {
		fmt.Println("GET /todo/all endpoint\n   returning all ")
		if counter == 0 {
			return c.String(http.StatusOK, "")
		}
		fmt.Println("Counter: %d", counter);
		
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return allTodos(todos).Render(c.Request().Context(), c.Response().Writer)
	})

	e.POST("/todo/", func(c echo.Context) error {
		fmt.Println("POST /todo/ endpoint\n   creating new todo")
		name := c.FormValue("name")
		if strings.Trim(name, " ") == "" {
			return c.String(http.StatusOK, "")
		}
		todos[counter] = Todo{name, false}
		todoComp := todo(counter, todos[counter])
		counter++
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return todoComp.Render(c.Request().Context(), c.Response().Writer)
	})

	e.PUT("/todo/:id", func(c echo.Context) error {
		fmt.Println("PUT /todo/:id endpoint\n   checking a todo")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(err)
		}
		if (id >= counter) || (todos[id].name == "") {
			return c.String(http.StatusBadRequest, "Does not exist")
		}

		todos[id].status = !todos[id].status
		todoComp := todo(id, todos[id])
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return todoComp.Render(c.Request().Context(), c.Response().Writer)

	})

	e.DELETE("/todo/:id", func(c echo.Context) error {
		fmt.Println("DELETE /todo/:id endpoint\n   delete a todo")
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
}
