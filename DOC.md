# Virtual Pet (GO Doc)

## Package Main

Package main implements a virtual pet simulation game where users can create, interact with, and care for digital pets.

The virtual pet system simulates various pet attributes including:

- Health and hunger levels
- Mood and happiness
- Energy and activity status

Users can perform various actions with their pets:

- Feed and care for their pet
- Play and interact with their pet
- Monitor their pet's status and well-being

## API Server

### NewServer

```go
func NewServer(pet *ai.Pet) *ApiServer
```

Creates a new Fiber API server instance with the specified pet. It initializes a new Fiber application and sets up middleware for logging and CORS. The server is configured with routes for handling pet actions, status requests, and serving static files from the public directory.

### Start

```go
func (server *ApiServer) Start(port string) error
```

Initializes and launches the API server on the specified port. It registers all defined routes with their respective HTTP methods and handlers. Each route in the server's route collection is mapped to the corresponding Fiber application method (GET, POST, PUT, DELETE, PATCH, OPTIONS).

### Stop

```go
func (server *ApiServer) Stop() error
```

Gracefully shuts down the API server by calling the underlying application's Shutdown method. This ensures all pending requests are completed and resources are properly released.

### AddStaticFilesRoute

```go
func (server *ApiServer) AddStaticFilesRoute()
```

Adds a route to serve the static index.html file from the public directory. It creates a new Route that handles GET requests to the root path "/" and responds by sending the contents of "./public/index.html".

### AddStatusRoute

```go
func (server *ApiServer) AddStatusRoute()
```

Adds a GET endpoint at /api/pet-status that returns the pet's name and current status. The response is a JSON object containing the pet's name and status.

### AddActionRoute

```go
func (server *ApiServer) AddActionRoute()
```

Adds a new route to handle pet actions in the API server. It creates a POST endpoint at "/api/pet-action" that processes ActionRequest payloads containing an action and optional text.

## Virtual Pet

### NewPet

```go
func NewPet(name string) *Pet
```

Creates and initializes a new virtual pet with the given name. It sets up the initial state with default values for hunger (50), energy (100), and happiness (50). It also initializes the chat state using a Gemini model.

### HandleAction

```go
func (p *Pet) HandleAction(ctx context.Context, action, text string) (string, error)
```

Processes different actions (feed, play, sleep) for the virtual pet and returns an AI-generated response. It updates the pet's state based on the action:

- feed: Decreases hunger by 20, increases energy by 10
- play: Increases happiness by 20, decreases energy by 15, increases hunger by 10
- sleep: Increases energy by 50, increases hunger by 10

### GetStatus

```go
func (p *Pet) GetStatus() map[string]int
```

Retrieves the current state of the pet's vital statistics. Returns a map containing the pet's hunger, energy, and happiness levels.

## AI Client

### GetResponseText

```go
func GetResponseText(resp *genai.GenerateContentResponse) string
```

Processes a GenerateContentResponse from the Gemini API and extracts the text content. It handles responses with single or multiple parts, logging relevant metadata and content information.

### NewGeminiModel

```go
func NewGeminiModel(ctx context.Context, systemInstruction string) (*genai.GenerativeModel, error)
```

Creates and configures a new Gemini generative AI model instance. It initializes the model with the API key and model name from environment variables, and sets up safety settings to block harmful content.

Required Environment Variables:

- GEMINI_API_KEY: API key for Gemini service authentication
- GEMINI_MODEL_NAME: Name of the Gemini model to use
