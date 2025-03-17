package ai

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GetResponseText processes a GenerateContentResponse from the Gemini API and extracts the text content.
// It handles responses with single or multiple parts, logging relevant metadata and content information.
// The function validates the response structure at each level and provides appropriate error handling.
//
// Parameters:
//   - resp: A pointer to genai.GenerateContentResponse containing the API response
//
// Returns:
//   - string: The extracted text content from the response. Returns empty string if no valid content is found.
//
// The function performs the following checks:
//   - Presence of candidates in the response
//   - Presence of content in the first candidate
//   - Presence of parts in the candidate content
//   - Handles both single and multiple part responses
//   - Processes text parts and logs non-text parts
func GetResponseText(resp *genai.GenerateContentResponse) string {
	var responseText string

	log.Printf("Cached Content Token Count: %d\n", resp.UsageMetadata.CachedContentTokenCount)
	log.Printf("Prompt Token Count: %d\n", resp.UsageMetadata.PromptTokenCount)
	log.Printf("Total Token Count: %d\n", resp.UsageMetadata.TotalTokenCount)

	log.Printf("Content: %s\n", resp.Candidates[0].Content)

	if len(resp.Candidates) == 0 {
		log.Println("No candidates found in the response.")
		return responseText
	}

	if resp.Candidates[0].Content == nil {
		log.Println("No content found in the first candidate.")
		return responseText
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("No parts found in the candidate content.")
		return responseText
	}

	if len(resp.Candidates[0].Content.Parts) > 1 {
		log.Println("Multiple parts found in the candidate content.")
		for _, part := range resp.Candidates[0].Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				responseText += string(textPart)
				responseText += "\n"
				log.Printf("Part: %s\n", textPart)
			} else {
				log.Printf("Non-text part found: %T\n", part)
				responseText += "Non-text part found."
			}
		}
		return responseText
	}

	if len(resp.Candidates[0].Content.Parts) == 1 {
		log.Println("Single part found in the candidate content.")

		for _, part := range resp.Candidates[0].Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				responseText = string(textPart)
				log.Printf("Part: %s\n", textPart)
			} else {
				log.Printf("Non-text part found: %T\n", part)
				responseText += "Non-text part found."
			}
		}

		return responseText
	}

	return responseText
}

// NewGeminiModel creates and configures a new Gemini generative AI model instance.
// It initializes the model with the API key from GEMINI_API_KEY environment variable
// and the model name from GEMINI_MODEL_NAME environment variable.
// The function sets up safety settings to block low and above thresholds for:
// dangerous content, hate speech, sexually explicit content, and harassment.
//
// Parameters:
//   - ctx: The context.Context for managing the client lifecycle
//
// Returns:
//   - *genai.GenerativeModel: Configured Gemini model instance
//   - error: Any error encountered during initialization
//
// Required Environment Variables:
//   - GEMINI_API_KEY: API key for Gemini service authentication
//   - GEMINI_MODEL_NAME: Name of the Gemini model to use
func NewGeminiModel(ctx context.Context) (*genai.GenerativeModel, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel(os.Getenv("GEMINI_MODEL_NAME"))

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockLowAndAbove,
		},
	}

	return model, nil
}
