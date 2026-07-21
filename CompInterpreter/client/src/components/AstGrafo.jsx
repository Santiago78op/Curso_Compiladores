import { useEffect, useRef } from 'react';
import { Network } from 'vis-network';
import { DataSet } from 'vis-data';

/* 6.3 AST: reporte como grafo (vis-network) a partir de { nodes, edges }
   que ya llegan armados desde el servidor (reportes/ast-grafo.js). */
export default function AstGrafo({ ast }) {
  const containerRef = useRef(null);
  const networkRef = useRef(null);

  useEffect(() => {
    if (!containerRef.current || !ast || !ast.nodes || ast.nodes.length === 0) {
      return;
    }

    const estilos = getComputedStyle(document.documentElement);
    const colorTexto = estilos.getPropertyValue('--text-h').trim() || '#08060d';
    const colorFondoNodo = estilos.getPropertyValue('--code-bg').trim() || '#f4f3ec';
    const colorBorde = estilos.getPropertyValue('--accent').trim() || '#aa3bff';
    const colorLinea = estilos.getPropertyValue('--border').trim() || '#e5e4e7';

    const nodes = new DataSet(
      ast.nodes.map((n) => ({ id: n.id, label: n.label }))
    );
    const edges = new DataSet(
      ast.edges.map((e, i) => ({ id: i, from: e.from, to: e.to, label: e.label }))
    );

    const options = {
      layout: {
        hierarchical: {
          enabled: true,
          direction: 'UD',
          sortMethod: 'directed',
          levelSeparation: 90,
          nodeSpacing: 150,
        },
      },
      nodes: {
        shape: 'box',
        margin: { top: 8, bottom: 8, left: 10, right: 10 },
        color: { background: colorFondoNodo, border: colorBorde },
        font: { color: colorTexto, face: 'ui-monospace, Consolas, monospace', size: 13, multi: false, align: 'left' },
      },
      edges: {
        color: { color: colorLinea, highlight: colorBorde },
        arrows: { to: { enabled: true, scaleFactor: 0.6 } },
        font: { color: colorTexto, size: 11, strokeWidth: 0, background: 'none' },
        smooth: { type: 'cubicBezier', forceDirection: 'vertical' },
      },
      physics: false,
      interaction: { hover: true, zoomView: true, dragView: true },
    };

    if (networkRef.current) networkRef.current.destroy();
    networkRef.current = new Network(containerRef.current, { nodes, edges }, options);

    return () => {
      if (networkRef.current) {
        networkRef.current.destroy();
        networkRef.current = null;
      }
    };
  }, [ast]);

  if (!ast || !ast.nodes || ast.nodes.length === 0) {
    return <div className="panel-vacio">Ejecutá un archivo para ver el AST.</div>;
  }
  return <div className="ast-grafo" ref={containerRef} />;
}
