# ConjAnalyzer — Curso guiado (presentación)

Presentación por etapas del proyecto **ConjAnalyzer** (OLC1, USAC): un intérprete de teoría de conjuntos — definición de conjuntos, operaciones en notación prefija, simplificación algebraica, diagrama de Venn y reportes — explicado etapa por etapa con el código real del proyecto.

- **Ver localmente:** abrir `index.html` con doble clic (no necesita servidor ni internet).
- **Navegación:** el índice enlaza a cada etapa; dentro de una etapa, flechas ← → o botones. Botón ◐ para tema claro/oscuro.
- **Publicar online:** activar GitHub Pages (Settings → Pages → Deploy from branch → `main` / root).

## Estructura modular

```
presentacion-conjanalyzer/
├── index.html          ← índice del curso (roadmap con links)
├── etapa0.html          ← conceptos base: los conjuntos como lenguaje
├── etapa1.html          ← tabla de tokens + analizador léxico (JFlex)
├── etapa2.html          ← gramática BNF + analizador sintáctico (CUP)
├── etapa3.html          ← ejecución: Entorno, Conjunto, el árbol NodoOperacion
├── etapa4.html          ← el Simplificador: las 5 leyes de la sección 7
├── etapa5.html          ← diagrama de Venn exacto, con Canvas
├── etapa6.html          ← reportes HTML + JSON de simplificación (Gson)
├── simplificador.html  ← profundización ★: el Simplificador en movimiento (demos animadas)
├── assets/
│   ├── estilo.css      ← sistema de diseño compartido (idéntico al de presentacion-dataforge)
│   └── deck.js         ← navegación de slides + tema claro/oscuro (idéntico, sin modificar)
└── README.md
```

**Agregar una etapa** = crear `etapaN.html` (copiar el esqueleto de una existente y reemplazar las `<section class="slide">`) y agregar su fila al roadmap de `index.html`. El CSS y JS no se tocan.

## Contenido

| Etapa | Tema | Estado |
|---|---|---|
| 0 | Conceptos base: universo ASCII 33–126, case sensitive real, notación prefija | ✅ |
| 1 | Tabla de tokens + analizador léxico (JFlex) — sin `%ignorecase` | ✅ |
| 2 | Gramática BNF + analizador sintáctico (CUP) — no terminal `operacion` tipado `NodoOperacion` | ✅ |
| 3 | Ejecución: `Entorno`, `Conjunto`, el árbol acotado `NodoOperacion` y por qué sí hace falta | ✅ |
| 4 | El Simplificador: 5 leyes (doble complemento, DeMorgan, idempotentes, absorción, distributivas) | ✅ auditado 2026-07-21 |
| 5 | Diagrama de Venn exacto (sombreado píxel a píxel vía `pertenenciaRegion`) | ✅ |
| 6 | Reportes HTML (reflexión sobre `sym`) + JSON de simplificación con Gson | ✅ |

**Proyecto funcional completo** — pasó una auditoría de calidad (21 de julio de 2026) que agregó la ley distributiva faltante en `Simplificador.java` y eliminó código muerto (`Entorno.getUniverso()`). Pendiente solo el empaquetado de entrega (JAR ejecutable, manuales PDF, repo `OLC1_Proyecto1_#Carnet`).

El código de referencia del proyecto vive en `../ConjAnalyzer/` (workspace Hades).
