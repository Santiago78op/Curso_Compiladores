---
tags: [concepto, intermedio]
aliases: [3 direcciones, three-address code, TAC, cuádruplos, tripletas, representación intermedia]
fuente: "Libro del Dragón, cap. 6.2"
fecha: 2026-07-10
---

# Código de tres direcciones

> **Contexto teórico, no requerido en los proyectos** (los 4 proyectos del curso son [[Compilador vs intérprete|intérpretes]] que ejecutan el [[Árbol de sintaxis abstracta (AST)|AST]] directamente). Se incluye porque es el tema nuclear del cap. 6 y aparece en exámenes de teoría.

Representación intermedia **lineal** (no arbórea) donde cada instrucción tiene a lo sumo **un operador** y **tres "direcciones"** (dos operandos + un resultado). Las expresiones complejas se descomponen usando **temporales** (`t1`, `t2`, …).

Ejemplo — `a + b * 3` se convierte en:

```
t1 = b * 3
t2 = a + t1
```

## Formas de implementarlo

| Forma | Qué guarda | Ventaja |
|-------|-----------|---------|
| **Cuádruplos** | `(op, arg1, arg2, resultado)` — el temporal es explícito | Fácil de reordenar (optimización) |
| **Tripletas** | `(op, arg1, arg2)` — el resultado es *la posición* de la tripleta | Más compacto, sin nombres de temporales |
| **Tripletas indirectas** | Lista de punteros a tripletas | Reordenable sin romper referencias |

## Backpatching (parcheo hacia atrás)

Técnica para generar **saltos** (`goto`, `if x goto L`) en una sola pasada: cuando la etiqueta destino aún no se conoce, el salto se emite **incompleto** y su posición se guarda en una lista; al conocerse el destino, se "parchea" la lista completa. Es como el flujo de control (`if`, `while`, `&&` con corto circuito) se traduce a código lineal.

## Relacionadas
- [[Cap 6 - Generación de código intermedio]]
- [[Árbol de sintaxis abstracta (AST)]]
- [[Fases de un compilador]]
- [[Flujo de control y switch]]
