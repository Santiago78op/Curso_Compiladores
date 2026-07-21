/* ============================================================
   Operaciones aritmeticas, relacionales y logicas.
   Las tablas de compatibilidad son LITERALES del enunciado (5.5).
   Propagacion por null: si un operando es JS null (ya hubo un error
   antes), se devuelve null en silencio (un error por causa raiz).
   OJO: Valor de tipo 'null' (el null DEL LENGUAJE) NO es JS null.
   ============================================================ */
const { TIPO } = require('./tipos');
const { Valor, aNumero, aTexto } = require('./valor');

// I=int, D=double, B=bool, C=char, S=string ; 'X' = combinacion invalida
// Cada tabla: filas = operando izquierdo, columnas = operando derecho.

// 5.5.1 Suma
const T_SUMA = {
  int:    { int: 'int',    double: 'double', bool: 'int', char: 'int',  string: 'string' },
  double: { int: 'double', double: 'double', bool: 'double', char: 'double', string: 'string' },
  bool:   { int: 'int',    double: 'double', bool: 'X',   char: 'X',   string: 'string' },
  char:   { int: 'int',    double: 'double', bool: 'X',   char: 'string', string: 'string' },
  string: { int: 'string', double: 'string', bool: 'string', char: 'string', string: 'string' }
};
// 5.5.2 Resta
const T_RESTA = {
  int:    { int: 'int',    double: 'double', bool: 'int', char: 'int' },
  double: { int: 'double', double: 'double', bool: 'double', char: 'double' },
  bool:   { int: 'int',    double: 'double', bool: 'X',   char: 'X' },
  char:   { int: 'int',    double: 'double', bool: 'X',   char: 'X' }
};
// 5.5.3 Multiplicacion
const T_MULT = {
  int:    { int: 'int',    double: 'double', char: 'int' },
  double: { int: 'double', double: 'double', char: 'double' },
  char:   { int: 'int',    double: 'double', char: 'X' }
};
// 5.5.4 Division (siempre Decimal donde esta definida)
const T_DIV = {
  int:    { int: 'double', double: 'double', char: 'double' },
  double: { int: 'double', double: 'double', char: 'double' },
  char:   { int: 'double', double: 'double', char: 'X' }
};
// 5.5.5 Potencia
const T_POT = {
  int:    { int: 'int',    double: 'double' },
  double: { int: 'double', double: 'double' }
};
// 5.5.6 Raiz
const T_RAIZ = {
  int:    { int: 'double', double: 'double' },
  double: { int: 'double', double: 'double' }
};
// 5.5.6 Modulo  -> NOTA: el enunciado define Entero % Entero = DECIMAL
// (identico a la tabla de Raiz). Es inusual (uno esperaria Entero), pero
// se implementa TAL CUAL lo dice la fuente, sin "corregirlo".
const T_MOD = {
  int:    { int: 'double', double: 'double' },
  double: { int: 'double', double: 'double' }
};

const NOMBRE_OP = { '+': 'suma', '-': 'resta', '*': 'multiplicación',
  '/': 'división', '^': 'potencia', '$': 'raíz', '%': 'módulo' };

function nombreTipo(t) {
  return { int: 'ENTERO', double: 'DECIMAL', bool: 'BOOLEAN',
    char: 'CARÁCTER', string: 'CADENA', null: 'NULO', vector: 'VECTOR' }[t] || t;
}

function errBinaria(errores, op, ti, td, l, c) {
  errores.semantico('No se puede realizar la ' + NOMBRE_OP[op] + ' entre ' +
    nombreTipo(ti) + ' y ' + nombreTipo(td), l, c);
  return null;
}

// Empaqueta el numero como int o double segun el tipo resultante de la tabla.
function empacar(tipoRes, numero) {
  if (tipoRes === TIPO.INT) return Valor.int(numero);
  return Valor.double(numero);
}

function aritmetica(op, izq, der, l, c, errores) {
  if (izq === null || der === null) return null;   // propagacion por null
  let tabla;
  switch (op) {
    case '+': tabla = T_SUMA; break;
    case '-': tabla = T_RESTA; break;
    case '*': tabla = T_MULT; break;
    case '/': tabla = T_DIV; break;
    case '^': tabla = T_POT; break;
    case '$': tabla = T_RAIZ; break;
    case '%': tabla = T_MOD; break;
  }
  const fila = tabla[izq.tipo];
  const tipoRes = fila ? fila[der.tipo] : undefined;
  if (!tipoRes || tipoRes === 'X') return errBinaria(errores, op, izq.tipo, der.tipo, l, c);

  // caso string -> concatenacion
  if (tipoRes === TIPO.STRING) {
    return Valor.string(aTexto(izq) + aTexto(der));
  }
  const a = aNumero(izq), b = aNumero(der);
  let r;
  switch (op) {
    case '+': r = a + b; break;
    case '-': r = a - b; break;
    case '*': r = a * b; break;
    case '/':
      if (b === 0) { errores.semantico('División entre cero', l, c); return null; }
      r = a / b; break;
    case '^': r = Math.pow(a, b); break;
    case '$':
      if (b === 0) { errores.semantico('Índice de raíz igual a cero', l, c); return null; }
      r = Math.pow(a, 1 / b); break;
    case '%':
      if (b === 0) { errores.semantico('Módulo entre cero', l, c); return null; }
      r = a % b; break;
  }
  if (tipoRes === TIPO.INT) r = Math.trunc(r);
  return empacar(tipoRes, r);
}

// 5.5.7 Negacion unaria: solo int/double
function negacion(v, l, c, errores) {
  if (v === null) return null;
  if (v.tipo === TIPO.INT) return Valor.int(-v.valor);
  if (v.tipo === TIPO.DOUBLE) return Valor.double(-v.valor);
  errores.semantico('No se puede aplicar la negación unaria sobre ' + nombreTipo(v.tipo), l, c);
  return null;
}

/* ---------- Relacionales ----------
   DECISION DE DISENO (la matriz 5.6 del enunciado quedo con las celdas
   vacias al convertir el PDF; no es recuperable). Criterio adoptado y
   documentado aqui como decision propia, NO como dato del enunciado:
     - Igualdad (==, !=):
         * Nulo contra CUALQUIER tipo (solo verifica igualdad de nulidad).
         * int/double/char entre si (por valor numerico / codigo ASCII).
         * string SOLO con string.
         * bool SOLO con bool.
         * cualquier otra combinacion -> error semantico.
     - Orden (<, <=, >, >=):
         * int/double/char entre si (por valor / ASCII).
         * string con string (orden lexicografico).
         * cualquier otra combinacion (incluye bool y null) -> error semantico.
*/
function esNumComparable(t) { return t === TIPO.INT || t === TIPO.DOUBLE || t === TIPO.CHAR; }

function relacional(op, izq, der, l, c, errores) {
  if (izq === null || der === null) return null;
  const ti = izq.tipo, td = der.tipo;

  if (op === '==' || op === '!=') {
    let iguales;
    if (ti === TIPO.NULL || td === TIPO.NULL) {
      iguales = (ti === TIPO.NULL && td === TIPO.NULL);
    } else if (esNumComparable(ti) && esNumComparable(td)) {
      iguales = aNumero(izq) === aNumero(der);
    } else if (ti === TIPO.STRING && td === TIPO.STRING) {
      iguales = izq.valor === der.valor;
    } else if (ti === TIPO.BOOL && td === TIPO.BOOL) {
      iguales = izq.valor === der.valor;
    } else {
      errores.semantico('No se puede comparar (' + op + ') entre ' +
        nombreTipo(ti) + ' y ' + nombreTipo(td), l, c);
      return null;
    }
    return Valor.bool(op === '==' ? iguales : !iguales);
  }

  // operadores de orden
  let a, b;
  if (esNumComparable(ti) && esNumComparable(td)) {
    a = aNumero(izq); b = aNumero(der);
  } else if (ti === TIPO.STRING && td === TIPO.STRING) {
    a = izq.valor; b = der.valor;   // comparacion lexicografica de JS
  } else {
    errores.semantico('No se puede comparar (' + op + ') entre ' +
      nombreTipo(ti) + ' y ' + nombreTipo(td), l, c);
    return null;
  }
  let res;
  switch (op) {
    case '<': res = a < b; break;
    case '<=': res = a <= b; break;
    case '>': res = a > b; break;
    case '>=': res = a >= b; break;
  }
  return Valor.bool(res);
}

function logica(op, izq, der, l, c, errores) {
  if (op === '!') {
    if (izq === null) return null;
    if (izq.tipo !== TIPO.BOOL) {
      errores.semantico('El operador ! requiere un valor BOOLEAN, no ' + nombreTipo(izq.tipo), l, c);
      return null;
    }
    return Valor.bool(!izq.valor);
  }
  if (izq === null || der === null) return null;
  if (izq.tipo !== TIPO.BOOL || der.tipo !== TIPO.BOOL) {
    errores.semantico('El operador ' + op + ' requiere valores BOOLEAN, no ' +
      nombreTipo(izq.tipo) + ' y ' + nombreTipo(der.tipo), l, c);
    return null;
  }
  const r = op === '&&' ? (izq.valor && der.valor) : (izq.valor || der.valor);
  return Valor.bool(r);
}

module.exports = { aritmetica, negacion, relacional, logica, nombreTipo };
