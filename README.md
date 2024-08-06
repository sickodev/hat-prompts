# Generative AI Chat API

## Overview

This project implements a simple REST API using the Gin web framework in Go to generate responses using a generative AI model. The API handles POST requests with a prompt and returns generated text based on the prompt using Google's Generative AI service.

## Prerequisites

1. Go 1.16 or later
2. An API key for Google's Generative AI service
3. Git
4. Go modules

## Setup

### 1. Clone the Repository

```sh
git clone github.com/sickodev/hat-prompts.git
cd hat-prompts
```

### 2. Install Dependencies

```sh
go mod tidy
```

### 3. Environment Variables

Create a `.env` file in the root directory of the project and add the following environment variables:

```
API_KEY=your_google_generative_ai_api_key
PORT=:8080
```

### 4. Run the Application

```sh
go run main.go
```

The server will start on the port specified in the `.env` file.

## Endpoints

### POST /api/v1/generate

**Description**: Generates a response based on the provided prompt.

**Request Body**:

```json
{
  "prompt": "Your prompt here"
}
```

**Response**:

- `200 OK`: Successfully generated response.
  ```json
  {
    "response": "Generated response based on the prompt"
  }
  ```
- `400 Bad Request`: Invalid JSON.
  ```json
  {
    "error": "Invalid JSON"
  }
  ```
- `500 Internal Server Error`: Error from the generative AI service.
  ```json
  {
    "error": "Error message"
  }
  ```

## Project Structure

- `main.go`: The main entry point of the application.
- `.env`: Environment variables file.

## Dependencies

- [Gin Web Framework](https://github.com/gin-gonic/gin): HTTP web framework for Go.
- [godotenv](https://github.com/joho/godotenv): Loads environment variables from a `.env` file.
- [Google Generative AI Go Client](https://github.com/google/generative-ai-go): Client library for accessing Google's Generative AI service.

## Detailed Explanation

1. **Environment Setup**:
    - Load environment variables from the `.env` file using `godotenv`.

2. **Gin Router Setup**:
    - Create a new Gin router.
    - Define the `/api/v1/generate` endpoint.

3. **Handler Function** (`getResults`):
    - Read the `Prompt` from the request body.
    - Initialize the Generative AI client using the provided API key.
    - Configure the generative model settings including safety settings.
    - Send the prompt to the model and retrieve the generated response.
    - Return the response as JSON.

## Safety Settings

The model is configured with safety settings to block content categorized as harassment, hate speech, sexually explicit, and dangerous content. Only content with a high threshold for these categories is blocked.

