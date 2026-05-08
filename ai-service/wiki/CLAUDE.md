# CBPI Wiki — schema and authoring contract

This wiki is the long-term knowledge base of the CBPI assistant. It follows
the **Karpathy LLM Wiki** pattern: plain markdown on disk, index-first
retrieval, and small enough that the entire index fits in a system prompt.
No embeddings or vector store until the wiki passes ~300 pages.

The on-disk tree is the source of truth. Every page is git-tracked. Reads go
through `app/services/wiki_reader.py`; writes go through
`app/services/wiki_writer.py` and never bypass `app/services/ingest_pipeline.py`.

## Categories (top-level folders)

| Folder       | What goes there                                                             |
|--------------|-----------------------------------------------------------------------------|
| `entities/`  | Concrete things: a currency, an exchange, a broker, a fund, a regulator.    |
| `concepts/`  | Reusable ideas: hedging, exposição cambial, PFIC, tax-deferred accounts.    |
| `sources/`   | External provenance: a paper, an Anthropic-fetched article, a B3 release.   |
| `syntheses/` | Cross-cutting analyses that combine several entities/concepts.              |
| `playbooks/` | Decision flows the assistant can walk a user through.                       |

A page must live under exactly one category. Cross-references use wikilinks.

## Page format

Every page **must** start with a YAML front-matter block:

```markdown
---
type: entity | concept | source | synthesis | playbook
status: draft | stable
last_reviewed: YYYY-MM-DD
sources:
  - <wikilink or URL>
---

# Título da página em PT-BR

(corpo da página em português)
```

Required keys: `type`, `status`, `last_reviewed`. `sources` is required for
`type: source` pages and recommended elsewhere. `status: stable` pages cannot
be overwritten by an ingest unless the edit sets `force: true` — this protects
hand-curated knowledge from automated routines.

### Section template (recommended)

```markdown
## TL;DR
Uma frase.

## Definição
...

## Por que importa para o investidor brasileiro
...

## Relacionado
- [[concepts/exposicao-cambial]]
- [[entities/USD]]
```

### Wikilinks

Use double-bracket syntax with the **path without the `.md` suffix**:
`[[entities/USD]]`, `[[concepts/PFIC]]`. The linter rejects links that don't
resolve to an existing page.

### Language

- Page **content** is Portuguese (PT-BR).
- **Structural** fields (front-matter keys, folder names, file slugs, tags) are
  English/ASCII so tooling stays predictable.
- File slugs use `kebab-case` ASCII (`exposicao-cambial.md`, not
  `exposição-cambial.md`).

## The weekly Anthropic ingest routine

The wiki is fed by a routine that runs once per week. It is implemented as an
Anthropic Managed Agent (configured outside this repo) that:

1. Searches the public web for fresh material on cross-border investing
   relevant to a Brazilian investor (BRL/USD exposure, PFIC, IRRF, ETF
   listings on B3, etc.).
2. Drafts new `sources/` pages and updates to `concepts/` / `syntheses/`.
3. Calls `POST /v1/wiki/ingest` on this service with a Bearer token
   (`WIKI_INGEST_TOKEN`) and the payload defined in
   `app/schemas/ingest.py::IngestRequest`.

The receiving endpoint is the only contract the routine has to honour. It:

- Validates each edit (path, front-matter, no traversal, no clobber of
  `status: stable` without `force=true`).
- Applies all edits **atomically** through a tempdir-then-rename swap with
  per-file `.bak` backups, so partial application is impossible.
- Appends a one-line entry to `wiki/log.md` for auditability.

The routine **must not** mark pages as `status: stable`; promotion to stable
is a human-only operation.

## Operational rules for the assistant

When answering a user question:

1. Read this schema and `index.md` (both are sent in the cached system block).
2. Pull in only the page bodies whose slugs appear relevant to the query.
3. If the wiki doesn't cover the question, say so — never fabricate.
4. Cite every page used as `[[category/slug]]`.
5. End every answer with the educational disclaimer.
