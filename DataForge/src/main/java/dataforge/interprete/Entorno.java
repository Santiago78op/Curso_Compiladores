package dataforge.interprete;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

/**
 * El ENTORNO de ejecución: todo el estado del programa que corre.
 * - tabla de símbolos (variables y arreglos)
 * - salida de consola (lo que verá el usuario)
 * - errores semánticos acumulados (NO detienen la ejecución)
 * - gráficas registradas con EXEC (se dibujan en la Etapa 5)
 *
 * Las acciones {: ... :} del parser.cup llaman a estos métodos:
 * el parser VALIDA la forma, el entorno le da SIGNIFICADO.
 */
public class Entorno {

    private final LinkedHashMap<String, Simbolo> simbolos = new LinkedHashMap<>();
    private final StringBuilder consola = new StringBuilder();
    private final List<RegistroError> errores = new ArrayList<>();
    private final List<Grafica> graficas = new ArrayList<>();
    private final List<Object[]> tokens = new ArrayList<>();  // {lexema, tipo, línea, col} para el reporte §6.1

    /* atributos válidos de los bloques de gráfica (decisión de diseño nº 2:
       son ID en la gramática y se validan acá, en la ejecución) */
    private static final List<String> ATTRS_VALIDOS =
            List.of("titulo", "ejex", "ejey", "titulox", "tituloy", "values", "label");

    /* ================= errores ================= */

    public void error(String tipo, String descripcion, int linea, int columna) {
        errores.add(new RegistroError(tipo, descripcion, linea + 1, columna + 1));
    }

    /* ================= tokens (los registra el lexer, reporte §6.1) ================= */

    /** El lexer llama esto por cada token: guarda el LEXEMA ORIGINAL
     *  (yytext), no el valor convertido — así "1" se reporta como 1, no 1.0. */
    public void registrarToken(String lexema, int tipo, int linea, int columna) {
        tokens.add(new Object[]{ lexema, tipo, linea + 1, columna + 1 });
    }

    public List<Object[]> getTokens() { return tokens; }

    /* ================= tabla de símbolos ================= */
    /* el lenguaje es case insensitive → la CLAVE se normaliza a minúsculas,
       pero se conserva el nombre original para los reportes */

    public void declararVariable(String id, String tipo, Object valor, int l, int c) {
        String clave = id.toLowerCase();
        if (simbolos.containsKey(clave)) {
            error("Semántico", "'" + id + "' ya fue declarada (línea "
                    + simbolos.get(clave).linea + ")", l, c);
            return;
        }
        if (valor == null) return;  // la expresión ya reportó su propio error
        boolean ok = (tipo.equals("double") && valor instanceof Double)
                  || (tipo.equals("char[]") && valor instanceof String);
        if (!ok) {
            error("Semántico", "no se puede asignar " + nombreTipo(valor)
                    + " a la variable '" + id + "' de tipo " + tipo, l, c);
            return;
        }
        simbolos.put(clave, new Simbolo(id, "variable", tipo, valor, l + 1, c + 1));
    }

    public void declararArreglo(String idArr, String tipo, ArrayList<Object> valores, int l, int c) {
        String clave = idArr.toLowerCase();
        if (simbolos.containsKey(clave)) {
            error("Semántico", "'" + idArr + "' ya fue declarado (línea "
                    + simbolos.get(clave).linea + ")", l, c);
            return;
        }
        for (Object v : valores) {
            if (v == null) return;  // error ya reportado río arriba
            boolean ok = (tipo.equals("double") && v instanceof Double)
                      || (tipo.equals("char[]") && v instanceof String);
            if (!ok) {
                error("Semántico", "el arreglo '" + idArr + "' es de tipo " + tipo
                        + " pero contiene " + nombreTipo(v), l, c);
                return;
            }
        }
        simbolos.put(clave, new Simbolo(idArr, "arreglo", tipo, valores, l + 1, c + 1));
    }

    /** Valor de una variable; null (+ error registrado) si no existe. */
    public Object valorDe(String id, int l, int c) {
        Simbolo s = simbolos.get(id.toLowerCase());
        if (s == null || !s.categoria.equals("variable")) {
            error("Semántico", "la variable '" + id + "' no ha sido declarada", l, c);
            return null;
        }
        return s.valor;
    }

    /** Valor de un arreglo por nombre (@id); null si no existe. */
    @SuppressWarnings("unchecked")
    public ArrayList<Object> valorArreglo(String idArr, int l, int c) {
        Simbolo s = simbolos.get(idArr.toLowerCase());
        if (s == null || !s.categoria.equals("arreglo")) {
            error("Semántico", "el arreglo '" + idArr + "' no ha sido declarado", l, c);
            return null;
        }
        return (ArrayList<Object>) s.valor;
    }

    /* ================= consola (5.9) ================= */

    /** console::print = e1, e2, ... → una línea, separada por comas. */
    public void imprimir(ArrayList<Object> exprs) {
        StringBuilder linea = new StringBuilder();
        for (int i = 0; i < exprs.size(); i++) {
            if (i > 0) linea.append(", ");
            linea.append(formatear(exprs.get(i)));
        }
        consola.append(linea).append('\n');
    }

    /** console::column = titulo -> arreglo → tabla de una columna (5.9.2). */
    public void imprimirColumna(Object titulo, ArrayList<Object> arr) {
        if (arr == null) return;
        consola.append("--------------\n")
               .append(formatear(titulo)).append('\n')
               .append("--------------\n");
        for (Object v : arr) consola.append(formatear(v)).append('\n');
    }

    /** Formato de CONSOLA: 15.0 se muestra "15"; 15.7 queda "15.7";
     *  las cadenas van sin comillas (console::print = "hola" → hola). */
    public static String formatear(Object v) {
        if (v == null) return "null";
        if (v instanceof Double d) {
            if (d == Math.floor(d) && !d.isInfinite()) return String.valueOf(d.longValue());
            return String.valueOf(d);
        }
        return String.valueOf(v);
    }

    /** Formato de REPORTE (§6.3): como el ejemplo del enunciado —
     *  cadenas CON comillas ("Hola Mundo"), arreglos elemento por
     *  elemento con el mismo criterio ([1, 2.5, 7] y no [1.0, 2.5, 7.0]). */
    public static String valorReporte(Object v) {
        if (v == null) return "null";
        if (v instanceof String s) return "\"" + s + "\"";
        if (v instanceof ArrayList<?> lista) {
            StringBuilder b = new StringBuilder("[");
            for (int i = 0; i < lista.size(); i++) {
                if (i > 0) b.append(", ");
                b.append(valorReporte(lista.get(i)));
            }
            return b.append("]").toString();
        }
        return formatear(v);   // Double con la regla del 15/15.7
    }

    private static String nombreTipo(Object v) {
        if (v instanceof Double) return "double";
        if (v instanceof String) return "char[]";
        if (v instanceof ArrayList) return "arreglo";
        return "desconocido";
    }

    /* ================= gráficas (5.10) ================= */

    /** Recibe la lista de atributos del bloque; aplica "la última gana"
     *  y solo registra la gráfica si apareció su EXEC. */
    public void registrarGrafica(String tipo, ArrayList<?> attrs, int l, int c) {
        LinkedHashMap<String, Object> mapa = new LinkedHashMap<>();
        boolean exec = false;
        for (Object o : attrs) {
            Object[] par = (Object[]) o;
            String nombre = (String) par[0];
            if (nombre.equals("EXEC")) { exec = true; continue; }
            String clave = nombre.toLowerCase();
            if (!ATTRS_VALIDOS.contains(clave)) {
                error("Semántico", "'" + nombre + "' no es un atributo de gráfica válido", l, c);
                continue;
            }
            mapa.put(clave, par[1]);   // put repetido = la última instrucción gana
        }
        if (!exec) return;             // sin EXEC no se muestra (5.10) — no es error
        if (!validarGrafica(tipo, mapa, l, c)) return;
        if (tipo.equals("Histogram")) tablaHistograma(mapa);
        graficas.add(new Grafica(tipo, mapa));
    }

    /** Chequeo semántico: cada tipo de gráfica exige SUS atributos,
     *  con el tipo correcto y (donde aplica) tamaños que coincidan. */
    private boolean validarGrafica(String tipo, LinkedHashMap<String, Object> m, int l, int c) {
        boolean ok = true;
        switch (tipo) {
            case "graphBar", "graphLine" -> {
                ok &= exigirTexto(tipo, m, "titulo", l, c);
                ok &= exigirTexto(tipo, m, "titulox", l, c);
                ok &= exigirTexto(tipo, m, "tituloy", l, c);
                ok &= exigirLista(tipo, m, "ejex", String.class, l, c);
                ok &= exigirLista(tipo, m, "ejey", Double.class, l, c);
                if (ok && ((ArrayList<?>) m.get("ejex")).size() != ((ArrayList<?>) m.get("ejey")).size()) {
                    error("Semántico", tipo + ": ejeX y ejeY deben tener la misma cantidad de elementos", l, c);
                    ok = false;
                }
            }
            case "graphPie" -> {
                ok &= exigirTexto(tipo, m, "titulo", l, c);
                ok &= exigirLista(tipo, m, "label", String.class, l, c);
                ok &= exigirLista(tipo, m, "values", Double.class, l, c);
                if (ok && ((ArrayList<?>) m.get("label")).size() != ((ArrayList<?>) m.get("values")).size()) {
                    error("Semántico", tipo + ": label y values deben tener la misma cantidad de elementos", l, c);
                    ok = false;
                }
            }
            case "Histogram" -> {
                ok &= exigirTexto(tipo, m, "titulo", l, c);
                ok &= exigirLista(tipo, m, "values", Double.class, l, c);
            }
        }
        return ok;
    }

    private boolean exigirTexto(String tipo, LinkedHashMap<String, Object> m, String attr, int l, int c) {
        Object v = m.get(attr);
        if (v instanceof String) return true;
        error("Semántico", tipo + ": falta el atributo '" + attr + "' (o no es una cadena)", l, c);
        return false;
    }

    private boolean exigirLista(String tipo, LinkedHashMap<String, Object> m, String attr,
                                Class<?> elemento, int l, int c) {
        Object v = m.get(attr);
        if (v instanceof ArrayList<?> lista && !lista.isEmpty()) {
            boolean todos = true;
            for (Object e : lista) todos &= elemento.isInstance(e);
            if (todos) return true;
        }
        String esperado = (elemento == Double.class) ? "double" : "char[]";
        error("Semántico", tipo + ": falta el atributo '" + attr
                + "' (o no es un arreglo de " + esperado + ")", l, c);
        return false;
    }

    /** Salida en consola del histograma (5.10.3): frecuencia,
     *  frecuencia acumulada y frecuencia relativa de cada valor. */
    private void tablaHistograma(LinkedHashMap<String, Object> m) {
        ArrayList<Double> datos = new ArrayList<>();
        for (Object v : (ArrayList<?>) m.get("values")) datos.add((Double) v);
        var frec = Operaciones.frecuencias(datos);
        int n = datos.size(), acum = 0;
        consola.append("--------------\n").append(m.get("titulo")).append('\n')
               .append("--------------\n")
               .append(String.format("%-8s %-6s %-6s %-6s%n", "Valor", "Frec.", "Acum.", "Rel."));
        for (var e : frec.entrySet()) {
            acum += e.getValue();
            consola.append(String.format(java.util.Locale.US, "%-8s %-6d %-6d %-6.2f%n",
                    formatear(e.getKey()), e.getValue(), acum, e.getValue() / (double) n));
        }
    }

    /* ================= getters para reportes/GUI ================= */

    public String getConsola() { return consola.toString(); }
    public List<RegistroError> getErrores() { return errores; }
    public LinkedHashMap<String, Simbolo> getSimbolos() { return simbolos; }
    public List<Grafica> getGraficas() { return graficas; }
}
