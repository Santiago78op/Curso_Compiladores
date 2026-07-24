# Auditoría Fase A — Capítulo 7: Entornos en tiempo de ejecución

**Fuente:** `doc\Boocks\dragon-md\cap07-entornos-ejecucion.md` (Dragón 2ª ed., pp. 427–503) leído completo, contrastado contra: `Cap 7`, `Registro de activación y pila de control`, `Entornos y alcance`, `Tabla de símbolos`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto general

La nota `Registro de activación y pila de control` tiene **el mejor puente teoría→proyecto del vault**: "registro de activación ≈ entorno, pila de control ≈ pila de entornos, la recursión funciona porque cada llamada tiene su propio entorno; los campos de bajo nivel no aplican a un intérprete". Eso es exactamente el marco correcto. Las brechas: la distinción enlace de acceso / enlace de control (que esconde EL bug clásico de intérpretes), el árbol de activación sin desarrollar, y la ausencia total de recolección de basura en el vault (la mitad del capítulo).

## Brechas (ordenadas por impacto)

### B1. Enlace de acceso ≠ enlace de control, y el bug del alcance dinámico accidental (§7.3.5–7.3.6, pp. 445–447) — **ALTA (conexión directa)**
La nota nombra ambos enlaces pero no los distingue operativamente: **enlace de control** = apunta al registro del *llamador* (orden dinámico de llamadas); **enlace de acceso** = apunta al registro del procedimiento donde fui *declarado* (alcance estático). En un intérprete esta distinción es EL bug clásico: si al invocar una función creás su entorno como hijo del **entorno actual (llamador)** en vez del **entorno de declaración** (global, o el de la función externa), obtenés *alcance dinámico accidental* — una variable local del llamador se "filtra" a la función llamada. Pregunta de defensa directa: "¿por qué el entorno de una función cuelga del global y no de quien la llama?"
**Acción:** sección en `Registro de activación…` (o `Entornos y alcance`) con la distinción, el bug, y el fragmento real de los proyectos donde se elige el padre correcto.

### B2. El árbol de activación y sus tres propiedades (§7.2.1–7.2.2, pp. 430–435) — **MEDIA-ALTA**
La nota lo menciona en una línea; el libro da las 3 propiedades que justifican que una PILA baste: (1) las llamadas siguen un recorrido en **preorden** del árbol; (2) los retornos, un **postorden**; (3) las activaciones vivas = el **camino de la raíz al nodo actual**. Ejercicio de examen típico: dado un programa recursivo (quicksort, fibonacci), dibujar el árbol y decir cuántos registros conviven en la pila en el peor momento (ej. 7.2.1–7.2.3).
**Acción:** ampliar la nota con el árbol de activación de fibonacci(5) en Mermaid + la instantánea de la pila.

### B3. Funciones como parámetros llevan su enlace de acceso = closures (§7.3.7, pp. 448–449) — **MEDIA**
Cuando se pasa una función como parámetro, el emisor pasa **⟨función, enlace de acceso⟩** — el par que en lenguajes modernos se llama *closure*. Si CompInterpreter (JS) soporta funciones como valores o callbacks, esta es su teoría; si no, es la respuesta a "¿qué pasaría si tu lenguaje permitiera devolver funciones?".
**Acción:** párrafo en la nota de registro de activación, citando §7.3.7 y la fig. 7.13 (`f: ⟨d, enlace⟩`).

### B4. No existe ninguna nota de recolección de basura (§7.4–7.8, pp. 452–499 — la mitad del capítulo) — **MEDIA (cultura + conexión JVM/V8)**
El vault no menciona GC en ningún lado, y los 5 proyectos CORREN sobre recolectores de basura (JVM para los Java, V8 para CompInterpreter). Respuesta a "¿quién libera la memoria de tus entornos?": el GC del anfitrión — un entorno se vuelve *inalcanzable* al retornar la función (salvo referencias vivas). Contenido mínimo para una nota nueva con encuadre "contexto, no requerido" (como la de tres direcciones): alcanzabilidad y conjunto raíz; conteo de referencias y por qué falla con ciclos (fig. 7.18); marcar-y-limpiar en 4 estados; generacional con el dato "80–98% de los objetos mueren jóvenes" (§7.7.3) — que es lo que usan la JVM y V8.
**Acción:** nota nueva `Recolección de basura` en `01-Conceptos`.

### B5. Menores — **BAJA**
- **Display** (§7.3.8): el arreglo `d[i]` → registro más alto de profundidad i, alternativa a las cadenas de enlaces (Dijkstra). Un párrafo si se quiere completar la teoría de acceso no local.
- **Organización de la memoria** (§7.1, fig. 7.1): código / estática / montículo / pila creciendo en sentidos opuestos — buen diagrama de cultura general para la nota Cap 7.
- **Fugas y punteros colgantes** (§7.4.5): vocabulario útil (los proyectos no los sufren gracias al GC — buen contraste).

## Fuera de alcance (correctamente ausente)

- Secuencias de llamada emisor/receptor y varargs (§7.2.3): back-end/ABI.
- Datos de longitud variable en la pila (§7.2.4), best-fit/coalescencia (§7.4.4), jerarquía de memoria y localidad (§7.4.2–7.4.3): sistemas operativos.
- GC incremental, tren, paralelo, conservador, referencias débiles (§7.7–7.8): especialización.

## Lo que está bien (sin acción)

- El mapeo intérprete de la nota (frame ≈ entorno, con la honestidad de "estado de máquina y temporales no aplican") es el marco exacto que el curso necesita.
- `Tabla de símbolos` y `Entornos y alcance` ya cubren la pila/árbol de entornos con `get` que sube al padre — coherente con §7.2.
- Aliases bien curados (frame, stack frame, activation record).

## Material aprovechable para presentaciones (Fase B)

- **El árbol de activación de quicksort** (fig. 7.4) sincronizado con la pila creciendo/decreciendo (fig. 7.6) — el stepper perfecto para explicar recursión en la página de tabla-simbolos o en la del cap. 7.
- El **bug del alcance dinámico accidental** como demo interactiva: mismo código, dos elecciones de entorno padre, dos resultados.
- Dato wow: "80–98% de los objetos mueren jóvenes" y por eso la JVM/V8 son generacionales.
- El diagrama de memoria (fig. 7.1) con pila y montículo creciendo uno hacia el otro.
