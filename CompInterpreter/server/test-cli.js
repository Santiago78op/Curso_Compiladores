/* Prueba por consola: corre un .ci por el pipeline y muestra el resultado.
   Uso: node test-cli.js ../entradas/ejemplo_anexo.ci */
const fs = require('fs');
const path = require('path');
const { analizar } = require('./src/analizar');

const archivo = process.argv[2] || path.join(__dirname, '..', 'entradas', 'ejemplo_anexo.ci');
const codigo = fs.readFileSync(archivo, 'utf8');
const r = analizar(codigo);

console.log('==== ARCHIVO:', archivo, '====');
console.log('\n--- CONSOLA ---\n' + (r.consola || '(vacia)'));
console.log('\n--- ERRORES (' + r.errores.length + ') ---');
r.errores.forEach(e => console.log(`  [${e.tipo}] L${e.linea}:C${e.columna} ${e.descripcion}`));
console.log('\n--- SIMBOLOS (' + r.simbolos.length + ') ---');
r.simbolos.forEach(s => console.log(`  ${s.id} | ${s.categoria} | ${s.tipoDato} | ${s.entorno} | ${s.valor}`));
console.log('\n--- AST ---  nodos:', r.ast.nodes.length, ' aristas:', r.ast.edges.length);
