# 30-Day Optimization Plan (Go SDK)

## Outcome Target

- Make Go SDK production-hardened with predictable behavior, measurable quality checks, and fast adoption path.
- Keep first API success under 10 minutes and first SSE integration under 30 minutes.

## P0 (Day 1-14): Reliability and Contract Hardening

| Workstream | Task | Files | Acceptance |
| --- | --- | --- | --- |
| API Contract Alignment | Confirm path/method coverage matches backend contract and docs | `client.go`, `openapi/ai-sandbox-v1.yaml`, `docs/INTEGRATION.md` | 11 API operations validated with no mismatch |
| Safe Retry Policy | Keep retries idempotent by default and control write retry behavior explicitly | `client.go`, `README.md` | No duplicate write operations in fault simulation |
| Error Observability | Extend error payload surface with trace/request id when provided | `client.go`, `README.md`, `docs/INTEGRATION.md` | Troubleshooting section includes trace workflow |
| QA Baseline | Add unit tests for request builder, retries, and SSE parser logic | `client.go`, `*_test.go` (new), `go.mod` | CI tests pass with critical-path coverage target |
| CI/CD Guardrails | Add workflow for `go test`, vet, lint, and module checks | `.github/workflows/ci.yml` (new), `go.mod` | PR checks enforced before merge |

## P1 (Day 15-30): Developer Experience and Growth

| Workstream | Task | Files | Acceptance |
| --- | --- | --- | --- |
| Example Expansion | Extend quickstart to full flow including SSE and KB interactions | `examples/quickstart/main.go`, `README.md` | Example runs with env vars only |
| Visual Docs Upgrade | Add troubleshooting matrix and production tuning guide | `README.md`, `docs/INTEGRATION.md` | Reduced onboarding support tickets |
| Release Quality | Add release checklist and compatibility notes per version | `CHANGELOG.md`, `CONTRIBUTING.md` | Every release has impact summary |
| Security Posture | Enable vulnerability and secret scanning workflows | `.github/workflows/ci.yml`, `SECURITY.md` | No unresolved high-severity issue at release gate |

## Language File Checklist

- `README.md`
- `docs/INTEGRATION.md`
- `docs/30D_OPTIMIZATION_PLAN.md`
- `client.go`
- `examples/quickstart/main.go`
- `openapi/ai-sandbox-v1.yaml`
- `go.mod`
- `CHANGELOG.md`
- `CONTRIBUTING.md`
- `SECURITY.md`

## Definition of Done (DoD)

- [ ] 11/11 API operations pass production integration validation.
- [ ] SSE stream parser handles chunking and stops correctly on `[DONE]`.
- [ ] Retry defaults prevent duplicate non-idempotent write operations.
- [ ] CI enforces build/lint/test gates on every PR.
- [ ] Quickstart runs from clean environment using required env vars only.
