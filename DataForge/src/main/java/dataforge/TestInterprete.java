package dataforge;

import java.io.FileReader;

import dataforge.analisis.Lexer;
import dataforge.analisis.Parser;
import dataforge.interprete.Entorno;
import dataforge.interprete.Grafica;
import dataforge.interprete.Simbolo;

/**
 * Prueba del INTÉRPRETE completo (Etapa 3): lexer → parser → ejecución.
 * Muestra lo que en la Etapa 4 irá a la GUI: la consola de salida,
 * la tabla de símbolos (reporte §6.3), los errores y las gráficas
 * registradas (que la Etapa 5 dibujará).
 *
 * Uso en IDEA: ▶ Run. Otro archivo: Run → Edit Configurations → args.
 */
public class TestInterprete {

    public static void main(String[] args) throws Exception {
        String ruta = args.length > 0 ? args[0] : "entradas/ejemplo3_errores.df";

        Entorno ent;
        /* try-with-resources: garantiza que el archivo se cierre aunque
           el parser aborte antes de llegar a EOF (el cierre automático
           del lexer generado por JFlex solo ocurre al leer hasta el final). */
        try (FileReader lector = new FileReader(ruta)) {
            Parser parser = new Parser(new Lexer(lector));
            ent = parser.entorno;
            try {
                parser.parse();
            } catch (Exception e) {
                System.err.println("(el análisis se detuvo por un error sintáctico)");
            }
        }

        System.out.println("=== CONSOLA ===");
        System.out.print(ent.getConsola());

        System.out.println("\n=== TABLA DE SÍMBOLOS ===");
        System.out.printf("%-3s %-14s %-9s %-7s %-26s %-5s %-3s%n",
                "#", "Nombre", "Categoría", "Tipo", "Valor", "Línea", "Col");
        int n = 1;
        for (Simbolo s : ent.getSimbolos().values()) {
            System.out.printf("%-3d %-14s %-9s %-7s %-26s %-5d %-3d%n",
                    n++, s.nombre, s.categoria, s.tipo,
                    Entorno.valorReporte(s.valor), s.linea, s.columna);
        }

        if (!ent.getErrores().isEmpty()) {
            System.out.println("\n=== ERRORES ===");
            int e = 1;
            for (var err : ent.getErrores())
                System.out.println((e++) + ". " + err);
        }

        if (!ent.getGraficas().isEmpty()) {
            System.out.println("\n=== GRÁFICAS REGISTRADAS (se dibujan en la Etapa 5) ===");
            for (Grafica g : ent.getGraficas())
                System.out.println("· " + g);
        }
    }
}
