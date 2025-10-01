# Vertex AI Resources - Implementation Tracking

## ADD first (Priority)

| Resource | externalname.go | externalnamenottested.go | AddResourceConfigurator | Notes |
|----------|----------------|-------------------------|------------------------|-------|
| `google_vertex_ai_deployment_resource_pool` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, region required |
| `google_vertex_ai_endpoint` | ✅ | ✅ (commented out) | ✅ | TemplatedStringAsIdentifier, location required |
| `google_vertex_ai_endpoint_iam_member` | ✅ | ✅ (not present) | ✅ | IdentifierFromProvider (IAM), refs endpoint |
| `google_vertex_ai_endpoint_with_model_garden_deployment` | ✅ | ✅ (not present) | ✅ | IdentifierFromProvider (No Import), location required |
| `google_vertex_ai_index` | ✅ | ✅ (commented out) | ✅ | TemplatedStringAsIdentifier, region required |
| `google_vertex_ai_index_endpoint` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, region required |
| `google_vertex_ai_index_endpoint_deployed_index` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, refs index+endpoint, region required |
| `google_vertex_ai_metadata_store` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, region required |
| `google_vertex_ai_rag_engine_config` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier (region singleton), region required |

## ADD last (Future)

| Resource | externalname.go | externalnamenottested.go | AddResourceConfigurator | Notes |
|----------|----------------|-------------------------|------------------------|-------|
| `google_vertex_ai_feature_group` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, region required |
| `google_vertex_ai_feature_group_feature` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, refs feature_group, region required |
| `google_vertex_ai_feature_group_iam_member` | ✅ | ✅ (not present) | ✅ | IdentifierFromProvider (IAM), refs feature_group |
| `google_vertex_ai_feature_online_store` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, region required |
| `google_vertex_ai_feature_online_store_featureview` | ✅ | ✅ (not present) | ✅ | TemplatedStringAsIdentifier, refs feature_online_store, region required |
| `google_vertex_ai_feature_online_store_featureview_iam_member` | ✅ | ✅ (not present) | ✅ | IdentifierFromProvider (IAM), refs feature_online_store+feature_view |
| `google_vertex_ai_feature_online_store_iam_member` | ✅ | ✅ (not present) | ✅ | IdentifierFromProvider (IAM), refs feature_online_store |

## Included in current Provider

| Resource | externalname.go | externalnamenottested.go | AddResourceConfigurator | Notes |
|----------|----------------|-------------------------|------------------------|-------|
| `google_vertex_ai_dataset` | ✅ | N/A | ✅ | IdentifierFromProvider (No Import) |
| `google_vertex_ai_featurestore` | ✅ | N/A | ✅ | IdentifierFromProvider |
| `google_vertex_ai_featurestore_entitytype` | ✅ | N/A | ✅ | IdentifierFromProvider |
| `google_vertex_ai_tensorboard` | ✅ | N/A | ✅ | TemplatedStringAsIdentifier |

## Legacy (Not implementing)

| Resource | Status | Notes |
|----------|--------|-------|
| `google_vertex_ai_featurestore_entitytype_feature` | ❌ Legacy | Deprecated |
| `google_vertex_ai_featurestore_entitytype_iam` | ❌ Legacy | Deprecated |
| `google_vertex_ai_featurestore_iam` | ❌ Legacy | Deprecated |