package compscript.interprete;

/**
 * Senales de control de flujo, implementadas como excepciones (patron
 * clasico de interpretes tree-walk). Los ciclos capturan Break/Continue;
 * las llamadas capturan Retorno. Sin stack-trace: son control, no errores.
 */
public final class Senales {

    private Senales() {}

    /** break; -> sale del ciclo mas cercano (5.17.1). */
    public static class Break extends RuntimeException {
        public final int linea, columna;
        public Break(int l, int c) { super(null, null, false, false); linea = l; columna = c; }
    }

    /** continue; -> salta a la siguiente iteracion (5.17.2). */
    public static class Continue extends RuntimeException {
        public final int linea, columna;
        public Continue(int l, int c) { super(null, null, false, false); linea = l; columna = c; }
    }

    /** return [expr]; -> termina la funcion/metodo y devuelve el valor (5.17.3). */
    public static class Retorno extends RuntimeException {
        public final Valor valor;   // null si es "return;"
        public final int linea, columna;
        public Retorno(Valor valor, int l, int c) {
            super(null, null, false, false);
            this.valor = valor; linea = l; columna = c;
        }
    }
}
