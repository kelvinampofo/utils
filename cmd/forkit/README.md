# forkit

`forkit` is a small CLI for embedded implementation prototypes.

Say you are changing a component and want to try two different prop APIs. Or say
you want an agent to take one pass while you, or another agent, try a different
one. Fork the file, edit the drafts, compare them in git or your editor, then
promote the one you want.

Start a fork folder:

```plain
forkit init src/components/Button.tsx
```

Create a couple of drafts:

```plain
forkit copy src/components/Button.tsx option-a
forkit copy src/components/Button.tsx option-b
```

Edit the copied files:

```plain
src/components/Button.forkit/option-a.tsx
src/components/Button.forkit/option-b.tsx
```

When one wins, copy it back over the real file:

```plain
forkit promote src/components/Button.tsx option-a --clean
```

The fork folder sits beside the source file. It contains `base.tsx`, your
drafts, and `manifest.json`, which lists the files in the set.

```plain
src/components/Button.tsx
src/components/Button.forkit/base.tsx
src/components/Button.forkit/option-a.tsx
src/components/Button.forkit/option-b.tsx
src/components/Button.forkit/manifest.json
```

Treat `*.forkit/` as scratch work. I usually ignore it:

```plain
*.forkit/
```

Use `git diff --no-index` or your editor to compare drafts. Use `forkit drop`
to delete one draft and `forkit clean` to delete the whole fork folder.

## Agents

For agents, the clean rule is: one agent owns one draft file. Tell the agent to
edit only `Button.forkit/agent-pass.tsx`, leave `Button.tsx` alone, and report
the tradeoffs when done. The human picks the winner and runs `forkit promote`.

There is a small agent skill next to the util:

```plain
cmd/forkit/SKILL.md
```

Copy `cmd/forkit/SKILL.md` into whatever
skill or instruction system your agent uses.

Then prompt the agent with the workflow and the target file:

```plain
# Prompt
Use the forkit agent workflow to make two drafts for src/components/Button.tsx:
simple-api and visual-polish. Don't promote; summarize tradeoffs.
```
