package compscript.interprete;

import java.util.LinkedHashMap;

/**
 * Un AMBITO de la pila de entornos. Cada bloque (funcion, if, ciclo...)
 * crea un Entorno hijo; la busqueda de una variable sube por la cadena de
 * padres (alcance estatico anidado). El estado global (consola, errores,
 * funciones, structs) vive en el Contexto compartido.
 */
public class Entorno {

    public final Contexto contexto;
    public final Entorno padre;       // null = entorno global
    public final String nombre;       // ambito, para el reporte 6.4
    private final LinkedHashMap<String, Simbolo> tabla = new LinkedHashMap<>();

    public Entorno(Contexto contexto, Entorno padre, String nombre) {
        this.contexto = contexto;
        this.padre = padre;
        this.nombre = nombre;
    }

    public Entorno crearHijo(String nombre) {
        return new Entorno(contexto, this, nombre);
    }

    /* ================= errores ================= */

    /** Registra el error semantico Y aborta la ejecucion (4.3). */
    public void errorSemantico(String descripcion, int linea, int columna) {
        contexto.error("Semantico", descripcion, linea, columna);
        throw new ErrorSemantico(descripcion, linea, columna);
    }

    public void imprimir(String texto) {
        contexto.consola.append(texto).append('\n');
    }

    /* ================= tabla de simbolos ================= */

    public void declarar(Simbolo s, int linea, int columna) {
        String clave = s.nombre.toLowerCase();
        if (tabla.containsKey(clave)) {
            errorSemantico("'" + s.nombre + "' ya fue declarada en este ambito", linea, columna);
        }
        tabla.put(clave, s);
        contexto.simbolos.add(s);   // registro plano para el reporte 6.4
    }

    /** Busca subiendo por los ambitos; null si no existe. */
    public Simbolo buscar(String id) {
        String clave = id.toLowerCase();
        for (Entorno e = this; e != null; e = e.padre) {
            Simbolo s = e.tabla.get(clave);
            if (s != null) return s;
        }
        return null;
    }

    /** Como buscar, pero registra error y aborta si no existe. */
    public Simbolo obtener(String id, int linea, int columna) {
        Simbolo s = buscar(id);
        if (s == null) {
            errorSemantico("la variable '" + id + "' no ha sido declarada", linea, columna);
        }
        return s;
    }

    /* ================= chequeo de tipos ================= */

    /** true si el valor v puede vivir en una variable declarada de tipo t. */
    public static boolean compatible(Tipo t, Valor v) {
        if (v == null) return false;
        if (v.tipo.cat == Tipo.Cat.NULL) return true;
        // Una lista/vector recien creado y vacio adopta el tipo declarado.
        if ((t.cat == Tipo.Cat.LIST || t.cat == Tipo.Cat.VECTOR)
                && (v.tipo.cat == t.cat) && v.lista().isEmpty()) return true;
        return t.equals(v.tipo);
    }

    public static String tipoDe(Valor v) {
        return v == null ? "null" : v.tipo.nombre();
    }
}
