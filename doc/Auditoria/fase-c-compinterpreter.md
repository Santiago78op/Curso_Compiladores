# Fase C.3 — Code review de CompInterpreter (JS + Jison + Express)

**Qué se revisó:** el servidor completo (~1,400 renglones propios, sin el parser generado): `interprete/` (interprete.js 672, operaciones, nativas, entorno, valor, tipos, errores), `reportes/ast-grafo.js`, `analizar.js`, `server.js`, y `grammar.jison` en los puntos sensibles (precedencia, `error`, `%options`, nodos LOGICA). El cliente React se revisó en B.3 (contrato) y no se repite.
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto

Runtime bien diseñado: tablas literales del enunciado con las 2 decisiones difíciles documentadas en el código (el `%` que da DECIMAL; la matriz relacional ilegible del PDF reconstruida como decisión propia), guardas `MAX_ITER`/`MAX_DEPTH` (único proyecto con ambas), los 2 fixes de la auditoría previa verificados presentes, y fall-through de switch correcto incluso cayendo al `default`. Pero comparte con VLangCherry **el bug del cortocircuito**, y el review encontró un hueco de validación de índices y dos silencios alrededor de métodos void.

## Hallazgos

### A1 🐞 (ALTA) — `&&` y `||` sin cortocircuito
`interprete.js:531-535` (case `'LOGICA'`): evalúa `izq` **y** `der` siempre, y recién entonces llama a `ops.logica`. Es el mismo bug que CompScript encontró, documentó como "bug real, no decisión" y corrigió — y que VLangCherry también tiene (A2 de C.1). Tres proyectos, mismo defecto, uno solo corregido.
- **Repro:** `let x: int = 0;` … `if (x != 0 && 10 / x > 1)` → reporta "División entre cero" en la guarda que debía evitarla.
- **Fix (~8 líneas):** en el case `'LOGICA'`, si `op` es `&&`/`||`: evaluar `izq`; si es BOOL y decide (`false`/`true` respectivamente), retornar sin evaluar `der`. El caso `!` queda igual.

### M1 🐞 (MEDIA) — índices de vector sin validar tipo
`evalAcceso` (interprete.js:587) y `execAsignacionVector` (:281,291) hacen `Math.trunc(aNumero(idx))` **sin verificar que el índice sea ENTERO** — a diferencia de `VECTOR_NEW` (:181), que sí lo valida (inconsistencia interna).
- Índice string/null → `aNumero` da `NaN` → `NaN < 0` y `NaN >= len` son ambos `false` → **pasa los chequeos de rango** → `arr[NaN]` = `undefined`, que aguas abajo o se vuelve "Error interno durante la ejecución" o se pierde en silencio.
- Índice `double` (1.9) → se trunca a 1 en silencio en vez de reportar.
- **Fix:** validar `idx.tipo === TIPO.INT` en los 4 puntos (lectura 1D/2D, asignación 1D/2D), como ya hace VECTOR_NEW.

### M2 🐞 (MEDIA) — método (void) usado como expresión: null silencioso
`evalLlamada` → `invocar` devuelve `null` para un MÉTODO **sin reportar nada**. `let x: int = miMetodo();` no genera ningún error: `x` queda declarada con valor JS-null (fila vacía en el reporte), y cualquier uso posterior se propaga por null como si "ya hubiera habido un error antes" — pero nunca lo hubo. Contrasta con `ejecutar f()` sobre una FUNCIÓN, que sí reporta (:87-90).
- **Fix:** en `evalLlamada` (o `invocar` con una bandera "comoExpresion"), si `fn.tipo === 'METODO'` y el resultado se necesita como valor → error semántico "un método no retorna valor".

### M3 (MEDIA-BAJA) — `return <valor>` dentro de un método se descarta en silencio
`invocar` valida la señal RETURN solo para FUNCION (:493-505); para un MÉTODO, un `return 5;` se consume sin error (CompScript sí lo reporta: "es void y no puede retornar un valor"). Mismo hallazgo M2 de VLangCherry — patrón repetido en los dos proyectos post-CompScript.
- **Fix:** si `fn.tipo === 'METODO' && senal?.tipo === 'RETURN' && senal.tieneValor` → error semántico.

### M4 (MEDIA-BAJA) — nativas numéricas sobre vector `new` → `NaN` silencioso
`new vector int[3]` llena con `Valor.nulo()` (por diseño, :184). Pero `sum`/`average`/`max`/`min` hacen `aNumero(elemento)` que para `null` da `NaN`: `echo sum(v);` sobre un vector recién creado imprime `NaN` sin ningún error. **Fix:** en las 4 nativas, si algún elemento es de tipo `null` → error semántico ("el vector contiene elementos nulos sin asignar") o documentar un tratamiento (contar como 0).

### Bajos
- **B1** Coerción asimétrica: `char`→`double` implícita permitida (`coercionar`, :639), `char`→`int` no. Si no viene de las tablas del enunciado, unificar.
- **B2** Cuando FALTA un argumento sin default, se reporta el error pero la función **igual se ejecuta** con el parámetro en valor por defecto (:466-469). Política "continuar" válida (opuesta al abort de CompScript) — documentarla en el Manual Técnico.
- **B3** El default de un parámetro se evalúa en el entorno `local` (:464), donde ya viven los parámetros anteriores → `f(a: int, b: int = a + 1)` funciona. Feature útil y sin documentar.
- **B4** `valorPorDefecto` de `char` es `' '` con el comentario ambiguo «"carácter 0" del enunciado» (tipos.js:38) — si el enunciado dice `\0` (código 0), difiere; verificar y alinear código o comentario.
- **B5** Tabla de símbolos por clave `ambito::id` (última gana), pero **actualizada en cada asignación** → muestra valores finales ✓. No es log histórico (a diferencia de CompScript): documentar la diferencia, es pregunta de defensa.

## Lo que está BIEN (pro-defensa)

- **Único proyecto con ambas guardas** `MAX_ITER` (1M, por ciclo, contador propio por entrada de bucle) y `MAX_DEPTH` (2000) — exactamente lo que le falta a VLangCherry (A3 de C.1).
- Los 2 fixes de la auditoría del 21/07 verificados: `coercionarElementoOValorDefecto` con su comentario-justificación (:658-669) y la validación de `break`/`continue` escapados en `invocar` (:486-491).
- Switch con fall-through fiel a 5.16.2, incluida la caída al `default` tras un case sin break — y el `default` también corre si ningún case coincidió.
- Decisiones documentadas EN el código: `%` → DECIMAL "tal cual la fuente, sin corregirlo"; matriz relacional reconstruida "como decisión propia, NO como dato del enunciado"; coerción double→int justificada con el Anexo 11.1.
- `registrarSimbolo` con `entornoNombre` fijado al declarar — reasignar una global dentro de una función no duplica la fila.
- `ast-grafo.js` escapa `\n` y `"` en el `.dot` (CompScript no escapa `\n` — B4 de C.2).
- Args evaluados en el entorno del llamador, params declarados en hijo del global (alcance estático) ✓.

## Insumos para Fase D (orden sugerido)

1. **A1** cortocircuito (~8 líneas) + caso de prueba en `entradas/` (`x != 0 && 10/x > 1`).
2. **M1** validación de índices (4 puntos, ~8 líneas).
3. **M2+M3** métodos void como expresión / con return valuado (juntos, tocan `evalLlamada`/`invocar`).
4. **M4** nulls en nativas; B1/B4 tras verificar el enunciado.
5. Actualizar Manual Técnico (políticas B2/B5) y, en la presentación, la diapositiva de señales puede sumar el contraste "CompScript corrigió el cortocircuito; acá se corrigió en Fase D" como continuación de la historia de auditorías.
