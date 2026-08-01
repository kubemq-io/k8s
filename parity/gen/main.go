// Command gen regenerates parity/server_env_keys.txt — the authoritative list of
// environment-variable names the kubemq-server honors. It is the server side of the
// CRD↔server config-parity contract (ADD-014).
//
// It is deliberately SELF-CONTAINED: it parses the kubemq-server config source with
// go/parser (stdlib only) and never imports the server module, so the k8s build has no
// cross-module dependency on the server. It extracts every dotted config key passed to
// bindViperEnv / bindCeEnvAliases / viper.BindEnv, applies the server's exact
// convertEnvFormat (copied verbatim from kubemq-server/config/env.go), adds the os.Getenv
// specials, and writes the sorted unique env-var names.
//
// Cluster.Replication.* keys are bound via the server's bindClusterReplicationEnv (not
// bindViperEnv/viper.BindEnv directly), and — unlike the CE aliases — the operator-expected
// CANONICAL name for these is the natural CLUSTER_REPLICATION_* form, not the acronym-collapsed
// convertEnvFormat form (which would drop the Cluster/Replication dot boundary into
// CLUSTER_REPLICATIONREPLICA_ID). Those literals are folded via naturalEnvName instead (also
// copied verbatim from kubemq-server/config/env.go), which snake-cases each dot segment
// independently — including the same acronym-splitting regexes — so ReplicaID -> REPLICA_ID,
// MutualTLS -> MUTUAL_TLS, RTTMillisecond -> RTT_MILLISECOND, CAFile -> CA_FILE.
//
// KEY ARGUMENTS ARE RESOLVED THROUGH PACKAGE-LEVEL SLICE VARS, and an unresolvable binder
// call is FATAL. Both matter: the server moved its replication keys out of the call site
// into `var clusterReplicationEnvKeys = []string{...}` and called
// `bindClusterReplicationEnv(clusterReplicationEnvKeys...)`. A literal-only extractor
// silently produced zero keys, so a regen quietly dropped all fourteen — the golden then
// under-reported what the server honors, and the parity gate it backs went partly blind
// on exactly the family the operator injects. Failing closed is the point: an extractor
// that finds nothing must stop, not shrink the golden.
//
// Usage (see `task parity:regen`):
//
//	go run ./parity/gen -server ../kubemq-server/config -out parity/server_env_keys.txt
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// --- verbatim from kubemq-server/config/env.go (keep in sync via `task parity:regen`) ---
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func convertEnvFormat(str string) string {
	return strings.ToUpper(strings.ReplaceAll(toSnakeCase(str), ".", ""))
}

// naturalEnvName converts a dotted config key into the operator-expected SCREAMING_SNAKE_CASE
// env name, preserving a separator at EVERY dot boundary — including boundaries around all-caps
// acronym segments (e.g. "Replication"/"ReplicaID"). convertEnvFormat collapses those acronym
// boundaries (Cluster.Replication.ReplicaID -> CLUSTER_REPLICATIONREPLICA_ID); this produces the
// natural CLUSTER_REPLICATION_REPLICA_ID instead. Each segment is independently snake-cased (via
// the same toSnakeCase acronym-splitting regexes above) so intra-segment camelCase/acronyms are
// still expanded (RTTMillisecond -> RTT_MILLISECOND, CAFile -> CA_FILE).
func naturalEnvName(key string) string {
	segments := strings.Split(key, ".")
	for i, seg := range segments {
		segments[i] = strings.ToUpper(toSnakeCase(seg))
	}
	return strings.Join(segments, "_")
}

// --- end verbatim ---

// os.Getenv specials the server reads directly (not through bindViperEnv), plus keys bound
// via viper.BindEnv with an explicit env name. HOST is bound via viper.BindEnv("Host","HOST")
// and is also picked up by the AST walk; the two below are pure os.Getenv reads.
var specials = []string{"KUBEMQ_TOKEN", "KUBEMQ_API_ADMIN_USERNAME"}

func main() {
	serverDir := flag.String("server", "../kubemq-server/config", "path to the kubemq-server config package")
	out := flag.String("out", "parity/server_env_keys.txt", "output golden file")
	flag.Parse()

	keys := map[string]bool{}        // dotted config keys -> convertEnvFormat (collapsed) name
	naturalKeys := map[string]bool{} // dotted config keys -> naturalEnvName (dot-preserving) name

	entries, err := os.ReadDir(*serverDir)
	if err != nil {
		fatal("read server dir %s: %v", *serverDir, err)
	}
	fset := token.NewFileSet()
	var files []*ast.File
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(*serverDir, e.Name()), nil, 0)
		if err != nil {
			fatal("parse %s: %v", e.Name(), err)
		}
		files = append(files, f)
	}

	// Package-level `var x = []string{...}` first: a binder call may spread one, and the
	// declaration can live in a different file from the call.
	sliceVars := map[string][]string{}
	for _, f := range files {
		collectStringSliceVars(f, sliceVars)
	}

	var unresolved []string
	for _, f := range files {
		unresolved = append(unresolved, extractKeys(f, keys, naturalKeys, sliceVars, fset)...)
	}
	if len(unresolved) > 0 {
		fatal("these binder calls resolved to no config keys — the extractor cannot see their "+
			"arguments, and generating anyway would silently shrink the golden:\n  %s",
			strings.Join(unresolved, "\n  "))
	}
	if len(keys) == 0 && len(naturalKeys) == 0 {
		fatal("no config keys found under %s — wrong path?", *serverDir)
	}

	envNames := buildEnvNames(keys, naturalKeys)
	for _, s := range specials {
		envNames[s] = true
	}

	list := make([]string, 0, len(envNames))
	for e := range envNames {
		list = append(list, e)
	}
	sort.Strings(list)

	var b strings.Builder
	b.WriteString("# Authoritative kubemq-server env-var names (ADD-014 config-parity golden).\n")
	b.WriteString("# GENERATED by `go run ./parity/gen` (task parity:regen). Do not hand-edit.\n")
	fmt.Fprintf(&b, "# %d keys.\n", len(list))
	for _, e := range list {
		b.WriteString(e)
		b.WriteString("\n")
	}
	if err := os.WriteFile(*out, []byte(b.String()), 0o644); err != nil {
		fatal("write %s: %v", *out, err)
	}
	fmt.Printf("wrote %d env keys to %s\n", len(list), *out)
}

// extractKeys walks a parsed server config .go file's AST and folds every dotted config key
// passed to bindViperEnv / bindCeEnvAliases / viper.BindEnv into keys, and every dotted config
// key passed to bindClusterReplicationEnv into naturalKeys (mutated in place; extracted from
// main's per-file loop so it is independently unit-testable against an in-memory AST).
// It returns a description of every binder call whose arguments it could not resolve to a
// single config key, so the caller can fail instead of emitting a short golden.
func extractKeys(f *ast.File, keys, naturalKeys map[string]bool, sliceVars map[string][]string, fset *token.FileSet) []string {
	var unresolved []string
	where := func(call *ast.CallExpr, name string) string {
		if fset == nil {
			return name
		}
		return fmt.Sprintf("%s at %s", name, fset.Position(call.Pos()))
	}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		name := callName(call.Fun)
		switch name {
		// bindKafkaOAuthBearerEnvAliases binds each key to BOTH convertEnvFormat and a
		// friendlier CONNECTORS_KAFKA_OAUTH_BEARER_* alias. Like the CE aliases, only the
		// collapsed form is listed: one canonical name per config key keeps the allowlist
		// from needing a second entry for the same underlying setting.
		case "bindViperEnv", "bindCeEnvAliases", "bindKafkaOAuthBearerEnvAliases":
			args := configKeyArgs(call, sliceVars)
			if len(args) == 0 {
				unresolved = append(unresolved, where(call, name))
			}
			for _, k := range args {
				keys[k] = true
			}
		case "viper.BindEnv":
			// arg0 is the dotted config key (the env-name args are derived, not keys).
			// Non-literal arg0 is normal here — the binder helpers call it in a loop over
			// their own parameter — so this case is deliberately not fail-closed.
			if lits := stringLits(call.Args); len(lits) > 0 {
				keys[lits[0]] = true
			}
		case "bindClusterReplicationEnv":
			// Variadic dotted config keys (e.g. "Cluster.Replication.ReplicaID"). These bind
			// to the natural CLUSTER_REPLICATION_* env name (see naturalEnvName), not the
			// acronym-collapsed convertEnvFormat name — kept in a separate set so the two
			// naming schemes don't cross-pollinate.
			args := configKeyArgs(call, sliceVars)
			if len(args) == 0 {
				unresolved = append(unresolved, where(call, name))
			}
			for _, k := range args {
				naturalKeys[k] = true
			}
		}
		return true
	})
	return unresolved
}

// collectStringSliceVars records every package-level `var name = []string{"a", "b"}` so a
// binder call that spreads one can be resolved back to its literals.
func collectStringSliceVars(f *ast.File, into map[string][]string) {
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, ident := range vs.Names {
				if i >= len(vs.Values) {
					continue
				}
				lit, ok := vs.Values[i].(*ast.CompositeLit)
				if !ok || !isStringSliceType(lit.Type) {
					continue
				}
				if vals := stringLits(lit.Elts); len(vals) > 0 {
					into[ident.Name] = vals
				}
			}
		}
	}
}

func isStringSliceType(expr ast.Expr) bool {
	arr, ok := expr.(*ast.ArrayType)
	if !ok || arr.Len != nil {
		return false
	}
	elt, ok := arr.Elt.(*ast.Ident)
	return ok && elt.Name == "string"
}

// configKeyArgs resolves a binder call's arguments to dotted config keys: inline string
// literals, plus any identifier (spread or not) naming a package-level string slice.
func configKeyArgs(call *ast.CallExpr, sliceVars map[string][]string) []string {
	out := stringLits(call.Args)
	for _, a := range call.Args {
		ident, ok := a.(*ast.Ident)
		if !ok {
			continue
		}
		out = append(out, sliceVars[ident.Name]...)
	}
	return out
}

// buildEnvNames folds the two dotted-key sets extractKeys collects into the final env-name set:
// keys via convertEnvFormat (the acronym-collapsed form — canonical for everything except
// Cluster.Replication.*), naturalKeys via naturalEnvName (the dot-preserving form — canonical
// for Cluster.Replication.*, F9). Does NOT add the os.Getenv specials (main adds those on top).
func buildEnvNames(keys, naturalKeys map[string]bool) map[string]bool {
	envNames := map[string]bool{}
	for k := range keys {
		// Canonical env name = convertEnvFormat (what the operator's SetConfig emits, incl. the
		// collapsed CONNECTORSCE_* form for CE). The server's natural-name CE aliases
		// (bindCeEnvAliases) are redundant compat aliases and are intentionally NOT listed —
		// the collapsed form is the single source of truth for parity.
		envNames[convertEnvFormat(k)] = true
	}
	for k := range naturalKeys {
		// Cluster.Replication.* is the inverse of the CE case above: the natural dot-preserving
		// name IS the canonical/operator-expected form (F9), so fold naturalEnvName, not
		// convertEnvFormat.
		envNames[naturalEnvName(k)] = true
	}
	return envNames
}

func callName(fun ast.Expr) string {
	switch f := fun.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		if x, ok := f.X.(*ast.Ident); ok {
			return x.Name + "." + f.Sel.Name
		}
	}
	return ""
}

func stringLits(args []ast.Expr) []string {
	var out []string
	for _, a := range args {
		if bl, ok := a.(*ast.BasicLit); ok && bl.Kind == token.STRING {
			out = append(out, strings.Trim(bl.Value, "`\""))
		}
	}
	return out
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "parity/gen: "+format+"\n", a...)
	os.Exit(1)
}
