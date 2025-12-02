# OpenAPI v2 to v3 Migration Guide

## Key Structural Changes

### Root Document
- `swagger: "2.0"` → `openapi: "3.x.x"`
- `host`, `basePath`, `schemes` → `servers[]`
- `consumes`, `produces` → moved to operation level or media types
- `definitions` → `components/schemas`
- `parameters` → `components/parameters`
- `responses` → `components/responses`
- `securityDefinitions` → `components/securitySchemes`

### Parameters
- `body` and `formData` parameters → `requestBody`
- Parameters now only: `path`, `query`, `header`, `cookie`
- Content types moved to `content` object with media types

### Responses
- Response content now under `content` with media types
- `schema` → `content[mediaType]/schema`

### Security
- Security scheme types updated
- OAuth2 flows restructured

### New Features in v3
- `servers` array for multiple server URLs
- `components` object for reusable components
- `requestBody` separate from parameters
- `webhooks` support
- `callbacks` support
- `links` in responses

## Migration Steps

1. Update root document structure
2. Convert host/basePath/schemes to servers
3. Move root-level definitions/parameters/responses to components
4. Extract body/formData parameters to requestBody
5. Update response schemas to use content/mediaType structure
6. Update security definitions structure
7. Convert all fixtures
8. Update all tests
