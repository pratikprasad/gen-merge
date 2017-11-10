package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"bytes"
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
