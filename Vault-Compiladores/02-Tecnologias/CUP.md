---
tags: [tecnologia, sintactico, java]
fuente: "Manual CUP (Constructor of Useful Parsers)"
fecha: 2026-07-10
---

# CUP

Generador de **analizadores sintácticos** [[Análisis sintáctico ascendente LR|LALR(1)]] para Java (equivalente Java de Yacc). Consume los `Symbol` de [[JFlex]] y ejecuta acciones semánticas al reducir.

## Estructura del `.cup`
```java
import java_cup.runtime.*;
parser code {: /* manejo de errores */ :};
scan with {: return lexer.next_token(); :};

terminal        PLUS, TIMES, LPAREN, RPAREN;
terminal Double NUMBER;
non terminal Object expr;

precedence left PLUS;
precedence left TIMES;

start with expr;

expr ::= expr:a PLUS expr:b {: RESULT = (Double)a + (Double)b; :}
       | NUMBER:n           {: RESULT = n; :}
       | LPAREN expr:e RPAREN {: RESULT = e; :}
       ;
```

- `:x` nombra el valor de un símbolo; `RESULT` = valor de la regla.
- `precedence left/right/nonassoc` resuelve [[Conflictos shift-reduce y reduce-reduce|conflictos]].
- **Comando manual:** `java -jar java-cup-11b.jar -parser Parser -symbols sym Sintaxis.cup` → genera `Parser.java` y `sym.java`.

## Integración con Maven (recomendado)

Mejor que el comando manual: el **cup-maven-plugin** genera el parser automáticamente en la fase `generate-sources` (ver [[Maven]]):

```xml
<plugin>
  <groupId>com.github.vbmacher</groupId>
  <artifactId>cup-maven-plugin</artifactId>
  <version>11b-20160615-3</version>
  <executions>
    <execution><goals><goal>generate</goal></goals></execution>
  </executions>
  <configuration>
    <className>Parser</className>
    <symbolsName>sym</symbolsName>
    <!-- lee src/main/cup/parser.cup → target/generated-sources/cup -->
  </configuration>
</plugin>
```

Y como dependencia de runtime (reemplaza el JAR suelto):
```xml
<dependency>
  <groupId>com.github.vbmacher</groupId>
  <artifactId>java-cup-runtime</artifactId>
  <version>11b-20160615-3</version>
</dependency>
```

## Usado en
[[DataForge]], [[ConjAnalyzer]], [[CompScript]]

## Relacionadas
- [[JFlex]]
- [[Ambigüedad, precedencia y asociatividad]]
- [[Cap 4 - Análisis sintáctico]]
- [[Codebase-Memory-MCP]] — código fuente real de cup-maven-plugin indexado como grafo consultable
