---
tags: [concepto, semantico]
aliases: ["síntesis de tipos", "inferencia de tipos", "sobrecarga de operadores", "sobrecarga por firma", "resolución de sobrecarga"]
fuente: "Libro del Dragón §6.5 (§6.5.3 sobrecarga) y §2.8.3"
fecha: 2026-07-24
---

# Comprobación de tipos

Parte del **análisis semántico**: verifica que los operadores reciban operandos compatibles. Dos formas:
- **Síntesis de tipos:** construye el tipo de una expresión desde sus subexpresiones (`si f: s→t y x: s, entonces f(x): t`). Requiere declarar antes de usar.
- **Inferencia de tipos:** deduce el tipo por el uso (lenguajes tipo ML).

Las instrucciones se tratan como funciones: `if (E) S` espera `(bool, void) → void`.

## Sobrecarga de operadores y funciones (§6.5.3)

Un símbolo **sobrecargado** tiene **distintos significados según su contexto**, y la sobrecarga se **resuelve** eligiendo un significado único por la **firma** (el operador/nombre **más** los tipos de sus operandos/argumentos). El caso emblema (ejemplo 6.13): el **`+`** de Java es **suma** entre números y **concatenación** entre cadenas — se decide mirando los tipos de los operandos.

> **Conexión directa con los proyectos:** las **tablas de compatibilidad de tipos** de [[CompScript]]/[[CompInterpreter]] y el método `Operaciones.suma` de [[DataForge]] (que resuelve `+` como suma numérica **o** concatenación según los tipos) **son sobrecarga de operadores** en el sentido de §6.5.3 — aunque el código nunca use la palabra. Cada celda de la tabla `⟨op, tipoIzq, tipoDer⟩ → tipoResultado` es una firma; una combinación sin entrada = **error semántico**. (El "sin sobrecarga" que a veces se dice de estos lenguajes se refiere a *funciones* de usuario con igual nombre, concepto distinto.)

La resolución por firma se implementa eficientemente con el **número de valor** (hash de `⟨etiqueta, hijos⟩`), y es síntesis de tipos: la regla (6.10) elige `f(x): tₖ` según el tipo `sₖ` del argumento.

## Relacionadas
- [[Conversión de tipos (coerción y cast)]]
- [[Manejo de errores (léxicos, sintácticos, semánticos)]]
- [[Atributos sintetizados y heredados]]
- [[Cap 6 - Generación de código intermedio]]
