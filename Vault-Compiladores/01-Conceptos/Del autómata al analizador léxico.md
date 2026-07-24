---
tags: [concepto, lexico]
aliases: ["longest match", "match más largo", "retroceso al último estado de aceptación", "prioridad de la primera regla", "tokenizar vs reconocer"]
fuente: "Libro del Dragón §3.8.1–3.8.3, pp. 166–171"
fecha: 2026-07-24
---

# Del autómata al analizador léxico

Un autómata ([[Autómata finito (AFN y AFD)|AFN/AFD]]) por sí solo dice **sí/no** a *una* cadena. Un analizador léxico hace algo más: **corta lexemas** de un flujo continuo y decide **qué token** es cada uno. Estos son los tres mecanismos que cierran ese salto (§3.8), y son la teoría detrás de las "recetas" que aparecen en [[JFlex]].

## 1. Combinar los autómatas de todos los patrones (§3.8.1, fig. 3.50/3.52)

Se toman los AFN de cada patrón (uno por token) y se unen bajo un **nuevo estado inicial** con transiciones **ε** hacia cada uno. El autómata combinado reconoce, a la vez, *cualquiera* de los patrones. Los estados de aceptación quedan **etiquetados** con el patrón al que pertenecen.

## 2. Longest match: avanzar y retroceder (§3.8.2)

El lexer **avanza** consumiendo entrada mientras el autómata tenga transiciones. Cuando se queda **sin transición** (conjunto de estados vacío), **retrocede** hasta el **último conjunto de estados que incluía un estado de aceptación**: ese prefijo es el **lexema más largo** posible. Así se implementa el *longest match* — la regla "el lexema más largo gana" no es una convención arbitraria, es la consecuencia directa de este retroceso.

## 3. Desempate por orden: gana el patrón listado primero (§3.8.2)

Si el conjunto de aceptación al que se retrocede contiene **varios** estados de aceptación (el lexema calza más de un patrón de **igual longitud**), se elige el patrón que aparece **primero** en el programa Lex/`.flex`. Por eso las palabras reservadas se declaran **antes** que `{Id}`: `if` calza tanto la reservada como el identificador (misma longitud), y el orden decide. *(Ojo: esto solo desempata igual longitud; entre longitudes distintas manda el longest match del punto 2 — ver [[Ambigüedad, precedencia y asociatividad]] para la misma distinción en sintaxis.)*

## Ejemplo 3.27 — entrada `aaba`, patrones `a`, `abb`, `a*b+`

Partiendo de `ε-cerradura(0) = {0,1,3,7}` y avanzando: tras la 4ª letra el conjunto de estados queda **vacío** (no hay salida sobre `a` desde el estado 8). Se retrocede:
- tras `a` → conjunto con el estado 2 ⇒ calza el patrón `a`;
- tras `aab` → conjunto con el estado 8 ⇒ calza `a*b+`, y es el **prefijo más largo** que llega a aceptación.

Se elige **`aab`** como lexema y se ejecuta la acción del patrón `a*b+`. (El `a` inicial, aunque también aceptaba, pierde por no ser el match más largo.)

## Relacionadas
- [[Autómata finito (AFN y AFD)]]
- [[Construcción de subconjuntos]]
- [[JFlex]]
- [[Cap 3 - Análisis léxico]]
