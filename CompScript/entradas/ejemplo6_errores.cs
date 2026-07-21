// Ejemplo 6: errores a proposito (los tres tipos aparecen en el reporte)

// 1) ERROR LEXICO: el caracter '#' no pertenece al lenguaje (se descarta)
let x: int = 5;
#

// 2) ERROR SINTACTICO: falta la expresion despues del '=' (modo panico -> ';')
let w: int = ;

// 3) ERROR SEMANTICO: no se puede restar dos cadenas (aborta la ejecucion)
let z: string = "a" - "b";
console.log(z);   // no llega a ejecutarse: la ejecucion termino en el error semantico
