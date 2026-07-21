---
tags: [proyecto, lexico, sintactico]
fuente: "Enunciado OLC1_PT1_1S2024 (Proyecto 1, 1er Sem 2024)"
fecha: 2026-07-10
---

# DataForge

**Objetivo:** aplicar análisis [[Cap 3 - Análisis léxico|léxico]] y [[Cap 4 - Análisis sintáctico|sintáctico]] para un lenguaje de operaciones aritméticas, estadísticas y **graficación** de datos. Extensión `.df`, *case insensitive*, encapsulado en `PROGRAM … END PROGRAM`.

**Etapa del compilador:** léxico + sintáctico (+ evaluación/ejecución).

**Tecnologías:** [[JFlex]] + [[CUP]] (Java), GUI con [[JavaFX y Scene Builder]], build con [[Maven]], gráficas con **JavaFX Charts**, reportes en HTML.

**Rasgos:** tipos `double` y `char[]`; arreglos (`@`); operaciones como funciones (`SUM/RES/MUL/DIV/MOD`, anidables); estadísticas (`Media, Mediana, Moda, Varianza, Max, Min`); impresión (`console::print`, `console::column`); 4 gráficas (barras, pie, línea, histograma). Reportes: tokens, errores, tabla de símbolos.

**Conceptos aplicados:** [[Token, lexema y patrón]] · [[Expresiones regulares]] · [[Gramática libre de contexto (BNF)]] · [[Ambigüedad, precedencia y asociatividad]] · [[Tabla de símbolos]] · [[Manejo de errores (léxicos, sintácticos, semánticos)]]

## Guía de elaboración
- [[Guía DataForge]]

## Relacionadas
- [[ConjAnalyzer]]
- [[Cap 3 - Análisis léxico]]
- [[Cap 4 - Análisis sintáctico]]
