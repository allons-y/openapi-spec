# OpenAPI v2 to v3 Conversion Status

## Completed

### Core Type Definitions ✅
- ✅ Created `server.go` with `Server` and `ServerVariable` structs
- ✅ Created `components.go` with `Components`, `Link`, and `Callback` structs
- ✅ Created `request_body.go` with `RequestBody` struct
- ✅ Created `media_type.go` with `MediaType`, `Encoding`, and `Example` structs
- ✅ Updated `swagger.go`:
  - Changed `SwaggerProps` from v2 to v3 structure
  - Replaced `Swagger`, `Host`, `BasePath`, `Schemes`, `Consumes`, `Produces` with `OpenAPI`, `Servers`
  - Replaced `Definitions`, `Parameters`, `Responses`, `SecurityDefinitions` with `Components`
  - Added `Webhooks` and `JSONSchemaDialect`

### Updated Structs ✅
- ✅ Updated `OperationProps` in `operation.go`:
  - Removed `Consumes`, `Produces`, `Schemes`
  - Added `RequestBody`, `Callbacks`, `Servers`
  - Deprecated helper methods `WithConsumes()` and `WithProduces()`

- ✅ Updated `ParamProps` in `parameter.go`:
  - Added OpenAPI v3 fields: `Deprecated`, `Style`, `Explode`, `AllowReserved`, `Example`, `Examples`, `Content`
  - Updated documentation to reflect only 4 parameter types (path, query, header, cookie)
  - Deprecated `BodyParam()`, `FormDataParam()`, `FileParam()` helpers
  - Added `CookieParam()` helper

- ✅ Updated `ResponseProps` in `response.go`:
  - Added `Content` map (MediaType)
  - Added `Links` map
  - Kept `Schema` and `Examples` for backward compatibility (marked as deprecated)

### Build System ✅
- ✅ Code compiles successfully
- ✅ Updated `expander.go` to work with `Components` instead of root-level definitions/parameters/responses
- ✅ Updated `spec.go` constants to prioritize OpenAPI v3

### Documentation ✅
- ✅ Created `MIGRATION_V3.md` with migration guidance
- ✅ Updated `README.md` to reflect OpenAPI v3 support
- ✅ Updated inline documentation and comments

### Sample Conversion ✅
- ✅ Converted `fixtures/local_expansion/spec.json` to OpenAPI v3 format

## Remaining Work

### Fixture Conversion ⚠️
- ⚠️ 140+ fixture files still need conversion from v2 to v3
- Files need updates:
  - `swagger: "2.0"` → `openapi: "3.0.0"`
  - `host`/`basePath`/`schemes` → `servers[]`
  - `definitions` → `components/schemas`
  - `parameters` → `components/parameters`
  - `responses` → `components/responses`
  - `securityDefinitions` → `components/securitySchemes`
  - Response `schema` → `content` with media type
  - Body/form parameters → `requestBody`

### Test Updates ⚠️
- ⚠️ All test files need updates to use OpenAPI v3 format
- Known failing test files:
  - `operation_test.go` - references `Consumes`, `Produces`, `Schemes`
  - `swagger_test.go` - references `Swagger`, `Host`, `BasePath`, `Definitions`, etc.
  - Many other test files likely affected

### Additional Struct Updates 🔄
- 🔄 `security.go` - SecurityScheme needs v3 updates (OAuth2 flows restructured)
- 🔄 `header.go` - May need Content field like parameters
- 🔄 Path item structs may need updates
- 🔄 Schema validation rules may need updates

### Reference Resolution 🔄
- 🔄 Update ref resolution to handle `#/components/*` paths instead of `#/definitions/*`
- 🔄 Update schema loader to recognize v3 structure
- 🔄 Update circular reference handling for v3

### Validator Updates 🔄
- 🔄 Validation rules need updating for v3 structure
- 🔄 Required field validation (OpenAPI, Info, Paths)
- 🔄 Parameter location validation (only path/query/header/cookie)
- 🔄 Media type validation

## Quick Start for Completing Migration

### 1. Update Tests (Priority: HIGH)
```bash
# Find all test files referencing old v2 fields
grep -r "Swagger:" . --include="*_test.go"
grep -r "Consumes:" . --include="*_test.go"
grep -r "Definitions:" . --include="*_test.go"
# Update each file systematically
```

### 2. Convert Fixtures (Priority: HIGH)
```bash
# Use a script or manual conversion
# Key conversions needed in every fixture:
# - swagger → openapi
# - host/basePath/schemes → servers
# - definitions → components/schemas
# - securityDefinitions → components/securitySchemes
# - response schema → content
```

### 3. Update Security (Priority: MEDIUM)
- Update SecurityScheme for v3 OAuth2 flows
- Update security reference paths

### 4. Test and Validate (Priority: HIGH)
```bash
go test ./...
```

## Breaking Changes

This is a MAJOR version change that breaks backward compatibility:

1. **Root document structure completely changed**
2. **All fixture files must be updated**
3. **Downstream packages using this library will need updates**
4. **Helper functions deprecated (WithConsumes, WithProduces, BodyParam, etc.)**

## Migration Path for Users

Users of this library should:
1. Update their OpenAPI specs from v2 to v3 format
2. Update code to use new struct fields:
   - Use `spec.Components.Schemas` instead of `spec.Definitions`
   - Use `operation.RequestBody` instead of body parameters
   - Use `response.Content` instead of `response.Schema`
   - Use `spec.Servers` instead of `host`/`basePath`/`schemes`
3. Update validation and analysis code for v3 structure
4. Update code generation templates for v3

## Estimated Remaining Effort

- **Test updates**: 4-8 hours
- **Fixture conversion**: 6-12 hours (could be automated with script)
- **Security updates**: 2-4 hours
- **Reference resolution updates**: 4-6 hours
- **Validation updates**: 4-6 hours
- **Testing and bug fixes**: 8-16 hours

**Total**: 28-52 hours of work remaining
