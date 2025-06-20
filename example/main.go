package main

import (
	"errors"
	"net/http"

	"github.com/linktomarkdown/htxp"
)

func main() {
	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"name": "张三",
			"age":  30,
		}
		htxp.Success(w, data)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("发生了一个错误")
		htxp.Error(w, err)
	})

	http.ListenAndServe(":8080", nil)
}