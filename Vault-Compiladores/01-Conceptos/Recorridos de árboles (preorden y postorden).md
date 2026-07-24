---
tags: [concepto, semantico, sintactico]
aliases: ["preorden", "postorden", "recorrido postorden", "notación postfija", "postfija", "por qué un intérprete evalúa"]
fuente: "Libro del Dragón §2.3.4 (recorridos) y §2.3.1 (notación postfija), pp. 53–58"
fecha: 2026-07-24
---

# Recorridos de árboles (preorden y postorden)

Un **recorrido** de un árbol lo visita desde la raíz hacia abajo. Preorden y postorden son los dos casos que importan (§2.3.4, fig. 2.11):

- **Preorden** de un (sub)árbol con raíz `N`: **primero `N`**, luego los recorridos preorden de sus hijos de izquierda a derecha (la acción se hace *al llegar* al nodo).
- **Postorden** de `N`: primero los recorridos postorden de los hijos, y **`N` al final** (la acción se hace *justo antes de dejar* el nodo por última vez).

El procedimiento `visitar(N)` de la fig. 2.11 es un recorrido **postorden**.

## Por qué esto ES la teoría de un intérprete

**Ejecutar el AST = recorrido en postorden.** Para evaluar un nodo (`E1 op E2`) hay que tener ya los valores de sus hijos, así que se evalúan primero los hijos y luego se aplica el operador del nodo. Los **[[Atributos sintetizados y heredados|atributos sintetizados]]** (los que suben) se calculan justo en cualquier recorrido de **abajo hacia arriba**. El método `evaluar()` / `ejecutar()` / `visit()` de [[CompScript]] y [[CompInterpreter]] es exactamente el `visitar(N)` postorden de la fig. 2.11: cada nodo del [[Árbol de sintaxis abstracta (AST)|AST]] "sabe" resolverse consultando primero a sus hijos.

## Notación postfija: la salida del postorden (§2.3.1, fig. 2.14)

La **notación postfija** (o polaca inversa) pone el operador **después** de sus operandos: `x - y` → `xy-`. Su definición es inductiva: un operando queda igual; `E1 op E2` postfija es `postfija(E1) postfija(E2) op`. Se **evalúa** buscando el operador de más a la izquierda cuyos operandos ya estén disponibles.

Colgar acciones `print` del operador en cada nodo y ejecutarlas en **postorden** produce la postfija. Ejemplo del libro (fig. 2.14): `9 - 5 + 2` traducido en postorden imprime **`95-2+`**.

> Contraste útil: si esas acciones se ejecutaran en **preorden**, saldría notación **prefija** (`+-952`). El orden del recorrido es lo que fija la notación.

## Relacionadas
- [[Árbol de sintaxis abstracta (AST)]]
- [[Traducción dirigida por la sintaxis]]
- [[Atributos sintetizados y heredados]]
- [[CompScript]]
- [[Cap 2 - Traductor simple orientado a la sintaxis]]
