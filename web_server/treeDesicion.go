package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

var number int
var num = 0
var s = 1

func readDataFunc(pathFile string, columns int) [][]string {

	dataRead := [][]string{}

	file, err := os.Open(pathFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = columns
	reader.Comment = '#'

	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}
		dataRead = append(dataRead, record)
	}
	return dataRead
}

//Esta  funcion encuentra los unicos valores de una matriz en un arreglo de string
func unique(data [][]string, nColumn int) []string {
	var column = []string{}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[0]); j++ {
			column = append(column, data[i][nColumn])
		}
	}

	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range column {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

//retorna un boleano de ser encontrado o no un elemento en un array
func elementNotInArray(element string, array []string) bool {
	for i := 0; i < len(array); i++ {
		if element == array[i] {
			return false
		}
	}
	return true
}

//Cuenta de acuerdo a la matriz cuantos valores repetidos hay y los devuelve en un map
func classCounts(data [][]string) map[string]int {
	mapCounts := make(map[string]int)
	var long = len(data[0]) - 1
	var uniques = unique(data, long)
	for i := 0; i < len(uniques); i++ {
		mapCounts[uniques[i]] = 0
	}
	for i := 0; i < len(data); i++ {
		var label = data[i][long]
		for j := 0; j < len(uniques); j++ {
			if label == uniques[j] {
				mapCounts[label]++
			}
		}
	}
	return mapCounts
}

// Se crea una structura llamada question para recrear las preguntas a realizar para el nodo de decision
type Question struct {
	value  string
	column int
}

// Se realiza la pregunta, de acuerdo a la pregunta por cada elemento de un arreglo
func match(example []string, question Question) bool {
	var val = example[question.column]
	if n, err := strconv.Atoi(question.value); err == nil {
		if n2, err := strconv.Atoi(val); err == nil {
			if n2 >= n {
				return true
			}
		}
	} else {
		if val == question.value {
			return true
		}
	}

	return false
}

// de acuerdo a la data en una matriz y la pregunta se realiza la particion entre las respuestas de Verdadero o falso en dos distintos arreglos
func partition(data [][]string, question Question) ([][]string, [][]string) {
	trueRows := [][]string{}
	falseRows := [][]string{}
	for i := 0; i < len(data); i++ {
		if match(data[i], question) {
			trueRows = append(trueRows, data[i])
		} else {
			falseRows = append(falseRows, data[i])
		}
	}
	return trueRows, falseRows
}

// Se halla la impureza o gini para elegir posteriormente la mejor opcion
func gini(data [][]string) float64 {
	counts := classCounts(data)
	var impurity = 1.0
	for j, _ := range counts {
		probOfLbl := float64(counts[j]) / float64(len(data))
		impurity -= math.Pow(probOfLbl, 2.0)
	}

	return impurity
}

// Retorna el ponderado de impureza entre los dos arreglos left y rigth o arreglo de true y false
func infoGain(left [][]string, rigth [][]string, currentUnCertainty float64) float64 {
	longLeft := float64(len(left))
	longRigth := float64(len(rigth))
	p := longLeft / (longLeft + longRigth)
	info := (currentUnCertainty - p*gini(left) - (1-p)*gini(rigth))
	return info
}

// Busca la mejor pregunta de acuerdo a una data y retorna la pregunta encontrada.
func findBestSplit(data [][]string) (float64, Question) {
	bestGain := 0.0
	bestQuestion := Question{value: "none"}

	currentUnCertainty := gini(data)
	nFeatures := len(data[0]) - 1

	for col := 0; col < nFeatures; col++ {
		values := unique(data, col)
		for _, val := range values {
			question := Question{value: val, column: col}

			trueRows, falseRows := partition(data, question)

			if len(trueRows) == 0 || len(falseRows) == 0 {
				continue
			}

			gain := infoGain(trueRows, falseRows, currentUnCertainty)
			if gain >= bestGain {
				bestGain, bestQuestion = gain, question

			}
		}
	}
	return bestGain, bestQuestion
}

func leafConstructor(data [][]string) map[string]int {
	return classCounts(data)
}

// Funcion para convertir un map en un arreglo bidimensional
func leafConvertArrBi(m map[string]int) [][]string {
	pairs := [][]string{}
	for key, value := range m {
		pairs = append(pairs, []string{key, strconv.Itoa(value)})
	}
	return pairs
}
func leafConvertArrBiString(m map[string]string) [][]string {
	pairs := [][]string{}
	for key, value := range m {
		pairs = append(pairs, []string{key, value})
	}
	return pairs
}

// La estructura del arbol la inclui en la misma estructura de Nodos de decision donde el nodo Raiz tiene como atributos la mejor pregunta
// con dos arreglos de falsos y verdaderos de dicha pregunta
// Posteriormente los atributos isNotDecision y predictions los inclui para trabajar con solo un tipo de dato
// Se evaluara en cada caso si el nodo es o no un de decision o si solo sera una hoja y tendra predicctions como arreglo de las predicciones de Leaf
type DecisionNode struct {
	question      Question
	trueBranch    []DecisionNode
	falseBranch   []DecisionNode
	isNotDecision bool
	predictions   [][]string
}

// Este funcion retornara un nodo de decision de forma recursiva una vez se encuentre cada hoja del arbol
// primero busca el gain o grado de impureza de la data y la mejor pregunta a realizar para la particion
// despues pregunta si el gain es 0 esto significaria que la impuresa es TOTALMENTE pura y para lo cual retorna un leaf
// de no ser asi realiza la particion en dos columnas que simularan las ramas para llamarse recursivamente hasta obtener las hojas
// luego retornara sucesivamente un nodo de decision
func buildTree(data [][]string) DecisionNode {
	gain, questionBT := findBestSplit(data)
	true_branch := []DecisionNode{}
	false_branch := []DecisionNode{}
	if gain == 0 {
		mapLeaf := leafConstructor(data)
		return DecisionNode{isNotDecision: true, predictions: leafConvertArrBi(mapLeaf)}
	}
	true_rows, false_rows := partition(data, questionBT)

	true_branch = append(true_branch, buildTree(true_rows))

	false_branch = append(false_branch, buildTree(false_rows))

	return DecisionNode{question: questionBT, trueBranch: true_branch, falseBranch: false_branch, isNotDecision: false}
}

// Esta funcion realiza la misma accion de la anterior solo que con algunos cambios
// Al buscar cada branch de true o false implemente dos gorutimes
// Para que no exista un conflicto al que termine la funcion antesde que acaben los procesos
// añadi un grupo de gorutimes y este espere su termino para retornar el valor.
func buildTreeFor(data [][]string) DecisionNode {
	var wg sync.WaitGroup

	gain, questionBT := findBestSplit(data)
	true_branch := []DecisionNode{}
	false_branch := []DecisionNode{}
	if gain == 0 {
		mapLeaf := leafConstructor(data)
		return DecisionNode{isNotDecision: true, predictions: leafConvertArrBi(mapLeaf)}
	}
	true_rows, false_rows := partition(data, questionBT)

	wg.Add(2)
	go paralem(&wg, &true_branch, &true_rows)
	go paralem(&wg, &false_branch, &false_rows)

	wg.Wait()

	return DecisionNode{question: questionBT, trueBranch: true_branch, falseBranch: false_branch, isNotDecision: false}
}

//Mediante esta funcion cada proceso funcionara y almacenara el valor con el cual se podra determinar si el proceso ya acabo o no
//Implementar esta funcion incremento la rapidez en crear el arbol pero podria crear un sobrecosto al crear tantos procesos
func paralem(wg *sync.WaitGroup, branch *[]DecisionNode, rows *[][]string) {
	defer wg.Done()
	*branch = append(*branch, buildTreeFor(*rows))

}

// Esta funcion solo permite ver como se va creando el arbol
func printTree(node DecisionNode, spacing string) {
	if node.isNotDecision {
		fmt.Println(spacing, "Predic", node.predictions)
		return
	}

	fmt.Println(spacing, node.question)

	fmt.Println(spacing, "--> True: ")
	printTree(node.trueBranch[0], spacing+"  ")

	fmt.Println(spacing, "--> False: ")
	printTree(node.falseBranch[0], spacing+"  ")
}

// Esta funcion retorna una matriz con la respuesta final de la ultima columna que seria la columna de Diagnostico con la respuesta
// de yes or no y el porcentaje de prediccion de cada una de las dos respuestas
func classify(row []string, node DecisionNode) [][]string {
	if node.isNotDecision {
		return node.predictions
	}

	if match(row, node.question) {
		return classify(row, node.trueBranch[0])
	} else {
		return classify(row, node.falseBranch[0])
	}
}

func printLeaf(counts [][]string) [][]string {

	total := 0.0
	var props = make(map[string]string)
	var countMapt = make(map[string]int)
	for _, row := range counts {
		countMapt[row[0]], _ = strconv.Atoi(row[1])
	}
	for i := range countMapt {
		total = float64(countMapt[i]) + total
	}
	for lbl := range countMapt {
		a := countMapt[lbl]
		c := float64((float64(a) / float64(total)) * 100.0)
		b := strconv.FormatFloat(c, 'f', 1, 64) + "%"
		props[lbl] = b
	}
	return leafConvertArrBiString(props)
}

func startTree() {
	// La data se lee con la funcion readDataFunc con los parametos del nombre del archivo y la cantidad de columnas
	// Si se desea cambiar de archivo se debera poner la misma cantidad de columnas
	var data = readDataFunc("datatest02Covid.csv", 21)
	var startTime = time.Now()

	// Para testing seria un total de 4 columnas
	//var data = readDataFunc("testing.csv", 4)
	//var startTime = time.Now()

	myTree := buildTreeFor(data)

	// Calculamos la duracion que toma en realizar el arbol
	var duration = time.Since(startTime)

	// este for recorre toda la ultima columna y escribiendo el porcentade de prediccion para cada pregunta "Yes or No"
	for _, row := range data {
		fmt.Println("Actual: ", (row[len(data[0])-1]), " Predicted: ", printLeaf(classify(row, myTree)))
	}

	printTree(myTree, "")
	fmt.Println(duration)
}
