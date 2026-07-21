package conjanalyzer.interprete;

import java.util.LinkedHashSet;
import java.util.Set;

/**
 * Una operacion definida con OPERA (4.6): su nombre, el arbol prefijo, el
 * conjunto resultante ya evaluado, los conjuntos base que referencia (para
 * el diagrama de Venn) y el resultado de intentar simplificarla (5.4 / 7).
 */
public class Operacion {

    public final String nombre;
    public final NodoOperacion arbol;
    public final Set<Character> resultado;           // caracteres del conjunto resultante
    public final LinkedHashSet<String> referencias;  // nombres de conjuntos base, en orden
    public final ResultadoSimplificacion simplificacion;
    public final int linea;
    public final int columna;

    public Operacion(String nombre, NodoOperacion arbol, Set<Character> resultado,
                     LinkedHashSet<String> referencias, ResultadoSimplificacion simplificacion,
                     int linea, int columna) {
        this.nombre = nombre;
        this.arbol = arbol;
        this.resultado = resultado;
        this.referencias = referencias;
        this.simplificacion = simplificacion;
        this.linea = linea;
        this.columna = columna;
    }
}
