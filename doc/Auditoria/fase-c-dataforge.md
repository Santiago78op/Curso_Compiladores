# Fase C.5 — Code review de DataForge (Java + JFlex + CUP, S-atribuido sin AST)

**Qué se revisó:** ~1,370 renglones leídos completos — `Entorno`, `Operaciones`, `parser.cup`, `Lexer.flex`, `Graficador`, `Reportes` — más `Interprete`/`Simbolo`/`Grafica`/`RegistroError` (triviales) y skim de GUI/tests (patrón familia ya verificado 2 veces).
**Fecha:** 2026-07-22 · **Auditor:** Claude (Fable 5)

## Veredicto: limpio — y el barrido cruzado confirma por qué

Cero hallazgos altos y medios. Además, el **barrido cruzado** de los bugs encontrados en los otros 4 proyectos da todo N/A por diseño: sin operadores lógicos no hay bug de cortocircuito posible; sin indexación no hay validación de índices; sin AST no hay `.dot` que escapar; sin ciclos ni funciones no hacen falta guardas anti-cuelgue. **La decisión arquitectónica "sin AST porque no hay control de flujo" también eliminó clases enteras de defectos** — un argumento de defensa que no estaba formulado así en ningún material.

Verificaciones que pasaron:
- Los **dos formateadores** (`formatear` consola sin comillas/15.0→"15" vs `valorReporte` §6.3 con comillas) implementados y usados cada uno en su lugar, exactamente como documenta el CLAUDE.md.
- **Variables y arreglos comparten espacio de nombres con chequeo cruzado** (`simbolos.containsKey` antes de declarar cualquiera) — mejor que ConjAnalyzer (B2 de C.4).
- El **typo `graphBar`/`grapBar` del enunciado** resuelto con doble alternativa en el `.flex` con su comentario-justificación.
- `validarGrafica` como contrato: el `Graficador` castea sin re-chequear y lo declara en su javadoc — responsabilidades limpias.
- "La última gana" en atributos repetidos + EXEC requerido-pero-no-error, fiel a 5.10.
- Estadística correcta: mediana par/impar, varianza poblacional documentada, moda con desempate "primera en aparecer" (comentario y `LinkedHashMap` coinciden).
- El comentario del `.flex` sobre orden de reglas es **correcto acá** ("a igual longitud gana la primera" — reservadas vs `{Id}` sí es el caso de empate), a diferencia de los comentarios de ConjAnalyzer/CompInterpreter.
- `esc()` en todas las celdas, reflexión sobre `sym`, entorno fresco por ejecución, modo pánico `error PUNTO_COMA`.

## Hallazgos (todos BAJOS)

- **B1** — `Cadena = \"[^\"]*\"` **admite saltos de línea**: una comilla sin cerrar se traga todo hasta la próxima comilla del archivo, produciendo un error sintáctico lejos del origen real. Los otros proyectos excluyen `\r\n` en sus cadenas. Fix de 4 caracteres: `\"[^\"\r\n]*\"` (una cadena sin cerrar pasa a ser error léxico en su propia línea).
- **B2** — En `attr ::= ID DOBLE_DOS_PUNTOS tipo IGUAL valor_attr ...` el **`tipo` declarado se parsea y se descarta** — `titulo :: double = "hola" end;` no reporta la incoherencia declarada-vs-real (la validación efectiva la hace `validarGrafica` por atributo, así que el error de fondo SÍ se detecta, pero con otro mensaje). Documentar la decisión o validar coherencia.
- **B3** — `EXEC tipo_graf` **no compara** el tipo del EXEC con el del bloque: `exec graphPie end;` dentro de un `graphBar(...)` se acepta. Decidir (¿error semántico o irrelevante?) y documentar en `docs/gramatica.txt`.
- **B4** — `exigirLista` rechaza listas vacías con el mensaje "falta el atributo" — engañoso para `ejex :: double = [] end;` (el atributo está, pero vacío). Separar el mensaje ("está vacío").

## Observación final de la Fase C

DataForge y ConjAnalyzer (los dos construidos/auditados más temprano y con menos superficie semántica) salieron limpios; CompScript salió sólido con decisiones a documentar; los dos proyectos más nuevos y grandes (CompInterpreter, VLangCherry) concentran los bugs reales — y **repiten defectos que CompScript ya había corregido** (cortocircuito, return en void). La lección para Fase D y para la defensa: los fixes deben propagarse entre proyectos, no quedarse donde se encontraron.
