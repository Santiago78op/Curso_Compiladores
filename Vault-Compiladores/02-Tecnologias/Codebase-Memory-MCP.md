---
tags: [tecnologia, herramienta, meta, grafo-de-codigo]
fuente: "https://github.com/DeusData/codebase-memory-mcp"
fecha: 2026-07-21
aliases: [CBM, codebase-memory, grafo de conocimiento]
---

# Codebase-Memory-MCP

Servidor MCP que indexa **código fuente real** (vía tree-sitter, 158 lenguajes) en un grafo de conocimiento persistente (SQLite en `~/.cache/codebase-memory-mcp/`). A diferencia de este vault (notas de prosa sobre teoría y decisiones), este grafo permite **consultar la implementación real** de una herramienta: quién llama a qué función, cómo está organizado un paquete, dónde vive un algoritmo concreto.

Complementa [[ANTLR]], [[CUP]], [[JFlex]] y [[Jison]]: las notas de esas 4 explican *cómo se usan* esas herramientas en los proyectos de Hades; este grafo permite explorar *cómo están implementadas por dentro*.

## Proyectos indexados (código fuente clonado en `Hades/Fuentes-Compiladores/`)

| Proyecto | Repo | Nodos | Edges | Qué mirar ahí |
|---|---|---|---|---|
| antlr4 | github.com/antlr/antlr4 | 20.799 | 98.573 | Generación de ATN (Augmented Transition Network) y el algoritmo ALL(\*) — la alternativa de ANTLR a las tablas LALR fijas |
| jflex | github.com/jflex-de/jflex | 10.527 | 40.070 | Construcción de AFN→AFD (subset construction) y minimización de autómatas a partir de expresiones regulares — la teoría de [[Autómata finito (AFN y AFD)]] llevada a código real |
| cup-maven-plugin | github.com/vbmacher/cup-maven-plugin | 333 | 455 | Wrapper Maven sobre CUP: fase `generate-sources`, cómo se invoca el generador LALR real por debajo |
| jison | github.com/zaach/jison | 1.214 | 2.436 | Generador LALR(1)/SLR/LR(1) en JS — construcción de tablas de parsing y su ejecución en el runtime generado |

## Cómo consultar el grafo (dentro de una sesión de Claude Code, tras reiniciar)

```
list_projects                                          # confirma que los 4 están indexados
get_graph_schema(project="...")                        # tipos de nodo/edge disponibles
search_graph(project="...", name_pattern=".*NFAState.*")  # ubicar una clase/función por nombre
trace_path(project="...", function_name="X", direction="both")  # quién llama y a quién llama
get_code_snippet(qualified_name="proyecto.paquete.Clase")  # leer el código real
get_architecture(project="...", aspects=["all"])       # panorama de paquetes, capas, hotspots
```

También disponible sin reiniciar, vía CLI directo (no requiere conexión MCP):
```bash
codebase-memory-mcp cli search_graph '{"project": "C-Users-72358-Desktop-Hades-Fuentes-Compiladores-jflex", "name_pattern": ".*NFA.*"}'
```

Nombres de proyecto exactos (los asigna el indexador a partir del path): `C-Users-72358-Desktop-Hades-Fuentes-Compiladores-{antlr4|cup-maven-plugin|jflex|jison}` — usar `list_projects` para confirmarlos si cambia el path.

## Por qué existe esta nota

El vault trae teoría y decisiones de los proyectos del curso; este grafo trae la **implementación de referencia** de las 4 herramientas que generan lexers/parsers para esos proyectos. Sirve para responder preguntas de "¿cómo lo resuelve de verdad JFlex/ANTLR/CUP/Jison?" con evidencia de código en vez de memoria o suposición — relevante para profundizar en [[Autómata finito (AFN y AFD)]], [[Análisis sintáctico ascendente LR]] y el algoritmo ALL(\*) de ANTLR.

## Relacionadas
- [[ANTLR]] · [[CUP]] · [[JFlex]] · [[Jison]]
- [[Autómata finito (AFN y AFD)]]
- [[Análisis sintáctico ascendente LR]]
