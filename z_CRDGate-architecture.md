# CRD Gate Architecture: Why 464 "gvk is ready" Messages in a KMS-Only Provider

## Date
April 13, 2026

## Context
After deploying `provider-gcp-kms` (Docker image `martinflemingdev/provider-gcp-kms-tf-v7.20.0:v1.0.0`) to a Kubernetes cluster, the pod log (`kms-pod.log`) showed **464 DEBUG "gvk is ready" messages** spanning many non-KMS API groups (compute, bigquery, cloudplatform, vertexai, monitoring, etc.).

The KMS family provider should only reconcile **8 KMS resource types** (cluster + namespaced = 16 controllers total). So why is it reporting readiness for 464 GVKs?

**TL;DR**: This is harmless. The CRD gate controller watches ALL CRDs on the cluster and logs "gvk is ready" for every established one. Only the 16 KMS controllers actually start. The other 448 `gate.Set()` calls are no-ops with no registered listener.

---

## Table of Contents
1. [The Three Source Files](#1-the-three-source-files)
2. [Source File 1: The CRD Gate Reconciler](#2-source-file-1-the-crd-gate-reconciler)
3. [Source File 2: The Gate Primitive](#3-source-file-2-the-gate-primitive)
4. [Source File 3: The Generated KMS Controller](#4-source-file-3-the-generated-kms-controller)
5. [The Wiring: zz_main.go](#5-the-wiring-zz_maingo)
6. [The Family Setup Files](#6-the-family-setup-files)
7. [Complete Data Flow](#7-complete-data-flow)
8. [Pod Log Analysis](#8-pod-log-analysis)
9. [Why ALL CRDs?](#9-why-all-crds)
10. [Key Takeaways](#10-key-takeaways)
11. [References](#11-references)

---

## 1. The Three Source Files

| # | File | Repo | Role |
|---|------|------|------|
| 1 | `pkg/reconciler/customresourcesgate/reconciler.go` | crossplane-runtime v2 | Watches ALL CRDs, calls `gate.Set(gvk, true/false)` for each, logs "gvk is ready" |
| 2 | `pkg/gate/gate.go` | crossplane-runtime v2 | Generic gated callback registry — `Register(fn, depends...)` + `Set(condition, bool)` |
| 3 | `internal/controller/cluster/kms/cryptokey/zz_controller.go` | provider-upjet-gcp (generated) | Individual controller that calls `gate.Register(Setup, myGVK)` to delay start until CRD exists |

Supporting files:
| File | Repo | Role |
|------|------|------|
| `pkg/reconciler/customresourcesgate/setup.go` | crossplane-runtime v2 | Sets up the CRD-watching controller |
| `pkg/controller/gate.go` | crossplane-runtime v2 | Defines the `Gate` interface |
| `pkg/pipeline/templates/controller.go.tmpl` | upjet v2 | Template that generates the `SetupGated` pattern |
| `cmd/provider/kms/zz_main.go` | provider-upjet-gcp (generated) | Wires everything together |
| `internal/controller/cluster/zz_kms_setup.go` | provider-upjet-gcp (generated) | Calls `SetupGated` for each KMS controller |

---

## 2. Source File 1: The CRD Gate Reconciler

**Repo**: `github.com/crossplane/crossplane-runtime/v2`
**File**: `pkg/reconciler/customresourcesgate/reconciler.go`

This is the file that **produces the "gvk is ready" log messages**.

### setup.go — How the controller is created

```go
// pkg/reconciler/customresourcesgate/setup.go (lines 29-46)

// Setup adds a controller that reconciles CustomResourceDefinitions to support
// delayed start of controllers.
// o.Gate is expected to be something like *gate.Gate[schema.GroupVersionKind].
func Setup(mgr ctrl.Manager, o controller.Options) error {
    if o.Gate == nil || reflect.ValueOf(o.Gate).IsNil() {
        return errors.New("gate is required")
    }

    r := &Reconciler{
        log:  o.Logger,
        gate: o.Gate,
    }

    return ctrl.NewControllerManagedBy(mgr).
        For(&apiextensionsv1.CustomResourceDefinition{}).  // ← WATCHES ALL CRDs
        Named("crd-gate").
        Complete(reconcile.AsReconciler[*apiextensionsv1.CustomResourceDefinition](
            mgr.GetClient(), r))
}
```

**Key point**: `.For(&apiextensionsv1.CustomResourceDefinition{})` means this single controller watches **every CRD** on the cluster. Every time a CRD is created, updated, or deleted, this reconciler fires.

### reconciler.go — The Reconcile function that logs "gvk is ready"

```go
// pkg/reconciler/customresourcesgate/reconciler.go (lines 27-65)

// Reconciler reconciles a CustomResourceDefinitions in order to gate and wait
// on CRD readiness to start downstream controllers.
type Reconciler struct {
    log  logging.Logger
    gate controller.Gate
}

// Reconcile reconciles CustomResourceDefinitions and reports ready and unready
// GVKs to the gate.
func (r *Reconciler) Reconcile(
    _ context.Context,
    crd *apiextensionsv1.CustomResourceDefinition,
) (ctrl.Result, error) {
    established := isEstablished(crd)
    gkvs := toGVKs(crd)

    switch {
    // CRD is not ready or being deleted.
    case !established || !crd.GetDeletionTimestamp().IsZero():
        for gvk := range gkvs {
            r.log.Debug("gvk is not ready", "gvk", gvk)
            r.gate.Set(gvk, false)
        }
        return ctrl.Result{}, nil

    // CRD is ready.
    default:
        for gvk, served := range gkvs {
            if served {
                r.log.Debug("gvk is ready", "gvk", gvk)  // ← THIS IS THE LOG LINE
                r.gate.Set(gvk, true)                     // ← signals the gate
            }
        }
    }

    return ctrl.Result{}, nil
}
```

**Key point**: For **every established CRD** on the cluster, it:
1. Logs `"gvk is ready"` at DEBUG level
2. Calls `r.gate.Set(gvk, true)`

It does NOT check whether anyone is listening for that GVK. It blindly reports readiness for everything.

### Helper functions

```go
// pkg/reconciler/customresourcesgate/reconciler.go (lines 66-86)

func toGVKs(crd *apiextensionsv1.CustomResourceDefinition) map[schema.GroupVersionKind]bool {
    gvks := make(map[schema.GroupVersionKind]bool, len(crd.Spec.Versions))
    for _, version := range crd.Spec.Versions {
        gvks[schema.GroupVersionKind{
            Group:   crd.Spec.Group,
            Version: version.Name,
            Kind:    crd.Spec.Names.Kind,
        }] = version.Served
    }
    return gvks
}

func isEstablished(crd *apiextensionsv1.CustomResourceDefinition) bool {
    if len(crd.Status.Conditions) > 0 {
        for _, cond := range crd.Status.Conditions {
            if cond.Type == apiextensionsv1.Established {
                return cond.Status == apiextensionsv1.ConditionTrue
            }
        }
    }
    return false
}
```

**Key point**: `toGVKs` returns a map keyed by GVK with `bool` indicating whether the version is served. A CRD with 3 versions (v1beta1, v1beta2, v1) produces 3 GVKs. This is why some resources appear multiple times in the log — they have multiple served versions.

---

## 3. Source File 2: The Gate Primitive

**Repo**: `github.com/crossplane/crossplane-runtime/v2`
**File**: `pkg/gate/gate.go`

This is a **generic, type-safe gated callback registry**. It has nothing to do with Kubernetes — it's a pure concurrency primitive.

```go
// pkg/gate/gate.go (full file)

package gate

import (
    "slices"
    "sync"
)

// Gate implements a gated function callback registration with comparable conditions.
type Gate[T comparable] struct {
    mux       sync.RWMutex
    satisfied map[T]bool
    fns       []gated[T]
}

// gated is an internal tracking resource.
type gated[T comparable] struct {
    fn       func()      // callback to invoke when all deps are true
    depends  []T         // conditions this function is waiting on (AND logic)
    released bool        // true = already called, pending garbage collection
}

// Register a callback function that will be called when all the provided
// dependent conditions are true. After all conditions are true, the callback
// function is removed from the registration and will not be called again.
func (g *Gate[T]) Register(fn func(), depends ...T) {
    g.mux.Lock()
    g.fns = append(g.fns, gated[T]{fn: fn, depends: depends})
    g.mux.Unlock()
    g.process()
}

// Set marks the associated condition to the given value. If the condition is
// already set as that value, then this is a no-op.
// Returns true if there was an update detected.
func (g *Gate[T]) Set(condition T, value bool) bool {
    g.mux.Lock()

    if g.satisfied == nil {
        g.satisfied = make(map[T]bool)
    }

    old, found := g.satisfied[condition]

    updated := false
    if !found || old != value {
        updated = true
        g.satisfied[condition] = value  // ← just writes to a map
    }
    g.mux.Unlock()

    if updated {
        g.process()  // ← only if something changed
    }

    return updated
}

func (g *Gate[T]) process() {
    g.mux.Lock()
    defer g.mux.Unlock()

    for i := range g.fns {
        release := true

        for _, dep := range g.fns[i].depends {
            if !g.satisfied[dep] {
                release = false  // ← at least one dep not satisfied
            }
        }

        if release {
            fn := g.fns[i].fn
            g.fns[i].released = true
            go fn()  // ← fires the callback in a goroutine
        }
    }

    // garbage collect released functions
    g.fns = slices.DeleteFunc(g.fns, func(a gated[T]) bool {
        return a.released
    })
}
```

### How it works for KMS

The Gate is instantiated as `*gate.Gate[schema.GroupVersionKind]` in `zz_main.go`.

**Registration phase** (at startup):
- 16 KMS controllers call `gate.Register(setupCallback, kmsGVK)` — one per controller
- Each registers a callback function + the single GVK it depends on
- The gate's `fns` slice now has 16 entries

**Discovery phase** (when CRDs are found):
- The CRD gate reconciler calls `gate.Set(gvk, true)` for **every** CRD it discovers
- For non-KMS GVKs (448 of them): `Set()` writes `satisfied[gvk] = true`, then `process()` iterates the `fns` slice. No fn depends on this GVK, so nothing fires. **This is essentially a map write + a short loop — negligible cost.**
- For KMS GVKs (16 of them): `Set()` writes `satisfied[kmsGVK] = true`, `process()` finds a matching callback, fires it in a goroutine, and garbage-collects the entry.

---

## 4. Source File 3: The Generated KMS Controller

**Repo**: `github.com/upbound/provider-gcp/v2` (this repo, generated by Upjet)
**File**: `internal/controller/cluster/kms/cryptokey/zz_controller.go`

This is auto-generated by Upjet from the template `pkg/pipeline/templates/controller.go.tmpl`.

### SetupGated — registers callback with the gate

```go
// SetupGated adds a controller that reconciles CryptoKey managed resources.
func SetupGated(mgr ctrl.Manager, o tjcontroller.Options) error {
    o.Options.Gate.Register(func() {
        if err := Setup(mgr, o); err != nil {
            mgr.GetLogger().Error(err, "unable to setup reconciler",
                "gvk", v1beta1.CryptoKey_GroupVersionKind.String())
        }
    }, v1beta1.CryptoKey_GroupVersionKind)  // ← depends on THIS one GVK
    return nil
}
```

**Key point**: `SetupGated` does NOT start the controller immediately. It registers a callback with the gate that says: "When `kms.gcp.upbound.io/v1beta1 CryptoKey` is marked as ready, THEN call `Setup()`."

### Setup — actually starts the controller

```go
// Setup adds a controller that reconciles CryptoKey managed resources.
func Setup(mgr ctrl.Manager, o tjcontroller.Options) error {
    name := managed.ControllerName(v1beta1.CryptoKey_GroupVersionKind.String())
    // ... initializers, event handler, API callbacks ...

    opts := []managed.ReconcilerOption{
        managed.WithExternalConnecter(
            tjcontroller.NewTerraformPluginSDKAsyncConnector(
                mgr.GetClient(),
                o.OperationTrackerStore,
                o.SetupFn,
                o.Provider.Resources["google_kms_crypto_key"],  // ← TF resource
                // ...
            )),
        // ... logger, recorder, finalizer, timeout, poll ...
    }

    r := managed.NewReconciler(mgr,
        xpresource.ManagedKind(v1beta1.CryptoKey_GroupVersionKind), opts...)

    return ctrl.NewControllerManagedBy(mgr).
        Named(name).
        WithOptions(o.ForControllerRuntime()).
        WithEventFilter(xpresource.DesiredStateChanged()).
        Watches(&v1beta1.CryptoKey{}, eventHandler).  // ← ONLY watches CryptoKey CRs
        Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}
```

**Key point**: `Setup()` is what actually registers a controller-runtime controller that watches `CryptoKey` custom resources. This function only runs AFTER the gate confirms the CRD exists.

### The Upjet Template

All generated controllers follow the same pattern from `pkg/pipeline/templates/controller.go.tmpl` in the Upjet repo:

```go
// Template excerpt (crossplane/upjet pkg/pipeline/templates/controller.go.tmpl)
func SetupGated(mgr ctrl.Manager, o tjcontroller.Options) error {
    o.Options.Gate.Register(func() {
        if err := Setup(mgr, o); err != nil {
            mgr.GetLogger().Error(err, "unable to setup reconciler",
                "gvk", {{ .CRD.GroupVersionKind }}.String())
        }
    }, {{ .CRD.GroupVersionKind }})
    return nil
}
```

Every Upjet-generated controller gets this identical pattern — the only variable is the GVK.

---

## 5. The Wiring: zz_main.go

**File**: `cmd/provider/kms/zz_main.go`

This is the KMS family provider entry point. The critical sections:

### Scheme Registration (ALL APIs)

```go
// zz_main.go lines 48-50 — imports
clusterapis    "github.com/upbound/provider-gcp/v2/apis/cluster"
namespacedapis "github.com/upbound/provider-gcp/v2/apis/namespaced"

// zz_main.go lines 175-179 — scheme registration
kingpin.FatalIfError(clusterapis.AddToScheme(mgr.GetScheme()), ...)
kingpin.FatalIfError(namespacedapis.AddToScheme(mgr.GetScheme()), ...)
```

**Why ALL types are registered to the scheme**: Cross-resource references. A KMS `CryptoKey` might reference a `cloudplatform.Project` or an `iam.ServiceAccount`. For the reference resolver to work, those types must be known to the runtime scheme. This is the monolith scheme pattern — all types registered, but only some controllers started.

### Controller Imports (KMS only)

```go
// zz_main.go lines 54-55 — imports
clustercontroller    "github.com/upbound/provider-gcp/v2/internal/controller/cluster"
namespacedcontroller "github.com/upbound/provider-gcp/v2/internal/controller/namespaced"
```

### Gate Creation and Wiring

```go
// zz_main.go lines 261-267

canSafeStart, err := canWatchCRD(ctx, mgr)
// ...
if canSafeStart {
    crdGate := new(gate.Gate[schema.GroupVersionKind])  // ← one shared gate instance

    clusterOpts.Gate = crdGate       // ← passed to KMS controllers
    namespacedOpts.Gate = crdGate    // ← passed to KMS controllers

    // Sets up the CRD-watching controller (ONE controller for ALL CRDs)
    kingpin.FatalIfError(customresourcesgate.Setup(mgr, namespacedOpts.Options), ...)

    // Registers 8 cluster KMS gated callbacks
    kingpin.FatalIfError(clustercontroller.SetupGated_kms(mgr, clusterOpts), ...)

    // Registers 8 namespaced KMS gated callbacks
    kingpin.FatalIfError(namespacedcontroller.SetupGated_kms(mgr, namespacedOpts), ...)
} else {
    // Fallback: start controllers immediately without gating
    kingpin.FatalIfError(clustercontroller.Setup_kms(mgr, clusterOpts), ...)
    kingpin.FatalIfError(namespacedcontroller.Setup_kms(mgr, namespacedOpts), ...)
}
```

**Key point**: There is ONE `gate.Gate` instance shared across all controllers. The CRD gate reconciler calls `Set()` on it for every CRD. The 16 KMS controllers register callbacks on it for their specific GVKs.

---

## 6. The Family Setup Files

### Cluster: `internal/controller/cluster/zz_kms_setup.go`

```go
func SetupGated_kms(mgr ctrl.Manager, o controller.Options) error {
    for _, setup := range []func(ctrl.Manager, controller.Options) error{
        cryptokey.SetupGated,
        cryptokeyiammember.SetupGated,
        cryptokeyversion.SetupGated,
        keyhandle.SetupGated,        // ← our new resource!
        keyring.SetupGated,
        keyringiammember.SetupGated,
        keyringimportjob.SetupGated,
        secretciphertext.SetupGated,
    } {
        if err := setup(mgr, o); err != nil {
            return err
        }
    }
    return nil
}
```

### Namespaced: `internal/controller/namespaced/zz_kms_setup.go`

Identical structure, 8 namespaced KMS controllers.

**Total**: 8 cluster + 8 namespaced = **16 `gate.Register()` calls**, each for one KMS GVK.

---

## 7. Complete Data Flow

```
STARTUP
═══════

zz_main.go
    │
    ├─ clusterapis.AddToScheme()          Register ALL GCP types to scheme
    ├─ namespacedapis.AddToScheme()       (for cross-resource references)
    │
    ├─ crdGate := new(gate.Gate[GVK])     Create single Gate instance
    │
    ├─ customresourcesgate.Setup()         Start ONE CRD-watching controller
    │   └─ .For(&CRD{})                   Watches ALL CRDs via controller-runtime
    │
    ├─ clustercontroller.SetupGated_kms()
    │   ├─ cryptokey.SetupGated()         → gate.Register(Setup, CryptoKey_GVK)
    │   ├─ cryptokeyversion.SetupGated()  → gate.Register(Setup, CryptoKeyVersion_GVK)
    │   ├─ keyhandle.SetupGated()         → gate.Register(Setup, KeyHandle_GVK)
    │   └─ ... (8 total)
    │
    └─ namespacedcontroller.SetupGated_kms()
        ├─ cryptokey.SetupGated()         → gate.Register(Setup, CryptoKey_GVK)
        └─ ... (8 total)

    Gate state after startup:
    ┌──────────────────────────────────────────────────┐
    │ fns: [16 entries, each waiting on 1 KMS GVK]     │
    │ satisfied: {}  (empty — no CRDs discovered yet)  │
    └──────────────────────────────────────────────────┘


RUNTIME (controller-manager starts)
════════════════════════════════════

CRD Gate Reconciler receives events for ALL CRDs on cluster
    │
    ├─ CRD: compute.gcp.upbound.io/Instance (established)
    │   └─ Reconcile():
    │       ├─ log.Debug("gvk is ready", "gvk", "compute.../Instance")  ← LOG LINE
    │       └─ gate.Set(compute_Instance_GVK, true)
    │           └─ process(): scans fns[] → no callback depends on this → NO-OP
    │
    ├─ CRD: kms.gcp.upbound.io/CryptoKey (established)
    │   └─ Reconcile():
    │       ├─ log.Debug("gvk is ready", "gvk", "kms.../CryptoKey")     ← LOG LINE
    │       └─ gate.Set(CryptoKey_GVK, true)
    │           └─ process(): finds matching callback → go Setup(mgr, o) → CONTROLLER STARTS
    │
    ├─ CRD: bigquery.gcp.upbound.io/Dataset (established)
    │   └─ Reconcile():
    │       ├─ log.Debug("gvk is ready", ...)                            ← LOG LINE
    │       └─ gate.Set(bigquery_Dataset_GVK, true) → NO-OP
    │
    └─ ... (464 total CRDs × 1 log line each = 464 "gvk is ready" messages)

    Result:
    ┌──────────────────────────────────────────────────────────┐
    │ 464 DEBUG log lines (harmless)                           │
    │ 464 gate.Set() calls (448 no-ops + 16 trigger callbacks)│
    │ 16 controllers started (8 cluster KMS + 8 namespaced)   │
    │ 0 controllers for compute, bigquery, vertexai, etc.     │
    └──────────────────────────────────────────────────────────┘
```

---

## 8. Pod Log Analysis

### Log distribution by API group (from `kms-pod.log`)

| API Group | Count | Notes |
|-----------|-------|-------|
| `compute.gcp.upbound.io` | 138 | Cluster-scoped compute CRDs |
| `compute.gcp.m.upbound.io` | 96 | Namespaced compute CRDs |
| `bigquery.gcp.upbound.io` | 34 | |
| `cloudplatform.gcp.upbound.io` | 20 | |
| `bigquery.gcp.m.upbound.io` | 20 | |
| `vertexai.gcp.upbound.io` | 19 | |
| `monitoring.gcp.upbound.io` | 16 | |
| `cloudplatform.gcp.m.upbound.io` | 16 | |
| `vertexai.gcp.m.upbound.io` | 15 | |
| `pkg.crossplane.io` | 12 | Crossplane core CRDs |
| **`kms.gcp.upbound.io`** | **12** | **Cluster-scoped KMS** (8 types × some with v1beta2) |
| `monitoring.gcp.m.upbound.io` | 9 | |
| **`kms.gcp.m.upbound.io`** | **8** | **Namespaced KMS** (8 types × v1beta1 only) |
| `apiextensions.crossplane.io` | 8 | Crossplane core CRDs |
| `secretmanager.gcp.upbound.io` | 5 | |
| `dataform.gcp-beta.upbound.io` | 4 | |
| `dataform.gcp-beta.m.upbound.io` | 4 | |
| ... | ... | (other groups with 1-3 entries each) |
| **Total** | **464** | |

### KMS log lines (20 total)

The 20 KMS entries break down as:
- **12 cluster-scoped** (`kms.gcp.upbound.io`): 8 resources, 4 of which have both v1beta1 and v1beta2 versions = 12 GVKs
- **8 namespaced** (`kms.gcp.m.upbound.io`): 8 resources × v1beta1 only = 8 GVKs

The cluster-scoped resources with v1beta2 versions (4 of them):
```
kms.gcp.upbound.io/v1beta1, Kind=CryptoKey
kms.gcp.upbound.io/v1beta2, Kind=CryptoKey        ← v1beta2
kms.gcp.upbound.io/v1beta1, Kind=CryptoKeyIAMMember
kms.gcp.upbound.io/v1beta2, Kind=CryptoKeyIAMMember  ← v1beta2
kms.gcp.upbound.io/v1beta1, Kind=CryptoKeyVersion
kms.gcp.upbound.io/v1beta2, Kind=CryptoKeyVersion    ← v1beta2
kms.gcp.upbound.io/v1beta1, Kind=KeyRingIAMMember
kms.gcp.upbound.io/v1beta2, Kind=KeyRingIAMMember    ← v1beta2
```

### Sample log output

```
2026-04-13T22:49:20Z  DEBUG  provider-gcp  Starting  {"sync-interval":"1h0m0s","poll-interval":"10m0s",...}
2026-04-13T22:49:25Z  INFO   provider-gcp  Beta feature enabled  {"flag":"EnableBetaManagementPolicies"}
2026-04-13T22:49:26Z  DEBUG  provider-gcp  gvk is ready  {"gvk":"dataform.gcp-beta.upbound.io/v1beta1, Kind=Repository"}
2026-04-13T22:49:26Z  DEBUG  provider-gcp  gvk is ready  {"gvk":"kms.gcp.upbound.io/v1beta1, Kind=SecretCiphertext"}
2026-04-13T22:49:26Z  DEBUG  provider-gcp  gvk is ready  {"gvk":"cloudplatform.gcp.upbound.io/v1beta1, Kind=ServiceNetworkingPeeredDNSDomain"}
...
(464 lines total, all within the same second — 22:49:26Z)
...
2026-04-13T22:49:26Z  DEBUG  provider-gcp  gvk is ready  {"gvk":"cloudplatform.gcp.upbound.io/v1beta2, Kind=ProjectIAMMember"}
```

All 464 messages appear within **one second** during the initial CRD cache sync.

---

## 9. Why ALL CRDs?

### Q: Can we make the CRD gate reconciler only watch KMS CRDs?

**No**. The CRD gate reconciler is designed as a generic, provider-agnostic component in `crossplane-runtime`. It uses:

```go
ctrl.NewControllerManagedBy(mgr).
    For(&apiextensionsv1.CustomResourceDefinition{}).  // watches ALL CRDs
```

There's no CRD-level label selector or field selector applied. It needs to be generic because:
1. It lives in `crossplane-runtime`, not in any specific provider
2. Different providers may share the same gate instance
3. CRDs can be installed/removed dynamically — the gate needs to react to CRD deletion too (`"gvk is not ready"`)

### Q: Is this wasteful?

**No**. The cost breakdown:
- **CRD watch**: One informer watches CRDs. This would exist anyway for any Kubernetes controller that needs CRDs.
- **Reconcile per CRD**: ~microseconds per CRD — reads status, calls `gate.Set()`, returns.
- **gate.Set()**: Lock a mutex, write to a map, scan a 16-entry slice, unlock. Negligible.
- **Log output**: 464 DEBUG lines at startup. Not logged at INFO or above. Only visible because you ran with debug logging.

### Q: Why does the KMS provider register ALL API types to the scheme?

Because of **cross-resource references**. Example: a KMS `CryptoKey` might have a `spec.forProvider.keyRingRef` that references a `KeyRing`, and the reference resolver needs to know the KeyRing's Go type. More broadly, any managed resource might reference resources from other API groups (projects, service accounts, networks, etc.).

The scheme registration in `zz_main.go`:
```go
clusterapis.AddToScheme(mgr.GetScheme())     // ALL cluster types
namespacedapis.AddToScheme(mgr.GetScheme())  // ALL namespaced types
```

This adds Go types to the runtime scheme. It does NOT create controllers, informers, or watches for those types. It just makes the types known to the serialization/deserialization layer.

---

## 10. Key Takeaways

| Aspect | What it looks like | What actually happens |
|--------|-------------------|----------------------|
| 464 "gvk is ready" | "Provider watching 464 resources!" | CRD gate acknowledges 464 CRDs exist on cluster |
| `gate.Set(gvk, true)` for non-KMS | "Starting controllers for compute!" | Map write + short loop scan → no-op |
| `gate.Set(gvk, true)` for KMS | "Starting KMS controller" | Fires callback → `Setup()` → real controller starts |
| `AddToScheme()` for all APIs | "Loading all GCP resources!" | Registers Go types for serialization, no controllers |
| Total controllers started | 464? | **16** (8 cluster KMS + 8 namespaced KMS) |

### The architecture in one sentence

> The CRD gate watches all CRDs and reports their readiness to a shared gate, but only controllers that registered callbacks for specific GVKs actually start — the rest are no-ops.

---

## 11. References

### Source Code (crossplane-runtime v2)

| File | URL |
|------|-----|
| `pkg/reconciler/customresourcesgate/reconciler.go` | [github.com/crossplane/crossplane-runtime/.../reconciler.go](https://github.com/crossplane/crossplane-runtime/blob/main/pkg/reconciler/customresourcesgate/reconciler.go) |
| `pkg/reconciler/customresourcesgate/setup.go` | [github.com/crossplane/crossplane-runtime/.../setup.go](https://github.com/crossplane/crossplane-runtime/blob/main/pkg/reconciler/customresourcesgate/setup.go) |
| `pkg/gate/gate.go` | [github.com/crossplane/crossplane-runtime/.../gate.go](https://github.com/crossplane/crossplane-runtime/blob/main/pkg/gate/gate.go) |
| `pkg/controller/gate.go` (interface) | [github.com/crossplane/crossplane-runtime/.../controller/gate.go](https://github.com/crossplane/crossplane-runtime/blob/main/pkg/controller/gate.go) |

### Source Code (upjet v2)

| File | URL |
|------|-----|
| `pkg/pipeline/templates/controller.go.tmpl` | [github.com/crossplane/upjet/.../controller.go.tmpl](https://github.com/crossplane/upjet/blob/main/pkg/pipeline/templates/controller.go.tmpl) |

### Source Code (this repo — provider-upjet-gcp)

| File | Path |
|------|------|
| KMS provider entry point | `cmd/provider/kms/zz_main.go` |
| Cluster KMS setup | `internal/controller/cluster/zz_kms_setup.go` |
| Namespaced KMS setup | `internal/controller/namespaced/zz_kms_setup.go` |
| Example generated controller | `internal/controller/cluster/kms/cryptokey/zz_controller.go` |

### Versions

| Component | Version |
|-----------|---------|
| crossplane-runtime | `v2.0.0-20250730220209-c306b1c8b181` |
| upjet | `v2.0.1-0.20251028081228-8d73164bb9bd` |
| provider-upjet-gcp | `tf-7.20.0` branch |
| Terraform google provider | `v7.20.0` |
