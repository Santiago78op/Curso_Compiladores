package dataforge.interprete;

import java.io.StringReader;

import dataforge.analisis.Lexer;
import dataforge.analisis.Parser;

/**
 * Fachada del pipeline completo: código fuente (String) → Entorno ejecutado.
 * La GUI (y cualquier otro cliente) no necesita saber que adentro hay
 * un lexer JFlex y un parser CUP — le entrega texto y recibe resultados.
 */
public class Interprete {

    public static Entorno ejecutar(String codigo) {
        /* mismo pipeline de siempre, otra fuente: StringReader en vez de
           FileReader — al lexer le da igual de dónde vienen los caracteres */
        Lexer lexer = new Lexer(new StringReader(codigo));
        Parser parser = new Parser(lexer);
        lexer.entorno = parser.entorno;   // el lexer registra tokens y errores léxicos
        try {
            parser.parse();
        } catch (Exception e) {
            /* error sintáctico irrecuperable: syntax_error() ya lo registró
               en el entorno con su línea/columna; acá solo evitamos que la
               excepción mate a la GUI */
        }
        return parser.entorno;
    }
}
