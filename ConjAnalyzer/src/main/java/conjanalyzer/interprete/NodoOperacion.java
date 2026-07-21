package conjanalyzer.interprete;

import java.util.LinkedHashSet;
import java.util.Map;
import java.util.Set;

/**
 * Nodo del arbol de una operacion entre conjuntos (4.6).
 *
 *  - HOJA      : referencia a un conjunto ya definido con CONJ ({nombre}).
 *  - UNARIO    : complemento '^' (un solo hijo, en {@link #izq}).
 *  - BINARIO   : union 'U', interseccion '&' o diferencia '-' (izq, der).
 *
 * El operador va ANTES de sus operandos (notacion prefija/polaca): el
 * parser CUP construye este arbol al reducir cada produccion. Se recorre
 * en postorden para EVALUAR (obtener el Set<Character> resultante) y es la
 * estructura que reescribe el {@link Simplificador} (seccion 7).
 */
public class NodoOperacion {

    public final char op;                // 'U' '&' '-' '^'  o  '\0' si es hoja
    public final String nombreConj;      // solo hojas: nombre del CONJ referenciado
    public final NodoOperacion izq;      // hijo izquierdo / unico hijo del complemento
    public final NodoOperacion der;      // hijo derecho (null en unario y hoja)

    private NodoOperacion(char op, String nombreConj, NodoOperacion izq, NodoOperacion der) {
        this.op = op;
        this.nombreConj = nombreConj;
        this.izq = izq;
        this.der = der;
    }

    /* -------- fabricas (las usan las acciones del parser.cup) -------- */

    public static NodoOperacion hoja(String nombre) {
        return new NodoOperacion('\0', nombre, null, null);
    }

    public static NodoOperacion unario(char op, NodoOperacion hijo) {
        return new NodoOperacion(op, null, hijo, null);
    }

    public static NodoOperacion binario(char op, NodoOperacion a, NodoOperacion b) {
        return new NodoOperacion(op, null, a, b);
    }

    public boolean esHoja()   { return op == '\0'; }
    public boolean esUnario() { return op == '^'; }

    /* ================= evaluacion ================= */

    /**
     * Evalua el arbol contra los conjuntos base y el universo, devolviendo
     * el conjunto de caracteres resultante. Retorna null si alguna hoja
     * referencia un conjunto inexistente (el llamador ya reporto el error).
     */
    public Set<Character> evaluar(Map<String, Set<Character>> conjuntos, Set<Character> universo) {
        if (esHoja()) {
            Set<Character> base = conjuntos.get(nombreConj);
            return (base == null) ? null : new LinkedHashSet<>(base);
        }
        Set<Character> a = izq.evaluar(conjuntos, universo);
        if (a == null) return null;
        if (esUnario()) {                     // complemento: universo menos A
            Set<Character> r = new LinkedHashSet<>(universo);
            r.removeAll(a);
            return r;
        }
        Set<Character> b = der.evaluar(conjuntos, universo);
        if (b == null) return null;
        Set<Character> r = new LinkedHashSet<>(a);
        switch (op) {
            case 'U' -> r.addAll(b);          // union (6.1)
            case '&' -> r.retainAll(b);       // interseccion (6.2)
            case '-' -> r.removeAll(b);       // diferencia A - B (6.3)
        }
        return r;
    }

    /**
     * Evalua el arbol como funcion booleana sobre un "vector de region":
     * dado, por cada conjunto base, si un elemento hipotetico le pertenece,
     * dice si pertenece al resultado. Es exacto porque U/&/-/^ operan punto
     * a punto sobre las funciones indicadoras. Lo usa el diagrama de Venn
     * para saber que regiones sombrear sin recorrer todo el universo.
     */
    public boolean pertenenciaRegion(Map<String, Boolean> region) {
        if (esHoja()) return region.getOrDefault(nombreConj, false);
        if (esUnario()) return !izq.pertenenciaRegion(region);
        boolean a = izq.pertenenciaRegion(region);
        boolean b = der.pertenenciaRegion(region);
        return switch (op) {
            case 'U' -> a || b;
            case '&' -> a && b;
            case '-' -> a && !b;
            default  -> false;
        };
    }

    /* ================= utilidades ================= */

    /** Nombres de los conjuntos base referenciados (en orden de aparicion). */
    public void referencias(Set<String> acumulador) {
        if (esHoja()) { acumulador.add(nombreConj); return; }
        izq.referencias(acumulador);
        if (der != null) der.referencias(acumulador);
    }

    /** Serializa de vuelta a la notacion prefija del lenguaje: "U & {A} {B} {C}". */
    public String toPrefijo() {
        if (esHoja()) return "{" + nombreConj + "}";
        if (esUnario()) return "^ " + izq.toPrefijo();
        return op + " " + izq.toPrefijo() + " " + der.toPrefijo();
    }

    /** Copia profunda (el Simplificador trabaja sobre una copia, no muta el original). */
    public NodoOperacion copia() {
        if (esHoja()) return hoja(nombreConj);
        if (esUnario()) return unario(op, izq.copia());
        return binario(op, izq.copia(), der.copia());
    }

    @Override
    public String toString() { return toPrefijo(); }
}
