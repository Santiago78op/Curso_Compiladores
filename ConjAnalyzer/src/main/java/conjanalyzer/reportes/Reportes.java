package conjanalyzer.reportes;

import java.io.File;
import java.lang.reflect.Field;
import java.nio.file.Files;
import java.util.LinkedHashMap;

import conjanalyzer.analisis.sym;
import conjanalyzer.interprete.Conjunto;
import conjanalyzer.interprete.Entorno;
import conjanalyzer.interprete.Operacion;
import conjanalyzer.interprete.RegistroError;

/**
 * Genera los reportes HTML del enunciado (5.2 y 5.3) a partir del Entorno de
 * la ULTIMA ejecucion: tabla de tokens y tabla de errores. Se agrega ademas
 * un reporte de conjuntos y operaciones para consulta. El diseno queda "a
 * discrecion del estudiante, solo legible": tabla simple con CSS embebido
 * (autocontenido, se abre sin internet).
 */
public class Reportes {

    /** Escribe tokens.html, errores.html y operaciones.html en la carpeta. */
    public static File[] generar(Entorno ent, File carpeta) throws Exception {
        carpeta.mkdirs();
        File t = new File(carpeta, "tokens.html");
        File e = new File(carpeta, "errores.html");
        File o = new File(carpeta, "operaciones.html");
        Files.writeString(t.toPath(), tokens(ent));
        Files.writeString(e.toPath(), errores(ent));
        Files.writeString(o.toPath(), operaciones(ent));
        return new File[]{ t, e, o };
    }

    /* ---------- 5.2: tokens (lexema + nombre del token + posicion) ---------- */

    private static String tokens(Entorno ent) throws Exception {
        LinkedHashMap<Integer, String> nombres = nombresTokens();
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (Object[] tk : ent.getTokens()) {
            filas.append(fila(n++, esc((String) tk[0]), nombres.get((Integer) tk[1]), tk[2], tk[3]));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='5'>— sin tokens —</td></tr>");
        return pagina("Reporte de Tokens",
                new String[]{"#", "Lexema", "Tipo", "Linea", "Columna"}, filas.toString());
    }

    /** id numerico del token → su nombre en sym (por reflexion, sin switch). */
    private static LinkedHashMap<Integer, String> nombresTokens() throws Exception {
        LinkedHashMap<Integer, String> m = new LinkedHashMap<>();
        for (Field f : sym.class.getFields()) {
            if (f.getType() == int.class) m.put(f.getInt(null), f.getName());
        }
        return m;
    }

    /* ---------- 5.3: errores (lexicos + sintacticos + semanticos) ---------- */

    private static String errores(Entorno ent) {
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (RegistroError e : ent.getErrores()) {
            filas.append(fila(n++, e.tipo, esc(e.descripcion), e.linea, e.columna));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='5'>sin errores</td></tr>");
        return pagina("Reporte de Errores",
                new String[]{"#", "Tipo", "Descripcion", "Linea", "Columna"}, filas.toString());
    }

    /* ---------- extra: conjuntos y operaciones ---------- */

    private static String operaciones(Entorno ent) {
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (Conjunto conj : ent.getConjuntos().values()) {
            filas.append(fila(n++, esc(conj.nombre), "CONJ", esc(conj.definicion),
                    esc(Entorno.formatearConjunto(conj.elementos))));
        }
        for (Operacion op : ent.getOperaciones().values()) {
            String simpl = op.simplificacion.seSimplifico
                    ? esc(op.simplificacion.simplificado.toPrefijo())
                    : "(no simplificable)";
            filas.append(fila(n++, esc(op.nombre), "OPERA", esc(op.arbol.toPrefijo()), simpl));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='5'>— sin definiciones —</td></tr>");
        return pagina("Conjuntos y Operaciones",
                new String[]{"#", "Nombre", "Categoria", "Definicion", "Resultado / Simplificacion"},
                filas.toString());
    }

    /* ---------- plantilla y utilitarios ---------- */

    private static String pagina(String titulo, String[] columnas, String filas) {
        StringBuilder ths = new StringBuilder();
        for (String c : columnas) ths.append("<th>").append(c).append("</th>");
        return """
            <!doctype html>
            <html lang="es"><head><meta charset="utf-8"><title>%s — ConjAnalyzer</title>
            <style>
              body { font-family: Segoe UI, sans-serif; margin: 32px; background: #f4f6fb; color: #1b2130; }
              h1 { font-size: 22px; } h1 + p { color: #5c6478; margin-top: -8px; }
              table { border-collapse: collapse; background: white; box-shadow: 0 1px 4px rgba(0,0,0,.12); }
              th, td { padding: 8px 16px; text-align: left; border-bottom: 1px solid #d9deea; }
              th { background: #3b4b8a; color: white; font-size: 13px; }
              tr:nth-child(even) td { background: #eef1f8; }
              td { font-family: Consolas, monospace; font-size: 14px; }
            </style></head><body>
            <h1>%s</h1><p>ConjAnalyzer — OLC1 Proyecto 1 · generado del ultimo analisis</p>
            <table><thead><tr>%s</tr></thead><tbody>%s</tbody></table>
            </body></html>
            """.formatted(titulo, titulo, ths, filas);
    }

    private static String fila(Object... celdas) {
        StringBuilder f = new StringBuilder("<tr>");
        for (Object c : celdas) f.append("<td>").append(c).append("</td>");
        return f.append("</tr>\n").toString();
    }

    /** Escapa HTML: un lexema como "&" o "<" no debe romper la tabla. */
    private static String esc(String s) {
        if (s == null) return "";
        return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;");
    }
}
