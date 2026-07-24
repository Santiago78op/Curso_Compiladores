# Fase C.0 — Las 4 verificaciones que dejó la Fase B.3

**Qué se hizo:** verificar contra el código real los 4 hechos que los informes de presentaciones dejaron marcados, antes de reescribir nada en Fase D.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## S3 — Dangling else en CompScript: CONFIRMADO (la diapositiva atribuye mal la causa)

`CompScript\src\main\cup\parser.cup:256`:

```
bloque ::= LLAVE_IZQ lista_inst:l LLAVE_DER {: RESULT = l; :} ;
if_stmt ::= IF PAR_IZQ expr PAR_DER bloque
          | IF PAR_IZQ expr PAR_DER bloque ELSE bloque
          | IF PAR_IZQ expr PAR_DER bloque ELSE if_stmt ;
```

Los cuerpos del `if` **exigen llaves**. El ejemplo de la diapositiva (`if (a) if (b) x; else y;`) ni siquiera parsea en CompScript: el else colgante **no puede ocurrir**, y por eso CUP no reporta conflicto — no por el "orden de las producciones".
**Acción D:** reescribir la diapositiva de etapa2: la ambigüedad del else colgante (Dragón §4.3.2 / §4.8.2) existe cuando el cuerpo puede ser una instrucción sin llaves; CompScript la **elimina por construcción** delimitando bloques con `{}` — que es una tercera estrategia (rediseñar la gramática), distinta de "resolver el conflicto a favor del shift" (C/Java) y digna de contarse como decisión.

## S4 — Los números 177/186 de ejemplo5.cs: RESUELTO (se reconcilian exactamente)

`ejemplo5.cs` llama: `factorial(n=5)`, `fib(n=10)`, `potencia(base=4)`, `potencia(base=2, exp=5)`.

| Fuente | Llamadas | Entradas que aporta al reporte 6.4 |
|---|---|---|
| `fib(10)`, base `n<2`, C(n)=1+C(n−1)+C(n−2) | C(10) = **177** | 177 × `n` |
| `factorial(5)` | 5 | 5 × `n` |
| `potencia` ×2 (el default `exp=2` también se declara) | 2 | 2 × `base` + 2 × `exp` |
| **Total** | | **186 filas** — de las cuales **182 son `n`** |

**Conclusión:** el "186" real es el **total de filas del reporte de símbolos**; la presentación lo atribuyó todo al "parámetro `n` de fib". Los dos números eran correctos por separado (177 llamadas de fib ✓, 186 filas ✓) pero la conexión estaba mal contada.
**Acción D:** en etapa4/etapa6 de presentacion-compscript: "el reporte muestra 186 filas: 177 `n` de fib + 5 `n` de factorial + 4 de `base`/`exp` de potencia".

## S6 — Sobrecarga del `+` en CompScript: CONFIRMADA

`CompScript\...\interprete\Operaciones.java:28-49` — la suma es la única operación que admite bool, char y cadena:
- `string + cualquier primitivo → string` (concatenación)
- `char + char → cadena`
- numérico + numérico → suma con dominancia de tipos

Es **sobrecarga de operadores resuelta por los tipos de los operandos** (Dragón §6.5.3, la brecha ALTA #10 de Fase A y H2 de B.1) implementada y sin nombrar en ningún deck.
**Acción D:** una tarjeta en presentacion-compscript etapa5 (o en cap6 del deck del Dragón, que sirve a todos): "el `+` de CompScript está sobrecargado: el intérprete elige suma o concatenación según la firma de tipos — §6.5.3". CompInterpreter tiene el mismo rasgo vía sus tablas 5.5 (verificar en C.4).

## V1 — Receptor por valor en VLangCherry: RESUELTO — y es un **BUG REAL** 🐞

**Hecho:** `ReceptorPuntero` se captura en el AST (`traductor.go:71` → `ast.go:56`) pero **el runtime jamás lo lee** — grep en todo `server/`: solo esas 2 apariciones. En `interprete.go:815`, `invocarFuncion` declara el receptor con `entFn.Declarar(fn.ReceptorNombre, *receptor)` — copia superficial del `Valor`, cuyo campo `Struct` es un puntero compartido (`StructVal{Campos map[string]*Valor}`).

**Consecuencia:** los receptores **por valor se comportan igual que por puntero**: `func (p Persona) M() { p.Edad = 30 }` SÍ muta el struct del llamador. La distinción `'*'?` de la gramática es código muerto en ejecución.

- La respuesta del quiz de etapa3 ("**25, sin cambiar**") es **falsa** con el código actual: imprimiría 26 igual.
- El bug es **latente**: ningún ejemplo de `entradas/` muta a través de un receptor por valor (`Saludar` solo lee), por eso las pruebas pasan.

**Acción (proyecto, Fase D):** en `invocarFuncion`, si `receptor != nil && !fn.ReceptorPuntero && receptor.Tipo.Base == TStruct`, clonar el `StructVal` antes de declarar (copia de `Campos` — son `map[string]*Valor`, hay que clonar también los `*Valor`; para semántica Go exacta, clonar structs anidados recursivamente; los slices dentro del struct pueden seguir compartiéndose — así es Go). Agregar un caso a `entradas/ejemplo2_structs.vch` que mute vía receptor por valor y verifique que NO cambia. **Después** corregir el quiz de etapa3 con la respuesta ya sin contradicción.

## Resumen

| # | Verificación | Resultado | Tipo de acción |
|---|---|---|---|
| S3 | Bloques con llaves eliminan el dangling else | ✅ Confirmado | Reescribir diapositiva (presentación) |
| S4 | 177/186 se reconcilian: 186 = total (182 `n` + 4 `base`/`exp`) | ✅ Resuelto | Corregir números (presentación) |
| S6 | `+` sobrecargado suma/concatenación | ✅ Confirmado | Tarjeta §6.5.3 (presentación) |
| V1 | `ReceptorPuntero` ignorado por el runtime | 🐞 **Bug real latente** | **Fix de código + caso de prueba** + corregir quiz |
