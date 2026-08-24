<div align="center">

# #010 Go (Golang) over Java

##### A record that describes the architectural decision, its context, and its consequences.

<img src="../imgs/architecture-design-record-logo.png" style="width:175px;"/>

</div>

## Context

An implementation language had to be selected for the Forms Hub services before development began. At Wells Fargo, Java — typically paired with Spring and the associated enterprise platform tooling — is the default choice for backend services. Selecting any other language is therefore a deliberate deviation from the organisational standard and warrants an explicit, recorded decision that future maintainers can revisit.

The runtime and operational profile of Forms Hub shaped this decision. The system is composed of two independently deployable microservices (the Tenants Service and the Forms Service) plus a shared `pkg/` library, deployed to OpenShift Container Platform with multiple replicas per service (see Deployment View). Several architectural characteristics place particular weight on the language runtime:

- **Horizontal scaling** — multiple replicas of each service run simultaneously and are frequently rescheduled by Kubernetes, so per-process startup time and memory footprint directly affect density, cost, and recovery speed.
- **Background workers** — both services run ticker-driven workers: the `DistributedWorker[J Job]` with Redis leader election (ADR-004) and the `PeriodicWorker[J Job]` used by the outbox relay (ADR-009). These rely heavily on lightweight concurrency primitives.
- **Asynchronous submission processing** — the submission pipeline is intentionally decoupled from intake (ADR-003), leaning on concurrent, non-blocking execution rather than a heavyweight threading model.

These same factors are summarised in the language entry of `architecture.md` §4.2 Technology Decisions; this ADR records the full reasoning behind that entry.

## Decision

Go is the implementation language for both the Tenants Service and the Forms Service, as well as the shared `pkg/` library, in preference to Java. The decision rests on four factors, described below in order of significance to this project.

**1. Operational simplicity.** Go compiles to a single statically linked binary. This produces minimal Docker images with no Java Virtual Machine, application server, or class-loading layer to ship or operate, and removes the need for ongoing garbage-collector and heap tuning per environment. On OCP, where each service is independently deployed and scaled, this materially reduces the operational surface area and the per-service configuration burden.

**2. Concurrency model.** Go's goroutines and channels map cleanly onto the system's concurrency needs. The generic background worker abstraction (`DistributedWorker[J Job]` and `PeriodicWorker[J Job]`) uses a goroutine pool with graceful drain-based shutdown, and the asynchronous submission pipeline (ADR-003) processes work off the request path. Expressing this with lightweight goroutines is simpler and lower-overhead than the equivalent thread-pool and executor machinery in Java.

**3. Runtime footprint and startup.** Go processes start in milliseconds and hold a small resident memory footprint compared to a warmed JVM. This suits an environment where replicas are scaled out, rescheduled, and restarted routinely — cold-start latency is low and each replica is inexpensive, improving density and resilience. Fast cold start additionally makes aggressive downscaling — or scale-to-zero in a cloud-hosted deployment — viable for request-driven load such as overnight idle periods, a posture the JVM's warm-up profile would penalise. This benefit is contingent on the always-on background workers (ADR-004, ADR-009) being handled separately, and is not yet realised on the current internal OCP deployment (§7).

**4. Ecosystem fit.** The libraries selected for the design integrate cleanly with Go and reinforce the Ports and Adapters structure (ADR-001): the `chi` router provides a composable middleware chain for cross-cutting concerns (auth, tenant extraction, idempotency, correlation ID); `expr-lang/expr` supplies safe, sandboxed rule evaluation (ADR-007); and the official MongoDB Go driver supports the document-oriented persistence model and multi-document transactions (ADR-006). Time-ordered UUID v7 identifiers and the generic `StrategyRegistry` (ADR-005) similarly fit idiomatic Go.

**Alternatives considered.** Java with Spring was the primary alternative, given its status as the enterprise default and the depth of internal platform support around it. It was not selected because its runtime and operational profile (JVM warm-up, larger memory footprint, application-server and framework configuration) works against the small-runtime, many-replica, worker-heavy design, and because the chosen libraries and concurrency model are a more natural fit for Go. The trade-off is a deliberate departure from the organisational default, accepted for the reasons above and reconsidered in the Consequences below.

## Consequences

- **Smaller internal talent pool.** The choice deviates from the Wells Fargo Java/Spring enterprise default and draws on a smaller internal pool of Go engineers. This may affect staffing, onboarding time, code review depth, and long-term maintenance, and should be weighed whenever the team scales or hands over ownership.
- **Less mature enterprise library support.** Internal enterprise library and framework support for Go is less mature than for Java. Fewer Wells-Fargo-standard, pre-approved Go libraries exist, so a number of cross-cutting concerns are implemented in-house within `pkg/` — auth middleware, the generic background worker and elector, and the cache and database abstractions — rather than inherited from an established platform stack. This increases in-house ownership of infrastructure code and its associated review and maintenance cost.
- **Uniform, low-overhead tooling.** Build, packaging, and runtime behaviour are consistent across both services and the shared library. A single toolchain and static-binary deployment model keeps the independently deployable, small-runtime microservice strategy coherent and easy to reason about.
- **Potential cloud cost savings under idle load.** The low runtime footprint and fast cold start open the door to scale-to-zero or heavy overnight downscaling in a cloud-hosted deployment, where request-driven form traffic is likely to be idle outside business hours. Realising this depends on decoupling the always-on background workers (ADR-004, ADR-009) from the request-serving replicas, and does not apply to the current internal OCP deployment (§7); it is recorded here as a forward-looking opportunity rather than a present benefit.
- **Contained blast radius for the choice.** The hexagonal boundaries (ADR-001) keep the language and its libraries isolated behind port interfaces. Should any individual library — or, in the extreme, the language decision itself for a given service — need to change, the domain and application layers are insulated, limiting the cost of revisiting this decision.
- **Revisit triggers.** This decision should be reconsidered if internal Go talent or approved-library support becomes a sustained constraint, or if org-wide platform mandates make Java materially cheaper to operate than the in-house Go infrastructure.
