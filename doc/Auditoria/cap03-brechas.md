# Auditoría Fase A — Capítulo 3: Análisis léxico

**Fuente:** `doc\Boocks\dragon-md\cap03-analisis-lexico.md` (Dragón 2ª ed., pp. 109–190) leído completo, contrastado contra: `Cap 3`, `Expresiones regulares`, `Definiciones regulares`, `Autómata finito (AFN y AFD)`, `Construcción de Thompson`, `Construcción de subconjuntos`, `Minimización de AFD`, `Manejo de errores`, `JFlex`, `Token, lexema y patrón`.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto general

Cobertura sólida y con **citas todas correctas** — este capítulo tiene la mejor familia de notas del vault (la de `Minimización de AFD`, con su tabla de rondas, es la nota modelo). Las 3 notas de algoritmos usan además el MISMO ejemplo canónico del libro (`(a|b)*abb`), lo que hace el pipeline ER→AFN→AFD→mínimo perfectamente trazable. Las brechas son de *ejecutabilidad* (una nota muestra el resultado pero no cómo llegar a él a mano) y del eslabón que une la teoría con JFlex.

## Brechas (ordenadas por impacto)

### B1. La construcción de subconjuntos no es ejecutable a mano (§3.7.1, Alg. 3.20, pp. 152–155) — **ALTA**
La nota `Construcción de subconjuntos` da las dos operaciones y el AFD *resultante*, pero omite el **algoritmo** (el `while` de estados sin marcar, fig. 3.32) y la **traza trabajada** del ejemplo 3.21: `A = ε-cerradura(0) = {0,1,2,4,7}` → `Dtran[A,a] = ε-cerradura({3,8}) = {1,2,3,4,6,7,8} = B` → … → tabla de la fig. 3.35. En el examen típico de OLC1 te dan un AFN y hay que convertirlo A MANO, mostrando conjuntos.
**Acción:** ampliar la nota con el pseudocódigo del algoritmo + la tabla de traza (estado AFD ↔ conjunto de estados AFN ↔ transiciones).

### B2. Cómo el autómata tokeniza (no solo reconoce) (§3.8.1–3.8.3, pp. 166–171) — **MEDIA-ALTA**
Ninguna nota cubre el salto de "autómata que dice sí/no" a "analizador léxico que corta lexemas": (1) se combinan los AFNs de TODOS los patrones con un nuevo estado inicial y transiciones ε (fig. 3.50); (2) el autómata avanza hasta quedarse sin transiciones y **retrocede hasta el último estado de aceptación visto** — así se implementa el *longest match*; (3) si ese estado acepta varios patrones, gana el listado primero. Esta es la teoría que fundamenta las dos reglas que la nota `JFlex` enuncia como recetas ("longest match + prioridad de la primera regla").
**Acción:** nota nueva `Del autómata al analizador léxico` (o sección en `Autómata finito`) con el ejemplo 3.27 (entrada `aaba` sobre los patrones `a`, `abb`, `a*b+`: avanza, muere en la 4ª letra, retrocede a `aab`).

### B3. Conversión directa ER→AFD: anulable / primerapos / ultimapos / siguientepos (§3.9.5, Alg. 3.36, pp. 173–180) — **MEDIA**
No existe en el vault. Es el algoritmo del árbol sintáctico de `(r)#` con las cuatro funciones — clásico de exámenes de compiladores (a menudo llamado *firstpos/followpos*) y es lo que Lex hace internamente (el libro lo dice: Lesk y Schmidt usaron el Alg. 3.36).
**Acción:** nota nueva con la tabla de reglas (fig. 3.58), el árbol anotado (fig. 3.59) y la tabla de siguientepos (fig. 3.60) para `(a|b)*abb#`. Vale confirmar si el curso lo evalúa; si sí, sube a ALTA.

### B4. Recuperación de errores léxicos mal fusionada con la sintáctica — **MEDIA-BAJA (corrección)**
La nota `Manejo de errores` describe el modo pánico como "descartar hasta token de sincronización de FOLLOW" — eso es la recuperación **sintáctica** (cap. 4). La recuperación **léxica** del §3.1.4 es otra: eliminar caracteres sucesivos hasta poder formar un token bien formado (+ las 4 transformaciones: borrar/insertar/sustituir/transponer un carácter). El propio CLAUDE.md de DataForge lo tiene bien ("léxico descarta el carácter") — la nota del vault debe distinguir los dos mecanismos.
**Acción:** editar la nota separando pánico léxico (descartar caracteres, §3.1.4) de pánico sintáctico (descartar tokens hasta sincronización, §4.1.4).

### B5. Por qué separar léxico de sintáctico (§3.1.1, pp. 110–111) — **BAJA**
Las 3 razones del libro (sencillez de diseño, eficiencia con técnicas especializadas/búfer, portabilidad) son pregunta teórica clásica y no están en ninguna nota. Un párrafo en la nota `Cap 3`.

### B6. Extras de cultura teórica — **BAJA (opcionales)**
- Leyes algebraicas de ER (fig. 3.7, p. 122): conmutatividad de `|`, ε identidad, `r** = r*`.
- Trade-offs AFN vs AFD (fig. 3.48, p. 165): por qué un lexer SIEMPRE quiere AFD y `grep` a veces prefiere simular el AFN. Explica una decisión real de ingeniería.
- Caso exponencial (ejemplo 3.25): `(a|b)*a(a|b)ⁿ⁻¹` necesita 2ⁿ estados de AFD — el "dato wow" para la presentación.
- Operador lookahead `/` de Lex (§3.5.4) con la anécdota del `IF(I,J)=3` de Fortran.

## Lo que está bien (sin acción)

- `Minimización de AFD`: **nota modelo** — algoritmo, tabla de rondas del ejemplo 3.40 fiel al libro, antes/después en Mermaid. Nada que tocar.
- `Construcción de Thompson`: fragmentos correctos (unión, cerradura, concatenación por fusión). El ejemplo incremental completo (figs. 3.43–3.46) sería bonito en la *presentación*, no hace falta en la nota.
- `Expresiones regulares` y `Definiciones regulares`: fieles, con el ejemplo exacto del libro (ej. 3.6) y el puente a macros JFlex.
- `Autómata finito`: fig. 3.24 bien reproducida; el pipeline de 4 pasos bien enunciado.
- `JFlex`: prácticamente correcta y consistente con §3.5 (estructura de 3 secciones, conflictos, yytext/yyline).
- Coherencia de las 3 notas de algoritmos usando el mismo ejemplo `(a|b)*abb`: mantener esa decisión.

## Material aprovechable para presentaciones (Fase B)

- **El pipeline completo animado con un solo ejemplo**: `(a|b)*abb` — Thompson incremental (figs. 3.43–3.46) → subconjuntos con la traza de conjuntos (fig. 3.35) → minimización por rondas (ej. 3.40). Verificar si `automatas.html` ya lo hace; si no, es la mejora #1 de esa página.
- Demo del **longest match con retroceso** (ejemplo 3.27, entrada `aaba`): stepper natural.
- Anécdotas: `DO 5 I = 1.25` de Fortran (p. 113) y el caso exponencial 2ⁿ.
- Ejercicios reutilizables: 3.3.2 (describir lenguajes de ER), 3.3.5 (escribir definiciones regulares), 3.7.1/3.7.3 (conversiones completas).
