import re, pathlib, unicodedata
def norm(s):
    s=re.sub(r"(\w)-\s*\n\s*(\w)", r"\1\2", s).lower()
    s=unicodedata.normalize("NFD",s); s="".join(c for c in s if unicodedata.category(c)!="Mn")
    return re.sub(r"\s+"," ",re.sub(r"[^a-z0-9 ]"," ",s))
def load(f):
    t=pathlib.Path(f).read_text(); p=re.split(r"<!-- page (\d+) -->", t)
    return {int(p[i]):norm(p[i+1]) for i in range(1,len(p),2)}
D={k:load(k) for k in ["vox.md","pp.md","frente-obrero.md"]}
CHECKS=[
("pp.md",[36],"300.000 millones"),("pp.md",[36],"15 años"),("pp.md",[37],"40.000 millones"),
("pp.md",[52],"4.500 médicos"),("pp.md",[36],"26.500 kilómetros"),("pp.md",[37],"46 puertos"),
("pp.md",[36],"1.000 limitaciones temporales"),("pp.md",[14],"25%"),("pp.md",[14],"35%"),
("pp.md",[14],"2029"),("pp.md",[47],"iva del 0"),("pp.md",[43],"del 10 al 4"),
("pp.md",[43],"20% del valor"),("pp.md",[43],"2.000 euros"),("pp.md",[43],"40.000 euros"),
("pp.md",[46],"85.000 euros"),("pp.md",[44],"cero a tres"),("pp.md",[43],"48 horas"),
("pp.md",[45],"menor de 14"),("pp.md",[45],"16"),("pp.md",[47],"75% el segundo"),
("pp.md",[48],"20% durante seis años"),("pp.md",[48],"siete años"),("pp.md",[48],"cien empleos"),
("pp.md",[49],"seis meses"),("pp.md",[50],"70%"),("pp.md",[34],"ipc"),
("pp.md",[35],"40%"),("pp.md",[35],"50% del coste publico"),("pp.md",[54],"7%"),
("pp.md",[54],"18 meses"),("pp.md",[54],"almaraz"),("pp.md",[54],"trillo"),
("pp.md",[20],"72 horas"),("pp.md",[18],"siete dias"),("pp.md",[16],"3.000 euros"),
("pp.md",[26],"noventa dias"),("pp.md",[27],"veinte años"),("pp.md",[39],"hispanidad"),
("pp.md",[29],"1968"),("pp.md",[17],"prision permanente revisable"),
("vox.md",[9],"50 euros"),("vox.md",[15],"5 millones"),("vox.md",[24],"24 horas"),
("vox.md",[24],"100% del itp"),("vox.md",[24],"0% de iva"),("vox.md",[33],"mercosur"),
("vox.md",[34],"pac"),("vox.md",[46],"ocon"),("vox.md",[48],"gibraltar"),
("vox.md",[55],"doce vocales"),("vox.md",[61],"cuencas"),("vox.md",[37],"tarjeta sanitaria unica"),
("vox.md",[27],"fin de la cita previa"),("vox.md",[28],"derogar tres normas"),
("vox.md",[42],"burka"),("vox.md",[55],"sedicion"),("vox.md",[45],"prision permanente revisable"),
("vox.md",[51],"0 a 3"),
("frente-obrero.md",[9],"diez trabajadores"),("frente-obrero.md",[34],"mas de dos"),
("frente-obrero.md",[53],"veinte años"),("frente-obrero.md",[9],"17/1977"),
("frente-obrero.md",[11],"1986"),("frente-obrero.md",[37],"23 millones de hectareas"),
("frente-obrero.md",[34],"1 1 hijos"),("frente-obrero.md",[42],"bucodental"),
("frente-obrero.md",[54],"cie"),("frente-obrero.md",[63],"45 años"),
]
bad=0; per={}
for doc,pgs,ph in CHECKS:
    ok=any(norm(ph) in D[doc][p] for p in pgs)
    per[doc]=per.get(doc,[0,0]); per[doc][0]+=1
    if not ok:
        bad+=1; per[doc][1]+=1; print("FAIL",doc,pgs,ph)
print("TOTAL", len(CHECKS), "FAIL", bad, per)
