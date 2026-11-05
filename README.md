# Manifesto Atlas

## Vista previa local en macOS

Ya puedes abrir una ejecución real de Frente Obrero: 64 páginas procesadas,
21 temas y 17 medidas principales con citas. El diseño reutiliza el HTML de
referencia aportado para el proyecto, sin modificar sus estilos ni scripts.

```bash
make render-preview
open runs/preview/frente-e2e/index.html
```

Es un **borrador privado**, no una publicación aprobada. Se conservan las llamadas
reales de generación y validación independiente; falta evidencia visual para las
páginas sin texto extraíble y el informe de cobertura señala posibles omisiones.
La pestaña de VOX indica que aún está pendiente de
este flujo. El JSON canónico, los informes y el Markdown están junto al HTML.

Para repetir el proceso con modelos y un nuevo directorio de resultados, consulta
el [flujo de macOS](docs/design/macos-preview.md). La regeneración ordinaria de
la vista almacenada funciona sin conexión ni llamadas a modelos.

> Encontrar información fiable, primaria y no editorializada sobre las promesas electorales es innecesariamente doloroso.

Manifesto Atlas es un archivo abierto y versionado de programas electorales y documentos políticos. Transforma las fuentes oficiales publicadas por los partidos en resúmenes breves y verificables, y las conserva en Git como memoria histórica.

El proyecto no resume artículos de prensa, comentarios de campaña, publicaciones en redes sociales ni interpretaciones de terceros. Los resúmenes parten de los documentos publicados por los propios partidos y permanecen vinculados a sus fuentes originales.

## TL;DR

Los programas políticos suelen estar dispersos entre webs de partidos, páginas temporales de campaña, portales regionales y documentos PDF difíciles de consultar. Muchos enlaces desaparecen después de las elecciones.

Manifesto Atlas pretende que esos documentos sean:

- Fáciles de encontrar.
- Fáciles de entender.
- Comparables entre partidos y territorios.
- Trazables hasta la fuente y la página originales.
- Conservados como registro histórico versionado.
- Consultables sin depender de comentarios políticos externos.

## Principios del proyecto

### Prioridad para las fuentes primarias

Cada resumen debe basarse en un programa electoral, programa político o documento equivalente publicado oficialmente por el partido u organización electoral.

Cada programa registra su fichero original, ámbito geográfico, elección, fecha de publicación y huella digital. Las medidas y demás afirmaciones factuales generadas incluyen referencias a las páginas de origen.

La primera versión trabaja con documentos oficiales seleccionados por el mantenedor y conserva la evidencia de su procedencia. No ofrece admisión automática de cualquier URL o fichero aportado por terceros. Cada afirmación factual generada debe tener citas verificables; un origen seleccionado no sustituye la validación independiente del contenido.

### Resumen factual y análisis separados

La capa de resumen describe lo que dice el documento. No juzga si una propuesta es deseable, realista o políticamente conveniente.

Las estimaciones y evaluaciones se muestran de forma separada y nunca se presentan como parte del contenido original del programa.

### El contexto territorial importa

Un programa nacional y otro autonómico del mismo partido no son intercambiables. Manifesto Atlas distingue entre documentos de ámbito:

- Europeo.
- Nacional.
- Autonómico.
- Provincial.
- Municipal.
- Político general, cuando no se trate de un programa electoral.

Las comparaciones de medidas comunes y diferencias regionales quedan para una ampliación posterior, cuando existan programas compatibles ya validados. La primera versión conserva sus identidades territoriales sin mezclarlas.

### Git como memoria histórica

Las fuentes, metadatos, resúmenes, metodologías y vistas generadas se versionan en Git. Esto permite comprobar:

- Qué publicó cada partido.
- Cuándo se incorporó la fuente.
- Cómo cambió un resumen.
- Qué versión de la metodología produjo una evaluación.
- Si una corrección posterior modificó la interpretación.

El objetivo es conservar la evidencia cuando desaparezcan las páginas de campaña o los enlaces de descarga.

## Flujo de publicación

El flujo previsto es:

1. Recuperar un PDF oficial concreto y conservar sus bytes, metadatos, procedencia y checksum.
2. Inspeccionarlo y extraer el texto por páginas con límites de recursos y aislamiento.
3. Utilizar un modelo configurado para generar un JSON de temas, medidas y citas en español.
4. Comprobar las afirmaciones y la cobertura completa mediante dos pases independientes contra la evidencia original.
5. Generar un resumen Markdown legible y verificar que todos sus resultados corresponden a las entradas actuales.
6. Repetir el proceso con el segundo documento real; añadir después HTML estático con assets suministrados y comprobaciones básicas de accesibilidad.

El objetivo inmediato es demostrar este recorrido con **un documento real**, antes de construir una plataforma de contribuciones, firmas, comparaciones o certificación de sitios web. La publicación en GitHub Pages conserva su objetivo, pero su despliegue se diseñará por separado.

El LLM que genera el resumen no es una autoridad ni puede validarse a sí mismo. Antes de publicar, otros pases con contexto limpio comprueban cada afirmación y referencia contra las páginas originales, y buscan omisiones o contradicciones en el documento completo. Las validaciones, los modelos, los prompts y las versiones de las herramientas quedan registrados en Git.

La validación es automática y probabilística: reduce errores, pero no garantiza que la interpretación de los modelos sea infalible. Si las comprobaciones no coinciden o la evidencia es insuficiente, el programa permanece sin publicar.

## Ampliaciones futuras: medición del cumplimiento

Estas evaluaciones no forman parte de la primera versión ni son un requisito para publicar un resumen factual. Su metodología necesitará una aprobación específica antes de implementarse.

Manifesto Atlas distingue tres conceptos que no deben mezclarse en un único dato sin contexto.

### Estimación de viabilidad

Estimación prospectiva sobre si una propuesta parece realizable durante la legislatura correspondiente. Puede considerar:

- Si la institución dispone de las competencias legales necesarias.
- Si la propuesta identifica financiación o recursos.
- Si el plazo de ejecución cabe dentro de la legislatura.
- Si la medida es suficientemente concreta para comprobarla.
- Si depende de otro gobierno, institución o acontecimiento externo.

Es una estimación, no un porcentaje de cumplimiento observado.

### Fiabilidad histórica

Historial del partido al cumplir promesas comparables. Debe calcularse usando programas anteriores versionados y resultados verificables, no su reputación general ni comentarios de terceros.

La evidencia reciente y directamente comparable debe tener más peso que promesas antiguas o no relacionadas.

### Cumplimiento observado

Evaluación retrospectiva realizada durante o después de la legislatura. Cada promesa podrá clasificarse como:

- Cumplida.
- Parcialmente cumplida.
- En curso.
- Incumplida.
- Bloqueada por otra institución.
- No verificable.

La evidencia y la fecha de evaluación deben quedar registradas para permitir revisiones posteriores.

### Confianza de cumplimiento

Podrá existir un indicador compuesto que combine la estimación de viabilidad con la fiabilidad histórica. Su fórmula, pesos, requisitos de evidencia e incertidumbre serán públicos y estarán versionados.

La interfaz mostrará siempre el desglose. Un porcentaje aislado sin evidencia produciría una falsa precisión contraria a los objetivos del proyecto.

## Lo que no es este proyecto

Manifesto Atlas no es:

- Una recomendación de voto.
- Una clasificación ideológica de partidos.
- Un sustituto de la lectura del programa original.
- Una verificación automática de todas las estadísticas citadas por un partido.
- Una afirmación de que la interpretación de un LLM es siempre correcta.

Es un archivo de fuentes, un sistema de resumen, una herramienta de comparación territorial y un registro verificable de compromisos políticos.

## Política de idiomas

El contenido dirigido al público español se escribe en español:

- Este README y la documentación introductoria.
- Los resúmenes de programas.
- La metodología explicada al público.
- El sitio HTML y todos sus textos visibles.

Los elementos técnicos se escriben en inglés:

- Source code and code comments.
- Make targets and command-line interfaces.
- Schemas, internal metadata fields, and prompts.
- Technical plans and architecture documents.
- Tests, commit messages, and developer documentation.

Los identificadores internos pueden estar en inglés, pero la interfaz debe traducirlos al español.

## Licencia

El código y la documentación originales del proyecto se publican bajo la [licencia MIT](LICENSE).

Los documentos de terceros, incluidos programas electorales y documentos políticos, mantienen los derechos y condiciones establecidos por sus editores originales. Su incorporación al archivo no los somete a la licencia MIT.

## Estado

Existen bibliotecas y pruebas de base, pero todavía no hay un recorrido funcional completo desde un PDF hasta un resumen público. La [especificación técnica](docs/SPEC.md) y el [contrato del MVP](docs/design/MVP.md) definen el plan reducido; actualizar tareas y contratos no significa que esas funciones ya estén implementadas.
