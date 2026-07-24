# Auditoría Fase A — Capítulo 6: Generación de código intermedio

**Fuente:** `doc\Boocks\dragon-md\cap06-codigo-intermedio.md` (Dragón 2ª ed., pp. 357–426) leído completo, contrastado contra: `Cap 6`, `Comprobación de tipos`, `Conversión de tipos (coerción y cast)`, `Flujo de control y switch`, `Código de tres direcciones`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto general

Las 4 notas están **bien enfocadas y con el encuadre correcto**: la nota de código de tres direcciones declara honestamente "contexto teórico, no requerido en los proyectos", y `Flujo de control` da el contraste clave ("en un intérprete no se generan etiquetas: se recorre el AST"). Las citas son todas correctas. Las brechas son las piezas semánticas que los proyectos SÍ implementan y que las notas mencionan sin el algoritmo: sobrecarga de operadores y la mecánica de coerción en operaciones binarias.

**Verificación pendiente del cap. 2 (B4): CONFIRMADA** — la sobrecarga NO está cubierta en ninguna nota del vault.

## Brechas (ordenadas por impacto)

### B1. Sobrecarga de operadores y funciones (§6.5.3, pp. 390–391 + §2.8.3) — **ALTA (conexión directa)**
El `+` que significa suma entre números y concatenación entre cadenas — exactamente la semántica que `Operaciones` (DataForge) y las tablas de compatibilidad (CompScript/CompInterpreter) implementan — es *sobrecarga*, y el libro da la regla de resolución: elegir el significado por la **firma** (operador + tipos de los operandos), implementable eficientemente con el método del número de valor. La nota `Comprobación de tipos` menciona las tablas de compatibilidad pero nunca usa ni explica el término "sobrecarga", que es como el tribunal va a preguntar.
**Acción:** sección "Sobrecarga" en `Comprobación de tipos` citando §6.5.3 (y §2.8.3): definición, resolución por firma, y el mapeo explícito a las tablas de compatibilidad de los proyectos.

### B2. La mecánica de coerción binaria: `max(t₁,t₂)` + `ampliar` (§6.5.2, figs. 6.25–6.27, pp. 388–390) — **MEDIA-ALTA (conexión directa)**
La nota `Conversión de tipos` da los conceptos (widening/narrowing, coerción/cast) pero no el **algoritmo** que todo intérprete ejecuta al evaluar `int + double`: (1) `max(t₁,t₂)` sube por la jerarquía de ampliación (fig. 6.25: char→short→int→long→float→double) para hallar el tipo del resultado; (2) `ampliar(valor, t, w)` promueve el operando menor. Es el esqueleto exacto del método de dominancia de tipos en `Operaciones`.
**Acción:** agregar a la nota el pseudocódigo de `max`/`ampliar` y la jerarquía de Java como diagrama Mermaid.

### B3. break / continue / return en el intérprete vs el compilador (§6.7.4, pp. 416–417) — **MEDIA**
El libro traduce `break` con listas de saltos sin destino (siglista de la construcción envolvente, backpatch al conocerla). Los proyectos lo implementan de otra forma (señales de control al recorrer el AST: el visitor devuelve/lanza un marcador Break/Continue/Return que los nodos de ciclo interceptan) — y esa dualidad no está documentada en ningún lado. Es pregunta doble de defensa: "¿cómo funciona tu break?" y "¿cómo lo haría un compilador?".
**Acción:** sección en `Flujo de control y switch` con ambos mecanismos, citando §6.7.4 para el lado compilador.

### B4. Estrategias de implementación del switch (§6.8.1, pp. 419–420) — **MEDIA-BAJA**
La nota dice "bifurcación de n vías" pero no las 3 estrategias del libro: saltos condicionales secuenciales (≤10 casos), tabla hash de valores→etiquetas, o arreglo de baldes si los valores caen en un rango denso. Aun en un intérprete la idea aplica (un `match` grande puede usar un HashMap en vez de if-else en cadena).
**Acción:** párrafo en la nota citando §6.8.1.

### B5. Cultura teórica de examen: GDA, número de valor y SSA (§6.1, §6.2.4) — **BAJA**
- **GDA** (DAG): árbol sintáctico donde las subexpresiones comunes comparten nodo (`a+a*(b-c)+(b-c)*d` → un solo nodo para `b-c`), construido con el **método del número de valor** (hash de firmas ⟨op,i,d⟩).
- **SSA** (asignación individual estática): cada variable se asigna una sola vez; la función **φ** combina definiciones de caminos distintos. Dato moderno: es la IR de LLVM.
**Acción:** dos párrafos cortos en `Código de tres direcciones` — mantienen su encuadre de "contexto teórico".

### B6. Expresiones de tipos y equivalencia (§6.3.1–6.3.2, pp. 371–373) — **BAJA**
`int[2][3]` = `arreglo(2, arreglo(3, integer))`; equivalencia estructural vs por nombre. Ya está parcialmente en la auditoría del cap. 5 (B6); relevante solo si algún proyecto tiene arreglos multidimensionales o structs con tipos nombrados.

## Fuera de alcance (correctamente ausente del vault)

- Inferencia de tipos con unificación y polimorfismo paramétrico (§6.5.4–6.5.5, ML): teoría de lenguajes funcionales, no OLC1.
- Direccionamiento de arreglos con fórmulas base+i×w y orden fila/columna (§6.4.3): back-end; los intérpretes usan listas/arrays del lenguaje anfitrión.
- Distribución de almacenamiento, alineación y relleno (§6.3.4).

## Lo que está bien (sin acción)

- `Código de tres direcciones`: el encuadre "no requerido, pero cae en exámenes" es exactamente el correcto; la tabla cuádruplos/tripletas/indirectas es fiel a §6.2.2–6.2.3; el resumen de backpatching es suficiente para su propósito.
- `Flujo de control y switch`: el contraste compilador (etiquetas/gotos) vs intérprete (recorrer el AST) es el marco didáctico correcto.
- `Conversión de tipos`: el mapeo `CAST(exp AS tipo)` = conversión explícita es exacto.
- `Comprobación de tipos`: síntesis vs inferencia fiel a §6.5.1; el tratamiento de instrucciones como funciones (`if: bool×void→void`) viene directo del libro.
- Citas correctas en las 4 notas.

## Material aprovechable para presentaciones (Fase B)

- **Backpatching paso a paso** (ejemplo 6.24): las instrucciones 100–105 de `x<100 || x>200 && x!=y` llenándose con backpatch(102,104) → backpatch(101,102) — stepper perfecto (aunque sea teórico para los proyectos).
- **El GDA de `a+a*(b-c)+(b-c)*d`** (fig. 6.3): visual inmediato de subexpresiones comunes.
- **Código de corto circuito** (fig. 6.34): la misma expresión booleana como cascada de saltos — contraste con cómo el intérprete simplemente no evalúa el operando derecho.
- **La jerarquía de coerciones de Java** (fig. 6.25) como diagrama: widening arriba, narrowing abajo.
