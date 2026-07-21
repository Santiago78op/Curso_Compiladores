---
tags: [concepto, semantico]
fuente: "Libro del Dragón, cap. 6.5"
fecha: 2026-07-10
---

# Comprobación de tipos

Parte del **análisis semántico**: verifica que los operadores reciban operandos compatibles. Dos formas:
- **Síntesis de tipos:** construye el tipo de una expresión desde sus subexpresiones (`si f: s→t y x: s, entonces f(x): t`). Requiere declarar antes de usar.
- **Inferencia de tipos:** deduce el tipo por el uso (lenguajes tipo ML).

Las instrucciones se tratan como funciones: `if (E) S` espera `(bool, void) → void`.

> Las **tablas de compatibilidad de tipos** de CompScript/CompInterpreter (qué da `+` entre int y string, etc.) son reglas de síntesis; una combinación no permitida = **error semántico**.

## Relacionadas
- [[Conversión de tipos (coerción y cast)]]
- [[Manejo de errores (léxicos, sintácticos, semánticos)]]
- [[Atributos sintetizados y heredados]]
