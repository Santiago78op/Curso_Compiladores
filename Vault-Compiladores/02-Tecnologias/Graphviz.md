---
tags: [tecnologia, visualizacion]
fuente: "Enunciado CompScript (sugerido) · guías de tecnologías"
fecha: 2026-07-10
---

# Graphviz

Herramienta y lenguaje **DOT** para dibujar grafos. Se genera un archivo `.dot` desde el [[Árbol de sintaxis abstracta (AST)|AST]] y Graphviz produce la imagen del árbol.

Ejemplo DOT:
```dot
digraph AST {
  suma -> id_x;
  suma -> mul;
  mul -> id_y;
  mul -> num_60;
}
```

Es la opción recomendada por el enunciado de [[CompScript]] para el **reporte AST** en la app de escritorio. En web se prefiere [[vis-network]]; en el vault, [[Mermaid]].

## Usado en
[[CompScript]]

## Relacionadas
- [[Mermaid]]
- [[vis-network]]
- [[Árbol de sintaxis abstracta (AST)]]
