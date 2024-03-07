package main

import (
	"log"
	"magic-routing-context/share"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber"
)

func DoSth(req share.SaferRequest) share.SaferResponse {

	time.Sleep(time.Duration(1000))
	return share.SaferResponse{
		ID:      req.ID,
		Message: req.Message,
	}
}

type WrapperContext struct {
	Context fiber.Ctx
}

var m map[string]fiber.Ctx
var lock = sync.RWMutex{}

var kafka = make(chan share.SaferResponse)

func Handler() gin.Engine {
	time.Sleep(time.Duration(1000))
	fn := func(c *gin.Context) error {

		context := *c
		var req share.SaferRequest
		if err := context.BodyParser(&req); err != nil {
			return context.Status(400).SendString(err.Error())
		} else {
			log.Printf("Requelst: %+v \n Context: %+v", req, context)
			lock.Lock()
			m[req.ID] = context
			lock.Unlock()
			res := DoSth(req)
			time.Sleep(time.Duration(1000))
			kafka <- res
		}
		return nil
	}
	return fn
}

func Consume() {
	for {
		for res := range kafka {
			// res := <-kafka
			log.Printf("Res: %v", res)
			lock.Lock()
			g := m[res.ID]
			lock.Unlock()
			log.Printf("Response: %+v \n Context: %+v", res, g)
			//err := g.Status(200).JSON(res)
			// if err != nil {
			// 	log.Println(err)
			// }
			// delete m[res.ID]
			// delete(m, res.ID)
			g.Status(300).JSON(res)
		}
	}

}

func main() {
	m = make(map[string]gin.Engine)
	app := gin.Default()

	forever := make(chan bool)
	go func() {
		app.Post("/api", Handler())

		app.Listen(":3000")
	}()

	go Consume()
	<-forever

}

// package main

// import (
// 	"log"
// 	"magic-routing-context/share"

// 	"github.com/gofiber/fiber/v2"
// )

// func Handler() fiber.Handler {
// 	fn := func(c *fiber.Ctx) error {
// 		// how to print the address of c

// 		log.Println(c)
// 		c.Status(200).JSON(share.SaferResponse{
// 			ID:      "123",
// 			Message: "Hello",
// 		})
// 		return nil
// 	}
// 	return fn
// }

// func main() {
// 	app := fiber.New()

// 	app.Post("/", Handler())

// 	app.Listen(":3000")

// }
