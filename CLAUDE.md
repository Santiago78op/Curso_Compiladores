# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Qué es este workspace

**Hades** es el workspace del curso *Organización de Lenguajes y Compiladores 1* (OLC1, USAC): un segundo cerebro de teoría de compiladores + el código de sus proyectos + el material didáctico del curso guiado.

```
Hades\
├── doc\
│   ├── Boocks\    ← Libro del Dragón (Aho/Lam/Sethi/Ullman, 2ª ed.) en PDF
│   └── Projects\  ← 4 PDFs de ENUNCIADOS: DataForge (PT1 1S2024), ConjAnalyzer (PT1 2S2024),
│                    CompScript (PT1 VD2024), CompInterpreter (PT2 2S2024)
│                    (los .clean.md junto a cada PDF son la conversión refinada — leé esos)
├── DataForge\               ← ✅ COMPLETO: intérprete Java funcional (ver sección DataForge)
├── presentacion-dataforge\  ← repo git propio: curso guiado en HTML (6 etapas + profundizaciones ★)
├── presentacion-libro-dragon\ ← repo git propio: el Dragón caps. 1-7 + panorama 8-12, COMPLETA (misma plantilla + capa de motion; complementa sin duplicar las profundizaciones ★ de arriba)
└── Vault-Compiladores\      ← vault Obsidian (57 notas), cerebro diamon "compiladores"
    ├── 00-MOC\           ← punto de entrada: enlaza TODAS las notas
    ├── 01-Conceptos\     ← teoría (léxico/sintáctico/semántico/intermedio)
    ├── 02-Tecnologias\   ← JFlex, CUP, Jison, Maven, JavaFX, Mermaid, Graphviz, vis-network
    ├── 03-Libro\         ← resúmenes del Dragón por capítulo (núcleo: caps 1-7)
    ├── 04-Proyectos\     ← ficha + "Guía <proyecto>" de cada uno (la Guía DataForge tiene el material REAL)
    └── 05-Recursos\      ← fichas de libros, Appel vs Dragón
```

## Reglas de trabajo del usuario

- **Español siempre**, tono didáctico: está aprendiendo compiladores, explicá el porqué teórico de cada paso, citando el Dragón (sección exacta) y las notas del vault.
- Trabajar **etapa por etapa**: resumir al terminar cada una y **esperar confirmación** antes de seguir.
- **Citar fuentes**. **Nunca inventar** contenido de archivos no leídos.
- **Modo de enseñanza acordado**: DataForge lo construyó Claude como demostración guiada (✅ hecho). En **ConjAnalyzer, CompScript y CompInterpreter los roles se INVIERTEN**: el usuario escribe el código, Claude tutorea y revisa.
- Estado actual: DataForge completo y verificado. Pendiente: empaquetado de entrega (manuales PDF, repo `OLC1_Proyecto1_#Carnet`) y/o arrancar ConjAnalyzer.
- El usuario trabaja los proyectos Java con **IntelliJ IDEA** (JDK lo gestiona el IDE; el `java` del PATH es 1.8 — no usarlo). Dar instrucciones en términos de IDEA; la terminal es alternativa.

## DataForge (`Hades\DataForge\`)

### Comandos

| Acción | En IDEA | Terminal |
|---|---|---|
| Regenerar lexer+parser y compilar | ventana Maven → Lifecycle → `compile` | `mvn compile` |
| Correr la GUI | ▶ Run sobre **`Lanzador`** (NUNCA sobre `EditorApp`: da "JavaFX runtime components are missing") | `mvn clean javafx:run` |
| Probar intérprete por consola | ▶ Run sobre `TestInterprete` (default: un ejemplo de `entradas/`) | `mvn compile exec:java -Dexec.args="entradas/ejemplo2.df"` |
| Tras tocar `pom.xml` | 🔄 Reload All Maven Projects | — |

Los WARNING de JavaFX al arrancar (unnamed module, native access, Unsafe) son esperados e inofensivos.

### Arquitectura (lo que no se ve en un archivo solo)

Pipeline: `Lexer.flex` (JFlex) → `parser.cup` (CUP, LALR) → **ejecución directa en las acciones** `{: RESULT = … :}` — gramática S-atribuida, **sin AST** (decisión defendida: DataForge no tiene control de flujo, toda instrucción corre exactamente 1 vez; CompScript sí necesitará AST).

- `dataforge.analisis` — **TODO GENERADO** (`Lexer`, `Parser`, `sym` en `target/generated-sources/`). Jamás editar el código generado: los cambios van en `Lexer.flex` / `parser.cup`. La clase `sym` la generan las declaraciones `terminal` del `.cup` (los nombres deben coincidir con lo que usa el `.flex`).
- `dataforge.interprete` — el runtime: `Entorno` (tabla de símbolos con claves lowercase —el case-insensitive alcanza a los identificadores—, consola, `RegistroError`, gráficas, tokens para reportes), `Operaciones` (aritmética/estadística con chequeo de tipos), `Interprete` (fachada String→Entorno; **conecta `lexer.entorno`** para que el lexer registre tokens y errores léxicos).
- `dataforge.gui` — `EditorApp` (UI por código, sin FXML), `Lanzador` (esquiva el chequeo del launcher de Java), `Graficador` (Grafica→Chart→Stage; los datos llegan YA validados por `Entorno.validarGrafica`).
- `dataforge.reportes` — `Reportes` genera los 3 HTML; nombres de token por reflexión sobre `sym`.

Convenciones críticas:
- **Entorno fresco por ejecución** (el enunciado §6 exige reportes solo del último análisis).
- **Propagación por null**: una expresión con error devuelve null y las operaciones que lo reciben callan — un error por causa, sin cascadas, sin abortar.
- **Dos formateadores**: `Entorno.formatear()` (consola: cadenas sin comillas, 15.0→"15") vs `valorReporte()` (reporte §6.3: cadenas con comillas, arreglos sin `.0`). No mezclarlos.
- Recuperación de errores: léxico descarta el carácter; sintáctico usa **modo pánico** (`instruccion ::= error PUNTO_COMA`); semántico propaga null.
- En `.cup`: declaraciones `non terminal` TODAS antes de las producciones, tipos raw (`ArrayList`, no genéricos ni arrays).
- En `.flex`: `%public` es obligatorio (sin él la clase queda package-private); reservadas ANTES de `{Id}`.
- `entradas/*.df` son los casos de prueba (ejemplo1 básico, ejemplo2 gráficas, ejemplo3 errores semánticos, ejemplo4 los 3 tipos de error). `docs/gramatica.txt` es la BNF **entregable** — mantenerla sincronizada si la gramática cambia.

## Presentación del curso (`Hades\presentacion-dataforge\`, repo git)

Material de estudio permanente del usuario; se abre con doble clic (autocontenido, sin fetch — el navegador lo bloquea en file://).

- **Estructura modular**: `index.html` = índice/roadmap; una página por etapa (`etapaN.html`) y por profundización (`gramaticas`, `automatas`, `tabla-simbolos`, `ast`, `fases`, `semantica`); comparten `assets/estilo.css` y `assets/deck.js`.
- **Agregar contenido** = crear la página copiando el esqueleto de una existente (solo cambian las `<section class="slide">`) + actualizar el roadmap del índice y el README + commit. Las demos animadas usan el patrón "stepper" (`[data-step]` + botón ▶ Paso) — ver `gramaticas.html` como referencia.
- **Ritual por etapa del curso**: lección en su `etapaN.html` + código real en `DataForge\` + al verificar el usuario, marcar completada (roadmap/README), actualizar la **Guía DataForge del vault** con el material real, y commitear.

## Cómo operar el vault (cerebro diamon `compiladores`, mundo personal)

- **Buscar/leer**: `mcp__diamon__brain_search` (FTS5: acentos-insensible, prefijos, los `aliases` actúan como sinónimos) y `brain_get_note`. Contexto completo: resource `brain://compiladores`.
- **Crear/editar notas**: preferí `brain_upsert_note` (indexa al instante). Si escribís con Write/Edit, corré `brain_reindex` después.
- **Tras cualquier tanda de cambios**: `brain_health` → política del vault: **0 enlaces rotos, 0 huérfanas**. Ojo: `\|` dentro de tablas rompe wikilinks para el indexador — usar enlaces sin alias dentro de tablas.
- **Convenciones de nota**: frontmatter con `tags`, `aliases` curados, `fuente` y `fecha`; wikilinks liberalmente (también en prosa); diagramas en **Mermaid**.

## Leer los PDFs de doc\

El Read nativo de PDF **no funciona en esta máquina** (falta poppler). Usá `mcp__diamon__doc_convert` (markitdown, cachea por mtime) o `/doc2md`, y refiná con `/refinar-md` (los PDFs de enunciados generan ~60 pseudo-tablas por conversión). Los `.clean.md` ya refinados son la fuente preferida.

## Stack de los proyectos (decidido en Fase 2; detalles en las Guías del vault)

| Proyecto | Stack | Build |
|---|---|---|
| DataForge ✅, ConjAnalyzer, CompScript | Java 17 + JFlex + CUP + JavaFX + Maven | `jflex-maven-plugin:1.9.1` + `cup-maven-plugin:11b-20160615-3` generan en `generate-sources` → `mvn clean javafx:run` (el `pom.xml` de DataForge es la referencia funcionando) |
| CompInterpreter | JS/TS + Jison + React (cliente) + Express (servidor REST) | `npx jison src/grammar.jison -o src/parser.js` · `node server.js` · `npm start` |

- Los 4 proyectos son **intérpretes**: generación de código/back-end fuera de alcance (caps. 8-12 del Dragón).
- La gramática BNF entregable es un documento aparte, **no** copia del `.cup`/`.jison`.
- Al construir cada proyecto, **actualizar su Guía en `04-Proyectos\`** con el material real (la Guía DataForge muestra el formato: secciones "CONSTRUIDO ✅" con decisiones, fragmentos y errores comunes fechados).
