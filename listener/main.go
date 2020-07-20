package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kyleoracle/listener/config"
	"github.com/kyleoracle/listener/model"

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
	}))

	// Routes
	e.POST("/", fn)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}

func fn(context echo.Context) error {
	bodyBytes, _ := ioutil.ReadAll(context.Request().Body)
	go post(bodyBytes)
	fmt.Println("\n\nnew tx, fn sends response")
	return context.String(http.StatusOK, "Hello, World!")
}

func post(bodyBytes []byte) {

	tx := gjson.GetBytes(bodyBytes, `eventMsg.0.payload.action.proposalResponsePayload.extension.results.nsRwset.#(namespace=="orderchain").rwset.writes.0.value.data`)
	fmt.Printf("tx=%s\n\n", tx)

	txBytes := []byte(tx.Raw)

	// data
	order := &model.OrderNetsuite{
		URI: "salesorder",
	}
	err := json.Unmarshal(txBytes, &order.Data)
	if err != nil {
		fmt.Println("Unmarshal err=", err.Error())
	}

	orderBytes, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Marshal err=", err.Error())
	}

	fmt.Printf("order bytes=%s\n\n", orderBytes)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(orderBytes)) // interface io.Reader, read bytes
	req.Header.Set("Content-Type", "application/json")

	// auth
	req.Header.Set("Authorization", config.Config["auth"])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// panic(err)
		fmt.Println("post err=", err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("post response Status=", resp.Status)
}
