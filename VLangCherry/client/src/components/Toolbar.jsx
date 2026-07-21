import { useRef } from 'react';

/* 4.2 Funcionalidades (nuevo/abrir/guardar) + 4.4 Herramientas (ejecutar). */
export default function Toolbar({ onNuevo, onAbrir, onGuardar, onEjecutar, ejecutando }) {
  const inputRef = useRef(null);

  const manejarArchivoElegido = (e) => {
    const file = e.target.files && e.target.files[0];
    if (file) onAbrir(file);
    e.target.value = '';
  };

  return (
    <div className="toolbar">
      <span className="toolbar-titulo">VLangCherry</span>
      <div className="toolbar-acciones">
        <button type="button" onClick={onNuevo} title="Nuevo archivo">
          Nuevo
        </button>
        <button type="button" onClick={() => inputRef.current.click()} title="Abrir archivo .vch">
          Abrir
        </button>
        <input
          ref={inputRef}
          type="file"
          accept=".vch,text/plain"
          hidden
          onChange={manejarArchivoElegido}
        />
        <button type="button" onClick={onGuardar} title="Guardar archivo actual">
          Guardar
        </button>
        <button
          type="button"
          className="boton-ejecutar"
          onClick={onEjecutar}
          disabled={ejecutando}
        >
          {ejecutando ? 'Ejecutando…' : '▶ Ejecutar'}
        </button>
      </div>
    </div>
  );
}
