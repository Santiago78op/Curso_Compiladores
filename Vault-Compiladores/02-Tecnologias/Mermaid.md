---
tags: [tecnologia, visualizacion, obsidian]
fuente: "Recomendación de guías de tecnologías"
fecha: 2026-07-10
---

# Mermaid

Lenguaje para crear **diagramas desde texto**. **Obsidian lo renderiza nativo**, por lo que es la herramienta elegida para dibujar [[Autómata finito (AFN y AFD)|AFDs]], diagramas de estados y [[Árbol de sintaxis abstracta (AST)|ASTs]] **directo dentro de las notas** de este vault.

Ejemplo (AFD):
```mermaid
stateDiagram-v2
    [*] --> S0
    S0 --> S1: letra
    S1 --> S1: letra | dígito
    S1 --> [*]
```

Para diagramas de estados muy complejos, **PlantUML** es el respaldo.

## Usado en
El propio vault (todas las notas con diagramas)

## Relacionadas
- [[Graphviz]]
- [[vis-network]]
- [[Árbol de sintaxis abstracta (AST)]]
