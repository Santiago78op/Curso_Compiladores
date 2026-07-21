/* 4.6 / 6.4: área de consola, solo lectura. */
export default function Consola({ lineas }) {
  if (!lineas || lineas.length === 0) {
    return <div className="panel-vacio">Ejecutá un archivo para ver la salida acá.</div>;
  }
  return (
    <pre className="consola" tabIndex={-1}>
      {lineas.join('\n')}
    </pre>
  );
}
