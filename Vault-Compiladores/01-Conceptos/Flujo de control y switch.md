---
tags: [concepto, semantico]
aliases: ["break y continue", "señales de control", "listas de saltos", "backpatch", "switch de n vías", "fall-through", "cortocircuito"]
fuente: "Libro del Dragón §6.6 (flujo), §6.7.4 (break/continue), §6.8 (switch)"
fecha: 2026-07-24
---

# Flujo de control y switch

- **`if`, `while`, `for`:** se evalúa una **expresión booleana** y se elige la rama. Con **cortocircuito**, `&&`/`||` solo evalúan lo necesario (importa cuando la derecha podría fallar, p. ej. la guarda `x != 0 && 10/x > 1` — bug real corregido en [[VLangCherry]] y [[CompInterpreter]]).
- **`switch`/`match`:** **bifurcación de n vías**: evaluar E, buscar el caso igual (o `default`), ejecutar su instrucción.

> En un [[Compilador vs intérprete|intérprete]] no se generan etiquetas/gotos: se recorre el [[Árbol de sintaxis abstracta (AST)|AST]] y se ejecuta la rama directamente. Base del `if/else if`, ciclos, `match` (CompScript) y `switch case` con *fall-through* (CompInterpreter).

## `break` / `continue` / `return`: compilador vs intérprete (§6.7.4)

Pregunta doble de defensa — *"¿cómo funciona tu `break`?"* y *"¿cómo lo haría un compilador?"*:

- **Compilador (§6.7.4):** `break` se traduce como un salto **sin destino conocido todavía** (aún no se sabe la etiqueta de salida del ciclo). Se acumula en una **lista de saltos** de la construcción envolvente y se rellena por **backpatch** cuando esa etiqueta se conoce. Igual `continue` (salto al inicio/incremento).
- **Intérprete (los proyectos):** al recorrer el AST no hay etiquetas; el cuerpo devuelve (o lanza) un **valor de señal** `Break`/`Continue`/`Return` que el **nodo de ciclo más cercano intercepta** y decide qué hacer. Es `Senales.Break` en [[CompScript]] (excepción de control), el objeto-señal `{tipo:'BREAK'}` en [[CompInterpreter]], y la `Senal{Tipo: SenalBreak}` en [[VLangCherry]]. Una señal que escapa hasta el cuerpo de la función sin que ningún ciclo la consuma = **error semántico** ("break fuera de un ciclo").

## Estrategias para implementar el `switch` (§6.8.1)

Cómo se compila la bifurcación de n vías (aplica también a un `match` grande en un intérprete):
1. **Saltos condicionales secuenciales** (una cadena de comparaciones) — bien para **pocos** casos (≲10).
2. **Tabla hash** valor→etiqueta — para muchos casos dispersos (un `HashMap` en el intérprete en vez de if-else en cadena).
3. **Arreglo de baldes** indexado por el valor — cuando los casos caen en un **rango denso**.

## Relacionadas
- [[Árbol de sintaxis abstracta (AST)]]
- [[Comprobación de tipos]]
- [[CompInterpreter]]
- [[Cap 6 - Generación de código intermedio]]
