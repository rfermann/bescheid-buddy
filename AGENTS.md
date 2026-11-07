# Agent guidance

## API deployment and logs

`apps/api` is deployed on a Coolify instance connected through the Coolify MCP.

After changing or debugging the API, inspect the logs for the `dev-bb-api`
service. Its container runs Air, so code changes trigger an automatic rebuild.
Use the log output to confirm that recompilation succeeded and to diagnose any
compiler errors.

## Agent skills

### Issue tracker

Issues and specs live in Jira project `BB` on `machkinfach.atlassian.net`. See `docs/agents/issue-tracker.md`.

### Triage labels

Use the default triage label vocabulary. See `docs/agents/triage-labels.md`.

### Domain docs

Use a single-context root glossary and ADR directory. See `docs/agents/domain.md`.
