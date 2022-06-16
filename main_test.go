package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func wantOut() string {
	out := `Created parking of 6 slots
Car with vehicle registration number "KA-01-HH-1234" has been parked at slot number 1
Car with vehicle registration number "PB-01-HH-1234" has been parked at slot number 2
1,2
Car with vehicle registration number "PB-01-TG-2341" has been parked at slot number 3
2
Slot number 2 vacated, the car with vehicle registration number "PB-01-HH-1234" left the space, the driver of the car was of age 21
Car with vehicle registration number "HR-29-TG-3098" has been parked at slot number 2
Unknown input command
`
	return out
}

func Test_main(t *testing.T) {
	//Save old settings before rewriting settings
	oldArgs := os.Args
	oldInputInteractive := inputInteractive
	oldOutStream := outStream
	defer func() {
		os.Args = oldArgs
		inputInteractive = oldInputInteractive
		outStream = oldOutStream
	}()

	//Setup redirection for interactive inputs
	inputInteractiveFile, err := os.Open("inputInteractive.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { inputInteractiveFile.Close() }()
	inputInteractive = inputInteractiveFile

	//Setup redirection for outputs
	var gotBuf bytes.Buffer
	outStream = &gotBuf

	//Setup expected output
	wantBuf := bytes.NewBufferString(wantOut()).Bytes()

	tests := []struct {
		name string
		args []string
	}{
		{name: "File input",
			args: []string{"cmd", "inputFile.txt"},
		},
		{name: "Interactive input",
			args: []string{"cmd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			main()
			if !bytes.Equal(gotBuf.Bytes(), wantBuf) {
				t.Errorf("main() = %v, want = %v", gotBuf.String(), string(wantBuf))
			}
		})
		gotBuf.Reset()
	}
}