# Fase A — Resumen consolidado de la auditoría del Libro del Dragón

**Qué se hizo:** lectura completa de los caps. 1–7 del Dragón (2ª ed., pp. 1–503) desde `doc\Boocks\dragon-md\` (conversión particionada y limpiada, ver su README), contrastada contra las **33 notas** del vault que los cubren (8 de `03-Libro` + 22 de `01-Conceptos` + JFlex y CUP de `02-Tecnologias`).
**Producto:** 7 informes por capítulo (`cap0N-brechas.md`) escritos como **backlog ejecutable**: cada brecha trae cita exacta (§ + página), nota destino y acción concreta — aplicables sin releer el libro.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Diagnóstico global del vault

- **Fortalezas:** citas casi todas correctas; los puentes teoría→proyecto son excelentes donde existen (frame≈entorno del cap. 7, "en un intérprete no se generan etiquetas" del cap. 6, el encuadre "no requerido pero cae en examen" de tres direcciones); las 3 notas de algoritmos del cap. 3 comparten el ejemplo canónico `(a|b)*abb`; `Minimización de AFD` es la nota modelo.
- **Patrón de brecha #1 — "resultado sin algoritmo":** varias notas dan la tabla final pero no el procedimiento ejecutable a mano que piden los exámenes (FIRST/FOLLOW, construcción de subconjuntos, tabla LL(1), motor LR).
- **Patrón de brecha #2 — conceptos que los proyectos implementan sin su nombre teórico:** sobrecarga, coerción max/ampliar, token `error`, enlace de acceso vs control, S-atribuida en la pila. Son las respuestas de defensa.
- **Desproporción:** el cap. 4 (corazón de CUP/Jison, 112 pp.) tiene la cobertura más delgada relativa a su peso.

## Las 12 brechas ALTAS (prioridad de la Fase D)

| # | Brecha | Cap./§ | Nota destino |
|---|--------|--------|--------------|
| 1 | Entorno vs estado (nombres→ubicaciones→valores, l/r-value) | 1 §1.6.2 | `Entornos y alcance` |
| 2 | Paso de parámetros (valor/referencia/nombre) + aliasing — **sin nota** | 1 §1.6.6–7 | nueva: `Paso de parámetros` |
| 3 | Precedencia por niveles de gramática (n+1 no terminales) — clave para la BNF entregable | 2 §2.2.6 | `Ambigüedad, precedencia…` |
| 4 | Tablas de símbolos encadenadas (`Ent`/`ant`/get-en-cadena) = la clase `Entorno` | 2 §2.7 | `Tabla de símbolos` |
| 5 | Construcción de subconjuntos ejecutable (algoritmo + traza) | 3 §3.7.1 | `Construcción de subconjuntos` |
| 6 | El motor LR: ítems LR(0), CERRADURA, ir_A, tabla ACCION, traza | 4 §4.6 | nueva: `Elementos LR(0) y la tabla SLR` |
| 7 | Reglas de cálculo de FIRST y FOLLOW (las 6 reglas) | 4 §4.4.2 | `FIRST y FOLLOW` |
| 8 | El token `error` de CUP/Yacc = teoría de `instruccion ::= error PUNTO_COMA` | 4 §4.9.4 | `Manejo de errores` + `CUP` |
| 9 | Atributos EN la pila del parser (la física de `RESULT`/`$1`) | 5 §5.4.2 | `Traducción dirigida por la sintaxis` |
| 10 | Sobrecarga de operadores (resolución por firma) = tablas de compatibilidad | 6 §6.5.3 | `Comprobación de tipos` |
| 11 | Coerción binaria: `max(t₁,t₂)` + `ampliar` = dominancia de tipos de `Operaciones` | 6 §6.5.2 | `Conversión de tipos` |
| 12 | Enlace de acceso vs control + bug del alcance dinámico accidental | 7 §7.3.5–6 | `Registro de activación…` |

**Medias destacadas:** grafo de dependencias (5), autómata→lexer con longest match (3), SLR vs LALR con el ejemplo L=R (4), gramática ambigua+precedence como decisión de ingeniería (4 §4.8.1), árbol de activación (7), break/continue dual compilador/intérprete (6 §6.7.4), nota nueva de GC con conexión JVM/V8 (7), recorridos preorden/postorden (2), "wcw no es libre de contexto" como respuesta de defensa (4 §4.3.5).

**Notas nuevas propuestas (6):** Paso de parámetros · Recorridos de árboles · Del autómata al analizador léxico · Elementos LR(0) y la tabla SLR · ER→AFD directo (siguientepos, opcional según pensum) · Recolección de basura. **Notas a ampliar/corregir: ~15.**

**Correcciones puntuales de citas/contenido:** `Entornos y alcance` (agregar §1.6.3–1.6.5), `Manejo de errores` (separar pánico léxico §3.1.4 del sintáctico §4.1.4), `Ambigüedad…` (agregar §2.2.5–2.2.6).

## Hallazgos pro-defensa (citas que VALIDAN los proyectos)

1. **DataForge sin AST**: §2.8.1 p. 92 ("es común… sin construir en realidad la estructura de datos tipo árbol") + §5.4.1–5.4.2 (SDT postfijo S-atribuido implementado en la pila LR).
2. **Gramática ambigua + `precedence` en los `.cup`/`.jison`**: §4.8.1 pp. 278–281 (menos estados, sin reducciones inútiles E→T, T→F).
3. **Validaciones en semántica y no en la gramática**: §4.3.5 (declarar-antes-de-usar = `wcw`, no libre de contexto).
4. **`instruccion ::= error PUNTO_COMA`**: §4.1.4 (producciones de error) + §4.9.4 (mecánica del token `error`).
5. **Recursión con pila de entornos**: §7.2 (árbol de activación) — cada llamada su frame ≈ su entorno.

## Insumos para la Fase B (auditoría de presentaciones)

Verificaciones marcadas en los informes: ¿`fases.html` anima el ejemplo guía de la fig. 1.7? ¿`automatas.html` anima Thompson→subconjuntos→minimización con `(a|b)*abb`? ¿`gramaticas.html` muestra el motor LR (autómata + tabla + traza)? Steppers candidatos nuevos: trazas gemelas LL/LR (figs. 4.21/4.38), longest-match con `aaba`, backpatching 100–105, pila con atributos (`RESULT`), árbol de activación de quicksort + pila, ejemplo 2.14 de shadowing, bug del alcance dinámico.

## Cómo ejecutar la Fase D (con Opus 4.8 + agentes)

1. Cambiar el modelo a Opus 4.8 y lanzar **un agente por capítulo** (7 agentes en paralelo), cada uno con su `capNN-brechas.md` como spec; escriben con `brain_upsert_note` (crear) o Edit (ampliar), citando siempre `Dragón §x.y, p. NN`.
2. Convenciones del vault: frontmatter con tags/aliases/fuente/fecha, wikilinks liberales, Mermaid, **sin `\|` con alias dentro de tablas**.
3. Verificación: agente `fable` (o la sesión principal) contrasta cada nota nueva/ampliada contra `dragon-md\capNN-*.md` antes de aceptar.
4. Cierre: `brain_reindex` → `brain_health` (0 rotos, 0 huérfanas) → actualizar el MOC si hay notas nuevas → commit.
