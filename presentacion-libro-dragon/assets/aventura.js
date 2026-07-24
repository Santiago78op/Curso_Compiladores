/* ============================================================
   aventura.js — las 6 interacciones del deck.

   Reconstrucción en JS plano del diseño de claude.ai/design,
   cuyo original corre sobre DCLogic (estado + sc-for + sc-if),
   un runtime que solo existe dentro de claude.ai. Acá el estado
   son variables y el re-render es innerHTML, con LOS MISMOS
   estilos en línea del diseño para que se vea idéntico.

   El escalado, la navegación y el HUD los pone escenario.js.
   ============================================================ */
(function () {
  'use strict';

  var MONO = "'JetBrains Mono',monospace";
  var DISPLAY = "'Bricolage Grotesque',sans-serif";

  /* ---------- datos (fig. 1.7 del Dragón) ---------- */

  var FASES = [
    { key: 'lex', num: 1, nombre: 'Análisis Léxico', en: 'Lexical', grupo: 'ANÁLISIS · front-end', acento: '#35d29b',
      desc: 'Lee los caracteres del código fuente y los agrupa en secuencias con significado llamadas tokens (lexemas ya clasificados).',
      ejemplo: '⟨id,initial⟩ ⟨+⟩ ⟨id,rate⟩ ⟨*⟩ ⟨num,60⟩',
      entra: 'Cadena de caracteres', sale: 'Flujo de tokens', con: 'JFlex · Jison · ANTLR' },
    { key: 'syn', num: 2, nombre: 'Análisis Sintáctico', en: 'Syntax', grupo: 'ANÁLISIS · front-end', acento: '#35d29b',
      desc: 'Verifica que los tokens sigan la gramática del lenguaje y construye un árbol sintáctico que refleja su estructura jerárquica.',
      ejemplo: '=\n├─ id position\n└─ + …',
      entra: 'Flujo de tokens', sale: 'Árbol sintáctico', con: 'CUP (LALR) · Jison · ANTLR' },
    { key: 'sem', num: 3, nombre: 'Análisis Semántico', en: 'Semantic', grupo: 'ANÁLISIS · front-end', acento: '#35d29b',
      desc: 'Revisa la coherencia del programa: tipos, declaraciones y ámbitos. Usa la tabla de símbolos e inserta las coerciones de tipo que hagan falta.',
      ejemplo: 'rate * inttofloat(60)',
      entra: 'Árbol sintáctico', sale: 'Árbol anotado', con: 'Código propio + tabla de símbolos' },
    { key: 'ir', num: 4, nombre: 'Cód. Intermedio', en: 'Intermediate', grupo: 'SÍNTESIS · back-end', acento: '#ff5a36',
      desc: 'Genera una representación intermedia sencilla e independiente de la máquina, como el código de tres direcciones.',
      ejemplo: 't1 = inttofloat(60)\nt2 = id rate * t1',
      entra: 'Árbol anotado', sale: 'Código de 3 direcciones', con: 'Fuera de alcance en OLC1 (§6.2)' },
    { key: 'opt', num: 5, nombre: 'Optimización', en: 'Optimization', grupo: 'SÍNTESIS · back-end', acento: '#ff5a36',
      desc: 'Mejora el código intermedio para que sea más rápido o más pequeño, sin cambiar lo que el programa hace.',
      ejemplo: 't1 = id rate * 60.0\nid position = …',
      entra: 'Código intermedio', sale: 'Código optimizado', con: 'Fuera de alcance (caps. 8–9)' },
    { key: 'gen', num: 6, nombre: 'Generación de Código', en: 'Code gen', grupo: 'SÍNTESIS · back-end', acento: '#ff5a36',
      desc: 'Traduce el código intermedio a código de máquina objetivo, asignando registros y seleccionando instrucciones.',
      ejemplo: 'LDF  R2, id rate\nMULF R2, R2, #60.0',
      entra: 'Código optimizado', sale: 'Código objeto', con: 'Fuera de alcance (cap. 8)' }
  ];

  var CONCEPTOS = [
    { key: 'trad', tag: 'IDEA 1', acento: '#35d29b', titulo: 'Es un traductor',
      cuerpo: 'Toma código escrito por personas y produce código equivalente que la máquina puede ejecutar. Equivalente = mismo comportamiento observable.' },
    { key: 'ana', tag: 'IDEA 2', acento: '#ffb020', titulo: 'Análisis',
      cuerpo: 'El front-end entiende el programa: lo descompone, lo estructura y verifica que tenga sentido. Es lo único que hacen los 4 proyectos del curso.' },
    { key: 'sin', tag: 'IDEA 3', acento: '#ff5a36', titulo: 'Síntesis',
      cuerpo: 'El back-end construye el programa destino a partir de esa representación. Un intérprete la reemplaza por ejecución directa.' }
  ];

  var PASOS = [
    { titulo: 'Código fuente', sub: 'Lo que escribe el programador', acento: '#e6e2d8',
      texto: 'position = initial + rate * 60' },
    { titulo: 'Análisis léxico', sub: 'Cadena → tokens', acento: '#35d29b',
      texto: '⟨id,position⟩ ⟨=⟩ ⟨id,initial⟩ ⟨+⟩ ⟨id,rate⟩ ⟨*⟩ ⟨num,60⟩' },
    { titulo: 'Análisis sintáctico', sub: 'Tokens → árbol sintáctico', acento: '#35d29b',
      texto: '        =\n       ╱ ╲\n position   +\n           ╱ ╲\n     initial   *\n             ╱ ╲\n         rate   60' },
    { titulo: 'Análisis semántico', sub: 'Coerción de tipos insertada', acento: '#ffb020',
      texto: '        =\n       ╱ ╲\n position   +\n           ╱ ╲\n     initial   *\n             ╱ ╲\n         rate  inttofloat\n                  │\n                  60' },
    { titulo: 'Código intermedio', sub: 'Tres direcciones', acento: '#ff5a36',
      texto: 't1 = inttofloat(60)\nt2 = id rate * t1\nt3 = id initial + t2\nid position = t3' },
    { titulo: 'Optimización', sub: 'Menos instrucciones, mismo resultado', acento: '#ff5a36',
      texto: 't1 = id rate * 60.0\nid position = id initial + t1' },
    { titulo: 'Código objeto', sub: 'Ensamblador destino', acento: '#ff5a36',
      texto: 'LDF   R2, id rate\nMULF  R2, R2, #60.0\nLDF   R1, id initial\nADDF  R1, R1, R2\nSTF   id position, R1' }
  ];

  var LEXEMAS = [
    { tipo: 'KEYWORD', val: 'if',     color: '#ff5a36' },
    { tipo: 'PUNCT',   val: '(',      color: '#9aa2b0' },
    { tipo: 'ID',      val: 'x',      color: '#35d29b' },
    { tipo: 'OP',      val: '&gt;=',  color: '#ffb020' },
    { tipo: 'NUM',     val: '10',     color: '#9b8cff' },
    { tipo: 'PUNCT',   val: ')',      color: '#9aa2b0' },
    { tipo: 'KEYWORD', val: 'return', color: '#ff5a36' },
    { tipo: 'ID',      val: 'x',      color: '#35d29b' },
    { tipo: 'PUNCT',   val: ';',      color: '#9aa2b0' }
  ];

  var PREGUNTAS = [
    { preg: '¿Qué fase agrupa los caracteres del código en tokens?',
      ops: ['Análisis sintáctico', 'Análisis léxico', 'Optimización'], ok: 1,
      exp: '✓ El análisis léxico (scanner) convierte la cadena de caracteres en un flujo de tokens. El sintáctico ya trabaja sobre esos tokens, no sobre caracteres.' },
    { preg: '¿Qué estructura produce el análisis sintáctico?',
      ops: ['Un árbol sintáctico', 'Código máquina', 'Una tabla hash'], ok: 0,
      exp: '✓ El parser construye un árbol sintáctico que refleja la estructura gramatical. La tabla hash es la tabla de símbolos, que acompaña a todas las fases.' },
    { preg: '¿Quién escribió el <span style="font-family:' + MONO + '">inttofloat(60)</span> del ejemplo?',
      ops: ['El programador', 'El análisis semántico', 'El optimizador'], ok: 1,
      exp: '✓ Nadie lo escribió en el fuente: lo insertó el análisis semántico al detectar que rate es float y 60 es int (coerción, §6.5.2).' },
    { preg: '¿Cuáles de las 6 fases implementan los proyectos de OLC1?',
      ops: ['Las 6', 'Solo las 3 de análisis', 'Solo léxico y sintáctico'], ok: 1,
      exp: '✓ Los 4 proyectos son intérpretes: hacen léxico, sintáctico y semántico, y después ejecutan. Los caps. 8–12 quedan fuera de alcance.' }
  ];

  var ORDEN_OK = ['lex', 'syn', 'sem', 'ir', 'opt', 'gen'];

  /* ============================================================
     1 · Tarjetas de concepto
     ============================================================ */
  (function () {
    var cont = document.getElementById('conceptos');
    if (!cont) return;
    var abierto = {};

    function pintar() {
      cont.innerHTML = '';
      CONCEPTOS.forEach(function (c) {
        var op = !!abierto[c.key];
        var d = document.createElement('div');
        d.setAttribute('role', 'button');
        d.setAttribute('tabindex', '0');
        d.style.cssText = 'cursor:pointer;transition:all .2s ease;border:1.5px solid ' +
          (op ? c.acento : 'rgba(255,255,255,.1)') + ';border-radius:18px;background:' +
          (op ? 'linear-gradient(160deg,rgba(255,255,255,.06),rgba(255,255,255,.01))' : '#15171f') +
          ';padding:26px 26px;min-height:150px;' + (op ? 'box-shadow:0 0 30px -8px ' + c.acento : '');
        d.innerHTML =
          '<div style="font-family:' + MONO + ';font-size:14px;letter-spacing:.1em;color:' + c.acento + ';margin-bottom:10px">' + c.tag + '</div>' +
          '<div style="font-family:' + DISPLAY + ';font-weight:700;font-size:30px;margin-bottom:6px">' + c.titulo + '</div>' +
          (op
            ? '<div style="font-size:19px;line-height:1.45;color:#c9cdd6">' + c.cuerpo + '</div>'
            : '<div style="font-size:16px;color:#6f7683;font-family:' + MONO + '">▸ clic para revelar</div>');
        function alternar() { abierto[c.key] = !abierto[c.key]; pintar(); }
        d.addEventListener('click', alternar);
        d.addEventListener('keydown', function (e) {
          if (e.key === 'Enter') { e.preventDefault(); e.stopPropagation(); alternar(); }
        });
        cont.appendChild(d);
      });
    }
    pintar();
  })();

  /* ============================================================
     2 · Selector de fases
     ============================================================ */
  (function () {
    var chips = document.getElementById('chips-fase');
    var panel = document.getElementById('panel-fase');
    if (!chips || !panel) return;
    var act = 0;

    function pintar() {
      chips.innerHTML = '';
      FASES.forEach(function (f, k) {
        var on = k === act;
        var b = document.createElement('div');
        b.setAttribute('role', 'button');
        b.setAttribute('tabindex', '0');
        b.style.cssText = 'flex:1;min-width:0;cursor:pointer;padding:16px 18px;border-radius:16px;transition:all .18s ease;border:1.5px solid ' +
          (on ? f.acento : 'rgba(255,255,255,.1)') + ';background:' +
          (on ? 'linear-gradient(160deg,rgba(255,255,255,.08),rgba(255,255,255,.02))' : '#15171f') + ';' +
          (on ? 'box-shadow:0 0 24px -6px ' + f.acento : '');
        b.innerHTML =
          '<div style="font-family:' + MONO + ';font-size:13px;color:' + f.acento + '">0' + f.num + '</div>' +
          '<div style="font-family:' + DISPLAY + ';font-weight:700;font-size:21px;line-height:1.05;margin-top:8px">' + f.nombre + '</div>' +
          '<div style="font-size:14px;color:#9aa2b0;font-family:' + MONO + ';margin-top:4px">' + f.en + '</div>';
        b.addEventListener('click', function () { act = k; pintar(); });
        b.addEventListener('keydown', function (e) {
          if (e.key === 'Enter') { e.preventDefault(); e.stopPropagation(); act = k; pintar(); }
        });
        chips.appendChild(b);
      });

      var f = FASES[act];
      panel.innerHTML =
        '<div style="min-width:0">' +
          '<div style="display:inline-flex;align-items:center;gap:10px;padding:6px 14px;border-radius:999px;background:rgba(255,255,255,.06);font-family:' + MONO + ';font-size:14px;color:' + f.acento + '">' +
            '<span style="width:8px;height:8px;border-radius:50%;background:' + f.acento + '"></span>' + f.grupo + '</div>' +
          '<h3 style="font-family:' + DISPLAY + ';font-weight:800;font-size:44px;margin:16px 0 4px">' + f.nombre + '</h3>' +
          '<div style="font-family:' + MONO + ';font-size:18px;color:#9aa2b0;margin-bottom:20px">' + f.en + '</div>' +
          '<p style="font-size:24px;line-height:1.5;color:#e6e2d8;margin:0 0 20px">' + f.desc + '</p>' +
          '<div style="font-family:' + MONO + ';font-size:15px;color:#6f7683;letter-spacing:.1em;margin-bottom:8px">// EJEMPLO</div>' +
          '<pre id="fase-ejemplo" style="margin:0;font-family:' + MONO + ';font-size:20px;color:#35d29b;background:#0d0e12;border:1px solid rgba(255,255,255,.08);border-radius:12px;padding:18px 20px;white-space:pre-wrap"></pre>' +
        '</div>' +
        '<div style="border-left:1px solid rgba(255,255,255,.08);padding-left:32px;display:flex;flex-direction:column;justify-content:center;min-width:0">' +
          '<div style="font-family:' + MONO + ';font-size:15px;color:#6f7683;letter-spacing:.1em;margin-bottom:14px">// ENTRA → SALE</div>' +
          '<div style="font-size:18px;color:#9aa2b0;margin-bottom:6px">Recibe</div>' +
          '<div style="font-size:22px;font-weight:700;color:#ffb020;margin-bottom:22px">' + f.entra + '</div>' +
          '<div style="font-family:' + DISPLAY + ';font-size:34px;color:#ff5a36;margin-bottom:22px">↓</div>' +
          '<div style="font-size:18px;color:#9aa2b0;margin-bottom:6px">Produce</div>' +
          '<div style="font-size:22px;font-weight:700;color:#35d29b">' + f.sale + '</div>' +
          '<div style="font-size:18px;color:#9aa2b0;margin:26px 0 6px">Con qué se hace</div>' +
          '<div style="font-family:' + MONO + ';font-size:16px;color:#c9cdd6">' + f.con + '</div>' +
        '</div>';
      panel.querySelector('#fase-ejemplo').textContent = f.ejemplo;
    }
    pintar();
  })();

  /* ============================================================
     3 · Tokenizador
     ============================================================ */
  (function () {
    var btn = document.getElementById('btn-lex');
    var reset = document.getElementById('btn-lex-reset');
    var salida = document.getElementById('lex-salida');
    var conteo = document.getElementById('lex-conteo');
    if (!btn || !salida) return;

    function correr() {
      var h = '<div style="display:flex;flex-wrap:wrap;gap:14px">';
      LEXEMAS.forEach(function (t, k) {
        h += '<div style="animation:popIn .35s cubic-bezier(.2,1.4,.4,1) both;animation-delay:' + (k * 0.07).toFixed(2) +
             's;padding:14px 20px;border-radius:12px;background:' + t.color + '22;border:1.5px solid ' + t.color + ';color:' + t.color + '">' +
               '<div style="font-family:' + MONO + ';font-size:12px;letter-spacing:.08em;opacity:.85;margin-bottom:2px">' + t.tipo + '</div>' +
               '<div style="font-family:' + MONO + ';font-size:24px;font-weight:700">' + t.val + '</div>' +
             '</div>';
      });
      salida.innerHTML = h + '</div>';
      if (conteo) conteo.textContent = LEXEMAS.length + ' tokens encontrados';
    }
    function limpiar() {
      salida.innerHTML = '<div style="border:1px dashed rgba(255,255,255,.16);border-radius:16px;padding:40px;text-align:center;color:#6f7683;font-family:' + MONO + ';font-size:20px">los tokens aparecerán aquí…</div>';
      if (conteo) conteo.textContent = '';
    }
    btn.addEventListener('click', correr);
    if (reset) reset.addEventListener('click', limpiar);
    limpiar();
  })();

  /* ============================================================
     4 · El viaje completo
     ============================================================ */
  (function () {
    var puntos = document.getElementById('puntos');
    var marco = document.getElementById('marco-codigo');
    var prev = document.getElementById('paso-prev');
    var next = document.getElementById('paso-next');
    if (!puntos || !marco) return;
    var k = 0;

    function pintar() {
      puntos.innerHTML = '';
      PASOS.forEach(function (p, j) {
        var on = j === k;
        var d = document.createElement('div');
        d.setAttribute('role', 'button');
        d.setAttribute('tabindex', '0');
        d.title = p.titulo;
        d.textContent = String(j);
        d.style.cssText = 'cursor:pointer;width:40px;height:40px;border-radius:10px;display:flex;align-items:center;justify-content:center;font-family:' + MONO +
          ';font-size:16px;font-weight:700;transition:all .18s;border:1.5px solid ' + (on ? p.acento : 'rgba(255,255,255,.12)') +
          ';background:' + (on ? p.acento : '#15171f') + ';color:' + (on ? '#0d0e12' : '#9aa2b0');
        d.addEventListener('click', function () { k = j; pintar(); });
        d.addEventListener('keydown', function (e) {
          if (e.key === 'Enter') { e.preventDefault(); e.stopPropagation(); k = j; pintar(); }
        });
        puntos.appendChild(d);
      });

      var p = PASOS[k];
      var pct = Math.round((k / (PASOS.length - 1)) * 100);
      marco.innerHTML =
        '<div style="position:absolute;top:0;left:0;height:4px;background:linear-gradient(90deg,#35d29b,#ffb020,#ff5a36);width:' + pct + '%;transition:width .35s ease"></div>' +
        '<div style="display:flex;align-items:flex-start;justify-content:space-between;gap:24px;margin-bottom:18px">' +
          '<div><div style="font-family:' + MONO + ';font-size:14px;color:' + p.acento + ';letter-spacing:.14em">PASO ' + (k + 1) + ' / ' + PASOS.length + '</div>' +
          '<h3 style="font-family:' + DISPLAY + ';font-weight:800;font-size:38px;margin:4px 0 0">' + p.titulo + '</h3></div>' +
          '<div style="font-family:' + MONO + ';font-size:16px;color:#9aa2b0;text-align:right;max-width:360px">' + p.sub + '</div>' +
        '</div>' +
        '<pre id="paso-texto" style="margin:0;font-family:' + MONO + ';font-size:27px;line-height:1.6;color:#e6e2d8;white-space:pre-wrap"></pre>';
      marco.querySelector('#paso-texto').textContent = p.texto;

      if (prev) prev.style.opacity = k === 0 ? '.4' : '1';
      if (next) next.style.opacity = k === PASOS.length - 1 ? '.4' : '1';
    }
    if (prev) prev.addEventListener('click', function () { if (k > 0) { k--; pintar(); } });
    if (next) next.addEventListener('click', function () { if (k < PASOS.length - 1) { k++; pintar(); } });
    pintar();
  })();

  /* ============================================================
     5 · Ordenar el pipeline
     ============================================================ */
  (function () {
    var fila = document.getElementById('fila-drag');
    var marcador = document.getElementById('marcador');
    var barajar = document.getElementById('barajar');
    if (!fila || !marcador) return;

    var orden = [3, 0, 5, 1, 4, 2];   // el desorden inicial del diseño
    var desde = null;

    function resuelto(o) { return o.every(function (fi, idx) { return FASES[fi].key === ORDEN_OK[idx]; }); }
    function revolver() {
      var a = [0, 1, 2, 3, 4, 5], x, y, t;
      for (x = a.length - 1; x > 0; x--) { y = Math.floor(Math.random() * (x + 1)); t = a[x]; a[x] = a[y]; a[y] = t; }
      return resuelto(a) ? revolver() : a;
    }

    function pintar() {
      var n = orden.reduce(function (acc, fi, idx) { return acc + (FASES[fi].key === ORDEN_OK[idx] ? 1 : 0); }, 0);
      var ok = n === 6;

      marcador.style.cssText = 'align-self:flex-start;margin-top:14px;font-family:' + MONO +
        ';font-size:19px;font-weight:700;padding:12px 22px;border-radius:999px;border:1.5px solid ' +
        (ok ? '#35d29b' : 'rgba(255,255,255,.16)') + ';background:' +
        (ok ? 'rgba(53,210,155,.14)' : 'rgba(255,255,255,.04)') + ';color:' + (ok ? '#35d29b' : '#c9cdd6');
      marcador.textContent = ok ? '🎉 ¡Perfecto! Las 6 fases en orden.' : n + ' / 6 en su lugar correcto — seguí arrastrando';

      fila.innerHTML = '';
      orden.forEach(function (fi, idx) {
        var f = FASES[fi];
        var bien = f.key === ORDEN_OK[idx];
        var d = document.createElement('div');
        d.draggable = true;
        d.style.cssText = 'position:relative;flex:1;min-width:0;cursor:grab;user-select:none;min-height:190px;display:flex;flex-direction:column;padding:22px 20px;border-radius:16px;transition:border-color .18s,background .18s;border:1.5px solid ' +
          (bien ? '#35d29b' : 'rgba(255,255,255,.14)') + ';background:' +
          (bien ? 'rgba(53,210,155,.1)' : 'linear-gradient(160deg,#191c25,#14161d)');
        d.innerHTML =
          '<div style="font-family:' + MONO + ';font-size:26px;color:#6f7683;margin-bottom:auto">⠿</div>' +
          '<div style="font-family:' + DISPLAY + ';font-weight:700;font-size:23px;line-height:1.1;margin-top:14px">' + f.nombre + '</div>' +
          '<div style="font-size:15px;color:#9aa2b0;font-family:' + MONO + ';margin-top:6px">' + f.en + '</div>' +
          (bien ? '<div style="position:absolute;top:14px;right:16px;color:#35d29b;font-size:22px">✓</div>' : '');

        d.addEventListener('dragstart', function (e) {
          desde = idx; d.style.opacity = '.45';
          if (e.dataTransfer) { e.dataTransfer.effectAllowed = 'move'; e.dataTransfer.setData('text/plain', String(idx)); }
        });
        d.addEventListener('dragover', function (e) { e.preventDefault(); d.style.borderColor = '#ffb020'; });
        d.addEventListener('dragleave', function () { d.style.borderColor = bien ? '#35d29b' : 'rgba(255,255,255,.14)'; });
        d.addEventListener('drop', function (e) {
          e.preventDefault();
          if (desde === null || desde === idx) return;
          var mov = orden.splice(desde, 1)[0];
          orden.splice(idx, 0, mov);
          desde = null; pintar();
        });
        d.addEventListener('dragend', function () { desde = null; d.style.opacity = '1'; });
        fila.appendChild(d);
      });
    }
    if (barajar) barajar.addEventListener('click', function () { orden = revolver(); pintar(); });
    pintar();
  })();

  /* ============================================================
     6 · Quiz
     ============================================================ */
  (function () {
    var grid = document.getElementById('grid-quiz');
    var score = document.getElementById('quiz-score');
    if (!grid) return;
    var elegidas = {};

    function pintar() {
      var aciertos = 0;
      grid.innerHTML = '';
      PREGUNTAS.forEach(function (q, qi) {
        var pick = elegidas[qi];
        var resp = pick !== undefined;
        if (resp && pick === q.ok) aciertos++;

        var carta = document.createElement('div');
        carta.style.cssText = 'border:1px solid rgba(255,255,255,.1);border-radius:18px;background:#15171f;padding:26px 24px;display:flex;flex-direction:column;min-width:0';
        var h =
          '<div style="font-family:' + MONO + ';font-size:14px;color:#ffb020;margin-bottom:10px">P' + (qi + 1) + '</div>' +
          '<div style="font-family:' + DISPLAY + ';font-weight:700;font-size:22px;line-height:1.2;margin-bottom:18px;min-height:80px">' + q.preg + '</div>' +
          '<div style="display:flex;flex-direction:column;gap:10px">';
        q.ops.forEach(function (label, oi) {
          var bg = '#1c1f2a', bd = 'rgba(255,255,255,.12)', col = '#e6e2d8';
          if (resp) {
            if (oi === q.ok) { bg = 'rgba(53,210,155,.16)'; bd = '#35d29b'; col = '#35d29b'; }
            else if (oi === pick) { bg = 'rgba(255,90,54,.14)'; bd = '#ff5a36'; col = '#ff5a36'; }
            else { col = '#6f7683'; }
          }
          h += '<div data-op="' + oi + '" role="button" tabindex="0" style="cursor:' + (resp ? 'default' : 'pointer') +
               ';padding:14px 18px;border-radius:12px;font-size:18px;transition:all .15s;border:1.5px solid ' + bd +
               ';background:' + bg + ';color:' + col + '">' + label + '</div>';
        });
        h += '</div>';
        if (resp) {
          h += '<div style="margin-top:16px;font-size:16px;line-height:1.4;color:#c9cdd6;border-top:1px solid rgba(255,255,255,.08);padding-top:14px">' + q.exp + '</div>';
        }
        carta.innerHTML = h;

        if (!resp) {
          carta.querySelectorAll('[data-op]').forEach(function (o) {
            function elegir(e) { e.stopPropagation(); elegidas[qi] = +o.dataset.op; pintar(); }
            o.addEventListener('click', elegir);
            o.addEventListener('keydown', function (e) { if (e.key === 'Enter') { e.preventDefault(); elegir(e); } });
          });
        }
        grid.appendChild(carta);
      });
      if (score) score.innerHTML = aciertos + '<span style="color:#6f7683;font-size:30px">/' + PREGUNTAS.length + '</span>';
    }
    pintar();
  })();

  /* el riel de escenario.js clona las diapositivas: avisarle que ya
     están rellenas para que las miniaturas no salgan vacías. */
  document.dispatchEvent(new CustomEvent('aventura:listo'));
})();
