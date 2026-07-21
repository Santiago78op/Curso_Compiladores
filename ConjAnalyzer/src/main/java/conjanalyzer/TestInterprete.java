package conjanalyzer;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;

import conjanalyzer.interprete.Conjunto;
import conjanalyzer.interprete.Entorno;
import conjanalyzer.interprete.Interprete;
import conjanalyzer.interprete.Operacion;
import conjanalyzer.reportes.JsonSalida;
import conjanalyzer.reportes.Reportes;

/**
 * Prueba del INTERPRETE completo por consola (sin GUI): lexer → parser →
 * evaluacion → simplificacion → reportes. Muestra lo que la GUI enseña:
 * la consola de salida (EVALUAR), los conjuntos y operaciones definidos, los
 * errores y el JSON de simplificacion (5.4).
 *
 * Uso: ▶ Run en IDEA, o  mvn compile exec:java -Dexec.args="entradas/ejemplo1.ca"
 */
public class TestInterprete {

    public static void main(String[] args) throws Exception {
        String ruta = args.length > 0 ? args[0] : "entradas/ejemplo1.ca";
        String codigo = Files.readString(new File(ruta).toPath(), StandardCharsets.UTF_8);

        System.out.println("### Archivo: " + ruta + " ###\n");
        Entorno ent = Interprete.ejecutar(codigo);

        System.out.println("=== CONSOLA (EVALUAR 4.7) ===");
        System.out.print(ent.getConsola());

        System.out.println("=== CONJUNTOS DEFINIDOS ===");
        for (Conjunto conj : ent.getConjuntos().values()) {
            System.out.printf("  %s (%s) = %s%n", conj.nombre, conj.definicion,
                    Entorno.formatearConjunto(conj.elementos));
        }

        System.out.println("\n=== OPERACIONES ===");
        for (Operacion op : ent.getOperaciones().values()) {
            System.out.printf("  %s = %s%n", op.nombre, op.arbol.toPrefijo());
            System.out.printf("     resultado (%d): %s%n", op.resultado.size(),
                    Entorno.formatearConjunto(op.resultado));
            System.out.printf("     simplificacion: %s%n", op.simplificacion.seSimplifico
                    ? op.simplificacion.simplificado.toPrefijo() + "  leyes=" + op.simplificacion.leyes
                    : "No se puede simplificar la operacion");
        }

        if (!ent.getErrores().isEmpty()) {
            System.out.println("\n=== ERRORES ===");
            int e = 1;
            for (var err : ent.getErrores()) System.out.println("  " + (e++) + ". " + err);
        }

        System.out.println("\n=== JSON DE SIMPLIFICACION (5.4) ===");
        System.out.println(JsonSalida.construir(ent));

        // deja los reportes en disco para verificar que se generan
        File carpeta = new File("reportes");
        Reportes.generar(ent, carpeta);
        JsonSalida.generar(ent, carpeta);
        System.out.println("\n(reportes escritos en " + carpeta.getAbsolutePath() + ")");
    }
}
