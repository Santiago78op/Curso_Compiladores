# Fase C.1 — Code review de VLangCherry (servidor Go completo)

**Qué se revisó:** los ~2,760 renglones del servidor leídos completos — `runtime/` (interprete, operaciones, nativas, entorno, valores, tipos, errores), `traductor/`, `analizar/`, más la gramática `.g4` y `docs/` en los puntos que los hallazgos exigieron. El cliente React se revisó solo por contrato (espeja al de CompInterpreter, ya auditado en B.3).
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5) · Incluye y supera el hallazgo V1 de C.0.

> ## ✅ LOS 5 ALTOS, CERRADOS Y VERIFICADOS EJECUTANDO — 2026-07-24
> No basta con leer el código: cada hallazgo se comprobó corriendo su **repro exacto** por `cmd/cli`. Resultados:
>
> | | Repro | Resultado observado |
> |---|---|---|
> | **A1** `mut` decorativo | `x := 5` · `x = 6` | `[Semántico] línea 5 col 5: no se puede reasignar "x": fue declarada sin "mut"`, y `x` queda en 5. **Aristas también cubiertas:** `+=` y `++` sobre variable sin `mut` reportan igual; los **parámetros** de función y las variables de `for` (clásico y rango) siguen siendo **mutables**, como recomendaba el informe. |
> | **A2** cortocircuito | `x := 0` · `if x != 0 && 10/x > 1 {}` | **Cero errores**: la división nunca se evalúa. El espejo con `\|\|` también cortocircuita. |
> | **A3** guardas | `for true {}` | `[Semántico] línea 3 col 9: posible ciclo infinito en el for (se superó el máximo de iteraciones)` — corta y devuelve. |
> | **A3** guardas | `func f() int { return f() }` | `[Semántico] recursión demasiado profunda al llamar "f"` — **el proceso sobrevive**, que era el punto: en Go el desbordamiento de pila es fatal y `recover()` no lo atrapa. Constantes en `interprete.go:38-39` (`maxIteraciones` 1 000 000, `maxProfundidad` 2 000). |
> | **A4** receptor por valor | método con receptor `(p Persona)` que muta | Clonado con `ClonarPorValor` (`valores.go:56`, recursivo en structs anidados); por valor no muta al llamador, por puntero sí. |
> | **A5** slice/struct sin init | `mut xs []int` · `xs = []int{1,2}` | Sin error: imprime `nil`, `len` 0, asigna y hace `append`. Ídem con struct: `nil` → `Punto{X: 1, Y: 2}` → acceso a campo. `ValorPorDefecto` conserva el tipo declarado (`tipos.go:101-105`). |
>
> **También cerrados (medios/bajos que el informe listaba):** M1 (`return` sin valor en función tipada → error), M2 (`return v` en void → error), M3 (`append`/`join` sobre slice nil, sin panic), B2 (campo duplicado en literal de struct → error). Existe `entradas/ejemplo7_semantica.vch` ejercitando los caminos felices.
>
> ### ⚠️ Lo que SIGUE abierto (reproducido el 2026-07-24)
>
> - **B3 — ✅ CORREGIDO 2026-07-24.** El síntoma era: con `func main() { x := }` la salida traía el error sintáctico correcto **y además** `[Semántico] línea 0 col 0: error interno durante la ejecución: runtime error: invalid memory address or nil pointer dereference`. Aparecía con cualquier typo y hacía parecer roto al intérprete.
>   **Fix aplicado** (`errores.go` + `analizar.go`, ~12 líneas): se agrega `ListaErrores.HayDeEntrada()` y el pipeline recuerda *antes* de traducir si el parseo ya había fallado. Si venía roto, el panic del traductor se traga en silencio — el error sintáctico ya explica todo. Si el código **parseó limpio** y aun así hubo panic, eso sí se reporta: esconderlo sería peor que el ruido.
>   Se descartó la opción «saltar la interpretación si hay errores sintácticos» porque habría matado el criterio *best-effort* documentado en §8.1: verificado que con entrada rota los errores **semánticos se siguen recolectando** (test `TestBestEffortSigueReportandoSemanticos`).
>   **Regresión:** `internal/analizar/analizar_test.go`, 13 casos que cubren B3, best-effort, y los altos A1/A2/A3/A5 + M1 — los que sobrevivieron tanto tiempo justamente porque ningún ejemplo los ejercitaba. La diapositiva del `recover()` de `presentacion-vlangcherry/etapa6.html` quedó sincronizada con el código nuevo.
> - **M4 — ✅ RESUELTO como decisión de diseño documentada** (no es un bug). El informe pedía «decidir y documentar»: está hecho en el comentario de `registrarSimbolo` (`interprete.go:169-180`). La tabla es un **registro de declaraciones — una fila por sitio de declaración en el fuente**, y la clave pasó de `ambito::nombre` a `ambito::nombre::linea:col`, con lo que (a) dos bloques hermanos que declaran el mismo nombre ya no se pisan y (b) una declaración dentro de un ciclo o de una función recursiva ocupa siempre la misma posición → una fila, no 177. La columna Valor es el valor **al declarar** (snapshot), deliberadamente: «la tabla describe el programa, no cada activación en runtime; para ver valores finales se usa `print`».
>   *Nota de lectura:* observar `0` donde la variable terminó en `999` es el comportamiento correcto bajo esa decisión, no una regresión.
> - **Observación nueva (no estaba en el informe):** cuando la guarda de ciclo infinito salta, la ejecución **continúa** después del `for` (imprime lo que sigue). Es coherente con la política «acumular y seguir» del proyecto, pero conviene decidirlo a conciencia: tras un ciclo cortado por guarda, lo que se imprima después puede confundir.

## Veredicto

Arquitectura **muy limpia** (traductor desacoplado, type-switch, reflection para el grafo, coerción aplicada con uniformidad notable en los 6 puntos de chequeo) y las 4 correcciones de la auditoría anterior están todas en su lugar. Pero el review encontró **4 bugs ALTOS** — dos de semántica del lenguaje que contradicen la propia documentación del proyecto, uno de robustez que puede matar el proceso del servidor, y el ya conocido de receptores — más una cola de medios/bajos. Ninguno lo ejercita ningún `entradas/*.vch`: todos son latentes.

## Hallazgos ALTOS

### A1 🐞 `mut` es decorativo: la inmutabilidad no se valida nunca
- La gramática acepta `'mut'?` en ambas declaraciones (`VLangCherry.g4:93-94`) y los ejemplos lo usan 18 veces, pero `traducirDeclVariable` (traductor.go:161-181) **nunca lee `MUT()`**, `ast.DeclVariable` no tiene campo de mutabilidad, y el runtime no la conoce.
- El propio `ManualUsuario.md:123` documenta: "`mut` indica que la variable puede reasignarse" — es decir, **reasignar una variable sin `mut` debería ser error semántico**, y hoy se acepta en silencio.
- **Repro:** `x := 5  x = 6` → corre sin error (debería reportar).
- **Fix:** campo `Mutable bool` en `DeclVariable` + capturarlo en el traductor; `Entorno` guarda la bandera junto a la celda (p. ej. `map[string]*celda{valor *Valor, mutable bool}` o un set paralelo); `ejecutarAsignacion`, `ejecutarIncDec` y los `+=`/`-=` validan contra ella. Los parámetros de función y las variables de `for` (init, índice/valor del rango) hay que decidirlos y documentarlos (sugerencia: mutables, como hasta ahora).

### A2 🐞 `&&` y `||` sin cortocircuito
- `interprete.go:661-670` (`case *ast.ExprBinaria`): SIEMPRE evalúa `Izquierda` y `Derecha` antes de despachar a `Y`/`O`. Es **el mismo bug que CompScript encontró y corrigió** en su auditoría (documentado en presentacion-compscript etapa5 como "bug real, no decisión").
- **Repro:** `x := 0  if x != 0 && 10/x > 1 { }` → reporta "división entre cero" en la guarda que debía evitarla.
- **Fix:** tratar `&&`/`||` antes del despacho genérico: evaluar `iz`; si decide (`false` para `&&`, `true` para `||`), retornar sin evaluar `der`. Verificar que `iz` sea bool antes de cortocircuitar.

### A3 🐞 Sin guardas de ciclo infinito ni de profundidad de recursión — puede matar el proceso
- No existe `MAX_ITER` ni `MAX_DEPTH` (contraste: CompInterpreter usa 1,000,000 y 2,000).
- `for true {}` deja la goroutine de esa petición girando para siempre: la petición nunca responde y un core queda al 100% — por cada petición así.
- Peor: recursión sin caso base → **stack overflow de Go, que es `fatal error`, NO un panic recuperable**: el `defer recover()` de `analizar.go` no lo atrapa y **el proceso del servidor muere entero**. La diapositiva de etapa6 presenta `recover()` como la protección del servidor — es insuficiente exactamente en este caso.
- **Repro:** `func f() int { return f() }  func main() { f() }` → mata el servidor.
- **Fix:** contador de profundidad en `invocarFuncion` (error semántico "recursión demasiado profunda" al pasar ~2,000) y contador de iteraciones por ciclo en `ejecutarFor` (~1,000,000). Ambos son ~10 líneas.

### A4 🐞 Receptor por valor comparte el struct (de C.0, consolidado acá)
- `ReceptorPuntero` capturado en el AST pero jamás leído por el runtime; `invocarFuncion` (interprete.go:815) declara `*receptor` con copia superficial → el método por valor muta el `StructVal` del llamador.
- **Fix:** si `!fn.ReceptorPuntero` y el receptor es struct, clonar `StructVal` (incl. los `*Valor` de `Campos`; structs anidados recursivo; slices internos pueden compartirse — semántica Go). + caso de prueba en `ejemplo2_structs.vch`.

### A5 🐞 Declarar slice/struct sin inicializar rompe toda asignación posterior
- `ValorPorDefecto` (tipos.go:101-103) devuelve para slice/struct un `Valor{Tipo: TipoNil()}` — **pierde el tipo declarado**.
- Luego `ejecutarAsignacion` valida `tipoCompatible(ptr.Tipo, v)` contra ese `nil`: asignar un `[]int` real da "no se puede asignar []int a variable de tipo nil".
- **Repro:** `mut xs []int` (sin `=`) y luego `xs = []int{1, 2}` → error semántico espurio.
- **Fix elegante:** el propio código ya anticipa la forma correcta — `EsNil` (valores.go:51) contempla `Valor{Tipo: TSlice, Slice: nil}`. Basta que `ValorPorDefecto` devuelva `Valor{Tipo: t}` (con `Slice`/`Struct` en nil) en vez de `TipoNil()`. Revisar de paso que `AImprimir` y `len` ya los manejan (sí: imprimen "nil" y len 0).
- Nota: este fix hace alcanzable el hallazgo M3 (append sobre slice nil) — corregir juntos.

## Hallazgos MEDIOS

### M1 — `return` sin valor en función no-void devuelve `0`
El zero-value de `Tipo` es `TInt` (iota 0): en `SentenciaReturn` sin expresión, `Senal.Valor = Valor{}` tipa como int 0. `func f() int { return }` retorna 0 **sin error** (y en `float64`, 0.0 vía coerción). Fix: bandera `TieneValor` en la señal (o `Valor *Valor`) y reportar "return sin valor en función que declara tipo de retorno".

### M2 — `return <valor>` en función void se descarta en silencio
`invocarFuncion` (interprete.go:834-836): si `TipoRetorno == nil` retorna Void ignorando la señal. CompScript reporta este caso como error semántico; acá pasa callado. Fix: si `senal.Tipo == SenalReturn && senal tiene valor` → error.

### M3 — `append`/`join` sobre slice nil → panic (inconsistente con `len`/`indexOf`)
`nativas.go:65` (`v.Slice.Elems`) y `:91` (`v.Slice.TipoElem`) desreferencian sin chequear nil; `len` (:45) e `indexOf` (:74) sí chequean. Con A5 corregido, `append(xs, 1)` sobre un slice declarado sin init haría panic → "error interno" genérico línea 0. Fix: chequeo nil (append sobre nil = crear slice de 1, como Go; join sobre nil = "").

### M4 — Tabla de símbolos: sobreescritura por clave y valor de declaración
`registrarSimbolo` usa clave `ambito::nombre` en un map: (a) las 177 activaciones de `fib` colapsan en UNA fila `fib::n`; (b) dos bloques hermanos dentro de la misma función comparten `ent.Nombre` → se pisan; (c) la columna Valor es el valor **al declarar** — las asignaciones posteriores no actualizan la fila. Es una decisión (snapshot vs. el log histórico de CompScript) pero hoy es la peor mezcla: ni historia ni estado final. Decidir y documentar; si 8.2 espera valores finales, actualizar la fila en cada asignación (el map de punteros ya lo permite) o regenerar al final.

## Hallazgos BAJOS

- **B1** `resolverLugar` y `evaluar` en su caso default devuelven `false` **sin mensaje** (interprete.go:447, :716) — un nodo no soportado falla en silencio. Agregar error defensivo "expresión no soportada".
- **B2** Campo duplicado en literal de struct (`Persona{Nombre:"a", Nombre:"b"}`) no reporta — el último gana en silencio.
- **B3** Tras errores sintácticos, traducir el parse tree roto puede acabar en panic→"error interno (línea 0)" — ruido junto a los errores reales. Opcional: si hubo errores sintácticos, saltar la interpretación (decisión de diseño; hoy es "best-effort" documentado).
- **B4** `print` y `println` son idénticos (cada entrada de consola ya es una línea) — inofensivo, pero documentarlo si el enunciado los distingue.

## Lo que está BIEN (para la defensa)

- Las 4 correcciones de la auditoría del 21/07 verificadas presentes: break/continue con `profLoop`/`profSwitch` **con save/restore por invocación** (mejor diseño que validar al final como CompInterpreter), colisión de nombres en las 2 direcciones, switch que reporta tipos incomparables, `Relacional` con strings/runes.
- Coerción int↔float64 aplicada con `tipoCompatible`+`coercionar` en LOS SEIS puntos (declaración, asignación, parámetros, retorno, elementos de slice, campos de struct) — uniformidad difícil de lograr y lograda.
- `for` con los ámbitos correctos: clásico con `entFor` compartido para `init`/actualización + hijo por iteración; rango con entorno fresco por vuelta.
- Switch: default corre solo si ningún caso matcheó, aunque esté declarado en medio; `continue` atraviesa el switch hasta el ciclo ✓.
- Escapes resueltos en el traductor byte a byte preservando UTF-8; posiciones línea/columna consistentes (+1) en todo el pipeline.
- `oyenteErrores` unificado lexer/parser, entorno fresco por petición, CORS, contrato JSON estable.

## Plan de corrección sugerido (Fase D, orden)

1. **A5** (ValorPorDefecto) + **M3** (nil en nativas) — juntos, ~15 líneas.
2. **A2** (cortocircuito) — ~15 líneas en `evaluar`.
3. **A3** (guardas) — ~10 líneas.
4. **A1** (`mut`) — el más transversal: AST + traductor + entorno + 3 validaciones.
5. **A4** (clonado de receptor por valor) + caso de prueba.
6. **M1/M2** (return) — juntos, tocan `Senal` e `invocarFuncion`.
7. M4/B1/B2 según tiempo. Agregar a `entradas/` un `ejemplo7_semantica.vch` que ejercite TODO lo anterior (hoy ningún ejemplo toca estos caminos — por eso sobrevivieron).
8. Tras los fixes: actualizar `ManualTecnico.md`, el quiz de etapa3 (V1) y la diapositiva del `recover()` de etapa6.
