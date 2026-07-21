---
tags: [tecnologia, visualizacion, javascript, web]
fuente: "Guías de tecnologías (alternativas web)"
fecha: 2026-07-10
---

# vis-network

Librería JavaScript para **grafos interactivos** (arrastrar nodos, zoom) en el navegador. Recomendada para el **reporte AST** de [[CompInterpreter]], que corre en web y luce mejor con un árbol interactivo que con una imagen estática.

Más simple que D3.js para un proyecto de curso.

## Ejemplo mínimo (AST de `a + 3`)
```html
<div id="ast" style="height:400px"></div>
<script src="https://unpkg.com/vis-network/standalone/umd/vis-network.min.js"></script>
<script>
  const nodes = new vis.DataSet([
    { id: 1, label: "+" },
    { id: 2, label: "a" },
    { id: 3, label: "3" },
  ]);
  const edges = new vis.DataSet([
    { from: 1, to: 2 },
    { from: 1, to: 3 },
  ]);
  new vis.Network(document.getElementById("ast"), { nodes, edges }, {
    layout: { hierarchical: { direction: "UD" } }  // árbol de arriba hacia abajo
  });
</script>
```

Patrón para el proyecto: al recorrer el [[Árbol de sintaxis abstracta (AST)|AST]], cada nodo agrega su entrada a `nodes` (con un `id` autoincremental) y una arista `{from: idPadre, to: idHijo}` a `edges`.

## Usado en
[[CompInterpreter]]

## Relacionadas
- [[Graphviz]]
- [[Mermaid]]
- [[Árbol de sintaxis abstracta (AST)]]
