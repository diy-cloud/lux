# Lux

```
╔╗   ╔╗ ╔╗╔═╗╔═╗
║║   ║║ ║║╚╗╚╝╔╝
║║   ║║ ║║ ╚╗╔╝ 
║║ ╔╗║║ ║║ ╔╝╚╗ 
║╚═╝║║╚═╝║╔╝╔╗╚╗
╚═══╝╚═══╝╚═╝╚═╝
                
                
```

Lux is a web framework for the Go programming language. It's designed to be convienent and easy to use, while still being fast runtime performance.

## Installation

### CLI Tool

To install the CLI tool, run the following command:

```bash
go install github.com/snowmerak/lux/v3@latest
```

### Library

```bash
go get github.com/snowmerak/lux/v3
```

## Usage

### Simple Example

1. Create a new project

```bash
go mod init myproject
```

2. Add `lux` to your project

```bash
go get github.com/snowmerak/lux/v3
```

3. Generate and edit sample api

```bash
lux g c -g api/index
```

This will generate a new controller in the `api` group called `index`.  
And edit `Route` in `metadata.controller.go`.

```go
// api/index/metadata.controller.go
package index

const (
	Route = "/index"
)
```

The below is the generated controller.

```go
// api/index/get.controller.go
package index

import (
	"github.com/snowmerak/lux/v3/context"
	"github.com/snowmerak/lux/v3/controller"
	"github.com/snowmerak/lux/v3/lux"
	"github.com/snowmerak/lux/v3/middleware"
)

type GetController struct{
	requestMiddlewares []middleware.Request
	responseMiddlewares []middleware.Response
	handler controller.RestHandler
}

func NewGetController() *GetController {
	return &GetController{
		requestMiddlewares: []middleware.Request{},
		responseMiddlewares: []middleware.Response{},
		handler: func(lc *context.LuxContext) error {
			// Write your handler here
			return lc.ReplyString("Hello, World!")
		},
	}
}

func RegisterGetController(c *GetController, l *lux.Lux) {
	l.AddRestController(Route, controller.GET, controller.RestController{
		RequestMiddlewares: c.requestMiddlewares,
		Handler: c.handler,
		ResponseMiddlewares: c.responseMiddlewares,
	})
}
```

4. Generate entrypoint

```bash
lux g e server
```

This will generate a new entrypoint called `server`.

5. Register API to entrypoint

```go
// server/main.go
package main

import (
	"context"
	"log"
	"playground/api/echo"
	"playground/api/index"
	"playground/client/redis"

	"github.com/snowmerak/lux/v3/lux"
	"github.com/snowmerak/lux/v3/provider"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	constructors := []any{
		lux.New,
		lux.GenerateListenAddress(":8080"),

		index.NewGetController,
	}

	updaters := []any{
		index.RegisterGetController,
	}

	p := provider.New()
	if err := p.Register(constructors...); err != nil {

		log.Fatal(err)
	}

	if err := p.Construct(ctx); err != nil {
		log.Fatal(err)
	}

	if err := provider.Update(p, updaters...); err != nil {
		log.Fatal(err)
	}

	if err := provider.JustRun(p, lux.ListenAndServe1); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
```

6. Run the server

```bash
go run server/main.go
```

### With Service

7. Generate and edit service

```bash
lux g s client/redis 
```

This will generate a new service in the `client` group called `redis`.  
And add `Get` method to `RedisService` in `service.go`.

```go
// client/redis/service.go
package redis

type RedisService struct{}

func NewService() *RedisService {
	return &RedisService{}
}

// Added method
func (r *RedisService) Get() string {
	return "Hello, Redis!"
}
```

8. Register service to entrypoint

Service should be registered before controller.  
You can use service in any controller via dependency injection.

```go
// server/main.go
package main

import (
	"context"
	"log"
	"playground/api/echo"
	"playground/api/index"
	"playground/client/redis"

	"github.com/snowmerak/lux/v3/lux"
	"github.com/snowmerak/lux/v3/provider"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	constructors := []any{
		lux.New,
		lux.GenerateListenAddress(":8080"),

		index.NewGetController,

		redis.NewService,
	}

	updaters := []any{
		index.RegisterGetController,
	}

	p := provider.New()
	if err := p.Register(constructors...); err != nil {

		log.Fatal(err)
	}

	if err := p.Construct(ctx); err != nil {
		log.Fatal(err)
	}

	if err := provider.Update(p, updaters...); err != nil {
		log.Fatal(err)
	}

	if err := provider.JustRun(p, lux.ListenAndServe1); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
```

9. Use service in controller

```go
// api/index/get.controller.go
package index

import (
	"playground/client/redis"

	"github.com/snowmerak/lux/v3/context"
	"github.com/snowmerak/lux/v3/controller"
	"github.com/snowmerak/lux/v3/lux"
	"github.com/snowmerak/lux/v3/middleware"
)

type GetController struct {
	requestMiddlewares  []middleware.Request
	responseMiddlewares []middleware.Response
	handler             controller.RestHandler
}

func NewGetController(redisService *redis.RedisService) *GetController {
	return &GetController{
		requestMiddlewares:  []middleware.Request{},
		responseMiddlewares: []middleware.Response{},
		handler: func(lc *context.LuxContext) error {
			message := redisService.Get()
			return lc.ReplyString(message)
		},
	}
}

func RegisterGetController(c *GetController, l *lux.Lux) {
	l.AddRestController(Route, controller.GET, controller.RestController{
		RequestMiddlewares:  c.requestMiddlewares,
		Handler:             c.handler,
		ResponseMiddlewares: c.responseMiddlewares,
	})
}
```

10. Run the server

```bash
go run server/main.go
```
