package compscript.interprete;

import java.io.StringReader;
import java.util.List;

import compscript.analisis.Lexer;
import compscript.analisis.Parser;
import compscript.ast.A;

/**
 * Fachada del pipeline completo: codigo fuente (String) -> Contexto ejecutado.
 *
 * Ejecuta el AST en 3 PASADAS (recomendacion del enunciado 5.23):
 *   1a) registra funciones, metodos y structs (permite forward-reference)
 *   2a) ejecuta declaraciones y asignaciones globales
 *   3a) ejecuta el RUN_MAIN (punto de entrada)
 *
 * Recuperacion de errores (4.3): los errores lexicos y sintacticos se
 * ACUMULAN durante el analisis (el lexer descarta el caracter intruso; el
 * parser usa modo panico hasta el ';'). Un error SEMANTICO en tiempo de
 * ejecucion lanza ErrorSemantico, que se captura aqui: el error ya quedo
 * en la tabla y la ejecucion TERMINA ordenadamente (no sigue corriendo
 * instrucciones), a diferencia de DataForge.
 */
public class Interprete {

    public static Contexto ejecutar(String codigo) {
        Contexto ctx = new Contexto();

        Lexer lexer = new Lexer(new StringReader(codigo));
        lexer.contexto = ctx;
        Parser parser = new Parser(lexer);
        parser.contexto = ctx;
        try {
            parser.parse();
        } catch (Exception e) {
            /* error sintactico irrecuperable: syntax_error() ya lo registro */
        }

        @SuppressWarnings("unchecked")
        List<A.Instruccion> raiz = (List<A.Instruccion>) (List<?>) parser.raiz;
        ctx.raiz = new A.Programa(raiz);
        ctx.global = new Entorno(ctx, null, "global");

        int erroresAntesPass1 = ctx.errores.size();
        try {
            // 1a pasada: registro de funciones/metodos/structs
            for (A.Instruccion i : raiz) {
                if (i instanceof A.DeclaracionFuncion f) ctx.registrarFuncion(f);
                else if (i instanceof A.DeclaracionStruct s) ctx.registrarStruct(s);
            }
            boolean registroOk = ctx.errores.size() == erroresAntesPass1;

            if (registroOk) {
                // 2a pasada: declaraciones y asignaciones globales
                for (A.Instruccion i : raiz) {
                    if (i instanceof A.DeclaracionFuncion
                            || i instanceof A.DeclaracionStruct
                            || i instanceof A.RunMain) continue;
                    i.ejecutar(ctx.global);
                }
                // 3a pasada: punto de entrada
                for (A.Instruccion i : raiz) {
                    if (i instanceof A.RunMain) i.ejecutar(ctx.global);
                }
            }
        } catch (ErrorSemantico ex) {
            /* error ya registrado; ejecucion terminada ordenadamente (4.3) */
        } catch (Senales.Break b) {
            ctx.error("Semantico", "'break' fuera de un ciclo", b.linea, b.columna);
        } catch (Senales.Continue c) {
            ctx.error("Semantico", "'continue' fuera de un ciclo", c.linea, c.columna);
        } catch (Senales.Retorno r) {
            ctx.error("Semantico", "'return' fuera de una funcion o metodo", r.linea, r.columna);
        } catch (StackOverflowError so) {
            ctx.error("Semantico", "desbordamiento de pila (posible recursion infinita)", 0, 0);
        }
        return ctx;
    }
}
