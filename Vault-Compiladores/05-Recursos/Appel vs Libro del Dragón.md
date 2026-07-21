---
tags: [recurso, comparacion, fuente-externa]
fuente: "https://www.cs.princeton.edu/~appel/modern/java/ vs Libro del Dragón 2ª ed. (Fase 4)"
fecha: 2026-07-10
---

# Appel vs Libro del Dragón

Comparación de la **fuente externa** ([[Modern Compiler Implementation in Java (Appel)]]) con el [[Libro del Dragón (ficha)|Libro del Dragón]], integrando aportes al segundo cerebro.

## ✅ Coincidencias
Ambos cubren el mismo **pipeline de front-end**: [[Fases de un compilador|fases]], [[Expresiones regulares|ER]] → [[Autómata finito (AFN y AFD)|AFN/AFD]] en el [[Cap 3 - Análisis léxico|léxico]], parsing [[Análisis sintáctico descendente LL(1)|LL]]/[[Análisis sintáctico ascendente LR|LR]], [[Comprobación de tipos]], y [[Registro de activación y pila de control|registros de activación]]. Ambos usan la tradición **lex/yacc** (JLex+CUP en Appel; Lex/Yacc como referencia en el Dragón).

## 🔀 Diferencias de enfoque
| Aspecto | Libro del Dragón | Appel (Tiger) |
|---|---|---|
| Estilo | Enciclopédico, teórico, de referencia | Guiado por **proyecto**, implementación-primero |
| Teoría de autómatas | Muy amplia (Thompson, subconjuntos, minimización, LR canónico) | Más concisa, lo justo para implementar |
| Código | Fragmentos y ejemplos | **Módulos Java completos** del compilador Tiger |
| Lenguaje objetivo | Ejemplos sueltos | **Tiger**, de principio a fin |
| Herramientas | Lex/Yacc (referencia) | **JLex + CUP** (= stack del curso) |

## 🆕 Material nuevo / valor añadido
- **JLex + CUP como stack real:** [[JFlex]] es el sucesor de JLex y [[CUP]] es idéntico → Appel es la **mejor plantilla de implementación** para los proyectos Java, más que el Dragón.
- **Estructura modular incremental** (Lexer → Parser → Abstract Syntax → Semantic → Activation Records → …) que calca el orden de construcción sugerido en las [[Guía CompScript|guías de proyecto]].
- Temas de back-end (selección de instrucciones, asignación de registros por coloreo, SSA, recolección de basura) → **fuera de alcance** del curso, igual que los [[Caps 8-12 - Panorama (fuera de alcance)|caps. 8–12]] del Dragón.

## Glosario (aportes de esta fuente)
- **Tiger:** lenguaje de juguete que se compila a lo largo del libro de Appel.
- **JLex:** generador léxico predecesor de [[JFlex]].
- **SPIM:** simulador MIPS para ejecutar el código generado (back-end).

## Conclusión para el curso
Usá el **Dragón** para la **teoría** (léxico/sintáctico/semántico, caps. 1–7) y **Appel** para la **implementación en Java** (JLex/JFlex + CUP, estructura modular). Para [[CompInterpreter]] (JS/[[Jison]]) ninguno aplica directo al stack, pero la teoría es la misma.

## Relacionadas
- [[Modern Compiler Implementation in Java (Appel)]]
- [[Libro del Dragón (ficha)]]
- [[Guías de tecnologías]]
