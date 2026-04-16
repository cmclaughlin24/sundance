# Tenants Service Code Review

## Issues Resolved Since 4/13 Review

1. ~~`SendErrorResponse` maps `ErrExists` to HTTP 404 instead of 409~~ (#2) -- `SendErrorResponse` now correctly passes `http.StatusConflict` as the HTTP status code (`http.go:61`).

2. ~~`ErrExists` error message typo `"already exits"`~~ (#4) -- Fixed to `"already exists"` (`error.go:7`).

3. ~~`compose.yaml` copy-paste error: tenants service rebuild watch references `./services/submissions/go.mod`~~ (#5) -- Now correctly references `./services/tenants/go.mod` (`compose.yaml:55`).

4. ~~Tenants DTOs mixed request and response types~~ (#7) -- Split into `dto/request.go` and `dto/response.go` in a dedicated `dto/` sub-package with exported types and mapper functions.

5. ~~DTO organization inconsistent between forms and tenants~~ (#27) -- Tenants now uses the same pattern as forms: exported types in a `dto/` sub-package with separate request/response files.

---

## Remaining Issues

### Bugs

1. **Handlers swallow JSON parse errors** (`handlers.go:82-84`, `119-121`, `238-240`, `276-278`) -- `createTenant`, `updateTenant`, `createDataSource`, and `updateDataSource` all `return` without writing an HTTP response when `ReadJsonPayload` fails. The client receives an empty response with a 200 status code.

2. **Services silently discard `domain.New*` constructor errors** (`tenants_service.go:32`, `43`; `data_sources_service.go:32`, `43`) -- Both services assign the error from `domain.NewTenant` / `domain.NewDataSource` to `_`. If domain constructors ever gain validation logic (they return `error` for a reason), those errors will be silently lost.

3. **`Remove` operations succeed silently for non-existent entities** (`tenant_repository.go:78-84`; `data_sources_repository.go:84-91`) -- Go's `delete` on a map is a no-op for missing keys, so deleting a non-existent tenant or data source returns `nil` (HTTP 204) instead of `ErrNotFound` (HTTP 404).

4. **`updateTenant` success message typo** (`handlers.go:145`) -- Says `"Success updated!"` instead of `"Successfully updated!"`.

### Architectural

5. **`core.go` imports from the adapters layer** (`core.go:7`, `13`) -- `ApplicationSettings` references `persistence.PersistanceSettings`, violating the hexagonal dependency rule. The core package should not depend on adapter types. Configuration types should be defined in core or injected as plain values. *(Unresolved from #7.)*

6. **Services receive the entire `*ports.Repository` struct** (`tenants_service.go:13`; `data_sources_service.go:13`) -- `TenantsService` gets access to `Database`, `Tenants`, AND `DataSources` repositories. Violates Interface Segregation. Each service should receive only the repository interface(s) it needs (e.g., `TenantsService` should take `TenantsRepository`, not `*Repository`). *(Unresolved from #8.)*

7. **REST handlers hold a reference to the full `Application`** (`handlers.go:18-20`) -- The `handlers` struct takes `*core.Application` rather than just `*ports.Services`. This gives the HTTP layer access to the logger, repository, and any future internal state, violating the principle of least privilege. *(Unresolved from #9.)*

8. **`DataSource` can be created without verifying its parent `Tenant` exists** (`data_sources_service.go:31-39`) -- `Create` and `Update` call `Upsert` directly without checking that the `TenantID` corresponds to an existing tenant. Allows orphaned data sources.

9. **`Tenant.DataSources` field is declared but never populated** (`tenant.go:15`) -- No code path loads or sets it, giving a false impression that `Tenant` is an aggregate root containing its data sources.

10. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:23-25`) -- Returns every tenant in a single unbounded response. *(Unresolved from #10.)*

11. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 11 handlers) -- `net/http` already runs each handler in its own goroutine. The extra goroutine + buffered channel + `select` pattern adds allocation overhead and complexity with no concurrency benefit. *(Unresolved from #25.)*

### Missing Functionality

12. **Zero test files** in the entire tenants service. *(Unresolved from #12.)*

13. **Domain validation unimplemented** (`tenant.go:18-23`; `data_source.go:24-30`) -- Constructors return `(*Entity, error)` but never validate. *(Unresolved from #13.)*

14. **No domain events** for cross-service communication. *(Unresolved from #14.)*

15. **Incomplete error-to-HTTP mapping** (`http.go:48-72`) -- `ErrInvalidID`, `ErrUnauthorized`, and validation errors from `go-playground/validator` all fall through to the default 500 case. Validation errors should map to 400; `ErrUnauthorized` should map to 401; `ErrInvalidID` should map to 400 or 422. *(Unresolved from #15.)*

16. **No real authentication**. *(Unresolved from #16.)*

17. **`Lookup` service method is a stub** (`data_sources_service.go:57-66`) -- Always returns `nil, nil` after verifying the data source exists.

18. **`DataSourceAttributes` concrete types are empty shells** (`data_source_attributes.go:5-9`) -- `StaticDataSourceAttributes`, `ScheduleDataSourceAttributes`, and `QueryDataSourceAttributes` have zero fields.

19. **`DataSourceLookup` value object has no constructor or validation** (`data_source_lookup.go`) -- Bare struct, no invariant enforcement.

### Code Quality

20. **`DataSourceAttributes interface{}` is an empty interface** (`data_source_attributes.go:3`) -- Equivalent to `any`. Should be a sealed interface with a marker method (e.g., `isDataSourceAttributes()`) to restrict implementations to known types. *(Unresolved from #19.)*

21. **Request DTO `Attributes` field uses `any`** (`dto/request.go:14`) -- `UpsertDataSourceRequest.Attributes` is typed as `any`, so JSON deserialization produces a `map[string]interface{}` at runtime. There is no mapping logic to convert inbound JSON into the correct `DataSourceAttributes` implementation based on the `Type` field.

22. **`DataSourceType` is never validated against known constants** -- `DataSourceType` is a `string` typedef with three constants (`static`, `scheduled`, `query`), but neither the command constructors nor the domain constructors validate that a given value matches one of these. Any arbitrary string is accepted.

23. **`CreateDataSourceCommand` / `UpdateDataSourceCommand` have no `validate` struct tags** (`commands.go:53-57`, `77-82`) -- The `ValidateStruct` call in their constructors is a no-op since there are no tags to validate against. Compare with `baseTenantCommand` which has `validate:"required,max=75"` tags.

24. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:57`; `data_sources_repository.go:61`) -- Couples persistence to wall-clock time, making deterministic tests impossible. Should be injected via a clock interface or function. *(Unresolved from #26.)*

25. **`"Persistance"` typo throughout** (`persistence.go:11`, `14`, `19`, `21`, `23`; `core.go:13`) -- `PersistanceDriver`, `PersistanceSettings`, `PersistanceOptions` should all be `Persistence*`. *(Unresolved from #24.)*

26. **`w.Write` error ignored** (`http.go:42`) -- `SendJsonResponse` calls `w.Write(out)` but discards the returned error. *(Unresolved from #28.)*

27. **In-memory transactions are no-ops** (`inmemory_database.go`) -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing. No atomicity guarantee for multi-step operations. *(Unresolved from #29.)*

28. **`ValidateStruct` has a redundant pattern** (`validate.go:9-14`) -- The function body `if err := v.Struct(s); err != nil { return err }; return nil` is equivalent to `return v.Struct(s)`. Idiomatic Go prefers the simpler form.

29. **Inconsistent response envelope** -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. Clients must handle two different response shapes.

30. **`getTenants` returns tenants with `DataSources` always `nil`** (`handlers.go:46-49`) -- The inconsistency between the domain model (which declares the field) and the actual data flow is confusing.

---

## Summary

The tenants service has a clean hexagonal structure with well-separated packages: domain entities in `core/domain/`, port interfaces in `core/ports/`, service implementations in `core/services/`, and adapters in `adapters/rest/` and `adapters/persistence/`. Several issues from the 4/13 review have been resolved -- the `ErrExists` HTTP status mapping is now correct, the DTO layer has been properly refactored into a sub-package consistent with the forms service, the `compose.yaml` references the correct module, and the error message typo is fixed.

However, the most impactful issues from the previous review remain open: zero test coverage, unimplemented domain validation, missing error-to-HTTP mapping for validation and auth errors, the core-imports-adapters dependency violation, and overly broad service dependencies. New issues identified in this review include handlers silently swallowing JSON parse errors (clients get empty 200 responses), services discarding domain constructor errors, `Remove` operations that never return 404, and the `DataSource` entity being orphan-prone since tenant existence is never verified on create/update.

The highest-impact improvements would be:
1. Fix the swallowed JSON parse errors in all four mutating handlers (immediate bug)
2. Add tenant existence verification in `DataSourcesService.Create`/`Update` (data integrity)
3. Implement domain validation in `NewTenant`/`NewDataSource` constructors (the signatures already expect it)
4. Add `validate` struct tags to the data source commands (currently no-ops)
5. Map validation errors to HTTP 400 in `SendErrorResponse` (user-facing correctness)
6. Add test coverage starting with the service layer (long-term reliability)
