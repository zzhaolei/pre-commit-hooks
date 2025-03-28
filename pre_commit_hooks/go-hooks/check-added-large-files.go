package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	maxMib float64
	hint   bool
)

var checkAddedLargeFilesCmd = &cobra.Command{
	Use:   "check-added-large-files",
	Short: "prevents giant files from being committed. Like `jujutsu`.",
	Run: func(cmd *cobra.Command, args []string) {
		if maxMib == 0 {
			os.Exit(0)
		}

		ret := findLargeAddedFiles()
		os.Exit(ret)
	},
}

func init() {
	checkAddedLargeFilesCmd.Flags().Float64Var(&maxMib, "max-mib", 1.0, "Maximum allowable MiB for added files")
	checkAddedLargeFilesCmd.Flags().BoolVar(&hint, "hint", false, "Whether to display hint")

	rootCmd.AddCommand(checkAddedLargeFilesCmd)
}

type RefusedFile struct {
	filename string
	mib      float64
	kib      float64
	bytes    float64
}

func findLargeAddedFiles() int {
	stagedFiles, err := getStagedFiles()
	if err != nil {
		fmt.Println(err)
	}

	filenamesFiltered := make(map[string]struct{}, len(stagedFiles))
	for _, s := range stagedFiles {
		filenamesFiltered[s] = struct{}{}
	}

	var ret int
	refusedFiles := make([]RefusedFile, 0, len(filenamesFiltered))
	maxKib := maxMib * 1024
	for filename := range filenamesFiltered {
		info, err := os.Stat(filename)
		if err != nil {
			fmt.Printf("Error accessing file %s: %v\n", filename, err)
			continue
		}

		size := float64(info.Size())
		kib := math.Ceil(size / 1024)
		if kib > maxKib {
			refusedFiles = append(refusedFiles, RefusedFile{
				filename: filename,
				mib:      kib / 1024,
				kib:      kib,
				bytes:    size,
			})
			ret = 1
		}
	}
	if len(refusedFiles) == 0 {
		return ret
	}

	fmt.Printf("%sWarning:%s ", yellowB, reset)
	fmt.Printf("%sRefused to commit some files:%s\n", bold, reset)

	var refusedMaxMib float64
	for _, r := range refusedFiles {
		refusedMaxMib = math.Max(refusedMaxMib, r.mib)

		mib := fmt.Sprintf("%s%.1f%s", redB, r.mib, reset)
		_maxMib := fmt.Sprintf("%s%.1f%s", greenB, maxMib, reset)
		fmt.Printf("  %s%s%s: %sMiB; the maximum size allowed is %sMiB\n", bold, r.filename, reset, mib, _maxMib)
	}
	if hint {
		fmt.Printf("%sHint:%s This is to prevent large files from being added by accident. You can fix this by:\n", cyanB, reset)
		fmt.Println("  - Adding the file to `.gitignore`")
		fmt.Printf("  - Change `.pre-commit-config.yaml` args `--max-mib=%.1f`\n", refusedMaxMib)
		fmt.Println("    This will increase the maximum file size allowed for new files, in this repository only.")
		fmt.Println("  - Run `git commit --no-verify`")
		fmt.Println("    This will not run pre-commit and commit-msg hooks.")
	}
	return ret
}

func getStagedFiles() ([]string, error) {
	cmd := []string{"git", "diff", "--staged", "--name-only", "--diff-filter=A"}
	command := exec.Command(cmd[0], cmd[1:]...)
	output, err := command.Output()

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			retCode := exitErr.ExitCode()
			return nil, fmt.Errorf("command '%v' returned unexpected exit code %d (expected %d): %s",
				cmd, retCode, 0, string(output))
		}
		return nil, fmt.Errorf("error running get staged files: %w", err)
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}
