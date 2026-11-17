# Run 29-N 2026 — tres programas, resúmenes factuales y verificación

Ejecución completa sobre los tres programas depositados en
`/Users/jparrill/Documents/programas-electorales-opencode-2026-10-30/29N-2026/`:
conversión a texto, resumen factual medida a medida, segunda gate de verificación
y una vista HTML navegable.

Este directorio es el **registro de la ejecución**. Las fuentes originales no se
copian aquí: su procedencia y huella están en [SOURCES.md](SOURCES.md).

## 1. Entradas y alcance

| Partido | Páginas | Propuestas numeradas | Naturaleza |
|---|---|---|---|
| Frente Obrero | 64 | 133 | Programa político, declarado **no electoral** (p. 6) |
| VOX | 66 | 115 | Programa electoral nacional («ES SENTIDO COMÚN») |
| PP | 64 | 249 | Programa electoral nacional («Reconstrucción Nacional»); 8 propuestas sin contenido |

## 2. Productos

| Fase | Fichero | Contenido |
|---|---|---|
| 1 · Conversión | `text/*.txt` | Texto plano con maquetación (`pdftotext -layout`) |
| 1 · Conversión | `frente-obrero.md`, `vox.md`, `pp.md` | Markdown con marcador `<!-- page N -->` por página |
| 2 · Resumen | `Resumen-Frente-Obrero-2026.md` | Resumen factual, página de origen por propuesta |
| 2 · Resumen | `Resumen-VOX-2026.md` | Ídem |
| 2 · Resumen | `Resumen-PP-2026.md` | Ídem |
| 2 · Resumen | `Resumen-*.pdf` | Versión imprimible A4 (Chrome headless) |
| 3 · Verificación | `Verificacion-Resumen-*.md` | Método, resultados, correcciones y límites |
| 4 · Vista | `Programas-electorales-2026.html` | Tres pestañas, 13 bloques desplegables, 106 fichas con contexto; enlaza los ficheros vecinos |
| 4 · Vista portátil | `Programas-electorales-2026-autocontenido.html` | La misma vista con los resúmenes y los informes **embebidos** (192 KB): se abre sola, sin los ficheros al lado ni conexión |
| 4 · Portable | `Programas-29N-2026-compendio.pdf` | Compendio A4 de 45 páginas con portada, índice, los tres resúmenes y los tres informes |
| Soporte | `SOURCES.md` | URL oficial, tamaño y SHA-256 de los PDF, con la comprobación de descarga |
| Herramientas | `tools/` | Extracción, comprobaciones, PDF de resúmenes, vista autocontenida y compendio |

## 3. Resultado de la segunda gate

| Documento | Citas literales | Líneas con página | Datos numéricos | Cobertura | Discrepancias |
|---|---|---|---|---|---|
| VOX | 79 | 117 | 18 | 115/115 | 0 |
| PP | 82 | 254 | 41 | 249/249 | 0 |
| Frente Obrero | 61 | 21 apartados | 10 | 133/133 | 0 |

Se aplicaron 12 correcciones durante la verificación, todas de literalidad de cita,
página de origen o inventario de páginas; ninguna cambia el sentido de una medida.
El detalle está en cada `Verificacion-*.md`.

## 4. Reglas aplicadas

- Se recoge lo que el documento afirma, sin valorar viabilidad, legalidad ni conveniencia.
- No se añaden datos externos a los resúmenes. Los recuadros «Contexto» de la vista
  HTML son la única capa añadida y están marcados como tales, también en el pie de página.
- Las frases entre «» son literales del documento y se comprueban una a una.
- Cada propuesta lleva la página del PDF de la que procede.
- Cuando el documento no dice algo (las 8 propuestas del PP remitidas a la campaña),
  se indica en lugar de rellenarlo.

## 5. Procedencia de las fuentes

Los tres PDF se descargaron de la URL oficial de cada partido y su SHA-256 coincide
byte a byte con el del fichero analizado (comprobado el 2026-11-17). El detalle y el
comando de verificación están en [SOURCES.md](SOURCES.md). Es una comprobación
puntual del mantenedor: si un partido sustituye el fichero en la misma ruta, la
huella dejará de coincidir y la ejecución quedará invalidada para ese documento.

## 6. Limitaciones

- La verificación es automática y de literalidad/localización: no juzga interpretación
  semántica ni la exactitud de las cifras que cada partido declara.
- El verificador es el mismo proceso que redactó los resúmenes. Reduce errores de
  transcripción, atribución y cobertura, pero no equivale a una revisión humana
  independiente ni a la validación con dos pases de contexto limpio que describe
  el contrato del proyecto.
- Los resúmenes condensan propuestas afines cuando el documento repite el mismo
  contenido; la cobertura por número de propuesta está comprobada al 100%.
- La vista HTML la ha escrito un agente sobre la plantilla aportada. El contrato del
  proyecto reserva la autoría de presentación al mantenedor: se conserva aquí como
  artefacto de la ejecución, no como activo aprobado.

## 7. Reproducción

```bash
cd runs/29n2026
python3 tools/verify.py     # citas literales y atribución de página
python3 tools/datos.py      # datos numéricos y nombres propios

# PDF de los resúmenes (requiere Chrome y el módulo markdown de Python)
python3 tools/md2pdf.py Resumen-VOX-2026.md

# Vista portátil con los textos embebidos
python3 tools/build_standalone.py --no-pdf

# Compendio PDF con índice y los seis documentos
python3 tools/build_compendium.py
```

Los dos generadores necesitan el módulo `markdown` de Python y Chrome para imprimir.
