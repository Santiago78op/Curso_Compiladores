// Ejemplo 2: sentencia match (funciona como if sucesivos, sin break ni fall-through)
string clasificar(n: int) {
    let r: string = "";
    match n {
        1 => { r = "uno"; }
        2 => { r = "dos"; }
        3 => { r = "tres"; }
        default => { r = "otro numero"; }
    }
    return r;
}

void main() {
    console.log(clasificar(n = 1));   // uno
    console.log(clasificar(n = 3));   // tres
    console.log(clasificar(n = 9));   // otro numero

    // match solo con default
    let x: int = 100;
    match x {
        default => { console.log("cayo en default"); }
    }
}

RUN_MAIN main();
