package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNaturalEnvName_ClusterReplicationAcronyms is a direct table-driven check that
// naturalEnvName reproduces the server's exact acronym-splitting semantics (copied verbatim from
// kubemq-server/config/env.go's own naturalEnvName): each dot segment is independently
// snake-cased via the matchFirstCap/matchAllCap regexes, so an all-caps acronym run is split from
// the segment around it (ReplicaID -> REPLICA_ID) while an all-caps *tail* run stays fused to its
// own boundary correctly (RTTMillisecond -> RTT_MILLISECOND, MutualTLS -> MUTUAL_TLS, CAFile ->
// CA_FILE) — the exact trap the plan (§3.2.7 F9) calls out: a naive per-segment snake-case would
// wrongly yield REPLICAID / MUTUALTLS.
func TestNaturalEnvName_ClusterReplicationAcronyms(t *testing.T) {
	cases := []struct {
		key  string
		want string
	}{
		{"Cluster.Replication.ReplicaID", "CLUSTER_REPLICATION_REPLICA_ID"},
		{"Cluster.Replication.Peers", "CLUSTER_REPLICATION_PEERS"},
		{"Cluster.Replication.Join", "CLUSTER_REPLICATION_JOIN"},
		{"Cluster.Replication.RTTMillisecond", "CLUSTER_REPLICATION_RTT_MILLISECOND"},
		{"Cluster.Replication.ElectionRTT", "CLUSTER_REPLICATION_ELECTION_RTT"},
		{"Cluster.Replication.HeartbeatRTT", "CLUSTER_REPLICATION_HEARTBEAT_RTT"},
		{"Cluster.Replication.SnapshotEntries", "CLUSTER_REPLICATION_SNAPSHOT_ENTRIES"},
		{"Cluster.Replication.CompactionOverhead", "CLUSTER_REPLICATION_COMPACTION_OVERHEAD"},
		{"Cluster.Replication.BootTimeoutSeconds", "CLUSTER_REPLICATION_BOOT_TIMEOUT_SECONDS"},
		{"Cluster.Replication.MutualTLS", "CLUSTER_REPLICATION_MUTUAL_TLS"},
		{"Cluster.Replication.CAFile", "CLUSTER_REPLICATION_CA_FILE"},
		{"Cluster.Replication.CertFile", "CLUSTER_REPLICATION_CERT_FILE"},
		{"Cluster.Replication.KeyFile", "CLUSTER_REPLICATION_KEY_FILE"},
	}
	for _, c := range cases {
		t.Run(c.key, func(t *testing.T) {
			require.Equal(t, c.want, naturalEnvName(c.key))
		})
	}
}

// fakeClusterReplicationSource mirrors the SHAPE of kubemq-server/config/cluster.go's
// defaultClusterReplicationConfig — variadic string-literal args passed to
// bindClusterReplicationEnv — without importing or depending on the real sibling repo, so this
// test is self-contained and stays deterministic even if the server file moves/reformats.
const fakeClusterReplicationSource = `package fakeconfig

func defaultClusterReplicationConfig() {
	bindClusterReplicationEnv(
		"Cluster.Replication.ReplicaID",
		"Cluster.Replication.Peers",
		"Cluster.Replication.Join",
		"Cluster.Replication.RTTMillisecond",
		"Cluster.Replication.ElectionRTT",
		"Cluster.Replication.HeartbeatRTT",
		"Cluster.Replication.SnapshotEntries",
		"Cluster.Replication.CompactionOverhead",
		"Cluster.Replication.BootTimeoutSeconds",
		"Cluster.Replication.MutualTLS",
		"Cluster.Replication.CAFile",
		"Cluster.Replication.CertFile",
		"Cluster.Replication.KeyFile",
	)
}
`

// fakeClusterReplicationSliceVarSource is the shape the server ACTUALLY uses today: the keys
// live in a package-level slice and the call spreads it. A literal-only extractor produced zero
// keys here and silently dropped the whole family from the golden — this pins the resolution.
const fakeClusterReplicationSliceVarSource = `package fakeconfig

var clusterReplicationEnvKeys = []string{
	"Cluster.Replication.ReplicaID",
	"Cluster.Replication.Peers",
	"Cluster.Replication.RaftAddress",
}

func defaultClusterReplicationConfig() {
	bindClusterReplicationEnv(clusterReplicationEnvKeys...)
}
`

// TestExtractKeys_BindClusterReplicationEnv is the F9 gen unit test: it feeds the
// "bindClusterReplicationEnv" case an in-memory AST shaped exactly like the server's
// defaultClusterReplicationConfig call and asserts the generated env-name set is EXACTLY the 13
// CLUSTER_REPLICATION_* names — no more, no less, and none of them leaking into the
// convertEnvFormat (collapsed) pathway.
func TestExtractKeys_BindClusterReplicationEnv(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fakeconfig.go", fakeClusterReplicationSource, 0)
	require.NoError(t, err)

	keys := map[string]bool{}
	naturalKeys := map[string]bool{}
	require.Empty(t, extractKeys(f, keys, naturalKeys, map[string][]string{}, fset))

	require.Empty(t, keys, "bindClusterReplicationEnv literals must not leak into the "+
		"convertEnvFormat (collapsed) key set")
	require.Len(t, naturalKeys, 13)

	envNames := buildEnvNames(keys, naturalKeys)

	want := []string{
		"CLUSTER_REPLICATION_REPLICA_ID",
		"CLUSTER_REPLICATION_PEERS",
		"CLUSTER_REPLICATION_JOIN",
		"CLUSTER_REPLICATION_RTT_MILLISECOND",
		"CLUSTER_REPLICATION_ELECTION_RTT",
		"CLUSTER_REPLICATION_HEARTBEAT_RTT",
		"CLUSTER_REPLICATION_SNAPSHOT_ENTRIES",
		"CLUSTER_REPLICATION_COMPACTION_OVERHEAD",
		"CLUSTER_REPLICATION_BOOT_TIMEOUT_SECONDS",
		"CLUSTER_REPLICATION_MUTUAL_TLS",
		"CLUSTER_REPLICATION_CA_FILE",
		"CLUSTER_REPLICATION_CERT_FILE",
		"CLUSTER_REPLICATION_KEY_FILE",
	}
	sort.Strings(want)

	got := make([]string, 0, len(envNames))
	for e := range envNames {
		got = append(got, e)
	}
	sort.Strings(got)

	require.Equal(t, want, got, "generated Cluster.Replication env names must be exactly the 13 "+
		"CLUSTER_REPLICATION_* names — no more, no less")
}

// TestExtractKeys_ResolvesPackageLevelSliceVar pins the shape the server uses today. The
// literal-only extractor silently returned zero keys here, so a regen dropped all fourteen
// CLUSTER_REPLICATION_* names and shrank the golden without a word.
func TestExtractKeys_ResolvesPackageLevelSliceVar(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fakeconfig.go", fakeClusterReplicationSliceVarSource, 0)
	require.NoError(t, err)

	sliceVars := map[string][]string{}
	collectStringSliceVars(f, sliceVars)
	require.Len(t, sliceVars["clusterReplicationEnvKeys"], 3)

	keys := map[string]bool{}
	naturalKeys := map[string]bool{}
	require.Empty(t, extractKeys(f, keys, naturalKeys, sliceVars, fset))

	require.Empty(t, keys)
	got := make([]string, 0, len(naturalKeys))
	for k := range buildEnvNames(keys, naturalKeys) {
		got = append(got, k)
	}
	sort.Strings(got)
	require.Equal(t, []string{
		"CLUSTER_REPLICATION_PEERS",
		"CLUSTER_REPLICATION_RAFT_ADDRESS",
		"CLUSTER_REPLICATION_REPLICA_ID",
	}, got)
}

// A binder call whose arguments the extractor cannot see must be REPORTED, not skipped.
// Silently emitting a shorter golden is how the replication family went missing.
func TestExtractKeys_UnresolvableBinderCallIsReported(t *testing.T) {
	const src = `package fakeconfig

func f(keys []string) {
	bindClusterReplicationEnv(keys...)
	bindViperEnv(someUnknownVar...)
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fakeconfig.go", src, 0)
	require.NoError(t, err)

	unresolved := extractKeys(f, map[string]bool{}, map[string]bool{}, map[string][]string{}, fset)
	require.Len(t, unresolved, 2, "both unresolvable binder calls must be reported")
}

// The OAUTHBEARER binder folds to the COLLAPSED name, like the CE aliases: one canonical env
// name per config key, so the allowlist needs a single entry per setting. The server also binds
// a friendlier CONNECTORS_KAFKA_OAUTH_BEARER_* alias, deliberately not listed.
func TestExtractKeys_KafkaOAuthBearerAliasesUseCollapsedName(t *testing.T) {
	const src = `package fakeconfig

func defaultKafkaConfig() {
	bindKafkaOAuthBearerEnvAliases(
		"Connectors.Kafka.OAuthBearer.Issuer",
		"Connectors.Kafka.OAuthBearer.ClientID",
	)
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fakeconfig.go", src, 0)
	require.NoError(t, err)

	keys := map[string]bool{}
	naturalKeys := map[string]bool{}
	require.Empty(t, extractKeys(f, keys, naturalKeys, map[string][]string{}, fset))
	require.Empty(t, naturalKeys)

	got := make([]string, 0, 2)
	for k := range buildEnvNames(keys, naturalKeys) {
		got = append(got, k)
	}
	sort.Strings(got)
	require.Equal(t, []string{
		"CONNECTORS_KAFKAO_AUTH_BEARER_CLIENT_ID",
		"CONNECTORS_KAFKAO_AUTH_BEARER_ISSUER",
	}, got)
}

// TestServerEnvKeysGolden_MatchesServer closes the staleness loop the previous pin could not:
// instead of hard-coding a key count, it re-runs the extractor against the real kubemq-server
// checkout and requires the committed golden to equal it exactly. A server that gains or renames
// an env key now fails HERE, naming the drift, rather than leaving the golden quietly behind and
// the parity gate it backs partly blind. Skips when the sibling is absent (standalone checkout),
// the same convention the CRD sibling tests use.
func TestServerEnvKeysGolden_MatchesServer(t *testing.T) {
	const serverDir = "../../../kubemq-server/config"
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		t.Skipf("SKIP sibling: %s not present (standalone checkout?)", serverDir)
	}

	entries, err := os.ReadDir(serverDir)
	require.NoError(t, err)

	fset := token.NewFileSet()
	var files []*ast.File
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(serverDir, e.Name()), nil, 0)
		require.NoError(t, err, "parse %s", e.Name())
		files = append(files, f)
	}

	sliceVars := map[string][]string{}
	for _, f := range files {
		collectStringSliceVars(f, sliceVars)
	}
	keys := map[string]bool{}
	naturalKeys := map[string]bool{}
	var unresolved []string
	for _, f := range files {
		unresolved = append(unresolved, extractKeys(f, keys, naturalKeys, sliceVars, fset)...)
	}
	require.Empty(t, unresolved, "the extractor cannot see these binder calls — it would emit a short golden")

	envNames := buildEnvNames(keys, naturalKeys)
	for _, s := range specials {
		envNames[s] = true
	}
	want := make([]string, 0, len(envNames))
	for e := range envNames {
		want = append(want, e)
	}
	sort.Strings(want)

	goldenSet := readGoldenKeys(t, "../server_env_keys.txt")
	got := make([]string, 0, len(goldenSet))
	for k := range goldenSet {
		got = append(got, k)
	}
	sort.Strings(got)

	require.Equal(t, want, got, "parity/server_env_keys.txt is out of date with the kubemq-server "+
		"checkout — regenerate it with `task parity:regen`")
}

// The replication family is the one the operator itself injects, and the one a regression
// already dropped once. Keep a direct floor on it so a partial regen cannot quietly halve it.
func TestServerEnvKeysGolden_CarriesClusterReplicationFamily(t *testing.T) {
	golden := readGoldenKeys(t, "../server_env_keys.txt")

	for _, w := range []string{
		"CLUSTER_REPLICATION_BOOT_TIMEOUT_SECONDS",
		"CLUSTER_REPLICATION_CA_FILE",
		"CLUSTER_REPLICATION_CERT_FILE",
		"CLUSTER_REPLICATION_COMPACTION_OVERHEAD",
		"CLUSTER_REPLICATION_ELECTION_RTT",
		"CLUSTER_REPLICATION_HEARTBEAT_RTT",
		"CLUSTER_REPLICATION_JOIN",
		"CLUSTER_REPLICATION_KEY_FILE",
		"CLUSTER_REPLICATION_MUTUAL_TLS",
		"CLUSTER_REPLICATION_PEERS",
		"CLUSTER_REPLICATION_RAFT_ADDRESS",
		"CLUSTER_REPLICATION_REPLICA_ID",
		"CLUSTER_REPLICATION_RTT_MILLISECOND",
		"CLUSTER_REPLICATION_SNAPSHOT_ENTRIES",
	} {
		require.True(t, golden[w], "parity/server_env_keys.txt missing %s — regen via `task parity:regen`", w)
	}
}

func readGoldenKeys(t *testing.T, path string) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err, "open %s", path)
	keys := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keys[line] = true
	}
	return keys
}
