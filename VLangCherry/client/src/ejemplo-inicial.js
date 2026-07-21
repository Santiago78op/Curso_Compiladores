/* Ejemplo inicial del editor: cubre structs, metodos, slices 1D/2D,
   control de flujo y las 3 formas del for (4.7.3). */
export default `struct Persona {
    string Nombre
    int Edad
}

func (p *Persona) Saludar() string {
    return "Hola, soy " + p.Nombre
}

func suma(a int, b int) int {
    return a + b
}

func main() {
    mut i := 10
    println("Valor de i:", i)

    numeros := []int{10, 20, 30, 40, 50}
    print("Elemento en indice 2:", numeros[2])
    numeros[2] = 100

    mtx := [][]int{
        {1, 2, 3},
        {4, 5, 6},
    }
    print("mtx[0][1] =", mtx[0][1])

    p := Persona{
        Nombre: "Alice",
        Edad: 25,
    }
    print(p.Saludar())

    if i > 5 {
        print("i es mayor a 5")
    } else {
        print("i es 5 o menor")
    }

    switch i {
    case 1:
        print("uno")
    default:
        print("otro valor")
    }

    for k := 1; k <= 5; k++ {
        if k % 2 == 0 {
            continue
        }
        print("impar:", k)
    }

    for indice, valor in numeros {
        print(indice, valor)
    }

    resultado := suma(3, 7)
    print("suma(3,7) =", resultado)
}
`;
