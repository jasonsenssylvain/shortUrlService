package shorturllib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Router struct {
	Processors map[int]Processor
}

const (
	SHORT_URL    = 0
	ORIGINAL_URL = 1
)

func (t *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.RequestURI
	fmt.Println("curr uri is " + requestUrl)

	err := r.ParseForm()
	if err != nil {
		return
	}

	params := make(map[string]string)
	for k, v := range r.Form {
		params[k] = v[0]
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		return
	}

	fmt.Println("body is " + string(body))

	action := 0
	if r.Method == "GET" {
		action = 0
	} else {
		action = 1
	}

	processor, _ := t.Processors[action]
	err = processor.ProcessRequest(r.Method, requestUrl, params, body, w, r)
	if err != nil {
		fmt.Printf("[ERROR] : %v\n", err)
	}
	return
}
