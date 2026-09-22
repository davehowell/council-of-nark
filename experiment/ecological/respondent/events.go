package respondent

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Usage struct {
	Input      int       `json:"input"`
	Output     int       `json:"output"`
	CacheRead  int       `json:"cache_read"`
	CacheWrite int       `json:"cache_write"`
	Total      int       `json:"total_tokens"`
	Cost       UsageCost `json:"cost"`
}

type UsageCost struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheRead  float64 `json:"cache_read"`
	CacheWrite float64 `json:"cache_write"`
	Total      float64 `json:"total"`
}

type EventSummary struct {
	EventCount               int              `json:"event_count"`
	AgentStarts              int              `json:"agent_starts"`
	AgentEnds                int              `json:"agent_ends"`
	AgentSettled             int              `json:"agent_settled"`
	AssistantMessages        int              `json:"assistant_messages"`
	ToolExecutions           int              `json:"tool_executions"`
	ExtensionErrors          int              `json:"extension_errors"`
	AutomaticRetries         int              `json:"automatic_retries"`
	SubmissionCount          int              `json:"submission_count"`
	SubmissionWasFinalTool   bool             `json:"submission_was_final_tool"`
	AssistantAfterSubmission bool             `json:"assistant_after_submission"`
	Usage                    Usage            `json:"usage"`
	Review                   EcologicalReview `json:"-"`
	ReviewJSON               []byte           `json:"-"`
}

type AuditSummary struct {
	Records              int `json:"records"`
	SessionStarts        int `json:"session_starts"`
	ProviderRequests     int `json:"provider_requests"`
	ProviderResponses    int `json:"provider_responses"`
	AssistantUsage       int `json:"assistant_usage_records"`
	FinalSubmissions     int `json:"final_submissions"`
	ToolsAfterSubmission int `json:"tools_after_submission"`
	BudgetStops          int `json:"budget_stops"`
}

func AnalyzeEvents(path string, maxFinalBytes int) (EventSummary, error) {
	file, err := os.Open(path)
	if err != nil {
		return EventSummary{}, err
	}
	defer file.Close()
	summary := EventSummary{}
	lastTool := ""
	submissionSeen := false
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal(line, &event); err != nil {
			return summary, fmt.Errorf("parse Pi event %d: %w", summary.EventCount+1, err)
		}
		summary.EventCount++
		typeName, _ := event["type"].(string)
		switch typeName {
		case "agent_start":
			summary.AgentStarts++
		case "agent_end":
			summary.AgentEnds++
		case "agent_settled":
			summary.AgentSettled++
		case "extension_error":
			summary.ExtensionErrors++
		case "auto_retry_start":
			summary.AutomaticRetries++
		case "message_end":
			message, _ := event["message"].(map[string]any)
			if role, _ := message["role"].(string); role == "assistant" {
				if submissionSeen {
					summary.AssistantAfterSubmission = true
				}
				summary.AssistantMessages++
				addUsage(&summary.Usage, message["usage"])
			}
		case "tool_execution_end":
			summary.ToolExecutions++
			toolName, _ := event["toolName"].(string)
			lastTool = toolName
			if toolName != "submit_ecological_review" {
				continue
			}
			summary.SubmissionCount++
			submissionSeen = true
			if isError, _ := event["isError"].(bool); isError {
				continue
			}
			result, _ := event["result"].(map[string]any)
			details, ok := result["details"]
			if !ok {
				continue
			}
			data, err := json.Marshal(details)
			if err != nil {
				return summary, fmt.Errorf("encode final submission: %w", err)
			}
			review, err := DecodeReview(data, maxFinalBytes)
			if err != nil {
				return summary, err
			}
			summary.Review = review
			summary.ReviewJSON = data
		}
	}
	if err := scanner.Err(); err != nil {
		return summary, fmt.Errorf("read Pi events: %w", err)
	}
	summary.SubmissionWasFinalTool = summary.SubmissionCount == 1 && lastTool == "submit_ecological_review"
	if summary.EventCount == 0 || summary.AgentStarts != 1 || summary.AgentSettled != 1 {
		return summary, fmt.Errorf("Pi lifecycle is incomplete: starts=%d settled=%d events=%d", summary.AgentStarts, summary.AgentSettled, summary.EventCount)
	}
	if summary.ExtensionErrors != 0 {
		return summary, fmt.Errorf("Pi emitted %d extension errors", summary.ExtensionErrors)
	}
	if summary.SubmissionCount != 1 || len(summary.ReviewJSON) == 0 {
		return summary, fmt.Errorf("expected exactly one valid final submission; got %d", summary.SubmissionCount)
	}
	if !summary.SubmissionWasFinalTool || summary.AssistantAfterSubmission {
		return summary, fmt.Errorf("structured submission was not the final action")
	}
	return summary, nil
}

func AnalyzeAudit(path string) (AuditSummary, error) {
	file, err := os.Open(path)
	if err != nil {
		return AuditSummary{}, err
	}
	defer file.Close()
	summary := AuditSummary{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 32*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var record struct {
			SchemaVersion int             `json:"schema_version"`
			Sequence      int             `json:"sequence"`
			Kind          string          `json:"kind"`
			Payload       json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(line, &record); err != nil {
			return summary, fmt.Errorf("parse provider audit record %d: %w", summary.Records+1, err)
		}
		summary.Records++
		if record.SchemaVersion != 1 || record.Sequence != summary.Records || record.Kind == "" {
			return summary, fmt.Errorf("provider audit sequence %d is invalid", summary.Records)
		}
		switch record.Kind {
		case "session_start":
			summary.SessionStarts++
		case "provider_request":
			if len(record.Payload) == 0 || bytes.Equal(record.Payload, []byte("null")) {
				return summary, fmt.Errorf("provider request %d has no serialized payload", summary.ProviderRequests+1)
			}
			summary.ProviderRequests++
		case "provider_response":
			summary.ProviderResponses++
		case "assistant_usage":
			summary.AssistantUsage++
		case "final_submission":
			summary.FinalSubmissions++
		case "tool_after_submission":
			summary.ToolsAfterSubmission++
		case "budget_stop":
			summary.BudgetStops++
		}
	}
	return summary, scanner.Err()
}

func addUsage(total *Usage, raw any) {
	usage, _ := raw.(map[string]any)
	input := intNumber(usage["input"])
	output := intNumber(usage["output"])
	cacheRead := intNumber(usage["cacheRead"])
	cacheWrite := intNumber(usage["cacheWrite"])
	messageTotal := intNumber(usage["totalTokens"])
	if messageTotal == 0 {
		messageTotal = input + output + cacheRead + cacheWrite
	}
	total.Input += input
	total.Output += output
	total.CacheRead += cacheRead
	total.CacheWrite += cacheWrite
	total.Total += messageTotal
	cost, _ := usage["cost"].(map[string]any)
	total.Cost.Input += floatNumber(cost["input"])
	total.Cost.Output += floatNumber(cost["output"])
	total.Cost.CacheRead += floatNumber(cost["cacheRead"])
	total.Cost.CacheWrite += floatNumber(cost["cacheWrite"])
	total.Cost.Total += floatNumber(cost["total"])
}

func intNumber(value any) int { return int(floatNumber(value)) }

func floatNumber(value any) float64 {
	switch number := value.(type) {
	case float64:
		return number
	case json.Number:
		value, _ := number.Float64()
		return value
	default:
		return 0
	}
}
