---
tags: [concepto, semantico, runtime]
aliases: ["llamada por valor", "llamada por referencia", "llamada por nombre", "paso de argumentos", "aliasing", "uso de alias"]
fuente: "Libro del Dragón §1.6.6 (mecanismos) y §1.6.7 (alias), pp. 33–35"
fecha: 2026-07-24
---

# Paso de parámetros

Cómo llega un **argumento (parámetro actual)** a su **parámetro formal** cuando se invoca un procedimiento. El libro (§1.6.6, pp. 33–34) describe tres mecanismos:

## Los tres mecanismos (§1.6.6)

- **Llamada por valor.** El parámetro actual se **evalúa** (si es expresión) o se **copia** (si es variable), y el formal recibe una copia independiente. Mutar el formal dentro del procedimiento **no** afecta al llamador. Es la opción de C, C++ (por defecto) y la mayoría de los lenguajes.
- **Llamada por referencia.** Se pasa la **dirección** del parámetro actual; el formal es un alias de la variable del llamador, así que asignarlo **sí** se ve afuera. Son los parámetros `ref` de C++.
- **Llamada por nombre.** El cuerpo se ejecuta como si el parámetro actual se **sustituyera literalmente** por el formal (estilo macro). De Algol 60; produce comportamientos no intuitivos con expresiones y hoy no se usa.

### El matiz de Java (y de los intérpretes del curso)

Java usa **exclusivamente llamada por valor** — pero las **referencias a objetos se pasan por valor**. El efecto: se comporta *como si* fuera por referencia para cualquier objeto (mutar sus campos se ve afuera), aunque reasignar el parámetro formal a otro objeto **no** cambia la variable del llamador. Es el matiz de defensa "¿cómo pasa parámetros tu lenguaje?": *los primitivos por valor; los objetos/estructuras, la referencia por valor*.

## Uso de alias (§1.6.7, ejemplo 1.9, p. 35)

Consecuencia del paso por referencia (o de su simulación en Java): dos parámetros formales pueden terminar apuntando a la **misma ubicación** — son **alias** uno del otro.

> **Ejemplo 1.9:** un procedimiento `p` tiene un arreglo `a` y llama a `q(x, y)` con `q(a, a)`. Como los nombres de arreglo son referencias, `x` y `y` quedan como alias. Una asignación `x[10] = 2` dentro de `q` hace que `y[10]` también valga 2.

El aliasing es **esencial para optimizar** (cap. 9): solo se puede sustituir `a = x+3` por `a = 5` si se está seguro de que ninguna otra variable es alias de `x`.

## Conexión con los proyectos

Los proyectos implementan funciones/métodos con parámetros, y cada uno tomó (y ahora documenta) una decisión de semántica de referencia:

- **[[VLangCherry]]:** un receptor de método **por valor** clona el struct (`ClonarPorValor`) — copia independiente; un receptor **por puntero** comparte el `StructVal`. Es la distinción valor/referencia de §1.6.6 hecha explícita.
- **[[CompScript]]:** declarar un vector con literal **copia**, pero una asignación variable-a-variable **comparte** la referencia (aliasing estilo Java, ejemplo 1.9). Decisión documentada en su Manual Técnico.
- **[[CompInterpreter]]:** argumentos evaluados en el entorno del llamador y pasados por valor; los vectores son objetos, así que aplica el matiz de Java.

Los argumentos se depositan en el [[Registro de activación y pila de control|registro de activación]] de la llamada, dentro de la pila de [[Entornos y alcance|entornos]].

## Relacionadas
- [[Entornos y alcance]]
- [[Registro de activación y pila de control]]
- [[Comprobación de tipos]]
- [[Cap 1 - Introducción]]
