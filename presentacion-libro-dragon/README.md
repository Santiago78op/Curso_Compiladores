# El Libro del Dragón — capítulo a capítulo (presentación)

Presentación por capítulos del **Libro del Dragón** (*Compilers: Principles, Techniques & Tools*, Aho · Lam · Sethi · Ullman, 2ª ed.) para OLC1 (USAC): la teoría del libro contada en el mismo formato del curso guiado de DataForge, conectando cada capítulo con los 4 proyectos del curso.

- **Ver localmente:** abrir `index.html` con doble clic (no necesita servidor ni internet).
- **Navegación:** el índice enlaza a cada capítulo; dentro de uno, flechas ← → o botones. Botón ◐ para tema claro/oscuro.
- **Publicar online:** activar GitHub Pages (Settings → Pages → Deploy from branch → `main` / root).

## Estructura modular

```
presentacion-libro-dragon/
├── index.html          ← índice (roadmap de capítulos)
├── cap1.html           ← una página por capítulo, solo sus diapositivas
├── ...                 ← cap2–cap7 + panorama (caps. 8–12) — se agregan por etapas
├── assets/
│   ├── estilo.css      ← sistema de diseño (derivado del de presentacion-dataforge + capa de motion propia)
│   └── deck.js         ← navegación de slides + tema claro/oscuro + dirección de la cascada
└── README.md
```

**Motion de la plantilla** (exclusivo de esta presentación): las diapositivas entran con cascada de bloques consciente de la dirección (avanzar ↑ / retroceder ↓), el índice tiene secuencia de carga escalonada, tarjetas y filas del roadmap tienen micro-interacciones al hover, y las respuestas del quiz se revelan animadas. Todo se desactiva automáticamente con `prefers-reduced-motion`.

**Agregar un capítulo** = crear `capN.html` (copiar el esqueleto de uno existente y reemplazar las `<section class="slide">`) y actualizar su fila en el roadmap de `index.html`. El CSS y JS compartidos no se tocan; las demos animadas ("stepper") llevan su CSS/JS inline en cada página.

## Contenido

| Página | Capítulo del Dragón | Estado |
|---|---|---|
| `cap1.html` | 1 · Introducción — compilador vs. intérprete, fases, ejemplo guía (fig. 1.7), tabla de símbolos, generadores, estático vs. dinámico | ✅ |
| `cap2.html` | 2 · Un traductor sencillo — BNF, derivaciones (demo), asociatividad/precedencia, acciones semánticas y postfijo (demo), descendente vs. LALR, lexer artesanal, ámbitos, AST | ✅ |
| `cap3.html` | 3 · Análisis léxico — rol del lexer, errores léxicos, ER y definiciones regulares, demo "match más largo + primera regla", anatomía del `.flex`, pipeline ER→AFN→AFD (enlaza a las demos de `automatas.html` del curso DataForge) | ✅ |
| `cap4.html` | 4 · Análisis sintáctico — 4 estrategias de recuperación de errores, else colgante, recetas LL, FIRST/FOLLOW, mangos, demo del autómata LR(0) con ítems, escalera SLR→LALR, checklist del `.cup` (complementa `gramaticas.html` del curso DataForge) | ✅ |
| `cap5.html` | 5 · Traducción dirigida por la sintaxis — SDD, sintetizados vs. heredados, demo del árbol anotado de `3*5+4`, grafo de dependencias, S/L-atribuidas, la pila de valores de CUP (`RESULT`), marcadores, el atributo-AST (DataForge vs. CompScript) | ✅ |
| `cap6.html` | 6 · Código intermedio leído con lente de intérprete — AST/DAG, tres direcciones (cultura), tipos y declaraciones, chequeo de tipos con demo de coerción (`inttofloat`), corto circuito, backpatching (por qué no aplica), switch, checklist semántico | ✅ |
| `cap7.html` | 7 · Entornos en tiempo de ejecución — zonas de memoria, árbol de activación, registro ≈ `Entorno`, demo de la pila con `fact(3)`, alcance estático vs. dinámico (el bug clásico), señales para return/break, heap y GC (cultura) | ✅ |
| `panorama.html` | 8–12 · Panorama — la frontera front-end/back-end (con la frase para la defensa), generación de código, optimización (y por qué ConjAnalyzer NO es el cap. 9), paralelismo, Apéndice A y a dónde seguir, mapa final capítulo → pieza del proyecto | ✅ |

**Presentación completa** — caps. 1–7 + panorama, con demo animada en cada capítulo.

Presentaciones hermanas (los 4 proyectos de OLC1, mismo formato de curso guiado): [`../presentacion-dataforge/`](../presentacion-dataforge/), [`../presentacion-conjanalyzer/`](../presentacion-conjanalyzer/), [`../presentacion-compscript/`](../presentacion-compscript/), [`../presentacion-compinterpreter/`](../presentacion-compinterpreter/). Las notas fuente viven en `../Vault-Compiladores/03-Libro/`.
