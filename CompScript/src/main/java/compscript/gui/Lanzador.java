package compscript.gui;

/**
 * Punto de entrada para correr la GUI desde IntelliJ IDEA con ▶ Run.
 *
 * El launcher de Java aborta con "JavaFX runtime components are missing"
 * cuando la clase main EXTIENDE Application y JavaFX no esta en el
 * module-path. Lanzar desde una clase que NO extiende Application esquiva
 * ese chequeo y funciona con el classpath plano que arma IDEA desde el pom.
 * (La alternativa oficial es `mvn clean javafx:run`.)
 */
public class Lanzador {
    public static void main(String[] args) {
        EditorApp.main(args);
    }
}
