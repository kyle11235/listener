package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/kyleoracle/listener/config"
	"github.com/labstack/echo" // https://echo.labstack.com/guide
	"github.com/labstack/echo/middleware"
	"github.com/tidwall/gjson" // https://programming.vip/docs/super-easy-to-parse-json-package-gjson-with-golang.html
)

var url string

func main() {

	// print config
	url = config.Config["url"]
	fmt.Printf("sending out url=%s\n", url)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		// fmt.Printf("body=\n\n%s\n", reqBody)

		tx := gjson.GetBytes(reqBody, `eventMsg.0.payload.action.proposalResponsePayload.extension.results.nsRwset.#(namespace=="orderchain").rwset.writes.0.value.data`)
		fmt.Printf("tx=%s\n\n", tx)

		txBytes := []byte(tx.Raw)                                          // string -> bytes = []byte(tx.Raw), bytes -> string = string(data[:])
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(txBytes)) // interface io.Reader, read bytes
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status=", resp.Status)

	}))

	// Routes
	e.POST("/", fn)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}

func fn(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
