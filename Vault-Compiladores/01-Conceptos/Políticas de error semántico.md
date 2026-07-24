---
tags: [concepto, semantico, defensa]
aliases: ["propagar null", "abortar en el primer error", "acumular errores", "recuperación semántica", "qué hacer tras un error de tipos", "política de errores"]
fuente: "Libro del Dragón cap. 6 (detección) · revisión de código de los 5 proyectos del curso (Fase C, 2026-07-22) contra sus enunciados"
fecha: 2026-07-24
---

# Políticas de error semántico

El [[Manejo de errores (léxicos, sintácticos, semánticos)|Libro del Dragón]] explica en detalle **cómo detectar** un error de tipos (§6.5) y cómo recuperarse de uno **sintáctico** (modo pánico, token `error`). Lo que **no** define es qué hacer *después* de detectar uno semántico: ahí no hay pánico ni sincronización, y la decisión la fija el enunciado de cada proyecto.

Los cinco intérpretes del curso tomaron **tres caminos distintos, y los tres son correctos** — cada uno por su requisito. Saber cuál se eligió y por qué es una pregunta de defensa casi garantizada.

## Las tres políticas

| Política | Proyecto | Mecanismo | Por qué es la correcta ahí |
|---|---|---|---|
| **Propagar y seguir** | [[DataForge]] | La expresión con error devuelve `null`; toda operación que reciba `null` calla y devuelve `null` | Su enunciado (§6) pide un reporte de errores **completo**: abortar en el primero daría un reporte de una línea |
| **Abortar** | [[CompScript]] | El error semántico corta la ejecución de esa instrucción | Su enunciado (§4.3) lo exige. Seguir con un valor inventado sería *menos* correcto |
| **Acumular y seguir** | [[CompInterpreter]] · [[VLangCherry]] | El error se agrega a una lista y la ejecución continúa con un valor neutro | Son servicios REST: una petición devuelve **todos** los errores de una vez, porque no hay segunda oportunidad interactiva |

## La disciplina que cada una exige

- **Propagar y seguir** solo funciona si la propagación es **silenciosa**: si una operación que recibe `null` también reportara, un error produciría una cascada de diez. La regla es *un error por causa*. Es más difícil de sostener de lo que parece — hay que revisar cada operación.
- **Abortar** es la más simple de implementar y la que da mensajes más limpios, pero obliga al usuario a corregir de a un error por corrida.
- **Acumular y seguir** necesita definir el **valor neutro** con el que se continúa, y ese valor puede generar errores derivados que confunden. Conviene marcarlos o tolerarlos conscientemente.

## Cómo se usa esto en la defensa

Si preguntan *«¿por qué tu intérprete no se detiene en el primer error?»*, la respuesta no es «así lo programé». Es:

> «Porque mi enunciado pide un reporte completo, y esa decisión se paga con la disciplina de propagar `null` sin generar errores en cascada.»

La política es una **decisión de diseño defendible**, no un accidente — pero solo si sabés cuál elegiste y qué te cuesta.

## La lección transversal: los fixes se propagan

La revisión de código encontró el **mismo bug** en varios proyectos a la vez. El caso emblema es el [[Flujo de control y switch|cortocircuito de `&&`/`||`]]: [[CompScript]] lo detectó y corrigió en su propia auditoría, pero [[VLangCherry]] y [[CompInterpreter]] —escritos *después*— lo repitieron. Lo mismo con `return` mal validado.

De ahí la regla de trabajo: **todo fix confirmado en un proyecto se verifica en los otros cuatro**. Un hallazgo no es «un bug de este proyecto», es «una clase de bug que mi forma de escribir intérpretes tiende a producir». Ver también [[Árbol de sintaxis abstracta (AST)]], donde la misma idea aparece del lado opuesto: qué clases de bug se vuelven *imposibles* según la arquitectura elegida.

## Relacionadas
- [[Manejo de errores (léxicos, sintácticos, semánticos)]]
- [[Comprobación de tipos]]
- [[Flujo de control y switch]]
- [[Árbol de sintaxis abstracta (AST)]]
- [[Cap 6 - Generación de código intermedio]]
