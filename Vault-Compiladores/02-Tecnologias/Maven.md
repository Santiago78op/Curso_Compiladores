---
tags: [tecnologia, build, java]
aliases: [pom.xml, POM, "ciclo de vida de Maven", "arquetipo", "archetype", "archetype:generate", "coordenadas GAV", "generate-sources", "mvn"]
fuente: "https://maven.apache.org/guides/getting-started/ · pom.xml real de DataForge · comandos ejecutados y verificados el 2026-07-27"
fecha: 2026-07-27
---

# Maven

Gestor de **build y dependencias** para Java. El **`pom.xml`** declara coordenadas (`groupId:artifactId:version`), dependencias y plugins. Descarga librerías del repositorio central a `~/.m2`.

**Estructura estándar:** `src/main/java`, `src/main/resources`, `src/test/java`.

**Comandos:** `mvn compile`, `mvn test`, `mvn package` (crea el JAR), `mvn clean`, `mvn javafx:run`.

## El ciclo de vida: la regla que lo explica todo

El ciclo *default* es una **secuencia ordenada de fases**: `validate → generate-sources → process-sources → compile → test-compile → test → package → verify → install → deploy`. (`clean` y `site` son ciclos aparte; por eso `mvn clean compile` nombra los dos.)

> **Pedir una fase ejecuta esa fase y todas las anteriores.** De ahí que `mvn compile` no solo compile: antes pasa por `generate-sources`, que es donde corren [[JFlex]] y [[CUP]].

Conviene no confundir tres palabras:

| Concepto | Qué es | Ejemplo |
|---|---|---|
| **Fase** | Un momento del build; por sí sola no hace nada | `generate-sources` |
| **Plugin** | Un artefacto que aporta capacidades, con sus propias coordenadas | `de.jflex:jflex-maven-plugin:1.9.1` |
| **Goal** | La acción concreta que se ejecuta | `jflex:generate` |

El build ocurre cuando un **goal se engancha a una fase**. Cada goal trae su fase por defecto (la de `jflex:generate` es `generate-sources`), y por eso el POM no necesita declararla. Un goal también se puede invocar suelto con `prefijo:goal` — pero entonces **no arrastra las fases anteriores**, que es por qué a veces hay que escribir `mvn compile exec:java`.

## Plugins para compiladores (generan código en `generate-sources`)

| Plugin | Coordenadas | Qué hace |
|--------|-------------|----------|
| **cup-maven-plugin** | `com.github.vbmacher:cup-maven-plugin:11b-20160615-3` | Corre [[CUP]] sobre `src/main/cup/parser.cup` → `target/generated-sources/cup` (ver XML en [[CUP]]) |
| **jflex-maven-plugin** | `de.jflex:jflex-maven-plugin` | Corre [[JFlex]] sobre `src/main/jflex/*.flex` → `target/generated-sources/jflex` |

Con ambos, `mvn compile` regenera lexer y parser solo — sin correr `jflex`/`java -jar` a mano.

> **Diagnóstico útil:** «no encuentra la clase `Parser`» casi nunca es un problema de import. Significa que `generate-sources` no produjo nada, normalmente por un error de sintaxis en el `.cup` o porque el archivo no está donde el plugin lo busca. La pista está más arriba en la salida de Maven.

> **Regla de `target/`:** todo lo que hay ahí es derivado y se regenera; por eso no se versiona. Verificado clonando el repo limpio: `mvn compile` reconstruye `Parser.java`, `sym.java` y `Lexer.java` desde cero.

## Arquetipos: plantillas de proyecto

Un **arquetipo** es una plantilla de proyecto parametrizada, **empaquetada y distribuida como un artefacto Maven más** (tiene sus propias coordenadas GAV, se versiona y se publica igual que una librería). Si una clase es un molde para objetos, un arquetipo es un molde para *proyectos*.

**Usar uno** — separando el *molde* (`-Darchetype*`) de la *pieza* (`-DgroupId`, `-DartifactId`…):

```
mvn archetype:generate -B \
  -DarchetypeGroupId=org.apache.maven.archetypes \
  -DarchetypeArtifactId=maven-archetype-quickstart \
  -DarchetypeVersion=1.5 \
  -DgroupId=com.olc1 -DartifactId=demo -Dversion=1.0.0 -Dpackage=com.olc1.demo
```

Sin `-B` y sin los `-D`, el comando abre un asistente interactivo. El proyecto que sale ya compila y pasa sus pruebas, y trae decisiones hechas: UTF-8 explícito, `maven.compiler.release`, y un BOM de JUnit en `<dependencyManagement>`.

**Fabricar el propio:** `mvn archetype:create-from-project` sobre un proyecto que ya funciona. Deja en `target/generated-sources/archetype` un descriptor `META-INF/maven/archetype-metadata.xml` y la plantilla en `archetype-resources/`. Después, `mvn install` para publicarlo en el repositorio local y `-DarchetypeCatalog=local` para usarlo.

**Los huecos van en dos sintaxis distintas:** `${variable}` dentro del contenido de un archivo, y `__variable__` en los nombres de archivo y carpeta (porque `$`, `{` y `}` no son portables en nombres).

> **La trampa para un proyecto de compiladores:** el descriptor marca cada conjunto de archivos como `filtered="true"` (pasa por el motor de plantillas) o sin filtrar (se copia byte a byte). Una **gramática nunca debe filtrarse**: está llena de símbolos que el motor interpretaría como suyos y la devolvería corrompida. Al convertir [[DataForge]] se comprobó que el plugin filtró los `.java` y `.txt` y dejó sin filtrar los `.flex` y `.cup` —porque decide por extensión y esas no las conoce—, lo cual es lo correcto acá, pero conviene verificarlo. La lista se controla con `archetype.filteredExtentions` (con ese error de tipeo en el nombre).

Material completo con los comandos ejecutados y su salida real: la presentación `presentacion-maven/` del workspace (`maven.html` y `arquetipos.html`).

> Unifica el build de los tres proyectos Java (incluye deps como JavaFX, CUP, Gson).

## Usado en
[[DataForge]], [[ConjAnalyzer]], [[CompScript]]

## Relacionadas
- [[JavaFX y Scene Builder]]
- [[CUP]]
- [[JFlex]]
