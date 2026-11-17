#!/usr/bin/env python3
"""Build a self-contained copy of the 29-N preview.

Takes the base HTML (tabs, sections, cards) and embeds the run documents inside
it: the summaries and verification reports as rendered HTML, and the PDFs as
base64 so the result opens from any location, with or without the original
files next to it.

    python3 build_standalone.py [--no-pdf] [--out FILE]
"""
import argparse
import base64
import json
import pathlib
import sys

import markdown

HERE = pathlib.Path(__file__).resolve().parent.parent
BASE = HERE / "Programas-electorales-2026.html"
# The source PDFs are not versioned (immutable, 28 MB). They are looked up in
# this order; --pdf-dir overrides it for a different checkout.
PDF_DIRS = [
    HERE / "pdfs",
    HERE,
    pathlib.Path("/Users/jparrill/Documents/programas-electorales-opencode-2026-10-30/29N-2026"),
]

DOCS = {
    "fo": {
        "name": "Frente Obrero",
        "resumen": "Resumen-Frente-Obrero-2026.md",
        "verificacion": "Verificacion-Resumen-Frente-Obrero-2026.md",
        "pdf": "Frente-Obrero-Programa-2026-10.pdf",
        "pdf_resumen": "Resumen-Frente-Obrero-2026.pdf",
    },
    "vox": {
        "name": "VOX",
        "resumen": "Resumen-VOX-2026.md",
        "verificacion": "Verificacion-Resumen-VOX-2026.md",
        "pdf": "Programa-completo-A5-1.pdf",
        "pdf_resumen": "Resumen-VOX-2026.pdf",
    },
    "pp": {
        "name": "Partido Popular",
        "resumen": "Resumen-PP-2026.md",
        "verificacion": "Verificacion-Resumen-PP-2026.md",
        "pdf": "PROGRAMA-ELECTORAL-PP-2026-RECONSTRUCCION-NACIONAL.pdf",
        "pdf_resumen": "Resumen-PP-2026.pdf",
    },
}

MD_EXT = ["tables", "sane_lists", "attr_list"]


def md_html(path: pathlib.Path) -> str:
    return markdown.markdown(path.read_text(encoding="utf-8"), extensions=MD_EXT)


def size_label(size: int) -> str:
    if size >= 1_000_000:
        return f"{size / 1_000_000:.1f} MB".replace(".", ",")
    return f"{round(size / 1000)} KB"


def find_pdf(name: str, extra: pathlib.Path | None) -> pathlib.Path:
    for base in ([extra] if extra else []) + PDF_DIRS:
        candidate = base / name
        if candidate.is_file():
            return candidate
    raise SystemExit(
        f"error: {name} not found. Copy the source PDFs into {PDF_DIRS[0]} "
        f"or pass --pdf-dir <dir>."
    )


def build(embed_pdf: bool, pdf_dir: pathlib.Path | None) -> dict:
    out = {}
    for key, spec in DOCS.items():
        resumen = HERE / spec["resumen"]
        verificacion = HERE / spec["verificacion"]
        docs = [
            {"label": "Resumen completo", "title": f"Resumen · {spec['name']}",
             "html": md_html(resumen)},
            {"label": "Informe de verificación", "title": f"Verificación · {spec['name']}",
             "html": md_html(verificacion)},
        ]
        if embed_pdf:
            for label, rel in (("PDF original", spec["pdf"]), ("PDF del resumen", spec["pdf_resumen"])):
                path = find_pdf(rel, pdf_dir)
                data = path.read_bytes()
                docs.append({
                    "label": f"{label} · {size_label(len(data))}",
                    "name": path.name,
                    "b64": base64.b64encode(data).decode("ascii"),
                })
        out[key] = docs
    out["sources"] = {
        "title": "Fuentes y procedencia",
        "html": md_html(HERE / "SOURCES.md"),
    }
    return out


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--no-pdf", action="store_true", help="embed only the text documents")
    ap.add_argument("--out", default=None)
    ap.add_argument("--pdf-dir", default=None, help="directory holding the source PDFs")
    args = ap.parse_args()

    html = BASE.read_text(encoding="utf-8")
    pdf_dir = pathlib.Path(args.pdf_dir) if args.pdf_dir else None
    payload = build(embed_pdf=not args.no_pdf, pdf_dir=pdf_dir)
    blob = json.dumps(payload, ensure_ascii=False, separators=(",", ":"))
    blob = blob.replace("</", "<\\/")  # keep the JSON inert inside <script>

    marker = '<script type="application/json" id="embedded-docs">{}</script>'
    if marker not in html:
        print("error: base HTML has no empty embedded-docs placeholder", file=sys.stderr)
        return 1
    html = html.replace(marker, f'<script type="application/json" id="embedded-docs">{blob}</script>')

    default_name = ("Programas-electorales-2026-standalone.html" if not args.no_pdf
                    else "Programas-electorales-2026-autocontenido.html")
    out = pathlib.Path(args.out) if args.out else BASE.with_name(default_name)
    out.write_text(html, encoding="utf-8")
    print(f"{out}  ({out.stat().st_size / 1_000_000:.1f} MB)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
