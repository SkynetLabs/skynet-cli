package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	// mainDocFile is the filepath to the top level doc file created by cobra
	mainDocFile = "./doc/skynet.md"

	// docREADMEFile is the filepath to the README file in /doc
	docREADMEFile = "./doc/README.md"
)

// Exit codes.
// inspired by sysexits.h
const (
	exitCodeGeneral = 1  // Not in sysexits.h, but is standard practice.
	exitCodeUsage   = 64 // EX_USAGE in sysexits.h
)

var rootCmd *cobra.Command

var (
	// generateDocs will trigger cobra to auto generate the documentation for
	// skynet commands
	generateDocs bool

	// Skynet Flags
	//
	// skynetPortal will define a Skynet Portal to use instead of the default.
	skynetPortal string

	// Upload Flags
	//
	// portalUploadPath is the relative URL path of the upload endpoint.
	portalUploadPath string

	// portalFileFieldName is the fieldName for files on the portal.
	portalFileFieldName string

	// portalDirectoryFileFieldName is the fieldName for directory files on the
	// portal.
	portalDirectoryFileFieldName string

	// customFilename is the custom filename to use for the upload. If this is
	// empty, the filename of the file being uploaded will be used by default.
	customFilename string
)

// copyDocFile will copy the top level auto generated doc file into a README for
// /doc
func copyDocFile() error {
	// Open the main doc file
	source, err := os.Open(mainDocFile)
	if err != nil {
		return err
	}
	defer source.Close()

	// Open the README file
	destination, err := os.Create(docREADMEFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	// Copy the main doc file into the README
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}
	return nil
}

// die prints its arguments to stderr, then exits the program with the default
// error code.
func die(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(exitCodeGeneral)
}

// generateSkyetDocs will have cobra auto generate the documentation for the
// skynet commands
func generateSkyetDocs() {
	// Build the docs
	err := doc.GenMarkdownTree(rootCmd, "./doc")
	if err != nil {
		die(err)
	}

	// Copy the main skynet.md file to the doc README file
	err = copyDocFile()
	if err != nil {
		die(err)
	}
}

func main() {
	rootCmd = &cobra.Command{
		Use:   "skynet",
		Short: "Perform actions related to Skynet",
		Long: `Perform actions related to Skynet, a file sharing and data publication platform
on top of Sia.`,
		Run: wrap(skynetcmd),
	}

	// Add Skynet Commands
	rootCmd.AddCommand(skynetBlacklistCmd, skynetConvertCmd, skynetDownloadCmd, skynetLsCmd, skynetPinCmd, skynetUnpinCmd, skynetUploadCmd)

	// Add Flags
	rootCmd.Flags().BoolVarP(&generateDocs, "", "d", false, "Generate the docs for skynet")
	rootCmd.PersistentFlags().StringVar(&skynetPortal, "portal", "", "Use a Skynet portal other than the default")
	skynetUploadCmd.Flags().StringVar(&portalUploadPath, "upload-path", "", "Relative URL path of the upload endpoint")
	skynetUploadCmd.Flags().StringVar(&portalFileFieldName, "file-field-name", "", "Defines the fieldName for the files on the portal")
	skynetUploadCmd.Flags().StringVar(&portalDirectoryFileFieldName, "directory-field-name", "", "Defines the fieldName for the directory files on the portal")
	skynetUploadCmd.Flags().StringVar(&customFilename, "filename", "", "Custom filename for the uploaded file")

	// run
	if err := rootCmd.Execute(); err != nil {
		// Since no commands return errors (all commands set Command.Run instead of
		// Command.RunE), Command.Execute() should only return an error on an
		// invalid command or flag. Therefore Command.Usage() was called (assuming
		// Command.SilenceUsage is false) and we should exit with exitCodeUsage.
		os.Exit(exitCodeUsage)
	}
}

// wrap wraps a generic command with a check that the command has been
// passed the correct number of arguments. The command must take only strings
// as arguments.
func wrap(fn interface{}) func(*cobra.Command, []string) {
	fnVal, fnType := reflect.ValueOf(fn), reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("wrapped function has wrong type signature")
	}
	for i := 0; i < fnType.NumIn(); i++ {
		if fnType.In(i).Kind() != reflect.String {
			panic("wrapped function has wrong type signature")
		}
	}

	return func(cmd *cobra.Command, args []string) {
		if len(args) != fnType.NumIn() {
			cmd.UsageFunc()(cmd)
			os.Exit(exitCodeUsage)
		}
		argVals := make([]reflect.Value, fnType.NumIn())
		for i := range args {
			argVals[i] = reflect.ValueOf(args[i])
		}
		fnVal.Call(argVals)
	}
}
