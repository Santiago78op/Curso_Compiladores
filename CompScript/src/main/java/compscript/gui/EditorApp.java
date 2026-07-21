package compscript.gui;

import java.io.File;
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
import javafx.scene.control.TreeItem;
import javafx.scene.control.TreeView;
import javafx.scene.layout.BorderPane;
import javafx.scene.layout.HBox;
import javafx.scene.text.Font;
import javafx.stage.FileChooser;
import javafx.stage.Stage;

import compscript.ast.A;
import compscript.interprete.Contexto;
import compscript.interprete.Interprete;

/**
 * Entorno de trabajo del enunciado (4): editor con pestanas, Nuevo/Abrir/
 * Guardar, Ejecutar, y botones de reportes (Tokens/Errores/Simbolos/AST).
 * La UI se construye por codigo (sin FXML).
 */
public class EditorApp extends Application {

    private TabPane pestanas;
    private TextArea consola;
    private Contexto ultimo;   // resultado de la ultima ejecucion (para reportes)
    private int contador = 1;

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
        Button bAst = new Button("Ver AST");
        bNuevo.setOnAction(e -> nuevaPestana("nuevo" + (contador++) + ".cs", plantilla(), null));
        bAbrir.setOnAction(e -> abrir(stage));
        bGuardar.setOnAction(e -> guardar(stage));
        bEjecutar.setOnAction(e -> ejecutar());
        bReportes.setOnAction(e -> reportes());
        bAst.setOnAction(e -> verAst());

        HBox barra = new HBox(8, bNuevo, bAbrir, bGuardar, new Separator(),
                bEjecutar, new Separator(), bReportes, bAst);
        barra.setPadding(new Insets(8));

        SplitPane centro = new SplitPane(pestanas, consola);
        centro.setOrientation(Orientation.VERTICAL);
        centro.setDividerPositions(0.68);

        BorderPane raiz = new BorderPane(centro);
        raiz.setTop(barra);

        nuevaPestana("nuevo" + (contador++) + ".cs", plantilla(), null);

        stage.setTitle("CompScript — OLC1 PT1 VD2024");
        stage.setScene(new Scene(raiz, 1050, 720));
        stage.show();
    }

    private String plantilla() {
        return "// CompScript\nint main() {\n    console.log(\"Hola CompScript\");\n    return 0;\n}\n\nRUN_MAIN main();\n";
    }

    /* ============ pestanas ============ */
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
        fc.setTitle("Abrir archivo CompScript");
        fc.getExtensionFilters().add(new FileChooser.ExtensionFilter("CompScript (*.cs)", "*.cs"));
        File f = fc.showOpenDialog(stage);
        if (f == null) return;
        try {
            nuevaPestana(f.getName(), Files.readString(f.toPath(), StandardCharsets.UTF_8), f);
        } catch (Exception ex) {
            consola.setText("No se pudo abrir: " + ex.getMessage());
        }
    }

    private void guardar(Stage stage) {
        Tab tab = pestanas.getSelectionModel().getSelectedItem();
        if (tab == null) return;
        File f = (File) tab.getUserData();
        if (f == null) {
            FileChooser fc = new FileChooser();
            fc.setTitle("Guardar como");
            fc.getExtensionFilters().add(new FileChooser.ExtensionFilter("CompScript (*.cs)", "*.cs"));
            fc.setInitialFileName(tab.getText());
            f = fc.showSaveDialog(stage);
            if (f == null) return;
            tab.setUserData(f);
            tab.setText(f.getName());
        }
        try {
            Files.writeString(f.toPath(), ((TextArea) tab.getContent()).getText(), StandardCharsets.UTF_8);
        } catch (Exception ex) {
            consola.setText("No se pudo guardar: " + ex.getMessage());
        }
    }

    /* ============ ejecutar (4.4) ============ */
    private void ejecutar() {
        TextArea area = areaActual();
        if (area == null) { consola.setText("No hay ninguna pestana abierta."); return; }
        ultimo = Interprete.ejecutar(area.getText());   // Contexto FRESCO por ejecucion

        StringBuilder salida = new StringBuilder(ultimo.consola);
        if (!ultimo.errores.isEmpty()) {
            salida.append("\n─── ERRORES (").append(ultimo.errores.size()).append(") ───\n");
            for (var e : ultimo.errores) salida.append(e).append('\n');
        }
        consola.setText(salida.toString());
    }

    /* ============ reportes (4.5 y 6) ============ */
    private void reportes() {
        if (ultimo == null) { consola.setText("Ejecuta un programa primero."); return; }
        try {
            File[] archivos = compscript.reportes.Reportes.generar(ultimo, new File("reportes"));
            for (File f : archivos) getHostServices().showDocument(f.toURI().toString());
        } catch (Exception ex) {
            consola.setText("No se pudieron generar los reportes: " + ex.getMessage());
        }
    }

    /* ============ AST en un TreeView, desde la interfaz (6.3) ============ */
    private void verAst() {
        if (ultimo == null || ultimo.raiz == null) { consola.setText("Ejecuta un programa primero."); return; }
        TreeItem<String> raiz = construir(ultimo.raiz);
        raiz.setExpanded(true);
        TreeView<String> arbol = new TreeView<>(raiz);
        Stage v = new Stage();
        v.setTitle("AST — CompScript");
        v.setScene(new Scene(arbol, 520, 640));
        v.show();
    }

    private TreeItem<String> construir(A.Nodo nodo) {
        TreeItem<String> item = new TreeItem<>(nodo.etiquetaAst());
        for (A.Nodo h : nodo.hijosAst()) if (h != null) item.getChildren().add(construir(h));
        item.setExpanded(true);
        return item;
    }

    public static void main(String[] args) { launch(args); }
}
