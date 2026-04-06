package main

import (
    "fmt"
    "os"

    aisandboxsdk "github.com/eGroupAI/ai-sandbox-sdk-go"
)

func main() {
    client := aisandboxsdk.NewClient(
        getenv("AI_SANDBOX_BASE_URL", "https://www.egroupai.com"),
        os.Getenv("AI_SANDBOX_API_KEY"),
    )
    result, err := client.CreateAgent(map[string]any{
        "agentDisplayName": "Go SDK Quickstart",
        "agentDescription": "Created by Go SDK",
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v\n", result)
}

func getenv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
