package compscript.ast;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import compscript.interprete.Contexto;
import compscript.interprete.Entorno;
import compscript.interprete.Operaciones;
import compscript.interprete.Senales;
import compscript.interprete.Simbolo;
import compscript.interprete.Tipo;
import compscript.interprete.Valor;

/**
 * AST de CompScript (todas las clases de nodo agrupadas). Se construye en
 * las acciones del parser (parser.cup) y se recorre para interpretar
 * (patron tree-walk: cada nodo sabe ejecutarse/evaluarse) y para el
 * reporte de AST como grafo (etiquetaAst + hijosAst).
 *
 * Los constructores reciben Object / ArrayList y castean internamente,
 * de modo que las acciones del .cup queden concisas (estilo DataForge).
 */
public final class A {

    private A() {}

    /* ============================================================
     * INTERFACES
     * ============================================================ */
    public interface Nodo {
        String etiquetaAst();
        default List<Nodo> hijosAst() { return new ArrayList<>(); }
    }
    public interface Instruccion extends Nodo { void ejecutar(Entorno e); }
    public interface Expresion  extends Nodo { Valor evaluar(Entorno e); }

    /* ============================================================
     * HELPERS
     * ============================================================ */
    @SuppressWarnings("unchecked")
    static List<Instruccion> insts(Object o) {
        return o == null ? new ArrayList<>() : (List<Instruccion>) o;
    }

    /** Ejecuta una lista de instrucciones en el entorno dado. */
    public static void ejecutar(List<Instruccion> lista, Entorno e) {
        if (lista == null) return;
        for (Instruccion i : lista) if (i != null) i.ejecutar(e);
    }

    static void add(List<Nodo> acc, Object o) {
        if (o == null) return;
        if (o instanceof Nodo n) acc.add(n);
        else if (o instanceof List<?> l) for (Object x : l) if (x instanceof Nodo n) acc.add(n);
    }
    static List<Nodo> hijos(Object... xs) {
        List<Nodo> acc = new ArrayList<>();
        for (Object o : xs) add(acc, o);
        return acc;
    }

    private static int indice(Valor v, Entorno e, int l, int c) {
        if (v == null || v.tipo.cat != Tipo.Cat.INT)
            e.errorSemantico("el indice debe ser un Entero", l, c);
        return (Integer) v.valor;
    }

    /** Valor por defecto de un tipo (5.3): defaults de primitivos + compuestos. */
    public static Valor defecto(Tipo t, Entorno e, int l, int c) {
        switch (t.cat) {
            case INT:    return Valor.vInt(0);
            case DOUBLE: return Valor.vDouble(0.0);
            case BOOL:   return Valor.vBool(true);      // default true (5.3)
            case CHAR:   return Valor.vChar(' ');
            case STRING: return Valor.vString("");
            case VECTOR: return new Valor(t, new ArrayList<Valor>());
            case LIST:   return new Valor(t, new ArrayList<Valor>());
            case STRUCT: return instanciarStructDefecto(t, e, l, c);
            default:     return Valor.VOID;
        }
    }

    private static Valor instanciarStructDefecto(Tipo t, Entorno e, int l, int c) {
        DeclaracionStruct def = e.contexto.structs.get(t.structName);
        if (def == null) e.errorSemantico("el struct '" + t.structName + "' no ha sido definido", l, c);
        LinkedHashMap<String, Valor> campos = new LinkedHashMap<>();
        for (CampoStruct cs : def.campos) campos.put(cs.id, defecto(cs.tipo, e, l, c));
        return new Valor(Tipo.struct(def.id), campos);
    }

    /* ============================================================
     * DATOS AUXILIARES
     * ============================================================ */
    public static class Parametro implements Nodo {
        public final String id; public final Tipo tipo; public final Expresion defecto;
        public Parametro(String id, Tipo tipo, Object defecto) {
            this.id = id; this.tipo = tipo; this.defecto = (Expresion) defecto;
        }
        public String etiquetaAst() { return "param " + id + ":" + tipo.nombre(); }
        public List<Nodo> hijosAst() { return hijos(defecto); }
    }
    public static class Argumento implements Nodo {
        public final String id; public final Expresion expr; public final int linea, columna;
        public Argumento(String id, Object expr, int l, int c) {
            this.id = id; this.expr = (Expresion) expr; linea = l; columna = c;
        }
        public String etiquetaAst() { return "arg " + id; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }
    public static class CampoStruct implements Nodo {
        public final String id; public final Tipo tipo;
        public CampoStruct(String id, Tipo tipo) { this.id = id; this.tipo = tipo; }
        public String etiquetaAst() { return "campo " + id + ":" + tipo.nombre(); }
    }
    public static class CampoValor implements Nodo {
        public final String campo; public final Expresion expr; public final int linea, columna;
        public CampoValor(String campo, Object expr, int l, int c) {
            this.campo = campo; this.expr = (Expresion) expr; linea = l; columna = c;
        }
        public String etiquetaAst() { return campo + ":"; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }
    public static class CasoMatch implements Nodo {
        public final Expresion valor; public final List<Instruccion> cuerpo;
        public CasoMatch(Object valor, Object cuerpo) {
            this.valor = (Expresion) valor; this.cuerpo = insts(cuerpo);
        }
        public String etiquetaAst() { return "case =>"; }
        public List<Nodo> hijosAst() { return hijos(valor, cuerpo); }
    }

    /* ============================================================
     * EXPRESIONES
     * ============================================================ */

    public static class Literal implements Expresion {
        public final Valor valor;
        public Literal(Valor valor) { this.valor = valor; }
        public Valor evaluar(Entorno e) { return valor; }
        public String etiquetaAst() { return valor.tipo.nombre() + ": " + valor.reporte(); }
    }

    public static class AccesoVariable implements Expresion {
        public final String id; public final int linea, columna;
        public AccesoVariable(String id, int l, int c) { this.id = id; linea = l; columna = c; }
        public Valor evaluar(Entorno e) { return e.obtener(id, linea, columna).valor; }
        public String etiquetaAst() { return "id: " + id; }
    }

    public static class Binaria implements Expresion {
        public final String op; public final Expresion izq, der; public final int linea, columna;
        public Binaria(String op, Object izq, Object der, int l, int c) {
            this.op = op; this.izq = (Expresion) izq; this.der = (Expresion) der; linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            switch (op) {
                case "+": case "-": case "*": case "/":
                case "^": case "$": case "%":
                    return Operaciones.aritmetica(op, izq.evaluar(e), der.evaluar(e), e, linea, columna);
                case "==": case "!=": case "<": case "<=": case ">": case ">=":
                    return Operaciones.relacional(op, izq.evaluar(e), der.evaluar(e), e, linea, columna);
                case "&&": {
                    // Cortocircuito (bug real corregido: antes se evaluaba SIEMPRE
                    // "der", incluso con izq ya false). Si izq es un Booleano falso,
                    // el resultado es false sin necesidad de evaluar der -- eso
                    // importa cuando der tiene efectos secundarios o falla en tiempo
                    // de ejecucion (p. ej. "x != 0 && 10/x > 1").
                    Valor a = izq.evaluar(e);
                    if (a.tipo.cat == Tipo.Cat.BOOL && !(Boolean) a.valor) return Valor.vBool(false);
                    return Operaciones.logica(op, a, der.evaluar(e), e, linea, columna);
                }
                case "||": {
                    // Cortocircuito simetrico: si izq es un Booleano verdadero, el
                    // resultado es true sin evaluar der.
                    Valor a = izq.evaluar(e);
                    if (a.tipo.cat == Tipo.Cat.BOOL && (Boolean) a.valor) return Valor.vBool(true);
                    return Operaciones.logica(op, a, der.evaluar(e), e, linea, columna);
                }
                default: return null;
            }
        }
        public String etiquetaAst() { return "op " + op; }
        public List<Nodo> hijosAst() { return hijos(izq, der); }
    }

    public static class Unaria implements Expresion {
        public final String op; public final Expresion expr; public final int linea, columna;
        public Unaria(String op, Object expr, int l, int c) {
            this.op = op; this.expr = (Expresion) expr; linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            Valor v = expr.evaluar(e);
            return op.equals("NEG")
                    ? Operaciones.negacion(v, e, linea, columna)
                    : Operaciones.negacionLogica(v, e, linea, columna);
        }
        public String etiquetaAst() { return op.equals("NEG") ? "neg (-)" : "not (!)"; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }

    public static class Cast implements Expresion {
        public final Expresion expr; public final Tipo destino; public final int linea, columna;
        public Cast(Object expr, Tipo destino, int l, int c) {
            this.expr = (Expresion) expr; this.destino = destino; linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            return Operaciones.cast(expr.evaluar(e), destino, e, linea, columna);
        }
        public String etiquetaAst() { return "cast as " + destino.nombre(); }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }

    public static class AccesoVector implements Expresion {
        public final String id; public final Expresion idx1, idx2; public final int linea, columna;
        public AccesoVector(String id, Object idx1, Object idx2, int l, int c) {
            this.id = id; this.idx1 = (Expresion) idx1; this.idx2 = (Expresion) idx2; linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            Simbolo s = e.obtener(id, linea, columna);
            if (s.tipo.cat != Tipo.Cat.VECTOR)
                e.errorSemantico("'" + id + "' no es un vector", linea, columna);
            List<Valor> data = s.valor.lista();
            int i = indice(idx1.evaluar(e), e, linea, columna);
            if (i < 0 || i >= data.size())
                e.errorSemantico("indice " + i + " fuera de rango en '" + id + "'", linea, columna);
            Valor fila = data.get(i);
            if (idx2 == null) return fila;
            if (fila.tipo.cat != Tipo.Cat.VECTOR)
                e.errorSemantico("'" + id + "' no es un vector de dos dimensiones", linea, columna);
            List<Valor> interno = fila.lista();
            int j = indice(idx2.evaluar(e), e, linea, columna);
            if (j < 0 || j >= interno.size())
                e.errorSemantico("indice " + j + " fuera de rango en '" + id + "'", linea, columna);
            return interno.get(j);
        }
        public String etiquetaAst() { return "acceso vector " + id + (idx2 == null ? "[]" : "[][]"); }
        public List<Nodo> hijosAst() { return hijos(idx1, idx2); }
    }

    public static class LiteralVector implements Expresion {
        public final List<Expresion> elementos;
        @SuppressWarnings("unchecked")
        public LiteralVector(Object elementos) { this.elementos = (List<Expresion>) elementos; }
        public Valor evaluar(Entorno e) {
            List<Valor> vals = new ArrayList<>();
            for (Expresion ex : elementos) vals.add(ex.evaluar(e));
            Tipo base; int dim;
            if (!vals.isEmpty() && vals.get(0).tipo.cat == Tipo.Cat.VECTOR) {
                dim = 2; base = vals.get(0).tipo.elemento;
            } else {
                dim = 1; base = vals.isEmpty() ? Tipo.NULL : vals.get(0).tipo;
            }
            return new Valor(Tipo.vector(base, dim), vals);
        }
        public String etiquetaAst() { return "vector literal [ ]"; }
        public List<Nodo> hijosAst() { return hijos(elementos); }
    }

    public static class LiteralStruct implements Expresion {
        public final List<CampoValor> campos;
        @SuppressWarnings("unchecked")
        public LiteralStruct(Object campos) { this.campos = (List<CampoValor>) campos; }
        // Nunca se evalua directamente: la Declaracion la interpreta con el tipo destino.
        public Valor evaluar(Entorno e) {
            e.errorSemantico("un literal de struct solo puede usarse al instanciar", 0, 0);
            return null;
        }
        public String etiquetaAst() { return "struct { }"; }
        public List<Nodo> hijosAst() { return hijos(campos); }
    }

    public static class AccesoCampo implements Expresion {
        public final String id, campo; public final int linea, columna;
        public AccesoCampo(String id, String campo, int l, int c) {
            this.id = id; this.campo = campo; linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            Simbolo s = e.obtener(id, linea, columna);
            if (s.tipo.cat != Tipo.Cat.STRUCT)
                e.errorSemantico("'" + id + "' no es un struct", linea, columna);
            Valor v = buscarCampo(s.valor.campos(), campo);
            if (v == null)
                e.errorSemantico("el struct '" + id + "' no tiene el campo '" + campo + "'", linea, columna);
            return v;
        }
        public String etiquetaAst() { return "acceso " + id + "." + campo; }
    }

    public static class Nativa implements Expresion {
        public final String func; public final Expresion arg; public final int linea, columna;
        public Nativa(String func, Object arg, int l, int c) {
            this.func = func; this.arg = (Expresion) arg; linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            Valor v = arg.evaluar(e);
            switch (func) {
                case "round": {
                    if (!v.tipo.esNumerico())
                        e.errorSemantico("round espera un valor numerico", linea, columna);
                    return Valor.vInt((int) Math.round(v.numero()));
                }
                case "length": {
                    if (v.tipo.cat == Tipo.Cat.VECTOR || v.tipo.cat == Tipo.Cat.LIST)
                        return Valor.vInt(v.lista().size());
                    if (v.tipo.cat == Tipo.Cat.STRING)
                        return Valor.vInt(((String) v.valor).length());
                    e.errorSemantico("length espera un vector, lista o cadena", linea, columna);
                    return null;
                }
                case "tostring": {
                    switch (v.tipo.cat) {
                        case INT: case DOUBLE: case BOOL: case CHAR: return Valor.vString(v.texto());
                        case STRUCT: return Valor.vString(v.textoStruct());
                        default:
                            e.errorSemantico("toString espera un valor numerico, char, bool o struct", linea, columna);
                            return null;
                    }
                }
                default: return null;
            }
        }
        public String etiquetaAst() { return func + "()"; }
        public List<Nodo> hijosAst() { return hijos(arg); }
    }

    /** Llamada a funcion o metodo (5.23): sirve como expresion (funcion) y,
     *  envuelta en ExpresionStmt, como instruccion (metodo). */
    public static class Llamada implements Expresion {
        public final String id; public final List<Argumento> args; public final int linea, columna;
        @SuppressWarnings("unchecked")
        public Llamada(String id, Object args, int l, int c) {
            this.id = id; this.args = (List<Argumento>) args; linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            return invocar(id, args, e, linea, columna);
        }
        public String etiquetaAst() { return "llamada " + id + "()"; }
        public List<Nodo> hijosAst() { return hijos(args); }
    }

    /** Operacion sobre lista dinamica (5.19). get/remove/pop devuelven valor;
     *  push/set/reverse son void. Sirve como expresion o como instruccion. */
    public static class OperacionLista implements Expresion {
        public final String id, metodo; public final List<Expresion> args; public final int linea, columna;
        @SuppressWarnings("unchecked")
        public OperacionLista(String id, String metodo, Object args, int l, int c) {
            this.id = id; this.metodo = metodo.toLowerCase();
            this.args = args == null ? new ArrayList<>() : (List<Expresion>) args;
            linea = l; columna = c;
        }
        public Valor evaluar(Entorno e) {
            Simbolo s = e.obtener(id, linea, columna);
            if (s.tipo.cat != Tipo.Cat.LIST)
                e.errorSemantico("'" + id + "' no es una lista dinamica", linea, columna);
            List<Valor> data = s.valor.lista();
            Tipo base = s.tipo.elemento;
            switch (metodo) {
                case "push": {
                    Valor v = args.get(0).evaluar(e);
                    if (!Entorno.compatible(base, v))
                        e.errorSemantico("no se puede agregar " + Entorno.tipoDe(v)
                                + " a una lista de " + base.nombre(), linea, columna);
                    data.add(v); return Valor.VOID;
                }
                case "get": {
                    int i = indice(args.get(0).evaluar(e), e, linea, columna);
                    checarRango(i, data.size(), e);
                    return data.get(i);
                }
                case "set": {
                    int i = indice(args.get(0).evaluar(e), e, linea, columna);
                    checarRango(i, data.size(), e);
                    Valor v = args.get(1).evaluar(e);
                    if (!Entorno.compatible(base, v))
                        e.errorSemantico("no se puede asignar " + Entorno.tipoDe(v)
                                + " en una lista de " + base.nombre(), linea, columna);
                    data.set(i, v); return Valor.VOID;
                }
                case "remove": {
                    int i = indice(args.get(0).evaluar(e), e, linea, columna);
                    checarRango(i, data.size(), e);
                    return data.remove(i);
                }
                case "pop": {
                    if (data.isEmpty()) e.errorSemantico("pop() sobre una lista vacia", linea, columna);
                    return data.remove(data.size() - 1);
                }
                case "reverse": {
                    Collections.reverse(data); return Valor.VOID;
                }
                default:
                    e.errorSemantico("metodo de lista desconocido: " + metodo, linea, columna);
                    return null;
            }
        }
        private void checarRango(int i, int n, Entorno e) {
            if (i < 0 || i >= n) e.errorSemantico("indice " + i + " fuera de rango en la lista '" + id + "'", linea, columna);
        }
        public String etiquetaAst() { return "lista " + id + "." + metodo + "()"; }
        public List<Nodo> hijosAst() { return hijos(args); }
    }

    /* ============================================================
     * INSTRUCCIONES
     * ============================================================ */

    public static class Declaracion implements Instruccion {
        public final boolean esConst; public final String id; public final Tipo tipo;
        public final Expresion init; public final int linea, columna;
        public Declaracion(boolean esConst, String id, Tipo tipo, Object init, int l, int c) {
            this.esConst = esConst; this.id = id; this.tipo = tipo;
            this.init = (Expresion) init; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Valor val;
            switch (tipo.cat) {
                case STRUCT:
                    val = instanciarStruct(e); break;
                case VECTOR:
                    val = (init == null) ? defecto(tipo, e, linea, columna)
                            : validarVector(init.evaluar(e), tipo, e); break;
                case LIST:
                    val = new Valor(tipo, new ArrayList<Valor>()); break;
                default: {
                    if (init == null) { val = defecto(tipo, e, linea, columna); break; }
                    Valor v = init.evaluar(e);
                    if (!Entorno.compatible(tipo, v))
                        e.errorSemantico("no se puede asignar " + Entorno.tipoDe(v)
                                + " a '" + id + "' de tipo " + tipo.nombre(), linea, columna);
                    val = v;
                }
            }
            e.declarar(new Simbolo(id, Simbolo.categoriaDe(tipo), tipo, !esConst,
                    e.nombre, val, linea + 1, columna + 1), linea, columna);
        }

        private Valor instanciarStruct(Entorno e) {
            DeclaracionStruct def = e.contexto.structs.get(tipo.structName);
            if (def == null) e.errorSemantico("el struct '" + tipo.structName + "' no ha sido definido", linea, columna);
            LinkedHashMap<String, Valor> campos = new LinkedHashMap<>();
            for (CampoStruct cs : def.campos) campos.put(cs.id, defecto(cs.tipo, e, linea, columna));
            if (init instanceof LiteralStruct lit) {
                for (CampoValor cv : lit.campos) {
                    CampoStruct cs = def.buscar(cv.campo);
                    if (cs == null)
                        e.errorSemantico("el struct '" + def.id + "' no tiene el campo '" + cv.campo + "'",
                                cv.linea, cv.columna);
                    Valor v = cv.expr.evaluar(e);
                    if (!Entorno.compatible(cs.tipo, v))
                        e.errorSemantico("el campo '" + cs.id + "' es " + cs.tipo.nombre()
                                + " y recibio " + Entorno.tipoDe(v), cv.linea, cv.columna);
                    campos.put(cs.id, v);
                }
            }
            return new Valor(Tipo.struct(def.id), campos);
        }

        private Valor validarVector(Valor crudo, Tipo declarado, Entorno e) {
            // Guarda de tipo: init es un <expresion> general (no solo un literal de
            // vector), asi que puede evaluar a un Entero/Cadena/etc. Sin esta
            // comprobacion, crudo.lista() lanza ClassCastException sin control
            // (bug real: revienta el interprete completo en vez de reportar un
            // error semantico y terminar de forma ordenada, como exige 4.3).
            if (crudo == null || crudo.tipo.cat != Tipo.Cat.VECTOR)
                e.errorSemantico("el vector '" + id + "' espera un literal de vector y recibio "
                        + Entorno.tipoDe(crudo), linea, columna);
            List<Valor> vals = crudo.lista();
            List<Valor> out = new ArrayList<>();
            if (declarado.dimensiones == 1) {
                for (Valor v : vals) {
                    if (!Entorno.compatible(declarado.elemento, v))
                        e.errorSemantico("el vector '" + id + "' es de " + declarado.elemento.nombre()
                                + " y contiene " + Entorno.tipoDe(v), linea, columna);
                    out.add(v);
                }
            } else { // 2D
                for (Valor fila : vals) {
                    if (fila.tipo.cat != Tipo.Cat.VECTOR)
                        e.errorSemantico("el vector '" + id + "' de dos dimensiones esta mal formado", linea, columna);
                    List<Valor> filaOut = new ArrayList<>();
                    for (Valor v : fila.lista()) {
                        if (!Entorno.compatible(declarado.elemento, v))
                            e.errorSemantico("el vector '" + id + "' es de " + declarado.elemento.nombre()
                                    + " y contiene " + Entorno.tipoDe(v), linea, columna);
                        filaOut.add(v);
                    }
                    out.add(new Valor(Tipo.vector(declarado.elemento, 1), filaOut));
                }
            }
            return new Valor(declarado, out);
        }
        public String etiquetaAst() {
            return (esConst ? "const " : "let ") + id + " : " + tipo.nombre();
        }
        public List<Nodo> hijosAst() { return hijos(init); }
    }

    public static class AsignacionVariable implements Instruccion {
        public final String id; public final Expresion expr; public final int linea, columna;
        public AsignacionVariable(String id, Object expr, int l, int c) {
            this.id = id; this.expr = (Expresion) expr; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Simbolo s = e.obtener(id, linea, columna);
            if (!s.mutable)
                e.errorSemantico("no se puede modificar la constante '" + id + "'", linea, columna);
            Valor v = expr.evaluar(e);
            if (!Entorno.compatible(s.tipo, v))
                e.errorSemantico("no se puede asignar " + Entorno.tipoDe(v)
                        + " a '" + id + "' de tipo " + s.tipo.nombre(), linea, columna);
            s.valor = v;
        }
        public String etiquetaAst() { return "asignacion " + id + " ="; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }

    public static class AsignacionVector implements Instruccion {
        public final String id; public final Expresion idx1, idx2, expr; public final int linea, columna;
        public AsignacionVector(String id, Object idx1, Object idx2, Object expr, int l, int c) {
            this.id = id; this.idx1 = (Expresion) idx1; this.idx2 = (Expresion) idx2;
            this.expr = (Expresion) expr; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Simbolo s = e.obtener(id, linea, columna);
            if (!s.mutable)
                e.errorSemantico("no se puede modificar el vector const '" + id + "'", linea, columna);
            if (s.tipo.cat != Tipo.Cat.VECTOR)
                e.errorSemantico("'" + id + "' no es un vector", linea, columna);
            List<Valor> data = s.valor.lista();
            int i = indice(idx1.evaluar(e), e, linea, columna);
            if (i < 0 || i >= data.size())
                e.errorSemantico("indice " + i + " fuera de rango en '" + id + "'", linea, columna);
            Valor v = expr.evaluar(e);
            if (idx2 == null) {
                if (!Entorno.compatible(s.tipo.elemento, v))
                    e.errorSemantico("no se puede asignar " + Entorno.tipoDe(v)
                            + " en un vector de " + s.tipo.elemento.nombre(), linea, columna);
                data.set(i, v);
            } else {
                Valor fila = data.get(i);
                if (fila.tipo.cat != Tipo.Cat.VECTOR)
                    e.errorSemantico("'" + id + "' no es un vector de dos dimensiones", linea, columna);
                List<Valor> interno = fila.lista();
                int j = indice(idx2.evaluar(e), e, linea, columna);
                if (j < 0 || j >= interno.size())
                    e.errorSemantico("indice " + j + " fuera de rango en '" + id + "'", linea, columna);
                if (!Entorno.compatible(s.tipo.elemento, v))
                    e.errorSemantico("no se puede asignar " + Entorno.tipoDe(v)
                            + " en un vector de " + s.tipo.elemento.nombre(), linea, columna);
                interno.set(j, v);
            }
        }
        public String etiquetaAst() { return "asignacion vector " + id; }
        public List<Nodo> hijosAst() { return hijos(idx1, idx2, expr); }
    }

    public static class AsignacionCampo implements Instruccion {
        public final String id, campo; public final Expresion expr; public final int linea, columna;
        public AsignacionCampo(String id, String campo, Object expr, int l, int c) {
            this.id = id; this.campo = campo; this.expr = (Expresion) expr; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Simbolo s = e.obtener(id, linea, columna);
            if (!s.mutable)
                e.errorSemantico("no se puede modificar el struct const '" + id + "'", linea, columna);
            if (s.tipo.cat != Tipo.Cat.STRUCT)
                e.errorSemantico("'" + id + "' no es un struct", linea, columna);
            DeclaracionStruct def = e.contexto.structs.get(s.tipo.structName);
            CampoStruct cs = def == null ? null : def.buscar(campo);
            if (cs == null)
                e.errorSemantico("el struct '" + id + "' no tiene el campo '" + campo + "'", linea, columna);
            Valor v = expr.evaluar(e);
            if (!Entorno.compatible(cs.tipo, v))
                e.errorSemantico("el campo '" + campo + "' es " + cs.tipo.nombre()
                        + " y recibio " + Entorno.tipoDe(v), linea, columna);
            buscarClave(s.valor.campos(), campo, v);
        }
        public String etiquetaAst() { return "asignacion " + id + "." + campo; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }

    public static class IncDec implements Instruccion {
        public final String id; public final boolean incremento; public final int linea, columna;
        public IncDec(String id, boolean incremento, int l, int c) {
            this.id = id; this.incremento = incremento; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Simbolo s = e.obtener(id, linea, columna);
            if (!s.mutable)
                e.errorSemantico("no se puede modificar la constante '" + id + "'", linea, columna);
            Valor v = s.valor;
            int d = incremento ? 1 : -1;
            if (v.tipo.cat == Tipo.Cat.INT) s.valor = Valor.vInt((Integer) v.valor + d);
            else if (v.tipo.cat == Tipo.Cat.DOUBLE) s.valor = Valor.vDouble((Double) v.valor + d);
            else e.errorSemantico("el incremento/decremento solo aplica a tipos numericos", linea, columna);
        }
        public String etiquetaAst() { return id + (incremento ? "++" : "--"); }
    }

    public static class ImprimirConsola implements Instruccion {
        public final Expresion expr; public final int linea, columna;
        public ImprimirConsola(Object expr, int l, int c) { this.expr = (Expresion) expr; linea = l; columna = c; }
        public void ejecutar(Entorno e) { e.imprimir(expr.evaluar(e).texto()); }
        public String etiquetaAst() { return "console.log"; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }

    /** Envuelve una expresion (llamada / op de lista) para usarla como instruccion. */
    public static class ExpresionStmt implements Instruccion {
        public final Expresion expr;
        public ExpresionStmt(Object expr) { this.expr = (Expresion) expr; }
        public void ejecutar(Entorno e) { expr.evaluar(e); }
        public String etiquetaAst() { return "expr-stmt"; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }

    public static class If implements Instruccion {
        public final Expresion cond; public final List<Instruccion> cuerpo, sino;
        public final If sinoIf; public final int linea, columna;
        public If(Object cond, Object cuerpo, Object sino, Object sinoIf, int l, int c) {
            this.cond = (Expresion) cond; this.cuerpo = insts(cuerpo); this.sino = insts(sino);
            this.sinoIf = (If) sinoIf; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Valor v = cond.evaluar(e);
            if (v.tipo.cat != Tipo.Cat.BOOL)
                e.errorSemantico("la condicion del if debe ser Booleana", linea, columna);
            if ((Boolean) v.valor) A.ejecutar(cuerpo, e.crearHijo("if"));
            else if (sinoIf != null) sinoIf.ejecutar(e);
            else if (!sino.isEmpty()) A.ejecutar(sino, e.crearHijo("else"));
        }
        public String etiquetaAst() { return "if"; }
        public List<Nodo> hijosAst() { return hijos(cond, cuerpo, sino, sinoIf); }
    }

    public static class Match implements Instruccion {
        public final Expresion selector; public final List<CasoMatch> casos;
        public final List<Instruccion> defecto; public final int linea, columna;
        @SuppressWarnings("unchecked")
        public Match(Object selector, Object casos, Object defecto, int l, int c) {
            this.selector = (Expresion) selector;
            this.casos = casos == null ? new ArrayList<>() : (List<CasoMatch>) casos;
            this.defecto = insts(defecto); linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Valor sel = selector.evaluar(e);
            for (CasoMatch cm : casos) {
                Valor cv = cm.valor.evaluar(e);
                Valor eq = Operaciones.relacional("==", sel, cv, e, linea, columna);
                if (eq != null && (Boolean) eq.valor) {
                    A.ejecutar(cm.cuerpo, e.crearHijo("match"));
                    return;   // sin fall-through (5.15.2)
                }
            }
            if (defecto != null && !defecto.isEmpty())
                A.ejecutar(defecto, e.crearHijo("match-default"));
        }
        public String etiquetaAst() { return "match"; }
        public List<Nodo> hijosAst() {
            List<Nodo> h = hijos(selector, casos);
            if (defecto != null && !defecto.isEmpty()) {
                h.add(new Nodo() {
                    public String etiquetaAst() { return "default =>"; }
                    public List<Nodo> hijosAst() { return hijos(defecto); }
                });
            }
            return h;
        }
    }

    public static class While implements Instruccion {
        public final Expresion cond; public final List<Instruccion> cuerpo; public final int linea, columna;
        public While(Object cond, Object cuerpo, int l, int c) {
            this.cond = (Expresion) cond; this.cuerpo = insts(cuerpo); linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            while (true) {
                Valor v = cond.evaluar(e);
                if (v.tipo.cat != Tipo.Cat.BOOL)
                    e.errorSemantico("la condicion del while debe ser Booleana", linea, columna);
                if (!(Boolean) v.valor) break;
                try {
                    A.ejecutar(cuerpo, e.crearHijo("while"));
                } catch (Senales.Break b) { break; }
                catch (Senales.Continue ct) { /* siguiente iteracion */ }
            }
        }
        public String etiquetaAst() { return "while"; }
        public List<Nodo> hijosAst() { return hijos(cond, cuerpo); }
    }

    public static class For implements Instruccion {
        public final Instruccion init, update; public final Expresion cond;
        public final List<Instruccion> cuerpo; public final int linea, columna;
        public For(Object init, Object cond, Object update, Object cuerpo, int l, int c) {
            this.init = (Instruccion) init; this.cond = (Expresion) cond;
            this.update = (Instruccion) update; this.cuerpo = insts(cuerpo); linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            Entorno amb = e.crearHijo("for");
            init.ejecutar(amb);
            while (true) {
                Valor v = cond.evaluar(amb);
                if (v.tipo.cat != Tipo.Cat.BOOL)
                    e.errorSemantico("la condicion del for debe ser Booleana", linea, columna);
                if (!(Boolean) v.valor) break;
                try {
                    A.ejecutar(cuerpo, amb.crearHijo("for-body"));
                } catch (Senales.Break b) { break; }
                catch (Senales.Continue ct) { /* cae a la actualizacion */ }
                update.ejecutar(amb);
            }
        }
        public String etiquetaAst() { return "for"; }
        public List<Nodo> hijosAst() { return hijos(init, cond, update, cuerpo); }
    }

    public static class DoWhile implements Instruccion {
        public final List<Instruccion> cuerpo; public final Expresion cond; public final int linea, columna;
        public DoWhile(Object cuerpo, Object cond, int l, int c) {
            this.cuerpo = insts(cuerpo); this.cond = (Expresion) cond; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) {
            while (true) {
                try {
                    A.ejecutar(cuerpo, e.crearHijo("do"));
                } catch (Senales.Break b) { break; }
                catch (Senales.Continue ct) { /* evalua condicion */ }
                Valor v = cond.evaluar(e);
                if (v.tipo.cat != Tipo.Cat.BOOL)
                    e.errorSemantico("la condicion del do-while debe ser Booleana", linea, columna);
                if (!(Boolean) v.valor) break;   // se repite MIENTRAS sea verdadera (5.16.3)
            }
        }
        public String etiquetaAst() { return "do-while"; }
        public List<Nodo> hijosAst() { return hijos(cuerpo, cond); }
    }

    /** No hace nada por si solo: ejecutar() lanza la senal que el while/for/do
     *  mas cercano atrapa (ver A.While/For/DoWhile) para cortar la iteracion. */
    public static class Break implements Instruccion {
        public final int linea, columna;
        public Break(int l, int c) { linea = l; columna = c; }
        public void ejecutar(Entorno e) { throw new Senales.Break(linea, columna); }
        public String etiquetaAst() { return "break"; }
    }
    /** Igual que Break, pero la senal que atrapa el ciclo indica "salta a la
     *  siguiente iteracion" en vez de "termina el ciclo". */
    public static class Continue implements Instruccion {
        public final int linea, columna;
        public Continue(int l, int c) { linea = l; columna = c; }
        public void ejecutar(Entorno e) { throw new Senales.Continue(linea, columna); }
        public String etiquetaAst() { return "continue"; }
    }
    /** Evalua la expresion (si la hay) ANTES de lanzar la senal: el valor viaja
     *  dentro de la excepcion y lo recupera invocar() al atraparla (ver abajo). */
    public static class Return implements Instruccion {
        public final Expresion expr; public final int linea, columna;
        public Return(Object expr, int l, int c) { this.expr = (Expresion) expr; linea = l; columna = c; }
        public void ejecutar(Entorno e) {
            Valor v = (expr == null) ? null : expr.evaluar(e);
            throw new Senales.Retorno(v, linea, columna);
        }
        public String etiquetaAst() { return "return"; }
        public List<Nodo> hijosAst() { return hijos(expr); }
    }

    public static class DeclaracionFuncion implements Instruccion {
        public final Tipo tipoRetorno; public final String id;
        public final List<Parametro> params; public final List<Instruccion> cuerpo;
        public final int linea, columna;
        @SuppressWarnings("unchecked")
        public DeclaracionFuncion(Tipo tipoRetorno, String id, Object params, Object cuerpo, int l, int c) {
            this.tipoRetorno = tipoRetorno; this.id = id;
            this.params = params == null ? new ArrayList<>() : (List<Parametro>) params;
            this.cuerpo = insts(cuerpo); linea = l; columna = c;
        }
        public void ejecutar(Entorno e) { /* se registra en la 1a pasada; no ejecuta aqui */ }
        public String etiquetaAst() {
            return (tipoRetorno.cat == Tipo.Cat.VOID ? "method " : "function " + tipoRetorno.nombre() + " ") + id;
        }
        public List<Nodo> hijosAst() { return hijos(params, cuerpo); }
    }

    public static class DeclaracionStruct implements Instruccion {
        public final String id; public final List<CampoStruct> campos; public final int linea, columna;
        @SuppressWarnings("unchecked")
        public DeclaracionStruct(String id, Object campos, int l, int c) {
            this.id = id; this.campos = (List<CampoStruct>) campos; linea = l; columna = c;
        }
        public void ejecutar(Entorno e) { /* se registra en la 1a pasada */ }
        public CampoStruct buscar(String nombre) {
            for (CampoStruct cs : campos) if (cs.id.equalsIgnoreCase(nombre)) return cs;
            return null;
        }
        public String etiquetaAst() { return "struct " + id; }
        public List<Nodo> hijosAst() { return hijos(campos); }
    }

    public static class RunMain implements Instruccion {
        public final String id; public final List<Argumento> args; public final int linea, columna;
        @SuppressWarnings("unchecked")
        public RunMain(String id, Object args, int l, int c) {
            this.id = id; this.args = args == null ? new ArrayList<>() : (List<Argumento>) args;
            linea = l; columna = c;
        }
        public void ejecutar(Entorno e) { invocar(id, args, e, linea, columna); }
        public String etiquetaAst() { return "RUN_MAIN " + id; }
        public List<Nodo> hijosAst() { return hijos(args); }
    }

    /** Raiz sintetica del AST (para el reporte 6.3). */
    public static class Programa implements Nodo {
        public final List<Instruccion> instrucciones;
        public Programa(Object instrucciones) { this.instrucciones = insts(instrucciones); }
        public String etiquetaAst() { return "PROGRAMA"; }
        public List<Nodo> hijosAst() { return hijos(instrucciones); }
    }

    /* ============================================================
     * INVOCACION (5.23): parametros por identificador, alcance global.
     * ============================================================ */
    private static Valor invocar(String id, List<Argumento> args, Entorno caller, int l, int c) {
        Contexto ctx = caller.contexto;
        DeclaracionFuncion f = ctx.funciones.get(id.toLowerCase());
        if (f == null)
            caller.errorSemantico("no existe la funcion o metodo '" + id + "'", l, c);

        // argumentos por nombre
        Map<String, Argumento> mapaArgs = new LinkedHashMap<>();
        for (Argumento a : args) {
            if (mapaArgs.containsKey(a.id.toLowerCase()))
                caller.errorSemantico("argumento '" + a.id + "' repetido en la llamada a '" + id + "'", a.linea, a.columna);
            mapaArgs.put(a.id.toLowerCase(), a);
        }
        // parametro desconocido
        for (Argumento a : args) {
            boolean existe = false;
            for (Parametro p : f.params) if (p.id.equalsIgnoreCase(a.id)) { existe = true; break; }
            if (!existe)
                caller.errorSemantico("'" + id + "' no tiene un parametro llamado '" + a.id + "'", a.linea, a.columna);
        }

        Entorno local = ctx.global.crearHijo(f.id);
        for (Parametro p : f.params) {
            Valor v;
            Argumento a = mapaArgs.get(p.id.toLowerCase());
            if (a != null) v = a.expr.evaluar(caller);
            else if (p.defecto != null) v = p.defecto.evaluar(caller);
            else { caller.errorSemantico("falta el argumento '" + p.id + "' en la llamada a '" + id + "'", l, c); return null; }
            if (!Entorno.compatible(p.tipo, v))
                caller.errorSemantico("el parametro '" + p.id + "' de '" + id + "' espera "
                        + p.tipo.nombre() + " y recibio " + Entorno.tipoDe(v), l, c);
            local.declarar(new Simbolo(p.id, Simbolo.categoriaDe(p.tipo), p.tipo, true,
                    f.id, v, f.linea + 1, f.columna + 1), l, c);
        }

        Valor retorno = null;
        try {
            A.ejecutar(f.cuerpo, local);
            // si el cuerpo termina sin "return", retorno queda null (revisado abajo)
        } catch (Senales.Retorno r) {
            retorno = r.valor;   // "return;" trae valor null; se valida segun tipoRetorno mas abajo
        }

        if (f.tipoRetorno.cat == Tipo.Cat.VOID) {
            if (retorno != null && retorno.tipo.cat != Tipo.Cat.VOID)
                caller.errorSemantico("el metodo '" + id + "' es void y no puede retornar un valor", l, c);
            return Valor.VOID;
        }
        if (retorno == null)
            caller.errorSemantico("la funcion '" + id + "' debe retornar un " + f.tipoRetorno.nombre(), l, c);
        if (!Entorno.compatible(f.tipoRetorno, retorno))
            caller.errorSemantico("la funcion '" + id + "' retorna " + Entorno.tipoDe(retorno)
                    + " pero se declaro " + f.tipoRetorno.nombre(), l, c);
        return retorno;
    }

    /* ---- utilidades de structs (busqueda case-insensitive) ---- */
    private static Valor buscarCampo(Map<String, Valor> campos, String nombre) {
        for (Map.Entry<String, Valor> e : campos.entrySet())
            if (e.getKey().equalsIgnoreCase(nombre)) return e.getValue();
        return null;
    }
    private static void buscarClave(Map<String, Valor> campos, String nombre, Valor nuevo) {
        for (String k : campos.keySet())
            if (k.equalsIgnoreCase(nombre)) { campos.put(k, nuevo); return; }
    }
}
