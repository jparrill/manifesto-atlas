# Informe de verificación — Resumen Frente Obrero 2026 (segunda gate)

| Dato | Valor |
|---|---|
| Objeto verificado | `Resumen-Frente-Obrero-2026.md` |
| Fuente de contraste | `frente-obrero.md` (texto extraído página a página de `Frente-Obrero-Programa-2026-10.pdf`) |
| Fecha de verificación | 2026-11-17 |
| Resultado global | **PASA** con 0 discrepancias abiertas |
| Verificador | Agente automatizado + cotejo manual por muestreo |

## 1. Método

Extracción con `pdftotext -layout`, división por páginas con marcadores `<!-- page N -->` y cuatro comprobaciones sobre texto normalizado:

- **A. Literalidad de citas.** Las 61 frases entre «» del resumen deben existir textualmente en el documento.
- **B. Página de cada apartado.** Cada uno de los 21 epígrafes del resumen indica el rango de páginas del que procede; se comprobó que el título del apartado aparece en alguna página de ese rango.
- **C. Cifras y datos singulares.** 10 datos numéricos y normativos comprobados uno a uno contra la página citada.
- **D. Cobertura.** Recuento de propuestas numeradas en el original frente a líneas de propuesta en el resumen.

## 2. Resultados

| Comprobación | Elementos | Fallos detectados | Fallos tras corrección |
|---|---|---|---|
| A. Citas literales | 61 | 1 | 0 |
| B. Página de cada apartado | 21 | 7 rangos incompletos | 0 |
| C. Cifras y datos singulares | 10 | 0 | 0 |
| D. Cobertura de propuestas | 133 | 0 | 0 |
| Estructura (21 apartados del índice) | 21 | 0 | 0 |
| Recuento de páginas del PDF | 64 | 0 | 0 |

Detalle de la comprobación C:

| Dato | Página | Estado |
|---|---|---|
| Representación obligatoria en empresas «con diez trabajadores o más» | p. 9 | OK |
| Real Decreto 17/1977 sobre relaciones de trabajo | p. 9 | OK |
| Entrada en la UE en 1986 | p. 11 | OK |
| España, «23 millones de hectáreas» de superficie agraria | p. 37 | OK |
| Objetivo de «más de dos» hijos por mujer | p. 34 | OK |
| Media actual de «1,1 hijos» por mujer | p. 34 | OK |
| Revisión de nacionalidades de «los últimos veinte años» | p. 53 | OK |
| Cuidado «bucodental, fisioterapéutico y oftalmológico» gratuito | p. 42 | OK |
| CIE: operatividad o sustitución | p. 54 | OK |
| Cese de contrato a los «45 años» en el ejército | p. 63 | OK |

Detalle de la comprobación D: el original contiene 133 elementos numerados (la numeración se reinicia en cada apartado, en 20 de los 21 apartados) y el resumen contiene 133 líneas de propuesta numeradas. Correspondencia 1:1.

## 3. Correcciones aplicadas durante la verificación

1. **Cita no literal corregida.** El resumen decía «ha rebasado el límite». El original dice «en España ya hemos rebasado el límite del volumen migratorio que el país puede soportar con creces» (p. 52). Corregido a «ya hemos rebasado el límite del volumen migratorio».
2. **Fecha del documento corregida.** El resumen afirmaba como fecha «7 de octubre de 2026» a partir de los metadatos. Verificado: el PDF **no tiene fecha de creación**; solo consta `ModDate: Wed Oct 7 01:54:50 2026 CEST`. Reescrito como «última modificación 7 de octubre de 2026 (no consta fecha de creación)».
3. **Rangos de página corregidos en 7 apartados.** Siete epígrafes citaban solo la página de las propuestas y omitían la página donde empieza el texto expositivo. Ampliados: 02 (pp. 8-9), 07 (pp. 21-22), 08 (pp. 24-26), 09 (pp. 28-29), 10 (pp. 31-32), 12 (pp. 37-38) y 17 (pp. 51-54).

## 4. Comprobaciones específicas de este documento

Por su naturaleza, se verificaron además tres afirmaciones que el resumen hace **sobre el propio documento** y que son las de mayor riesgo:

| Afirmación del resumen | Verificación | Estado |
|---|---|---|
| «El presente documento no pretende ser un programa electoral, es el programa político del Frente Obrero» | Localizado literalmente en la p. 6 | OK |
| El apartado 01 no contiene propuestas numeradas | Revisadas pp. 5-6: no hay lista numerada | OK |
| Las medidas migratorias se condicionan a «una nueva constitución» | «Siendo conscientes de que habrá que tener una nueva constitución para poder llevarlas a cabo», p. 53 | OK |
| El índice enumera 21 apartados | Índice de la p. 3: 01 a 21 | OK |
| La numeración se reinicia en cada apartado | Comprobado: 20 secuencias independientes que empiezan en 1 | OK |

## 5. Cotejo manual por muestreo

Releídas íntegras en el PDF, palabra por palabra, 10 propuestas de mayor riesgo interpretativo: 3 y 8 del apartado 02 (contratación de extranjeros y derecho de huelga), 1-3 del apartado 05 (fin de la transición ecológica y nuclear), 6 del apartado 10 (vientres de alquiler), 2 del apartado 12 (expropiación de tierras), 3, 6 y 7 del apartado 17 (expulsiones, nacionalidades y retirada de nacionalidad), 4 del apartado 19 (indultos). Las diez coinciden con el resumen; no se detectó ninguna afirmación ausente en el original.

## 6. Limitaciones

- La verificación es **de literalidad y localización**, no de interpretación.
- El documento intercala pasajes de opinión con propuestas. El resumen distingue ambos casos citando el texto expositivo solo cuando explica el sentido de la medida; la atribución de página permite comprobarlo.
- El verificador es el mismo proceso que redactó el resumen; reduce errores de transcripción, pero no sustituye una revisión humana independiente.

## 7. Reproducción

```bash
cd /Users/jparrill/Projects/manifesto-atlas/runs/29n2026
pdftotext -layout "/Users/jparrill/Documents/programas-electorales-opencode-2026-10-30/29N-2026/Frente-Obrero-Programa-2026-10.pdf" text/frente-obrero.txt
python3 tools/pagemark.py text/frente-obrero.txt frente-obrero.md
python3 tools/verify.py
python3 tools/datos.py
```
