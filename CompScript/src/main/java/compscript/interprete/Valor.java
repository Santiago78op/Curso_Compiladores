package compscript.interprete;

import java.util.List;
import java.util.Map;

/**
 * Un valor en tiempo de ejecucion: lleva SU tipo y el objeto Java que lo
 * representa. La ejecucion (tree-walk sobre el AST) produce y consume Valor.
 *
 *   INT    -> Integer            CHAR   -> Character
 *   DOUBLE -> Double             STRING -> String
 *   BOOL   -> Boolean            VOID/NULL -> null
 *   VECTOR / LIST -> java.util.List<Valor>
 *   STRUCT -> java.util.LinkedHashMap<String,Valor> (claves = nombre de campo, orden de declaracion)
 */
public class Valor {

    public final Tipo tipo;
    public final Object valor;

    public Valor(Tipo tipo, Object valor) {
        this.tipo = tipo;
        this.valor = valor;
    }

    /* ---- fabricas ---- */
    public static Valor vInt(int v)       { return new Valor(Tipo.INT, v); }
    public static Valor vDouble(double v) { return new Valor(Tipo.DOUBLE, v); }
    public static Valor vBool(boolean v)  { return new Valor(Tipo.BOOL, v); }
    public static Valor vChar(char v)     { return new Valor(Tipo.CHAR, v); }
    public static Valor vString(String v) { return new Valor(Tipo.STRING, v); }
    public static final Valor VOID = new Valor(Tipo.VOID, null);

    /* ---- acceso numerico uniforme (char -> codigo ascii, bool -> 1/0) ---- */
    public double numero() {
        switch (tipo.cat) {
            case INT:    return (Integer) valor;
            case DOUBLE: return (Double) valor;
            case CHAR:   return (int) (Character) valor;
            case BOOL:   return ((Boolean) valor) ? 1 : 0;
            default:     return 0;
        }
    }

    @SuppressWarnings("unchecked")
    public List<Valor> lista() { return (List<Valor>) valor; }

    @SuppressWarnings("unchecked")
    public Map<String, Valor> campos() { return (Map<String, Valor>) valor; }

    /* ================= formato para console.log (5.24) ================= */

    public String texto() {
        switch (tipo.cat) {
            case INT:    return Integer.toString((Integer) valor);
            case DOUBLE: return dbl((Double) valor);
            case BOOL:   return ((Boolean) valor) ? "true" : "false";
            case CHAR:   return String.valueOf((Character) valor);
            case STRING: return (String) valor;
            case VECTOR:
            case LIST: {
                StringBuilder sb = new StringBuilder("[");
                List<Valor> l = lista();
                for (int i = 0; i < l.size(); i++) {
                    if (i > 0) sb.append(", ");
                    sb.append(l.get(i).texto());
                }
                return sb.append("]").toString();
            }
            case STRUCT: return textoStruct();
            default:     return "null";
        }
    }

    /** toString de un struct (5.27): NombreStruct { campo: valor, ... }
     *  con las cadenas entre comillas dobles y los char entre comillas simples. */
    public String textoStruct() {
        StringBuilder sb = new StringBuilder(tipo.structName == null ? "struct" : tipo.structName);
        sb.append(" { ");
        boolean primero = true;
        for (Map.Entry<String, Valor> e : campos().entrySet()) {
            if (!primero) sb.append(", ");
            primero = false;
            sb.append(e.getKey()).append(": ").append(e.getValue().campoTexto());
        }
        return sb.append(" }").toString();
    }

    /** Valor de un campo dentro de un struct: cadenas y char con comillas. */
    private String campoTexto() {
        if (tipo.cat == Tipo.Cat.STRING) return "\"" + valor + "\"";
        if (tipo.cat == Tipo.Cat.CHAR)   return "'" + valor + "'";
        return texto();
    }

    /** Formato para la tabla de simbolos (6.4): cadenas con comillas. */
    public String reporte() {
        if (tipo.cat == Tipo.Cat.STRING) return "\"" + valor + "\"";
        if (tipo.cat == Tipo.Cat.CHAR)   return "'" + valor + "'";
        return texto();
    }

    /** 16.0 se muestra "16.0"; 3.281 queda "3.281". */
    private static String dbl(double d) {
        if (d == Math.rint(d) && !Double.isInfinite(d)) {
            return (long) d + ".0";
        }
        return Double.toString(d);
    }

    @Override
    public String toString() { return texto(); }
}
