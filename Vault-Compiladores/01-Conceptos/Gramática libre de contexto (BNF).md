---
tags: [concepto, sintactico]
aliases: [GLC, CFG, BNF, gramática, Backus-Naur, producciones, "wcw no es libre de contexto"]
fuente: "Libro del Dragón, caps. 2 y 4 (§4.3.5 límites de las GLC)"
fecha: 2026-07-24
---

# Gramática libre de contexto (BNF)

Notación para describir la sintaxis de un lenguaje. Tiene **4 componentes**: terminales (tokens), no terminales (variables), símbolo inicial y **producciones** `cabeza → cuerpo`.

Ejemplo: `instr → if ( expr ) instr else instr`.

**BNF** (Backus-Naur Form) es el formato exacto del **archivo de gramática** que se entrega en los proyectos (limpio, no copiado del archivo de CUP/Jison, y **no ambiguo por sí solo** — ver la construcción por niveles de [[Ambigüedad, precedencia y asociatividad]]).

## Qué NO puede expresar una GLC (§4.3.5) — respuesta de defensa

Pregunta clásica: *"¿por qué validás variables no declaradas / la aridad de funciones en el análisis semántico y no en la gramática?"*. Porque **no se puede hacer con una gramática**:

- `w c w` (declarar un nombre y luego usarlo idéntico) **no es un lenguaje libre de contexto**.
- `aⁿ bᵐ cⁿ dᵐ` (un parámetro por cada argumento) **tampoco**.

Una GLC puede "contar" y emparejar **dos** cosas (como `aⁿbⁿ` = paréntesis balanceados), pero **no tres ni comparar dos cadenas arbitrarias por igualdad**. Por eso el "declarar antes de usar" y la comprobación de tipos se resuelven en semántica, consultando la [[Tabla de símbolos]] — no en el parser. Es exactamente la frontera que separa el [[Análisis sintáctico ascendente LR|análisis sintáctico]] del [[Comprobación de tipos|semántico]].

## Relacionadas
- [[Derivaciones y árbol de análisis sintáctico]]
- [[Ambigüedad, precedencia y asociatividad]]
- [[Recursividad por la izquierda y factorización]]
- [[Manejo de errores (léxicos, sintácticos, semánticos)]]
