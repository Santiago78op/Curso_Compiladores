// Ejemplo 4: vectores 1D/2D, listas dinamicas y structs
struct persona {
    nombre: string;
    edad: int;
};

void run() {
    // vector 1D
    let v: int[] = [10, 20, 30];
    console.log(v);              // [10, 20, 30]
    console.log(v[1]);           // 20
    v[1] = 99;
    console.log(v);              // [10, 99, 30]

    // vector 2D
    let m: int[][] = [[1, 2], [3, 4]];
    console.log(m[1][0]);        // 3
    m[0][1] = 200;
    console.log(m[0][1]);        // 200

    // lista dinamica
    let lista: List<int>;
    lista.push(5);
    lista.push(7);
    lista.push(9);
    console.log(lista);          // [5, 7, 9]
    console.log(lista.get(1));   // 7
    lista.set(1, 70);
    console.log(lista.remove(0));// 5  (elimina y retorna)
    console.log(lista);          // [70, 9]
    lista.reverse();
    console.log(lista);          // [9, 70]
    console.log(length(lista));  // 2

    // struct: instanciacion por nombre (el orden no importa)
    let p: persona = { edad: 30, nombre: "Ana" };
    console.log(p.nombre);       // Ana
    p.edad = 31;
    console.log(toString(p));    // persona { nombre: "Ana", edad: 31 }
}

RUN_MAIN run();
