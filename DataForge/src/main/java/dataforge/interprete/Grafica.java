package dataforge.interprete;

import java.util.LinkedHashMap;

/**
 * Una gráfica lista para renderizar: tipo + sus atributos ya resueltos.
 * En esta etapa solo se ACUMULAN (la regla "la última instrucción gana"
 * ya está aplicada); el dibujo con JavaFX Charts llega en la Etapa 5.
 */
public class Grafica {

    public final String tipo;  // "graphBar" | "graphPie" | "graphLine" | "Histogram"
    public final LinkedHashMap<String, Object> atributos;

    public Grafica(String tipo, LinkedHashMap<String, Object> atributos) {
        this.tipo = tipo;
        this.atributos = atributos;
    }

    @Override
    public String toString() {
        return tipo + " " + atributos;
    }
}
