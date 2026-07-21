/* ============================================================
   Reporte AST como grafo (seccion 6.3).
   Recorre el AST generico (objetos con propiedad .tipo) y produce
   nodos/aristas listos para vis-network. Los campos escalares del
   nodo (op, nombre, valor, ...) se pegan a la etiqueta; los campos que
   son nodos hijos o arreglos de nodos se convierten en aristas.
   Al ser generico, cualquier nodo nuevo del AST aparece sin mantenimiento.
   ============================================================ */

function esNodo(x) {
  return x !== null && typeof x === 'object' && typeof x.tipo === 'string';
}

function construirGrafo(raiz) {
  const nodes = [];
  const edges = [];
  let contador = 0;

  // Visita un nodo del AST por reflection sobre sus propiedades (sin saber
  // nada de la gramática): separa cada campo en escalar (se concatena a la
  // etiqueta del nodo), nodo-hijo (arista) o arreglo (nodo-hijo por cada
  // elemento que sea nodo; el resto se junta como lista escalar). Recorre
  // los hijos DESPUES de armar la etiqueta para no pisar el orden de arriba
  // hacia abajo esperado por el layout jerárquico de vis-network.
  function visitar(n) {
    const miId = 'n' + (contador++);
    let etiqueta = n.tipo;
    const hijos = [];   // {label, ref} para procesar despues (mantiene orden)

    for (const clave of Object.keys(n)) {
      if (clave === 'tipo' || clave === 'linea' || clave === 'columna') continue;
      const v = n[clave];
      if (esNodo(v)) {
        hijos.push({ etiqueta: clave, nodo: v });
      } else if (Array.isArray(v)) {
        // arreglo: puede ser de nodos, de strings (ids) o mixto
        const escalares = [];
        v.forEach((el, i) => {
          if (esNodo(el)) hijos.push({ etiqueta: clave + '[' + i + ']', nodo: el });
          else if (Array.isArray(el)) {
            el.forEach((sub, j) => { if (esNodo(sub)) hijos.push({ etiqueta: clave + '[' + i + '][' + j + ']', nodo: sub }); });
          } else escalares.push(el);
        });
        if (escalares.length) etiqueta += '\n' + clave + ': [' + escalares.join(', ') + ']';
      } else if (v !== null && v !== undefined) {
        // escalar -> a la etiqueta
        etiqueta += '\n' + clave + ': ' + String(v);
      }
    }

    nodes.push({ id: miId, label: etiqueta });
    for (const h of hijos) {
      const hijoId = visitar(h.nodo);
      edges.push({ from: miId, to: hijoId, label: h.etiqueta });
    }
    return miId;
  }

  if (raiz) visitar(raiz);
  return { nodes, edges };
}

/* Genera codigo fuente DOT (Graphviz) equivalente, por si se quiere
   graficar del lado servidor o descargar. */
function aDot(grafo) {
  let dot = 'digraph AST {\n  node [shape=box, style=rounded, fontname="Consolas"];\n';
  for (const n of grafo.nodes) {
    const lbl = String(n.label).replace(/"/g, '\\"').replace(/\n/g, '\\n');
    dot += '  ' + n.id + ' [label="' + lbl + '"];\n';
  }
  for (const e of grafo.edges) {
    dot += '  ' + e.from + ' -> ' + e.to + ';\n';
  }
  dot += '}\n';
  return dot;
}

module.exports = { construirGrafo, aDot };
