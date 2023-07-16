package nodes

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mrhid6/dockerorc/utils"
)

type Node struct {
	Name   string   `json:"name"`
	Status bool     `json:"status"`
	IP     string   `json:"ip"`
	Tasks  []string `json:"tasks"`
}

type Nodes struct {
	Nodes []Node `json:"nodes"`
}

var nodes Nodes
var nodesConfigFileName = "nodes.json"

func InitNodes() {
	err := LoadNodesFromConfig()

	if err != nil {
		panic(err)
	}
}

func GetAllNodes() []Node {
	return nodes.Nodes
}

func LoadNodesFromConfig() error {

	err := utils.CreateFolder("configs")
	if err != nil {
		return err
	}

	configFilePath := filepath.Join("configs", nodesConfigFileName)

	if !utils.CheckFileExists(configFilePath) {
		file, err := os.Create(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	f, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byteValue, _ := io.ReadAll(f)

	json.Unmarshal(byteValue, &nodes)

	SaveNodesToConfig()

	return nil
}

func SaveNodesToConfig() {
	file, _ := json.MarshalIndent(nodes, "", "    ")

	configFilePath := filepath.Join("configs", nodesConfigFileName)
	err := os.WriteFile(configFilePath, file, 0755)

	if err != nil {
		panic(err)
	}

	absPath, _ := filepath.Abs(configFilePath)
	log.Printf("Saving nodes to config %s\r\n", absPath)
}

func HasNode(nodeName string) bool {
	for _, node := range nodes.Nodes {
		if node.Name == nodeName {
			return true
		}
	}

	return false
}

func RegisterNode(nodeName string, ip string) error {
	log.Printf("Registering node: %s\r\n", nodeName)
	if HasNode(nodeName) {
		log.Println("Skipped node already registered")
		return nil
	}

	newNode := Node{
		Name:   nodeName,
		IP:     ip,
		Status: true,
	}

	nodes.Nodes = append(nodes.Nodes, newNode)

	SaveNodesToConfig()

	return nil

}

func UpdateNodeStatus(nodeName string, status bool) error {

	foundNode := false
	for idx := range nodes.Nodes {
		node := &nodes.Nodes[idx]

		if node.Name == nodeName {
			node.Status = status
			foundNode = true
		}
	}

	if !foundNode {
		return errors.New("node was not found")
	}

	return nil

}

func API_RegisterNode(c *gin.Context) {
	type RegisterNodeBody struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	var dataBody RegisterNodeBody

	if err := c.BindJSON(&dataBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		c.Abort()
		return
	}

	if err := RegisterNode(dataBody.Name, dataBody.IP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func API_SetStatus(c *gin.Context) {
	type NodeStatusBody struct {
		Name   string `json:"name"`
		Status bool   `json:"status"`
	}

	var dataBody NodeStatusBody

	if err := c.BindJSON(&dataBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		c.Abort()
		return
	}

	if err := UpdateNodeStatus(dataBody.Name, dataBody.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
