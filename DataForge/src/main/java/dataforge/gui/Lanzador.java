package dataforge.gui;

/**
 * Punto de entrada para correr la GUI desde IDEA con ▶ Run.
 *
 * ¿Por qué existe? El launcher de Java revisa si la clase main
 * EXTIENDE Application: si es así y JavaFX no está en el module-path,
 * aborta con "JavaFX runtime components are missing". Lanzar desde una
 * clase que NO extiende Application esquiva ese chequeo y funciona
 * con el classpath plano que arma IDEA desde el pom.
 *
 * (La alternativa oficial es `mvn javafx:run`, que configura todo solo.)
 */
public class Lanzador {

    public static void main(String[] args) {
        EditorApp.main(args);
    }
}
