---
tags: [proyecto, lexico, sintactico, semantico, web, olc2]
fuente: "Enunciado OLC2_P1_EV2025 (Proyecto 1, 1er Sem 2025)"
fecha: 2026-07-21
---

# VLangCherry

**Objetivo:** intérprete para **V-lang Cherry**, un lenguaje con sintaxis inspirada en Go (subconjunto reducido, pensado para correr rápido con pocos recursos). Extensión `.vch`, *case sensitive*. A diferencia de los 4 proyectos de OLC1 del workspace, este es el **Proyecto 1 de OLC2** (curso distinto) y es **grupal, de 3 integrantes** (secciones 10.4/12.3 del enunciado).

**Etapa del compilador:** front-end completo ([[ANTLR]]) + AST propio + intérprete semántico de dos pasadas, con **arquitectura cliente-servidor** vía REST (Go + React), igual patrón que [[CompInterpreter]] pero con backend en Go en vez de Node.

**Tecnologías:** [[ANTLR]] (gramática `.g4`, target Go) para léxico/sintáctico, **Go** para todo el intérprete (obligatorio por el enunciado, sección 2.3), servidor REST con `net/http` de la librería estándar (sin frameworks), cliente **React + Vite**, AST interactivo con [[vis-network]]. Restricción real del enunciado: ejecución nativa en **Linux** (10.2) — se compila cruzado desde Windows con `GOOS=linux GOARCH=amd64 go build`.

**Rasgos:** tipos primitivos `int, float64, string, bool, rune`; slices (`[]T`, multidimensionales); structs con métodos por valor/puntero (`func (p Persona) Saludar() string` / `func (p *Persona) ...`); funciones con retorno único; sentencias `if/else if/else`, `switch/case` **sin fall-through**, `for` en sus 3 formas (condición, clásico, `for i, v in slice`); `break/continue/return`; funciones nativas `print, println, len, append, indexOf, join, Atoi, parseFloat, typeOf` y conversión de tipo estilo Go (`int(x)`, `float64(x)`, ...). Reportes: errores (léxico/sintáctico/semántico, recolectados TODOS sin abortar), tabla de símbolos, AST (grafo).

**Conceptos aplicados:** [[Árbol de sintaxis abstracta (AST)]] · [[Comprobación de tipos]] · [[Conversión de tipos (coerción y cast)]] · [[Tabla de símbolos]] · [[Entornos y alcance]] · [[Manejo de errores (léxicos, sintácticos, semánticos)]] · [[Registro de activación y pila de control]]

## Guía de elaboración
- [[Guía VLangCherry]]

## Relacionadas
- [[CompInterpreter]]
- [[ANTLR]]
- [[Cap 4 - Análisis sintáctico]]
- [[Cap 7 - Entornos en tiempo de ejecución]]
