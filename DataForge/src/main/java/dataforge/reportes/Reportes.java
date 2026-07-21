package dataforge.reportes;

import java.io.File;
import java.lang.reflect.Field;
import java.nio.file.Files;
import java.util.LinkedHashMap;

import dataforge.analisis.sym;
import dataforge.interprete.Entorno;
import dataforge.interprete.RegistroError;
import dataforge.interprete.Simbolo;

/**
 * Genera los tres reportes HTML del enunciado (§6) a partir del
 * Entorno de la ÚLTIMA ejecución: tokens, errores y tabla de símbolos.
 * El diseño queda "a discreción del estudiante, solo legible" — tabla
 * simple con CSS embebido (autocontenido: se abre sin internet).
 */
public class Reportes {

    /** Escribe tokens.html, errores.html y simbolos.html en la carpeta. */
    public static File[] generar(Entorno ent, File carpeta) throws Exception {
        carpeta.mkdirs();
        File t = new File(carpeta, "tokens.html");
        File e = new File(carpeta, "errores.html");
        File s = new File(carpeta, "simbolos.html");
        Files.writeString(t.toPath(), tokens(ent));
        Files.writeString(e.toPath(), errores(ent));
        Files.writeString(s.toPath(), simbolos(ent));
        return new File[]{ t, e, s };
    }

    /* ---------- §6.1: tokens (lexema ORIGINAL + nombre del token) ---------- */

    private static String tokens(Entorno ent) throws Exception {
        LinkedHashMap<Integer, String> nombres = nombresTokens();
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (Object[] tk : ent.getTokens()) {
            filas.append(fila(n++, esc((String) tk[0]), nombres.get((Integer) tk[1]), tk[2], tk[3]));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='5'>— sin tokens —</td></tr>");
        return pagina("Reporte de Tokens",
                new String[]{"#", "Lexema", "Tipo", "Línea", "Columna"}, filas.toString());
    }

    /** id numérico del token → su nombre en sym (por reflexión, sin switch). */
    private static LinkedHashMap<Integer, String> nombresTokens() throws Exception {
        LinkedHashMap<Integer, String> m = new LinkedHashMap<>();
        for (Field f : sym.class.getFields()) {
            if (f.getType() == int.class) m.put(f.getInt(null), f.getName());
        }
        return m;
    }

    /* ---------- §6.2: errores (léxicos + sintácticos + semánticos) ---------- */

    private static String errores(Entorno ent) {
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (RegistroError e : ent.getErrores()) {
            filas.append(fila(n++, e.tipo, esc(e.descripcion), e.linea, e.columna));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='5'>✓ sin errores</td></tr>");
        return pagina("Reporte de Errores",
                new String[]{"#", "Tipo", "Descripción", "Línea", "Columna"}, filas.toString());
    }

    /* ---------- §6.3: tabla de símbolos ---------- */

    private static String simbolos(Entorno ent) {
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (Simbolo s : ent.getSimbolos().values()) {
            filas.append(fila(n++, esc(s.nombre), s.categoria, s.tipo,
                    esc(Entorno.valorReporte(s.valor)), s.linea, s.columna));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='7'>— sin símbolos declarados —</td></tr>");
        return pagina("Tabla de Símbolos",
                new String[]{"#", "Nombre", "Categoría", "Tipo", "Valor", "Línea", "Columna"},
                filas.toString());
    }

    /* ---------- plantilla y utilitarios ---------- */

    private static String pagina(String titulo, String[] columnas, String filas) {
        StringBuilder ths = new StringBuilder();
        for (String c : columnas) ths.append("<th>").append(c).append("</th>");
        return """
            <!doctype html>
            <html lang="es"><head><meta charset="utf-8"><title>%s — DataForge</title>
            <style>
              body { font-family: Segoe UI, sans-serif; margin: 32px; background: #f5f7f6; color: #1c2321; }
              h1 { font-size: 22px; } h1 + p { color: #5c6b66; margin-top: -8px; }
              table { border-collapse: collapse; background: white; box-shadow: 0 1px 4px rgba(0,0,0,.12); }
              th, td { padding: 8px 16px; text-align: left; border-bottom: 1px solid #d8e0dd; }
              th { background: #0e7c6e; color: white; font-size: 13px; }
              tr:nth-child(even) td { background: #eef3f1; }
              td { font-family: Consolas, monospace; font-size: 14px; }
            </style></head><body>
            <h1>%s</h1><p>DataForge — OLC1 Proyecto 1 · generado del último análisis</p>
            <table><thead><tr>%s</tr></thead><tbody>%s</tbody></table>
            </body></html>
            """.formatted(titulo, titulo, ths, filas);
    }

    private static String fila(Object... celdas) {
        StringBuilder f = new StringBuilder("<tr>");
        for (Object c : celdas) f.append("<td>").append(c).append("</td>");
        return f.append("</tr>\n").toString();
    }

    /** Escapa HTML: un lexema como "<-" no debe romper la tabla. */
    private static String esc(String s) {
        return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;");
    }
}
