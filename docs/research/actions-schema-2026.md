# GitHub Actions schema audit: 2026-04-18 through 2026-10-18

Baseline: `011a6d15e749bb3f2d771eed9c7aa0e7e3e10ee7` (also local `main`).
Research branch: `investigate/actions-schema`.
Scope: workflow syntax, permission scopes, context properties/availability, action and reusable-workflow references, events, cache access. Concurrency and runner labels were deliberately excluded.

Read `CONTRIBUTING.md` before making changes. Research used `curl` against the official Actions changelog feeds/articles and `github/docs` raw/rendered pages, and `gh api` for official docs history and upstream open issues/PRs. Upstream proposals were leads, not authority for semantics. Documentation-addition dates below are distinguished from release-announcement dates; neither is assumed to be the first runtime availability date.

## Missing at baseline; implemented independently

Every implementation branch below starts directly at `main` / `011a6d1`, not on another fix. Worktrees are siblings of this worktree under `/Users/nigelstuke/Documents/repos/zew1me/action-lint-worktrees/`. Nothing was pushed; no PR was opened.

### Permission scopes

Branch `feat/token-permissions-2026`, worktree `token-permissions`.
Commit **`513317ece6e90bf527540684fdff0658124ad61f`**.

The permission-scope map rejected these documented scopes:

| Scope | Accepted levels | Official documentation addition |
| --- | --- | --- |
| `vulnerability-alerts` | `read`, `none` (not `write`) | [2026-04-23, github/docs commit](https://github.com/github/docs/commit/9e93d7ca5d7854d83b836785da362701ef1ea01e) |
| `code-quality` | `read`, `write`, `none` | [2026-05-11, github/docs commit](https://github.com/github/docs/commit/21f355ec39a11d9078ec3ee81c5023f1d328ebe7) |
| `copilot-requests` | `write`, `none` (not `read`) | [2026-10-13, github/docs commit](https://github.com/github/docs/commit/d1a0081f2407dbe4a08ba40b9cf3ae5e83659a9f) |

The [2026-09-03 Actions changelog](https://github.blog/changelog/2026-09-03-github-actions-early-september-2026-updates) also announces `vulnerability-alerts`. Canonical [permissions reference](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#permissions); [source enum](https://github.com/github/docs/blob/main/data/reusables/actions/github-token-available-permissions.md); [scope descriptions](https://github.com/github/docs/blob/main/data/reusables/actions/github-token-scope-descriptions.md).

Implementation changes only the scope table and its diagnostic snapshots. `TestPermissionsNewScopes` covers all three scopes at workflow, ordinary-job, and reusable-call level, including every valid level, prohibited levels, uppercase values and unknown values (45 subtests). Existing unknown-scope and `read-all`/`write-all` handling remains intact.

Upstream leads: [#666](https://github.com/rhysd/actionlint/pull/666), [#674](https://github.com/rhysd/actionlint/pull/674), [#689](https://github.com/rhysd/actionlint/pull/689), [#713](https://github.com/rhysd/actionlint/issues/713), [#720](https://github.com/rhysd/actionlint/issues/720).

### Context properties

Branch `feat/workflow-contexts-2026`, worktree `workflow-contexts`.
Commit **`71b6b5b2f116e2f237a95629771da90b0397d550`**.

The strict builtin context objects lacked six documented strings:

* `job.workflow_ref`, `job.workflow_sha`, `job.workflow_repository`, `job.workflow_file_path`: [official docs addition 2026-04-21](https://github.com/github/docs/commit/94ac73455f2e39623d1c976bbdc655acb2963ce7), also announced [2026-09-03](https://github.blog/changelog/2026-09-03-github-actions-early-september-2026-updates). These identify the file defining the current job, including a reusable workflow's own repository/commit, not necessarily the caller's. The docs explicitly distinguish them from `github.workflow_ref` / `github.workflow_sha`.
* `github.artifacts`, `github.artifacts_list`: [official docs addition 2026-08-12](https://github.com/github/docs/commit/f219414e4552acd7ee8aafe4e9da3516aab4a890). Both are **file paths**, not arrays or metadata objects. `artifacts_list` points to a file containing JSON.

All six are present in the current rendered [contexts reference](https://docs.github.com/en/actions/reference/workflows-and-actions/contexts). The corresponding [artifact-file version gate](https://github.com/github/docs/blob/main/data/features/actions-artifacts-file.yml) includes github.com. No context-availability rules were guessed or widened.

Tests add eight expression/type cases (dot and index access), a reusable workflow using checkout identity and artifact paths, and negative fixtures proving typo detection and string-vs-object/array errors still work. Existing strict context diagnostics are regenerated.

Upstream identity leads: [#647](https://github.com/rhysd/actionlint/issues/647), [#696](https://github.com/rhysd/actionlint/pull/696), [#707](https://github.com/rhysd/actionlint/pull/707), [#757](https://github.com/rhysd/actionlint/issues/757).

### Issue-field activity types

Branch `feat/workflow-events-2026`, worktree `workflow-events`.
Commit **`665972be35f67c9b8e39c2d05c7e6b8696b076ca`**.

`issues.field_added` and `issues.field_removed` were missing from `AllWebhookTypes`. Official [2026-06-04 docs change](https://github.com/github/docs/commit/bbeaf6d1102cc28151da211a0009d784904ad683) and current [Actions issues-event reference](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#issues) explicitly support both. `field_added` includes updates to an existing field value; no invented `field_changed` event was added.

Ran the existing generator:

```sh
go run ./scripts/generate-webhook-events ./all_webhooks.go
```

The **only generated change** was these two issue activity types. Tests accept both, reject misspellings, and reject these types on `pull_request`. Existing diagnostic fixtures and checks documentation are updated.

### `cache-mode`

Branch `feat/cache-mode`, worktree `cache-mode`.
Commit **`6292411fed169d7e01770e1d9f9a91188ae23247`**.

Official announcement: [2026-09-10](https://github.blog/changelog/2026-09-10-control-github-actions-cache-access-with-cache-mode). Documentation landed [2026-09-10](https://github.com/github/docs/commit/be60fe442256), covering [workflow syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#cache-mode), [job syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#jobsjob_idcache-mode), and [supported reusable-call keywords](https://docs.github.com/en/actions/reference/workflows-and-actions/reusing-workflow-configurations#supported-keywords-for-jobs-that-call-a-reusable-workflow).

Baseline rejected the key everywhere. The fix adds documented AST fields and parses the four literal modes at workflow, ordinary-job and reusable-call level:

| Value | Restore | Save |
| --- | --- | --- |
| `read` | yes | no |
| `write` | yes | yes |
| `write-only` | no | yes |
| `none` | no | no |

Job-level overrides and omitted values are preserved in the AST. Tests cover every value in all three locations, quoting, positions, anchors/aliases, omission/override behavior, invalid scalar values and invalid YAML shapes (mapping/sequence).

Deliberately **not** implemented: trigger-dependent effective defaults, caller/callee cache-capability mismatch analysis, or guessed expression-context availability. In particular, `read` and `write-only` are non-overlapping capabilities, not points on a simple ordered permission ladder. The [official reusable-workflow guidance](https://docs.github.com/en/actions/how-tos/reuse-automations/reuse-workflows#controlling-cache-access-in-reusable-workflows) documents these subtleties; this fix is parser/schema support, not a new cross-workflow policy checker.

Upstream lead: [#745](https://github.com/rhysd/actionlint/pull/745).

### Self-repository `$/` references

Branch `feat/self-repository-uses`, worktree `self-repository-uses`.
Commit **`a1dc5ac3fc96456ade2a8381c81e5c1d7bad8ed1`**.

Official announcement: [2026-07-30](https://github.blog/changelog/2026-07-30-reference-same-repository-actions-with-self-repository-syntax). Official [docs addition 2026-07-29](https://github.com/github/docs/commit/6a8580f09ba9), [workflow syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#jobsjob_idstepsuses), and [composite-action syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/metadata-syntax#runsstepsuses).

Baseline interpreted `$/...` as a malformed remote reference. The feature resolves against the repository of the **defining file**, at its running commit, without requiring checkout. It must not carry an `@ref` suffix. The announcement also covers reusable workflow calls and nested composition; it requires runner 2.336.0 or later on github.com.

The fix recognizes action and reusable-workflow `$/` references and reads metadata from the project being linted. Shared canonical cache keys are essential: reusable metadata may already have been written from an AST under the `./` spelling. Tests verify both lookup orders, AST-written metadata, cached missing-workflow failures, malformed refs/empty paths, required/unknown action inputs, reusable-workflow inputs/secrets, and strict action-output types. A multi-file fixture uses a self-repository action inside a reusable workflow as well as in its caller. The existing conservative behavior for absent local action directories is unchanged.

This does not add a new recursive composite-action step linter or download metadata for remote repositories. Those are separate capabilities, not necessary to recognize and validate references in the workflow/project paths actionlint already checks.

Upstream leads: [#711](https://github.com/rhysd/actionlint/issues/711), [#732](https://github.com/rhysd/actionlint/pull/732), [#760](https://github.com/rhysd/actionlint/issues/760). Unlike the described limitation in #732, this implementation rejects an action reference with a forbidden `@ref` suffix, as the official syntax explicitly requires.

## Already supported; no implementation needed

* **Merge-group branch filters**, documented [2026-09-21](https://github.com/github/docs/commit/63859c481b2195607ff94c4dda765fb131a34759): baseline `RuleEvents.checkWebhookEvent` already includes `merge_group` in the branch/branch-ignore allowlist and exclusion checks. Generic parser/glob handling already applies. Do not confuse this documentation update with an unsupported schema feature.
* **Matrix variable names are case-insensitive**, clarification [2026-09-23](https://github.com/github/docs/commit/03d2e24b34bd): baseline raw object/matrix keys are normalized to lower case and expression property dereferences already use the normalized keys. No new matrix schema field is required. This is not a claim about every historical bracket-index edge case.
* **`on.workflow_run.workflows`**, dedicated section added [2026-10-01](https://github.com/github/docs/commit/c67a9362ec63): baseline already parses `Workflows`, requires it for `workflow_run`, and rejects it for unrelated events.
* **Context availability**: reran `scripts/generate-availability` against the current official contexts source into a temporary file; it is byte-identical to baseline `availability.go`. No availability-table change is needed.
* Existing `github.workflow_ref`, `github.workflow_sha`, `job.check_run_id`, `artifact-metadata` permission, schedule timezones, and YAML anchors are already modeled. Their presence does not imply support for the new properties/scopes above.

## Unsupported/deferred or insufficiently substantiated

* **Parallel steps**, [2026-06-25 announcement](https://github.blog/changelog/2026-06-25-actions-steps-can-now-be-run-in-parallel), [docs addition 2026-06-24](https://github.com/github/docs/commit/f6aa36b7c56d), [limitations documented 2026-06-30](https://github.com/github/docs/commit/c5d56d694e3e): baseline lacks `background`, `wait`, `wait-all`, `cancel`, `parallel`. Deferred as a large AST/parser/visitor/expression/reference-validation change. Official [workflow syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax) is the semantic authority; the context-availability table still has no entries for these keys. Upstream [#694](https://github.com/rhysd/actionlint/pull/694), [#695](https://github.com/rhysd/actionlint/pull/695) were reviewed as leads. No partial parser-only acceptance was added that would silently skip linting grouped steps.
* **`merge_group.types: [destroyed]`**: upstream [#726](https://github.com/rhysd/actionlint/issues/726), [#740](https://github.com/rhysd/actionlint/pull/740) request this. The [webhook reference](https://docs.github.com/en/webhooks/webhook-events-and-payloads#merge_group) describes `destroyed`, but the [Actions trigger reference](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#merge_group) still lists only `checks_requested`; regenerating the Actions table does not add it. A webhook payload type alone is insufficient proof of a workflow trigger. No speculative change made.
* **`github.job_workflow_sha`**: upstream [#671](https://github.com/rhysd/actionlint/pull/671) itself notes the distinction between an OIDC token claim and a documented expression-context property. Current contexts reference still does not list it. Added the documented `job.workflow_sha` instead; no OIDC-claim-to-context inference.
* **`jobs.<job_id>.result` in reusable outputs**: documented by the [jobs context reference](https://docs.github.com/en/actions/reference/workflows-and-actions/contexts#jobs-context), and baseline currently lacks it. Upstream [#724](https://github.com/rhysd/actionlint/pull/724) is a real older-gap lead, but the property predates this audit window (the issue it closes is #343), so it was not presented as a new April–October schema change or implemented here.

## Out of scope / no workflow schema change

* Concurrency queue changes and every runner-label/image migration: explicitly assigned elsewhere; untouched.
* [2026-06-26 read-only cache defaults](https://github.blog/changelog/2026-06-26-read-only-actions-cache-for-untrusted-triggers) and [2026-06-18 safer pull_request_target checkout defaults](https://github.blog/changelog/2026-06-18-safer-pull_request_target-defaults-for-github-actions-checkout): runtime/security-default changes, not new workflow keys. Cache-mode syntax is separately implemented above.
* Workflow execution protections ([2026-06-18 preview](https://github.blog/changelog/2026-06-18-control-who-and-what-triggers-github-actions-workflows), [2026-09-17 GA](https://github.blog/changelog/2026-09-17-workflow-execution-protections-in-github-actions-generally-available)), bot-run approvals, malicious-workflow approvals: administrator policies/runtime admission, not workflow YAML additions.
* [Workflow dependency lockfiles, 2026-10-14 preview](https://github.blog/changelog/2026-10-14-github-actions-workflow-lockfiles-in-public-preview): [official concept reference](https://docs.github.com/en/actions/concepts/security/workflow-dependency-lockfiles) specifies a separate `.github/workflows/actions.lock` file and CLI-managed dependency resolution/policy. Not a new workflow key, and not the same feature as agentic workflow `.lock.yml` files. No lockfile validation added.
* [Agentic Workflows public preview, 2026-06-11](https://github.blog/changelog/2026-06-11-github-agentic-workflows-is-now-in-public-preview): Markdown/frontmatter compilation is a separate authoring format; not a reason to accept arbitrary new keys in ordinary workflow YAML.
* OIDC immutable subject claims (2026-04-23), GitHub App token-format changes (2026-04-24 / 2026-10-02), runner enforcement, Node 20 removal, artifact/run/status retention and API pagination changes, action-version metadata refreshes: no corresponding workflow syntax/context changes substantiated by this audit. These were not mixed into schema PRs.

## Validation and known baseline failures

On **each of the five implementation branches**, ran and passed:

```sh
go test . -count=1
go test -race .
go test ./... -skip '^TestDetectErrorBadRequest$' -count=1
go build ./...
go vet ./...
staticcheck ./...
go run ./scripts/check-checks -quiet docs/checks.md
git diff --check
```

(The core uncached test was also exercised by the uncached all-package run on branches where the earlier core run used the cache.) `pyflakes` was installed in an external temporary virtualenv for documentation generation/testing; `shellcheck` was already available. No repository dependency changes were made. The docs maintenance script was run with `-fix docs/checks.md` on every branch before its final quiet validation.

`go test ./...` was also attempted without exclusions on all five branches. The sole failure was **`scripts/generate-popular-actions/TestDetectErrorBadRequest`**, whose live malformed-URL request now succeeds instead of producing the expected exit status. Reproduced the same failure on untouched `011a6d1`; it is not caused by these schema changes. Related upstream lead: [#758](https://github.com/rhysd/actionlint/issues/758). The `-skip` run above excludes precisely that one test, not the generator package.

`govulncheck ./...` on the unchanged baseline reports **14 reachable vulnerabilities** in `golang.org/x/net@v0.52.0` / the installed Go 1.26.6 standard library (plus non-reachable dependency advisories). Thus the repository-wide vulnerability gate is not clean independently of this work. Dependency/toolchain updates were deliberately not bundled with schema fixes. The implementation branches retain baseline `go.mod` / `go.sum`.

Built a baseline binary and confirmed that the new positive context/event/cache/self-reference fixtures all exit 1 before the corresponding fixes. They pass in their respective branch suites. None of the implementation branches depends on another fix branch.

Raw downloaded research, command logs, generated temporary files and external tools were kept outside Git under `/tmp/actions-schema-research/`; the dated official URLs and immutable docs commits above are the durable evidence.
