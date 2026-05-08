import os
import tempfile
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path

from app.core.wiki_paths import LOG_FILE, resolve_within
from app.schemas.ingest import IngestRequest, WikiEdit


@dataclass(frozen=True)
class WriteResult:
    path: str
    op: str
    bytes_written: int


class WikiWriter:
    """Atomic batch writer. Every edit is staged into a tempdir then renamed
    onto the live tree. If any rename fails, completed renames are reversed
    using pre-write backups so the on-disk wiki is never half-applied.
    """

    def __init__(self, *, wiki_root: Path) -> None:
        self._root = wiki_root

    def apply(self, request: IngestRequest) -> list[WriteResult]:
        staged: list[tuple[Path, Path, WikiEdit]] = []  # (tmp, target, edit)
        with tempfile.TemporaryDirectory(prefix="wiki-stage-") as stage_dir:
            stage_root = Path(stage_dir)
            for edit in request.edits:
                target = resolve_within(self._root, edit.path)
                tmp = stage_root / edit.path
                tmp.parent.mkdir(parents=True, exist_ok=True)
                tmp.write_text(edit.content, encoding="utf-8")
                staged.append((tmp, target, edit))

            results, completed, backups = [], [], {}
            try:
                for tmp, target, edit in staged:
                    target.parent.mkdir(parents=True, exist_ok=True)
                    if target.exists():
                        backup = target.with_suffix(target.suffix + ".bak")
                        os.replace(target, backup)
                        backups[target] = backup
                    os.replace(tmp, target)
                    completed.append(target)
                    results.append(
                        WriteResult(
                            path=edit.path,
                            op=edit.op,
                            bytes_written=len(edit.content.encode("utf-8")),
                        )
                    )
            except Exception:
                for done in reversed(completed):
                    if done in backups:
                        os.replace(backups[done], done)
                    else:
                        try:
                            done.unlink()
                        except FileNotFoundError:
                            pass
                raise
            else:
                for backup in backups.values():
                    backup.unlink(missing_ok=True)

        self._append_log(request)
        return results

    def _append_log(self, request: IngestRequest) -> None:
        log_path = self._root / LOG_FILE
        ts = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        paths = ", ".join(e.path for e in request.edits)
        line = f"- {ts} · source={request.source.name} · {len(request.edits)} edits ({paths}) — {request.log_entry}\n"
        with log_path.open("a", encoding="utf-8") as f:
            f.write(line)
