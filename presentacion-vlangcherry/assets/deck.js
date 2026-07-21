/* ============================================================
   deck.js — navegación de diapositivas + tema claro/oscuro.
   Compartido por todas las páginas: en el índice (sin .stage)
   solo activa el botón de tema.
   ============================================================ */

/* ---- tema: respeta el sistema hasta que el usuario toca ◐; se recuerda ---- */
(function () {
  const root = document.documentElement;
  const saved = localStorage.getItem('df-theme');
  if (saved === 'dark' || saved === 'light') root.dataset.theme = saved;

  const btn = document.getElementById('theme');
  if (!btn) return;
  btn.addEventListener('click', () => {
    const dark = root.dataset.theme
      ? root.dataset.theme === 'dark'
      : matchMedia('(prefers-color-scheme: dark)').matches;
    const nuevo = dark ? 'light' : 'dark';
    root.dataset.theme = nuevo;
    localStorage.setItem('df-theme', nuevo);
  });
})();

/* ---- navegación de slides (solo si la página es un deck) ---- */
(function () {
  const slides = Array.from(document.querySelectorAll('.slide'));
  if (!slides.length) return;

  const cur = document.getElementById('cur');
  const tot = document.getElementById('tot');
  const bar = document.getElementById('bar');
  const prev = document.getElementById('prev');
  const next = document.getElementById('next');
  let i = 0;
  tot.textContent = slides.length;

  function show(n) {
    i = Math.max(0, Math.min(slides.length - 1, n));
    slides.forEach((s, k) => s.classList.toggle('active', k === i));
    cur.textContent = i + 1;
    bar.style.width = ((i + 1) / slides.length * 100) + '%';
    prev.disabled = i === 0;
    next.disabled = i === slides.length - 1;
    slides[i].scrollTop = 0;
  }
  prev.addEventListener('click', () => show(i - 1));
  next.addEventListener('click', () => show(i + 1));
  document.addEventListener('keydown', (e) => {
    if (e.key === 'ArrowRight' || e.key === 'PageDown' || e.key === ' ') { e.preventDefault(); show(i + 1); }
    if (e.key === 'ArrowLeft' || e.key === 'PageUp') { e.preventDefault(); show(i - 1); }
  });
  show(0);
})();
