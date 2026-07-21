package dataforge;

import java.io.FileReader;
import java.lang.reflect.Field;

import java_cup.runtime.Symbol;

import dataforge.analisis.Lexer;
import dataforge.analisis.sym;

/**
 * Prueba del analizador léxico (Etapa 1): lee un archivo .df y muestra
 * en consola la tabla de tokens reconocidos — el embrión del
 * "Reporte de Tokens" que pide el enunciado (sección 6.1).
 *
 * Uso:  mvn compile exec:java
 *       mvn compile exec:java -Dexec.args="entradas/otro.df"
 */
public class TestLexer {

    public static void main(String[] args) throws Exception {
        String ruta = args.length > 0 ? args[0] : "entradas/ejemplo1.df";

        Lexer lexer = new Lexer(new FileReader(ruta));

        System.out.printf("%-4s %-22s %-18s %-6s %-7s%n",
                "#", "Lexema", "Token", "Línea", "Columna");
        System.out.println("-".repeat(62));

        int n = 1;
        Symbol s;
        while ((s = lexer.next_token()).sym != sym.EOF) {
            /* Symbol guarda: sym = tipo de token, left = línea, right = columna,
               value = el lexema (o el valor ya convertido, ej. Double) */
            System.out.printf("%-4d %-22s %-18s %-6d %-7d%n",
                    n++, s.value, nombreToken(s.sym), s.left + 1, s.right + 1);
        }
    }

    /** Busca por reflexión el nombre de la constante en sym.java
     *  (así no mantenemos un switch gigante a mano). */
    private static String nombreToken(int id) throws Exception {
        for (Field f : sym.class.getFields()) {
            if (f.getType() == int.class && f.getInt(null) == id) {
                return f.getName();
            }
        }
        return "?";
    }
}
