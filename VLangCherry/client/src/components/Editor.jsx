import { forwardRef, useImperativeHandle, useRef, useState, useCallback } from 'react';

const LINE_HEIGHT = 21;

/* Editor de código: textarea monoespaciada + gutter de líneas sincronizado,
   resalte de la línea actual (4.1: "deberá mostrar la línea actual") y
   marcas rojas en el gutter para las líneas con error (6.1). */
const Editor = forwardRef(function Editor({ value, onChange, erroresPorLinea }, ref) {
  const textareaRef = useRef(null);
  const gutterRef = useRef(null);
  const [linea, setLinea] = useState(1);
  const [scrollTop, setScrollTop] = useState(0);

  const lineas = value.length ? value.split('\n') : [''];

  const actualizarLineaActual = useCallback(() => {
    const ta = textareaRef.current;
    if (!ta) return;
    const hastaCursor = ta.value.slice(0, ta.selectionStart);
    setLinea(hastaCursor.split('\n').length);
  }, []);

  const sincronizarScroll = (e) => {
    const top = e.target.scrollTop;
    setScrollTop(top);
    if (gutterRef.current) gutterRef.current.scrollTop = top;
  };

  useImperativeHandle(ref, () => ({
    irALinea(num) {
      const ta = textareaRef.current;
      if (!ta) return;
      const offsets = value.split('\n').reduce(
        (acc, l) => {
          acc.push(acc[acc.length - 1] + l.length + 1);
          return acc;
        },
        [0]
      );
      const pos = offsets[Math.max(0, num - 1)] || 0;
      ta.focus();
      ta.setSelectionRange(pos, pos);
      setLinea(num);
      const destino = (num - 1) * LINE_HEIGHT - ta.clientHeight / 2;
      ta.scrollTop = Math.max(0, destino);
      setScrollTop(ta.scrollTop);
      if (gutterRef.current) gutterRef.current.scrollTop = ta.scrollTop;
    },
  }));

  return (
    <div className="editor">
      <div className="editor-gutter" ref={gutterRef}>
        {lineas.map((_, i) => {
          const num = i + 1;
          const tieneError = erroresPorLinea && erroresPorLinea.has(num);
          return (
            <div
              key={num}
              className={
                'editor-gutter-linea' +
                (num === linea ? ' es-actual' : '') +
                (tieneError ? ' tiene-error' : '')
              }
            >
              {num}
            </div>
          );
        })}
      </div>
      <div className="editor-cuerpo">
        <div
          className="editor-linea-resaltada"
          style={{ top: (linea - 1) * LINE_HEIGHT - scrollTop }}
        />
        <textarea
          ref={textareaRef}
          className="editor-textarea"
          spellCheck={false}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onScroll={sincronizarScroll}
          onKeyUp={actualizarLineaActual}
          onClick={actualizarLineaActual}
          onFocus={actualizarLineaActual}
        />
      </div>
    </div>
  );
});

export default Editor;
export { LINE_HEIGHT };
