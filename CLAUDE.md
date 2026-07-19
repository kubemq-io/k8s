# CLAUDE.md — kubemq-io/k8s (CRD lib) engineering rules

Governing guide for the `github.com/kubemq-io/k8s` library. It holds the KubeMQ CRD
Go types, the object builders / config emitters the operator vendors, and a
bootstrap deployer with the CRD YAML embedded as Go string literals. The
kubemq-operator consumes this repo through a local `replace` directive and
re-vendors it after every change here — so a change in this repo is a change to the
operator's build. Keep it lean.

The **LEAN-CRD doctrine** (§2) is the load-bearing rule of this repo; read it before
adding any field or config key.

---

## 1. Project & layout

`github.com/kubemq-io/k8s` (go 1.17 in `go.mod`; compiled inside the operator's
go 1.25 build) is three things in one module:

- **CRD API types** — `api/v1beta1/` (served/current) and `api/v1alpha1/` (kept as the
  storage version; see §2.4). `KubemqCluster`, `KubemqConnector`, `KubemqDashboard`.
- **Object builders + config emitters** — `api/v1beta1/kubemqcluster/{config,deployment}/`.
  Each connector/subsystem is a `*Config` struct with a `SetConfig(*deployment.Config)`
  method that writes the **exact server env-var name** into the operator-owned ConfigMap
  (`SetConfigMapStringValues`). `deployment/` holds the StatefulSet/Service/ConfigMap
  templates.
- **Bootstrap deployer** — `controller/` applies KubeMQ into a cluster directly (still
  used, Non-Goal 6). The CRD schemas ship here as backtick Go literals in
  `controller/objects/manifests/*.go`, and the deployer's own typed structs live in
  `controller/objects/v1beta1/` and `controller/objects/v1alpha1/`.

| Area | Path |
|---|---|
| CRD spec/status types + top-level deepcopy | `api/v1beta1/kubemqcluster_types.go`, `api/v1beta1/zz_generated.deepcopy.go` |
| Per-subsystem config structs + `SetConfig` emitters | `api/v1beta1/kubemqcluster/config/*.go` |
| StatefulSet / Service / ConfigMap templates & builders | `api/v1beta1/kubemqcluster/deployment/*.go` |
| v1alpha1 mirror (storage version) | `api/v1alpha1/**` |
| Parity generator + golden + LEAN-CRD ledger | `parity/gen/main.go`, `parity/server_env_keys.txt`, `parity/configdata_allowlist.txt` |
| Bootstrap deployer + embedded CRD YAML literals | `controller/`, `controller/objects/manifests/*.go` |

The operator's `../kubemq-operator/CLAUDE.md` holds the reconcile/billing rules and the
operator-side half of this doctrine (`spec.env` application, the establishment guard).

---

## 2. LEAN-CRD doctrine (the load-bearing rule)

KubeMQ's promise is "a license and that's it." The server has ~370 env-tunables. If
every one had to become a typed CRD field, adding a single server config key would mean
edits in **five hand-maintained places** — server config → this k8s lib (type + deepcopy
+ `SetConfig`) → operator env generation → three CRD YAML copies → the parity golden.
That "5-place drift tax" is exactly what this doctrine exists to stop.

### 2.1 The rule: type Day-1 fields only, everything else via the env overlay

**The CRD types ONLY Day-1 fields:** connector `enable`, `port`(s), `advertisedHost`/
`advertisedPort`, `expose` (+ `nodePort` on single-port connectors), and credential/secret
references. **Every other server tunable is reachable without a schema change**, through
two generic escape hatches on `KubemqClusterSpec`:

- `Env map[string]string` — an overlay of raw `SERVER_ENV_NAME: value` pairs applied over
  the operator's computed ConfigMap (last-wins; the operator uppercases keys, drops empty
  values, and rejects the operator-computed keys — see §2.2).
- `EnvFromSecrets []string` — a list of Secret names projected onto the pods via `envFrom`
  (wired through `deployment.StatefulSetConfig.ExtraSecretRefs` / `SetExtraSecretRefs`), so
  credentials never transit the operator.

**Add a new TYPED field only when the value is one of:**
1. **Day-1** — part of the minimal "make this connector work" surface a first-time user sets.
2. **Must be validated** — needs an enum/range/format the API server enforces at admission
   (e.g. `expose` enum, `nodePort` 30000-32767).
3. **Must be operator-computed** — the operator derives or injects it (engine, cluster
   identity, replication wiring); a user must NOT be able to set it (see the denylist, §2.2).

If none of those hold, it is a tunable: leave it untyped and let it flow through `spec.env`.
Adding a typed field for a plain tunable is the over-engineering this doctrine forbids —
it re-incurs the 5-place tax for zero benefit. Existing typed fields are frozen (never
removed — Non-Goal 4); this doctrine governs what NEW surface to add.

### 2.2 OperatorComputedEnvKeys — the denylist that protects operator identity

The env overlay's escape hatch cannot be allowed to override the keys the operator itself
owns (engine selection, cluster name/routes, checksum, pod identity, replication) — a user
setting those would corrupt the deployment. The single source of truth is the exported
denylist in `api/v1beta1/kubemqcluster_types.go`:

- `OperatorComputedEnvKeys` — exact names: `STORE_ENGINE`, `CLUSTER_ENABLE`, `CLUSTER_NAME`,
  `CLUSTER_ROUTES`, `API_BIND_ADDRESS`, `CHECKSUM`, `POD_NAME`.
- `OperatorComputedEnvKeyPrefix` — `CLUSTER_REPLICATION_` (all next-engine raft membership).

The operator vendors these constants and **rejects** (does not silently strip) any
`spec.env` key that matches — declared intent must never silently diverge. This list is
the ONLY reason to be careful about what `spec.env` can reach; everything else is fair game.

### 2.3 Enforcement loop: parity generator → golden → allowlist ledger

The doctrine is not honor-system; it is machine-enforced so a server key can never quietly
fall through the cracks. Three artifacts + the parity test close the loop:

- **`parity/gen/main.go`** — parses the kubemq-server `config/` package with `go/parser`
  (stdlib only, no cross-module import) and emits **`parity/server_env_keys.txt`**, the
  authoritative sorted list of every env-var name the server honors. Regenerate with:

  ```bash
  task parity:regen          # = go run ./parity/gen -server ../kubemq-server/config -out parity/server_env_keys.txt
  ```

  It needs a sibling `../kubemq-server` checkout. The golden is GENERATED — never hand-edit
  it. (Note the two naming schemes it folds: `convertEnvFormat` for most keys vs the
  dot-preserving `naturalEnvName` for `CLUSTER_REPLICATION_*`; both are copied verbatim from
  the server's `config/env.go` and must stay in sync.)

- **`parity/configdata_allowlist.txt`** — the **LEAN-CRD ledger**: every server key that is
  intentionally NOT a typed CRD field, and by what mechanism it is instead reachable
  (`spec.env`/`spec.configData`, or operator-computed). Adding a key here is a **deliberate
  LEAN-CRD decision**, not a TODO to type it later. `operator-computed`-tagged entries are
  cross-checked against the denylist (§2.2).

- **`parity_test.go`** (in `api/v1beta1/kubemqcluster/config`) is the gate. Every server key
  in the golden must be **either** emitted by some `SetConfig` (i.e. typed) **or** listed in
  the allowlist. Every allowlist entry must be a real server key (no stale entries), and no
  exact allowlist entry may ALSO be emitted by a `SetConfig` (no overlap — that would mean it
  should just be typed). `parity_denylist_coupling_test.go` (an external `config_test`
  package, to avoid the `v1beta1 → config` import cycle) asserts every `operator-computed`-
  tagged allowlist entry is covered by the §2.2 denylist.

**So the workflow when the server gains a new config key is:** run `task parity:regen`; the
parity test now fails with the new key. Decide per the §2.1 rule — type it (add a field +
deepcopy + `SetConfig` emit + CRD schema) **or** ledger it in `configdata_allowlist.txt` with
a disposition comment. Either way the test goes green; drift cannot ship silently.

### 2.4 CRD rule: never drop a storedVersion

The CRDs are multi-version. **`v1alpha1` is kept as the storage version** (bootstrap-deployed
clusters stored at it) and must never be dropped from a served/stored list — dropping a
version any object was stored at makes those objects unreadable and blocks `kubectl apply`.
`v1beta1` is the current served schema. New fields land in the v1beta1 schema **and** are
lockstep-mirrored into the v1alpha1 schema blocks (spec fields, enums — but NO CEL/status-
only additions in v1alpha1). The `Env`/`EnvFromSecrets`/`Expose` fields also exist as typed
Go struct fields on the deployer's v1alpha1 structs so a deployer-applied manifest is not a
silent field sink. Any schema edit is made in every copy in the **same change**; the CRD sync
tests compare each copy's per-version schema block as canonical JSON.

---

## 3. Conventions & operational notes

- **Module path** `github.com/kubemq-io/k8s`. Consumed by the operator via a `replace` to
  `../k8s`; **any change here requires a re-vendor in the operator** (`go mod vendor` there),
  or its build breaks with "inconsistent vendoring".
- **Deepcopy is hand-maintained**, not `make generate`d in-loop. A new field on a CRD type
  needs a matching hand-add in `api/v1beta1/zz_generated.deepcopy.go` (top-level types) and in
  the per-struct `DeepCopy()` methods under `kubemqcluster/config/` (each config struct owns
  its own `DeepCopy`, e.g. `StoreConfig.DeepCopy`). A field present in the type but missing
  from deepcopy is a silent data-loss bug.
- **A field is only real when type + deepcopy + `SetConfig` emit + all CRD schema copies
  agree.** Anything present in one but not the others is pruned by the API server.
- **Adding a `SetConfig` emit** = write the exact server env-var name the server binds
  (`SetConfigMapStringValues(config.Name, "SERVER_ENV_NAME", value)`); env-var names are a
  wire contract, locked by the parity test.
- **Tests:** `go test ./...` (`task test` runs with a 30m timeout). **Lint:** `task lint`
  (`golangci-lint run --disable gocritic --enable misspell`). **Vet:** `go vet ./...`.
- **Never auto-commit or push** — git is the user's explicit action.
