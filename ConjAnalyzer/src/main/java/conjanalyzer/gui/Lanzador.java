package conjanalyzer.gui;

/**
 * Punto de entrada para correr la GUI desde IDEA con ▶ Run.
 *
 * ¿Por que existe? El launcher de Java revisa si la clase main EXTIENDE
 * Application: si es asi y JavaFX no esta en el module-path, aborta con
 * "JavaFX runtime components are missing". Lanzar desde una clase que NO
 * extiende Application esquiva ese chequeo y funciona con el classpath plano
 * que arma IDEA desde el pom.
 *
 * (La alternativa oficial es `mvn clean javafx:run`, que configura todo solo.)
 */
public class Lanzador {

    public static void main(String[] args) {
        EditorApp.main(args);
    }
}
