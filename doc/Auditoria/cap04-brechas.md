# Auditoría Fase A — Capítulo 4: Análisis sintáctico

**Fuente:** `doc\Boocks\dragon-md\cap04-analisis-sintactico.md` (Dragón 2ª ed., pp. 191–302) leído completo, contrastado contra: `Cap 4`, `GLC (BNF)`, `Derivaciones…`, `Ambigüedad, precedencia y asociatividad`, `Recursividad por la izquierda y factorización`, `FIRST y FOLLOW`, `Análisis sintáctico descendente LL(1)`, `Análisis sintáctico ascendente LR`, `Conflictos shift-reduce y reduce-reduce`, `Manejo de errores`, `CUP`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto general

Es el capítulo más largo y denso del libro (112 páginas) y el corazón teórico de CUP/Jison — y el vault lo cubre con notas correctas pero **desproporcionadamente delgadas**: la nota `Análisis sintáctico ascendente LR` resume en 10 líneas lo que el libro desarrolla en 60 páginas (§4.5–4.8). El patrón de brecha se repite: las notas dan *resultados* y *clasificaciones*, no los *algoritmos ejecutables a mano* que piden los exámenes. Además hay dos conexiones de oro con los proyectos que ninguna nota registra (B3 y B6).

## Brechas (ordenadas por impacto)

### B1. El motor LR por dentro: elementos, CERRADURA, ir_A, autómata y tabla (§4.6, pp. 241–258) — **ALTA**
La nota LR explica shift/reduce y el mango, pero nada del mecanismo: **elementos LR(0)** (producción con punto: `E → E·+ T`), `CERRADURA(I)`, `ir_A(I,X)`, la colección canónica, el **autómata LR(0)** (fig. 4.31), la tabla ACCION/ir_A (fig. 4.37, con s5/r2/acc) y la traza con pila de estados (fig. 4.38). Sin esto no se puede: (a) construir una tabla SLR a mano en un examen, (b) leer el `y.output`/dump de conflictos que CUP reporta cuando la gramática falla.
**Acción:** nota nueva `Elementos LR(0) y la tabla SLR` con el ejemplo canónico `E→E+T…` completo: autómata + tabla + traza de `id*id+id`. El recuadro "los elementos como estados de un AFN" (p. 257) regala la conexión con el cap. 3: la colección canónica ES la construcción de subconjuntos aplicada a los ítems.

### B2. FIRST y FOLLOW no son ejecutables a mano (§4.4.2, pp. 220–222) — **ALTA**
La nota da la tabla de resultados del ejemplo 4.30 (correcta) pero **no las reglas de cálculo**: las 3 reglas de PRIMERO (terminal → él mismo; X→Y₁…Yₖ propaga mientras haya ε; X→ε agrega ε) y las 3 de SIGUIENTE ($ al inicial; A→αBβ mete PRIMERO(β)−ε en SIGUIENTE(B); si β⇒ε o A→αB, SIGUIENTE(A) ⊆ SIGUIENTE(B)). Cálculo clásico de examen.
**Acción:** ampliar la nota con las 6 reglas + un segundo ejemplo trabajado paso a paso.

### B3. El token `error` de CUP/Yacc = la teoría del modo pánico de los proyectos (§4.9.4, pp. 295–297 + §4.1.4) — **ALTA (conexión directa)**
La línea `instruccion ::= error PUNTO_COMA` del `parser.cup` de DataForge es exactamente el mecanismo que el libro describe: **producción de error** (§4.1.4) implementada por Yacc/CUP así: al fallar, el parser saca estados de la pila hasta hallar uno con ítem `A → ·error α`, "desplaza" el token ficticio `error`, y descarta entrada hasta poder continuar (hasta el `;`). Ninguna nota del vault explica esto — y es LA pregunta de defensa sobre recuperación de errores.
**Acción:** ampliar `Manejo de errores` con la mecánica del token `error` + `yyerrok`, citando §4.9.4, y enlazar desde la nota CUP.

### B4. Construcción de la tabla LL(1) y el parser predictivo con pila (§4.4.3–4.4.4, pp. 224–228) — **MEDIA**
La nota LL(1) define la clase y muestra el orden de expansión, pero falta el **Algoritmo 4.31** (A→α va en M[A,a] para a∈PRIMERO(α); si ε∈PRIMERO(α), va en M[A,b] para b∈SIGUIENTE(A)) y el analizador **no recursivo con pila** (traza fig. 4.21: pila / entrada / acción). También el caso del else colgante en la tabla (M[S′,e] con doble entrada, se resuelve a favor de `eS`).
**Acción:** ampliar la nota LL(1) con el algoritmo de la tabla + una traza con pila.

### B5. SLR vs LALR: el porqué, no solo el diagrama (§4.6.4, §4.7, pp. 252–275) — **MEDIA**
La nota da la jerarquía LR(0)⊂SLR⊂LALR⊂LR(1) sin el porqué: SLR reduce con **todo** SIGUIENTE(A) (demasiado grueso — la gramática `S→L=R|R, L→*R|id, R→L` falla en SLR por eso, ej. 4.48); LR(1) lleva el lookahead **en el ítem** `[A→α·, a]`; LALR fusiona estados LR(1) con el mismo corazón (I₄+I₇→I₄₇) — nunca introduce conflictos shift/reduce nuevos, pero **sí puede introducir reduce/reduce** (ej. 4.58).
**Acción:** ampliar la nota LR con el ejemplo L=R (por qué SLR falla, por qué LALR no) y la fusión de estados.

### B6. Por qué los `.cup`/`.jison` usan gramática ambigua + `precedence` (§4.8.1, pp. 278–281) — **MEDIA (conexión directa)**
El libro da la justificación de ingeniería que valida el diseño de los 5 proyectos: con `E→E+E|E*E|…` + declaraciones de precedencia, el parser tiene **menos estados** y **no pierde tiempo reduciendo producciones simples** (`E→T`, `T→F`) cuya única función es codificar precedencia. También documenta las reglas por defecto de Yacc/CUP: shift/reduce → **desplazar** (resuelve el else colgante gratis); reduce/reduce → **primera producción** del archivo.
**Acción:** agregar esto a `Ambigüedad, precedencia y asociatividad` (complementa la brecha B1 del cap. 2: gramática por niveles PARA la BNF entregable, ambigua+precedence PARA el `.cup`) y las reglas por defecto a `Conflictos…`.

### B7. "¿Por qué validás variables no declaradas en semántica y no en la gramática?" (§4.3.5, pp. 215–216) — **MEDIA-BAJA (respuesta de defensa)**
El libro demuestra que `wcw` (declarar-antes-de-usar) y `aⁿbᵐcⁿdᵐ` (aridad de funciones) **no son lenguajes libres de contexto** — por eso NINGUNA gramática puede chequearlos y se validan en semántica con la tabla de símbolos. Respuesta teórica exacta a una pregunta de defensa frecuente. Cabe como párrafo en `Manejo de errores` o en la nota GLC ("una gramática puede contar dos cosas, no tres").

### B8. Menores — **BAJA**
- Eliminación **general** (no inmediata) de recursividad izquierda (Alg. 4.19, con orden de no terminales) — la nota solo tiene el caso inmediato. Irrelevante para LALR, pero es teoría de examen.
- Las 4 estrategias de recuperación (§4.1.4): pánico, nivel de frase, producciones de error, corrección global — la nota solo nombra pánico.
- Else colgante: la gramática no ambigua matched/open (fig. 4.10) como alternativa teórica al shift-default.

## Lo que está bien (sin acción)

- `FIRST y FOLLOW`: la tabla de resultados es fiel al ejemplo 4.30 (solo falta el algoritmo — B2).
- `CUP`: nota práctica excelente (precedence, RESULT, maven-plugin) — coincide con §4.9 modulo sintaxis Java.
- `Análisis LL(1)`: definición y diagrama de expansión correctos.
- `Conflictos`: conceptos correctos; el else colgante bien identificado como shift/reduce.
- Todas las citas de fuente correctas (4.4.2, 4.4, 4.5–4.7, 4.5.4/4.8, 4.3).

## Material aprovechable para presentaciones (Fase B)

- **Las dos trazas gemelas**: parser LL con pila (fig. 4.21) y parser LR con pila de estados (fig. 4.38) sobre entradas similares — el slide comparativo definitivo descendente vs ascendente, ideal como doble stepper.
- El **autómata LR(0)** de la gramática de expresiones (fig. 4.31) — pieza central para la página de gramáticas/parsing.
- La calculadora Yacc completa (fig. 4.59) es un mini-DataForge: gramática ambigua + precedence + acciones que EVALÚAN (S-atribuida sin AST) + producción `error` — como demo puente teoría→proyecto.
- Anécdota: Panini especificó la gramática del sánscrito con una notación equivalente a BNF ~400 a.C. (p. 300).
