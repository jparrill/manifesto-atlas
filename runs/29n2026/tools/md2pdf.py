#!/usr/bin/env python3
"""Convert the manifesto summaries (Markdown) to print-ready A4 PDF via headless Chrome."""
import subprocess, sys, pathlib, tempfile

import markdown

CHROME = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

CSS = """
@page { size: A4; margin: 18mm 16mm 16mm 16mm; }
* { box-sizing: border-box; }
body { font-family: "Iowan Old Style", Georgia, "Times New Roman", serif;
       font-size: 10.2pt; line-height: 1.45; color: #14161a; margin: 0; }
h1 { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
     font-size: 21pt; line-height: 1.15; margin: 0 0 4mm; letter-spacing: -0.02em;
     color: #0d1b2a; border-bottom: 2.5pt solid #0d1b2a; padding-bottom: 3mm; }
h2 { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
     font-size: 12.5pt; margin: 8mm 0 2.5mm; color: #1b2a41;
     border-left: 3pt solid #8a9bb5; padding-left: 3mm;
     page-break-after: avoid; break-after: avoid; }
h2:first-of-type { margin-top: 6mm; }
p { margin: 0 0 2.4mm; orphans: 3; widows: 3; }
ul { margin: 0 0 2.6mm; padding-left: 6mm; }
li { margin: 0 0 1.2mm; }
strong { color: #0d1b2a; }
hr { border: 0; border-top: 0.6pt solid #c8ced8; margin: 6mm 0; }
table { width: 100%; border-collapse: collapse; margin: 0 0 4mm; font-size: 9.2pt; }
th, td { border: 0.5pt solid #b9c1cd; padding: 1.6mm 2.2mm; text-align: left;
         vertical-align: top; }
th { background: #eef1f6; font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
     font-size: 8.8pt; }
tr { page-break-inside: avoid; break-inside: avoid; }
code { font-family: "SF Mono", Menlo, monospace; font-size: 8.6pt;
       background: #f1f3f7; padding: 0.3mm 1mm; border-radius: 1.5pt; }
"""

def convert(md_path: pathlib.Path, pdf_path: pathlib.Path) -> None:
    body = markdown.markdown(
        md_path.read_text(encoding="utf-8"),
        extensions=["tables", "sane_lists", "attr_list"],
    )
    html = (
        "<!doctype html><html lang=\"es\"><head><meta charset=\"utf-8\">"
        f"<title>{md_path.stem}</title><style>{CSS}</style></head><body>"
        f"{body}</body></html>"
    )
    with tempfile.NamedTemporaryFile("w", suffix=".html", delete=False, encoding="utf-8") as fh:
        fh.write(html)
        html_path = fh.name
    subprocess.run(
        [CHROME, "--headless=new", "--disable-gpu", "--no-pdf-header-footer",
         "--virtual-time-budget=4000",
         f"--print-to-pdf={pdf_path}", pathlib.Path(html_path).as_uri()],
        check=True, capture_output=True,
    )
    pathlib.Path(html_path).unlink(missing_ok=True)

if __name__ == "__main__":
    for arg in sys.argv[1:]:
        src = pathlib.Path(arg).resolve()
        convert(src, src.with_suffix(".pdf"))
        print("PDF:", src.with_suffix(".pdf"))
