package api

import (
	"github.com/dblbee/govitualpet/ai"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ActionRequest struct {
	Action string `json:"action"`
	Text   string `json:"text"`
}

type ActionResponse struct {
	Response string `json:"response"`
}

type Route struct {
	Path    string
	Method  string
	Handler fiber.Handler
}

type ApiServer struct {
	routes []Route
	app    *fiber.App
	pet    *ai.Pet
}

// NewServer creates a new Fiber API server instance with the specified pet.
// It initializes a new Fiber application and sets up middleware for logging and CORS.
// The server is configured with routes for handling pet actions, status requests,
// and serving static files from the public directory.
//
// Parameters:
//   - pet: A pointer to the virtual pet instance
//
// Returns:
//   - *ApiServer: A pointer to the newly created API server
func NewServer(pet *ai.Pet) *ApiServer {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	server := &ApiServer{
		app: app,
		pet: pet,
	}

	server.AddActionRoute()
	server.AddStatusRoute()
	server.AddStaticFilesRoute()

	return server
}

// Start initializes and launches the API server on the specified port.
// It registers all defined routes with their respective HTTP methods and handlers.
// Each route in the server's route collection is mapped to the corresponding Fiber
// application method (GET, POST, PUT, DELETE, PATCH, OPTIONS).
// If a route's method is not recognized, it is registered to handle all HTTP methods.
//
// Parameters:
//   - port: A string representing the port number to listen on (e.g., ":3000")
//
// Returns:
//   - error: Returns any error encountered while starting the server
func (server *ApiServer) Start(port string) error {
	for _, route := range server.routes {
		switch route.Method {
		case fiber.MethodGet:
			server.app.Get(route.Path, route.Handler)
		case fiber.MethodPost:
			server.app.Post(route.Path, route.Handler)
		case fiber.MethodPut:
			server.app.Put(route.Path, route.Handler)
		case fiber.MethodDelete:
			server.app.Delete(route.Path, route.Handler)
		case fiber.MethodPatch:
			server.app.Patch(route.Path, route.Handler)
		case fiber.MethodOptions:
			server.app.Options(route.Path, route.Handler)
		default:
			server.app.All(route.Path, route.Handler)
		}
	}

	return server.app.Listen(port)
}

// Stop gracefully shuts down the API server by calling the underlying
// application's Shutdown method. This ensures all pending requests are
// completed and resources are properly released.
// Returns any error encountered during shutdown.
func (server *ApiServer) Stop() error {
	return server.app.Shutdown()
}

// AddStaticFilesRoute adds a route to serve the static index.html file from the public directory.
// It creates a new Route that handles GET requests to the root path "/" and responds
// by sending the contents of "./public/index.html".
// The route is appended to the server's routes slice.
func (server *ApiServer) AddStaticFilesRoute() {
	staticRoute := Route{
		"/",
		fiber.MethodGet,
		func(c *fiber.Ctx) error {
			return c.SendFile("./public/index.html")
		},
	}
	server.routes = append(server.routes, staticRoute)
}

// AddStatusRoute adds a GET endpoint at /api/pet-status that returns the pet's name and current status.
// The response is a JSON object containing:
//   - name: string - The name of the pet
//   - status: object - The current status of the pet
//
// This route is appended to the server's routes collection.
func (server *ApiServer) AddStatusRoute() {
	statusRoute := Route{
		"/api/pet-status",
		fiber.MethodGet,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"name":   server.pet.Name,
				"status": server.pet.GetStatus(),
			})
		},
	}

	server.routes = append(server.routes, statusRoute)
}

// AddActionRoute adds a new route to handle pet actions in the API server.
// It creates a POST endpoint at "/api/pet-action" that processes ActionRequest
// payloads containing an action and optional text. The route parses the JSON
// request body, forwards the action to the pet handler, and returns an
// ActionResponse with the pet's response. If JSON parsing fails, it returns
// a 400 Bad Request status code.
func (server *ApiServer) AddActionRoute() {
	actionRoute := Route{
		"/api/pet-action",
		fiber.MethodPost,
		func(c *fiber.Ctx) error {
			var req ActionRequest
			if err := c.BodyParser(&req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "cannot parse JSON",
				})
			}

			response, err := server.pet.HandleAction(c.Context(), req.Action, req.Text)
			if err != nil {
				return err
			}

			return c.JSON(ActionResponse{Response: response})
		},
	}
	server.routes = append(server.routes, actionRoute)
}
