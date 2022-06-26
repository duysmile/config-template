package actions

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/vnteamopen/config-template/transform"
	"io"
	"os"
)

func Merge(inputPath string, listOutputPath []string) error {
	// printInfo(inputPath, outputPath)
	if _, err := os.Stat(inputPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("input '%s' doesn't exist", inputPath)
	}
	outputContent, err := parse(inputPath)
	if err != nil {
		return err
	}

	for _, outputPath := range listOutputPath {
		if err := write(outputPath, outputContent); err != nil {
			return err
		}
	}
	return nil
}

func printInfo(templatePath, outputPath string) {
	fmt.Printf("* Template path: %s\n* Output path: %s\n", templatePath, outputPath)
}

func write(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func CharByCharMerge(inputPath string, outputPath string) error {
	if _, err := os.Stat(inputPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("input '%s' doesn't exist", inputPath)
	}

	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file '%s': %s", inputPath, err.Error())
	}
	defer func() {
		if err := input.Close(); err != nil {
			fmt.Printf("cannot close file '%s': %s", inputPath, err.Error())
		}
	}()
	in := bufio.NewReader(input)

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file '%s': %s", outputPath, err.Error())
	}
	defer func() {
		if err := output.Close(); err != nil {
			fmt.Printf("cannot close file '%s': %s", outputPath, err.Error())
		}
	}()
	out := bufio.NewWriter(output)

	buf := make([]byte, 1)
	transformer := transform.NewTransformer()
	for {
		n, err := in.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("cannot read file '%s': %s", inputPath, err.Error())
		}

		if err == io.EOF {
			break
		}

		if n == 0 {
			continue
		}

		transformerBuf, err := transformer.Transform(buf[0])
		if err != nil {
			return err
		}
		if transformerBuf == nil {
			continue
		}
		if _, err := out.Write(transformerBuf); err != nil {
			return fmt.Errorf("cannot write file '%s': %s", outputPath, err.Error())
		}
	}

	if err := out.Flush(); err != nil {
		return fmt.Errorf("cannot write file '%s': %s", outputPath, err.Error())
	}

	return nil
}
