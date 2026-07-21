package conjanalyzer.interprete;

import java.util.Set;

/**
 * Un conjunto definido con CONJ (4.5): su nombre, los caracteres que
 * contiene y una descripcion legible de como fue declarado (para reportes).
 */
public class Conjunto {

    public final String nombre;
    public final Set<Character> elementos;
    public final String definicion;   // p.ej. "a~z" o "1, 2, 3, a, b"
    public final int linea;
    public final int columna;

    public Conjunto(String nombre, Set<Character> elementos, String definicion, int linea, int columna) {
        this.nombre = nombre;
        this.elementos = elementos;
        this.definicion = definicion;
        this.linea = linea;
        this.columna = columna;
    }
}
