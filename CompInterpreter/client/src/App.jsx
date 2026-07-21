import { useMemo, useRef, useState } from 'react';
import Toolbar from './components/Toolbar';
import FileTabs from './components/FileTabs';
import Editor from './components/Editor';
import Consola from './components/Consola';
import TablaErrores from './components/TablaErrores';
import TablaSimbolos from './components/TablaSimbolos';
import AstGrafo from './components/AstGrafo';
import { interpretar } from './api';
import EJEMPLO_ANEXO from './ejemplo-anexo.js';
import './App.css';

let siguienteId = 1;
function nuevoId() {
  return siguienteId++;
}

function archivoNuevo(nombre, contenido) {
  return { id: nuevoId(), nombre, contenido, sinGuardar: false };
}

const PANELES = [
  { clave: 'consola', etiqueta: 'Consola' },
  { clave: 'errores', etiqueta: 'Errores' },
  { clave: 'simbolos', etiqueta: 'Símbolos' },
  { clave: 'ast', etiqueta: 'AST' },
];

export default function App() {
  const [archivos, setArchivos] = useState(() => [
    archivoNuevo('principal.ci', EJEMPLO_ANEXO),
  ]);
  const [activoId, setActivoId] = useState(() => archivos[0].id);
  const [resultado, setResultado] = useState(null);
  const [ejecutando, setEjecutando] = useState(false);
  const [errorRed, setErrorRed] = useState(null);
  const [panel, setPanel] = useState('consola');
  const editorRef = useRef(null);

  const archivoActivo = archivos.find((a) => a.id === activoId) || archivos[0];

  const erroresPorLinea = useMemo(() => {
    const set = new Set();
    if (resultado && resultado.errores) {
      resultado.errores.forEach((e) => set.add(e.linea));
    }
    return set;
  }, [resultado]);

  const actualizarContenido = (contenido) => {
    setArchivos((prev) =>
      prev.map((a) => (a.id === activoId ? { ...a, contenido, sinGuardar: true } : a))
    );
  };

  const manejarNuevo = () => {
    const a = archivoNuevo('sin-titulo-' + (archivos.length + 1) + '.ci', '');
    setArchivos((prev) => [...prev, a]);
    setActivoId(a.id);
  };

  const manejarAbrir = (file) => {
    const lector = new FileReader();
    lector.onload = () => {
      const a = archivoNuevo(file.name, String(lector.result));
      setArchivos((prev) => [...prev, a]);
      setActivoId(a.id);
    };
    lector.readAsText(file);
  };

  const manejarGuardar = () => {
    if (!archivoActivo) return;
    const blob = new Blob([archivoActivo.contenido], { type: 'text/plain;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const enlace = document.createElement('a');
    enlace.href = url;
    enlace.download = archivoActivo.nombre.endsWith('.ci')
      ? archivoActivo.nombre
      : archivoActivo.nombre + '.ci';
    enlace.click();
    URL.revokeObjectURL(url);
    setArchivos((prev) =>
      prev.map((a) => (a.id === activoId ? { ...a, sinGuardar: false } : a))
    );
  };

  const manejarCerrar = (id) => {
    setArchivos((prev) => {
      const restantes = prev.filter((a) => a.id !== id);
      if (id === activoId && restantes.length) setActivoId(restantes[0].id);
      return restantes;
    });
  };

  const manejarEjecutar = async () => {
    if (!archivoActivo) return;
    setEjecutando(true);
    setErrorRed(null);
    try {
      const res = await interpretar(archivoActivo.contenido);
      setResultado(res);
      setPanel(res.errores && res.errores.length ? 'errores' : 'consola');
    } catch (e) {
      setErrorRed(e.message || String(e));
    } finally {
      setEjecutando(false);
    }
  };

  const irALinea = (num) => {
    if (editorRef.current) editorRef.current.irALinea(num);
  };

  return (
    <div className="app">
      <Toolbar
        onNuevo={manejarNuevo}
        onAbrir={manejarAbrir}
        onGuardar={manejarGuardar}
        onEjecutar={manejarEjecutar}
        ejecutando={ejecutando}
      />
      <FileTabs
        archivos={archivos}
        activoId={activoId}
        onSeleccionar={setActivoId}
        onCerrar={manejarCerrar}
      />
      <div className="app-cuerpo">
        <div className="panel-editor">
          {archivoActivo && (
            <Editor
              ref={editorRef}
              value={archivoActivo.contenido}
              onChange={actualizarContenido}
              erroresPorLinea={erroresPorLinea}
            />
          )}
        </div>
        <div className="panel-resultados">
          <div className="pestanas-resultado">
            {PANELES.map((p) => (
              <button
                key={p.clave}
                type="button"
                className={'pestana' + (panel === p.clave ? ' es-activa' : '')}
                onClick={() => setPanel(p.clave)}
              >
                {p.etiqueta}
                {p.clave === 'errores' && resultado && resultado.errores.length > 0
                  ? ' (' + resultado.errores.length + ')'
                  : ''}
              </button>
            ))}
          </div>
          <div className="contenido-resultado">
            {errorRed && <div className="aviso-error-red">{errorRed}</div>}
            {!errorRed && panel === 'consola' && (
              <Consola lineas={resultado ? resultado.consolaLineas : null} />
            )}
            {!errorRed && panel === 'errores' && (
              <TablaErrores errores={resultado ? resultado.errores : null} onIrALinea={irALinea} />
            )}
            {!errorRed && panel === 'simbolos' && (
              <TablaSimbolos simbolos={resultado ? resultado.simbolos : null} onIrALinea={irALinea} />
            )}
            {!errorRed && panel === 'ast' && <AstGrafo ast={resultado ? resultado.ast : null} />}
          </div>
        </div>
      </div>
    </div>
  );
}
