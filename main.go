package main

//Paquetes importados
import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//Nombre para los datos de los héries
//Estructuras para poder leer el archivo JSON
type Nombre struct {
	ID          int    `json:"id"`
	Nombre      string `json:"name"`
	Descripcion string `json:"description"`
}

//Resultado yHéroe se crearon porque se debe acceder a ellos antes de los datos de los superhéroes
//El tipo Resultado consiste de un array del tipo Nombre
type Resultado struct {
	Resultado []Nombre `json:"results"`
}

//El tipo héroe contiene los demás tipos
type Heroe struct {
	Data Resultado `json:"data"`
}

//Constantes con las llaves públicas y privadas
const llavePrivada, llavePublica = "f16a809c4072d7d075fc8df3d03773734d82750a", "477476fe99e4e2eff346fcd265b8b894"

func main() {
	nombre := ""
	fmt.Println("Bienvenido a este buscador de superhéroes")
	fmt.Println("Tenemos dos opciones")
	fmt.Println("1. Buscar por nombre")
	fmt.Println("2. Listar los 20 primeros héroes")
	ok := true
	//Se usa un loop para permitir solo una opción válida
	for ok {
		//Se registra la opción que elija el usuario
		fmt.Println("Elija una opción:")
		reader := bufio.NewReader(os.Stdin)
		opcion, _ := reader.ReadString('\n')
		opcion = strings.TrimRight(opcion, "\r\n")
		switch opcion {
		case "2":
			ok = false
		case "1":
			//Si es la opción dos se recibe el nombre
			fmt.Print("Escriba el nombre del superhéroe que desea buscar:")
			reader = bufio.NewReader(os.Stdin)
			nombre, _ = reader.ReadString('\n')
			nombre = strings.TrimRight(nombre, "\r\n")
			ok = false
		case "":
			fmt.Println("No se permiten vacíos")
			fmt.Println("Ingrese un valor válido [1 o 2]")
		default:
			fmt.Println("Ingrese un valor válido [1 o 2]")
		}
	}
	now := time.Now()
	//Se obtiene un entero como en JS
	tiempo := now.UnixNano()
	//Se convierte el entero en sistema decimal a una cadena
	ts := strconv.FormatInt(tiempo, 10)
	//Se inicializa un objeto md5 para encriptar la cadena
	str := md5.New()
	//Se asigna a str la cadena sin encriptar
	io.WriteString(str, ts+llavePrivada+llavePublica)
	//Se transforma a texto el hash
	hash := hex.EncodeToString(str.Sum(nil))
	//Se eige un URL según la opción de consola
	URL := "http://gateway.marvel.com/v1/public/characters?ts=" + ts + "&apikey=" + llavePublica + "&hash=" + hash
	if nombre != "" {
		//Se parsea por si el nombre contiene espacios, ara que la API entienda el nombre
		nombrecodificado := &url.URL{Path: nombre}
		nombre := nombrecodificado.String()
		URL = "http://gateway.marvel.com/v1/public/characters?name=" + nombre + "&ts=" + ts + "&apikey=" + llavePublica + "&hash=" + hash
	}
	//Se crea un cliente http que dura tres segundos
	spaceClient := http.Client{
		Timeout: time.Second * 3,
	}
	//Se crea un nuevo Request con el método y URL, y se verifica si hay error
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	//Para que el API entienda el tipo de petición que se está haciendo
	req.Header.Set("User-Agent", "Cliente para API de Marvel")

	//El cliente Cliente http realiza el Request
	//Mostrará el error en consola si se presentases
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	//Se leen los resultados con un formato de []byte
	//Mostrará el error en consola si se presentase
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	//Se inicializa la variable heroes es del tipo Heroe para almacenar los resaultados
	heroes := Heroe{}
	//Se generará el error si no se puede asignar el resultado a la variable heroes
	jsonErr := json.Unmarshal(body, &heroes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	//Si no hay resultado pero se asigna un valor vacío a la variable heroes
	if len(heroes.Data.Resultado) == 0 {
		fmt.Println("Nombre incorrecto de superhéroe")
	} else {
		//Se recorre la variable heroes y se muestran los datos obtenidos según la estructura del
		//tipo Heroe, se convierte los int a cadenas
		fmt.Println("--Héroes ordenados alfabéticamente--")
		fmt.Println()
		for i := 0; i < len(heroes.Data.Resultado); i++ {
			fmt.Println("Héroe número", i+1)
			fmt.Println("Nombre de superhéroe: " + heroes.Data.Resultado[i].Nombre)
			fmt.Println("ID de superhéroe: " + strconv.Itoa(heroes.Data.Resultado[i].ID))
			if heroes.Data.Resultado[i].Descripcion == "" {
				fmt.Println("Descripción: No tiene descripcion")
			} else {
				fmt.Println("Descripción: " + heroes.Data.Resultado[i].Descripcion)
			}
			fmt.Println()

		}
		//Se elije el mensaje final según la opción elegida
		fmt.Println("--Fin de la lista--")
		if nombre != "" {
			fmt.Println("--Se mostró el superhéroe que buscó--")
		} else {
			fmt.Println("--Se mostraron los primeros 20 superhéroes--")
		}
		fmt.Println()
	}
	salir()
}

//Salir espera que se presiona una tecla para acabar la ejecución del programa
func salir() {
	fmt.Println("Presione cualquier tecla para salir")
	reader := bufio.NewReader(os.Stdin)
	opcion, _ := reader.ReadString('\n')
	switch opcion {
	case "\r\n":
		break
	default:
		break
	}

}
