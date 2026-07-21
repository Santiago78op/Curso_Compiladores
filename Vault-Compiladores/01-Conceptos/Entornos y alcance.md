---
tags: [concepto, semantico]
fuente: "Libro del Dragón, cap. 7.3"
fecha: 2026-07-10
---

# Entornos y alcance

El **alcance (scope)** define dónde es visible un nombre.
- **Alcance estático (léxico):** depende de **dónde está escrito** (bloques anidados). Es el modelo de Java, JS y de los proyectos.
- **Alcance dinámico:** depende de quién llamó en ejecución (poco común).

Se implementa con una **[[Tabla de símbolos|tabla de símbolos]] como pila/árbol de entornos**: buscar un nombre = mirar el entorno actual y subir al **entorno padre** (enlace de acceso). Habilita variables globales vs locales y *shadowing*.

## Relacionadas
- [[Tabla de símbolos]]
- [[Registro de activación y pila de control]]
- [[CompScript]]
