---
name: forkit
description: Use when an agent is asked to work with forkit embedded implementation prototypes, run an implementation pass in a .forkit folder, coordinate multiple humans or agents exploring variants, or promote/discard forkit drafts while keeping the canonical source file protected until selection.
---

# Forkit Agent Workflow

## Purpose

Use `forkit` to create embedded implementation prototypes beside a source file. The canonical source file stays unchanged until a human explicitly asks to promote a draft.

## Default Workflow

1. Confirm the target source file.
2. Run `forkit init <file>` if `<stem>.forkit/manifest.json` does not exist.
3. Create or claim one variant with `forkit copy <file> <variant>`.
4. Edit only the assigned draft file, such as `Button.forkit/agent-pass.tsx`.
5. Use normal tools to compare drafts: IDE split view or `git diff --no-index`.
6. Do not run `forkit promote` unless the user explicitly selects a winning variant.
7. After promotion, run the project's normal build/test/check commands unless the user says not to.
8. Use `forkit drop <file> <variant>` only when asked to discard one draft.
9. Use `forkit clean <file>` only when asked to remove all drafts.

## Agent Rules

- Announce which variant file you own before editing.
- Treat `*.forkit/` as scratch space.
- Keep edits scoped to your assigned draft file unless the user explicitly asks for broader changes.
- Do not edit the canonical source file while exploring.
- Do not promote your own draft by default; summarize the tradeoffs and wait for selection.
- If asked to promote, validate the promoted canonical source file with the repo's usual commands.
- If another agent owns a variant, do not modify or remove it.
- If supporting files must change to make a draft meaningful, ask first and explain why the draft cannot stay file-local.

## Useful Commands

```plain
forkit init src/components/Button.tsx
forkit copy src/components/Button.tsx agent-pass
forkit copy src/components/Button.tsx visual-pass --from base
forkit list src/components/Button.tsx
git diff --no-index src/components/Button.forkit/agent-pass.tsx src/components/Button.forkit/visual-pass.tsx
forkit promote src/components/Button.tsx agent-pass --clean
```

## Good Final Summary

When finishing a draft, report:

- the variant file edited
- the implementation direction
- important tradeoffs
- any commands run
- validation after promotion, if promotion happened
- whether the canonical source file was left untouched
