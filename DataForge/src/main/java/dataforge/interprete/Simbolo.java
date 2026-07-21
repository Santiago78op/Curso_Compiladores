package dataforge.interprete;

/**
 * Una entrada de la tabla de símbolos: variable o arreglo declarado.
 * Guarda todo lo que pide el reporte del enunciado (§6.3):
 * nombre, categoría, tipo, valor y posición de declaración.
 */
public class Simbolo {

    public final String nombre;     // como lo escribió el usuario (para mostrar)
    public final String categoria;  // "variable" | "arreglo"
    public final String tipo;       // "double" | "char[]"
    public final Object valor;      // Double, String o ArrayList<Object>
    public final int linea;         // 1-based, listo para mostrar
    public final int columna;

    public Simbolo(String nombre, String categoria, String tipo,
                   Object valor, int linea, int columna) {
        this.nombre = nombre;
        this.categoria = categoria;
        this.tipo = tipo;
        this.valor = valor;
        this.linea = linea;
        this.columna = columna;
    }
}
