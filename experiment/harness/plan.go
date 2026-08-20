package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var allRoles = []string{"simplicity", "correctness", "architecture", "security", "cost", "entropy", "language"}

type builder struct {
	config Config
	calls  []Call
	sets   []OutputSet
}

func semanticParts(prefix []any, extra ...any) []any {
	out := append([]any{}, prefix...)
	return append(out, extra...)
}
func (b *builder) call(semantic []any, packet string, provider Provider, spec map[string]any, depends []string, phase string) Call {
	if phase == "" {
		phase = "review"
	}
	idParts := append([]any{"call"}, semantic...)
	reviewerParts := append([]any{"reviewer"}, semantic...)
	call := Call{CallID: "c-" + opaqueID(b.config.Seed, 16, idParts...), ReviewerID: "r-" + opaqueID(b.config.Seed, 10, reviewerParts...), Packet: packet, Provider: provider, PromptSpec: spec, DependsOn: append([]string{}, depends...), Phase: phase, Semantic: semantic}
	b.calls = append(b.calls, call)
	return call
}
func (b *builder) outputSet(semantic []any, packet string, callIDs, costIDs []string, metadata map[string]any) {
	parts := append([]any{"set"}, semantic...)
	b.sets = append(b.sets, OutputSet{SetID: "s-" + opaqueID(b.config.Seed, 16, parts...), Packet: packet, CallIDs: append([]string{}, callIDs...), CostCallIDs: append([]string{}, costIDs...), Metadata: metadata, Semantic: semantic})
}

func buildStageA(b *builder) error {
	c := b.config
	for _, packet := range c.Packets {
		for repeat := 0; repeat < c.Repetitions; repeat++ {
			for _, arm := range c.Arms {
				prefix := []any{"stage_a", packet, repeat, arm}
				if arm == "S0" || arm == "S1" || arm == "S2" {
					var spec map[string]any
					if arm == "S0" {
						spec = map[string]any{"kind": "generic"}
					} else if arm == "S1" {
						spec = map[string]any{"kind": "omnibus", "wrapper": "functional"}
					} else {
						spec = map[string]any{"kind": "omnibus", "wrapper": "fictional"}
					}
					review := b.call(semanticParts(prefix, "review"), packet, c.Provider, spec, nil, "")
					b.outputSet(semanticParts(prefix, "final"), packet, []string{review.CallID}, []string{review.CallID}, map[string]any{"design": "stage_a", "arm": arm, "repeat": repeat, "kind": "final"})
					continue
				}
				reviews := []Call{}
				for index, role := range allRoles {
					var spec map[string]any
					switch arm {
					case "M0":
						spec = map[string]any{"kind": "omnibus", "wrapper": "functional", "sample": index}
					case "M1":
						spec = map[string]any{"kind": "specialist", "wrapper": "functional", "role": role}
					case "M2":
						spec = map[string]any{"kind": "specialist", "wrapper": "fictional", "role": role}
					default:
						return fmt.Errorf("unknown Stage A arm %s", arm)
					}
					reviews = append(reviews, b.call(semanticParts(prefix, "review", index, role), packet, c.Provider, spec, nil, ""))
				}
				ids := callIDs(reviews)
				fused := b.call(semanticParts(prefix, "fuse"), packet, c.Provider, map[string]any{"kind": "fuser"}, ids, "fuse")
				b.outputSet(semanticParts(prefix, "raw_union"), packet, ids, ids, map[string]any{"design": "stage_a", "arm": arm, "repeat": repeat, "kind": "raw_union"})
				b.outputSet(semanticParts(prefix, "fused"), packet, []string{fused.CallID}, append(append([]string{}, ids...), fused.CallID), map[string]any{"design": "stage_a", "arm": arm, "repeat": repeat, "kind": "fused"})
			}
		}
	}
	return nil
}
func callIDs(calls []Call) []string {
	out := make([]string, len(calls))
	for i, c := range calls {
		out[i] = c.CallID
	}
	return out
}

func permutations(values []string) [][]string {
	if len(values) == 0 {
		return [][]string{{}}
	}
	out := [][]string{}
	for i, v := range values {
		rest := append([]string{}, values[:i]...)
		rest = append(rest, values[i+1:]...)
		for _, tail := range permutations(rest) {
			out = append(out, append([]string{v}, tail...))
		}
	}
	return out
}
func buildTopology(b *builder) error {
	c := b.config
	orders := [][]string{append([]string{}, c.Roles...)}
	if c.AllRoleOrders {
		orders = permutations(c.Roles)
	}
	for _, packet := range c.Packets {
		for repeat := 0; repeat < c.Repetitions; repeat++ {
			for _, order := range orders {
				orderName := strings.Join(order, "-")
				for _, topology := range c.Topologies {
					prefix := []any{"topology", packet, repeat, orderName, topology}
					reviews := []Call{}
					previous := []string{}
					for hop, role := range order {
						dynamic := topology == "chain" && hop > 0
						kind := "specialist"
						if dynamic {
							kind = "chain"
						}
						depends := []string{}
						if dynamic {
							depends = []string{previous[len(previous)-1]}
						}
						review := b.call(semanticParts(prefix, "review", hop, role), packet, c.Provider, map[string]any{"kind": kind, "wrapper": c.Wrapper, "role": role, "hop": hop + 1}, depends, "")
						reviews = append(reviews, review)
						previous = append(previous, review.CallID)
						if topology == "chain" {
							b.outputSet(semanticParts(prefix, fmt.Sprintf("hop-%d", hop+1)), packet, []string{review.CallID}, append([]string{}, previous...), map[string]any{"design": "topology", "topology": topology, "order": orderName, "repeat": repeat, "kind": fmt.Sprintf("hop-%d", hop+1)})
						}
					}
					ids := callIDs(reviews)
					fused := b.call(semanticParts(prefix, "fuse"), packet, c.Provider, map[string]any{"kind": "fuser"}, ids, "fuse")
					if topology == "fanout" {
						b.outputSet(semanticParts(prefix, "raw_union"), packet, ids, ids, map[string]any{"design": "topology", "topology": topology, "order": orderName, "repeat": repeat, "kind": "raw_union"})
					}
					b.outputSet(semanticParts(prefix, "fused"), packet, []string{fused.CallID}, append(append([]string{}, ids...), fused.CallID), map[string]any{"design": "topology", "topology": topology, "order": orderName, "repeat": repeat, "kind": "fused"})
				}
			}
		}
	}
	return nil
}
func buildProviderPair(b *builder) error {
	c := b.config
	for _, packet := range c.Packets {
		for repeat := 0; repeat < c.Repetitions; repeat++ {
			for providerIndex, provider := range c.Providers {
				for _, wrapper := range c.Wrappers {
					prefix := []any{"provider_pair", packet, repeat, providerIndex, wrapper}
					review := b.call(semanticParts(prefix, "review"), packet, provider, map[string]any{"kind": "specialist", "wrapper": wrapper, "role": c.Role}, nil, "")
					b.outputSet(semanticParts(prefix, "final"), packet, []string{review.CallID}, []string{review.CallID}, map[string]any{"design": "provider_pair", "provider_index": providerIndex, "adapter": provider.Adapter, "model": provider.Model, "wrapper": wrapper, "role": c.Role, "repeat": repeat, "kind": "final"})
				}
			}
		}
	}
	return nil
}

func buildPersonaFactorial(b *builder) error {
	c := b.config
	for _, packet := range c.Packets {
		for repeat := 0; repeat < c.Repetitions; repeat++ {
			for _, role := range c.Roles {
				for _, wrapper := range c.Wrappers {
					prefix := []any{"persona_factorial", packet, repeat, role, wrapper}
					spec := map[string]any{"kind": "specialist", "wrapper": wrapper, "role": role}
					if role == "omnibus" {
						spec = map[string]any{"kind": "omnibus", "wrapper": wrapper}
					}
					review := b.call(semanticParts(prefix, "review"), packet, c.Provider, spec, nil, "")
					b.outputSet(semanticParts(prefix, "final"), packet, []string{review.CallID}, []string{review.CallID}, map[string]any{
						"design": "persona_factorial", "provider_index": 0, "adapter": c.Provider.Adapter,
						"model": c.Provider.Model, "wrapper": wrapper, "role": role, "repeat": repeat, "kind": "final",
					})
				}
			}
		}
	}
	return nil
}

func BuildPlan(config Config) (Plan, error) {
	b := &builder{config: config}
	var err error
	switch config.Design {
	case "stage_a":
		err = buildStageA(b)
	case "topology":
		err = buildTopology(b)
	case "provider_pair":
		err = buildProviderPair(b)
	case "persona_factorial":
		err = buildPersonaFactorial(b)
	default:
		err = fmt.Errorf("unknown design %s", config.Design)
	}
	if err != nil {
		return Plan{}, err
	}
	sort.SliceStable(b.calls, func(i, j int) bool {
		return shaText(config.Seed+"\x1forder\x1f"+b.calls[i].CallID) < shaText(config.Seed+"\x1forder\x1f"+b.calls[j].CallID)
	})
	return Plan{SchemaVersion: 2, CreatedAt: utcNow(), Design: config.Design, Seed: config.Seed, Calls: b.calls, OutputSets: b.sets, Counts: map[string]int{"calls": len(b.calls), "output_sets": len(b.sets)}}, nil
}
func (h *Harness) Plan(run string) error {
	resolved, err := h.resolveRun(run)
	if err != nil {
		return err
	}
	path := filepath.Join(resolved, "plan.json")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("plan.json already exists; do not replace a frozen run plan")
	}
	var config Config
	if err := readJSON(filepath.Join(resolved, "config.json"), &config); err != nil {
		return err
	}
	plan, err := BuildPlan(config)
	if err != nil {
		return err
	}
	if err := writeJSON(path, plan); err != nil {
		return err
	}
	fmt.Printf("Planned %d calls and %d output sets.\n", plan.Counts["calls"], plan.Counts["output_sets"])
	return nil
}
