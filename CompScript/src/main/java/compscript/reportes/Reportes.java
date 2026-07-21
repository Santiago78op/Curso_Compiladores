package compscript.reportes;

import java.io.File;
import java.lang.reflect.Field;
import java.nio.file.Files;
import java.util.LinkedHashMap;

import compscript.analisis.sym;
import compscript.ast.A;
import compscript.interprete.Contexto;
import compscript.interprete.RegistroError;
import compscript.interprete.Simbolo;

/**
 * Genera los reportes del enunciado (6) a partir del Contexto de la ULTIMA
 * ejecucion: tokens (6.1), errores (6.2), AST como grafo (6.3) y tabla de
 * simbolos (6.4). Todo en HTML autocontenido (se abre sin internet).
 *
 * El AST se representa de DOS formas complementarias:
 *   - ast.html : grafo en forma de arbol (nodo raiz -> hijos), legible sin
 *                dependencias externas (decision: 100% autocontenido).
 *   - ast.dot  : fuente Graphviz para quien quiera un render de alta fidelidad
 *                (dot -Tpng ast.dot -o ast.png).
 * Ademas la GUI muestra el AST en un TreeView "desde la interfaz".
 */
public class Reportes {

    /** Escribe todos los reportes y devuelve los HTML a abrir. */
    public static File[] generar(Contexto ctx, File carpeta) throws Exception {
        carpeta.mkdirs();
        File t = new File(carpeta, "tokens.html");
        File e = new File(carpeta, "errores.html");
        File s = new File(carpeta, "simbolos.html");
        File a = new File(carpeta, "ast.html");
        File dot = new File(carpeta, "ast.dot");
        Files.writeString(t.toPath(), tokens(ctx));
        Files.writeString(e.toPath(), errores(ctx));
        Files.writeString(s.toPath(), simbolos(ctx));
        Files.writeString(a.toPath(), astHtml(ctx));
        Files.writeString(dot.toPath(), astDot(ctx));
        return new File[]{ t, e, s, a };
    }

    /* ---------- 6.1 tokens ---------- */
    private static String tokens(Contexto ctx) throws Exception {
        LinkedHashMap<Integer, String> nombres = nombresTokens();
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (Object[] tk : ctx.tokens) {
            filas.append(fila(n++, esc((String) tk[0]), nombres.get((Integer) tk[1]), tk[2], tk[3]));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='5'>— sin tokens —</td></tr>");
        return pagina("Reporte de Tokens",
                new String[]{"#", "Lexema", "Tipo", "Linea", "Columna"}, filas.toString());
    }

    private static LinkedHashMap<Integer, String> nombresTokens() throws Exception {
        LinkedHashMap<Integer, String> m = new LinkedHashMap<>();
        for (Field f : sym.class.getFields()) {
            if (f.getType() == int.class) m.put(f.getInt(null), f.getName());
        }
        return m;
    }

    /* ---------- 6.2 errores ---------- */
    private static String errores(Contexto ctx) {
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (RegistroError e : ctx.errores) {
            filas.append(fila(n++, e.tipo, esc(e.descripcion), e.linea, e.columna));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='5'>&#10003; sin errores</td></tr>");
        return pagina("Reporte de Errores",
                new String[]{"#", "Tipo", "Descripcion", "Linea", "Columna"}, filas.toString());
    }

    /* ---------- 6.4 tabla de simbolos ---------- */
    private static String simbolos(Contexto ctx) {
        StringBuilder filas = new StringBuilder();
        int n = 1;
        for (Simbolo s : ctx.simbolos) {
            filas.append(fila(n++, esc(s.nombre), s.categoria, esc(s.tipo.nombre()),
                    esc(s.ambito), esc(s.valor == null ? "null" : s.valor.reporte()),
                    s.linea, s.columna));
        }
        if (filas.length() == 0) filas.append("<tr><td colspan='8'>— sin simbolos —</td></tr>");
        return pagina("Tabla de Simbolos",
                new String[]{"#", "Id", "Tipo", "Tipo de Dato", "Entorno", "Valor", "Linea", "Columna"},
                filas.toString());
    }

    /* ---------- 6.3 AST (arbol HTML) ---------- */
    private static String astHtml(Contexto ctx) {
        StringBuilder arbol = new StringBuilder();
        if (ctx.raiz != null) astHtmlNodo(ctx.raiz, arbol);
        else arbol.append("<li>— sin AST —</li>");
        return """
            <!doctype html>
            <html lang="es"><head><meta charset="utf-8"><title>AST — CompScript</title>
            <style>
              body { font-family: Segoe UI, sans-serif; margin: 32px; background: #f5f7f6; color: #1c2321; }
              h1 { font-size: 22px; } h1 + p { color: #5c6b66; margin-top: -8px; }
              ul { list-style: none; padding-left: 22px; border-left: 1px dashed #9fb3ad; }
              li { margin: 3px 0; }
              .n { display: inline-block; background: #0e7c6e; color: #fff; padding: 2px 10px;
                   border-radius: 5px; font: 13px Consolas, monospace; }
              body > ul { border-left: none; padding-left: 0; }
            </style></head><body>
            <h1>Reporte de AST</h1>
            <p>CompScript — OLC1 PT1 · arbol de sintaxis abstracta del ultimo analisis</p>
            <ul>%s</ul>
            </body></html>
            """.formatted(arbol.toString());
    }

    private static void astHtmlNodo(A.Nodo nodo, StringBuilder sb) {
        sb.append("<li><span class='n'>").append(esc(nodo.etiquetaAst())).append("</span>");
        java.util.List<A.Nodo> hijos = nodo.hijosAst();
        if (!hijos.isEmpty()) {
            sb.append("<ul>");
            for (A.Nodo h : hijos) if (h != null) astHtmlNodo(h, sb);
            sb.append("</ul>");
        }
        sb.append("</li>");
    }

    /* ---------- 6.3 AST (Graphviz .dot, alta fidelidad opcional) ---------- */
    private static String astDot(Contexto ctx) {
        StringBuilder sb = new StringBuilder();
        sb.append("digraph AST {\n  node [shape=box, style=rounded, fontname=Consolas];\n");
        if (ctx.raiz != null) astDotNodo(ctx.raiz, sb, new int[]{0});
        sb.append("}\n");
        return sb.toString();
    }

    private static int astDotNodo(A.Nodo nodo, StringBuilder sb, int[] contador) {
        int id = contador[0]++;
        sb.append("  n").append(id).append(" [label=\"")
          .append(nodo.etiquetaAst().replace("\\", "\\\\").replace("\"", "\\\"")).append("\"];\n");
        for (A.Nodo h : nodo.hijosAst()) {
            if (h == null) continue;
            int hijoId = astDotNodo(h, sb, contador);
            sb.append("  n").append(id).append(" -> n").append(hijoId).append(";\n");
        }
        return id;
    }

    /* ---------- plantilla y utilitarios ---------- */
    private static String pagina(String titulo, String[] columnas, String filas) {
        StringBuilder ths = new StringBuilder();
        for (String c : columnas) ths.append("<th>").append(c).append("</th>");
        return """
            <!doctype html>
            <html lang="es"><head><meta charset="utf-8"><title>%s — CompScript</title>
            <style>
              body { font-family: Segoe UI, sans-serif; margin: 32px; background: #f5f7f6; color: #1c2321; }
              h1 { font-size: 22px; } h1 + p { color: #5c6b66; margin-top: -8px; }
              table { border-collapse: collapse; background: white; box-shadow: 0 1px 4px rgba(0,0,0,.12); }
              th, td { padding: 8px 16px; text-align: left; border-bottom: 1px solid #d8e0dd; }
              th { background: #0e7c6e; color: white; font-size: 13px; }
              tr:nth-child(even) td { background: #eef3f1; }
              td { font-family: Consolas, monospace; font-size: 14px; }
            </style></head><body>
            <h1>%s</h1><p>CompScript — OLC1 PT1 · generado del ultimo analisis</p>
            <table><thead><tr>%s</tr></thead><tbody>%s</tbody></table>
            </body></html>
            """.formatted(titulo, titulo, ths, filas);
    }

    private static String fila(Object... celdas) {
        StringBuilder f = new StringBuilder("<tr>");
        for (Object c : celdas) f.append("<td>").append(c).append("</td>");
        return f.append("</tr>\n").toString();
    }

    private static String esc(String s) {
        if (s == null) return "";
        return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;");
    }
}
