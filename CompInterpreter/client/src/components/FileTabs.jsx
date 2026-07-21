/* 4.1/4.2: el editor permite tener varios archivos abiertos a la vez. */
export default function FileTabs({ archivos, activoId, onSeleccionar, onCerrar }) {
  return (
    <div className="file-tabs">
      {archivos.map((a) => (
        <div
          key={a.id}
          className={'file-tab' + (a.id === activoId ? ' es-activa' : '')}
          onClick={() => onSeleccionar(a.id)}
        >
          <span className="file-tab-nombre">
            {a.nombre}
            {a.sinGuardar ? ' •' : ''}
          </span>
          {archivos.length > 1 && (
            <button
              type="button"
              className="file-tab-cerrar"
              title="Cerrar archivo"
              onClick={(e) => {
                e.stopPropagation();
                onCerrar(a.id);
              }}
            >
              ×
            </button>
          )}
        </div>
      ))}
    </div>
  );
}
