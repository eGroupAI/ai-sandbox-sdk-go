# Integration Guide (Go)

This SDK is designed for low-change, low-touch customer integration.

## Goals
- Stable API surface for v1.
- Explicit timeout and retry controls.
- Streaming chat support (`text/event-stream`).

## Retry safety
- **429 / 5xx** automatic retries apply only to **GET** and **HEAD**. **POST / PUT / PATCH** are not retried on those status codes to avoid duplicate side effects.
- **Transport** errors may still be retried for all methods, up to `MaxRetries`.

## Install
`go get github.com/eGroupAI/ai-sandbox-sdk-go`

## First Steps
1. Configure `BaseURL` and `APIKey` on `Client` (`NewClient(baseURL, apiKey)`).
2. Call `CreateAgent(...)`.
3. Create a chat channel with `CreateChatChannel(...)` and send the first message with `SendChat(...)` or `SendChatStream(...)`.
