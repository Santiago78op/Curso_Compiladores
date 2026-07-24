---
tags: [concepto, sintactico]
aliases: ["precedencia por niveles", "gramática no ambigua", "n+1 no terminales", "expr term factor"]
fuente: "Libro del Dragón §2.2.5–2.2.6 (precedencia estructural) y cap. 4.8 (precedencia por declaraciones)"
fecha: 2026-07-24
---

# Ambigüedad, precedencia y asociatividad

Una gramática es **ambigua** si una cadena tiene ≥2 árboles de análisis (p. ej. `E → E+E | E*E`, o `9+5*2` que sin reglas admite `(9+5)*2` y `9+(5*2)`). Se resuelve fijando:
- **Precedencia:** qué operador agrupa primero (`*` recibe sus operandos antes que `+`).
- **Asociatividad:** izquierda / derecha / no asociativa.

Hay **dos formas** de imponerlas, y conviene no confundirlas:

## 1. Por declaraciones (en el generador)

`precedence left/right` en [[CUP]] y `%left/%right/%nonassoc` en [[Jison]] (y `%prec` para casos como el menos unario). El generador usa esas declaraciones para resolver los [[Conflictos shift-reduce y reduce-reduce|conflictos shift-reduce]] de una gramática ambigua, sin cambiar la gramática. El *else* colgante se resuelve así, asociando el `else` con el `if` abierto más cercano (SHIFT).

## 2. Por la estructura de la gramática (§2.2.6 + recuadro "Generalización", pp. 49–50)

Se codifica la precedencia **en la gramática misma**, sin declaraciones. Un `factor` es "una expresión que **no puede separarse** por ningún operador" (un operando o una `( expr )` protegida por paréntesis); un `term` solo puede separarse por los operadores de mayor precedencia (`* /`); una `expr` puede separarse por cualquiera:

```
expr   → expr + term | expr - term | term
term   → term * factor | term / factor | factor
factor → dígito | ( expr )
```

**Regla general (p. 50):** para *n* niveles de precedencia se usan **n+1 no terminales**. El primero (`factor`) nunca se separa; cada nivel siguiente representa las expresiones separables solo por operadores de ese nivel o superior, y su producción termina apuntando al nivel de arriba. La asociatividad sale de la forma: `expr → expr + term` (recursión **izquierda**) da asociación por la izquierda.

> **Por qué importa para el curso:** la **BNF entregable** de cada proyecto debe ser **no ambigua por sí sola** — no lleva tablas de precedencia adjuntas. Escribirla con la construcción por niveles (`expr`/`term`/`factor`, n+1 no terminales) es la única forma correcta de entregarla, aunque el `.cup`/`.jison` real use `%left` sobre una gramática ambigua más corta. Ver [[Gramática libre de contexto (BNF)]].

## Relacionadas
- [[Gramática libre de contexto (BNF)]]
- [[Conflictos shift-reduce y reduce-reduce]]
- [[Recursividad por la izquierda y factorización]]
- [[Análisis sintáctico ascendente LR]]
- [[CompScript]]
