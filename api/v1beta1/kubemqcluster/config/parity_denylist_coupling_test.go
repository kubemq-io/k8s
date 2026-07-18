// Package config_test is deliberately an EXTERNAL test package (not `package config`):
// the coupling assertion needs both github.com/kubemq-io/k8s/api/v1beta1 (for
// OperatorComputedEnvKeys / OperatorComputedEnvKeyPrefix) and this config package, and
// v1beta1 already imports config — an internal `package config` test importing v1beta1
// would create an import cycle. External test packages compile as a separate unit during
// `go test`, so this is cycle-free.
package config_test

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kubemq-io/k8s/api/v1beta1"
)

// parityDenylistCouplingAllowlistPath is relative to this package dir
// (api/v1beta1/kubemqcluster/config → four levels up to the k8s module root), matching
// parityAllowlistPath in parity_test.go.
const parityDenylistCouplingAllowlistPath = "../../../../parity/configdata_allowlist.txt"

// operatorComputedTag is the disposition-comment marker (F9) identifying allowlist entries
// the operator itself derives/injects — the same set that must be denylisted from spec.env
// (F6.3). See parity/configdata_allowlist.txt's header comment.
const operatorComputedTag = "operator-computed"

// TestParity_OperatorComputedDispositionsAreDenylisted is the F6.3 drift guard: every
// configdata_allowlist.txt entry tagged "operator-computed" must be covered by
// v1beta1.OperatorComputedEnvKeys (exact-match entries) or v1beta1.OperatorComputedEnvKeyPrefix
// (prefix-match entries). Without this, a future F9 allowlist split could tag a new
// operator-computed key/prefix without also denylisting it from spec.env, silently reopening
// the collision the F6 overlay guard exists to close.
func TestParity_OperatorComputedDispositionsAreDenylisted(t *testing.T) {
	exact, prefixes := readOperatorComputedDispositions(t, parityDenylistCouplingAllowlistPath)

	tagged := append(append([]string{}, exact...), prefixes...)
	require.NotEmpty(t, tagged, "no operator-computed-tagged disposition found in %s — "+
		"has the F9 tag text changed, or did the allowlist lose its operator-computed entries?",
		parityDenylistCouplingAllowlistPath)

	allowedExact := map[string]bool{}
	for _, k := range v1beta1.OperatorComputedEnvKeys {
		allowedExact[k] = true
	}

	var bad []string
	for _, k := range exact {
		if allowedExact[k] || strings.HasPrefix(k, v1beta1.OperatorComputedEnvKeyPrefix) {
			continue
		}
		bad = append(bad, k+" (exact)")
	}
	for _, p := range prefixes {
		if p == v1beta1.OperatorComputedEnvKeyPrefix || strings.HasPrefix(p, v1beta1.OperatorComputedEnvKeyPrefix) {
			continue
		}
		bad = append(bad, p+"* (prefix)")
	}
	sort.Strings(bad)
	require.Empty(t, bad, "configdata_allowlist.txt entries tagged %q are not covered by "+
		"v1beta1.OperatorComputedEnvKeys / v1beta1.OperatorComputedEnvKeyPrefix — an "+
		"operator-computed key must also be denylisted from spec.env (F6.3):\n%s",
		operatorComputedTag, strings.Join(bad, "\n"))
}

// readOperatorComputedDispositions parses configdata_allowlist.txt and returns the key
// portion of every entry whose trailing comment contains operatorComputedTag, split into
// exact entries (no trailing `*`) and prefix entries (trailing `*`, stripped).
func readOperatorComputedDispositions(t *testing.T, path string) (exact []string, prefix []string) {
	t.Helper()
	f, err := os.Open(path)
	require.NoError(t, err, "open %s", path)
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		hashIdx := strings.Index(line, "#")
		if hashIdx < 0 {
			continue // no comment => no disposition to tag
		}
		key := strings.TrimSpace(line[:hashIdx])
		comment := line[hashIdx+1:]
		if key == "" || !strings.Contains(comment, operatorComputedTag) {
			continue
		}
		if strings.HasSuffix(key, "*") {
			prefix = append(prefix, strings.TrimSuffix(key, "*"))
		} else {
			exact = append(exact, key)
		}
	}
	require.NoError(t, sc.Err())
	return exact, prefix
}
