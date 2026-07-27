# Maven y los arquetipos — presentación

Presentación en HTML autocontenida (se abre con doble clic, sin servidor) sobre la
herramienta que construye los tres proyectos Java del curso.

## Páginas

| Archivo | Contenido |
|---|---|
| `index.html` | Índice |
| `maven.html` | **Cómo funciona Maven** — el problema que resuelve, coordenadas GAV, anatomía del POM, convención sobre configuración, ciclo de vida y fases, el trío fase/plugin/goal, dependencias y *scopes*, repositorios, y la chuleta de comandos. 13 diapositivas. |
| `arquetipos.html` | **El arquetipo** — qué es, cómo se usa (`archetype:generate`), qué genera, aplicaciones, cómo fabricar uno propio (`create-from-project`), el descriptor y el filtrado, publicarlo y usarlo. 13 diapositivas. |

Cada página cierra con un quiz de 4 preguntas con respuesta revelable.

## Todo el material está verificado

A diferencia de una presentación escrita de memoria, acá los comandos se ejecutaron
y la salida es real:

- El `mvn archetype:generate` con `maven-archetype-quickstart` 1.5 se corrió de
  verdad; la estructura, el `App.java` y el `pom.xml` que se muestran son los que
  salieron. El proyecto generado compila y pasa sus pruebas (`mvn test`, exit 0).
- El `archetype:create-from-project` se corrió sobre **DataForge**; los dos bloques
  `<fileSet>` de la diapositiva del descriptor son salida literal del
  `archetype-metadata.xml` generado.
- El nombre de la propiedad `archetype.filteredExtentions` —con su error de tipeo
  histórico— se confirmó con `mvn archetype:help -Ddetail=true`.
- Los fragmentos de POM salen del `DataForge/pom.xml` real del repositorio.

La diapositiva 11 de `arquetipos.html` documenta un fallo real que ocurrió al
probar `create-from-project` (el *invoker* anidado buscando un `~/.m2/settings.xml`
inexistente) y lo usa para enseñar a distinguir un problema del entorno de uno del
proyecto.

## Diseño

Comparte el sistema de diseño de las otras presentaciones (`assets/estilo.css`,
las tres tipografías servidas desde `assets/fuentes/`, y `assets/deck.js` para la
navegación). Su color de acento es el naranja tostado `#f0864a`, distinto del de
los otros seis decks.

Verificado con Playwright en escritorio (1600×900) y móvil (390×844): sin errores
de consola y sin desbordes horizontales.
