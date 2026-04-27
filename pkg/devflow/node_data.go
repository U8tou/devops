package devflow

import (
	"encoding/json"
	"sort"
	"strings"
)

// ParseNodeData 解析节点 data（kind + params）
func ParseNodeData(data json.RawMessage) (kind string, params map[string]interface{}, err error) {
	var raw struct {
		Kind   string          `json:"kind"`
		Params json.RawMessage `json:"params"`
	}
	if len(data) == 0 {
		return "git_repo", map[string]interface{}{}, nil
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return "", nil, err
	}
	kind = NormalizeNodeKind(raw.Kind)
	params = map[string]interface{}{}
	if len(raw.Params) > 0 {
		_ = json.Unmarshal(raw.Params, &params)
	}
	return kind, params, nil
}

// NodeDisplayNameForLog 节点在日志中的展示名：优先 data.label，空则回退为 kind（内部类型名）
func NodeDisplayNameForLog(data json.RawMessage, kind string) string {
	var raw struct {
		Label string `json:"label"`
	}
	if len(data) > 0 {
		_ = json.Unmarshal(data, &raw)
	}
	s := strings.TrimSpace(raw.Label)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	if s != "" {
		return s
	}
	return kind
}

// SshParams 执行用 SSH 凭据（脱敏日志勿打印密码/私钥）
type SshParams struct {
	Host, Port, Username, Password, PrivateKey, AuthType string
}

func sshParamsFromNode(kind string, params map[string]interface{}) *SshParams {
	switch kind {
	case "ssh_connection":
		port := strParam(params, "port")
		if strings.TrimSpace(port) == "" {
			port = "22"
		}
		return &SshParams{
			Host:       strings.TrimSpace(strParam(params, "host")),
			Port:       port,
			Username:   strParam(params, "username"),
			Password:   strParam(params, "password"),
			PrivateKey: strParam(params, "privateKey"),
			AuthType:   strParam(params, "authType"),
		}
	case "upload_servers":
		host := strings.TrimSpace(strParam(params, "host"))
		port := strParam(params, "port")
		if strings.TrimSpace(port) == "" {
			port = "22"
		}
		user := strParam(params, "username")
		// 与前端一致：结构化字段 host 不含 @
		if host != "" && strings.TrimSpace(user) != "" && !strings.Contains(host, "@") {
			return &SshParams{
				Host:       host,
				Port:       port,
				Username:   user,
				Password:   strParam(params, "password"),
				PrivateKey: strParam(params, "privateKey"),
				AuthType:   strParam(params, "authType"),
			}
		}
		// 兼容旧版单字段 host：user@host:port
		t, ok := ParseUploadHost(strParam(params, "host"))
		if !ok || strings.TrimSpace(t.Host) == "" {
			return nil
		}
		if t.Port != "" {
			port = t.Port
		}
		if port == "" {
			port = "22"
		}
		user = t.Username
		if user == "" {
			user = "root"
		}
		return &SshParams{
			Host:     t.Host,
			Port:     port,
			Username: user,
			AuthType: "key",
		}
	default:
		return nil
	}
}

func findNodeByID(nodes []FlowNode, id string) *FlowNode {
	for i := range nodes {
		if nodes[i].ID == id {
			return &nodes[i]
		}
	}
	return nil
}

func predsOf(target string, edges []FlowEdge) []string {
	var out []string
	for _, e := range edges {
		if e.Target == target {
			out = append(out, e.Source)
		}
	}
	return out
}

// SshFromPredecessorBranch 仅沿「一条入边」解析 SSH：先看该上游节点自身（上传/SSH），否则再沿其上游回溯。
func SshFromPredecessorBranch(nodes []FlowNode, edges []FlowEdge, predecessorID string) *SshParams {
	n := findNodeByID(nodes, predecessorID)
	if n == nil {
		return nil
	}
	kind, params, err := ParseNodeData(n.Data)
	if err != nil {
		return nil
	}
	if direct := sshParamsFromNode(kind, params); direct != nil && strings.TrimSpace(direct.Host) != "" {
		return direct
	}
	return ResolveUpstreamSsh(predecessorID, nodes, edges, nil)
}

func sshTargetDedupKey(sp *SshParams) string {
	u := strings.TrimSpace(sp.Username)
	h := strings.TrimSpace(sp.Host)
	p := strings.TrimSpace(sp.Port)
	if p == "" {
		p = "22"
	}
	return strings.ToLower(u + "@" + h + ":" + p)
}

// CollectSshTargetsForRemoteScript 收集远程脚本节点每条入边对应的 SSH；同一 user@host:port 去重，依次在多机执行时各跑一遍。
func CollectSshTargetsForRemoteScript(nodeID string, nodes []FlowNode, edges []FlowEdge) []*SshParams {
	preds := predsOf(nodeID, edges)
	sort.Strings(preds)
	seen := make(map[string]bool)
	var out []*SshParams
	for _, sid := range preds {
		sp := SshFromPredecessorBranch(nodes, edges, sid)
		if sp == nil || strings.TrimSpace(sp.Host) == "" {
			continue
		}
		k := sshTargetDedupKey(sp)
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, sp)
	}
	return out
}

// ResolveUpstreamSsh 与 flow-graph-ssh.ts resolveUpstreamSsh 一致
func ResolveUpstreamSsh(nodeID string, nodes []FlowNode, edges []FlowEdge, visited map[string]bool) *SshParams {
	if visited == nil {
		visited = make(map[string]bool)
	}
	if visited[nodeID] {
		return nil
	}
	visited[nodeID] = true
	preds := predsOf(nodeID, edges)
	sort.Strings(preds)
	for _, sid := range preds {
		n := findNodeByID(nodes, sid)
		if n == nil {
			continue
		}
		kind, params, err := ParseNodeData(n.Data)
		if err != nil {
			continue
		}
		if direct := sshParamsFromNode(kind, params); direct != nil && strings.TrimSpace(direct.Host) != "" {
			return direct
		}
		if up := ResolveUpstreamSsh(sid, nodes, edges, visited); up != nil && strings.TrimSpace(up.Host) != "" {
			return up
		}
	}
	return nil
}
