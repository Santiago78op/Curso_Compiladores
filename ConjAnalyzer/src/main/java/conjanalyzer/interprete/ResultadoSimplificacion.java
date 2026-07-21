package conjanalyzer.interprete;

import java.util.List;

/**
 * Resultado de intentar simplificar una operacion (seccion 5.4 y 7):
 *   - leyes        : nombres (en espanol, como la tabla de la seccion 7)
 *                    de las propiedades que se aplicaron, en orden.
 *   - simplificado : el arbol resultante, ya reescrito.
 *   - seSimplifico : true si se aplico al menos una ley (el arbol cambio).
 *
 * Si {@code seSimplifico} es false, el JSON de salida usa el string literal
 * "No se puede simplificar la operacion" en vez de un objeto.
 */
public class ResultadoSimplificacion {

    public final List<String> leyes;
    public final NodoOperacion simplificado;
    public final boolean seSimplifico;

    public ResultadoSimplificacion(List<String> leyes, NodoOperacion simplificado, boolean seSimplifico) {
        this.leyes = leyes;
        this.simplificado = simplificado;
        this.seSimplifico = seSimplifico;
    }
}
