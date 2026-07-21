package dataforge;

import java.io.FileReader;

import dataforge.analisis.Lexer;
import dataforge.analisis.Parser;

/**
 * Prueba del analizador sintáctico (Etapa 2): lee un archivo .df y lo
 * valida contra la gramática. Todavía no ejecuta nada — solo responde
 * la pregunta "¿este programa está bien escrito?".
 *
 * Uso en IDEA: ▶ Run sobre esta clase.
 * Terminal:    mvn compile exec:java
 *              mvn compile exec:java -Dexec.args="entradas/ejemplo2.df"
 */
public class TestParser {

    public static void main(String[] args) {
        String ruta = args.length > 0 ? args[0] : "entradas/ejemplo2.df";

        try {
            /* el pipeline completo: el parser PIDE tokens al lexer
               uno por uno (next_token) mientras valida la estructura */
            Parser parser = new Parser(new Lexer(new FileReader(ruta)));
            parser.parse();
            System.out.println("[OK] Análisis sintáctico exitoso: '" + ruta + "' es un programa válido.");
        } catch (Exception e) {
            System.err.println("[X] El análisis se detuvo: el programa NO es válido.");
        }
    }
}
