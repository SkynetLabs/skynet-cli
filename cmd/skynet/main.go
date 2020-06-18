package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
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

// Exit codes.
// inspired by sysexits.h
const (
	exitCodeGeneral = 1  // Not in sysexits.h, but is standard practice.
	exitCodeUsage   = 64 // EX_USAGE in sysexits.h
)

// die prints its arguments to stderr, then exits the program with the default
// error code.
func die(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(exitCodeGeneral)
}

func main() {
	root := &cobra.Command{
		Use:   os.Args[0],
		Short: "Perform actions related to Skynet",
		Long: `Perform actions related to Skynet, a file sharing and data publication platform
on top of Sia.`,
		Run: wrap(skynetcmd),
	}

	// Add Skynet Commands
	root.AddCommand(skynetBlacklistCmd, skynetConvertCmd, skynetDownloadCmd, skynetLsCmd, skynetPinCmd, skynetUnpinCmd, skynetUploadCmd)

	// Add Flags
	root.PersistentFlags().StringVar(&skynetPortal, "portal", "", "Use a Skynet portal other than the default")
	skynetUploadCmd.Flags().StringVar(&portalUploadPath, "upload-path", "", "Relative URL path of the upload endpoint")
	skynetUploadCmd.Flags().StringVar(&portalFileFieldName, "file-field-name", "", "Defines the fieldName for the files on the portal")
	skynetUploadCmd.Flags().StringVar(&portalDirectoryFileFieldName, "directory-field-name", "", "Defines the fieldName for the directory files on the portal")
	skynetUploadCmd.Flags().StringVar(&customFilename, "filename", "", "Custom filename for the uploaded file")

	// Build the docs
	err := doc.GenMarkdownTree(root, "./tmp")
	if err != nil {
		die(err)
	}

	// run
	if err := root.Execute(); err != nil {
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
