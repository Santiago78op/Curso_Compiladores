---
tags: [proyecto, lexico, sintactico]
fuente: "Enunciado OLC1_PT1_S2024 (Proyecto 1, 2do Sem 2024)"
fecha: 2026-07-10
---

# ConjAnalyzer

**Objetivo:** analizador [[Cap 3 - Análisis léxico|léxico]] y [[Cap 4 - Análisis sintáctico|sintáctico]] para **operaciones entre conjuntos** en **notación prefija/polaca**. Extensión `.ca`, *case sensitive*, encapsulado en `{ }`.

**Etapa del compilador:** léxico + sintáctico (+ evaluación y "simplificación").

**Tecnologías:** [[JFlex]] + [[CUP]] (Java), GUI [[JavaFX y Scene Builder]], [[Maven]], diagrama de **Venn** con JavaFX Canvas, salida **JSON** con **Gson**, reportes HTML.

**Rasgos:** define conjuntos (`CONJ`), universo ASCII 33–126, operaciones `U`, `&`, `^`, `-` prefijas, `EVALUAR(...)`; aplica **propiedades de teoría de conjuntos** (DeMorgan, doble complemento…) para simplificar y volcar a `.json`.

**Conceptos aplicados:** [[Gramática libre de contexto (BNF)]] · [[Derivaciones y árbol de análisis sintáctico]] · [[Análisis sintáctico ascendente LR]] · [[Tabla de símbolos]]. La "simplificación" rima con la optimización algebraica ([[Caps 8-12 - Panorama (fuera de alcance)]]).

## Guía de elaboración
- [[Guía ConjAnalyzer]]

## Relacionadas
- [[DataForge]]
- [[CompScript]]
- [[Cap 4 - Análisis sintáctico]]
