/* ============================================================
   Orquestador: fuente .ci -> { errores, consola, simbolos, ast }.
   Conecta el parser Jison (lexico+sintactico), el interprete
   (semantico+ejecucion) y los reportes. Entorno FRESCO por llamada.
   ============================================================ */
const parserMod = require('./parser.js');
const { ListaErrores } = require('./interprete/errores');
const { Interprete } = require('./interprete/interprete');
const { construirGrafo, aDot } = require('./reportes/ast-grafo');

function obtenerParser() {
  // jison exporta { parser, Parser, parse, ... } o el parser directo
  if (parserMod.parser) return parserMod.parser;
  return parserMod;
}

function analizar(codigo) {
  const errores = new ListaErrores();
  const parser = obtenerParser();

  // yy compartido entre lexer y parser. El lexer empuja errores lexicos
  // directamente a este arreglo (misma referencia que errores.errores).
  parser.yy = { errores: errores.errores };

  // Recuperacion de errores sintacticos (modo panico via `error` en la gramatica).
  parser.yy.parseError = function (msg, hash) {
    let descripcion;
    if (hash && hash.token) {
      const encontrado = hash.text !== undefined && hash.text !== '' ? ('"' + hash.text + '"') : hash.token;
      descripcion = 'Se encontró ' + encontrado;
      if (hash.expected && hash.expected.length) {
        descripcion += ' y se esperaba: ' + hash.expected.slice(0, 6).join(', ');
      }
    } else {
      descripcion = msg;
    }
    const l = hash && hash.loc ? hash.loc.first_line : (hash ? hash.line + 1 : 0);
    const c = hash && hash.loc ? hash.loc.first_column + 1 : 0;
    errores.sintactico(descripcion, l, c);
    if (!hash || !hash.recoverable) {
      throw new Error(descripcion);   // lo atrapamos abajo
    }
  };
  // jison usa this.parseError; lo enlazamos tambien
  parser.parseError = parser.yy.parseError;

  let ast = null;
  try {
    ast = parser.parse(codigo);
  } catch (e) {
    // error sintactico no recuperable: los errores ya quedaron registrados.
    // Si no se registro ninguno (fallo inesperado), lo dejamos como sintactico.
    if (!errores.hay()) errores.sintactico(e.message || 'Error de análisis', 0, 0);
    ast = ast || null;
  }

  // Interprete (semantico + ejecucion). Best-effort aunque haya errores previos.
  const interprete = new Interprete(errores);
  if (ast) {
    try {
      interprete.interpretar(ast);
    } catch (e) {
      errores.semantico('Error interno durante la ejecución: ' + (e.message || e), 0, 0);
    }
  }

  const grafo = ast ? construirGrafo(ast) : { nodes: [], edges: [] };
  const listaErrores = errores.errores.slice().sort((a, b) =>
    (a.linea - b.linea) || (a.columna - b.columna));

  return {
    errores: listaErrores,
    consola: interprete.consola.join('\n'),
    consolaLineas: interprete.consola,
    simbolos: Array.from(interprete.simbolos.values()),
    ast: grafo,
    dot: aDot(grafo)
  };
}

module.exports = { analizar };
