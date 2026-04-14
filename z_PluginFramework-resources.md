# Implementing Terraform Plugin Framework Resources in Upjet Providers

> **Status**: Research / Planning  
> **Date**: 2026-04-13
> **Branch**: `tf-7.20.0`  
> **Context**: Bumping `hashicorp/terraform-provider-google` from v6.47.0 → v7.20.0 exposed that two resources migrated from Plugin SDK to Plugin Framework in v7.0.0, causing a runtime panic.

---

## Table of Contents

1. [Problem Statement](#problem-statement)
2. [Root Cause](#root-cause)
3. [Background: SDK vs Framework Resources](#background-sdk-vs-framework-resources)
4. [GCP TF Provider Framework Resources](#gcp-tf-provider-framework-resources)
5. [Upjet v2 Framework Support](#upjet-v2-framework-support)
6. [Working Reference: provider-upjet-aws](#working-reference-provider-upjet-aws)
7. [Cross-Provider Comparison: AWS vs Azure vs GCP](#cross-provider-comparison-aws-vs-azure-vs-gcp)
8. [Implementation Steps for provider-upjet-gcp](#implementation-steps-for-provider-upjet-gcp)
9. [Blockers and Open Questions](#blockers-and-open-questions)
10. [Temporary Workaround](#temporary-workaround)
11. [File Reference Matrix](#file-reference-matrix)

---

## Problem Statement

After bumping the Terraform GCP provider to v7.20.0, the `kms` family provider panicked at startup:

```
panic: the resource google_apigee_keystores_aliases_key_cert_file is configured
to be reconciled with Terraform Plugin SDK but the resource's Go schema does not exist
```

The resource `google_apigee_keystores_aliases_key_cert_file` was migrated from the **Terraform Plugin SDK** to the **Terraform Plugin Framework** in the GCP TF provider v7.0.0. Upjet's SDK-based reconciler cannot find its schema because the schema now only exists in the Framework provider.

---

## Root Cause

Terraform providers have been migrating resources from the legacy **Plugin SDK** (`terraform-plugin-sdk/v2`) to the newer **Plugin Framework** (`terraform-plugin-framework`). When a resource is migrated:

- Its `schema.Resource` struct is removed from `Provider().ResourcesMap` (SDK)
- A new `resource.Resource` implementation is registered via `provider.Resources()` (Framework)
- The Go schema that Upjet's SDK connector (`TerraformPluginSDKConnector`) expects **no longer exists**

Upjet resolves this at startup: it tries to find the resource in `sdkProvider.ResourcesMap`, fails, and panics.

---

## Background: SDK vs Framework Resources

| Aspect | Plugin SDK | Plugin Framework |
|--------|-----------|-----------------|
| Provider interface | `*schema.Provider` | `provider.Provider` (implements `fwprovider.Provider`) |
| Resource definition | `*schema.Resource` in `ResourcesMap` | `resource.Resource` via `Resources()` method |
| Schema | `schema.Schema` (SDK) | `resource.Schema` (Framework) |
| Upjet connector | `TerraformPluginSDKConnector` | `TerraformPluginFrameworkConnector` |
| Protocol | terraform-plugin-go `tfprotov5` | terraform-plugin-go `tfprotov6` |
| Upjet config option | `WithTerraformPluginSDKIncludeList()` | `WithTerraformPluginFrameworkIncludeList()` |

---

## GCP TF Provider Framework Resources

As of `hashicorp/terraform-provider-google` v7.20.0, there are exactly **2 Framework resources**, registered in [`google/fwprovider/framework_provider.go`](https://github.com/hashicorp/terraform-provider-google/blob/main/google/fwprovider/framework_provider.go):

```go
func (p *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        apigee.NewApigeeKeystoresAliasesKeyCertFileResource,
        storage.NewStorageNotificationResource,
    }
}
```

Terraform resource names:
1. `google_apigee_keystores_aliases_key_cert_file`
2. `google_storage_notification`

---

## Upjet v2 Framework Support

Upjet v2 **does** have full support for Terraform Plugin Framework resources. Key APIs in `github.com/crossplane/upjet/v2`:

| API | Location | Purpose |
|-----|----------|---------|
| `config.WithTerraformPluginFrameworkProvider(fwProvider)` | `pkg/config` | Registers the Framework provider instance |
| `config.WithTerraformPluginFrameworkIncludeList(list)` | `pkg/config` | Specifies which resources use the Framework connector |
| `config.FrameworkResourceWithComputedIdentifier()` | `pkg/config` | External name config for Framework resources with computed IDs |
| `config.ExternalName.TFPluginFrameworkOptions` | `pkg/config` | Framework-specific external name settings |
| `TerraformPluginFrameworkConnector` | `pkg/controller` | The reconciler that manages Framework resources via tfprotov6 |

---

## Working Reference: provider-upjet-aws

The **[crossplane-contrib/provider-upjet-aws](https://github.com/crossplane-contrib/provider-upjet-aws)** provider has a complete, working implementation of Framework resource support. Below is a file-by-file walkthrough.

### Key Architectural Difference: `xpprovider` Package

The **critical** difference between the AWS and GCP providers is where the TF providers are instantiated:

| | AWS Provider | GCP Provider (current) |
|---|---|---|
| **Import** | `"github.com/hashicorp/terraform-provider-aws/xpprovider"` | `"github.com/hashicorp/terraform-provider-google/google/provider"` |
| **Call** | `fwProvider, sdkProvider, err := xpprovider.GetProvider(ctx)` | `sdkProvider := provider.Provider()` |
| **Returns** | Both `fwprovider.Provider` + `*schema.Provider` | Only `*schema.Provider` |

The `xpprovider` package lives **inside the upstream terraform-provider-aws repo** and provides a unified entry point that constructs both the Framework and SDK provider instances. The GCP TF provider does **not** have an equivalent `xpprovider` package.

### AWS File: `config/externalname.go`

**Lines 0-33**: Defines the `TerraformPluginFrameworkExternalNameConfigs` map alongside the existing SDK map:

```go
// TerraformPluginFrameworkExternalNameConfigs contains all external name
// configurations belonging to Terraform Plugin Framework resources to be
// reconciled under the no-fork architecture.
var TerraformPluginFrameworkExternalNameConfigs = map[string]config.ExternalName{
    // auditmanager
    "aws_auditmanager_account_registration": config.FrameworkResourceWithComputedIdentifier(),
    // ...
}
```

> Source: [`config/externalname.go` L0-33](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/externalname.go#L1-L33)

**Lines 3292-3309**: The `ResourceConfigurator()` function applies with **precedence**: Framework > SDK > CLI:

```go
func ResourceConfigurator() config.ResourceOption {
    return func(r *config.Resource) {
        if e, ok := TerraformPluginFrameworkExternalNameConfigs[r.Name]; ok {
            r.ExternalName = e
        } else if e, ok := terraformPluginSDKExternalNameConfigs[r.Name]; ok {
            r.ExternalName = e
        } else if e, ok := CLIReconciledExternalNameConfigs[r.Name]; ok {
            r.ExternalName = e
        }
    }
}
```

> Source: [`config/externalname.go` L3292-3309](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/externalname.go#L3292-L3309)

**Lines 3586-3603**: Helper function `frameworkNameAsIdentifier()`:

```go
func frameworkNameAsIdentifier() config.ExternalName {
    e := config.FrameworkNameAsIdentifier("name")
    // ...
    return e
}
```

> Source: [`config/externalname.go` L3586-3603](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/externalname.go#L3586-L3603)

### AWS File: `config/registry_common.go`

**Lines 90-100**: `TerraformPluginFrameworkResourceList()` function, parallel to `TerraformPluginSDKResourceList()`:

```go
func TerraformPluginFrameworkResourceList() []string {
    l := make([]string, len(TerraformPluginFrameworkExternalNameConfigs))
    i := 0
    for n := range TerraformPluginFrameworkExternalNameConfigs {
        l[i] = n + "$"
        i++
    }
    return l
}
```

> Source: [`config/registry_common.go` L90-100](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/registry_common.go#L90-L100)

### AWS File: `config/registry_cluster.go`

**Lines 101-126**: `GetProvider()` accepts **both** providers and passes Framework options to `config.NewProvider`:

```go
func GetProvider(ctx context.Context, fwProvider fwprovider.Provider,
    sdkProvider *schema.Provider, generationProvider bool, skipDefaultTags bool) (*config.Provider, error) {
    // ...
    pc := config.NewProvider([]byte(providerSchema), "aws",
        modulePath, providerMetadata,
        // ...
        config.WithTerraformPluginSDKIncludeList(TerraformPluginSDKResourceList()),
        config.WithTerraformPluginFrameworkIncludeList(TerraformPluginFrameworkResourceList()),
        // ...
        config.WithTerraformProvider(sdkProvider),
        config.WithTerraformPluginFrameworkProvider(fwProvider),
        // ...
    )
    // ...
}
```

> Source: [`config/registry_cluster.go` L101-126](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/registry_cluster.go#L101-L126)

### AWS File: `config/registry_namespaced.go`

Same pattern as `registry_cluster.go`:

```go
func GetProviderNamespaced(ctx context.Context, fwProvider fwprovider.Provider,
    sdkProvider *schema.Provider, generationProvider bool, skipDefaultTags bool) (*config.Provider, error) {
    // ...
    pc := config.NewProvider([]byte(providerSchema), "aws",
        modulePath, providerMetadata,
        // ...
        config.WithTerraformPluginSDKIncludeList(TerraformPluginSDKResourceList()),
        config.WithTerraformPluginFrameworkIncludeList(TerraformPluginFrameworkResourceList()),
        // ...
        config.WithTerraformProvider(sdkProvider),
        config.WithTerraformPluginFrameworkProvider(fwProvider),
        // ...
    )
    // ...
}
```

> Source: [`config/registry_namespaced.go` L59-84](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/registry_namespaced.go#L59-L84)

### AWS File: `hack/main.go.tmpl`

The template that generates `cmd/provider/*/zz_main.go`:

```go
import (
    // ...
    "github.com/hashicorp/terraform-provider-aws/xpprovider"
    // ...
)

func main() {
    // ...
    fwProvider, sdkProvider, err := xpprovider.GetProvider(ctx)
    kingpin.FatalIfError(err, "Cannot get the Terraform framework and SDK providers")
    clusterProvider, err := config.GetProvider(ctx, fwProvider, sdkProvider, false, *skipDefaultTags)
    // ...
}
```

> Source: [`hack/main.go.tmpl` L195](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/hack/main.go.tmpl#L195)

### AWS File: `cmd/generator/main.go`

```go
import (
    "github.com/hashicorp/terraform-provider-aws/xpprovider"
    // ...
)

func main() {
    // ...
    fwProvider, sdkProvider, err := xpprovider.GetProvider(ctx)
    pc, err := config.GetProvider(context.Background(), fwProvider, sdkProvider, true, false)
    pns, err := config.GetProviderNamespaced(context.Background(), fwProvider, sdkProvider, true, false)
    // ...
}
```

> Source: [`cmd/generator/main.go` L38-41](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/cmd/generator/main.go#L38-L41)

---

## Cross-Provider Comparison: AWS vs Azure vs GCP

Researching all three major Upjet providers reveals distinct architectural patterns. Only **AWS** has implemented Framework resource support. Understanding the differences informs what GCP needs to adopt.

### Architecture Summary

| Aspect | AWS | Azure | GCP (current) |
|--------|-----|-------|---------------|
| **TF Provider Version** | Latest | v4.54.0 | v7.20.0 |
| **`go.mod replace`** | `upbound/terraform-provider-aws` | `upbound/terraform-provider-azurerm` | None (uses `hashicorp` directly) |
| **`xpprovider` Package** | ✅ In Upbound's fork | ✅ In Upbound's fork (SDK-only) | ❌ None |
| **`xpprovider` Returns** | `(fwProvider, sdkProvider, err)` | `(*schema.Provider, error)` | N/A — calls `provider.Provider()` |
| **Framework Resources** | ✅ Full support | ❌ Not needed (all SDK at v4.54.0) | ❌ Not implemented (2 resources affected) |
| **External Name Maps** | SDK + Framework + CLI | SDK + CLI | SDK only |
| **Resource Tiers** | Framework > SDK > CLI | SDK > CLI | SDK only |
| **Upjet Version** | v2.2.1-0.20251128 | v2.2.1-0.20251217 | v2.0.1 |

### The `go.mod replace` Pattern

All three providers import the upstream TF provider, but AWS and Azure redirect to **Upbound-maintained forks** via `go.mod replace`:

**AWS** (`go.mod`):
```
replace github.com/hashicorp/terraform-provider-aws =>
    github.com/upbound/terraform-provider-aws v0.0.0-20260305123303-f7691456b787
```

**Azure** (`go.mod`):
```
replace github.com/hashicorp/terraform-provider-azurerm =>
    github.com/upbound/terraform-provider-azurerm v0.0.0-20251127122522-9029e3f708c4
```

**GCP** (`go.mod`):
```
github.com/hashicorp/terraform-provider-google v1.20.1-0.20260217183254-340083741c96
```
No `replace` directive — imports `hashicorp/terraform-provider-google` directly.

The forks contain the `xpprovider` package that the Upjet provider imports. **GCP's lack of a fork means there's nowhere for an `xpprovider` package to live** (in the upstream-import style). This is why Step 1 in the implementation plan proposes creating a local wrapper.

### Azure Deep Dive: `xpprovider` (SDK-Only)

Azure's fork contains an `xpprovider` package generated from [`hack/provider.go.txt`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/hack/provider.go.txt):

```go
package xpprovider

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

func Provider() *schema.Provider {
    return provider.AzureProvider()
}
```

This wraps the SDK provider **only** — no Framework provider. The generated `zz_main.go` files call:

```go
sdkProvider, err := xpprovider.GetProviderSchema(context.Background())
// ... only sdkProvider, no fwProvider
clusterProvider, err := config.GetProvider(ctx, sdkProvider, false)
```

> Source: [`hack/main.go.tmpl`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/hack/main.go.tmpl)

### Azure: CLI-Reconciled Resources

A notable Azure pattern **not present in GCP** is the CLI-reconciled resource tier. Azure defines two external name maps:

```go
var TerraformPluginSDKExternalNameConfigs = map[string]config.ExternalName{ /* ... */ }
var CLIReconciledExternalNameConfigs = map[string]config.ExternalName{ /* ... */ }
```

And two corresponding resource list functions:

```go
func TerraformPluginSDKResourceList() []string { /* ... */ }
func CLIReconciledResourceList() []string { /* ... */ }
```

The `ResourceConfigurator()` applies with **SDK > CLI** precedence:

```go
func ResourceConfigurator() config.ResourceOption {
    return func(r *config.Resource) {
        if e, ok := TerraformPluginSDKExternalNameConfigs[r.Name]; ok {
            r.ExternalName = e
        } else if e, ok := CLIReconciledExternalNameConfigs[r.Name]; ok {
            r.ExternalName = e
        }
    }
}
```

CLI-reconciled resources fall back to executing `terraform` CLI commands instead of direct SDK calls. This could be useful for GCP resources that are problematic with direct SDK reconciliation.

> Source: [`config/externalname.go`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/config/externalname.go), [`config/registry_common.go`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/config/registry_common.go)

### Azure: Other Notable Patterns

| Pattern | Description | GCP Equivalent |
|---------|-------------|----------------|
| `hack/provider.go.txt` | Template placed into fork's `xpprovider/` dir | None — could adopt if GCP creates a fork |
| `common.RemoveIndex(r.ExternalName.IdentifierFields, "field")` | Removes fields from external name identifier list | Not used in GCP |
| `bumpVersionsWithEmbeddedLists()` | Singleton list → embedded object migration for API compatibility | GCP has `old-singleton-list-apis.txt` but different approach |
| `SchemaElementOptions.SetInitProviderOverrides()` with `TagOverrides` | Fine-grained init provider overrides | Not used in GCP |

### Key Takeaway

To implement Framework support for GCP, the **AWS pattern is the model to follow** (not Azure). Specifically:

1. **AWS** is the only Upjet provider with working Framework plumbing
2. **Azure** is in the same boat as GCP (SDK-only) but hasn't needed Framework support yet because Azure TF provider v4.54.0 hasn't migrated any resources
3. **GCP** is unique in having **no fork and no `xpprovider` package at all**, making it the least scaffolded of the three
4. GCP could additionally adopt the **CLI-reconciled tier from Azure** as a third reconciliation strategy for problematic resources

---

## Implementation Steps for provider-upjet-gcp

### Step 1: Create a `GetProvider` Function That Returns Both Providers

The GCP TF provider currently has:
- **SDK provider**: `provider.Provider()` in `github.com/hashicorp/terraform-provider-google/google/provider`
- **Framework provider**: `fwprovider.New()` in `github.com/hashicorp/terraform-provider-google/google/fwprovider`

**Option A**: Create a local `xpprovider` package in the Upjet GCP provider that wraps both:

```go
// internal/xpprovider/provider.go
package xpprovider

import (
    "context"

    fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    googlefwprovider "github.com/hashicorp/terraform-provider-google/google/fwprovider"
    googleprovider "github.com/hashicorp/terraform-provider-google/google/provider"
)

func GetProvider(ctx context.Context) (fwprovider.Provider, *schema.Provider, error) {
    sdkProvider := googleprovider.Provider()
    fwProvider := googlefwprovider.New() // or whatever the actual constructor is
    return fwProvider, sdkProvider, nil
}
```

**Option B**: Wait for the upstream `hashicorp/terraform-provider-google` to add an `xpprovider` package (like `terraform-provider-aws` has).

> **⚠️ IMPORTANT**: Investigate the actual GCP Framework provider constructor. The AWS TF provider exports this via `xpprovider.GetProvider()` which handles internal initialization. The GCP TF provider's `fwprovider` package may need specific initialization (e.g., provider version string, transport configuration).

### Step 2: Add `TerraformPluginFrameworkExternalNameConfigs` to `config/externalname.go`

```go
// TerraformPluginFrameworkExternalNameConfigs contains all external name
// configurations belonging to Terraform Plugin Framework resources.
var TerraformPluginFrameworkExternalNameConfigs = map[string]config.ExternalName{
    // Apigee
    "google_apigee_keystores_aliases_key_cert_file": config.IdentifierFromProvider,
    // Storage
    "google_storage_notification": config.IdentifierFromProvider,
}
```

> **Note**: The specific `ExternalName` strategy for each resource needs to be determined by examining how the resource's ID is computed. `config.FrameworkResourceWithComputedIdentifier()` may be appropriate, or a `TemplatedStringAsIdentifier()` could work if the ID format is known. Use `config.IdentifierFromProvider` as a safe fallback.

### Step 3: Add `TerraformPluginFrameworkResourceList()` to `config/registry_common.go`

```go
// TerraformPluginFrameworkResourceList returns the list of resources that
// use the Terraform Plugin Framework connector.
func TerraformPluginFrameworkResourceList() []string {
    l := make([]string, len(TerraformPluginFrameworkExternalNameConfigs))
    i := 0
    for n := range TerraformPluginFrameworkExternalNameConfigs {
        l[i] = n + "$"
        i++
    }
    return l
}
```

### Step 4: Update `config/registry_cluster.go` — `GetProvider()`

Change the function signature to accept both providers:

```go
// Before:
func GetProvider(_ context.Context, sdkProvider *schema.Provider, generationProvider bool) (*ujconfig.Provider, error) {

// After:
func GetProvider(_ context.Context, fwProvider fwprovider.Provider, sdkProvider *schema.Provider, generationProvider bool) (*ujconfig.Provider, error) {
```

Add Framework options to `ujconfig.NewProvider`:

```go
pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix,
    modulePath, providerMetadata,
    // ... existing options ...
    ujconfig.WithTerraformPluginSDKIncludeList(resourceList(terraformPluginSDKExternalNameConfigs)),
    ujconfig.WithTerraformPluginFrameworkIncludeList(TerraformPluginFrameworkResourceList()),  // NEW
    // ... existing options ...
    ujconfig.WithTerraformProvider(sdkProvider),
    ujconfig.WithTerraformPluginFrameworkProvider(fwProvider),  // NEW
)
```

Add the new import:

```go
fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
```

### Step 5: Update `config/registry_namespaced.go` — `GetNamespacedProvider()`

Same changes as Step 4:

```go
// Before:
func GetNamespacedProvider(_ context.Context, sdkProvider *schema.Provider, generationProvider bool) (*ujconfig.Provider, error) {

// After:
func GetNamespacedProvider(_ context.Context, fwProvider fwprovider.Provider, sdkProvider *schema.Provider, generationProvider bool) (*ujconfig.Provider, error) {
```

And add the same `WithTerraformPluginFrameworkIncludeList` + `WithTerraformPluginFrameworkProvider` options.

### Step 6: Update `hack/main.go.tmpl`

Change the import and provider instantiation:

```go
// Before:
"github.com/hashicorp/terraform-provider-google/google/provider"

// After (Option A - local xpprovider):
xpprovider "github.com/upbound/provider-gcp/v2/internal/xpprovider"
```

Change the provider initialization:

```go
// Before:
sdkProvider := provider.Provider()
clusterProvider, err := config.GetProvider(ctx, sdkProvider, false)

// After:
fwProvider, sdkProvider, err := xpprovider.GetProvider(ctx)
kingpin.FatalIfError(err, "Cannot get the Terraform framework and SDK providers")
clusterProvider, err := config.GetProvider(ctx, fwProvider, sdkProvider, false)
```

### Step 7: Update `cmd/generator/main.go`

Same changes as Step 6:

```go
// Before:
sdkProvider := provider.Provider()
pc, err := config.GetProvider(context.Background(), sdkProvider, true)
pns, err := config.GetNamespacedProvider(context.Background(), sdkProvider, true)

// After:
fwProvider, sdkProvider, err := xpprovider.GetProvider(context.Background())
kingpin.FatalIfError(err, "Cannot get the Terraform framework and SDK providers")
pc, err := config.GetProvider(context.Background(), fwProvider, sdkProvider, true)
pns, err := config.GetNamespacedProvider(context.Background(), fwProvider, sdkProvider, true)
```

### Step 8: Add `ResourceConfigurator()` Precedence Logic

In `config/externalname.go` or a config override, ensure Framework external names take precedence. See the AWS `ResourceConfigurator()` pattern:

```go
func ResourceConfigurator() config.ResourceOption {
    return func(r *config.Resource) {
        if e, ok := TerraformPluginFrameworkExternalNameConfigs[r.Name]; ok {
            r.ExternalName = e
        } else if e, ok := terraformPluginSDKExternalNameConfigs[r.Name]; ok {
            r.ExternalName = e
        }
    }
}
```

### Step 9: Uncomment the Framework Resources

Un-comment the two resources in `config/externalname.go` and move them from `terraformPluginSDKExternalNameConfigs` to `TerraformPluginFrameworkExternalNameConfigs`.

Un-comment the configurators in `config/cluster/apigee/config.go` and `config/namespaced/apigee/config.go` (and equivalent storage config files if they exist).

### Step 10: Run `make generate` and `make build`

After all changes:

```bash
make generate
make build-provider.kms    # Test a family provider
make build-provider.apigee # Test the apigee provider specifically
```

---

## Blockers and Open Questions

### 1. GCP TF Provider Does Not Expose an `xpprovider` Package

The AWS TF provider (`hashicorp/terraform-provider-aws`) has a dedicated `xpprovider` package that provides:

```go
func GetProvider(ctx context.Context) (fwprovider.Provider, *schema.Provider, error)
```

The GCP TF provider (`hashicorp/terraform-provider-google`) does **not** have this. We need to:

- **Investigate** the `google/fwprovider` package exports to understand how to construct the Framework provider
- **Determine** if additional initialization (version string, transport, credentials setup) is required
- **Create** a local wrapper function that properly constructs both providers

### 2. Framework Provider Initialization Complexity

The Framework provider may require initialization with specific configuration (e.g., the provider version string for User-Agent, credential/transport helpers). This needs careful investigation of:

- `google/fwprovider/framework_provider.go` — the `New()` or constructor function
- What `Configure()` method expects at runtime
- Whether the Upjet `TerraformPluginFrameworkConnector` handles this automatically

### 3. External Name Strategy for Framework Resources

Each Framework resource needs the correct `ExternalName` configuration. For the 2 GCP Framework resources:

- `google_apigee_keystores_aliases_key_cert_file` — ID format needs investigation
- `google_storage_notification` — ID format needs investigation

Use `config.IdentifierFromProvider` as a safe default, or investigate `config.FrameworkResourceWithComputedIdentifier()`.

### 4. Only 2 Resources Currently Affected

As of GCP TF provider v7.20.0, only 2 resources use the Framework. This is a small surface area, which makes the cost/benefit of implementing full Framework support debatable for now. However, HashiCorp is actively migrating more resources to Framework, so this number will grow with future provider versions.

---

## Temporary Workaround

Both Framework resources have been **commented out** from `terraformPluginSDKExternalNameConfigs` in `config/externalname.go` and their configurators disabled:

- `config/externalname.go` — `google_apigee_keystores_aliases_key_cert_file` commented out
- `config/externalname.go` — `google_storage_notification` commented out  
- `config/cluster/apigee/config.go` — configurator commented out
- `config/namespaced/apigee/config.go` — configurator commented out

This means these 2 resources will not be generated as CRDs. They can be re-enabled once the Framework provider plumbing described above is implemented.

---

## File Reference Matrix

| File (GCP) | File (AWS Reference) | File (Azure Reference) | What Changes |
|------------|---------------------|----------------------|-------------|
| `config/externalname.go` | [`config/externalname.go`](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/externalname.go) | [`config/externalname.go`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/config/externalname.go) | Add `TerraformPluginFrameworkExternalNameConfigs` map |
| `config/registry_common.go` | [`config/registry_common.go`](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/registry_common.go) | [`config/registry_common.go`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/config/registry_common.go) | Add `TerraformPluginFrameworkResourceList()` |
| `config/registry_cluster.go` | [`config/registry_cluster.go`](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/registry_cluster.go) | [`config/registry_cluster.go`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/config/registry_cluster.go) | Add `fwProvider` param + Framework options |
| `config/registry_namespaced.go` | [`config/registry_namespaced.go`](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/config/registry_namespaced.go) | [`config/registry_namespaced.go`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/config/registry_namespaced.go) | Add `fwProvider` param + Framework options |
| `hack/main.go.tmpl` | [`hack/main.go.tmpl`](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/hack/main.go.tmpl) | [`hack/main.go.tmpl`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/hack/main.go.tmpl) | Import xpprovider, call `GetProvider()` for both |
| `cmd/generator/main.go` | [`cmd/generator/main.go`](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/cmd/generator/main.go) | [`cmd/generator/main.go`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/cmd/generator/main.go) | Import xpprovider, call `GetProvider()` for both |
| `go.mod` | [`go.mod`](https://github.com/crossplane-contrib/provider-upjet-aws/blob/main/go.mod) | [`go.mod`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/go.mod) | Add `replace` directive to point to Upbound fork |
| *(new)* `internal/xpprovider/provider.go` | N/A (lives in TF provider fork) | [`hack/provider.go.txt`](https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/hack/provider.go.txt) | Create local wrapper to construct both providers |

---

## References

- **Upjet Adding a New Resource**: https://github.com/crossplane/upjet/blob/main/docs/adding-new-resource.md
- **provider-upjet-aws (working Framework impl)**: https://github.com/crossplane-contrib/provider-upjet-aws
- **provider-upjet-azure (SDK + CLI, no Framework)**: https://github.com/crossplane-contrib/provider-upjet-azure
- **Upjet Framework Connector Source**: https://github.com/crossplane/upjet/tree/main/pkg/controller
- **GCP TF Provider Framework Resources**: https://github.com/hashicorp/terraform-provider-google/blob/main/google/fwprovider/framework_provider.go
- **AWS TF Provider xpprovider Package**: https://github.com/hashicorp/terraform-provider-aws (search for `xpprovider/` directory)
- **Azure hack/provider.go.txt**: https://github.com/crossplane-contrib/provider-upjet-azure/blob/main/hack/provider.go.txt
- **Upbound AWS TF Fork (contains xpprovider)**: https://github.com/upbound/terraform-provider-aws
- **Upbound Azure TF Fork (contains xpprovider, SDK-only)**: https://github.com/upbound/terraform-provider-azurerm
