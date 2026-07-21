package conjanalyzer.gui;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.ArrayList;
import java.util.List;

import javafx.application.Application;
import javafx.geometry.Insets;
import javafx.geometry.Orientation;
import javafx.scene.Scene;
import javafx.scene.control.Button;
import javafx.scene.control.ComboBox;
import javafx.scene.control.Label;
import javafx.scene.control.ScrollPane;
import javafx.scene.control.Separator;
import javafx.scene.control.SplitPane;
import javafx.scene.control.Tab;
import javafx.scene.control.TabPane;
import javafx.scene.control.TextArea;
import javafx.scene.layout.BorderPane;
import javafx.scene.layout.HBox;
import javafx.scene.layout.Priority;
import javafx.scene.layout.StackPane;
import javafx.scene.layout.VBox;
import javafx.scene.text.Font;
import javafx.stage.FileChooser;
import javafx.stage.Stage;

import conjanalyzer.interprete.Entorno;
import conjanalyzer.interprete.Interprete;
import conjanalyzer.interprete.Operacion;
import conjanalyzer.reportes.JsonSalida;
import conjanalyzer.reportes.Reportes;

/**
 * Entorno de trabajo del enunciado (seccion 3): editor con pestañas,
 * Nuevo / Abrir / Guardar, boton Ejecutar, consola de salida no editable
 * (3.5) y panel de diagramas de Venn navegable (5.1).
 *
 * La UI se arma POR CODIGO (sin FXML): el scene graph se construye a mano.
 */
public class EditorApp extends Application {

    private TabPane pestanas;
    private TextArea consola;
    private Entorno ultimoEntorno;
    private int contadorNuevos = 1;

    // panel de Venn navegable
    private final List<Operacion> operaciones = new ArrayList<>();
    private int indiceVenn = 0;
    private StackPane lienzoVenn;
    private ComboBox<String> selectorVenn;
    private Label contadorVenn;
    private Button btnAnterior, btnSiguiente;

    private static final String PLANTILLA =
            "{\n" +
            "    # Definicion de conjuntos\n" +
            "    CONJ : conjuntoA -> 1,2,3,a,b;\n" +
            "    CONJ : conjuntoB -> a~z;\n" +
            "    CONJ : conjuntoC -> 0~9;\n\n" +
            "    # Definicion de operaciones\n" +
            "    OPERA : operacion1 -> & {conjuntoA} {conjuntoB};\n" +
            "    OPERA : operacion2 -> & U {conjuntoB} {conjuntoC} {conjuntoA};\n\n" +
            "    # Evaluar conjuntos de datos\n" +
            "    EVALUAR ( {a, b, c} , operacion1 );\n" +
            "    EVALUAR ( {1, b} , operacion1 );\n" +
            "}\n";

    @Override
    public void start(Stage stage) {
        consola = new TextArea();
        consola.setEditable(false);
        consola.setFont(Font.font("Consolas", 13));
        consola.setPromptText("Consola de salida — ejecuta un programa para ver resultados");

        pestanas = new TabPane();

        Button bNuevo = new Button("Nuevo");
        Button bAbrir = new Button("Abrir");
        Button bGuardar = new Button("Guardar");
        Button bEjecutar = new Button("▶ Ejecutar");
        Button bReportes = new Button("Reportes");
        bNuevo.setOnAction(e -> nuevaPestana("nuevo" + (contadorNuevos++) + ".ca", "{\n\n}\n", null));
        bAbrir.setOnAction(e -> abrir(stage));
        bGuardar.setOnAction(e -> guardar(stage));
        bEjecutar.setOnAction(e -> ejecutar());
        bReportes.setOnAction(e -> reportes());

        HBox barra = new HBox(8, bNuevo, bAbrir, bGuardar, new Separator(), bEjecutar, bReportes);
        barra.setPadding(new Insets(8));

        // izquierda: editor arriba, consola abajo
        SplitPane izquierda = new SplitPane(pestanas, consola);
        izquierda.setOrientation(Orientation.VERTICAL);
        izquierda.setDividerPositions(0.62);

        // derecha: panel de diagramas de Venn navegable
        BorderPane panelVenn = construirPanelVenn();

        SplitPane centro = new SplitPane(izquierda, panelVenn);
        centro.setOrientation(Orientation.HORIZONTAL);
        centro.setDividerPositions(0.52);

        BorderPane raiz = new BorderPane(centro);
        raiz.setTop(barra);

        nuevaPestana("ejemplo.ca", PLANTILLA, null);

        stage.setTitle("ConjAnalyzer — OLC1 Proyecto 1");
        stage.setScene(new Scene(raiz, 1180, 760));
        stage.show();
    }

    /* ============ panel de Venn ============ */

    private BorderPane construirPanelVenn() {
        Label titulo = new Label("Diagramas de Venn");
        titulo.setFont(Font.font("Segoe UI", 14));

        selectorVenn = new ComboBox<>();
        selectorVenn.setOnAction(e -> {
            int i = selectorVenn.getSelectionModel().getSelectedIndex();
            if (i >= 0 && i != indiceVenn) { indiceVenn = i; mostrarVenn(); }
        });

        btnAnterior = new Button("◀");
        btnSiguiente = new Button("▶");
        contadorVenn = new Label("0 / 0");
        btnAnterior.setOnAction(e -> { if (indiceVenn > 0) { indiceVenn--; mostrarVenn(); } });
        btnSiguiente.setOnAction(e -> {
            if (indiceVenn < operaciones.size() - 1) { indiceVenn++; mostrarVenn(); }
        });

        HBox nav = new HBox(8, btnAnterior, contadorVenn, btnSiguiente, new Separator(), selectorVenn);
        nav.setPadding(new Insets(8));

        lienzoVenn = new StackPane();
        lienzoVenn.getChildren().add(new Label("Ejecuta un programa con operaciones para ver sus diagramas."));

        VBox arriba = new VBox(4, titulo, nav);
        arriba.setPadding(new Insets(8, 8, 0, 8));

        ScrollPane scroll = new ScrollPane(lienzoVenn);
        scroll.setFitToWidth(true);
        VBox.setVgrow(scroll, Priority.ALWAYS);

        BorderPane panel = new BorderPane(scroll);
        panel.setTop(arriba);
        return panel;
    }

    private void refrescarPanelVenn() {
        operaciones.clear();
        if (ultimoEntorno != null) operaciones.addAll(ultimoEntorno.getOperaciones().values());
        selectorVenn.getItems().clear();
        for (Operacion op : operaciones) selectorVenn.getItems().add(op.nombre);
        indiceVenn = 0;
        mostrarVenn();
    }

    private void mostrarVenn() {
        lienzoVenn.getChildren().clear();
        if (operaciones.isEmpty()) {
            lienzoVenn.getChildren().add(new Label("No hay operaciones en el ultimo analisis."));
            contadorVenn.setText("0 / 0");
            btnAnterior.setDisable(true);
            btnSiguiente.setDisable(true);
            return;
        }
        Operacion op = operaciones.get(indiceVenn);
        lienzoVenn.getChildren().add(DiagramaVenn.crear(op));
        contadorVenn.setText((indiceVenn + 1) + " / " + operaciones.size());
        selectorVenn.getSelectionModel().select(indiceVenn);
        btnAnterior.setDisable(indiceVenn == 0);
        btnSiguiente.setDisable(indiceVenn == operaciones.size() - 1);
    }

    /* ============ pestañas ============ */

    private void nuevaPestana(String titulo, String contenido, File archivo) {
        TextArea area = new TextArea(contenido);
        area.setFont(Font.font("Consolas", 14));
        Tab tab = new Tab(titulo, area);
        tab.setUserData(archivo);
        pestanas.getTabs().add(tab);
        pestanas.getSelectionModel().select(tab);
    }

    private TextArea areaActual() {
        Tab t = pestanas.getSelectionModel().getSelectedItem();
        return (t == null) ? null : (TextArea) t.getContent();
    }

    /* ============ abrir / guardar (3.2) ============ */

    private void abrir(Stage stage) {
        FileChooser fc = new FileChooser();
        fc.setTitle("Abrir archivo ConjAnalyzer");
        fc.getExtensionFilters().add(new FileChooser.ExtensionFilter("ConjAnalyzer (*.ca)", "*.ca"));
        File f = fc.showOpenDialog(stage);
        if (f == null) return;
        try {
            nuevaPestana(f.getName(), Files.readString(f.toPath(), StandardCharsets.UTF_8), f);
        } catch (Exception ex) {
            consola.setText("No se pudo abrir el archivo: " + ex.getMessage());
        }
    }

    private void guardar(Stage stage) {
        Tab tab = pestanas.getSelectionModel().getSelectedItem();
        if (tab == null) return;
        File f = (File) tab.getUserData();
        if (f == null) {
            FileChooser fc = new FileChooser();
            fc.setTitle("Guardar como");
            fc.getExtensionFilters().add(new FileChooser.ExtensionFilter("ConjAnalyzer (*.ca)", "*.ca"));
            fc.setInitialFileName(tab.getText());
            f = fc.showSaveDialog(stage);
            if (f == null) return;
            tab.setUserData(f);
            tab.setText(f.getName());
        }
        try {
            Files.writeString(f.toPath(), ((TextArea) tab.getContent()).getText(), StandardCharsets.UTF_8);
        } catch (Exception ex) {
            consola.setText("No se pudo guardar el archivo: " + ex.getMessage());
        }
    }

    /* ============ ejecutar (3.3): la pestaña actual → el interprete ============ */

    private void ejecutar() {
        TextArea area = areaActual();
        if (area == null) {
            consola.setText("No hay ninguna pestaña abierta.");
            return;
        }
        // Entorno FRESCO por ejecucion (seccion 5)
        ultimoEntorno = Interprete.ejecutar(area.getText());

        StringBuilder salida = new StringBuilder(ultimoEntorno.getConsola());
        if (!ultimoEntorno.getErrores().isEmpty()) {
            salida.append("\n--- ERRORES (").append(ultimoEntorno.getErrores().size()).append(") ---\n");
            for (var e : ultimoEntorno.getErrores()) salida.append(e).append('\n');
        }
        if (!ultimoEntorno.getOperaciones().isEmpty()) {
            salida.append("\n--- JSON DE SIMPLIFICACION (5.4) ---\n")
                  .append(JsonSalida.construir(ultimoEntorno)).append('\n');
        }
        consola.setText(salida.toString());
        refrescarPanelVenn();
    }

    /* ============ reportes (3.4 y 5): del ultimo analisis ============ */

    private void reportes() {
        if (ultimoEntorno == null) {
            consola.setText("Ejecuta un programa primero — los reportes son del ultimo analisis.");
            return;
        }
        try {
            File carpeta = new File("reportes");
            File[] htmls = Reportes.generar(ultimoEntorno, carpeta);
            File json = JsonSalida.generar(ultimoEntorno, carpeta);
            for (File f : htmls) getHostServices().showDocument(f.toURI().toString());
            consola.appendText("\n(Reportes generados en " + carpeta.getAbsolutePath()
                    + " — incluye " + json.getName() + ")\n");
        } catch (Exception ex) {
            consola.setText("No se pudieron generar los reportes: " + ex.getMessage());
        }
    }

    public static void main(String[] args) {
        launch(args);
    }
}
