package compscript.interprete;

import compscript.interprete.Tipo.Cat;

/**
 * Semantica de los operadores (5.5, 5.6, 5.7) y de los casteos (5.13).
 * Implementa LITERALMENTE las 8 tablas de compatibilidad aritmetica y la
 * matriz relacional del enunciado: toda combinacion no listada es un error
 * de tipo semantico que ABORTA la ejecucion (via Entorno.errorSemantico).
 */
public class Operaciones {

    /* ================= aritmetica (operadores infijos) ================= */

    public static Valor aritmetica(String op, Valor a, Valor b, Entorno e, int l, int c) {
        switch (op) {
            case "+": return suma(a, b, e, l, c);
            case "-": return restaMult(op, a, b, e, l, c);
            case "*": return restaMult(op, a, b, e, l, c);
            case "/": return division(a, b, e, l, c);
            case "^": return potencia(a, b, e, l, c);
            case "$": return raiz(a, b, e, l, c);
            case "%": return modulo(a, b, e, l, c);
            default:  return null;
        }
    }

    /* ---- 5.5.1 Suma (unica que admite bool, char y cadena) ---- */
    private static Valor suma(Valor a, Valor b, Entorno e, int l, int c) {
        Cat x = a.tipo.cat, y = b.tipo.cat;
        Cat res = null;
        if (x == Cat.STRING || y == Cat.STRING) {
            // Cadena con cualquiera de los otros primitivos -> Cadena
            if (esPrimitivo(x) && esPrimitivo(y)) res = Cat.STRING;
        } else if (x == Cat.CHAR && y == Cat.CHAR) {
            res = Cat.STRING;                       // char + char -> cadena
        } else if (x == Cat.INT && y == Cat.INT) {
            res = Cat.INT;
        } else if ((x == Cat.INT || x == Cat.DOUBLE || x == Cat.BOOL || x == Cat.CHAR)
                && (y == Cat.INT || y == Cat.DOUBLE || y == Cat.BOOL || y == Cat.CHAR)) {
            // numerico/bool/char: si alguno es double -> Decimal, si no -> Entero
            boolean hayDouble = (x == Cat.DOUBLE || y == Cat.DOUBLE);
            boolean boolConBool = (x == Cat.BOOL && y == Cat.BOOL);
            boolean boolConChar = (x == Cat.BOOL && y == Cat.CHAR) || (x == Cat.CHAR && y == Cat.BOOL);
            if (boolConBool || boolConChar) res = null;   // invalidas en la tabla
            else res = hayDouble ? Cat.DOUBLE : Cat.INT;
        }
        if (res == null) return error("suma", a, b, e, l, c);
        if (res == Cat.STRING) return Valor.vString(concat(a) + concat(b));
        double r = a.numero() + b.numero();
        return res == Cat.INT ? Valor.vInt((int) r) : Valor.vDouble(r);
    }

    /* ---- 5.5.2 Resta / 5.5.3 Multiplicacion (Entero, Decimal, Caracter) ---- */
    private static Valor restaMult(String op, Valor a, Valor b, Entorno e, int l, int c) {
        Cat res = tablaRestaMult(a.tipo.cat, b.tipo.cat);
        if (res == null) return error(op.equals("-") ? "resta" : "multiplicacion", a, b, e, l, c);
        double r = op.equals("-") ? a.numero() - b.numero() : a.numero() * b.numero();
        return res == Cat.INT ? Valor.vInt((int) r) : Valor.vDouble(r);
    }

    private static Cat tablaRestaMult(Cat x, Cat y) {
        if (!enGrupo(x, Cat.INT, Cat.DOUBLE, Cat.CHAR) || !enGrupo(y, Cat.INT, Cat.DOUBLE, Cat.CHAR))
            return null;
        if (x == Cat.CHAR && y == Cat.CHAR) return null;       // char - char invalido
        return (x == Cat.DOUBLE || y == Cat.DOUBLE) ? Cat.DOUBLE : Cat.INT;
    }

    /* ---- 5.5.4 Division (siempre Decimal) ---- */
    private static Valor division(Valor a, Valor b, Entorno e, int l, int c) {
        Cat x = a.tipo.cat, y = b.tipo.cat;
        if (!enGrupo(x, Cat.INT, Cat.DOUBLE, Cat.CHAR) || !enGrupo(y, Cat.INT, Cat.DOUBLE, Cat.CHAR)
                || (x == Cat.CHAR && y == Cat.CHAR))
            return error("division", a, b, e, l, c);
        if (b.numero() == 0) e.errorSemantico("division entre cero", l, c);
        return Valor.vDouble(a.numero() / b.numero());
    }

    /* ---- 5.5.5 Potencia (Entero/Decimal) ---- */
    private static Valor potencia(Valor a, Valor b, Entorno e, int l, int c) {
        Cat x = a.tipo.cat, y = b.tipo.cat;
        if (!enGrupo(x, Cat.INT, Cat.DOUBLE) || !enGrupo(y, Cat.INT, Cat.DOUBLE))
            return error("potencia", a, b, e, l, c);
        double r = Math.pow(a.numero(), b.numero());
        return (x == Cat.INT && y == Cat.INT) ? Valor.vInt((int) r) : Valor.vDouble(r);
    }

    /* ---- 5.5.6 Raiz (a $ b = raiz b-esima de a; siempre Decimal) ---- */
    private static Valor raiz(Valor a, Valor b, Entorno e, int l, int c) {
        Cat x = a.tipo.cat, y = b.tipo.cat;
        if (!enGrupo(x, Cat.INT, Cat.DOUBLE) || !enGrupo(y, Cat.INT, Cat.DOUBLE))
            return error("raiz", a, b, e, l, c);
        if (b.numero() == 0) e.errorSemantico("raiz de indice cero", l, c);
        return Valor.vDouble(Math.pow(a.numero(), 1.0 / b.numero()));
    }

    /* ---- 5.5.7 Modulo (Entero/Decimal; resultado Decimal) ---- */
    private static Valor modulo(Valor a, Valor b, Entorno e, int l, int c) {
        Cat x = a.tipo.cat, y = b.tipo.cat;
        if (!enGrupo(x, Cat.INT, Cat.DOUBLE) || !enGrupo(y, Cat.INT, Cat.DOUBLE))
            return error("modulo", a, b, e, l, c);
        if (b.numero() == 0) e.errorSemantico("modulo entre cero", l, c);
        return Valor.vDouble(a.numero() % b.numero());
    }

    /* ---- 5.5.8 Negacion unaria (Entero/Decimal) ---- */
    public static Valor negacion(Valor a, Entorno e, int l, int c) {
        if (a.tipo.cat == Cat.INT) return Valor.vInt(-(Integer) a.valor);
        if (a.tipo.cat == Cat.DOUBLE) return Valor.vDouble(-(Double) a.valor);
        e.errorSemantico("no se puede negar un valor de tipo " + a.tipo.nombre().toUpperCase(), l, c);
        return null;
    }

    /* ================= relacionales (5.6) ================= */

    public static Valor relacional(String op, Valor a, Valor b, Entorno e, int l, int c) {
        Cat x = a.tipo.cat, y = b.tipo.cat;
        boolean ok = (esNum(x) && esNum(y))            // Entero/Decimal/Caracter entre si
                || (x == Cat.BOOL && y == Cat.BOOL)     // Booleano solo con Booleano
                || (x == Cat.STRING && y == Cat.STRING); // Cadena solo con Cadena
        if (!ok) return error("comparacion (" + op + ")", a, b, e, l, c);

        int cmp;
        if (x == Cat.STRING) cmp = ((String) a.valor).compareTo((String) b.valor);
        else cmp = Double.compare(a.numero(), b.numero());

        boolean r;
        switch (op) {
            case "==": r = cmp == 0; break;
            case "!=": r = cmp != 0; break;
            case "<":  r = cmp < 0;  break;
            case "<=": r = cmp <= 0; break;
            case ">":  r = cmp > 0;  break;
            case ">=": r = cmp >= 0; break;
            default:   r = false;
        }
        return Valor.vBool(r);
    }

    /* ================= logicos (5.7) ================= */

    public static Valor logica(String op, Valor a, Valor b, Entorno e, int l, int c) {
        if (a.tipo.cat != Cat.BOOL || b.tipo.cat != Cat.BOOL)
            return error(op.equals("&&") ? "AND (&&)" : "OR (||)", a, b, e, l, c);
        boolean ba = (Boolean) a.valor, bb = (Boolean) b.valor;
        return Valor.vBool(op.equals("&&") ? (ba && bb) : (ba || bb));
    }

    public static Valor negacionLogica(Valor a, Entorno e, int l, int c) {
        if (a.tipo.cat != Cat.BOOL)
            e.errorSemantico("el operador ! requiere un Booleano, no " + a.tipo.nombre().toUpperCase(), l, c);
        return Valor.vBool(!(Boolean) a.valor);
    }

    /* ================= casteos (5.13) ================= */

    public static Valor cast(Valor v, Tipo destino, Entorno e, int l, int c) {
        Cat o = v.tipo.cat, d = destino.cat;
        if (o == d) return new Valor(destino, v.valor);          // identidad
        if (o == Cat.INT  && d == Cat.DOUBLE) return Valor.vDouble((Integer) v.valor);
        if (o == Cat.DOUBLE && d == Cat.INT)  return Valor.vInt((int) (double) (Double) v.valor);
        if (o == Cat.INT  && d == Cat.CHAR)   return Valor.vChar((char) (int) (Integer) v.valor);
        if (o == Cat.CHAR && d == Cat.INT)    return Valor.vInt((Character) v.valor);
        if (o == Cat.CHAR && d == Cat.DOUBLE) return Valor.vDouble((int) (Character) v.valor);
        e.errorSemantico("no se puede castear " + v.tipo.nombre().toUpperCase()
                + " a " + destino.nombre().toUpperCase(), l, c);
        return null;
    }

    /* ================= helpers ================= */

    private static boolean esPrimitivo(Cat c) {
        return c == Cat.INT || c == Cat.DOUBLE || c == Cat.BOOL || c == Cat.CHAR || c == Cat.STRING;
    }
    private static boolean esNum(Cat c) {
        return c == Cat.INT || c == Cat.DOUBLE || c == Cat.CHAR;
    }
    private static boolean enGrupo(Cat c, Cat... permitidos) {
        for (Cat p : permitidos) if (c == p) return true;
        return false;
    }

    /** Representacion de un valor primitivo para la concatenacion de cadenas. */
    private static String concat(Valor v) {
        return v.texto();
    }

    private static Valor error(String operacion, Valor a, Valor b, Entorno e, int l, int c) {
        e.errorSemantico("No se puede realizar " + operacion + " entre "
                + a.tipo.nombre().toUpperCase() + " y " + b.tipo.nombre().toUpperCase(), l, c);
        return null;
    }
}
