package dataforge.interprete;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;

/**
 * La SEMÁNTICA de las operaciones del lenguaje (5.7 y 5.8).
 * Todas validan tipos: DataForge solo opera aritmética/estadística
 * sobre double (enunciado 5.7) — violar eso es error SEMÁNTICO,
 * la tercera familia de errores (ni léxico ni sintáctico).
 *
 * Convención de propagación: si un operando llega null es que una
 * expresión interior YA reportó su error — devolvemos null sin
 * reportar de nuevo (evita cascadas de errores duplicados).
 */
public class Operaciones {

    /* ============ aritméticas: SUM, RES, MUL, DIV, MOD ============ */

    public static Object aritmetica(String op, Object a, Object b,
                                    Entorno ent, int l, int c) {
        if (a == null || b == null) return null;
        if (!(a instanceof Double) || !(b instanceof Double)) {
            ent.error("Semántico", op + " solo acepta valores double (recibió "
                    + Entorno.formatear(a) + " y " + Entorno.formatear(b) + ")", l, c);
            return null;
        }
        double x = (Double) a, y = (Double) b;
        switch (op) {
            case "SUM": return x + y;
            case "RES": return x - y;
            case "MUL": return x * y;
            case "MOD":
                if (y == 0) {
                    ent.error("Semántico", "módulo entre cero", l, c);
                    return null;
                }
                return x % y;
            case "DIV":
                if (y == 0) {
                    ent.error("Semántico", "división entre cero", l, c);
                    return null;
                }
                return x / y;
        }
        return null;  // inalcanzable: la gramática solo produce esos 5
    }

    /* ============ estadísticas sobre arreglo double (5.8) ============ */

    public static Object estadistica(String fn, ArrayList<Object> arr,
                                     Entorno ent, int l, int c) {
        if (arr == null) return null;
        ArrayList<Double> datos = new ArrayList<>();
        for (Object v : arr) {
            if (v == null) return null;
            if (!(v instanceof Double)) {
                ent.error("Semántico", fn + " requiere un arreglo de tipo double", l, c);
                return null;
            }
            datos.add((Double) v);
        }
        if (datos.isEmpty()) {
            ent.error("Semántico", fn + " recibió un arreglo vacío", l, c);
            return null;
        }
        switch (fn) {
            case "Media":    return media(datos);
            case "Mediana":  return mediana(datos);
            case "Moda":     return moda(datos);
            case "Varianza": return varianza(datos);
            case "Max":      return Collections.max(datos);
            case "Min":      return Collections.min(datos);
        }
        return null;
    }

    private static double media(ArrayList<Double> d) {
        double s = 0;
        for (double v : d) s += v;
        return s / d.size();
    }

    private static double mediana(ArrayList<Double> d) {
        ArrayList<Double> ord = new ArrayList<>(d);
        Collections.sort(ord);
        int n = ord.size();
        return (n % 2 == 1) ? ord.get(n / 2)
                            : (ord.get(n / 2 - 1) + ord.get(n / 2)) / 2.0;
    }

    /** Valor más frecuente; en empate, el primero que alcanzó esa frecuencia. */
    private static double moda(ArrayList<Double> d) {
        LinkedHashMap<Double, Integer> freq = new LinkedHashMap<>();
        for (double v : d) freq.merge(v, 1, Integer::sum);
        double moda = d.get(0);
        int max = 0;
        for (var e : freq.entrySet())
            if (e.getValue() > max) { max = e.getValue(); moda = e.getKey(); }
        return moda;
    }

    /** Varianza poblacional: Σ(x−μ)² / n. */
    private static double varianza(ArrayList<Double> d) {
        double m = media(d), s = 0;
        for (double v : d) s += (v - m) * (v - m);
        return s / d.size();
    }

    /* ============ frecuencias para el histograma (5.10.3) ============ */

    /** Frecuencia de cada valor, ordenado ascendente (2→3, 5→2, 7→1…).
     *  La usan la tabla de consola del histograma Y su gráfica. */
    public static LinkedHashMap<Double, Integer> frecuencias(ArrayList<Double> datos) {
        ArrayList<Double> ord = new ArrayList<>(datos);
        Collections.sort(ord);
        LinkedHashMap<Double, Integer> f = new LinkedHashMap<>();
        for (double v : ord) f.merge(v, 1, Integer::sum);
        return f;
    }
}
