package conjanalyzer.interprete;

import java.io.StringReader;

import conjanalyzer.analisis.Lexer;
import conjanalyzer.analisis.Parser;

/**
 * Fachada del pipeline completo: codigo fuente (String) → Entorno ejecutado.
 * La GUI (y cualquier otro cliente) no necesita saber que adentro hay un
 * lexer JFlex y un parser CUP — entrega texto y recibe resultados.
 */
public class Interprete {

    public static Entorno ejecutar(String codigo) {
        Lexer lexer = new Lexer(new StringReader(codigo));
        Parser parser = new Parser(lexer);
        lexer.entorno = parser.entorno;   // el lexer registra tokens y errores lexicos
        try {
            parser.parse();
        } catch (Exception e) {
            /* error sintactico irrecuperable: syntax_error() ya lo registro en
               el entorno con su linea/columna; aca solo evitamos que la
               excepcion mate a la GUI */
        }
        return parser.entorno;
    }
}
