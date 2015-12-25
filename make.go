package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	ffmpegFolderName := "ffmpeg-2.8.4"
	path, err := filepath.Abs("")
	if err != nil {
		fmt.Println("Error locating absulte file paths")
		os.Exit(1)
	}

	folderPath := filepath.Join(path, ffmpegFolderName)

	_, err = folderExists(folderPath)
	if err != nil {
		fmt.Println("The folder: %s either does not exist or is not in the same directory as make.go", folderPath)
		os.Exit(1)
	}
	// change the working directory
	err = os.Chdir(folderPath)
	if err != nil {
		fmt.Println("File Path Could not be changed")
		os.Exit(1)
	}

	var b bytes.Buffer
	fmt.Println("Configuring data and compiling data...please wait, this takes a while")
	if err := Execute(&b,
		//exec.Command("cd", folderPath),
		exec.Command("./configure", "--disable-yasm"),
		exec.Command("make"),
	); err != nil {
		log.Fatalln(err)
	}
	io.Copy(os.Stdout, &b)
	fmt.Println("Done!")

}

func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Execute(output_buffer *bytes.Buffer, stack ...*exec.Cmd) (err error) {
	var error_buffer bytes.Buffer
	pipe_stack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdin_pipe, stdout_pipe := io.Pipe()
		stack[i].Stdout = stdout_pipe
		stack[i].Stderr = &error_buffer
		stack[i+1].Stdin = stdin_pipe
		pipe_stack[i] = stdout_pipe
	}
	stack[i].Stdout = output_buffer
	stack[i].Stderr = &error_buffer

	if err := call(stack, pipe_stack); err != nil {
		log.Fatalln(string(error_buffer.Bytes()), err)
	}
	return err
}

func call(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = call(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}
