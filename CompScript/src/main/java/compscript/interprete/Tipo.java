package compscript.interprete;

/**
 * Tipo de dato de CompScript. Cubre los primitivos (5.3), los tipos
 * compuestos (vectores 1D/2D, listas) y los structs. Se usa tanto en el
 * chequeo semantico (comparar tipos) como en la tabla de simbolos (6.4).
 */
public class Tipo {

    public enum Cat { INT, DOUBLE, BOOL, CHAR, STRING, VOID, NULL, VECTOR, LIST, STRUCT }

    public final Cat cat;
    public final Tipo elemento;    // base de VECTOR / LIST
    public final int dimensiones;  // VECTOR: 1 o 2; resto: 0
    public final String structName; // STRUCT

    private Tipo(Cat cat, Tipo elemento, int dimensiones, String structName) {
        this.cat = cat;
        this.elemento = elemento;
        this.dimensiones = dimensiones;
        this.structName = structName;
    }

    /* ---- primitivos ---- */
    public static final Tipo INT    = new Tipo(Cat.INT, null, 0, null);
    public static final Tipo DOUBLE = new Tipo(Cat.DOUBLE, null, 0, null);
    public static final Tipo BOOL   = new Tipo(Cat.BOOL, null, 0, null);
    public static final Tipo CHAR   = new Tipo(Cat.CHAR, null, 0, null);
    public static final Tipo STRING = new Tipo(Cat.STRING, null, 0, null);
    public static final Tipo VOID   = new Tipo(Cat.VOID, null, 0, null);
    public static final Tipo NULL   = new Tipo(Cat.NULL, null, 0, null);

    public static Tipo vector(Tipo base, int dim) {
        return new Tipo(Cat.VECTOR, base, dim, null);
    }
    public static Tipo lista(Tipo base) {
        return new Tipo(Cat.LIST, base, 0, null);
    }
    public static Tipo struct(String nombre) {
        return new Tipo(Cat.STRUCT, null, 0, nombre.toLowerCase());
    }

    public boolean esNumerico() {
        return cat == Cat.INT || cat == Cat.DOUBLE || cat == Cat.CHAR;
    }

    /** Igualdad ESTRUCTURAL (para el chequeo de tipos). */
    @Override
    public boolean equals(Object o) {
        if (!(o instanceof Tipo t)) return false;
        if (cat != t.cat) return false;
        switch (cat) {
            case VECTOR:
                return dimensiones == t.dimensiones && java.util.Objects.equals(elemento, t.elemento);
            case LIST:
                return java.util.Objects.equals(elemento, t.elemento);
            case STRUCT:
                return java.util.Objects.equals(structName, t.structName);
            default:
                return true;
        }
    }

    @Override
    public int hashCode() { return cat.hashCode(); }

    /** Nombre legible en espanol para errores y reportes. */
    public String nombre() {
        switch (cat) {
            case INT:    return "Entero";
            case DOUBLE: return "Decimal";
            case BOOL:   return "Booleano";
            case CHAR:   return "Caracter";
            case STRING: return "Cadena";
            case VOID:   return "void";
            case NULL:   return "null";
            case VECTOR: return "Vector<" + (elemento == null ? "?" : elemento.nombre())
                    + ">" + (dimensiones == 2 ? "[][]" : "[]");
            case LIST:   return "Lista<" + (elemento == null ? "?" : elemento.nombre()) + ">";
            case STRUCT: return "Struct " + structName;
            default:     return "?";
        }
    }

    @Override
    public String toString() { return nombre(); }
}
