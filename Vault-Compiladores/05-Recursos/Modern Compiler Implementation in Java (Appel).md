---
tags: [recurso, libro, java]
fuente: "https://www.cs.princeton.edu/~appel/modern/java/ — A. Appel, Cambridge University Press, 1998"
fecha: 2026-07-10
---

# Modern Compiler Implementation in Java (Appel)

**Autor:** Andrew W. Appel (Princeton). **Editorial:** Cambridge University Press, 1998 (edición preliminar 1997). Conocido como el **"libro del Tigre"**.

## Enfoque
Libro **guiado por proyecto**: el estudiante construye, **módulo a módulo**, un compilador completo para el lenguaje de juguete **Tiger**. Es implementación-primero (trae código Java), a diferencia del estilo enciclopédico del [[Libro del Dragón (ficha)|Libro del Dragón]].

## Herramientas (¡las del curso!)
- **JLex** — generador léxico, **antecesor directo de [[JFlex]]**.
- **CUP** — el **mismo** [[CUP]] que usan los proyectos Java.
- **SPIM** — simulador MIPS/RISC (para el back-end; fuera de alcance del curso).

## Estructura publicada del libro
> La TOC en línea (`contents.html`) devolvió 404; esta es la estructura bibliográfica del libro.

**Parte I — Fundamentos de la compilación** (relevante al curso): Introducción · Análisis léxico · Análisis sintáctico · Sintaxis abstracta · Análisis semántico · Registros de activación · Traducción a código intermedio · Bloques básicos y trazas · Selección de instrucciones · Análisis de vivacidad · Asignación de registros · Integración final.

**Parte II — Temas avanzados** (fuera de alcance, como los [[Caps 8-12 - Panorama (fuera de alcance)|caps. 8–12]] del Dragón): Recolección de basura · Lenguajes OO · Lenguajes funcionales · Tipos polimórficos · Análisis de flujo de datos · Optimización de ciclos · Forma SSA · Pipelining · Jerarquía de memoria.

## Recursos disponibles
Table of Contents y Preface, **lista de erratas**, software de apoyo (JLex, CUP, SPIM) y los **módulos del compilador Tiger** para los ejercicios. Existen versiones gemelas en **ML** y en **C**.

## Por qué te sirve
Como usa **JLex + CUP**, es la **mejor referencia de implementación** para [[DataForge]], [[ConjAnalyzer]] y [[CompScript]] (Java). Ver la comparación: [[Appel vs Libro del Dragón]].

## Relacionadas
- [[Appel vs Libro del Dragón]]
- [[CUP]]
- [[JFlex]]
