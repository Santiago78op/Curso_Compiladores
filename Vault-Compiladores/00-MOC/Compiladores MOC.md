---
tags: [moc, compiladores]
fuente: "Curso OLC1 (USAC) · Libro del Dragón 2ª ed. · Guías de tecnologías"
fecha: 2026-07-10
---

# 🧭 Compiladores MOC

Mapa de contenido del segundo cerebro de **teoría de compiladores** (curso *Organización de Lenguajes y Compiladores 1*, USAC). Punto de entrada a todas las notas.

## 🧩 Conceptos (fases del compilador)
**Panorama:** [[Fases de un compilador]] · [[Compilador vs intérprete]] · [[Manejo de errores (léxicos, sintácticos, semánticos)]]

**Léxico** `#lexico`: [[Token, lexema y patrón]] · [[Expresiones regulares]] · [[Definiciones regulares]] · [[Autómata finito (AFN y AFD)]] · [[Construcción de Thompson]] · [[Construcción de subconjuntos]] · [[Minimización de AFD]]

**Sintáctico** `#sintactico`: [[Gramática libre de contexto (BNF)]] · [[Derivaciones y árbol de análisis sintáctico]] · [[Ambigüedad, precedencia y asociatividad]] · [[Recursividad por la izquierda y factorización]] · [[FIRST y FOLLOW]] · [[Análisis sintáctico descendente LL(1)]] · [[Análisis sintáctico ascendente LR]] · [[Conflictos shift-reduce y reduce-reduce]]

**Semántico** `#semantico`: [[Traducción dirigida por la sintaxis]] · [[Atributos sintetizados y heredados]] · [[Árbol de sintaxis abstracta (AST)]] · [[Comprobación de tipos]] · [[Conversión de tipos (coerción y cast)]] · [[Flujo de control y switch]] · [[Tabla de símbolos]] · [[Entornos y alcance]] · [[Paso de parámetros]] · [[Registro de activación y pila de control]]

**Intermedio** `#intermedio` *(contexto, no requerido en proyectos)*: [[Código de tres direcciones]]

## 🛠️ Tecnologías
[[JFlex]] · [[CUP]] · [[Jison]] · [[ANTLR]] · [[JavaFX y Scene Builder]] · [[Maven]] · [[Mermaid]] · [[Graphviz]] · [[vis-network]] · [[Codebase-Memory-MCP]]

## 📖 Libro del Dragón (resúmenes por capítulo)
[[Cap 1 - Introducción]] · [[Cap 2 - Traductor simple orientado a la sintaxis]] · [[Cap 3 - Análisis léxico]] · [[Cap 4 - Análisis sintáctico]] · [[Cap 5 - Traducción dirigida por la sintaxis]] · [[Cap 6 - Generación de código intermedio]] · [[Cap 7 - Entornos en tiempo de ejecución]] · [[Caps 8-12 - Panorama (fuera de alcance)]]

## 🚀 Proyectos del curso (OLC1)
[[DataForge]] · [[ConjAnalyzer]] · [[CompScript]] · [[CompInterpreter]]

**Guías de elaboración (paso a paso):** [[Guía DataForge]] · [[Guía ConjAnalyzer]] · [[Guía CompScript]] · [[Guía CompInterpreter]]

## 🚀 Proyectos de OLC2 (curso hermano, mismo cerebro)
[[VLangCherry]] — Proyecto 1 de *Organización de Lenguajes y Compiladores 2*, grupal (3 integrantes): intérprete con [[ANTLR]] + Go, arquitectura cliente-servidor REST igual patrón que [[CompInterpreter]].

**Guía de elaboración:** [[Guía VLangCherry]]

## 📚 Recursos
[[Libro del Dragón (ficha)]] · [[Guías de tecnologías]] · [[Modern Compiler Implementation in Java (Appel)]] · [[Appel vs Libro del Dragón]]

## 🗺️ Orden de estudio sugerido
1. [[Fases de un compilador]] → [[Cap 1 - Introducción]] y [[Cap 2 - Traductor simple orientado a la sintaxis]]
2. Léxico → [[Cap 3 - Análisis léxico]] + [[JFlex]] / [[Jison]]
3. Sintáctico → [[Cap 4 - Análisis sintáctico]] + [[CUP]] / [[Jison]]
4. Semántico → [[Cap 5 - Traducción dirigida por la sintaxis]] + [[Cap 6 - Generación de código intermedio]]
5. Entornos → [[Cap 7 - Entornos en tiempo de ejecución]]
6. Proyectos: [[DataForge]] → [[ConjAnalyzer]] → [[CompScript]] → [[CompInterpreter]]
