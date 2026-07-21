package dataforge.gui;

import java.util.ArrayList;

import javafx.scene.Scene;
import javafx.scene.chart.BarChart;
import javafx.scene.chart.CategoryAxis;
import javafx.scene.chart.Chart;
import javafx.scene.chart.LineChart;
import javafx.scene.chart.NumberAxis;
import javafx.scene.chart.PieChart;
import javafx.scene.chart.XYChart;
import javafx.stage.Stage;

import dataforge.interprete.Entorno;
import dataforge.interprete.Grafica;
import dataforge.interprete.Operaciones;

/**
 * Convierte cada Grafica registrada por el Entorno en un JavaFX Chart
 * y lo muestra en su propia ventana (el "EXEC muestra la gráfica en
 * pantalla" del enunciado 5.10).
 *
 * Los atributos ya llegan VALIDADOS (Entorno.validarGrafica): acá
 * los casts son seguros — la GUI dibuja, no vuelve a chequear.
 */
public class Graficador {

    public static void mostrar(Grafica g) {
        Chart chart = switch (g.tipo) {
            case "graphBar"  -> barras(g);
            case "graphPie"  -> pie(g);
            case "graphLine" -> linea(g);
            case "Histogram" -> histograma(g);
            default -> null;
        };
        if (chart == null) return;
        Stage ventana = new Stage();
        ventana.setTitle(g.tipo + " — " + g.atributos.get("titulo"));
        ventana.setScene(new Scene(chart, 640, 480));
        ventana.show();
    }

    /* ---- 5.10.1: barras — ejes de categoría (X) y numérico (Y) ---- */
    private static Chart barras(Grafica g) {
        CategoryAxis x = new CategoryAxis();
        x.setLabel((String) g.atributos.get("titulox"));
        NumberAxis y = new NumberAxis();
        y.setLabel((String) g.atributos.get("tituloy"));

        BarChart<String, Number> chart = new BarChart<>(x, y);
        chart.setTitle((String) g.atributos.get("titulo"));
        chart.setLegendVisible(false);
        chart.getData().add(serie(g.atributos.get("ejex"), g.atributos.get("ejey")));
        return chart;
    }

    /* ---- 5.10.4: línea — mismos ejes, otra geometría ---- */
    private static Chart linea(Grafica g) {
        CategoryAxis x = new CategoryAxis();
        x.setLabel((String) g.atributos.get("titulox"));
        NumberAxis y = new NumberAxis();
        y.setLabel((String) g.atributos.get("tituloy"));

        LineChart<String, Number> chart = new LineChart<>(x, y);
        chart.setTitle((String) g.atributos.get("titulo"));
        chart.setLegendVisible(false);
        chart.getData().add(serie(g.atributos.get("ejex"), g.atributos.get("ejey")));
        return chart;
    }

    /* ---- 5.10.2: pie — pares label/value, sin ejes ---- */
    private static Chart pie(Grafica g) {
        PieChart chart = new PieChart();
        chart.setTitle((String) g.atributos.get("titulo"));
        ArrayList<?> labels = (ArrayList<?>) g.atributos.get("label");
        ArrayList<?> values = (ArrayList<?>) g.atributos.get("values");
        for (int i = 0; i < labels.size(); i++) {
            chart.getData().add(new PieChart.Data(
                    (String) labels.get(i), (Double) values.get(i)));
        }
        return chart;
    }

    /* ---- 5.10.3: histograma — barras de valor vs frecuencia ---- */
    private static Chart histograma(Grafica g) {
        ArrayList<Double> datos = new ArrayList<>();
        for (Object v : (ArrayList<?>) g.atributos.get("values")) datos.add((Double) v);

        CategoryAxis x = new CategoryAxis();
        x.setLabel("Valor");
        NumberAxis y = new NumberAxis();
        y.setLabel("Frecuencia");

        BarChart<String, Number> chart = new BarChart<>(x, y);
        chart.setTitle((String) g.atributos.get("titulo"));
        chart.setLegendVisible(false);
        chart.setCategoryGap(2);

        XYChart.Series<String, Number> s = new XYChart.Series<>();
        for (var e : Operaciones.frecuencias(datos).entrySet()) {
            s.getData().add(new XYChart.Data<>(Entorno.formatear(e.getKey()), e.getValue()));
        }
        chart.getData().add(s);
        return chart;
    }

    /** Empareja ejeX[i] con ejeY[i] en una serie de puntos. */
    private static XYChart.Series<String, Number> serie(Object ejeX, Object ejeY) {
        ArrayList<?> xs = (ArrayList<?>) ejeX;
        ArrayList<?> ys = (ArrayList<?>) ejeY;
        XYChart.Series<String, Number> s = new XYChart.Series<>();
        for (int i = 0; i < xs.size(); i++) {
            s.getData().add(new XYChart.Data<>((String) xs.get(i), (Double) ys.get(i)));
        }
        return s;
    }
}
