package compscript.interprete;

/**
 * Una entrada de la tabla de simbolos (reporte 6.4). Guarda todo lo que
 * pide el enunciado: nombre, categoria (Variable/Vector/Lista/Struct),
 * tipo de dato, entorno (ambito) donde se declaro, valor y posicion.
 *
 * El valor es MUTABLE: la asignacion (5.12) actualiza este campo.
 */
public class Simbolo {

    public final String nombre;     // como lo escribio el usuario
    public final String categoria;  // "Variable" | "Vector" | "Lista" | "Struct"
    public final Tipo tipo;
    public final boolean mutable;   // let = true, const = false
    public final String ambito;     // entorno donde vive (global, nombre de funcion...)
    public final int linea;
    public final int columna;
    public Valor valor;

    public Simbolo(String nombre, String categoria, Tipo tipo, boolean mutable,
                   String ambito, Valor valor, int linea, int columna) {
        this.nombre = nombre;
        this.categoria = categoria;
        this.tipo = tipo;
        this.mutable = mutable;
        this.ambito = ambito;
        this.valor = valor;
        this.linea = linea;
        this.columna = columna;
    }

    public static String categoriaDe(Tipo t) {
        switch (t.cat) {
            case VECTOR: return "Vector";
            case LIST:   return "Lista";
            case STRUCT: return "Struct";
            default:     return "Variable";
        }
    }
}
