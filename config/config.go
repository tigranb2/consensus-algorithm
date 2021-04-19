package config
import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

func readConfig() (map[int]string, string){
	line := 0
	nodeCounts := ""
	file, err := os.Open("config.json") //opens the file
	if err != nil {
		fmt.Println(err) //handles error
	}

	defer file.Close()

	scanner := bufio.NewScanner(file) //new scanner for scanning file word by word

	nodes := make(map[int]string)
	for scanner.Scan() {
		if line > 0 { //scans id info
			nodes[line] = scanner.Text() //each line is an individual id; uses id as key
		}  else {
			nodeCounts = (scanner.Text())
		}
		
		line++
	}
	return nodes, nodeCounts //returns scaned lines, not formatted
}

var nodesInfo, nodeCounts = readConfig()

//formats node count and fault count
func defineCounts(nodeCounts string) (int, int) {
	nc := strings.Fields(nodeCounts)
	nodeCount, _ := strconv.Atoi(nc[0])
	faultCount, _ := strconv.Atoi(nc[1])
	return nodeCount, faultCount
}

var NodeCount, FaultCount = defineCounts(nodeCounts)

//formats ip:port for nodes
func defineNodes(nodesInfo map[int]string) map[int]string {
	nodesCONNECT := make(map[int]string)
	for node, info := range nodesInfo {
		nI := strings.Fields(info)
		nodesCONNECT[node] = nI[1]+":"+nI[2]

		if len(nodesCONNECT) == NodeCount {
			break
		}
	}

	return nodesCONNECT
}

var NodesCONNECT = defineNodes(nodesInfo)