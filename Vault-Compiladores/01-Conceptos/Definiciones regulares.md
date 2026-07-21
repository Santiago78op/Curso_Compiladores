---
tags: [concepto, lexico]
fuente: "Libro del Dragón, cap. 3.3.4"
fecha: 2026-07-10
---

# Definiciones regulares

Secuencia de nombres `d₁→r₁, …, dₙ→rₙ` **no recursiva** (cada `rᵢ` usa solo el alfabeto y los `d` previos). Permiten dar nombre a [[Expresiones regulares]] reutilizables.

Ejemplo del libro:
```
digito  → [0-9]
digitos → digito+
numero  → digitos ( . digitos )? ( E [+-]? digitos )?
```

> Son exactamente las **macros** de [[JFlex]] (`Digit = [0-9]`) y las definiciones del `%lex` de [[Jison]].

## Relacionadas
- [[Expresiones regulares]]
- [[JFlex]]
- [[Jison]]
