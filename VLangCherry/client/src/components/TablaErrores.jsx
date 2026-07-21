/* 6.1 Tabla de errores: #, Tipo, Descripción, Línea, Columna. */
export default function TablaErrores({ errores, onIrALinea }) {
  if (!errores || errores.length === 0) {
    return <div className="panel-vacio">Sin errores léxicos, sintácticos ni semánticos.</div>;
  }
  return (
    <table className="tabla-reporte">
      <thead>
        <tr>
          <th>#</th>
          <th>Tipo</th>
          <th>Descripción</th>
          <th>Línea</th>
          <th>Columna</th>
        </tr>
      </thead>
      <tbody>
        {errores.map((e, i) => (
          <tr
            key={i}
            className={'fila-error fila-error-' + e.tipo.toLowerCase()}
            onClick={() => onIrALinea && onIrALinea(e.linea)}
          >
            <td>{i + 1}</td>
            <td>{e.tipo}</td>
            <td>{e.descripcion}</td>
            <td>{e.linea}</td>
            <td>{e.columna}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
