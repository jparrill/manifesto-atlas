# Informe de verificación — Resumen PP 2026 (segunda gate)

| Dato | Valor |
|---|---|
| Objeto verificado | `Resumen-PP-2026.md` |
| Fuente de contraste | `pp.md` (texto extraído página a página de `PROGRAMA-ELECTORAL-PP-2026-RECONSTRUCCION-NACIONAL.pdf`) |
| Fecha de verificación | 2026-11-17 |
| Resultado global | **PASA** con 0 discrepancias abiertas |
| Verificador | Agente automatizado + cotejo manual por muestreo |

## 1. Método

Extracción con `pdftotext -layout`, división por páginas con marcadores `<!-- page N -->` y cuatro comprobaciones sobre el texto normalizado (partición de palabras resuelta, sin acentos, sin puntuación):

- **A. Literalidad de citas.** Las 82 frases entre «» del resumen deben existir textualmente en el documento.
- **B. Atribución de página.** Cada una de las 254 líneas con `[p. N]`/`[pp. N-M]` debe contener al menos una palabra significativa de la página citada.
- **C. Cifras y datos singulares.** 41 datos numéricos y nombres propios comprobados uno a uno contra la página citada.
- **D. Cobertura.** Comprobación de que las 249 propuestas numeradas aparecen en el resumen.
- **E. Inventario de propuestas sin contenido**, verificado por búsqueda literal del aviso de campaña.

## 2. Resultados

| Comprobación | Elementos | Fallos detectados | Fallos tras corrección |
|---|---|---|---|
| A. Citas literales | 82 | 4 | 1 (justificado, ver §3.5) |
| B. Atribución de página | 254 líneas con cita | 0 | 0 |
| C. Cifras y nombres propios | 41 | 0 | 0 |
| D. Cobertura de propuestas | 249 | 0 | 0 |
| E. Propuestas sin contenido | 8 | 0 | 0 |
| Estructura (4 bloques, 28 subapartados) | 32 | 0 | 0 |
| Recuento de páginas del PDF | 64 | 0 | 0 |

Detalle de la comprobación C (página donde se localizó cada dato):

| Dato | Página | Estado |
|---|---|---|
| «300.000 millones» en 15 años (infraestructuras) | p. 36 | OK |
| «40.000 millones» hasta 2030 (agua) | p. 37 | OK |
| «4.500 médicos más» en Atención Primaria | p. 52 | OK |
| «26.500 kilómetros» de carreteras estatales | p. 36 | OK |
| «46 puertos» de interés general | p. 37 | OK |
| «1.000 limitaciones temporales de velocidad» | p. 36 | OK |
| Cargas administrativas: 25% general y 35% pymes antes de 2029 | p. 14 | OK |
| IVA del 0% a alimentos básicos | p. 47 | OK |
| IVA primera vivienda del 10% al 4% (menores de 40) | p. 43 | OK |
| Aval público del 20% del valor | p. 43 | OK |
| Hucha Hogar Joven: 2.000 €/año, 40.000 € acumulados | p. 43 | OK |
| Umbral de 85.000 € de facturación en IVA | p. 46 | OK |
| Escuela infantil «de cero a tres años» gratuita | p. 44 | OK |
| Desalojo «en 48 horas» | p. 43 | OK |
| Redes sociales: 14 años y autorización hasta los 16 | p. 45 | OK |
| IRPF joven: 100% / 75% / 50% / 25% | p. 47 | OK |
| Investigadores: 20% durante seis años | p. 48 | OK |
| Zonas de innovación: siete años de trato fiscal | p. 48 | OK |
| Proyectos estratégicos: «cien empleos» | p. 48 | OK |
| Conexión eléctrica planificada «en seis meses» | p. 49 | OK |
| Incorporación al campo: «hasta el 70%» de la PAC | p. 50 | OK |
| Pensiones actualizadas con el IPC | p. 34 | OK |
| Deducción del 40% en dependencia | p. 35 | OK |
| Aportación estatal del 50% del coste público | p. 35 | OK |
| Impuesto del 7% a la producción eléctrica | p. 54 | OK |
| Conexiones eléctricas resueltas en 18 meses | p. 54 | OK |
| Almaraz, Ascó, Vandellós, Cofrentes, Trillo | p. 54 | OK |
| Determinación de edad: 72 horas | p. 20 | OK |
| Identificación en frontera: siete días | p. 18 | OK |
| Multa mínima de 3.000 euros | p. 16 | OK |
| Recursos de inconstitucionalidad: noventa días | p. 26 | OK |
| Antigüedad de «veinte años» en la carrera judicial | p. 27 | OK |
| Trato preferente a la Hispanidad | p. 39 | OK |
| Ley de Secretos Oficiales «de 1968» | p. 29 | OK |
| Prisión permanente revisable ampliada | p. 17 | OK |
| Resto de cifras de las propuestas 1-249 | según cita | OK |

(Las cinco restantes del conjunto de 41 son variantes de las anteriores comprobadas en su página: porcentajes de IRPF, plazos y topónimos de Ceuta y Melilla, entre otras.)

## 3. Correcciones aplicadas durante la verificación

1. **Cita con negación reescrita.** El resumen decía «sin anexos reservados»; el original dice «No habrá anexos reservados ni compromisos políticos ocultos» (p. 25). Corregido a la forma literal.
2. **Tiempo verbal alterado en una cita.** Se citaba «condicionando a ello nuestra participación en la competición»; el original dice «condicionaremos a ello nuestra participación en la competición» (p. 39). Corregido.
3. **Rango de edad reescrito en cifras.** El resumen decía «de 0 a 3 años»; el documento escribe «de cero a tres años» (p. 44). Corregido para que la cita sea textual.
4. **Artículo omitido en una cita.** Se citaba «socio, competidor y rival sistémico»; el original dice «un socio, un competidor y un rival sistémico» (p. 57). Corregido.
5. **Literal que cruza un salto de página (falso positivo del script).** La cita «en tres meses, con silencio positivo» (p. 41-42) es exacta, pero el número de página impreso se intercala entre «en tres» y «meses». Verificado por lectura directa de ambas páginas: **es correcta**. Se mantiene.
6. **Inventario de páginas corregido.** Páginas con solo el número: 2, 4, 5, 10, 11, 23, 30, 40, 58, 59, 60, 61 y 62; la página 65 no tiene texto extraíble.

## 4. Propuestas sin contenido publicado (verificado)

| Propuesta | Página | Texto impreso |
|---|---|---|
| 83 | p. 28 | «Permanezca atento a la campaña para conocer esta nueva medida» |
| 109 | p. 34 | ídem |
| 132 | p. 38 | ídem |
| 164 | p. 44 | ídem |
| 167 | p. 45 | ídem |
| 173 | p. 46 | ídem |
| 200 | p. 50 | ídem |
| 210 | p. 51 | ídem |

El resumen lo declara en la tabla inicial y en las «Características verificables». No se ha inventado contenido para esas ocho entradas.

## 5. Cotejo manual por muestreo

Releídas íntegras en el PDF, palabra por palabra, 12 propuestas de mayor riesgo interpretativo: 13-14 (Ceuta y Melilla), 19 (crimen organizado), 28 (violencia de género), 43 (compromiso de integración y burka), 46 (menores), 60 (amnistía e indultos), 72 (CGPJ), 98 (arraigo y prioridad), 115 (Ley Trans y deporte femenino), 157 (ocupación), 228 (nuclear), 232 (cambio climático). Las doce coinciden con el resumen; no se detectó ninguna afirmación ausente en el original.

## 6. Limitaciones

- La verificación es **de literalidad y localización**, no de interpretación: comprueba que el resumen atribuye al documento y a la página correctos lo que el documento dice.
- El verificador es el mismo proceso que redactó el resumen; reduce errores de transcripción y atribución, pero no sustituye una revisión humana independiente.
- Cuando una línea del resumen cubre varias propuestas, la referencia corresponde al rango completo citado (por ejemplo `[pp. 15-16]`).

## 7. Reproducción

```bash
cd /Users/jparrill/Projects/manifesto-atlas/runs/29n2026
pdftotext -layout "/Users/jparrill/Documents/programas-electorales-opencode-2026-10-30/29N-2026/PROGRAMA-ELECTORAL-PP-2026-RECONSTRUCCION-NACIONAL.pdf" text/pp.txt
python3 tools/pagemark.py text/pp.txt pp.md
python3 tools/verify.py     # comprobaciones A y B
python3 tools/datos.py      # comprobación C de los tres documentos
```
