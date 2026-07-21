package conjanalyzer.interprete;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * El ENTORNO de ejecucion de ConjAnalyzer: todo el estado del programa que
 * corre.
 *   - universo (caracteres ASCII 33..126, seccion 4.4)
 *   - conjuntos definidos con CONJ (4.5)
 *   - operaciones definidas con OPERA (4.6), ya evaluadas y simplificadas
 *   - salida de consola (lo que ve el usuario: resultados de EVALUAR, 4.7)
 *   - errores acumulados (NO detienen la ejecucion)
 *   - tokens (los registra el lexer, para el reporte 5.2)
 *
 * Las acciones {: ... :} del parser.cup llaman a estos metodos: el parser
 * VALIDA la forma, el entorno le da SIGNIFICADO.
 *
 * IMPORTANTE (case sensitive, 4.1): las claves NO se normalizan —
 * "conjuntoA" y "conjuntoa" son conjuntos distintos.
 *
 * Se crea un Entorno FRESCO por ejecucion (seccion 5: los reportes son solo
 * del ultimo analisis, no se arrastra estado de corridas anteriores).
 */
public class Entorno {

    /** Universo: simbolos ASCII 33 ('!') a 126 ('~') inclusive (4.4). */
    public static final int UNIVERSO_MIN = 33;
    public static final int UNIVERSO_MAX = 126;

    private final Set<Character> universo = new LinkedHashSet<>();
    private final LinkedHashMap<String, Conjunto> conjuntos = new LinkedHashMap<>();
    private final LinkedHashMap<String, Operacion> operaciones = new LinkedHashMap<>();
    private final StringBuilder consola = new StringBuilder();
    private final List<RegistroError> errores = new ArrayList<>();
    private final List<Object[]> tokens = new ArrayList<>();  // {lexema, tipo, linea, col}

    public Entorno() {
        for (int i = UNIVERSO_MIN; i <= UNIVERSO_MAX; i++) universo.add((char) i);
    }

    /* ================= errores y tokens ================= */

    public void error(String tipo, String descripcion, int linea, int columna) {
        errores.add(new RegistroError(tipo, descripcion, linea + 1, columna + 1));
    }

    /** El lexer llama esto por cada token (reporte de tokens 5.2). */
    public void registrarToken(String lexema, int tipo, int linea, int columna) {
        tokens.add(new Object[]{ lexema, tipo, linea + 1, columna + 1 });
    }

    public List<Object[]> getTokens() { return tokens; }

    /* ================= definicion de conjuntos (4.5) ================= */

    /** CONJ : id -> a~b;  → todos los caracteres entre a y b inclusive. */
    public void definirConjuntoRango(String id, String a, String b, int l, int c) {
        if (existe(id, l, c)) return;
        Character ini = charUnico(a, l, c);
        Character fin = charUnico(b, l, c);
        if (ini == null || fin == null) return;
        if (ini > fin) {
            error("Semantico", "rango invalido en '" + id + "': '" + a
                    + "' debe ser menor o igual que '" + b + "'", l, c);
            return;
        }
        LinkedHashSet<Character> set = new LinkedHashSet<>();
        for (char ch = ini; ch <= fin; ch++) set.add(ch);
        conjuntos.put(id, new Conjunto(id, set, a + "~" + b, l + 1, c + 1));
    }

    /** CONJ : id -> e1, e2, ...;  → cada elemento es un unico caracter. */
    public void definirConjuntoLista(String id, ArrayList<String> elementos, int l, int c) {
        if (existe(id, l, c)) return;
        LinkedHashSet<Character> set = new LinkedHashSet<>();
        for (String elem : elementos) {
            Character ch = charUnico(elem, l, c);
            if (ch == null) return;   // el error ya quedo registrado
            set.add(ch);
        }
        conjuntos.put(id, new Conjunto(id, set, String.join(", ", elementos), l + 1, c + 1));
    }

    /* ================= definicion de operaciones (4.6) ================= */

    /** OPERA : id -> operacion;  valida referencias, evalua y simplifica. */
    public void definirOperacion(String id, NodoOperacion arbol, int l, int c) {
        if (existeOperacion(id, l, c)) return;

        LinkedHashSet<String> refs = new LinkedHashSet<>();
        arbol.referencias(refs);
        for (String ref : refs) {
            if (!conjuntos.containsKey(ref)) {
                error("Semantico", "la operacion '" + id + "' referencia el conjunto '"
                        + ref + "', que no ha sido definido", l, c);
                return;
            }
        }

        Map<String, Set<Character>> mapa = new LinkedHashMap<>();
        for (var e : conjuntos.entrySet()) mapa.put(e.getKey(), e.getValue().elementos);

        Set<Character> resultado = arbol.evaluar(mapa, universo);
        if (resultado == null) resultado = new LinkedHashSet<>();

        ResultadoSimplificacion simpl = new Simplificador().simplificar(arbol);

        operaciones.put(id, new Operacion(id, arbol, resultado, refs, simpl, l + 1, c + 1));
    }

    /* ================= evaluar pertenencia (4.7) ================= */

    /**
     * EVALUAR ( {e1, e2, ...}, operacion );  Evalua cada elemento contra el
     * conjunto RESULTANTE de la operacion y escribe en consola el formato
     * exacto de la seccion 4.8 (exitoso / fallo por elemento).
     *
     * Nota de diseno: la seccion 4.7 dice explicitamente "validar que cierto
     * conjunto de datos pertenezca al conjunto resultante", asi que la
     * pertenencia se mide contra el conjunto ya evaluado. (El segundo bloque
     * de salida del ejemplo 4.8 del enunciado muestra "1 -> exitoso" para una
     * interseccion que no contiene al 1: eso es una inconsistencia del
     * enunciado; aca se respeta la semantica de conjunto resultante.)
     */
    public void evaluar(ArrayList<String> datos, String operacion, int l, int c) {
        Operacion op = operaciones.get(operacion);
        if (op == null) {
            error("Semantico", "la operacion '" + operacion + "' no ha sido definida", l, c);
            return;
        }
        consola.append("===============\n")
               .append("Evaluar: ").append(operacion).append('\n')
               .append("===============\n");
        for (String dato : datos) {
            Character ch = charUnico(dato, l, c);
            if (ch == null) continue;   // elemento invalido: ya se reporto el error
            boolean pertenece = op.resultado.contains(ch);
            consola.append(ch).append(" -> ").append(pertenece ? "exitoso" : "fallo").append('\n');
        }
        consola.append('\n');
    }

    /* ================= helpers ================= */

    /** Convierte un lexema de elemento en un unico Character validado. */
    private Character charUnico(String lexema, int l, int c) {
        if (lexema == null || lexema.length() != 1) {
            error("Semantico", "'" + lexema + "' no es un elemento valido: "
                    + "cada elemento debe ser un unico caracter", l, c);
            return null;
        }
        char ch = lexema.charAt(0);
        if (ch < UNIVERSO_MIN || ch > UNIVERSO_MAX) {
            error("Semantico", "el elemento '" + ch + "' no pertenece al universo (ASCII 33..126)", l, c);
            return null;
        }
        return ch;
    }

    private boolean existe(String id, int l, int c) {
        if (conjuntos.containsKey(id)) {
            error("Semantico", "el conjunto '" + id + "' ya fue definido (linea "
                    + conjuntos.get(id).linea + ")", l, c);
            return true;
        }
        return false;
    }

    private boolean existeOperacion(String id, int l, int c) {
        if (operaciones.containsKey(id)) {
            error("Semantico", "la operacion '" + id + "' ya fue definida (linea "
                    + operaciones.get(id).linea + ")", l, c);
            return true;
        }
        return false;
    }

    /** Formatea un conjunto de caracteres como {a, b, c} para mostrar. */
    public static String formatearConjunto(Set<Character> s) {
        StringBuilder b = new StringBuilder("{");
        int i = 0;
        for (Character ch : s) {
            if (i++ > 0) b.append(", ");
            b.append(ch);
        }
        return b.append("}").toString();
    }

    /* ================= getters para reportes / GUI ================= */

    public String getConsola() { return consola.toString(); }
    public List<RegistroError> getErrores() { return errores; }
    public LinkedHashMap<String, Conjunto> getConjuntos() { return conjuntos; }
    public LinkedHashMap<String, Operacion> getOperaciones() { return operaciones; }
}
