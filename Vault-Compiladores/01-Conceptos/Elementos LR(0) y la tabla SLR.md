---
tags: [concepto, sintactico, algoritmo]
aliases: ["elementos LR(0)", "ítem LR(0)", "CERRADURA", "ir_A", "goto", "autómata LR(0)", "tabla SLR", "colección canónica", "ACCION ir_A"]
fuente: "Libro del Dragón §4.6 (SLR), pp. 241–258"
fecha: 2026-07-24
---

# Elementos LR(0) y la tabla SLR

El motor **por dentro** de un parser ascendente ([[Análisis sintáctico ascendente LR|LR]]): cómo se pasa de una gramática a la **tabla ACCION/ir_A** que decide *shift* vs *reduce*. Es lo que [[CUP]]/Jison generan, y lo que hay que leer cuando reportan un conflicto.

Gramática canónica (de expresiones), **aumentada** con `E' → E`:
```
E → E + T | T      T → T * F | F      F → ( E ) | id
```

## Elementos LR(0) (§4.6.2)

Un **elemento** es una producción con un **punto** que marca lo ya visto: `E → E · + T` significa "ya reconocí `E`, espero ver `+ T`". Un elemento con el punto al final (`E → E + T ·`) señala un **mango** listo para reducir.

**CERRADURA(I):** si `A → α · B β` está en `I` y `B → γ` es producción, agrega `B → · γ` (repetir hasta punto fijo). El estado inicial es `CERRADURA({E' → · E})`:
```
E' → · E     E → · E + T     E → · T
T → · T * F  T → · F         F → · ( E )   F → · id
```

**ir_A(I, X):** desde el conjunto `I`, avanza el punto sobre el símbolo `X` en todos los elementos que lo tienen justo después del punto, y cierra el resultado. Es la transición del autómata.

**Colección canónica y autómata LR(0):** partir de `CERRADURA({E'→·E})` y aplicar `ir_A` sobre cada símbolo hasta no generar conjuntos nuevos. Cada conjunto = un **estado**; las `ir_A` = las transiciones (fig. 4.31).

> **Puente con el léxico (recuadro p. 257):** los elementos funcionan como los estados de un AFN, y construir la colección canónica **es** la [[Construcción de subconjuntos|construcción de subconjuntos]] aplicada a los ítems. La misma idea del cap. 3, un nivel más arriba.

## La tabla SLR (Algoritmo 4.46)

Para cada estado `i` (conjunto de elementos `Iᵢ`):
- **shift:** si `[A → α · a β] ∈ Iᵢ` con `a` terminal e `ir_A(Iᵢ, a) = Iⱼ` → `ACCION[i, a] = desplazar j`.
- **reduce:** si `[A → α ·] ∈ Iᵢ` (y `A ≠ E'`) → `ACCION[i, a] = reducir A → α` para **cada `a ∈ SIGUIENTE(A)`** (ver [[FIRST y FOLLOW]]). *(Usar SIGUIENTE es lo que hace "simple" al SLR — y también lo que lo hace fallar en algunas gramáticas donde LALR sí funciona.)*
- **aceptar:** si `[E' → E ·] ∈ Iᵢ` → `ACCION[i, $] = aceptar`.
- **ir_A:** si `ir_A(Iᵢ, A) = Iⱼ` para un no terminal `A` → `ir_A[i, A] = j`.

Una celda con dos acciones = **conflicto** shift/reduce o reduce/reduce (la gramática no es SLR).

## Traza de `id * id + id` (pila de símbolos ↔ acción)

La pila real guarda **estados** del autómata; abajo se muestran los símbolos gramaticales que esos estados representan, que es lo que importa para seguir el reduce:

| Pila (símbolos) | Entrada | Acción |
|---|---|---|
| (vacía) | id * id + id $ | desplazar id |
| id | * id + id $ | reducir F → id |
| F | * id + id $ | reducir T → F |
| T | * id + id $ | desplazar * |
| T * | id + id $ | desplazar id |
| T * id | + id $ | reducir F → id |
| T * F | + id $ | reducir T → T * F |
| T | + id $ | reducir E → T |
| E | + id $ | desplazar + |
| E + | id $ | desplazar id |
| E + id | $ | reducir F → id |
| E + F | $ | reducir T → F |
| E + T | $ | reducir E → E + T |
| E | $ | aceptar |

Observá que `T * F` se reduce (mango) **antes** de tocar el `+`: así la tabla codifica que `*` liga más fuerte que `+`, sin declaraciones de precedencia.

## Relacionadas
- [[Análisis sintáctico ascendente LR]]
- [[Conflictos shift-reduce y reduce-reduce]]
- [[FIRST y FOLLOW]]
- [[Construcción de subconjuntos]]
- [[CUP]]
- [[Cap 4 - Análisis sintáctico]]
