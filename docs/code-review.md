# Codebase Review: Sundance SaaS Form Builder

## Issues Resolved by Recent Changes

1. ~~Most handlers had `// TODO: Send response`~~ -- `getForm`, `createForm`, `updateForm`, `publishVersion`, and `retireVersion` now all send proper JSON responses with DTOs.

2. ~~`getVersions`, `getVersion`, `createVersion`, `updateVersion` were empty stubs~~ -- All four handlers are now fully implemented with the goroutine/channel pattern, DTO mapping, and proper error handling.

3. ~~`FindVersions` and `FindVersion` missing from primary port~~ -- Added to the `FormsService` interface in `primary.go` and implemented in `form_service.go` with tenant isolation checks.

4. ~~Query objects lacked constructors~~ -- `FindByIDQuery`, `FindVersionsQuery`, and `FindVersionByIDQuery` now have `New*` constructors in `query.go`.

5. ~~DTOs were a single flat file~~ -- Refactored into a dedicated `dto/` sub-package with separate files per domain concept and bidirectional mapping functions (`DtoToPage`, `PageToResponseDto`, etc.).

6. ~~`CreateVersionCommand` required a pre-generated VersionID~~ -- Now only requires `FormID` and `TenantID`; the version ID is generated internally.

7. ~~Tenants DTOs mixed request and response types~~ -- Split into `request_dto.go` and `response_dto.go`, and response DTOs/mappers renamed to `*ResponseDto`/`*ToResponseDto`.

8. ~~Unused `removeVersion` handler/route~~ -- The `DELETE` route for versions was removed.

---

## Remaining Issues

### Bugs

1. **`getForms` sends `res.data` instead of `dtos`** (`handlers.go:51`) -- The DTO slice is built on lines 46-49, but line 51 still passes `res.data` to `SendJsonResponse`.

2. **`SendErrorResponse` maps `ErrExists` to HTTP 404 instead of 409** (`pkg/common/http.go`) -- The body says `StatusCode: 409` and `"Conflict"`, but the actual HTTP status code written is `404`.

3. **Handlers swallow JSON parse errors** -- `createForm` (line 97), `updateForm`, and `createVersion` (line 249) return without writing an HTTP response when `ReadJsonPayload` fails. Only `updateVersion` (line 293) correctly calls `SendErrorResponse`.

4. **`ErrExists` error message typo** (`pkg/common/error.go`): `"already exits"` should be `"already exists"`.

5. **`compose.yaml` copy-paste error**: Tenants service rebuild watch references `./services/submissions/go.mod` instead of `./services/tenants/go.mod`.

6. **Receiver name mismatch** in `section.go`: `UpdateFields` uses receiver `v` instead of `s`.

### Architectural

7. **`core.go` imports from the adapters layer** -- `ApplicationSettings` references `persistence.PersistanceSettings`, violating the hexagonal dependency rule. Configuration types should be defined in core or injected as plain values.

8. **Services receive the entire `Repository` struct** -- e.g., `TenantsService` gets `Database`, `Tenants`, AND `DataSources`. Violates Interface Segregation. Each service should receive only the interfaces it needs.

9. **REST handlers hold a reference to the full `Application`** -- The `handlers` struct takes `*core.Application` rather than just `*ports.Services`.

10. **`Find()` in forms service has no tenant filtering** -- Returns all forms across all tenants, unlike every other query which enforces tenant isolation.

11. **Aggregate boundaries unclear** -- `Form` has no `Versions` field; `Version` can be loaded/modified independently without going through `Form`. Inconsistent with `Tenant` which has a `DataSources` field.

### Missing Functionality

12. **Zero test files** in the entire repository.

13. **Domain validation unimplemented** -- All entity constructors contain `// TODO: Implement domain specific validation`.

14. **No domain events** for cross-service communication.

15. **Incomplete error-to-HTTP mapping** -- `ErrUnauthorized`, `ErrMissingTenantID`, and all domain errors (`ErrVersionLocked`, `ErrInvalidVersion`, etc.) fall through to the default 500 case.

16. **No real authentication** -- `X-Tenant-ID` is blindly trusted; `"placeholder"` hardcoded for user ID in publish/retire.

17. **No `Delete` operation for forms.**

18. **`ConditionalRule` is an empty stub** -- Contains only an ID; no rule type, conditions, or actions.

### Code Quality

19. **`FieldAttributes interface{}` is an empty interface** -- No type safety; any value can be assigned. Should use a sealed interface with a marker method.

20. **`DateFieldAttributes` missing `BaseFieldAttributes` embedding** -- Every other field attribute type embeds it; this one is missing `IsReadOnly` and `IsRequired`.

21. **`validator.New()` created per command constructor call** -- Expensive; should be a package-level singleton.

22. **Forms commands have no `validate` struct tags** -- `CreateFormCommand`, `UpdateFormCommand`, etc. have no tags, so the validator call is a no-op.

23. **Inconsistent constructor signatures** -- Forms constructors return `(*Entity, error)`; tenants constructors return just `*Entity` with no error.

24. **"Persistance" typo throughout** -- `PersistanceDriver`, `PersistanceSettings`, `PersistanceOptions`, and file names.

25. **Unnecessary goroutine/channel pattern** in every handler -- `net/http` already runs each handler in its own goroutine; the extra goroutine adds allocation overhead with no benefit.

26. **`time.Now()` called directly in service layer** -- Couples domain logic to wall-clock time; should be injected via a clock interface for testability.

27. **DTO organization still inconsistent** -- Forms uses exported types in a `dto/` sub-package; tenants uses unexported types inline in the `rest` package.

28. **`w.Write` error ignored** in `SendJsonResponse`.

29. **In-memory transactions are no-ops** -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing; the `CreateVersion` flow has a race condition since read and write locks are not held atomically.

---

## Summary

The project has a solid architectural foundation -- clean hexagonal layering, consistent package structure across services, well-defined ports, and a good CQRS-lite pattern with explicit command/query objects. The `Version` state machine is well-implemented, and the recent commits made significant progress: all forms service handlers are now functional, DTOs were properly refactored into a sub-package with bidirectional mapping, and query objects gained constructors.

However, several bugs remain in the handler/error-mapping layer (the `getForms` DTO bug, swallowed parse errors, wrong HTTP status for `ErrExists`). The biggest structural risks are zero test coverage, unimplemented domain validation, and missing error-to-HTTP mapping for domain/auth errors. The architecture has a few violations of its own principles (core importing adapters, overly broad service dependencies). Addressing the bugs, adding tests, implementing domain validation, and fixing the error mapping would be the highest-impact improvements.
