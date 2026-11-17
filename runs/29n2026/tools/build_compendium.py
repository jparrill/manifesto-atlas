#!/usr/bin/env python3
"""Build one portable PDF that contains every text document of the 29-N run.

Cover + index + the three summaries + the three verification reports, all in a
single A4 file with internal links. The source PDFs are not embedded: they are
identified by their official URL and SHA-256 in SOURCES.md, which is also
printed on the cover.

    python3 tools/build_compendium.py [--out FILE]
"""
import argparse
import datetime
import pathlib
import subprocess
import tempfile

import markdown

HERE = pathlib.Path(__file__).resolve().parent.parent
CHROME = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

DOCS = [
    ("resumen-fo", "Resumen · Frente Obrero", "Resumen-Frente-Obrero-2026.md",
     "64 páginas · 133 propuestas · programa político, no electoral"),
    ("resumen-vox", "Resumen · VOX", "Resumen-VOX-2026.md",
     "66 páginas · 115 propuestas · programa electoral nacional"),
    ("resumen-pp", "Resumen · Partido Popular", "Resumen-PP-2026.md",
     "64 páginas · 249 propuestas · 8 remitidas a la campaña sin contenido"),
    ("verif-fo", "Verificación · Frente Obrero", "Verificacion-Resumen-Frente-Obrero-2026.md",
     "61 citas literales · 21 apartados · 133/133 propuestas cubiertas"),
    ("verif-vox", "Verificación · VOX", "Verificacion-Resumen-VOX-2026.md",
     "79 citas literales · 117 referencias de página · 115/115 propuestas cubiertas"),
    ("verif-pp", "Verificación · Partido Popular", "Verificacion-Resumen-PP-2026.md",
     "82 citas literales · 254 referencias de página · 249/249 propuestas cubiertas"),
]

CSS = """
@page { size: A4; margin: 20mm 17mm 18mm 17mm; }
* { box-sizing: border-box; }
body { font-family: "Iowan Old Style", Georgia, "Times New Roman", serif;
       font-size: 10.2pt; line-height: 1.45; color: #14161a; margin: 0; }
.cover { page-break-after: always; padding-top: 18mm; }
.cover .eyebrow { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
  font-size: 8.5pt; letter-spacing: .18em; text-transform: uppercase; color: #5a6478; margin: 0 0 6mm; }
.cover h1 { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
  font-size: 30pt; line-height: 1.05; letter-spacing: -0.03em; margin: 0 0 5mm; color: #0d1b2a; }
.cover .lead { font-size: 12pt; color: #34405a; margin: 0 0 10mm; max-width: 130mm; }
.cover h2 { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
  font-size: 10pt; letter-spacing: .1em; text-transform: uppercase; color: #5a6478;
  margin: 8mm 0 3mm; border: 0; padding: 0; }
table.meta { width: 100%; border-collapse: collapse; font-size: 9pt; }
table.meta td { border-top: 0.5pt solid #ccd2dd; padding: 2mm 2mm 2mm 0; vertical-align: top; }
table.meta td:first-child { width: 34mm; color: #5a6478; }
.hash { font-family: "SF Mono", Menlo, monospace; font-size: 7pt; word-break: break-all; }
.warn { margin-top: 8mm; padding: 4mm 5mm; border-left: 2.5pt solid #b45309;
        background: #fdf6ec; font-size: 9pt; color: #4a3a20; }
.toc { page-break-after: always; }
.toc h1 { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
  font-size: 18pt; margin: 0 0 6mm; color: #0d1b2a; border-bottom: 1.5pt solid #0d1b2a; padding-bottom: 2mm; }
.toc ol { list-style: none; margin: 0; padding: 0; counter-reset: d; }
.toc li { counter-increment: d; border-bottom: 0.5pt solid #dfe3ea; padding: 3mm 0; }
.toc li::before { content: counter(d) ". "; color: #8b95ab;
  font-family: -apple-system, "Helvetica Neue", Arial, sans-serif; font-size: 9pt; }
.toc a { color: #14161a; text-decoration: none; font-weight: 600; }
.toc span { display: block; color: #5a6478; font-size: 8.6pt; margin-top: 1mm; }
.doc { page-break-before: always; }
.doc > .back { font-size: 8.4pt; color: #8b95ab; margin: 0 0 4mm;
  font-family: -apple-system, "Helvetica Neue", Arial, sans-serif; }
.doc > .back a { color: #5a6478; text-decoration: none; }
h1 { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
     font-size: 20pt; line-height: 1.15; margin: 0 0 5mm; letter-spacing: -0.02em;
     color: #0d1b2a; border-bottom: 2pt solid #0d1b2a; padding-bottom: 2.5mm; }
h2 { font-family: -apple-system, "Helvetica Neue", Arial, sans-serif;
     font-size: 12pt; margin: 7mm 0 2mm; color: #1b2a41;
     border-left: 2.5pt solid #8a9bb5; padding-left: 2.5mm;
     page-break-after: avoid; break-after: avoid; }
p { margin: 0 0 2.2mm; orphans: 3; widows: 3; }
ul, ol { margin: 0 0 2.4mm; padding-left: 6mm; }
li { margin: 0 0 1.1mm; }
hr { border: 0; border-top: 0.6pt solid #c8ced8; margin: 5mm 0; }
table { width: 100%; border-collapse: collapse; margin: 0 0 3.5mm; font-size: 8.8pt; }
th, td { border: 0.5pt solid #b9c1cd; padding: 1.4mm 2mm; text-align: left; vertical-align: top; }
th { background: #eef1f6; font-family: -apple-system, "Helvetica Neue", Arial, sans-serif; font-size: 8.4pt; }
tr { page-break-inside: avoid; break-inside: avoid; }
code { font-family: "SF Mono", Menlo, monospace; font-size: 8.4pt; background: #f1f3f7;
       padding: 0.2mm 1mm; border-radius: 1.5pt; }
"""

SOURCES = [
    ("Frente Obrero", "https://frenteobrero.es/documentos/programa-frente-obrero.pdf",
     1_202_338, "7b9f423f2c1ada187b6f42ed8dbf4e025356254d7a6b81830e114b818ef2e282"),
    ("VOX", "https://www.voxespana.es/wp-content/uploads/2026/11/Programa-completo-A5-1.pdf",
     14_536_418, "065abf6720be2f7b19b8e911920c2ee2ce5c2059d0ae8bf09770bc790ed11b65"),
    ("PP", "https://www.pp.es/wp-content/uploads/2026/11/PROGRAMA-ELECTORAL-PP-2026-RECONSTRUCCION-NACIONAL.pdf",
     12_304_990, "ae0c35d7ef3faf540acaa79d47abbe3ae2bf56a3a755c4e78a8f5bf3ea49789c"),
]


def cover() -> str:
    rows = "".join(
        f"<tr><td>{party}</td><td><span class='hash'>{url}</span><br>"
        f"{size:,} B · SHA-256 <span class='hash'>{sha}</span></td></tr>"
        .replace(",", ".")
        for party, url, size, sha in SOURCES
    )
    today = datetime.date.today().isoformat()
    return f"""
<section class="cover">
  <p class="eyebrow">Programas electorales · 29 de noviembre de 2026</p>
  <h1>Compendio de resúmenes<br>y verificación</h1>
  <p class="lead">Tres programas nacionales resumidos medida a medida, con la página
  de origen de cada propuesta, y el informe que contrasta cada resumen contra el
  documento original.</p>
  <h2>Qué contiene</h2>
  <table class="meta">
    <tr><td>Documentos</td><td>3 resúmenes + 3 informes de verificación</td></tr>
    <tr><td>Propuestas</td><td>133 (Frente Obrero), 115 (VOX) y 249 (PP)</td></tr>
    <tr><td>Comprobaciones</td><td>222 citas literales verificadas y 69 datos numéricos cotejados</td></tr>
    <tr><td>Discrepancias abiertas</td><td>0</td></tr>
    <tr><td>Generado</td><td>{today}</td></tr>
  </table>
  <h2>Fuentes oficiales</h2>
  <table class="meta">{rows}</table>
  <p class="warn"><b>Alcance.</b> Este compendio describe lo que cada documento afirma;
  no verifica las cifras que declaran los partidos ni valora la viabilidad o
  conveniencia de sus propuestas. La verificación comprueba literalidad de citas,
  atribución de página y cobertura de propuestas. El mismo proceso que redactó los
  resúmenes ejecutó la verificación: reduce errores de transcripción y atribución,
  pero no equivale a una revisión humana independiente. Los tres PDF de origen se
  descargaron de las URL indicadas y sus huellas SHA-256 coinciden con las de los
  ficheros analizados.</p>
</section>"""


def toc() -> str:
    items = "".join(
        f'<li><a href="#{anchor}">{title}</a><span>{note}</span></li>'
        for anchor, title, _, note in DOCS
    )
    return f"""<section class="toc"><h1>Índice</h1><ol>{items}</ol>
<p style="margin-top:8mm;font-size:8.6pt;color:#5a6478">Los títulos del índice son
enlaces internos. Cada documento empieza en una página nueva.</p></section>"""


def body() -> str:
    out = []
    for anchor, title, filename, _ in DOCS:
        html = markdown.markdown((HERE / filename).read_text(encoding="utf-8"),
                                 extensions=["tables", "sane_lists", "attr_list"])
        out.append(f'<section class="doc" id="{anchor}">'
                   f'<p class="back"><a href="#indice">↑ Índice</a> · {filename}</p>'
                   f'{html}</section>')
    return "".join(out)


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--out", default=str(HERE / "Programas-29N-2026-compendio.pdf"))
    args = ap.parse_args()

    html = (f'<!doctype html><html lang="es"><head><meta charset="utf-8">'
            f'<title>Compendio · programas 29-N 2026</title><style>{CSS}</style></head>'
            f'<body><div id="indice">{cover()}{toc()}</div>{body()}</body></html>')

    with tempfile.NamedTemporaryFile("w", suffix=".html", delete=False, encoding="utf-8") as fh:
        fh.write(html)
        tmp = fh.name

    subprocess.run([CHROME, "--headless=new", "--disable-gpu", "--no-pdf-header-footer",
                    "--virtual-time-budget=8000", f"--print-to-pdf={args.out}",
                    pathlib.Path(tmp).as_uri()], check=True, capture_output=True)
    pathlib.Path(tmp).unlink(missing_ok=True)
    size = pathlib.Path(args.out).stat().st_size
    print(f"{args.out}  ({size / 1_000_000:.1f} MB)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
