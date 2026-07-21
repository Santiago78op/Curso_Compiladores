// Ejemplo 3: ciclos while / for / do-while con break y continue
void demo() {
    // while
    let i: int = 0;
    while (i < 3) {
        console.log(i);          // 0 1 2
        i++;
    }

    // for con continue (salta i==2) y break (corta en i==4)
    for (i = 0; i < 10; i++) {
        if (i == 2) { continue; }
        if (i == 4) { break; }
        console.log(i * 10);     // 0 10 30
    }

    // do-while: se ejecuta al menos una vez, repite MIENTRAS la condicion sea true
    let k: int = 0;
    do {
        console.log(k);          // 0 1
        k++;
    } while (k < 2);
}

RUN_MAIN demo();
