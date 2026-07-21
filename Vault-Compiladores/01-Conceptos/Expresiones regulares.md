---
tags: [concepto, lexico]
fuente: "Libro del Dragón, cap. 3.3"
fecha: 2026-07-10
---

# Expresiones regulares

Notación para describir los **patrones** de los [[Token, lexema y patrón|tokens]]. Operaciones **básicas** (precedencia mayor→menor): `*` (cerradura de Kleene, cero o más), concatenación, `|` (unión).

**Extensiones** (abreviaturas definibles con las básicas): `+` (una o más, `r+ = rr*`), `?` (cero o una, `r? = r|ε`), y las que usan Lex/JFlex/Jison: clases `[abc]`, rangos `[a-z]`, negación `[^…]`.

Ejemplo (identificador): `letra_ ( letra_ | digito )*`.

Toda ER se puede convertir en un [[Autómata finito (AFN y AFD)|autómata finito]] (ver [[Construcción de Thompson]]).

## Relacionadas
- [[Definiciones regulares]]
- [[Autómata finito (AFN y AFD)]]
- [[JFlex]]
