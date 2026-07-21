/* ============================================================
   Servidor Express (REST) de CompInterpreter.
   Endpoint principal: POST /interpretar  { codigo: "<fuente .ci>" }
   Responde JSON: { errores, consola, simbolos, ast, dot }.
   ============================================================ */
const express = require('express');
const cors = require('cors');
const { analizar } = require('./src/analizar');

const app = express();
const PUERTO = process.env.PORT || 4000;

app.use(cors());                              // CORS para el cliente de dev (Vite/CRA)
app.use(express.json({ limit: '5mb' }));

app.get('/', (req, res) => {
  res.json({ nombre: 'CompInterpreter', estado: 'ok', endpoint: 'POST /interpretar { codigo }' });
});

app.get('/salud', (req, res) => res.json({ estado: 'ok' }));

app.post('/interpretar', (req, res) => {
  const codigo = (req.body && typeof req.body.codigo === 'string') ? req.body.codigo : '';
  try {
    const resultado = analizar(codigo);
    res.json(resultado);
  } catch (e) {
    res.status(500).json({
      errores: [{ tipo: 'Interno', descripcion: e.message || String(e), linea: 0, columna: 0 }],
      consola: '', simbolos: [], ast: { nodes: [], edges: [] }, dot: ''
    });
  }
});

app.listen(PUERTO, () => {
  console.log('CompInterpreter server escuchando en http://localhost:' + PUERTO);
});
