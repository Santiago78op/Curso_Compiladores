package conjanalyzer.interprete;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashSet;
import java.util.List;

/**
 * Motor de simplificacion de operaciones (seccion 7 del enunciado).
 *
 * Reescribe el arbol de una operacion aplicando propiedades de la teoria de
 * conjuntos, de forma REPETIDA y en POSTORDEN (primero los hijos, luego el
 * nodo actual), hasta llegar a un punto fijo o al tope de iteraciones.
 *
 * Leyes que APLICA como transformacion (reducen la complejidad del arbol):
 *   - Ley del doble complemento  : ^^X                        → X
 *   - Leyes de DeMorgan          : ^(X U Y)                   → ^X & ^Y   (y su dual)
 *   - Propiedades idempotentes   : X U X, X & X                → X
 *   - Propiedades de absorcion   : X U (X & Y)                 → X   (y sus variantes)
 *   - Propiedades distributivas : (X & Y) U (X & Z)           → X & (Y U Z)   (y su dual)
 *
 * Decisiones de diseno (la parte con mas ambiguedad del enunciado):
 *
 *  1. DeMorgan solo se aplica cuando NETO simplifica: exigimos que el operando
 *     del complemento (X o Y) sea a su vez un complemento, de modo que al
 *     empujar el "^" hacia adentro se genere un "^^" que luego cancela por
 *     doble complemento. Asi nunca "infla" un arbol como ^(A U B) que no gana
 *     nada; solo dispara en cadenas tipo ^(^(A&B) & ^C) → (A&B) U C, que es
 *     exactamente el ejemplo de la seccion 5.4.
 *
 *  2. Las propiedades CONMUTATIVA y ASOCIATIVA no se aplican como reescritura
 *     por si solas (reordenar/regrupar no simplifica y podria ciclar). Se usan
 *     como AYUDA DE COMPARACION: para reconocer X U X o la absorcion aunque los
 *     operandos vengan en distinto orden o anidados, se comparan los subarboles
 *     por su forma canonica (operandos de U/& aplanados y ordenados). Cuando esa
 *     comparacion "no trivial" es la que habilita una idempotencia o absorcion,
 *     se reportan tambien "Propiedades conmutativas" y/o "Propiedades asociativas"
 *     para dejar constancia de que se usaron.
 *
 *  3. Se comparan las HOJAS por nombre de conjunto (un conjunto no puede
 *     definirse en terminos de otro, asi que los nombres son los atomos).
 *
 *  4. Distributiva: igual que con DeMorgan, solo se aplica en el sentido que
 *     REDUCE el arbol (factorizar el termino comun), nunca en el sentido de
 *     "expandir" (X & (Y U Z) → (X & Y) U (X & Z) aumenta la cantidad de hojas
 *     y no simplifica nada). Se detecta cuando los dos operandos de un U/&
 *     son a su vez nodos del operador opuesto y comparten un operando comun
 *     (comparado con {@link #equivalentes}, o sea admitiendo conmutativa /
 *     asociativa igual que el resto de las reglas).
 */
public class Simplificador {

    private static final int MAX_ITER = 100;
    private LinkedHashSet<String> leyes;

    public ResultadoSimplificacion simplificar(NodoOperacion original) {
        leyes = new LinkedHashSet<>();
        NodoOperacion actual = original.copia();
        String antes;
        int guarda = 0;
        do {
            antes = actual.toPrefijo();
            actual = paso(actual);
        } while (!actual.toPrefijo().equals(antes) && ++guarda < MAX_ITER);
        return new ResultadoSimplificacion(new ArrayList<>(leyes), actual, !leyes.isEmpty());
    }

    /* Un "paso" = una pasada bottom-up que aplica a lo sumo una regla por nodo.
       El bucle de arriba repite pasadas hasta que el arbol deja de cambiar. */
    private NodoOperacion paso(NodoOperacion n) {
        if (n.esHoja()) return n;

        if (n.esUnario()) {                       // nodo complemento
            NodoOperacion hijo = paso(n.izq);

            if (hijo.esUnario()) {                // ^^X → X
                leyes.add("Ley del doble complemento");
                return hijo.izq;
            }
            // DeMorgan guardado: ^(X U Y) → ^X & ^Y  solo si X o Y ya es un ^
            if (esBinarioUnionInter(hijo) && (hijo.izq.esUnario() || hijo.der.esUnario())) {
                leyes.add("Leyes de DeMorgan");
                char nuevo = (hijo.op == 'U') ? '&' : 'U';
                return NodoOperacion.binario(nuevo,
                        NodoOperacion.unario('^', hijo.izq),
                        NodoOperacion.unario('^', hijo.der));
            }
            return NodoOperacion.unario('^', hijo);
        }

        // nodo binario: simplificar hijos primero
        NodoOperacion a = paso(n.izq);
        NodoOperacion b = paso(n.der);

        if (n.op == 'U' || n.op == '&') {
            if (equivalentes(a, b)) {             // X U X → X  /  X & X → X
                leyes.add("Propiedades idempotentes");
                registrarCommAssoc(a, b);
                return a;
            }
            NodoOperacion abs = absorcion(n.op, a, b);
            if (abs != null) return abs;
            NodoOperacion dist = distributiva(n.op, a, b);
            if (dist != null) return dist;
        }
        return NodoOperacion.binario(n.op, a, b);
    }

    /* (X interno Y) op (X interno Z) → X interno (Y op Z), con interno el
       operador opuesto de op. Solo factoriza (reduce hojas); nunca expande. */
    private NodoOperacion distributiva(char op, NodoOperacion a, NodoOperacion b) {
        char interno = (op == 'U') ? '&' : 'U';
        if (!esBinario(a, interno) || !esBinario(b, interno)) return null;

        NodoOperacion factor, restoA, restoB;
        if (equivalentes(a.izq, b.izq))      { factor = a.izq; restoA = a.der; restoB = b.der; }
        else if (equivalentes(a.izq, b.der)) { factor = a.izq; restoA = a.der; restoB = b.izq; }
        else if (equivalentes(a.der, b.izq)) { factor = a.der; restoA = a.izq; restoB = b.der; }
        else if (equivalentes(a.der, b.der)) { factor = a.der; restoA = a.izq; restoB = b.izq; }
        else return null;

        leyes.add("Propiedades distributivas");
        registrarDistributiva(factor, a, b);
        return NodoOperacion.binario(interno, factor, NodoOperacion.binario(op, restoA, restoB));
    }

    /* A op (A op2 B) → A , con op/op2 el par union/interseccion de la absorcion */
    private NodoOperacion absorcion(char op, NodoOperacion a, NodoOperacion b) {
        char interno = (op == 'U') ? '&' : 'U';   // U absorbe sobre &, & absorbe sobre U
        // A op (A interno B)
        if (esBinario(b, interno) && (equivalentes(b.izq, a) || equivalentes(b.der, a))) {
            leyes.add("Propiedades de absorcion");
            registrarAbsorcion(a, b);
            return a;
        }
        // (A interno B) op A   (variante conmutada)
        if (esBinario(a, interno) && (equivalentes(a.izq, b) || equivalentes(a.der, b))) {
            leyes.add("Propiedades de absorcion");
            registrarAbsorcion(b, a);
            return b;
        }
        return null;
    }

    /* ---------- reporte de conmutativa/asociativa (decision de diseno nº 2) ---------- */

    private void registrarCommAssoc(NodoOperacion a, NodoOperacion b) {
        if (!a.toPrefijo().equals(b.toPrefijo())) {          // hubo reordenamiento
            leyes.add("Propiedades conmutativas");
        }
        if (tieneAnidacionMismoOp(a) || tieneAnidacionMismoOp(b)) {
            leyes.add("Propiedades asociativas");
        }
    }

    /* si el factor comun no cayo en la posicion "canonica" (izq de ambos) se
       uso conmutativa; si alguno de los dos lados ya viene anidado con el
       mismo operador interno, se uso asociativa para llegar a esa forma. */
    private void registrarDistributiva(NodoOperacion factor, NodoOperacion a, NodoOperacion b) {
        if (!(equivalentes(a.izq, factor) && equivalentes(b.izq, factor))) {
            leyes.add("Propiedades conmutativas");
        }
        if (tieneAnidacionMismoOp(a) || tieneAnidacionMismoOp(b)) {
            leyes.add("Propiedades asociativas");
        }
    }

    private void registrarAbsorcion(NodoOperacion simple, NodoOperacion compuesto) {
        // si el termino que absorbio no era el operando izquierdo directo, se uso conmutativa
        if (!equivalentes(compuesto.izq, simple)) {
            leyes.add("Propiedades conmutativas");
        }
        if (tieneAnidacionMismoOp(compuesto)) {
            leyes.add("Propiedades asociativas");
        }
    }

    /* ================= equivalencia estructural (comm + assoc) ================= */

    /** Dos subarboles son equivalentes si comparten forma canonica: los
     *  operandos de las cadenas U/& se aplanan y ordenan, asi A U B == B U A
     *  y (A U B) U C == A U (B U C). */
    private boolean equivalentes(NodoOperacion a, NodoOperacion b) {
        return canon(a).equals(canon(b));
    }

    private String canon(NodoOperacion n) {
        if (n.esHoja()) return "'" + n.nombreConj + "'";
        if (n.esUnario()) return "^(" + canon(n.izq) + ")";
        if (n.op == '-') return "-(" + canon(n.izq) + "," + canon(n.der) + ")";
        // U o &: aplanar operandos del mismo operador, canonizar, ordenar y deduplicar
        List<String> ops = new ArrayList<>();
        aplanar(n, n.op, ops);
        List<String> unicos = new ArrayList<>(new LinkedHashSet<>(ops));
        Collections.sort(unicos);
        return n.op + "[" + String.join(",", unicos) + "]";
    }

    private void aplanar(NodoOperacion n, char op, List<String> acum) {
        if (!n.esHoja() && !n.esUnario() && n.op == op) {
            aplanar(n.izq, op, acum);
            aplanar(n.der, op, acum);
        } else {
            acum.add(canon(n));
        }
    }

    /* ================= helpers ================= */

    private boolean esBinario(NodoOperacion n, char op) {
        return !n.esHoja() && !n.esUnario() && n.op == op;
    }

    private boolean esBinarioUnionInter(NodoOperacion n) {
        return !n.esHoja() && !n.esUnario() && (n.op == 'U' || n.op == '&');
    }

    private boolean tieneAnidacionMismoOp(NodoOperacion n) {
        if (n.esHoja() || n.esUnario()) return false;
        if (n.op == 'U' || n.op == '&') {
            if (esBinario(n.izq, n.op) || esBinario(n.der, n.op)) return true;
        }
        return tieneAnidacionMismoOp(n.izq) || (n.der != null && tieneAnidacionMismoOp(n.der));
    }
}
