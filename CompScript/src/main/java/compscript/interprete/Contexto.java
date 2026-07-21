package compscript.interprete;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

import compscript.ast.A;

/**
 * Estado GLOBAL de una ejecucion (fresco por corrida, enunciado 6):
 *  - consola de salida
 *  - errores acumulados (lexico + sintactico + semantico) -> reporte 6.2
 *  - tokens reconocidos por el lexer -> reporte 6.1
 *  - registro de funciones/metodos y structs (declarados en la 1a pasada)
 *  - registro plano de TODOS los simbolos declarados -> reporte 6.4
 *
 * A diferencia del Entorno (que es UN ambito de la pila de scopes), el
 * Contexto es unico y lo comparten todos los entornos de una ejecucion.
 */
public class Contexto {

    public final StringBuilder consola = new StringBuilder();
    public final List<RegistroError> errores = new ArrayList<>();
    public final List<Object[]> tokens = new ArrayList<>();  // {lexema, tipoInt, linea, col}
    public final List<Simbolo> simbolos = new ArrayList<>();  // todos, con su ambito

    public final LinkedHashMap<String, A.DeclaracionFuncion> funciones = new LinkedHashMap<>();
    public final LinkedHashMap<String, A.DeclaracionStruct> structs = new LinkedHashMap<>();

    public A.Nodo raiz;   // raiz del AST (para el reporte 6.3)
    public Entorno global; // entorno raiz: padre de todo ambito de funcion (alcance estatico)

    /* ================= errores ================= */

    public void error(String tipo, String descripcion, int linea, int columna) {
        errores.add(new RegistroError(tipo, descripcion, linea + 1, columna + 1));
    }

    public boolean hayErrores() { return !errores.isEmpty(); }

    public boolean hayErroresGraves() {   // lexico o sintactico
        for (RegistroError e : errores)
            if (e.tipo.startsWith("Lex") || e.tipo.startsWith("Sint")) return true;
        return false;
    }

    /* ================= tokens (reporte 6.1) ================= */

    public void registrarToken(String lexema, int tipo, int linea, int columna) {
        tokens.add(new Object[]{ lexema, tipo, linea + 1, columna + 1 });
    }

    /* ================= registro de funciones y structs (1a pasada) ================= */

    public void registrarFuncion(A.DeclaracionFuncion f) {
        String clave = f.id.toLowerCase();
        if (funciones.containsKey(clave) || structs.containsKey(clave)) {
            error("Semantico", "ya existe una funcion, metodo o struct con el id '"
                    + f.id + "' (no hay sobrecarga)", f.linea, f.columna);
            return;
        }
        funciones.put(clave, f);
    }

    public void registrarStruct(A.DeclaracionStruct s) {
        String clave = s.id.toLowerCase();
        if (structs.containsKey(clave) || funciones.containsKey(clave)) {
            error("Semantico", "ya existe un struct, funcion o metodo con el id '"
                    + s.id + "'", s.linea, s.columna);
            return;
        }
        structs.put(clave, s);
    }
}
