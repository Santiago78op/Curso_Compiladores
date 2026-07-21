package dataforge.gui;

import java.io.File;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;

import javafx.application.Application;
import javafx.geometry.Insets;
import javafx.geometry.Orientation;
import javafx.scene.Scene;
import javafx.scene.control.Button;
import javafx.scene.control.Separator;
import javafx.scene.control.SplitPane;
import javafx.scene.control.Tab;
import javafx.scene.control.TabPane;
import javafx.scene.control.TextArea;
import javafx.scene.layout.BorderPane;
import javafx.scene.layout.HBox;
import javafx.scene.text.Font;
import javafx.stage.FileChooser;
import javafx.stage.Stage;

import dataforge.interprete.Entorno;
import dataforge.interprete.Interprete;

/**
 * El entorno de trabajo del enunciado (§4): editor con pestañas,
 * Nuevo/Abrir/Guardar, botón Ejecutar y consola de salida no editable.
 *
 * La UI está construida POR CÓDIGO (sin FXML): el scene graph se arma
 * a mano, lo que hace visible la estructura Stage → Scene → nodos.
 */
public class EditorApp extends Application {

    private TabPane pestanas;
    private TextArea consola;
    private Entorno ultimoEntorno;   // lo consumirán las Etapas 5 (gráficas) y 6 (reportes)
    private int contadorNuevos = 1;

    @Override
    public void start(Stage stage) {
        /* ---- consola (enunciado 4.6: solo muestra, no editable) ---- */
        consola = new TextArea();
        consola.setEditable(false);
        consola.setFont(Font.font("Consolas", 13));
        consola.setPromptText("Consola de salida — ejecutá un programa para ver resultados");

        /* ---- editor con pestañas (4.1-4.3) ---- */
        pestanas = new TabPane();

        /* ---- barra de herramientas (4.2 y 4.4) ---- */
        Button bNuevo = new Button("Nuevo");
        Button bAbrir = new Button("Abrir");
        Button bGuardar = new Button("Guardar");
        Button bEjecutar = new Button("▶ Ejecutar");
        Button bReportes = new Button("Reportes");
        bNuevo.setOnAction(e -> nuevaPestana("nuevo" + (contadorNuevos++) + ".df",
                "PROGRAM\n\n\nEND PROGRAM\n", null));
        bAbrir.setOnAction(e -> abrir(stage));
        bGuardar.setOnAction(e -> guardar(stage));
        bEjecutar.setOnAction(e -> ejecutar());
        bReportes.setOnAction(e -> reportes());

        HBox barra = new HBox(8, bNuevo, bAbrir, bGuardar, new Separator(), bEjecutar, bReportes);
        barra.setPadding(new Insets(8));

        /* ---- composición: editor arriba, consola abajo ---- */
        SplitPane centro = new SplitPane(pestanas, consola);
        centro.setOrientation(Orientation.VERTICAL);
        centro.setDividerPositions(0.68);

        BorderPane raiz = new BorderPane(centro);
        raiz.setTop(barra);

        nuevaPestana("nuevo" + (contadorNuevos++) + ".df",
                "PROGRAM\n\n\nEND PROGRAM\n", null);

        stage.setTitle("DataForge — OLC1 Proyecto 1");
        stage.setScene(new Scene(raiz, 1000, 700));
        stage.show();
    }

    /* ============ pestañas ============ */

    /** Crea una pestaña con su editor; el File asociado viaja en userData
     *  (null = archivo nuevo sin guardar). Cerrable en cualquier momento:
     *  si no se guardó, se descarta (enunciado 4.2). */
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

    /* ============ abrir / guardar (4.2) ============ */

    private void abrir(Stage stage) {
        FileChooser fc = new FileChooser();
        fc.setTitle("Abrir archivo DataForge");
        fc.getExtensionFilters().add(
                new FileChooser.ExtensionFilter("DataForge (*.df)", "*.df"));
        File f = fc.showOpenDialog(stage);
        if (f == null) return;
        try {
            nuevaPestana(f.getName(),
                    Files.readString(f.toPath(), StandardCharsets.UTF_8), f);
        } catch (IOException ex) {
            consola.setText("No se pudo abrir el archivo: " + ex.getMessage());
        }
    }

    private void guardar(Stage stage) {
        Tab tab = pestanas.getSelectionModel().getSelectedItem();
        if (tab == null) return;
        File f = (File) tab.getUserData();
        if (f == null) {                       // archivo nuevo → "Guardar como"
            FileChooser fc = new FileChooser();
            fc.setTitle("Guardar como");
            fc.getExtensionFilters().add(
                    new FileChooser.ExtensionFilter("DataForge (*.df)", "*.df"));
            fc.setInitialFileName(tab.getText());
            f = fc.showSaveDialog(stage);
            if (f == null) return;
            tab.setUserData(f);
            tab.setText(f.getName());
        }
        try {
            Files.writeString(f.toPath(),
                    ((TextArea) tab.getContent()).getText(), StandardCharsets.UTF_8);
        } catch (IOException ex) {
            consola.setText("No se pudo guardar el archivo: " + ex.getMessage());
        }
    }

    /* ============ ejecutar (4.4): la pestaña actual → el intérprete ============ */

    private void ejecutar() {
        TextArea area = areaActual();
        if (area == null) {
            consola.setText("No hay ninguna pestaña abierta.");
            return;
        }
        /* Entorno FRESCO por ejecución: el estado de una corrida
           no se arrastra a la siguiente (enunciado §6: los reportes
           son solo del último análisis) */
        ultimoEntorno = Interprete.ejecutar(area.getText());

        StringBuilder salida = new StringBuilder(ultimoEntorno.getConsola());
        if (!ultimoEntorno.getErrores().isEmpty()) {
            salida.append("\n─── ERRORES (")
                  .append(ultimoEntorno.getErrores().size()).append(") ───\n");
            for (var e : ultimoEntorno.getErrores())
                salida.append(e).append('\n');
        }
        if (!ultimoEntorno.getGraficas().isEmpty()) {
            salida.append("\n(").append(ultimoEntorno.getGraficas().size())
                  .append(" gráfica(s) mostrada(s) en pantalla)\n");
            /* estamos en el hilo de JavaFX (evento del botón):
               se pueden abrir ventanas directamente */
            ultimoEntorno.getGraficas().forEach(Graficador::mostrar);
        }
        consola.setText(salida.toString());
    }

    /* ============ reportes (4.5 y §6): del último análisis ============ */

    private void reportes() {
        if (ultimoEntorno == null) {
            consola.setText("Ejecutá un programa primero — los reportes son del último análisis.");
            return;
        }
        try {
            File[] archivos = dataforge.reportes.Reportes.generar(
                    ultimoEntorno, new File("reportes"));
            for (File f : archivos) {
                /* getHostServices: la forma JavaFX de abrir el navegador */
                getHostServices().showDocument(f.toURI().toString());
            }
        } catch (Exception ex) {
            consola.setText("No se pudieron generar los reportes: " + ex.getMessage());
        }
    }

    public static void main(String[] args) {
        launch(args);
    }
}
