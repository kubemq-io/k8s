package main

import (
	"go/parser"
	"go/token"
	"os"
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
// defaultClusterReplicationConfig — the exact 13 variadic string-literal args passed to
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
	extractKeys(f, keys, naturalKeys)

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

// TestServerEnvKeysGolden_ContainsClusterReplication13 pins the committed golden
// (parity/server_env_keys.txt) itself: it must carry exactly the 13 CLUSTER_REPLICATION_* names
// this F9 change adds, alongside the CLUSTER_* prefix's pre-existing keys — regold drift (missing
// or extra CLUSTER_REPLICATION_* entries) fails loudly here instead of silently in the downstream
// parity test.
func TestServerEnvKeysGolden_ContainsClusterReplication13(t *testing.T) {
	golden := readGoldenKeys(t, "../server_env_keys.txt")

	want := []string{
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
		"CLUSTER_REPLICATION_REPLICA_ID",
		"CLUSTER_REPLICATION_RTT_MILLISECOND",
		"CLUSTER_REPLICATION_SNAPSHOT_ENTRIES",
	}
	for _, w := range want {
		require.True(t, golden[w], "parity/server_env_keys.txt missing %s — regen via `task parity:regen`", w)
	}

	var gotReplication []string
	for k := range golden {
		if strings.HasPrefix(k, "CLUSTER_REPLICATION_") {
			gotReplication = append(gotReplication, k)
		}
	}
	sort.Strings(gotReplication)
	sort.Strings(want)
	require.Equal(t, want, gotReplication, "parity/server_env_keys.txt CLUSTER_REPLICATION_* keys "+
		"must be exactly the 13 F9 names — no stragglers, none missing")
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
