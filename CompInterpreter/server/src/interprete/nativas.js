/* ============================================================
   Funciones nativas (secciones 5.22 - 5.25).
   Reciben un Valor ya evaluado y devuelven un Valor (o null si error).
   ============================================================ */
const { TIPO } = require('./tipos');
const { Valor, aNumero, aTexto } = require('./valor');
const { nombreTipo } = require('./operaciones');

function err(errores, msg, l, c) { errores.semantico(msg, l, c); return null; }

// M4: max/min/sum/average sobre un vector con elementos SIN asignar (null, como
// los que crea "new vector int[n]") producirian NaN en silencio (aNumero(null)
// es NaN). En su lugar se reporta un error semantico.
function tieneNulos(vec) { return vec.some(e => !e || e.tipo === TIPO.NULL); }

function ejecutarNativa(fn, arg, l, c, errores) {
  if (arg === null) return null;

  switch (fn) {
    case 'lower':
      if (arg.tipo !== TIPO.STRING) return err(errores, 'lower espera una CADENA, no ' + nombreTipo(arg.tipo), l, c);
      return Valor.string(arg.valor.toLowerCase());

    case 'upper':
      if (arg.tipo !== TIPO.STRING) return err(errores, 'upper espera una CADENA, no ' + nombreTipo(arg.tipo), l, c);
      return Valor.string(arg.valor.toUpperCase());

    case 'round':
      if (arg.tipo !== TIPO.INT && arg.tipo !== TIPO.DOUBLE)
        return err(errores, 'round espera un valor numérico, no ' + nombreTipo(arg.tipo), l, c);
      // >= .5 sube, < .5 baja (Math.round hace exactamente eso para el .5)
      return Valor.int(Math.round(arg.valor));

    case 'truncate':
      if (arg.tipo !== TIPO.INT && arg.tipo !== TIPO.DOUBLE)
        return err(errores, 'truncate espera un valor numérico, no ' + nombreTipo(arg.tipo), l, c);
      return Valor.int(Math.trunc(arg.valor));

    case 'len':
      if (arg.tipo === TIPO.STRING) return Valor.int(arg.valor.length);
      if (arg.tipo === TIPO.VECTOR) return Valor.int(arg.valor.length);
      return err(errores, 'len espera una CADENA o VECTOR, no ' + nombreTipo(arg.tipo), l, c);

    case 'tostring':
      if (arg.tipo === TIPO.INT || arg.tipo === TIPO.DOUBLE || arg.tipo === TIPO.BOOL)
        return Valor.string(aTexto(arg));
      return err(errores, 'toString espera numérico o boolean, no ' + nombreTipo(arg.tipo), l, c);

    case 'tochararray': {
      if (arg.tipo !== TIPO.STRING) return err(errores, 'toCharArray espera una CADENA, no ' + nombreTipo(arg.tipo), l, c);
      const chars = arg.valor.split('').map(ch => Valor.char(ch));
      return Valor.vector(TIPO.CHAR, 1, chars);
    }

    case 'reverse': {
      if (arg.tipo !== TIPO.VECTOR || arg.dimension !== 1)
        return err(errores, 'reverse espera un VECTOR de una dimensión', l, c);
      const copia = arg.valor.slice().reverse();
      return Valor.vector(arg.tipoBase, 1, copia);
    }

    case 'max':
    case 'min': {
      if (arg.tipo !== TIPO.VECTOR || arg.dimension !== 1)
        return err(errores, fn + ' espera un VECTOR de una dimensión', l, c);
      if (arg.valor.length === 0) return err(errores, fn + ' sobre un vector vacío', l, c);
      if (tieneNulos(arg.valor)) return err(errores, fn + ' no puede operar sobre un vector con elementos nulos sin asignar', l, c);
      const base = arg.tipoBase;
      const buscarMax = fn === 'max';
      let mejor = arg.valor[0];
      for (let i = 1; i < arg.valor.length; i++) {
        const e = arg.valor[i];
        let comp;
        if (base === TIPO.STRING) comp = e.valor < mejor.valor ? -1 : (e.valor > mejor.valor ? 1 : 0);
        else comp = aNumero(e) - aNumero(mejor);
        if ((buscarMax && comp > 0) || (!buscarMax && comp < 0)) mejor = e;
      }
      return mejor;
    }

    case 'sum': {
      if (arg.tipo !== TIPO.VECTOR || arg.dimension !== 1)
        return err(errores, 'sum espera un VECTOR de una dimensión', l, c);
      if (tieneNulos(arg.valor)) return err(errores, 'sum no puede operar sobre un vector con elementos nulos sin asignar', l, c);
      const base = arg.tipoBase;
      if (base === TIPO.STRING) {
        return Valor.string(arg.valor.map(v => v.valor).join(''));
      }
      let acc = 0;
      for (const e of arg.valor) acc += aNumero(e);
      if (base === TIPO.DOUBLE) return Valor.double(acc);
      return Valor.int(acc);   // int, char(ASCII) y bool -> entero
    }

    case 'average': {
      if (arg.tipo !== TIPO.VECTOR || arg.dimension !== 1)
        return err(errores, 'average espera un VECTOR de una dimensión', l, c);
      if (arg.tipoBase === TIPO.STRING)
        return err(errores, 'average no es aplicable a un VECTOR de CADENA', l, c);
      if (arg.valor.length === 0) return err(errores, 'average sobre un vector vacío', l, c);
      if (tieneNulos(arg.valor)) return err(errores, 'average no puede operar sobre un vector con elementos nulos sin asignar', l, c);
      let acc = 0;
      for (const e of arg.valor) acc += aNumero(e);
      return Valor.double(acc / arg.valor.length);
    }
  }
  return err(errores, 'Función nativa desconocida: ' + fn, l, c);
}

module.exports = { ejecutarNativa };
