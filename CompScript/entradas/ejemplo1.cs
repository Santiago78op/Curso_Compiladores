// Ejemplo 1: tipos de dato, expresiones, casteos y console.log (fase basica)
let a: int = 10;
let b: int = 3;
let c: double = 2.5;
let flag: bool = true;
let letra: char = 'F';
let saludo: string = "Hola" + " " + "Mundo";

console.log(saludo);              // Hola Mundo
console.log(a + b);               // 13
console.log(a - b);               // 7
console.log(a * b);               // 30
console.log(a / b);               // division siempre Decimal
console.log(a % b);               // 1.0
console.log(a ^ b);               // 1000  (potencia)
console.log(9 $ 2);               // 3.0   (raiz cuadrada de 9)
console.log(letra);               // F
console.log(cast(letra as int));  // 70
console.log(cast(70 as char));    // F
console.log(cast(18.6 as int));   // 18
console.log(cast(16 as double));  // 16.0
console.log(flag && (a > b));     // true
console.log(!flag);               // false
console.log(-a);                  // -10  (negacion unaria)
console.log('a' + 'b');           // ab   (char + char -> cadena)
