package compscript.interprete;

/**
 * Excepcion que ABORTA la ejecucion ante un error semantico.
 *
 * Decision de diseno (enunciado 4.3): a diferencia de DataForge -que
 * seguia ejecutando-, CompScript debe "recuperarse del error y TERMINAR
 * la ejecucion del programa". El error ya quedo registrado en la tabla
 * (6.2) antes de lanzarse; esta excepcion solo desenrolla la pila de
 * llamadas hasta el Interprete, que corta ordenadamente la ejecucion.
 *
 * Se construye sin stack-trace (rendimiento): es control de flujo, no un bug.
 */
public class ErrorSemantico extends RuntimeException {
    public final int linea;
    public final int columna;

    public ErrorSemantico(String mensaje, int linea, int columna) {
        super(mensaje, null, false, false);
        this.linea = linea;
        this.columna = columna;
    }
}
