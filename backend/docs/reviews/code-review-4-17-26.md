# Tenants Service Code Review

## Issues Resolved Since 4/16 Review

1. ~~Handlers swallow JSON parse errors~~ (#1) -- All four mutating handlers (`createTenant`, `updateTenant`, `createDataSource`, `updateDataSource`) now call `h.sendErrorResponse(w, err)` and `return` when `ReadJsonPayload` fails (`handlers.go:84-87`, `122-125`, `242-245`, `287-290`).

2. ~~`updateTenant` success message typo `"Success updated!"`~~ (#4) -- Now correctly says `"Successfully updated!"` (`handlers.go:148`).

3. ~~Request DTO `Attributes` field uses `any` with no mapping logic~~ (#21) -- A strategy pattern in `dto/request.go` now maps inbound JSON to typed `DataSourceAttributes` implementations (`StaticDataSourceAttributes`, `ScheduledDataSourceAttributes`, `QueryDataSourceAttributes`) based on the `Type` field via `RequestToDataSourceAttributes`.

4. ~~`CreateDataSourceCommand` / `UpdateDataSourceCommand` have no `validate` struct tags~~ (#23) -- `baseDataSourceCommand` now has `validate:"required"` on `TenantID`, `validate:"oneof=static scheduled query"` on `Type`, and `validate:"required"` on `Attributes` (`commands.go:53-57`). `Attributes` has been moved from the individual command structs into the shared base. `UpdateDataSourceCommand.ID` now has `validate:"required"` (`commands.go:85`).

5. ~~`"Persistance"` typo throughout~~ (#25) -- All names now use the correct spelling: `PersistenceDriver`, `PersistenceSettings`, `PersistenceOptions` (`persistence.go:11`, `19`, `21`).

6. ~~`updateDataSource` sends `nil` to error handler~~ (caught during review) -- The `select` case was using the outer-scope `err` variable instead of `res.err`. Now correctly passes `res.err` (`handlers.go:315`).

7. ~~`createTenant` / `updateTenant` write two HTTP responses for validation errors~~ (caught during review) -- The redundant inner `if validate.IsValidationErr(err)` blocks that sent a duplicate 400 response before falling through to `sendErrorResponse` have been removed. Both handlers now delegate directly to `sendErrorResponse` (`handlers.go:91`, `129`).

8. ~~Inconsistent response mapper naming~~ (#29 partial) -- `DataSourceToResponseDto` renamed to `DataSourceToResponse` for consistency with `TenantToResponse` (`dto/response.go:36`).

---

## Remaining Issues

### Bugs

1. **JSON parse errors map to 500 instead of 400** (`handlers.go:84-87`, `122-125`, `242-245`, `287-290`) -- While handlers no longer swallow JSON parse errors (resolved #1), the `json.Decoder` error from `ReadJsonPayload` is not a `validator.ValidationErrors` and is not `ErrDataSourceAttrParse`, so `sendErrorResponse` (`handlers.go:379-391`) delegates to `common.SendErrorResponse` which falls through to the default 500 case. Malformed JSON from the client produces a 500 Internal Server Error instead of a 400 Bad Request.

2. **Services silently discard `domain.New*` constructor errors** (`tenants_service.go:32`, `43`; `data_sources_service.go:32`, `43`) -- Both services assign the error from `domain.NewTenant` / `domain.NewDataSource` to `_`. If domain constructors gain validation logic, those errors will be silently lost. *(Unresolved from #2.)*

3. **`Remove` operations succeed silently for non-existent entities** (`tenant_repository.go:78-84`; `data_sources_repository.go:84-91`) -- Go's `delete` on a map is a no-op for missing keys, so deleting a non-existent tenant or data source returns `nil` (HTTP 204) instead of `ErrNotFound` (HTTP 404). *(Unresolved from #3.)*

4. **`ErrDataSourceAttrParse` is defined but never returned** (`dto/request.go:9`) -- The sentinel error is defined and checked in `sendErrorResponse` (`handlers.go:381`), but `RequestToDataSourceAttributes` and `parseAttributes` never return or wrap it. The actual errors from `json.Marshal`, `strategy.Get`, and `json.Unmarshal` are returned unwrapped, so the `errors.Is(err, dto.ErrDataSourceAttrParse)` check is dead code. Attribute parsing errors currently fall through to the default 500 case in `common.SendErrorResponse`.

### Architectural

5. **`core.go` imports from the adapters layer** (`core.go:7`, `13`) -- `ApplicationSettings` references `persistence.PersistenceSettings`, violating the hexagonal dependency rule. The core package should not depend on adapter types. *(Unresolved from #5.)*

6. **Services receive the entire `*ports.Repository` struct** (`tenants_service.go:13`; `data_sources_service.go:13`) -- `TenantsService` gets access to `Database`, `Tenants`, AND `DataSources` repositories. Violates Interface Segregation. Each service should receive only the repository interface(s) it needs. *(Unresolved from #6.)*

7. **REST handlers hold a reference to the full `Application`** (`handlers.go:20-22`) -- The `handlers` struct takes `*core.Application` rather than just `*ports.Services`. This gives the HTTP layer access to the logger, repository, and any future internal state. *(Unresolved from #7.)*

8. **`DataSource` can be created without verifying its parent `Tenant` exists** (`data_sources_service.go:31-39`) -- `Create` and `Update` call `Upsert` directly without checking that the `TenantID` corresponds to an existing tenant. Allows orphaned data sources. *(Unresolved from #8.)*

9. **`Tenant.DataSources` field is declared but never populated** (`tenant.go:15`) -- No code path loads or sets it, giving a false impression that `Tenant` is an aggregate root containing its data sources. *(Unresolved from #9.)*

10. **`Find()` in tenants service has no pagination or filtering** (`tenants_service.go:23-25`) -- Returns every tenant in a single unbounded response. *(Unresolved from #10.)*

11. **Unnecessary goroutine/channel pattern in every handler** (`handlers.go`, all 11 handlers) -- `net/http` already runs each handler in its own goroutine. The extra goroutine + buffered channel + `select` pattern adds allocation overhead and complexity with no concurrency benefit. *(Unresolved from #11.)*

### Missing Functionality

12. **Zero test files** in the entire tenants service. *(Unresolved from #12.)*

13. **Domain validation unimplemented** (`tenant.go:18-23`; `data_source.go:24-30`) -- Constructors return `(*Entity, error)` but never validate. *(Unresolved from #13.)*

14. **No domain events** for cross-service communication. *(Unresolved from #14.)*

15. **Incomplete error-to-HTTP mapping** (`http.go:48-72`) -- `ErrInvalidID` and `ErrUnauthorized` are defined in `error.go` but not mapped in `SendErrorResponse`. They fall through to the default 500 case. `ErrInvalidID` should map to 400 or 422; `ErrUnauthorized` should map to 401. *(Unresolved from #15.)*

16. **No real authentication**. *(Unresolved from #16.)*

17. **`Lookup` service method is a stub** (`data_sources_service.go:57-66`) -- Always returns `nil, nil` after verifying the data source exists. *(Unresolved from #17.)*

18. **`DataSourceAttributes` concrete types are incomplete** (`data_source_attributes.go`) -- `ScheduledDataSourceAttributes` has zero fields. `StaticDataSourceAttributes` and `QueryDataSourceAttributes` have fields but lack `json` struct tags, so JSON marshaling uses Go's default capitalized field names (e.g., `"Data"` instead of `"data"`, `"Endpoint"` instead of `"endpoint"`). *(Unresolved from #18.)*

19. **`DataSourceLookup` value object has no constructor or validation** (`data_source_lookup.go`) -- Bare struct, no invariant enforcement. *(Unresolved from #19.)*

### Code Quality

20. **`DataSourceAttributes interface{}` is an empty interface** (`data_source_attributes.go:3`) -- Equivalent to `any`. Should be a sealed interface with a marker method (e.g., `isDataSourceAttributes()`) to restrict implementations to known types. *(Unresolved from #20.)*

21. **`DataSourceType` is never validated in domain constructors** -- `DataSourceType` is a `string` typedef with three constants (`static`, `scheduled`, `query`). The `validate:"oneof=..."` tag on `baseDataSourceCommand.Type` now covers command validation, but the domain constructor `NewDataSource` still accepts any arbitrary string. *(Partially resolved from #22.)*

22. **`time.Now()` called directly in the repository layer** (`tenant_repository.go:57`; `data_sources_repository.go:61`) -- Couples persistence to wall-clock time, making deterministic tests impossible. Should be injected via a clock interface or function. *(Unresolved from #24.)*

23. **`w.Write` error ignored** (`http.go:42`) -- `SendJsonResponse` calls `w.Write(out)` but discards the returned error. *(Unresolved from #26.)*

24. **`SendJsonResponse` accepts `headers` parameter but never applies them** (`http.go:33`) -- The `headers ...http.Header` variadic parameter is accepted in the function signature but the body never iterates or sets them. Callers expecting custom headers to be applied will be silently ignored.

25. **In-memory transactions are no-ops** (`inmemory_database.go`) -- `BeginTx`/`CommitTx`/`RollbackTx` do nothing. No atomicity guarantee for multi-step operations. *(Unresolved from #27.)*

26. **`ValidateStruct` has a redundant pattern** (`validate.go:19-25`) -- The function body `if err := v.Struct(s); err != nil { return err }; return nil` is equivalent to `return v.Struct(s)`. *(Unresolved from #28.)*

27. **Inconsistent response envelope** -- List endpoints (`getTenants`, `getDataSources`, `getDataSourceLookup`) return a bare JSON array, while create/update endpoints return an `ApiResponse[T]` wrapper with a `message` field. Clients must handle two different response shapes. *(Unresolved from #29.)*

28. **`getTenants` returns tenants with `DataSources` always `nil`** (`handlers.go:48-51`) -- The inconsistency between the domain model (which declares the field) and the actual data flow is confusing. *(Unresolved from #30.)*

---

## Summary

The tenants service is showing an accelerated pace of improvement. Eleven issues have been resolved since the 4/16 review -- more than double the five resolved in the prior cycle. Key fixes include handlers now properly responding to JSON parse failures instead of leaving clients with empty 200s, validation struct tags added to all data source commands with `Attributes` consolidated into the shared base command, the strategy pattern for attribute deserialization, and several bugs caught during review (wrong error variable in `updateDataSource`, double HTTP responses for validation errors) fixed before they reached production.

However, the fix for JSON parse error handling introduced a gap: while handlers now call `sendErrorResponse`, the `json.Decoder` error type is not recognized by any case in the error-handling chain, so malformed JSON from clients still produces a 500 Internal Server Error instead of the expected 400. Additionally, `ErrDataSourceAttrParse` is defined and checked in `sendErrorResponse` but is never actually returned by any code path, making that branch dead code. Both issues mean certain client errors surface as 500s.

The core architectural and missing functionality issues from prior reviews remain open: zero test coverage, unimplemented domain validation, the core-imports-adapters dependency violation, overly broad service dependencies, missing error-to-HTTP mapping for `ErrInvalidID`/`ErrUnauthorized`, and no tenant existence verification for data source operations.

The highest-impact improvements would be:
1. Map `json.Decoder` errors and attribute parsing errors to HTTP 400 in `sendErrorResponse` (immediate -- client errors produce 500s)
2. Either return `ErrDataSourceAttrParse` from `RequestToDataSourceAttributes` or remove the dead `errors.Is` check (correctness)
3. Implement domain validation in `NewTenant`/`NewDataSource` constructors (the signatures already expect it)
4. Add tenant existence verification in `DataSourcesService.Create`/`Update` (data integrity)
5. Make `Remove` return `ErrNotFound` for missing entities (API correctness)
6. Add test coverage starting with the service and handler layers (long-term reliability)
