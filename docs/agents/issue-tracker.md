# Issue tracker: Jira

Issues and specs live in Jira project `BB` at `https://machkinfach.atlassian.net/`.

Use the project-local `atlassian_machkinfach` MCP server to search, read, create, and update issues. Keep implementation plans, acceptance criteria, dependencies, and discussion on the Jira issue.

Use Jira issue keys such as `BB-123` when referring to work in code, documentation, or agent conversations.

## Work types

- **Epic**: a discovery effort or a delivery effort.
- **Story**: an independently deliverable user-facing slice with acceptance criteria.
- **Task**: a decision, investigation, prototype, or technical or operational change.
- **Bug**: a defect to correct.
- **Subtask**: an optional, separately tracked part of a Story, Task, or Bug.

Project `BB` supports all five types. Labels are available on every type. Stories, Tasks, and Bugs can have an Epic parent; Subtasks require a parent. Epics have no Parent field.

## Wayfinding operations

- Create the Wayfinder map as a **Discovery Epic**, labelled `wayfinder:map`. Reuse it when resuming the same map.
- Create decision tickets as child **Tasks** of the Discovery Epic. Classify them with `wayfinder:research`, `wayfinder:prototype`, `wayfinder:grilling`, or `wayfinder:task`.
- Keep the map as an index: each Task holds its question and records its answer in a resolution comment. Link the resolved Task by title from the map's Decisions so far.
- Claim a Task by assigning it to the developer driving the map before starting work.
- Express dependencies with Jira's **Blocks** link (`blocks` / `is blocked by`).
- Find candidates with `project = BB AND parent = "<discovery-epic-key>" AND statusCategory != Done AND assignee IS EMPTY ORDER BY created ASC`. Inspect their `is blocked by` links; the frontier contains only candidates whose blockers are all in the Done status category. Take the first frontier Task unless the user names another.
- Complete the Discovery Epic when its destination is reached, all decisions are resolved, and no in-scope questions remain. Track delivery separately.

## Delivery handoff

- When implementation is approved, create or identify a separate **Delivery Epic** for the resulting spec and implementation tickets.
- Link the Discovery and Delivery Epics using **Relates** (`relates to`). Give both the same `effort:<slug>` label. Project `BB` has no `implements` link type.
- `/to-tickets` publishes independently deliverable Stories, Tasks, and Bugs directly under the Delivery Epic. Link relevant discovery decisions and use Blocks links for implementation dependencies.
- Use Subtasks only when a delivery ticket needs a smaller, separately tracked breakdown. Keep implementation tickets under the Delivery Epic rather than under discovery Tasks.

## Grilling outside Wayfinder

Update an existing relevant issue with the agreed conclusions. When a separate durable decision record is needed, use a **Task** labelled `decision`; record the question and final answer there. Create a Story when the work describes a deliverable user outcome. A grilling session alone does not require a new issue.
