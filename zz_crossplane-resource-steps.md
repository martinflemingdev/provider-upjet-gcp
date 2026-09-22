# Steps to Add a New Resource to an Upjet-based Crossplane Provider

source: https://github.com/crossplane/upjet/blob/main/docs/adding-new-resource.md

## 1. Add External Name Configuration
**Edit:** `config/externalname.go`  

**Action:**  
Add an entry for your resource with either `config.IdentifierFromProvider` or `config.TemplatedStringAsIdentifier`, based on the import format described in the Terraform docs.  

**Reference:**  
Use the “Import” section from the Terraform Registry for the exact import path template.  

> **Note:** You may have to add a new section to externalname.go for a new service.


Example:
```go
// Imported by using the following projects/{{project}}/locations/{{location}}/endpoints/{{name}}
"google_vertex_ai_endpoint": config.TemplatedStringAsIdentifier(
  "name",
  "projects/{{ .setup.configuration.project }}/locations/{{ .parameters.location }}/endpoints/{{ .external_name }}",
),

// Imported by using the following projects/{{project}}/locations/{{location}}/endpoints/{{name}} roles/viewer user:jane@example.com
"google_vertex_ai_endpoint_iam_member": config.IdentifierFromProvider,
```
## 2. Remove or Comment Out in Not-Tested File
**Edit**: `config/externalnamenottested.go`

**Action**:
If your resource appears here, comment it out or remove it.
This allows Upjet to include and generate it.

## 3. Handle TF Registry Warnings and Mutually Exclusive Fields (If Needed) and All IAM Resources
**Where**: Check the “Warning”/“Note” boxes on the Terraform Registry page for your resource.

**Edit**: `config/namespace|cluster/<service>/config.go` (e.g., config/namespace/bigquery/config.go, config/cluster/vertexai/config.go)

**Action**:

If a field in the resource should not be managed directly (e.g., because there’s a separate CRD for it), move that field to status:

```go
p.AddResourceConfigurator("resource_name", func(r *config.Resource) {
  config.MoveToStatus(r.TerraformResource, "field_name")
})
```
  
Add Crossplane references (r.References) if the resource has fields referring to other managed resources.

**Always do this for IAM resources.**

Reference schema:
https://github.com/hashicorp/terraform-provider-google/tree/main/google/services

Example:

```go
p.AddResourceConfigurator("google_vertex_ai_endpoint_iam_member", func(r *config.Resource) {
  r.References["endpoint"] = config.Reference{
    TerraformName: "google_vertex_ai_endpoint",
  }
})
```

### 3.1 If the resource is for a new service, do these additional steps:

1. **Create the config files for both scopes:**
   - `config/cluster/<service>/config.go`
   - `config/namespaced/<service>/config.go`

2. **Add `AddResourceConfigurator` blocks** in both config files for each resource in the service:
   ```go
   func Configure(p *config.Provider) {
       p.AddResourceConfigurator("google_<service>_<resource>", func(r *config.Resource) {
           // Add any resource-specific configuration here
           // e.g., config.MarkAsRequired(r.TerraformResource, "field_name")
       })
   }
   ```

3. **Import and register in provider.go files:**
   - Add the import in both `config/cluster/provider.go` and `config/namespaced/provider.go`:
     ```go
     "github.com/upbound/provider-gcp/config/cluster/<service>"  // or /namespaced/<service>
     ```
   - Add the `Configure` call to the `init()` function in both files:
     ```go
     ProviderConfiguration.AddConfig(<service>.Configure)
     ```
   - The configuration will be automatically included in the loops in `registry_cluster.go` and `registry_namespaced.go` via the `cluster.ProviderConfiguration` and `namespaced.ProviderConfiguration` slices

4. **Add group override for multi-word service names** (if applicable):
   - If your service has multiple words in its name (e.g., `google_model_armor`, `google_vertex_ai`), add an entry to the `groupMap` in `config/overrides.go`:
     ```go
     "google_<service_name>.+": ReplaceGroupWords("", 2),
     ```
   - This ensures the correct group name is used (e.g., "modelarmor" instead of "model", "vertexai" instead of "vertex")
   - The number `2` indicates how many underscore-separated words to skip after "google_" when calculating the group name


## 4. Run Code Generation and Commit
Terminal commands:

```bash
make submodules
make generate
git add .
git commit -m "make generate"
git push
```
  
## 5. Validate Types and CRD
Check generated files.

Validate the new CRD(s) appear in the **package/crds/** directory and check their GVK.  

### Manually Configuring Cross-Resource References

If certain resource fields didn't generate Ref and Selector fields automatically, you can force Upjet to generate them by adding a reference configuration in the resource's `config.go` file.

**When is this needed?**
Upjet's automatic reference detection uses heuristics based on field naming patterns and descriptions. It may miss references for:
- Singular fields nested in complex structures (e.g., `all_updates_rule.pubsub_topic`)
- Fields where the name doesn't clearly match the resource type
- Fields with non-standard naming conventions

**How to add manual references:**

1. Open the appropriate config file (e.g., `config/cluster/billing/config.go`)

2. Add a `References` configuration to the resource configurator:

```go
p.AddResourceConfigurator("google_billing_budget", func(r *config.Resource) {
    // ... existing config ...
    
    // Cross-resource reference to Pub/Sub Topic
    // Use dot notation for nested fields: "parent.child.field"
    r.References["all_updates_rule.pubsub_topic"] = config.Reference{
        TerraformName: "google_pubsub_topic",
        Extractor:     "github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()",
    }
})
```

Run `make generate` to regenerate the CRDs with the new Ref/Selector fields

Verify the generation:
Check the generated type file (e.g., zz_budget_types.go) to confirm the reference fields were added:

```go
PubsubTopic *string `json:"pubsubTopic,omitempty" tf:"pubsub_topic,omitempty"`

// Reference to a Topic in pubsub to populate pubsubTopic.
PubsubTopicRef *v1.Reference `json:"pubsubTopicRef,omitempty" tf:"-"`

// Selector for a Topic in pubsub to populate pubsubTopic.
PubsubTopicSelector *v1.Selector `json:"pubsubTopicSelector,omitempty" tf:"-"`
```

## 6. Build, Push and Test Provider
Load and push provider package to a registry for deployment:


```bash
make build-provider.<service>
```

Validate **\_output/xpkg files** and locate the .xpkg file to docker load.

```bash
docker load -i /home/martinfleming/src/github.com/martinflemingdev/provider-upjet-gcp/_output/xpkg/linux_amd64/provider-gcp-bigquery-v0.0.0-1357.g5381bec1.dirty.xpkg

docker tag sha256:6636c90cd5a8eab8f3e6a08c049bacf1313b0a1e4dff8936a7b6c17a5b4bd0b3 \
  martinflemingdev/provider-gcp-modelarmor-tf-v6.50.0:v1.0.0

GA
docker push martinflemingdev/provider-gcp-modelarmor-tf-v6.50.0:v1.0.0

Beta
docker push martinflemingdev/provider-gcp-beta-dataform-tf-v6.48.0:v1.0.0
  ```