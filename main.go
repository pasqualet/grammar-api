package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/joshlf13/grammar"
	"github.com/labstack/echo"
)

func getGrammar() *grammar.Grammar {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <grammar file>\n", os.Args[0])
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(2)
	}

	g, err := grammar.New(file)

	if err != nil {
		fmt.Printf("Error creating grammar: \n%v\n", err)
		os.Exit(3)
	}

	return g
}

func getServer(g *grammar.Grammar) *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		buf := new(bytes.Buffer)
		if err := g.Speak(buf); err != nil {
			e.Logger.Error("Error getting new sentence: %v\n", err)
			c.String(http.StatusInternalServerError, "")
		}

		return c.String(http.StatusOK, buf.String())
	})
	return e
}

func main() {
	g := getGrammar()
	e := getServer(g)
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":1323"
	}
	e.Logger.Fatal(e.Start(serverAddr))
}
