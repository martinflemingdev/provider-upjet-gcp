# Vertex AI Resources - Implementation Tracking

## ADD first (Priority)

| Resource | externalname.go | externalnamenottested.go | create config.go and AddResourceConfigurator | Add .Configure call to GetProvider() | Notes
|----------|----------------|-------------------------|------------------------|-------|-------|
| `google_model_armor_floorsetting` | ✅ | ✅ (not present) | ⬜ | ⬜ | TemplatedStringAsIdentifier (singleton), parent+location required |
| `google_model_armor_template` | ✅ | ✅ (not present) | ⬜ | ⬜ | TemplatedStringAsIdentifier, location required |
