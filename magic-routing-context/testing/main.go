package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"magic-routing-context/share"
	"math/rand"
	"net/http"
	"sync"
)

func main() {
	concurrent := 10
	var wg sync.WaitGroup

	// send post requests to localhost:7201/api/v1/moneytransfer and get response
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Specify the body of the POST
			num := fmt.Sprintf("%d", rand.Intn(1000))
			req := share.SaferRequest{
				ID:      num,
				Message: num,
			}
			bodyBytes, _ := json.Marshal(req)

			res, err := http.Post("http://localhost:3000/", "application/json",
				bytes.NewBuffer(bodyBytes),
			)

			if err != nil {
				log.Printf(err.Error())
			} else {
				// Read the response body in type share.Response
				var response share.SaferResponse
				body, _ := io.ReadAll(res.Body)
				json.Unmarshal(body, &response)
				log.Printf("Request: %v. Response: %v", req, response)

			}
		}(i)
	}

	wg.Wait()
}

// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"magic-routing-context/share"
// 	"math/rand"
// 	"net/http"

// 	"github.com/pborman/uuid"
// )

// func main() {
// 	concurrent := 5

// 	// send post requests to localhost:7201/api/v1/moneytransfer and get response
// 	for i := 0; i < concurrent; i++ {
// 		// Specify the body of the POST
// 		num := rand.Intn(1000)
// 		req := share.SaferRequest{
// 			ID:      uuid.New(),
// 			Message: fmt.Sprintf("%d", num),
// 		}
// 		bodyBytes, _ := json.Marshal(req)

// 		res, err := http.Post("http://localhost:3000/", "application/json",
// 			bytes.NewBuffer(bodyBytes),
// 		)

// 		if err != nil {
// 			log.Printf(err.Error())
// 		} else {
// 			// Read the response body in type share.Response
// 			var response share.SaferResponse
// 			body, _ := io.ReadAll(res.Body)
// 			json.Unmarshal(body, &response)
// 			log.Printf("Request: %v. Response: %v", req, response)

// 		}

// 	}

// }
