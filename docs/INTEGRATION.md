# Integration Guide (Go)

This SDK is designed for low-change, low-touch customer integration.

## Goals
- Stable API surface for v1.
- Explicit timeout and retry controls.
- Streaming chat support (`text/event-stream`).

## Install
`go get github.com/eGroupAI/ai-sandbox-sdk-go`

## First Steps
1. Configure `base_url` and `api_key`.
2. Call `create_agent`.
3. Create chat channel and send the first message.
