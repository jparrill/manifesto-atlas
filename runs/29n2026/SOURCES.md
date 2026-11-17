# Fuentes del run 29-N 2026 — procedencia y huellas

Los tres PDF son **inmutables**: no se copian al repositorio ni se modifican. Se
registra su URL oficial, tamaño y SHA-256 para poder comprobar que el texto
extraído y los resúmenes corresponden exactamente a esos bytes.

## Procedencia verificada

Descarga realizada el **2026-11-17** desde las URL oficiales indicadas por el
mantenedor. En las tres, la huella del fichero descargado coincide byte a byte
con la del PDF analizado: la procedencia es verificable, no declarativa.

| Partido | URL oficial | Tamaño | SHA-256 | Descarga |
|---|---|---|---|---|
| Frente Obrero | https://frenteobrero.es/documentos/programa-frente-obrero.pdf | 1.202.338 B | `7b9f423f2c1ada187b6f42ed8dbf4e025356254d7a6b81830e114b818ef2e282` | 200 · idéntica |
| VOX | https://www.voxespana.es/wp-content/uploads/2026/11/Programa-completo-A5-1.pdf | 14.536.418 B | `065abf6720be2f7b19b8e911920c2ee2ce5c2059d0ae8bf09770bc790ed11b65` | 200 · idéntica |
| PP | https://www.pp.es/wp-content/uploads/2026/11/PROGRAMA-ELECTORAL-PP-2026-RECONSTRUCCION-NACIONAL.pdf | 12.304.990 B | `ae0c35d7ef3faf540acaa79d47abbe3ae2bf56a3a755c4e78a8f5bf3ea49789c` | 200 · idéntica |

Comprobación reproducible:

```bash
for u in \
  "https://frenteobrero.es/documentos/programa-frente-obrero.pdf" \
  "https://www.voxespana.es/wp-content/uploads/2026/11/Programa-completo-A5-1.pdf" \
  "https://www.pp.es/wp-content/uploads/2026/11/PROGRAMA-ELECTORAL-PP-2026-RECONSTRUCCION-NACIONAL.pdf"; do
  curl -sSL "$u" | shasum -a 256
done
```

Qué significa y qué no:

- **Significa** que los bytes analizados son los que el partido publica hoy en esa
  URL, y que la huella permite detectar cualquier sustitución posterior.
- **No significa** certificación de origen por un tercero: es una comprobación
  puntual del mantenedor, no una firma. Si el partido reemplaza el fichero en la
  misma ruta, la huella dejará de coincidir y la ejecución quedará invalidada para
  ese documento.

Ruta original de los ficheros analizados:

```
/Users/jparrill/Documents/programas-electorales-opencode-2026-10-30/29N-2026/
```

## Metadatos declarados por los PDF

| Documento | Páginas | Fecha en el documento | Metadatos del fichero |
|---|---|---|---|
| Frente Obrero | 64 | sin fecha impresa | sin `CreationDate`; `ModDate` 2026-10-07 |
| VOX | 66 | «29 de noviembre de 2026» (p. 3) | productor 3-Heights PDF Merge Split Shell 6.12.1.11 |
| PP | 64 | «próximo 29 de noviembre» (p. 6) | maquetación InDesign; PDF/X |

## Identidad editorial

- **Frente Obrero**: el propio documento declara en la p. 6 que «no pretende ser
  un programa electoral», sino el programa político del movimiento. Ámbito nacional.
- **VOX**: programa electoral de ámbito nacional con cabecera «ES SENTIDO COMÚN».
  No es el documento autonómico de Castilla y León que figuraba en el HTML de
  referencia.
- **PP**: «Reconstrucción Nacional — el programa del cambio necesario para España»,
  4 bloques y 249 propuestas numeradas, de las que 8 remiten a la campaña sin
  contenido publicado.

## Extracción

```bash
pdftotext -layout "<pdf>" text/<nombre>.txt
python3 tools/pagemark.py text/<nombre>.txt <nombre>.md
```

No fue necesario OCR: los tres documentos tienen capa de texto. Las páginas sin
texto extraíble son imágenes a sangre, portadillas o páginas en blanco y se
enumeran en cada resumen.
