# Informe de verificación — Resumen VOX 2026 (segunda gate)

| Dato | Valor |
|---|---|
| Objeto verificado | `Resumen-VOX-2026.md` |
| Fuente de contraste | `vox.md` (texto extraído página a página de `Programa-completo-A5-1.pdf`) |
| Fecha de verificación | 2026-11-17 |
| Resultado global | **PASA** con 0 discrepancias abiertas |
| Verificador | Agente automatizado + cotejo manual por muestreo (el mismo proceso que generó el resumen; ver «Limitaciones») |

## 1. Método

El texto del PDF se extrajo con `pdftotext -layout` y se dividió en páginas mediante los caracteres de salto de página, insertando un marcador `<!-- page N -->` por página. Sobre ese texto se aplicaron cuatro comprobaciones:

- **A. Literalidad de citas.** Toda frase entre «» del resumen debe aparecer textualmente en el documento, tras normalizar guiones de partición de palabra, espacios, mayúsculas y acentos. Detecta invenciones o paráfrasis presentadas como cita.
- **B. Atribución de página.** Cada línea del resumen que cita `[p. N]` o `[pp. N-M]` debe contener al menos una palabra significativa presente en esa página o en el rango citado. Detecta páginas equivocadas.
- **C. Cifras y datos singulares.** Lista curada de 19 datos numéricos y nombres propios (tipos impositivos, plazos, topónimos, importes) comprobados uno a uno contra la página citada.
- **D. Cobertura.** Comprobación de que las 115 propuestas numeradas del documento aparecen individualmente en el resumen y de que ningún número falta.

## 2. Resultados

| Comprobación | Elementos | Fallos detectados | Fallos tras corrección |
|---|---|---|---|
| A. Citas literales | 79 | 2 | 0 |
| B. Atribución de página | 117 líneas con cita | 0 | 0 |
| C. Cifras y nombres propios | 19 | 0 | 0 |
| D. Cobertura de propuestas | 115 | 0 | 0 |
| Estructura (13 apartados del índice) | 13 | 0 | 0 |
| Recuento de páginas del PDF | 66 | 0 | 0 |

Detalle de la comprobación C (página donde se localizó cada dato):

| Dato del resumen | Localizado en | Estado |
|---|---|---|
| Cuota de autónomos «desde los 50 euros» | p. 9 | OK |
| «5 millones de personas en 5 años» | p. 15 | OK |
| Desalojo «en menos de 24 horas» | p. 24 | OK |
| Bonificación del «100% del ITP» | p. 24 | OK |
| IVA «0%» primera vivienda | p. 24 | OK |
| Rechazo al acuerdo UE-Mercosur | p. 33 | OK |
| Rechazo a la reducción del presupuesto de la PAC | p. 34 | OK |
| Restablecimiento del OCON Sur | p. 46 | OK |
| Gibraltar: presión internacional | p. 48 | OK |
| «doce vocales» del CGPJ elegidos por jueces | p. 55 | OK |
| Agenda 2030 (blindaje del sector primario) | pp. 12, 33 | OK |
| Interconexión de «todas las cuencas» | p. 61 | OK |
| Tarjeta sanitaria única con historia clínica compartida | p. 37 | OK |
| Fin de la cita previa | p. 27 | OK |
| Derogar «tres normas» por cada nueva | p. 28 | OK |
| Prohibición del burka y el niqab | p. 42 | OK |
| Retirada del delito de sedición (recuperarlo) | p. 55 | OK |
| Prisión permanente revisable ampliada | p. 45 | OK |
| Educación infantil 0-3 gratuita | p. 51 | OK |

## 3. Correcciones aplicadas durante la verificación

1. **Cita no literal corregida.** El resumen decía «quien entre de forma ilegal jamás podrá regularizar su situación». El original dice «Quien entra de forma ilegal jamás podrá regularizar su situación» (p. 15). Corregido a la forma textual.
2. **Titular convertido indebidamente en cita.** Se citaba «limpiarla de comisarios políticos», que era una paráfrasis. El titular original es «Limpiar RTVE de comisarios políticos» (p. 58). Reescrita la medida 110 con el texto literal.
3. **Citas de encabezado añadidas** a las medidas 54, 56 y 68 para que el resumen conserve la formulación original y no solo la descripción («Impedir la destrucción de la riqueza agrícola y el patrimonio natural», «Fin a las trabas interminables que arruinan nuestro campo y nuestro mar», «Ningún español debe quedarse sin poder estudiar por falta de recursos económicos»).
4. **Inventario de páginas en blanco corregido.** El resumen enumeraba mal las páginas sin texto. Verificado: sin texto extraíble 2, 4, 6, 8, 14, 18, 50, 64, 65, 66 y 67; con solo cabecera y número 17, 21, 22, 26, 31, 32, 36, 39, 40, 44, 49, 54, 59, 60 y 63.

## 4. Cotejo manual por muestreo

Se releyeron íntegras en el PDF, palabra por palabra, 12 propuestas consideradas de mayor riesgo interpretativo: 1 (Ceuta), 8 (energía y nuclear), 13 (remigración), 18 (prioridad nacional), 31 (okupación), 32 (alquiler a inmigrantes irregulares), 41 (excepciones al silencio positivo), 47 (leyes de género), 62 (colaboración público-privada en sanidad), 93 (Gibraltar), 100 (CGPJ), 107 (responsabilidad patrimonial de cargos públicos). Las doce coinciden con el resumen; no se detectó ninguna afirmación ausente en el original.

## 5. Limitaciones

- La verificación es **de literalidad y localización**, no de interpretación política: comprueba que el resumen dice lo que el documento dice y en la página indicada; no juzga si la medida es viable, legal o cierta.
- El verificador es el mismo proceso que redactó el resumen. Reduce errores de transcripción y atribución, pero no sustituye una revisión humana independiente.
- El documento usa maquetación a dos columnas con partición de palabras; la normalización resuelve los casos detectados, pero una cita que cruzara un salto de página con número impreso intercalado exigiría lectura manual. No se detectó ninguno abierto en este resumen.

## 6. Reproducción

```bash
cd /Users/jparrill/Projects/manifesto-atlas/runs/29n2026
pdftotext -layout "/Users/jparrill/Documents/programas-electorales-opencode-2026-10-30/29N-2026/Programa-completo-A5-1.pdf" text/vox.txt
python3 /tmp/pagemark.py text/vox.txt vox.md   # inserta <!-- page N -->
python3 /tmp/verify.py                                  # comprobaciones A y B
```

Los scripts `/tmp/pagemark.py` y `/tmp/verify.py` se conservan en este directorio como `tools/pagemark.py` y `tools/verify.py`.
