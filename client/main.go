package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"time"

	"go.elastic.co/apm/module/apmhttp"
)

var serverURL = "http://localhost:8080"

func serveSlowly(w http.ResponseWriter, req *http.Request) {
	// Wrap the HTTP client with apmhttp.WrapClient. When using the
	// wrapped client, any request whose context contains a transaction
	// will have a span reported.
	client := apmhttp.WrapClient(http.DefaultClient)
	slowReq, _ := http.NewRequest("GET", serverURL+"/slow", nil)
	errorReq, _ := http.NewRequest("GET", "http://testing.invalid", nil)

	// Propagate context with the outgoing request.
	resp, err := client.Do(slowReq.WithContext(req.Context()))
	if err != nil {
		log.Fatal(err)
	}

	// In the case where the request succeeds (i.e. no error
	// was returned above; unrelated to the HTTP status code),
	// the span is not ended until the body is consumed.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("response: %s\n", body)

	// Send a request to a URL with an unresolvable host. This
	// will cause the entire request to fail, immediately
	// ending the span.
	resp, err = client.Do(errorReq.WithContext(req.Context()))
	if err != nil {
		fmt.Println("error occurred")
	} else {
		resp.Body.Close()
	}

	// if len(spans) != 2 {
	// 	fmt.Println(len(spans), "spans")
	// } else {
	// 	for i, span := range spans {
	// 		const expectedFloor = 250 * time.Millisecond
	// 		if time.Duration(span.Duration*float64(time.Millisecond)) >= expectedFloor {
	// 			// This is the expected case (see output below). As noted
	// 			// previously, the span is only ended once the response body
	// 			// has been consumed (or closed).
	// 			fmt.Printf("span #%d duration >= %s\n", i+1, expectedFloor)
	// 		} else {
	// 			fmt.Printf("span #%d duration < %s\n", i+1, expectedFloor)
	// 		}
	// 	}
	// }

	w.WriteHeader(http.StatusTeapot)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	time.Sleep(250 * time.Millisecond)
	w.Write([]byte("*yawn*"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", serveSlowly)
	http.ListenAndServe(":3030", apmhttp.Wrap(mux))
}
