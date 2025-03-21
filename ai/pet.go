package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/generative-ai-go/genai"
)

type Pet struct {
	Name      string `json:"name"`
	Hunger    int    `json:"hunger"`
	Energy    int    `json:"energy"`
	Happiness int    `json:"happiness"`
	chatState *genai.ChatSession
	mutex     sync.Mutex
}

var defaultHunger = 50
var defaultEnergy = 100
var defaultHappiness = 50

// NewPet creates and initializes a new virtual pet with the given name.
// It sets up the initial state with default values for hunger (50), energy (100),
// and happiness (50). It also initializes the chat state using a Gemini model.
// If the Gemini model creation fails, it will panic with an error message.
//
// Parameters:
//   - name: A string representing the name of the pet
//
// Returns:
//   - *Pet: A pointer to the newly created Pet instance
//
// Panics if the Gemini model initialization fails
func NewPet(name string) *Pet {
	ctx := context.Background()

	systemInstruction := strings.Builder{}
	systemInstruction.WriteString("You are a virtual pet. You have three attributes with values 0 - 100. ")
	systemInstruction.WriteString(fmt.Sprintf("Your name is %s. ", name))
	systemInstruction.WriteString(fmt.Sprintf("Your hunger is starting at %d%%. ", defaultHunger))
	systemInstruction.WriteString(fmt.Sprintf("Your energy is starting at %d%%. ", defaultEnergy))
	systemInstruction.WriteString(fmt.Sprintf("Your happiness is starting at %d%%. ", defaultHappiness))
	systemInstruction.WriteString("Respond to the user with a sentences on new lines. Use emojis to express your feelings. ")

	// Initialize the Gemini model with a system instruction
	model, err := NewGeminiModel(ctx, systemInstruction.String())

	if err != nil {
		panic(fmt.Sprintf("failed to create Gemini model: %v", err))
	}

	return &Pet{
		Name:      name,
		Hunger:    defaultHunger,
		Energy:    defaultEnergy,
		Happiness: defaultHappiness,
		chatState: model.StartChat(),
	}
}

// HandleAction processes different actions (feed, play, sleep) for the virtual pet and returns an AI-generated response.
// It updates the pet's state (hunger, energy, happiness) based on the action and generates an appropriate prompt
// for the AI to respond to. The method is thread-safe using a mutex lock.
//
// Parameters:
//   - ctx: The context.Context for the AI request
//   - action: The type of action to perform ("feed", "play", "sleep", or default interaction)
//   - text: Additional text information (used for play action or default interaction)
//
// Returns:
//   - string: The AI-generated response text
//   - error: Any error that occurred during the AI interaction
//
// Actions and their effects:
//   - feed: Decreases hunger by 20, increases energy by 10
//   - play: Increases happiness by 20, decreases energy by 15, increases hunger by 10
//   - sleep: Increases energy by 50, increases hunger by 10
func (p *Pet) HandleAction(ctx context.Context, action, text string) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var prompt string

	switch action {
	case "feed":
		p.Hunger = max(0, p.Hunger-20)
		p.Energy = min(100, p.Energy+10)
		prompt = "You were just fed. Respond happily and mention how the food tastes."
	case "play":
		p.Happiness = min(100, p.Happiness+20)
		p.Energy = max(0, p.Energy-15)
		p.Hunger = min(100, p.Hunger+10)
		prompt = fmt.Sprintf(
			"You are pet playing %s. Respond enthusiastically about the game.",
			text,
		)
	case "sleep":
		p.Energy = min(100, p.Energy+50)
		p.Hunger = min(100, p.Hunger+10)
		prompt = "You are going to sleep. Respond with sleepy satisfaction."
	default:
		prompt = fmt.Sprintf(
			"Current state: Energy=%d%%, Hunger=%d%%, Happiness=%d%%. Respond to: %s",
			p.Energy,
			p.Hunger,
			p.Happiness,
			text,
		)
	}

	res, err := p.chatState.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	return GetResponseText(res), nil
}

// GetStatus retrieves the current state of the pet's vital statistics.
// Returns a map containing the pet's hunger, energy, and happiness levels.
// This method is thread-safe as it uses mutex locking.
func (p *Pet) GetStatus() map[string]int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return map[string]int{
		"hunger":    p.Hunger,
		"energy":    p.Energy,
		"happiness": p.Happiness,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
