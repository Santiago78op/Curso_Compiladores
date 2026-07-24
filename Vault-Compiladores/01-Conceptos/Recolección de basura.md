---
tags: [concepto, runtime, contexto]
aliases: ["garbage collection", "GC", "alcanzabilidad", "conjunto raíz", "conteo de referencias", "marcar y limpiar", "generacional", "quién libera la memoria"]
fuente: "Libro del Dragón §7.4–7.7, pp. 452–499"
fecha: 2026-07-24
---

# Recolección de basura

> **Contexto teórico, no requerido en los proyectos** (como [[Código de tres direcciones]]): pero los 5 proyectos **corren sobre un recolector de basura** —la JVM para los de Java, V8 para [[CompInterpreter]]—, así que esta es la respuesta a *"¿quién libera la memoria de tus entornos?"*: el GC del anfitrión. Un entorno se vuelve **inalcanzable** al retornar la función (salvo que algo vivo lo siga referenciando), y el GC lo reclama.

## Alcanzabilidad y conjunto raíz (§7.4.3)

Un objeto es **basura** cuando ya no es **alcanzable**. El GC parte del **conjunto raíz** (variables globales y las de la pila de activaciones vivas) y sigue los punteros: todo lo que se alcanza está vivo; el resto se recicla. Reformula "liberar memoria" como "¿queda algún camino a este objeto?".

## Familias de algoritmos

- **Conteo de referencias (§7.5).** Cada objeto lleva un contador de cuántos lo referencian; al llegar a 0 se libera de inmediato. Simple e incremental, pero **falla con ciclos**: dos objetos que se apuntan mutuamente pero que nadie más referencia mantienen su contador en 1 y **nunca** se liberan (fuga).
- **Marcar y limpiar / basados en rastreo (§7.6).** Periódicamente: **marcar** todo lo alcanzable desde las raíces, y luego **limpiar** (reclamar) lo no marcado. Sí recupera ciclos, a costa de pausas.
- **Generacional (§7.7.3).** Aprovecha que **los objetos "mueren jóvenes"** —entre el **80 % y el 98 %** de los recién creados se vuelven basura enseguida—: se recolecta con más frecuencia la generación **joven** (barata, mucha basura) y rara vez la vieja. Es lo que usan **la JVM y V8**, así que es literalmente el recolector bajo el que corren los proyectos.

## Contraste para la defensa

Los lenguajes sin GC sufren **fugas de memoria** y **punteros colgantes** (§7.4.5); los proyectos no, precisamente porque delegan la memoria en el GC del anfitrión. La contraparte es que el GC introduce **pausas** no deterministas — irrelevante para un intérprete de curso, pero es la razón por la que un lenguaje de tiempo real evitaría este modelo.

## Relacionadas
- [[Registro de activación y pila de control]]
- [[Entornos y alcance]]
- [[Compilador vs intérprete]]
- [[Cap 7 - Entornos en tiempo de ejecución]]
