# CompInterpreter — Curso guiado (presentación)

Presentación por etapas del proyecto **CompInterpreter** (OLC1, USAC, Proyecto 2): un intérprete web cliente-servidor para el lenguaje `.ci`, construido con Jison, Express y React, explicado en narración de instructor a partir del código real del proyecto.

- **Ver localmente:** abrir `index.html` con doble clic (no necesita servidor ni internet).
- **Navegación:** el índice enlaza a cada etapa; dentro de una etapa, flechas ← → o botones. Botón ◐ para tema claro/oscuro.
- **Publicar online:** activar GitHub Pages (Settings → Pages → Deploy from branch → `main` / root).

## Estructura modular

```
presentacion-compinterpreter/
├── index.html              ← índice del curso (roadmap con links)
├── etapa0.html              ← una página por etapa, solo sus diapositivas
├── etapa1.html
├── etapa2.html
├── etapa3.html
├── etapa4.html
├── etapa5.html
├── etapa6.html
├── etapa7.html
├── arquitectura-rest.html   ← profundización ★: arquitectura cliente-servidor REST (demo interactiva)
├── assets/
│   ├── estilo.css           ← sistema de diseño compartido (idéntico al de presentacion-dataforge)
│   └── deck.js               ← navegación de slides + tema claro/oscuro (idéntico, sin modificar)
└── README.md
```

**Agregar una etapa** = crear `etapaN.html` (copiar el esqueleto de una existente y reemplazar las `<section class="slide">`) y agregar su fila al roadmap de `index.html`. El CSS y el JS compartidos no se tocan.

## Contenido

| Etapa | Tema | Estado |
|---|---|---|
| 0 | Conceptos base: por qué JS/Jison en vez de Java/JFlex+CUP, y por qué esta vez sí hace falta un AST | ✅ |
| 1 | La gramática Jison: tabla de tokens, `%lex`, escapes y recuperación léxica | ✅ verificado: `ejemplo_errores.ci` produce el error léxico esperado |
| 2 | El parser: precedencia (tabla 5.10), la fábrica de nodos del AST, recuperación sintáctica | ✅ verificado: `ejemplo_vectores2d.ci` → 106 nodos / 105 aristas |
| 3 | El intérprete: entorno, tipos, coerción vs. cast, dos pasadas (forward-reference), propagación por null | ✅ |
| 4 | Señales de control (`break`/`continue`/`return`), *fall-through* del switch, y los 2 bugs reales corregidos en la auditoría del 21/07/2026 | ✅ |
| 5 | El servidor Express: contrato REST completo, `analizar.js` como orquestador, entorno fresco por petición, CORS | ✅ |
| 6 | El cliente React: `forwardRef`+`useImperativeHandle` en el editor, gutter sincronizado, AST con `vis-network` | ✅ |
| 7 | Suite Playwright: 3 pruebas de extremo a extremo (camino feliz, los 3 tipos de error, nuevo archivo) | ✅ verificado: 3/3 pasan |

**Profundización ★:** arquitectura cliente-servidor REST — el contrato de la API paso a paso, con una demo interactiva del ciclo completo petición/respuesta, y las dos consecuencias de diseño que esa arquitectura obliga (entorno fresco por concurrencia, CORS) que no existen en los otros 3 proyectos del curso (monolitos JavaFX).

**Proyecto funcional completo** (servidor y cliente) — pendiente solo el empaquetado de entrega (manuales PDF + repo GitLab `OLC1_Proyecto2_#Carnet`).

El código de referencia del proyecto vive en `../CompInterpreter/` (workspace Hades): gramática real en `server/src/grammar.jison`, intérprete en `server/src/interprete/`, cliente en `client/src/`, casos de prueba en `entradas/*.ci`, tests E2E en `client/e2e/compinterpreter.spec.js`.
