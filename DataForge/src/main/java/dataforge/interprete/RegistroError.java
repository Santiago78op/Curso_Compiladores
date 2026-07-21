package dataforge.interprete;

/**
 * Un error del análisis/ejecución, estructurado para la tabla del
 * reporte (§6.2): tipo (Léxico/Sintáctico/Semántico), descripción
 * y posición. Línea y columna se guardan ya en base 1, listas para mostrar.
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
        return String.format("[%s] %s (línea %d, columna %d)",
                tipo, descripcion, linea, columna);
    }
}
