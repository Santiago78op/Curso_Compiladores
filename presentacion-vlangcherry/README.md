# VLangCherry — Curso guiado (presentación)

Presentación por etapas del proyecto **VLangCherry** (OLC2, USAC): material de estudio para el equipo sobre la arquitectura real del intérprete del lenguaje V-Lang Cherry, construido con ANTLR4 sobre Go.

- **Ver localmente:** abrir `index.html` con doble clic (no necesita servidor ni internet).
- **Navegación:** el índice enlaza a cada etapa; dentro de una etapa, flechas ← → o botones. Botón ◐ para tema claro/oscuro.
- **Publicar online:** activar GitHub Pages (Settings → Pages → Deploy from branch → `main` / root).

## Estructura modular

```
presentacion-vlangcherry/
├── index.html          ← índice del curso (roadmap con links)
├── etapa0.html         ← conceptos base: por qué ANTLR4 + Go
├── etapa1.html         ← la gramática ANTLR4 (lexer + parser en un .g4)
├── etapa2.html         ← del parse tree al AST propio, sin Visitor generado
├── etapa3.html         ← tipos estáticos, structs y métodos por referencia
├── etapa4.html         ← slices 1D/2D
├── etapa5.html         ← control de flujo, recursión y auditoría de validación semántica
├── etapa6.html         ← servidor REST con net/http
├── etapa7.html         ← cliente React, reusando componentes
├── generadores.html    ← profundización ★: ANTLR4 vs JFlex+CUP vs Jison
├── assets/
│   ├── estilo.css      ← sistema de diseño compartido (copiado tal cual de presentacion-dataforge)
│   └── deck.js         ← navegación de slides + tema claro/oscuro (idem)
└── README.md
```

**Agregar una etapa** = crear `etapaN.html` (copiar el esqueleto de una existente y reemplazar las `<section class="slide">`) y agregar su fila al roadmap de `index.html`. El CSS y JS no se tocan.

## Contenido

| Etapa | Tema | Estado |
|---|---|---|
| 0 | Conceptos base: por qué ANTLR4 + Go, y no JFlex/CUP + Java | ✅ |
| 1 | La gramática ANTLR4: lexer y parser combinados en un `.g4` | ✅ |
| 2 | Del parse tree al AST propio: por qué no se usó el Visitor generado | ✅ |
| 3 | Tipos estáticos, structs y métodos "por referencia" gratis vía punteros Go | ✅ |
| 4 | Slices 1D/2D | ✅ |
| 5 | Control de flujo, recursión y los 4 arreglos de la auditoría (2026-07-21) | ✅ |
| 6 | Servidor REST con `net/http` | ✅ |
| 7 | Cliente React, reusando componentes de CompInterpreter | ✅ |
| ★ | Profundización: ANTLR4 vs JFlex+CUP vs Jison, los 3 generadores del curso | ✅ |

**Proyecto de referencia completo y verificado** (gramática, intérprete, servidor y cliente funcionando; 6 ejemplos + suite Playwright en verde). Esta presentación es material de estudio para la defensa oral del equipo — no reemplaza `docs/ManualTecnico.md` ni `docs/gramatica.txt`, que son la fuente de verdad técnica.

El código de referencia del proyecto vive en `../VLangCherry/` (workspace Hades).

## Sobre este proyecto

VLangCherry es el primer proyecto de **OLC2** (curso distinto de OLC1, donde viven DataForge/ConjAnalyzer/CompScript/CompInterpreter) y es trabajo **grupal de 3 integrantes**, con defensa oral en vivo. Es también el único de los 5 proyectos del curso guiado que usa ANTLR4 en vez de JFlex+CUP o Jison, y Go en vez de Java o JavaScript — decisión que la Etapa 0 explica en detalle.
