package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileGen(t *testing.T) {
	pkg, err := ParsePackage("teststructs")
	assert.Nil(t, err)

	testGeneratedFile, err := ioutil.ReadFile("teststructs/genmerge_teststructs.go")
	assert.Nil(t, err)

	buf := bytes.NewBufferString("")
	err = PrintMergePackage(buf, pkg)
	assert.Nil(t, err)

	assert.Equal(t, string(testGeneratedFile), buf.String())
}
