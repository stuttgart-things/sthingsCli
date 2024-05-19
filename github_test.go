/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestGetGitTree_WithInvalidFiles(t *testing.T) {
// 	// CREATE A MOCK GITHUB.REFERENCE OBJECT
// 	ref := &github.Reference{
// 		Ref: github.String("refs/heads/master"),
// 		Object: &github.GitObject{
// 			SHA: github.String("6dcb09b5b57875f334f61aebed695e2e4193db5e"),
// 		},
// 	}

// 	// CREATE AN INVALID SOURCEFILES STRING
// 	sourceFiles := "go.sum:test/go.sum,go.mod:test/go.mod"
// 	sourceOwner := "patrick-hermann-sva"
// 	sourceRepo := "stuttgart-things"
// 	// CALL THE GETGITTREE FUNCTION
// 	tree, err := GetGitTree(ref, sourceFiles, sourceOwner, sourceRepo)

// 	// ASSERT THAT THE RETURNED GITHUB.TREE OBJECT IS NIL
// 	assert.Nil(t, tree)

// 	// ASSERT THAT THE RETURNED ERROR IS NOT NIL
// 	assert.NotNil(t, err)
// }

func TestReadFileToVar_WithValidFile(t *testing.T) {
	// CREATE A TEMPORARY FILE
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// WRITE SOME DATA TO THE FILE
	text := []byte("Hello, World!")
	if _, err := tmpfile.Write(text); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// CALL THE READFILETOVAR FUNCTION
	content, err := ReadFileToVar(tmpfile.Name())

	// ASSERT THAT THE RETURNED CONTENT IS THE SAME AS THE DATA WRITTEN TO THE FILE
	assert.Equal(t, text, content)

	// ASSERT THAT THE RETURNED ERROR IS NIL
	assert.Nil(t, err)
}

func TestReadFileToVar_WithNonExistentFile(t *testing.T) {
	// CALL THE READFILETOVAR FUNCTION WITH A NON-EXISTENT FILE PATH
	content, err := ReadFileToVar("non_existent_file.txt")

	// ASSERT THAT THE RETURNED CONTENT IS NIL
	assert.Nil(t, content)

	// ASSERT THAT THE RETURNED ERROR IS NOT NIL
	assert.NotNil(t, err)
}
