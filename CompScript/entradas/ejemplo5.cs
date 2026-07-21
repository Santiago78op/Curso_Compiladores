// Ejemplo 5: funciones recursivas, parametros por defecto y RUN_MAIN
// (prueba la pila de entornos por llamada)

int factorial(n: int) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n = n - 1);
}

int fib(n: int) {
    if (n < 2) {
        return n;
    }
    return fib(n = n - 1) + fib(n = n - 2);
}

// parametro con valor por defecto (5.21): se puede omitir en la llamada
int potencia(base: int, exp: int = 2) {
    return base ^ exp;
}

void main() {
    console.log(factorial(n = 5));    // 120
    console.log(fib(n = 10));         // 55
    console.log(potencia(base = 4));      // 16 (usa exp por defecto = 2)
    console.log(potencia(base = 2, exp = 5)); // 32
}

RUN_MAIN main();
