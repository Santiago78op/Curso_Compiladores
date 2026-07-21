package compscript;

import java.nio.file.Files;
import java.nio.file.Path;

import compscript.interprete.Contexto;
import compscript.interprete.Interprete;
import compscript.interprete.RegistroError;
import compscript.interprete.Simbolo;

/**
 * Prueba del INTERPRETE completo por consola (lexer -> parser -> AST ->
 * 3 pasadas -> ejecucion). Muestra lo que en la GUI iria a la consola, la
 * tabla de simbolos (6.4) y los errores (6.2).
 *
 * Uso: mvn compile exec:java -Dexec.args="entradas/ejemplo1.cs"
 *      (sin args corre entradas/ejemplo1.cs)
 */
public class TestInterprete {

    public static void main(String[] args) throws Exception {
        String ruta = args.length > 0 ? args[0] : "entradas/ejemplo1.cs";
        String codigo = Files.readString(Path.of(ruta));

        System.out.println("============================================================");
        System.out.println(" ARCHIVO: " + ruta);
        System.out.println("============================================================");

        Contexto ctx = Interprete.ejecutar(codigo);

        System.out.println("--- CONSOLA ---");
        System.out.print(ctx.consola);

        System.out.println("\n--- TABLA DE SIMBOLOS ---");
        System.out.printf("%-3s %-14s %-9s %-16s %-10s %-22s %-5s %-3s%n",
                "#", "Id", "Tipo", "TipoDato", "Entorno", "Valor", "Lin", "Col");
        int n = 1;
        for (Simbolo s : ctx.simbolos) {
            System.out.printf("%-3d %-14s %-9s %-16s %-10s %-22s %-5d %-3d%n",
                    n++, s.nombre, s.categoria, s.tipo.nombre(), s.ambito,
                    s.valor == null ? "null" : s.valor.reporte(), s.linea, s.columna);
        }

        System.out.println("\n--- ERRORES (" + ctx.errores.size() + ") ---");
        int e = 1;
        for (RegistroError err : ctx.errores) System.out.println((e++) + ". " + err);

        System.out.println("\n--- TOKENS: " + ctx.tokens.size() + " reconocidos ---");
    }
}
