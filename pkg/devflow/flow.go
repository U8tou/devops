package devflow

import (
	"encoding/json"
	"sort"
	"strings"
)

// FlowEdge Vue Flow 边
type FlowEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// FlowNode 画布节点（仅解析执行所需字段）
type FlowNode struct {
	ID       string `json:"id"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"`
	Data json.RawMessage `json:"data"`
}

type flowDoc struct {
	Nodes []FlowNode `json:"nodes"`
	Edges []FlowEdge `json:"edges"`
}

// ParseFlow 解析流程 JSON（含 nodes/edges）
func ParseFlow(flowJSON string) ([]FlowNode, []FlowEdge, error) {
	var doc flowDoc
	if err := json.Unmarshal([]byte(flowJSON), &doc); err != nil {
		return nil, nil, err
	}
	if doc.Nodes == nil {
		doc.Nodes = []FlowNode{}
	}
	if doc.Edges == nil {
		doc.Edges = []FlowEdge{}
	}
	return doc.Nodes, doc.Edges, nil
}

func cmpNodeVisual(a, b *FlowNode) int {
	ax, ay := a.Position.X, a.Position.Y
	bx, by := b.Position.X, b.Position.Y
	if ax != bx {
		if ax < bx {
			return -1
		}
		return 1
	}
	if ay != by {
		if ay < by {
			return -1
		}
		return 1
	}
	return strings.Compare(a.ID, b.ID)
}

// TopoOrderedNodes 拓扑序 + 同层按画布坐标（与 process-run-dialog.vue 一致）
func TopoOrderedNodes(nodes []FlowNode, edges []FlowEdge) []FlowNode {
	byID := make(map[string]*FlowNode, len(nodes))
	for i := range nodes {
		byID[nodes[i].ID] = &nodes[i]
	}
	inc := make(map[string]int, len(nodes))
	for _, n := range nodes {
		inc[n.ID] = 0
	}
	for _, e := range edges {
		if _, ok := byID[e.Target]; ok {
			inc[e.Target]++
		}
	}
	next := make(map[string][]string)
	for _, e := range edges {
		next[e.Source] = append(next[e.Source], e.Target)
	}
	for _, targets := range next {
		sort.Slice(targets, func(i, j int) bool {
			return cmpNodeVisual(byID[targets[i]], byID[targets[j]]) < 0
		})
	}

	var ready []*FlowNode
	for i := range nodes {
		if inc[nodes[i].ID] == 0 {
			ready = append(ready, &nodes[i])
		}
	}
	sort.Slice(ready, func(i, j int) bool { return cmpNodeVisual(ready[i], ready[j]) < 0 })

	out := make([]FlowNode, 0, len(nodes))
	processed := make(map[string]bool)

	for len(ready) > 0 {
		n := ready[0]
		ready = ready[1:]
		if processed[n.ID] {
			continue
		}
		processed[n.ID] = true
		out = append(out, *n)

		for _, t := range next[n.ID] {
			inc[t]--
			if inc[t] == 0 {
				if node := byID[t]; node != nil {
					ready = append(ready, node)
				}
			}
		}
		sort.Slice(ready, func(i, j int) bool { return cmpNodeVisual(ready[i], ready[j]) < 0 })
	}

	var rest []FlowNode
	for i := range nodes {
		if !processed[nodes[i].ID] {
			rest = append(rest, nodes[i])
		}
	}
	sort.Slice(rest, func(i, j int) bool { return cmpNodeVisual(&rest[i], &rest[j]) < 0 })
	return append(out, rest...)
}

// nodeFlowData 仅解 flowEnabled
type nodeFlowData struct {
	FlowEnabled *bool `json:"flowEnabled"`
}

// NodeFlowEnabled 与前端 AutomationNodePersistData.flowEnabled 一致，缺省为 true
func NodeFlowEnabled(data json.RawMessage) bool {
	if len(data) == 0 {
		return true
	}
	var t nodeFlowData
	if err := json.Unmarshal(data, &t); err != nil {
		return true
	}
	if t.FlowEnabled == nil {
		return true
	}
	return *t.FlowEnabled
}

// BuildDownstreamFlowSkip 返回应跳过的节点 id（本节点被关闭的节点 + 其沿出边可达的下游闭包）
func BuildDownstreamFlowSkip(nodes []FlowNode, edges []FlowEdge) map[string]struct{} {
	var roots []string
	for i := range nodes {
		if !NodeFlowEnabled(nodes[i].Data) {
			roots = append(roots, nodes[i].ID)
		}
	}
	next := make(map[string][]string)
	for _, e := range edges {
		next[e.Source] = append(next[e.Source], e.Target)
	}
	out := make(map[string]struct{})
	q := append([]string(nil), roots...)
	for len(q) > 0 {
		id := q[0]
		q = q[1:]
		if _, ok := out[id]; ok {
			continue
		}
		out[id] = struct{}{}
		for _, t := range next[id] {
			if _, ok := out[t]; !ok {
				q = append(q, t)
			}
		}
	}
	return out
}
