# Virtual Pet (AI Agent)

A virtual pet application powered by Google's Gemini AI that creates an interactive, intelligent pet companion.

## Features

- **AI-Powered Interactions**: Uses Gemini AI model for natural and contextual responses
- **Real-time Status**: Monitors pet's hunger, energy, and happiness levels
- **Interactive Web Interface**: Simple web UI for interacting with your pet
- **Multiple Actions**:
  - Feed your pet
  - Play games with your pet
  - Let your pet sleep
  - Chat with your pet

## Technical Stack

[GO Documentation](DOC.md)

- **Backend**: Go (Golang)
  - Fiber web framework for API endpoints
  - Google Generative AI SDK for Gemini
- **Frontend**: HTML/JavaScript
  - Single-page application
  - Real-time status updates
  - Interactive UI elements

## Setup

1. Clone the repository
2. Create a `.env` file with:

   ```sh
   GEMINI_API_KEY=your_api_key_here
   ```

3. Run the application:

   ```sh
   go run main.go
   ```

4. Access the web interface at `http://localhost:3000`

## Architecture

- `ai/`: Contains AI-related functionality
  - `client.go`: Gemini API client setup
  - `pet.go`: Pet logic and state management
- `api/`: HTTP server and route handlers
- `public/`: Frontend assets
  - `index.html`: Web interface

## Usage

The pet responds to various actions through the web interface:

- Feed: Decreases hunger and increases energy
- Play: Increases happiness but consumes energy
- Sleep: Restores energy but increases hunger
- Chat: Free-form interaction with your pet

The pet's responses are context-aware and take into account its current state.
