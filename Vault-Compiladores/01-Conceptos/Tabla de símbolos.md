---
tags: [concepto, semantico]
fuente: "Libro del Dragón, caps. 2 y 7"
fecha: 2026-07-10
---

# Tabla de símbolos

Estructura de datos que guarda información de cada **nombre** del programa (tipo, valor, alcance, línea/columna). La usan **todas** las [[Fases de un compilador|fases del compilador]].

En un intérprete se organiza como una **pila/árbol de [[Entornos y alcance|entornos]]**: cada bloque/función apila un entorno; buscar un nombre = mirar el entorno actual y subir al padre.

> El **reporte de tabla de símbolos** de los proyectos (id, tipo, entorno, valor, línea, columna) es esta estructura volcada tras la ejecución.

## Relacionadas
- [[Entornos y alcance]]
- [[Registro de activación y pila de control]]
- [[Fases de un compilador]]
