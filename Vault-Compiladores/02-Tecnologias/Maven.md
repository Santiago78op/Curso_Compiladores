---
tags: [tecnologia, build, java]
fuente: "https://maven.apache.org/guides/getting-started/"
fecha: 2026-07-10
---

# Maven

Gestor de **build y dependencias** para Java. El **`pom.xml`** declara coordenadas (`groupId:artifactId:version`), dependencias y plugins. Descarga librerías del repositorio central a `~/.m2`.

**Estructura estándar:** `src/main/java`, `src/main/resources`, `src/test/java`.

**Comandos:** `mvn compile`, `mvn test`, `mvn package` (crea el JAR), `mvn clean`, `mvn javafx:run`.

## Plugins para compiladores (generan código en `generate-sources`)

| Plugin | Coordenadas | Qué hace |
|--------|-------------|----------|
| **cup-maven-plugin** | `com.github.vbmacher:cup-maven-plugin:11b-20160615-3` | Corre [[CUP]] sobre `src/main/cup/parser.cup` → `target/generated-sources/cup` (ver XML en [[CUP]]) |
| **jflex-maven-plugin** | `de.jflex:jflex-maven-plugin` | Corre [[JFlex]] sobre `src/main/jflex/*.flex` → `target/generated-sources/jflex` |

Con ambos, `mvn compile` regenera lexer y parser solo — sin correr `jflex`/`java -jar` a mano.

> Unifica el build de los tres proyectos Java (incluye deps como JavaFX, CUP, Gson).

## Usado en
[[DataForge]], [[ConjAnalyzer]], [[CompScript]]

## Relacionadas
- [[JavaFX y Scene Builder]]
- [[CUP]]
- [[JFlex]]
