import re, pathlib, unicodedata, json, sys

def load_pages(f):
    t = pathlib.Path(f).read_text(encoding="utf-8")
    parts = re.split(r"<!-- page (\d+) -->", t)
    return {int(parts[i]): parts[i+1] for i in range(1, len(parts), 2)}

def dehyph(s):
    s = re.sub(r"(\w)-\s*\n\s*(\w)", r"\1\2", s)
    s = re.sub(r"\s+", " ", s)
    return s

def norm(s):
    s = dehyph(s)
    s = s.lower()
    s = unicodedata.normalize("NFD", s)
    s = "".join(c for c in s if unicodedata.category(c) != "Mn")
    s = re.sub(r"[^a-z0-9 ]", " ", s)
    s = re.sub(r"\s+", " ", s)
    return s.strip()

STOP = set("""de la el los las que en y a del al con por para se su sus un una unos unas es son como
mas menos entre sin sobre tras durante desde hasta cuando donde quien quienes cual cuales este esta estos
estas ese esa eso otro otra otros otras todo toda todos todas ser esta estan sera no ni o pero si ya
lo le les nos me te mi tu mis tus nuestro nuestra""".split())

def keywords(text, n=4, minlen=6):
    words = [w for w in re.findall(r"[A-Za-zÁÉÍÓÚÜÑáéíóúüñ]+", text)]
    words = [w for w in words if len(w) >= minlen and norm(w) not in STOP]
    seen = []
    for w in words:
        nw = norm(w)
        if nw and nw not in seen:
            seen.append(nw)
    return seen[:n]

def main(src_md, sum_md):
    pages = load_pages(src_md)
    npages = {p: norm(b) for p, b in pages.items()}
    full = norm(" ".join(pages[p] for p in sorted(pages)))
    text = pathlib.Path(sum_md).read_text(encoding="utf-8")

    # Check 1: quoted literals «...»
    quotes = re.findall(r"«([^»]{8,200})»", text)
    qfail = []
    for q in quotes:
        nq = norm(q)
        if not nq:
            continue
        if nq not in full:
            qfail.append(q)
    # Check 2: page attribution
    line_fail = []
    total_lines = 0
    for line in text.split("\n"):
        m = re.findall(r"\[pp?\.\s*([0-9\-–,\s]+)\]", line)
        if not m:
            continue
        total_lines += 1
        cited = []
        for grp in m:
            for part in re.split(r"[,]", grp):
                part = part.strip()
                if re.match(r"^\d+\s*[-–]\s*\d+$", part):
                    a, b = re.split(r"[-–]", part)
                    cited += list(range(int(a), int(b) + 1))
                elif part.isdigit():
                    cited.append(int(part))
        kws = keywords(line)
        ok = None
        for pg in cited:
            body = npages.get(pg, "")
            if any(k in body for k in kws):
                ok = pg
                break
        if ok is None:
            # try neighbour pages (off-by-one in pagination)
            alt = [pg for pg in range(min(cited)-1, max(cited)+2) if 1 <= pg <= max(npages)]
            for pg in alt:
                body = npages.get(pg, "")
                if any(k in body for k in kws):
                    ok = pg
                    break
            if ok is None:
                line_fail.append((cited, kws, line[:90]))
            else:
                line_fail.append(("OFFBY", cited, ok, kws, line[:70]))
    print(f"== {sum_md} ==")
    print(f"citas literales: {len(quotes)}, no encontradas: {len(qfail)}")
    for q in qfail:
        print("   QUOTE-FAIL:", q[:120])
    print(f"lineas con cita de pagina: {total_lines}, sin coincidencia: {len([x for x in line_fail if x[0]!='OFFBY'])}, desalineadas: {len([x for x in line_fail if x[0]=='OFFBY'])}")
    for f in line_fail:
        print("   ", f)

for a, b in [("vox.md","Resumen-VOX-2026.md"), ("pp.md","Resumen-PP-2026.md"), ("frente-obrero.md","Resumen-Frente-Obrero-2026.md")]:
    main(a, b)
