package aisandboxsdk

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

type ApiError struct {
    Status int
    Body   string
}

func (e *ApiError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.Status, e.Body)
}

type Client struct {
    BaseURL    string
    APIKey     string
    HTTPClient *http.Client
    MaxRetries int
}

func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        BaseURL: strings.TrimRight(baseURL, "/"),
        APIKey: apiKey,
        HTTPClient: &http.Client{Timeout: 30 * time.Second},
        MaxRetries: 2,
    }
}

func (c *Client) request(method, path string, body any, accept string) (*http.Response, error) {
    for attempt := 0; ; attempt++ {
        var payload io.Reader
        if body != nil {
            raw, err := json.Marshal(body)
            if err != nil {
                return nil, err
            }
            payload = bytes.NewReader(raw)
        }
        req, err := http.NewRequest(method, c.BaseURL+"/api/v1"+path, payload)
        if err != nil {
            return nil, err
        }
        req.Header.Set("Authorization", "Bearer "+c.APIKey)
        req.Header.Set("Accept", accept)
        if body != nil {
            req.Header.Set("Content-Type", "application/json")
        }

        resp, err := c.HTTPClient.Do(req)
        if err != nil {
            if attempt < c.MaxRetries {
                time.Sleep(time.Duration(200*(attempt+1)) * time.Millisecond)
                continue
            }
            return nil, err
        }
        if (resp.StatusCode == 429 || resp.StatusCode >= 500) && attempt < c.MaxRetries {
            _ = resp.Body.Close()
            time.Sleep(time.Duration(200*(attempt+1)) * time.Millisecond)
            continue
        }
        if resp.StatusCode >= 400 {
            raw, _ := io.ReadAll(resp.Body)
            _ = resp.Body.Close()
            return nil, &ApiError{Status: resp.StatusCode, Body: string(raw)}
        }
        return resp, nil
    }
}

func (c *Client) json(method, path string, body any) (map[string]any, error) {
    resp, err := c.request(method, path, body, "application/json")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    out := map[string]any{}
    err = json.NewDecoder(resp.Body).Decode(&out)
    return out, err
}

func (c *Client) CreateAgent(payload map[string]any) (map[string]any, error) { return c.json(http.MethodPost, "/agents", payload) }
func (c *Client) UpdateAgent(agentID int, payload map[string]any) (map[string]any, error) { return c.json(http.MethodPut, fmt.Sprintf("/agents/%d", agentID), payload) }
func (c *Client) ListAgents(query string) (map[string]any, error) {
    suffix := ""
    if query != "" {
        suffix = "?" + query
    }
    return c.json(http.MethodGet, "/agents"+suffix, nil)
}
func (c *Client) GetAgentDetail(agentID int) (map[string]any, error) { return c.json(http.MethodGet, fmt.Sprintf("/agents/%d", agentID), nil) }
func (c *Client) CreateChatChannel(agentID int, payload map[string]any) (map[string]any, error) {
    return c.json(http.MethodPost, fmt.Sprintf("/agents/%d/channels", agentID), payload)
}
func (c *Client) SendChat(agentID int, payload map[string]any) (map[string]any, error) {
    return c.json(http.MethodPost, fmt.Sprintf("/agents/%d/chat", agentID), payload)
}
func (c *Client) GetChatHistory(agentID int, channelID, query string) (map[string]any, error) {
    if query == "" {
        query = "limit=50&page=0"
    }
    return c.json(http.MethodGet, fmt.Sprintf("/agents/%d/channels/%s/messages?%s", agentID, channelID, query), nil)
}
func (c *Client) GetKnowledgeBaseArticles(agentID, collectionID int, query string) (map[string]any, error) {
    if query == "" {
        query = "startIndex=0"
    }
    return c.json(http.MethodGet, fmt.Sprintf("/agents/%d/collections/%d/articles?%s", agentID, collectionID, query), nil)
}
func (c *Client) CreateKnowledgeBase(agentID int, payload map[string]any) (map[string]any, error) {
    return c.json(http.MethodPost, fmt.Sprintf("/agents/%d/collections", agentID), payload)
}
func (c *Client) UpdateKnowledgeBaseStatus(agentCollectionID int, payload map[string]any) (map[string]any, error) {
    return c.json(http.MethodPatch, fmt.Sprintf("/agent-collections/%d/status", agentCollectionID), payload)
}
func (c *Client) ListKnowledgeBases(agentID int, query string) (map[string]any, error) {
    if query == "" {
        query = "activeOnly=false"
    }
    return c.json(http.MethodGet, fmt.Sprintf("/agents/%d/collections?%s", agentID, query), nil)
}

func (c *Client) SendChatStream(agentID int, payload map[string]any) ([]string, error) {
    resp, err := c.request(http.MethodPost, fmt.Sprintf("/agents/%d/chat", agentID), payload, "text/event-stream")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var chunks []string
    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        line := scanner.Text()
        if !strings.HasPrefix(line, "data: ") {
            continue
        }
        data := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
        if data == "[DONE]" {
            break
        }
        chunks = append(chunks, data)
    }
    return chunks, scanner.Err()
}
