/* ============================================================
   Tipos del lenguaje CompInterpreter.
   Internamente usamos los mismos nombres que las palabras clave
   (int, double, bool, char, string, null) + 'vector'.
   Los nombres "bonitos" en espanol son para los reportes (seccion 6).
   ============================================================ */

const TIPO = {
  INT: 'int',
  DOUBLE: 'double',
  BOOL: 'bool',
  CHAR: 'char',
  STRING: 'string',
  NULL: 'null',
  VECTOR: 'vector'
};

const NOMBRE_ES = {
  int: 'Entero',
  double: 'Decimal',
  bool: 'Boolean',
  char: 'Carácter',
  string: 'Cadena',
  null: 'Nulo',
  vector: 'Vector'
};

function nombreEspanol(t) {
  return NOMBRE_ES[t] || t;
}

// Valor por defecto de cada tipo (seccion 5.3)
function valorPorDefecto(tipo) {
  switch (tipo) {
    case TIPO.INT: return 0;
    case TIPO.DOUBLE: return 0.0;
    case TIPO.BOOL: return true;
    case TIPO.CHAR: return ' ';   // caracter espacio (codigo 32) - "carácter 0" del enunciado
    case TIPO.STRING: return '';
    case TIPO.NULL: return null;
    default: return null;
  }
}

const NUMERICOS = [TIPO.INT, TIPO.DOUBLE, TIPO.CHAR, TIPO.BOOL];

function esNumerico(t) {
  return t === TIPO.INT || t === TIPO.DOUBLE || t === TIPO.CHAR || t === TIPO.BOOL;
}

module.exports = { TIPO, NOMBRE_ES, nombreEspanol, valorPorDefecto, NUMERICOS, esNumerico };
