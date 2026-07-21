package compscript.interprete;

/**
 * Un error del analisis/ejecucion para la tabla del reporte (6.2):
 * tipo (Lexico/Sintactico/Semantico), descripcion y posicion.
 * Linea y columna se guardan ya en base 1, listas para mostrar.
 */
public class RegistroError {

    public final String tipo;
    public final String descripcion;
    public final int linea;
    public final int columna;

    public RegistroError(String tipo, String descripcion, int linea, int columna) {
        this.tipo = tipo;
        this.descripcion = descripcion;
        this.linea = linea;
        this.columna = columna;
    }

    @Override
    public String toString() {
        return String.format("[%s] %s (linea %d, columna %d)",
                tipo, descripcion, linea, columna);
    }
}
