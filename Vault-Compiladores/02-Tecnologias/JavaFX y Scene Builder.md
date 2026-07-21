---
tags: [tecnologia, gui, java]
fuente: "https://openjfx.io/ · https://gluonhq.com/products/scene-builder/"
fecha: 2026-07-10
---

# JavaFX y Scene Builder

**JavaFX:** framework de GUI de escritorio para Java. La app hereda de `Application`; se arranca con `launch()` y `start(Stage)` recibe la ventana (`Stage`), que contiene una `Scene` con un árbol de `Node`s.

Piezas útiles: `BorderPane` (menú/editor/consola), `TabPane` (múltiples archivos), `TextArea` (editor), `TableView` (reportes), `MenuBar`. **JavaFX Charts** incluye gráficas (barras/pie/línea) sin dependencias extra → usado en [[DataForge]].

**Scene Builder:** editor visual *drag-and-drop* que genera **FXML** (XML de la UI) enlazado a un controlador Java con `@FXML`.

**Ejecución:** con [[Maven]] `mvn javafx:run` (módulos `javafx.controls`, `javafx.fxml`).

## Usado en
[[DataForge]], [[ConjAnalyzer]], [[CompScript]]

## Relacionadas
- [[Maven]]
- [[Graphviz]]
- [[DataForge]]
