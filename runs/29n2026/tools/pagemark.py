import sys, pathlib
src = pathlib.Path(sys.argv[1])
dst = pathlib.Path(sys.argv[2])
raw = src.read_text(encoding="utf-8", errors="replace")
pages = raw.split("\f")
out = []
for i, p in enumerate(pages, start=1):
    body = p.strip("\n")
    if not body.strip():
        body = "_(sin texto extraíble)_"
    out.append(f"<!-- page {i} -->\n\n{body}\n")
dst.write_text("\n".join(out), encoding="utf-8")
print(dst, len(pages))
