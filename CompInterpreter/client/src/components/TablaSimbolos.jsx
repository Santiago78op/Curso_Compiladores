/* 6.2 Tabla de símbolos: #, ID, Tipo, Tipo de Dato, Entorno, Valor, Línea, Columna. */
export default function TablaSimbolos({ simbolos, onIrALinea }) {
  if (!simbolos || simbolos.length === 0) {
    return <div className="panel-vacio">Sin símbolos declarados todavía.</div>;
  }
  return (
    <table className="tabla-reporte">
      <thead>
        <tr>
          <th>#</th>
          <th>ID</th>
          <th>Tipo</th>
          <th>Tipo de Dato</th>
          <th>Entorno</th>
          <th>Valor</th>
          <th>Línea</th>
          <th>Columna</th>
        </tr>
      </thead>
      <tbody>
        {simbolos.map((s, i) => (
          <tr key={i} onClick={() => onIrALinea && onIrALinea(s.linea)}>
            <td>{i + 1}</td>
            <td>{s.id}</td>
            <td>{s.categoria}</td>
            <td>{s.tipoDato}</td>
            <td>{s.entorno}</td>
            <td>{s.valor}</td>
            <td>{s.linea}</td>
            <td>{s.columna}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
