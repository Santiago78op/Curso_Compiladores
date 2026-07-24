# Auditoría Fase A — Capítulo 2: Un traductor simple orientado a la sintaxis

**Fuente:** `doc\Boocks\dragon-md\cap02-traductor-simple.md` (Dragón 2ª ed., pp. 39–107) leído completo, contrastado contra: `Cap 2`, `Gramática libre de contexto (BNF)`, `Derivaciones y árbol…`, `Ambigüedad, precedencia y asociatividad`, `Traducción dirigida por la sintaxis`, `Atributos sintetizados y heredados`, `Árbol de sintaxis abstracta (AST)`, `Recursividad por la izquierda y factorización`, `Tabla de símbolos`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto general

Este capítulo es el mejor cubierto por el vault hasta ahora: las 8 notas tocan casi todas sus secciones y las citas son mayormente correctas. Las brechas están en (1) la técnica de **precedencia codificada en la gramática misma** — crítica porque la BNF entregable no puede usar `%left` —, (2) el patrón de **tablas de símbolos encadenadas** que es literalmente la clase `Entorno` de los proyectos, y (3) piezas didácticas sin nota (recorridos, postfija, léxico a mano).

## Hallazgo positivo (para la defensa de DataForge)

§2.8.1 (p. 92) valida textualmente la decisión de DataForge de **no construir AST**: "es común que los compiladores emitan el código… mientras el analizador sintáctico 'avanza' …, sin construir en realidad la estructura de datos tipo árbol completa". DataForge = el patrón del cap. 2 (traducción S-atribuida ejecutada en las acciones); CompScript = el patrón del §2.8.2 (construir el AST con clases `Nodo`/`If`/`While` y recorrerlo). **Cita ideal para justificar la arquitectura en ambas defensas.**

## Brechas (ordenadas por impacto)

### B1. Precedencia por niveles de gramática (§2.2.6 + recuadro "Generalización", pp. 49–50) — **ALTA**
La nota `Ambigüedad…` solo describe la solución con **declaraciones** (`precedence left` en CUP, `%left` en Jison). Falta la solución **estructural**: una gramática con *n* niveles de precedencia usa *n+1* no terminales (`expr → expr + term | term`, `term → term * factor | factor`, `factor → dígito | ( expr )`), donde el nivel más profundo "no puede separarse". **La BNF entregable de cada proyecto debe ser no ambigua por sí sola** (no lleva tablas de precedencia): esta técnica es la única forma de escribirla bien.
**Acción:** ampliar `Ambigüedad, precedencia y asociatividad` con la construcción por niveles + regla n+1; ajustar cita (hoy dice solo "cap. 4.8"; esto es §2.2.5–2.2.6).

### B2. Tablas de símbolos encadenadas = la clase `Entorno` (§2.7, pp. 85–91) — **ALTA**
La nota `Tabla de símbolos` menciona la "pila/árbol de entornos" pero no el patrón canónico que TODOS los proyectos implementan: clase `Ent` con hash table propia + puntero `ant` al entorno padre; `get` sube la cadena hasta null; entrar a bloque = `sup = new Ent(sup)`, salir = restaurar `guardado` (disciplina de pila). Tampoco está el ejemplo 2.14 (`{ int x; char y; { bool y; x; y; } x; y; }` → shadowing trabajado) ni el recuadro "¿Quién crea las entradas?" (el parser, no el lexer — pregunta de examen).
**Acción:** ampliar `Tabla de símbolos` con el patrón `Ent`/`ant`/`get`-en-cadena y el ejemplo 2.14; conecta directo con la brecha B3 del cap. 1 (regla del bloque anidado más cercano).

### B3. Recorridos de árboles: preorden / postorden / DFS (§2.3.4, pp. 56–58) — **MEDIA**
No existe nota. Es la teoría de POR QUÉ un intérprete funciona: **ejecutar el AST = recorrido en postorden** (evaluar hijos, luego aplicar el nodo); los atributos sintetizados se evalúan en cualquier recorrido de abajo hacia arriba. El `visit()`/`ejecutar()` de CompScript/CompInterpreter es exactamente el `visitar(N)` de la fig. 2.11.
**Acción:** nota nueva `Recorridos de árboles (preorden y postorden)` o sección en la nota AST.

### B4. Comprobación estática: l-value/r-value y sobrecarga (§2.8.3, pp. 97–99) — **MEDIA**
El capítulo define: comprobación sintáctica más allá de la gramática (ej. `break` solo dentro de ciclo — ¡los proyectos lo validan en semántica!), l-value vs r-value (conecta con brecha B1 del cap. 1), coerciones y **sobrecarga** (`+` = suma o concatenación según tipos — regla que DataForge/CompScript implementan literalmente en `Operaciones`).
**Acción:** verificar en la auditoría del cap. 6 si `Comprobación de tipos` y `Conversión de tipos` cubren sobrecarga; si no, agregarla ahí citando §2.8.3.

### B5. Transformar gramáticas sin romper las acciones (§2.5.2 + §2.5.4, pp. 70–74) — **MEDIA-BAJA**
La nota `Recursividad por la izquierda…` da la transformación pero no las dos sutilezas del libro: (a) al eliminar recursividad izquierda, las acciones semánticas **viajan como si fueran terminales** — moverlas al final produce traducciones incorrectas (9−5+2 → 952+− ¡erróneo!); (b) la recursividad **por la cola** se convierte en iteración (`while`), que es como el parser real integra `resto` en `expr`. Nota: para CUP/Jison (LALR) la recursividad izquierda NO es problema — es la forma *preferida* — esta técnica aplica solo si escribieras un parser descendente a mano.
**Acción:** ampliar la nota con ambas sutilezas + la aclaración LALR-vs-LL.

### B6. Notación postfija (§2.3.1, pp. 53–54) — **BAJA**
Definición inductiva + algoritmo de evaluación (buscar el primer operador, tomar operandos a su izquierda). Clásico de examen, sin nota. Cabría como sección corta en `Traducción dirigida por la sintaxis`.

### B7. Analizador léxico a mano (§2.6, pp. 76–84) — **BAJA**
El capítulo construye a mano lo que JFlex genera: variable `vistazo` (lookahead de 1 carácter), `v = v*10 + dígito` para números, y **palabras reservadas via tabla sembrada** (buscar el lexema antes de declararlo id — la teoría detrás de "reservadas ANTES de {Id}" en el `.flex` de DataForge).
**Acción:** párrafo en la nota `JFlex` citando §2.6.4 como fundamento de esa convención.

## Lo que está bien (sin acción)

- `GLC (BNF)`: los 4 componentes correctos, y el recordatorio de que la BNF entregable no es copia del `.cup` — alineado con §2.2.1.
- `Derivaciones y árbol…`: definición formal fiel (§2.2.2–2.2.3), buen Mermaid, la conexión AST-sin-andamiaje es exactamente §2.5.1.
- `AST`: fiel a §2.5.1/§2.8.2, distinción concreto/abstracto correcta.
- `Traducción dirigida por la sintaxis` y `Atributos…`: correctas para el nivel cap. 2 (el detalle fino es cap. 5 — se auditará ahí).

## Material aprovechable para presentaciones (Fase B)

- **Números romanos** (ejercicios 2.2.6 y 2.3.3–2.3.4): gramática + traducción bidireccional — ejercicio memorable para la página de gramáticas.
- **Fig. 2.46** (p. 106): la misma instrucción traducida a AST *y* a tres direcciones — buen slide comparativo de IRs.
- **Ejemplo 2.14** (shadowing con tipos): demo stepper natural para la página tabla-simbolos.
- La secuencia postfija 9−5+2 → 95−2+ animada paso a paso (fig. 2.14: acciones print colgando del árbol).
