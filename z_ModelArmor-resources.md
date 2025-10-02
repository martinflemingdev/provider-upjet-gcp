# Vertex AI Resources - Implementation Tracking

## ADD first (Priority)

| Resource | externalname.go | externalnamenottested.go | create config.go's and AddResourceConfigurator | register in provider.go's | add group override | Notes
|----------|----------------|-------------------------|------------------------|-------|-------|-------|
| `google_model_armor_floorsetting` | ✅ | ✅ (not present) | ✅ | ✅ | ✅ | TemplatedStringAsIdentifier (singleton), parent+location required |
| `google_model_armor_template` | ✅ | ✅ (not present) | ✅ | ✅ | ✅ | TemplatedStringAsIdentifier, location required |
