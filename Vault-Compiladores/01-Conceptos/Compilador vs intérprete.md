---
tags: [concepto, compiladores]
fuente: "Libro del Dragón, cap. 1"
fecha: 2026-07-10
---

# Compilador vs intérprete

- **Compilador:** traduce el programa fuente a un programa destino equivalente (p. ej. código máquina), que luego se ejecuta. Suele ser más rápido en ejecución.
- **Intérprete:** no produce traducción; **ejecuta directamente** las operaciones del fuente sobre las entradas. Da mejores diagnósticos de error.
- **Híbrido (Java):** compila a *bytecodes* y una máquina virtual los interpreta (con JIT los traduce a máquina justo antes de ejecutar).

> Los 4 proyectos del curso son **intérpretes**: recorren el [[Árbol de sintaxis abstracta (AST)]] y ejecutan, sin generar código máquina.

## Relacionadas
- [[Fases de un compilador]]
- [[Árbol de sintaxis abstracta (AST)]]
- [[CompScript]]
