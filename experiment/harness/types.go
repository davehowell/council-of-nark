package harness

import "encoding/json"

type Provider struct {
	Adapter string `json:"adapter"`
	Model   string `json:"model"`
	Effort  string `json:"effort"`
}

type Config struct {
	SchemaVersion         int            `json:"schema_version"`
	Design                string         `json:"design"`
	Label                 string         `json:"label"`
	Seed                  string         `json:"seed"`
	Packets               []string       `json:"packets"`
	Arms                  []string       `json:"arms"`
	Repetitions           int            `json:"repetitions"`
	Provider              Provider       `json:"provider"`
	Providers             []Provider     `json:"providers"`
	Role                  string         `json:"role"`
	Roles                 []string       `json:"roles"`
	Wrappers              []string       `json:"wrappers"`
	Wrapper               string         `json:"wrapper"`
	Topologies            []string       `json:"topologies"`
	AllRoleOrders         bool           `json:"all_role_orders"`
	MaxConcurrency        int            `json:"max_concurrency"`
	RequestTimeoutSeconds int            `json:"request_timeout_seconds"`
	MaxAttempts           int            `json:"max_attempts"`
	RetryBackoffSeconds   float64        `json:"retry_backoff_seconds"`
	Decoding              map[string]any `json:"decoding"`
	Status                string         `json:"status"`
	Interpretation        string         `json:"interpretation"`
}

type Call struct {
	CallID     string         `json:"call_id"`
	ReviewerID string         `json:"reviewer_id"`
	Packet     string         `json:"packet"`
	Provider   Provider       `json:"provider"`
	PromptSpec map[string]any `json:"prompt_spec"`
	DependsOn  []string       `json:"depends_on"`
	Phase      string         `json:"phase"`
	Semantic   []any          `json:"semantic"`
}

type OutputSet struct {
	SetID       string         `json:"set_id"`
	Packet      string         `json:"packet"`
	CallIDs     []string       `json:"call_ids"`
	CostCallIDs []string       `json:"cost_call_ids"`
	Metadata    map[string]any `json:"metadata"`
	Semantic    []any          `json:"semantic"`
}

type Plan struct {
	SchemaVersion int            `json:"schema_version"`
	CreatedAt     string         `json:"created_at"`
	Design        string         `json:"design"`
	Seed          string         `json:"seed"`
	Calls         []Call         `json:"calls"`
	OutputSets    []OutputSet    `json:"output_sets"`
	Counts        map[string]int `json:"counts"`
}

type Asset struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type Freeze struct {
	SchemaVersion   int              `json:"schema_version"`
	CreatedAt       string           `json:"created_at"`
	SourceCommit    string           `json:"source_commit"`
	ConfigSource    string           `json:"config_source"`
	ConfigSHA256    string           `json:"config_sha256"`
	Assets          []Asset          `json:"assets"`
	ExternalRuntime []map[string]any `json:"external_runtime"`
	Isolation       map[string]any   `json:"isolation"`
	PublicAudit     map[string]any   `json:"public_audit"`
}

type Invocation struct {
	Command        []string
	CommandDisplay []string
	ReturnCode     int
	Stdout         string
	Stderr         string
	LatencySeconds float64
	Parsed         map[string]any
	ParseError     string
	Usage          map[string]any
	Sandbox        map[string]any
}

func providerMap(p Provider) map[string]any {
	data, _ := json.Marshal(p)
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	return m
}
