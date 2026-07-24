# Auditoría Fase A — Capítulo 1: Introducción

**Fuente:** `doc\Boocks\dragon-md\cap01-introduccion.md` (Dragón 2ª ed., pp. 1–38) leído completo, contrastado contra las notas del vault: `Cap 1 - Introducción`, `Fases de un compilador`, `Compilador vs intérprete`, `Token, lexema y patrón`, `Tabla de símbolos`, `Entornos y alcance`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto general

El vault cubre bien el **§1.1–§1.2** (procesadores de lenguaje, fases, front/back-end) a nivel conceptual. Las brechas fuertes están en **§1.6 (Fundamentos de los lenguajes de programación)**, que es justamente la sección más relevante para *intérpretes*: ahí el libro define formalmente conceptos que los 5 proyectos implementan a diario (entorno, alcance, paso de parámetros) y el vault los cubre a medias o nada.

## Brechas (ordenadas por impacto)

### B1. Entorno vs estado — la doble asignación (§1.6.2, fig. 1.8, p. 26) — **ALTA**
El libro define: **entorno** = asignación nombres → ubicaciones (l-values); **estado** = asignación ubicaciones → valores (r-values). Es la teoría que explica por qué la clase central de los 5 proyectos se llama `Entorno`. Ninguna nota del vault la recoge — `Entornos y alcance` salta directo a scope y cita solo cap. 7.3.
**Acción:** ampliar `Entornos y alcance` (o nota nueva) con la figura 1.8 en Mermaid y la conexión l-value/r-value → asignación en los intérpretes.

### B2. Paso de parámetros (§1.6.6, pp. 33–34) + aliasing (§1.6.7, p. 35) — **ALTA**
Llamada por valor / por referencia / por nombre, y el matiz clave: Java pasa **referencias por valor** ("llamada por referencia efectiva" para objetos). No existe NINGUNA nota del vault sobre esto, y CompScript/CompInterpreter/VLangCherry implementan funciones con parámetros — cuando el usuario defienda "¿cómo pasa parámetros tu lenguaje?", esta es la teoría.
**Acción:** nota nueva `Paso de parámetros` con los 3 mecanismos + aliasing (ejemplo 1.9: `q(a,a)`).

### B3. Regla de alcance por bloques (§1.6.3, ejemplo 1.6, pp. 28–31) — **MEDIA**
El vault enuncia estático vs dinámico pero no la regla operativa: "el uso de x se refiere a la declaración de x en el bloque circundante B_i más interno que declare x", ni un ejemplo trabajado de shadowing (el B1–B4 del libro). Es el algoritmo que implementa `Entorno.buscar()` subiendo por el padre.
**Acción:** agregar la regla + un ejemplo trabajado con la pila de entornos de los proyectos.
**Corrección de cita:** la nota `Entornos y alcance` cita "cap. 7.3"; el contenido de scope estático/dinámico está en §1.6.3–1.6.5 (cap. 7 cubre la *implementación* runtime).

### B4. El ejemplo guía completo de la fig. 1.7 (pp. 5–7) — **MEDIA (didáctica)**
`posicion = inicial + velocidad * 60` atravesando las 6 fases con salidas concretas: tokens (1.2) → árbol → `inttofloat` → tres direcciones (1.3) → optimizado (1.4) → máquina (1.5). El vault lo menciona en una línea; ninguna nota lo traza. Es LA demo canónica para presentaciones (patrón stepper).
**Acción:** para Fase B — verificar si `presentacion-dataforge/fases.html` ya lo anima; si no, es la mejora #1 de esa página.

### B5. Fases vs pasadas (§1.2.8, p. 11) — **BAJA**
Distinción clásica de examen (fase = organización lógica; pasada = lee entrada y escribe salida, agrupa fases). No está en el vault.
**Acción:** un párrafo en `Fases de un compilador`.

### B6. Alcance dinámico: los 2 ejemplos reales (§1.6.5, pp. 31–33) — **BAJA**
La nota lo despacha como "poco común", pero el libro da 2 casos vigentes: macros de C y **despacho de métodos virtuales** (ejemplo 1.8: `x.m()` se resuelve en ejecución). El segundo es relevante si algún proyecto tiene OO.

## Lo que está bien (sin acción)

- `Compilador vs intérprete`: fiel al §1.1, incluye el híbrido Java/JIT correctamente.
- `Fases de un compilador`: el Mermaid refleja la fig. 1.6, incluida la tabla de símbolos como estructura transversal (no fase).
- `Token, lexema y patrón`: correcto (el detalle fino viene en cap. 3).
- §1.3–§1.5 (historia, aplicaciones, RISC/CISC): cultura general, no requerido por OLC1 — razonable que el vault no lo indexe.

## Material aprovechable para presentaciones (Fase B)

- Ejercicios del libro reutilizables como auto-evaluación: 1.1.1–1.1.5 (compilador vs intérprete), 1.3.1 (clasificar lenguajes), 1.6.1–1.6.4 (alcance por bloques).
- El recuadro "Nombres, identificadores y variables" (p. 28) y "Declaraciones y definiciones" (p. 32): distinciones finas, buenas para slides.
