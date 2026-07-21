# CompScript — Curso guiado (presentación)

Presentación por etapas del proyecto **CompScript** (OLC1, USAC, Vacaciones Dic. 2024): un intérprete completo con análisis léxico, sintáctico y semántico que, a diferencia de su proyecto hermano DataForge, construye explícitamente un **árbol de sintaxis abstracta (AST)** y ejecuta el programa recorriéndolo.

- **Ver localmente:** abrir `index.html` con doble clic (no necesita servidor ni internet).
- **Navegación:** el índice enlaza a cada etapa; dentro de una etapa, flechas ← → o botones. Botón ◐ para tema claro/oscuro.
- **Publicar online:** activar GitHub Pages (Settings → Pages → Deploy from branch → `main` / root).

## Estructura modular

```
presentacion-compscript/
├── index.html          ← índice del curso (roadmap con links)
├── etapa0.html         ← una página por etapa, solo sus diapositivas
├── etapa1.html
├── etapa2.html
├── etapa3.html
├── etapa4.html
├── etapa5.html
├── etapa6.html
├── etapa7.html
├── ast.html            ← profundización ★ obligatoria: el AST real de CompScript
├── assets/
│   ├── estilo.css      ← sistema de diseño compartido (idéntico al de presentacion-dataforge)
│   └── deck.js         ← navegación de slides + tema claro/oscuro (idéntico al de presentacion-dataforge)
└── README.md
```

**Agregar una etapa** = crear `etapaN.html` (copiar el esqueleto de una existente y reemplazar las `<section class="slide">`) y agregar su fila al roadmap de `index.html`. El CSS y JS no se tocan.

## Contenido

| Etapa | Tema | Estado |
|---|---|---|
| 0 | Conceptos base: por qué CompScript necesita un AST (a diferencia de DataForge) | ✅ |
| 1 | Tabla de tokens + analizador léxico (JFlex) — 37 reservadas + 28 símbolos | ✅ verificado: 209 tokens de `ejemplo1.cs` |
| 2 | Gramática BNF + analizador sintáctico (CUP) — expresiones infijas, precedencia, `UMENOS` | ✅ verificado: sin conflictos shift-reduce, modo pánico con `ejemplo6_errores.cs` |
| 3 | El árbol de sintaxis abstracta: `ast/A.java` (~25 clases de nodo) | ✅ verificado: incluye el bug real de `validarVector` corregido en la auditoría 2026-07-21 |
| 4 | Tabla de símbolos: `Contexto` y `Entorno`, ámbitos anidados, alcance estático | ✅ verificado: 186 entradas de `n` en la recursión de `fibonacci` |
| 5 | Flujo de control: `if`/`match`/`while`/`for`/`do-while`, `break`/`continue`/`return` como excepciones | ✅ verificado con `ejemplo3.cs`; incluye el bug real de falta de cortocircuito en `&&`/`\|\|` corregido en la auditoría 2026-07-21 |
| 6 | Funciones: argumentos por nombre, valores por defecto, recursión, 3 pasadas de ejecución | ✅ verificado: `factorial(5)`→120, `fib(10)`→55 |
| 7 | Editor JavaFX + los 5 reportes (tokens, errores, símbolos, AST en HTML y Graphviz) | ✅ código completo (no se abrió la ventana en la sesión de verificación) |

**Proyecto funcional completo** — pendiente solo el empaquetado de entrega (manuales PDF + repo GitHub).

### Profundización

| Página | Tema |
|---|---|
| `ast.html` ★ | El AST real de un ejemplo completo (`factorial` de `ejemplo5.cs`): demo paso a paso de cómo el parser construye los nodos, el árbol completo dibujado, y la comparación honesta evaluación directa (DataForge) vs. recorrido de árbol (CompScript) |

El código de referencia del proyecto vive en `../CompScript/` (workspace Hades); su Manual Técnico en `../CompScript/docs/ManualTecnico.md` y la gramática entregable en `../CompScript/docs/gramatica.txt`.
