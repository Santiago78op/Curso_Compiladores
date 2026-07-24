---
tags: [concepto, semantico]
aliases: [SDT, SDD, TDS, syntax-directed translation, reglas semánticas, acciones semánticas, "SDT postfijo", "atributos en la pila", "S-atribuida"]
fuente: "Libro del Dragón cap. 5 (§5.4.1 SDT postfijo, §5.4.2 atributos en la pila)"
fecha: 2026-07-24
---

# Traducción dirigida por la sintaxis

Asocia a cada producción de la gramática **reglas semánticas** (una **SDD**) o **acciones** escritas dentro de la producción (un **SDT**) que calculan [[Atributos sintetizados y heredados|atributos]].

Una SDD con **solo atributos sintetizados** (**S-atribuida**) se implementa directo con un parser LR → por eso las acciones `{: RESULT = … :}` de [[CUP]] y `{ $$ = … }` de [[Jison]] son reglas semánticas.

## SDT postfijo: la formalización de DataForge (§5.4.1)

Un **SDT postfijo** pone **todas** las acciones al **final** de la producción, y se ejecutan justo al **reducir**. Como toda SDD S-atribuida se puede volver postfija, es implementable en la **pila del parser LR sin construir árbol**. Esa es la definición exacta de la arquitectura de [[DataForge]]: *"gramática S-atribuida, ejecución directa en las acciones, sin AST"* (Dragón §5.4.1–5.4.2, pp. 324–327). [[CompScript]], en cambio, usa las acciones para **construir** el AST y lo recorre después.

## Los atributos viven EN la pila (§5.4.2, fig. 5.20)

La mecánica física detrás de `RESULT`/`$$`/`$1`: cada registro de la pila LR guarda **estado + atributo(s)** (un campo `val`). Al reducir `E → E₁ + T`, el cuerpo ocupa las 3 posiciones superiores —`E₁` en `tope−2`, `+` en `tope−1`, `T` en `tope`— y la acción del libro es:

```text
E → E + T   { pila[tope-2].val = pila[tope-2].val + pila[tope].val; }
```

que en CUP se escribe `expr:a PLUS expr:b {: RESULT = a + b; :}`. Es exactamente lo mismo: `RESULT` es el `val` del registro que quedará representando a `E`, y `:a`/`:b` son los `val` de las posiciones del cuerpo. Esto explica:
- por qué `$1`/`$3` (o los `:nombre`) se numeran **por posición** en el cuerpo;
- por qué los atributos **sintetizados son "gratis"** en LALR (ya están en la pila al reducir);
- por qué DataForge **no necesita ninguna estructura extra** para evaluar: la pila del parser YA es su almacén de resultados.

Aplicación típica: construir el [[Árbol de sintaxis abstracta (AST)]] o **evaluar** expresiones al vuelo (ver [[Recorridos de árboles (preorden y postorden)|postorden]]).

## Relacionadas
- [[Atributos sintetizados y heredados]]
- [[Árbol de sintaxis abstracta (AST)]]
- [[Elementos LR(0) y la tabla SLR]]
- [[Cap 5 - Traducción dirigida por la sintaxis]]
