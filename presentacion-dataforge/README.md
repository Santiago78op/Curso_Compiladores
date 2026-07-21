# DataForge — Curso guiado (presentación)

Presentación por etapas del proyecto **DataForge** (OLC1, USAC): teoría de compiladores explicada desde cero mientras se construye el intérprete paso a paso.

- **Ver localmente:** abrir `index.html` con doble clic (no necesita servidor ni internet).
- **Navegación:** el índice enlaza a cada etapa; dentro de una etapa, flechas ← → o botones. Botón ◐ para tema claro/oscuro.
- **Publicar online:** activar GitHub Pages (Settings → Pages → Deploy from branch → `main` / root).

## Estructura modular

```
presentacion-dataforge/
├── index.html          ← índice del curso (roadmap con links)
├── etapa0.html         ← una página por etapa, solo sus diapositivas
├── etapa1.html
├── etapa2.html
├── etapa3.html
├── etapa4.html
├── etapa5.html
├── etapa6.html
├── gramaticas.html     ← profundización: gramáticas, LL/LR, tabla M, conflictos (demos animadas)
├── automatas.html      ← profundización: ER → AFN → AFD con simulaciones en vivo
├── tabla-simbolos.html ← profundización: tabla de símbolos y ámbitos encadenados (§2.7)
├── ast.html            ← profundización: traza real del intérprete + AST vs evaluación directa
├── fases.html          ← profundización: las fases del compilador (Cap. 1) mapeadas al proyecto
├── semantica.html      ← profundización: sistema de tipos formalizado + manejo de errores (Cap. 6)
├── assets/
│   ├── estilo.css      ← sistema de diseño compartido
│   └── deck.js         ← navegación de slides + tema claro/oscuro
└── README.md
```

**Agregar una etapa** = crear `etapaN.html` (copiar el esqueleto de una existente y reemplazar las `<section class="slide">`) y agregar su fila al roadmap de `index.html`. El CSS y JS no se tocan.

## Contenido

| Etapa | Tema | Estado |
|---|---|---|
| 0 | Conceptos base: intérprete, tokens, lexema, patrón | ✅ |
| 1 | Tabla de tokens + analizador léxico (JFlex) | ✅ verificado: 73 tokens de ejemplo1.df |
| 2 | Gramática BNF + analizador sintáctico (CUP) | ✅ verificado: ejemplos válidos [OK], error sintáctico con línea/columna |
| 3 | Ejecución: variables, arreglos, aritmética, estadísticas | ✅ verificado: 3 ejemplos, errores semánticos con línea/columna |
| 4 | Editor JavaFX | ✅ verificado: editor funcionando |
| 5 | Gráficas | ✅ verificado: las 4 ventanas abren |
| 6 | Reportes HTML + modo pánico | ✅ verificado: 3 reportes con los 3 tipos de error |

**Proyecto funcional completo** — pendiente solo el empaquetado de entrega (manuales PDF + repo GitHub).

El código de referencia del proyecto vive en `../DataForge/` (workspace Hades).
