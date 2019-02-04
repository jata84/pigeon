package core

import (
	"fmt"
	"testing"
)

func TestCreateMultipleNamespaceNode(t *testing.T) {
	nodeTree := NewNodeTree()

	nodeTree.NewNamespace(stringToNamespaces("+prueba.test.cuentas"))
	nodeTree.NewNamespace(stringToNamespaces("-prueba.test.mensajes"))
	nodeTree.NewNamespace(stringToNamespaces("+prueba.test.alertas"))
	fmt.Println(nodeTree.ToJson())
}

func TestCreateMultiplesNamespaceNode(t *testing.T) {
	nodeTree := NewNodeTree()

	nodeTree.NewNamespace(stringToNamespaces("-prueba.test.cuentas"))
	nodeTree.NewNamespace(stringToNamespaces("-prueba.test.mensajes"))
	nodeTree.NewNamespace(stringToNamespaces("-prueba.test.alertas"))
	nodeTree.NewNamespace(stringToNamespaces("+prueba.publico"))
	nodeTree.NewNamespace(stringToNamespaces("+prueba.test.pruebas.mensajes.publicos"))
	nodeTree.NewNamespace(stringToNamespaces("-prueba.test.pruebas.mensajes.privados"))
	nodeTree.NewNamespace(stringToNamespaces("-prueba.test.pruebas.mensajes.privados.test"))

	fmt.Println(nodeTree.ToJson())

}
