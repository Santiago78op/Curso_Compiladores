# Hades — Segundo cerebro de Organización de Lenguajes y Compiladores

Material completo de dos cursos de la Facultad de Ingeniería (USAC): **OLC1** (Organización de Lenguajes y Compiladores 1) y **OLC2** (Organización de Lenguajes y Compiladores 2). Incluye el código funcional de 5 intérpretes, sus presentaciones de clase, sus manuales de programación, y un vault de teoría de compiladores.

**Punto de entrada para dar clase:** [`GuiaCurso.md`](GuiaCurso.md) — el plan completo de sesiones, intercalando teoría con cada proyecto.

## Qué hay acá

| Carpeta | Qué es |
|---|---|
| `DataForge/`, `ConjAnalyzer/`, `CompScript/`, `CompInterpreter/` | Los 4 proyectos de **OLC1**: intérpretes completos en Java (JFlex+CUP+JavaFX) salvo CompInterpreter (JS+Jison+Express+React). Cada uno con `docs/gramatica.txt`, `docs/ManualUsuario.md`, `docs/ManualTecnico.md` y `docs/GuiaProgramacion.md`. |
| `VLangCherry/` | El proyecto de **OLC2**: intérprete en Go+ANTLR4, trabajo grupal con defensa oral. Mismos 4 documentos en `docs/`. |
| `presentacion-dataforge/`, `presentacion-conjanalyzer/`, `presentacion-compscript/`, `presentacion-compinterpreter/`, `presentacion-vlangcherry/` | Una presentación HTML por proyecto (abrir `index.html` con doble clic, sin servidor), con su `GuionClase.md` — la chuleta de cómo darla. |
| `presentacion-libro-dragon/` | Presentación teórica del Libro del Dragón (Aho/Lam/Sethi/Ullman), capítulo a capítulo, conectada explícitamente con los 4 proyectos de OLC1. |
| `Vault-Compiladores/` | Vault Obsidian con la teoría de compiladores organizada (conceptos, tecnologías, resúmenes del libro, guías de proyecto), consultable también vía el MCP `diamon`. |
| `doc/` | Los enunciados originales de cada proyecto (PDF + conversión a Markdown) y el Libro del Dragón. |

## Los 3 documentos de cada proyecto, y cuándo abrir cada uno

1. **`GuionClase.md`** (en cada `presentacion-X/`) — tu chuleta para dar la clase: tiempos, qué decir, demos en vivo con comando exacto, preguntas probables.
2. **`index.html`** (en cada `presentacion-X/`) — lo que ve la audiencia: diapositivas en 3ª persona, con quiz.
3. **`docs/GuiaProgramacion.md`** (en cada carpeta de proyecto) — el tutorial paso a paso para quien va a escribir el código, organizado por tema.

Detalle completo de cómo se conectan en [`GuiaCurso.md`](GuiaCurso.md).

## Estado

Los 5 proyectos compilan y corren limpio, auditados con `codebase-memory-mcp` (bugs reales corregidos, no cosméticos). Las 6 presentaciones están en 3ª persona, enlazadas entre sí. Pendiente (a mano, fuera del alcance de este repo): capturas de pantalla reales en los manuales, exportar a PDF, y armar los repos de entrega en GitHub/GitLab.
