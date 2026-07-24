/* ============================================================
   escenario.js — el motor del "deck-stage" del diseño.

   Reproduce el comportamiento del web component de
   claude.ai/design sin shadow DOM ni dependencias:

     · escala el lienzo 1920×1080 con transform: scale() para
       que quepa siempre en la ventana (letterboxing);
     · navegación por teclado ←/→ ↑/↓ PgUp/PgDn Espacio
       Home/End, dígitos 1-9, R para volver al inicio;
     · deep-link por hash (#3 abre la diapositiva 3);
     · HUD inferior que aparece al mover el mouse y se
       desvanece a los 2.6 s de quietud;
     · riel de miniaturas (T) con clones escalados;
     · notas del orador (N) desde data-speaker-notes;
     · en móvil, tocar la mitad izquierda/derecha navega.

   Expone window.Escenario { ir, actual, total } para que
   otros scripts (aventura.js) puedan reaccionar.
   ============================================================ */
(function () {
  'use strict';

  var DISENO_W = 1920, DISENO_H = 1080;

  var escenario = document.getElementById('escenario');
  var lienzo = document.getElementById('lienzo');
  if (!escenario || !lienzo) return;

  var slides = Array.prototype.slice.call(lienzo.querySelectorAll(':scope > section'));
  if (!slides.length) return;

  var hud = document.getElementById('hud');
  var riel = document.getElementById('riel');
  var notas = document.getElementById('notas');
  var notasTexto = document.getElementById('notas-texto');
  var elCur = document.getElementById('hud-cur');
  var elTot = document.getElementById('hud-tot');

  var i = 0;
  var rielAbierto = false;

  /* ---------- escalado: el lienzo entero cabe en la ventana ---------- */
  function ajustar() {
    var dispW = window.innerWidth - (rielAbierto ? parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--riel-w')) || 196 : 0);
    var dispH = window.innerHeight;
    var s = Math.min(dispW / DISENO_W, dispH / DISENO_H);
    lienzo.style.transform = 'scale(' + s + ')';
  }

  /* ---------- navegación ---------- */
  function ir(n, desdeHash) {
    i = Math.max(0, Math.min(slides.length - 1, n));
    slides.forEach(function (s, k) {
      if (k === i) s.setAttribute('data-activa', '');
      else s.removeAttribute('data-activa');
    });
    if (elCur) elCur.textContent = String(i + 1);
    if (riel) {
      Array.prototype.forEach.call(riel.children, function (m, k) {
        if (k === i) { m.setAttribute('data-activa', ''); }
        else m.removeAttribute('data-activa');
      });
      var act = riel.children[i];
      if (act && rielAbierto) act.scrollIntoView({ block: 'nearest' });
    }
    if (notasTexto) {
      var t = slides[i].getAttribute('data-speaker-notes') || 'Esta diapositiva no tiene notas.';
      notasTexto.textContent = t;
    }
    if (!desdeHash) {
      try { history.replaceState(null, '', '#' + (i + 1)); } catch (e) { /* file:// estricto */ }
    }
    document.dispatchEvent(new CustomEvent('escenario:cambio', { detail: { indice: i } }));
  }
  function avanzar(d) { ir(i + d); }

  /* ---------- HUD que se desvanece ---------- */
  var temporizador = null, hudFijo = false;
  function mostrarHud() {
    if (!hud) return;
    hud.setAttribute('data-visible', '');
    clearTimeout(temporizador);
    if (!hudFijo) temporizador = setTimeout(function () { hud.removeAttribute('data-visible'); }, 2600);
  }
  if (hud) {
    hud.addEventListener('mouseenter', function () { hudFijo = true; clearTimeout(temporizador); });
    hud.addEventListener('mouseleave', function () { hudFijo = false; mostrarHud(); });
    hud.addEventListener('focusin', function () { hudFijo = true; clearTimeout(temporizador); });
    hud.addEventListener('focusout', function () { hudFijo = false; mostrarHud(); });
  }
  window.addEventListener('mousemove', mostrarHud, { passive: true });

  /* ---------- riel de miniaturas ---------- */
  function construirRiel() {
    if (!riel) return;
    var anchoMini = 176;                       // 196 - padding
    var escala = anchoMini / DISENO_W;
    riel.innerHTML = '';
    slides.forEach(function (s, k) {
      var mini = document.createElement('div');
      mini.className = 'mini';
      mini.style.width = anchoMini + 'px';
      mini.style.height = Math.round(DISENO_H * escala) + 'px';
      mini.title = s.getAttribute('data-label') || ('Diapositiva ' + (k + 1));

      var lupa = document.createElement('div');
      lupa.className = 'lupa';
      lupa.style.transform = 'scale(' + escala + ')';
      var clon = s.cloneNode(true);
      clon.removeAttribute('data-activa');
      // el clon es estático: sin ids duplicados ni interacción
      clon.querySelectorAll('[id]').forEach(function (n) { n.removeAttribute('id'); });
      lupa.appendChild(clon);
      mini.appendChild(lupa);

      var n = document.createElement('span');
      n.className = 'n'; n.textContent = String(k + 1);
      mini.appendChild(n);

      var rot = document.createElement('span');
      rot.className = 'rot'; rot.textContent = s.getAttribute('data-label') || '';
      mini.appendChild(rot);

      mini.addEventListener('click', function () { ir(k); });
      riel.appendChild(mini);
    });
  }
  function alternarRiel() {
    rielAbierto = !rielAbierto;
    if (riel) riel.classList.toggle('abierto', rielAbierto);
    escenario.classList.toggle('con-riel', rielAbierto);
    ajustar();
    ir(i);
  }

  /* ---------- notas ---------- */
  function alternarNotas() { if (notas) notas.classList.toggle('abierto'); }

  /* ---------- teclado ---------- */
  document.addEventListener('keydown', function (e) {
    if (e.metaKey || e.ctrlKey || e.altKey) return;
    var t = e.target;
    if (t && (t.tagName === 'INPUT' || t.tagName === 'TEXTAREA' || t.isContentEditable)) return;

    switch (e.key) {
      case 'ArrowRight': case 'ArrowDown': case 'PageDown': case ' ':
        e.preventDefault(); avanzar(1); mostrarHud(); break;
      case 'ArrowLeft': case 'ArrowUp': case 'PageUp':
        e.preventDefault(); avanzar(-1); mostrarHud(); break;
      case 'Home': e.preventDefault(); ir(0); mostrarHud(); break;
      case 'End': e.preventDefault(); ir(slides.length - 1); mostrarHud(); break;
      case 'r': case 'R': ir(0); mostrarHud(); break;
      case 't': case 'T': alternarRiel(); mostrarHud(); break;
      case 'n': case 'N': alternarNotas(); mostrarHud(); break;
      default:
        if (/^[1-9]$/.test(e.key)) { ir(parseInt(e.key, 10) - 1); mostrarHud(); }
    }
  });

  /* ---------- táctil: mitad izquierda / derecha ---------- */
  escenario.addEventListener('click', function (e) {
    if (!matchMedia('(pointer: coarse)').matches) return;
    var interactivo = e.target.closest('a, button, input, [role="button"], [draggable="true"]');
    if (interactivo) return;
    avanzar(e.clientX < window.innerWidth / 2 ? -1 : 1);
    mostrarHud();
  });

  /* ---------- cableado del HUD ---------- */
  function conectar(id, fn) { var b = document.getElementById(id); if (b) b.addEventListener('click', fn); }
  conectar('hud-prev', function () { avanzar(-1); });
  conectar('hud-next', function () { avanzar(1); });
  conectar('hud-reset', function () { ir(0); });
  conectar('hud-riel', alternarRiel);
  conectar('hud-notas', alternarNotas);

  /* ---------- arranque ---------- */
  if (elTot) elTot.textContent = String(slides.length);
  construirRiel();
  window.addEventListener('resize', ajustar, { passive: true });
  ajustar();

  var inicial = parseInt((location.hash || '').replace('#', ''), 10);
  ir(isFinite(inicial) && inicial >= 1 ? inicial - 1 : 0, true);
  window.addEventListener('hashchange', function () {
    var n = parseInt((location.hash || '').replace('#', ''), 10);
    if (isFinite(n) && n >= 1 && n - 1 !== i) ir(n - 1, true);
  });
  mostrarHud();

  /* el riel clona las diapositivas al arrancar; si aventura.js las
     rellena después, se reconstruye para que las miniaturas no queden
     vacías. */
  document.addEventListener('aventura:listo', function () {
    construirRiel();
    ir(i, true);
  });

  window.Escenario = {
    ir: ir,
    actual: function () { return i; },
    total: slides.length
  };
})();
