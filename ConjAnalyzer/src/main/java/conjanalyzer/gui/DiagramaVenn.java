package conjanalyzer.gui;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import javafx.geometry.Insets;
import javafx.geometry.Pos;
import javafx.scene.canvas.Canvas;
import javafx.scene.canvas.GraphicsContext;
import javafx.scene.control.Label;
import javafx.scene.image.PixelWriter;
import javafx.scene.layout.VBox;
import javafx.scene.paint.Color;
import javafx.scene.text.Font;
import javafx.scene.text.FontWeight;
import javafx.scene.text.TextAlignment;

import conjanalyzer.interprete.Entorno;
import conjanalyzer.interprete.Operacion;

/**
 * Diagrama de Venn de una operacion (5.1). Dibuja un circulo por cada
 * conjunto base referenciado (1, 2 o 3) y SOMBREA la region que abarca el
 * conjunto resultante.
 *
 * Como funciona el sombreado (exacto, no aproximado en la logica):
 * la pertenencia de un elemento al resultado es una funcion booleana de a
 * que conjuntos base pertenece (union/interseccion/diferencia/complemento
 * operan punto a punto sobre las indicadoras). Por cada pixel se calcula en
 * que circulos cae y se evalua el arbol como funcion booleana
 * (NodoOperacion.pertenenciaRegion): si da true, el pixel pertenece al
 * resultado y se pinta. El complemento pinta tambien la region "fuera de
 * todos los circulos" (con un tono mas claro) porque el resultado incluye
 * elementos del universo que no estan en ningun conjunto base.
 *
 * Para operaciones con mas de 3 conjuntos base no hay un Venn geometrico
 * razonable: se muestra el resultado en texto.
 */
public class DiagramaVenn {

    private static final int ANCHO = 470, ALTO = 360;
    private static final Color BASE_BG = Color.web("#f4f6fb");
    private static final Color DENTRO  = Color.web("#2ea043", 0.55);  // resultado dentro de circulos
    private static final Color FUERA   = Color.web("#2ea043", 0.16);  // resultado fuera (complemento)

    private static final Color[] TRAZOS = {
            Color.web("#3b4b8a"), Color.web("#c2410c"), Color.web("#0f766e")
    };

    /** Construye el panel (titulo + canvas + info) para una operacion. */
    public static VBox crear(Operacion op) {
        VBox caja = new VBox(8);
        caja.setPadding(new Insets(10));
        caja.setAlignment(Pos.TOP_CENTER);

        Label titulo = new Label("Operacion: " + op.nombre);
        titulo.setFont(Font.font("Segoe UI", FontWeight.BOLD, 16));

        List<String> bases = new ArrayList<>(op.referencias);
        Canvas canvas = new Canvas(ANCHO, ALTO);
        GraphicsContext gc = canvas.getGraphicsContext2D();

        if (bases.size() >= 1 && bases.size() <= 3) {
            dibujar(gc, op, bases);
        } else {
            gc.setFill(BASE_BG);
            gc.fillRect(0, 0, ANCHO, ALTO);
            gc.setFill(Color.web("#1b2130"));
            gc.setTextAlign(TextAlignment.CENTER);
            gc.setFont(Font.font("Segoe UI", 14));
            gc.fillText(bases.isEmpty()
                    ? "La operacion no referencia conjuntos"
                    : "Venn no disponible para " + bases.size() + " conjuntos base",
                    ANCHO / 2.0, ALTO / 2.0);
        }

        Label prefijo = new Label("Expresion:  " + op.arbol.toPrefijo());
        prefijo.setFont(Font.font("Consolas", 13));

        String resTxt = Entorno.formatearConjunto(op.resultado);
        if (resTxt.length() > 70) resTxt = resTxt.substring(0, 70) + " ...}";
        Label resultado = new Label("Resultado (" + op.resultado.size() + " elementos):  " + resTxt);
        resultado.setFont(Font.font("Consolas", 12));
        resultado.setWrapText(true);
        resultado.setMaxWidth(ANCHO);

        String simplTxt = op.simplificacion.seSimplifico
                ? op.simplificacion.simplificado.toPrefijo()
                  + "   [leyes: " + String.join(", ", op.simplificacion.leyes) + "]"
                : "No se puede simplificar la operacion";
        Label simpl = new Label("Simplificacion:  " + simplTxt);
        simpl.setFont(Font.font("Consolas", 12));
        simpl.setWrapText(true);
        simpl.setMaxWidth(ANCHO);

        caja.getChildren().addAll(titulo, canvas, prefijo, resultado, simpl);
        return caja;
    }

    /* ---------- pintado por pixel + trazos de los circulos ---------- */

    private static void dibujar(GraphicsContext gc, Operacion op, List<String> bases) {
        int n = bases.size();
        double[][] centros = centros(n);
        double r = (n == 3) ? 100 : 110;

        gc.setFill(BASE_BG);
        gc.fillRect(0, 0, ANCHO, ALTO);

        PixelWriter pw = gc.getPixelWriter();
        Map<String, Boolean> region = new HashMap<>();
        for (int y = 0; y < ALTO; y++) {
            for (int x = 0; x < ANCHO; x++) {
                boolean dentroDeAlguno = false;
                for (int i = 0; i < n; i++) {
                    double dx = x - centros[i][0], dy = y - centros[i][1];
                    boolean adentro = dx * dx + dy * dy <= r * r;
                    region.put(bases.get(i), adentro);
                    dentroDeAlguno |= adentro;
                }
                if (op.arbol.pertenenciaRegion(region)) {
                    pw.setColor(x, y, dentroDeAlguno ? DENTRO : FUERA);
                }
            }
        }

        // trazos y etiquetas de cada circulo encima del sombreado
        gc.setLineWidth(2.5);
        gc.setFont(Font.font("Segoe UI", FontWeight.BOLD, 15));
        gc.setTextAlign(TextAlignment.CENTER);
        for (int i = 0; i < n; i++) {
            gc.setStroke(TRAZOS[i]);
            gc.strokeOval(centros[i][0] - r, centros[i][1] - r, 2 * r, 2 * r);
            double[] et = etiqueta(n, i, centros[i], r);
            gc.setFill(TRAZOS[i]);
            gc.fillText(bases.get(i), et[0], et[1]);
        }
    }

    private static double[][] centros(int n) {
        double cx = ANCHO / 2.0;
        return switch (n) {
            case 1 -> new double[][]{ {cx, 175} };
            case 2 -> new double[][]{ {cx - 55, 175}, {cx + 55, 175} };
            default -> new double[][]{ {cx - 60, 145}, {cx + 60, 145}, {cx, 245} };
        };
    }

    private static double[] etiqueta(int n, int i, double[] c, double r) {
        if (n == 1) return new double[]{ c[0], c[1] - r - 8 };
        if (n == 2) return new double[]{ c[0] + (i == 0 ? -r * 0.6 : r * 0.6), c[1] - r - 8 };
        // n == 3
        return switch (i) {
            case 0 -> new double[]{ c[0] - r * 0.5, c[1] - r - 8 };
            case 1 -> new double[]{ c[0] + r * 0.5, c[1] - r - 8 };
            default -> new double[]{ c[0], c[1] + r + 20 };
        };
    }
}
