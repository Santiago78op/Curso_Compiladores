/* ============================================================
   Valor: representacion de un valor en tiempo de ejecucion,
   etiquetado con su tipo. Los vectores guardan tipoBase y dimension.
   ============================================================ */
const { TIPO } = require('./tipos');

class Valor {
  constructor(tipo, valor) {
    this.tipo = tipo;      // 'int' | 'double' | 'bool' | 'char' | 'string' | 'null' | 'vector'
    this.valor = valor;    // JS number/boolean/string/null, o Array<Valor> (o Array<Array<Valor>>) para vectores
    this.tipoBase = null;  // solo vectores: tipo de los elementos
    this.dimension = 0;    // solo vectores: 1 o 2
  }

  static int(v) { return new Valor(TIPO.INT, Math.trunc(v)); }
  static double(v) { return new Valor(TIPO.DOUBLE, v); }
  static bool(v) { return new Valor(TIPO.BOOL, !!v); }
  static char(v) { return new Valor(TIPO.CHAR, v); }
  static string(v) { return new Valor(TIPO.STRING, v); }
  static nulo() { return new Valor(TIPO.NULL, null); }

  static vector(tipoBase, dimension, contenido) {
    const v = new Valor(TIPO.VECTOR, contenido);
    v.tipoBase = tipoBase;
    v.dimension = dimension;
    return v;
  }
}

/* Convierte un Valor numerico (int/double/char/bool) a numero JS. */
function aNumero(v) {
  switch (v.tipo) {
    case TIPO.INT:
    case TIPO.DOUBLE: return v.valor;
    case TIPO.BOOL: return v.valor ? 1 : 0;
    case TIPO.CHAR: return v.valor.charCodeAt(0);
    default: return NaN;
  }
}

/* Representacion textual para consola / concatenacion / reportes. */
function aTexto(v) {
  if (v === null || v === undefined) return 'null';
  switch (v.tipo) {
    case TIPO.NULL: return 'null';
    case TIPO.BOOL: return v.valor ? 'true' : 'false';
    case TIPO.DOUBLE: {
      // mostrar 15.0 como "15.0" pero mantener decimales reales
      if (Number.isInteger(v.valor)) return v.valor.toFixed(1);
      return String(v.valor);
    }
    case TIPO.INT: return String(v.valor);
    case TIPO.CHAR: return v.valor;
    case TIPO.STRING: return v.valor;
    case TIPO.VECTOR: {
      if (v.dimension === 2) {
        return '[' + v.valor.map(fila => '[' + fila.map(aTexto).join(', ') + ']').join(', ') + ']';
      }
      return '[' + v.valor.map(aTexto).join(', ') + ']';
    }
    default: return String(v.valor);
  }
}

module.exports = { Valor, aNumero, aTexto };
