# Auditoría Fase A — Capítulo 5: Traducción orientada por la sintaxis

**Fuente:** `doc\Boocks\dragon-md\cap05-traduccion-sintaxis.md` (Dragón 2ª ed., pp. 303–355) leído completo, contrastado contra: `Cap 5`, `Traducción dirigida por la sintaxis`, `Atributos sintetizados y heredados`, `Árbol de sintaxis abstracta (AST)`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto general

Las 3 notas de concepto capturan bien el núcleo (SDD vs SDT, sintetizado vs heredado, S-atribuida↔LR, L-atribuida↔LL, construcción del AST con `Nodo`/`Hoja` — esta última calca la fig. 5.10 del libro). Lo que falta es el **cómo funciona por debajo**: la implementación de los atributos en la pila del parser (que es la explicación física de `RESULT`/`$$`) y la teoría de los órdenes de evaluación (grafos de dependencias). Este capítulo también formaliza el término "S-atribuida" con el que se defiende la arquitectura de DataForge.

## Hallazgo positivo (para la defensa de DataForge)

§5.4.1 define el **SDT postfijo**: todas las acciones al final de la producción, ejecutadas al reducir — implementable en la pila LR sin construir árbol. Es la formalización exacta de DataForge ("gramática S-atribuida, ejecución directa en las acciones"). La cita completa para la defensa: Dragón §5.4.1–5.4.2, pp. 324–327.

## Brechas (ordenadas por impacto)

### B1. Los atributos viven EN la pila del parser (§5.4.2, pp. 325–327) — **ALTA (conexión directa)**
El libro muestra la mecánica física detrás de `RESULT` (CUP) y `$$`/`$1` (Yacc/Jison): cada registro de la pila LR guarda estado + atributo(s); al reducir `E → E₁ + T`, la acción lee `pila[tope−2].val` (E₁) y `pila[tope].val` (T), y deja el resultado donde quedará E (fig. 5.20). Esto explica: (a) por qué `$1`/`$3` se numeran por posición en el cuerpo, (b) por qué los atributos sintetizados son "gratis" en LALR, y (c) por qué DataForge no necesita ninguna estructura extra para evaluar.
**Acción:** ampliar `Traducción dirigida por la sintaxis` con la fig. 5.20 traducida a términos CUP (`RESULT = a + b` ≡ manipulación de pila).

### B2. Grafo de dependencias, orden topológico y circularidad (§5.2, pp. 310–313) — **MEDIA-ALTA**
Sin nota. Es la teoría de POR QUÉ existen las clases S/L-atribuidas: el grafo de dependencias (flechas entre instancias de atributos del árbol anotado) determina los órdenes de evaluación válidos (topológicos); un ciclo (`A.s = B.i; B.i = A.s+1`) = SDD inevaluable, y decidir circularidad en general es exponencial. S-atribuida y L-atribuida son las subclases que **garantizan** grafo acíclico. Ejercicio de examen típico: dibujar el grafo para un árbol anotado.
**Acción:** sección nueva en `Atributos sintetizados y heredados` con el grafo del ejemplo 5.5 (el de `3*5` con her/sin) en Mermaid.

### B3. Definición precisa de L-atribuida + el límite LR y los marcadores (§5.2.4, §5.5.4, recuadro p. 348) — **MEDIA**
La nota dice "L-atribuida encaja con LL" sin la definición: los heredados solo pueden depender del **padre y de los hermanos a la IZQUIERDA** (de ahí la L). Además el recuadro de p. 348 responde una pregunta fina: NO toda L-atribuida sobre gramática LR se puede evaluar en ascenso (al reducir B en `A→BC` aún no sabés qué producción es). El truco de los **marcadores** (`M→ε` en lugar de cada acción intermedia, §5.5.4) tiene una consecuencia práctica directa: **las acciones a mitad de producción en CUP/Yacc se implementan con marcadores ocultos y pueden introducir conflictos** — por eso la convención de poner las acciones al final.
**Acción:** completar la definición en la nota + un párrafo "acciones a mitad de regla" en la nota CUP.

### B4. El patrón while con etiquetas (ejemplo 5.19, S.siguiente / C.true / C.false) — **MEDIA (teórico)**
El ejemplo canónico de SDD L-atribuida: generar código de flujo de control heredando etiquetas (`C.true`, `C.false`, `S.siguiente`) y sintetizando `S.codigo`. Para tus proyectos-intérprete no aplica directamente (ejecutan el AST, no generan etiquetas), pero es el puente al cap. 6 y tema clásico de examen teórico. Cabe como sección en la nota TDS con la aclaración "en un intérprete esto se sustituye por recorrer el AST".

### B5. Propagación de tipos en declaraciones (ejemplo 5.10, `D → T L`) — **MEDIA-BAJA**
El patrón de pasar el tipo como atributo heredado por la lista de identificadores (`int a,b,c` → `L.her = T.tipo` → `agregarTipo(id.entrada, L.her)`). Los proyectos lo resuelven con una lista sintetizada de ids + una pasada en la acción final — vale documentar la equivalencia de ambos patrones en la Guía correspondiente.

### B6. Menores — **BAJA**
- Estructura de tipos `int[2][3]` → `arreglo(2, arreglo(3, integer))` (§5.3.2): representación de tipos de arreglo como árboles — relevante si algún proyecto tiene arreglos multidimensionales.
- Generación "al instante" con atributo principal (§5.5.2): emitir en vez de concatenar cadenas.
- Composición tipográfica Eqn/TeX (§5.4.5): anécdota — este libro se compuso con TeX, cuya teoría de cuadros es una SDD L-atribuida.

## Lo que está bien (sin acción)

- `Traducción dirigida por la sintaxis`: el mapeo SDD/SDT → `RESULT`/`$$` es correcto y es el corazón del capítulo.
- `Atributos sintetizados y heredados`: definiciones fieles a §5.1.1, mapeo S↔LR / L↔LL correcto.
- `AST`: la construcción `Nodo`/`Hoja` calca la fig. 5.10; la conexión con los reportes de los proyectos es valiosa.
- Citas correctas (cap. 5, §5.1).

## Material aprovechable para presentaciones (Fase B)

- **Árbol anotado con evaluación paso a paso** (fig. 5.3: `3*5+4 n` con val subiendo) — stepper natural para la página de semántica/TDS.
- El **contraste her/sin** de la fig. 5.5 (`3*5` en gramática LL): cómo el operando izquierdo "viaja" por atributos heredados cuando el árbol no coincide con la sintaxis abstracta.
- La pila con atributos (fig. 5.19–5.20) como animación: qué pasa físicamente al reducir con `RESULT`.
- Anécdota TeX/Eqn (p. 331): la composición tipográfica de fórmulas es una SDD — los compiladores están en todos lados.
