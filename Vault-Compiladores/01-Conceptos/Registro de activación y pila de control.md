---
tags: [concepto, semantico]
aliases: [frame, activation record, stack frame, pila de llamadas, árbol de activación]
fuente: "Libro del Dragón, cap. 7.2"
fecha: 2026-07-10
---

# Registro de activación y pila de control

Cada llamada a función crea un **registro de activación** (frame) con: parámetros, valor de retorno, **enlace de control** (al llamador), **enlace de acceso** (al entorno de declaración), estado de máquina, datos locales y temporales. La **pila de control** guarda las activaciones vivas; el **árbol de activación** representa todas las llamadas de la ejecución.

> En un intérprete, el registro de activación ≈ **entorno** y la pila de control ≈ **pila de entornos**. La recursión funciona porque cada llamada tiene su propio entorno. Los campos de bajo nivel (estado de máquina, temporales) no aplican a un intérprete.

## Relacionadas
- [[Entornos y alcance]]
- [[Tabla de símbolos]]
- [[Cap 7 - Entornos en tiempo de ejecución]]
