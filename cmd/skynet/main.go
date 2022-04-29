package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"gitlab.com/NebulousLabs/errors"
)

const (
	// binDescription is the binary description.
	binDescription = `Perform actions related to Skynet, a file sharing and data publication platform
built on top of Sia.`

	// binName is the binary name of skynet-cli.
	binName = "skynet"

	// docREADMEFile is the filepath to the README file in /doc
	docREADMEFile = "./doc/README.md"

	// mainDocFile is the filepath to the top level doc file created by cobra
	mainDocFile = "./doc/skynet.md"

	// version is the current version of skynet-cli.
	version = "2.0.1"
)

// Exit codes.
// inspired by sysexits.h
const (
	exitCodeGeneral = 1  // Not in sysexits.h, but is standard practice.
	exitCodeUsage   = 64 // EX_USAGE in sysexits.h
)

var (
	fullDescription = fmt.Sprintf("%s\n\n%s", versionString, binDescription)

	rootCmd *cobra.Command

	// versionString is the version string of the executable.
	versionString = fmt.Sprintf("%s %s", binName, version)
)

var (
	// generateDocs will trigger cobra to auto generate the documentation for
	// skynet commands
	generateDocs bool

	// Common Skynet Flags
	//
	// skynetPortal will define a Skynet Portal to use instead of the default.
	skynetPortal string

	// endpointPath is the relative URL path of the endpoint.
	endpointPath string

	// apiKey is the skyd authentication password to use for a single Skynet
	// node.
	apiKey string

	// skynetAPIKey is the authentication API key to use for a Skynet portal
	// (sets the "Skynet-Api-Key" header).
	skynetAPIKey string

	// customUserAgent is the custom user agent to use.
	customUserAgent string

	// Upload Flags
	//
	// portalFileFieldName is the fieldName for files on the portal.
	portalFileFieldName string

	// portalDirectoryFileFieldName is the fieldName for directory files on the
	// portal.
	portalDirectoryFileFieldName string

	// customFilename is the custom filename to use for the upload. If this is
	// empty, the filename of the file being uploaded will be used by default.
	customFilename string

	// customDirname is the custom directory filename to use for the upload. If
	// this is empty, the name of the directory being uploaded will be used by
	// default.
	customDirname string

	// uploadSkykeyName is the name of the skykey used to encrypt the upload.
	uploadSkykeyName string

	// uploadSkykeyID is the ID of the skykey used to encrypt the upload.
	uploadSkykeyID string

	// Download Flags
	//
	// downloadSkykeyName is the name of the skykey used to decrypt the upload.
	downloadSkykeyName string

	// downloadSkykeyID is the ID of the skykey used to decrypt the upload.
	downloadSkykeyID string
)

// copyDocFile will copy the top level auto generated doc file into a README for
// /doc
func copyDocFile() (err error) {
	// Open the main doc file
	source, err := os.Open(mainDocFile)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Compose(err, source.Close())
	}()

	// Open the README file
	destination, err := os.Create(docREADMEFile)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Compose(err, destination.Close())
	}()

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
		Long:  fullDescription,
		Run:   wrap(skynetcmd),
	}

	rootCmd.AddCommand(versionCmd)

	// Add Skynet Commands
	rootCmd.AddCommand(skynetDownloadCmd, skynetSkykeyCmd, skynetUploadCmd)
	skynetSkykeyCmd.AddCommand(skynetSkykeyAddCmd, skynetSkykeyCreateCmd, skynetSkykeyGetCmd, skynetSkykeyGetSkykeysCmd)
	skynetSkykeyGetCmd.AddCommand(skynetSkykeyGetIDCmd, skynetSkykeyGetNameCmd)

	// Add flags.
	rootCmd.Flags().BoolVarP(&generateDocs, "generate-docs", "d", false, "Generate the docs for Skynet")
	rootCmd.PersistentFlags().StringVar(&skynetPortal, "portal", "", "Use a Skynet portal other than the default")
	rootCmd.PersistentFlags().StringVar(&endpointPath, "endpoint-path", "", "Relative URL path of the endpoint to use")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Authentication password to use for a single Skynet node.")
	rootCmd.PersistentFlags().StringVar(&skynetAPIKey, "skynet-api-key", "", "Authentication API key to use for a Skynet portal")
	rootCmd.PersistentFlags().StringVar(&customUserAgent, "custom-user-agent", "", "Custom user agent to use")
	// Upload flags.
	skynetUploadCmd.Flags().StringVar(&portalFileFieldName, "file-field-name", "", "Defines the fieldName for the files on the portal")
	skynetUploadCmd.Flags().StringVar(&portalDirectoryFileFieldName, "directory-field-name", "", "Defines the fieldName for the directory files on the portal")
	skynetUploadCmd.Flags().StringVar(&customFilename, "filename", "", "Custom filename for the uploaded file")
	skynetUploadCmd.Flags().StringVar(&customDirname, "dirname", "", "Custom dirname for the uploaded directory")
	skynetUploadCmd.Flags().StringVar(&uploadSkykeyName, "skykey-name", "", "Name of the skykey on the portal used to encrypt the upload")
	skynetUploadCmd.Flags().StringVar(&uploadSkykeyID, "skykey-id", "", "ID of the skykey on the portal used to encrypt the upload")
	// Download flags.
	skynetDownloadCmd.Flags().StringVar(&downloadSkykeyName, "skykey-name", "", "Name of the skykey on the portal used to decrypt the download")
	skynetDownloadCmd.Flags().StringVar(&downloadSkykeyID, "skykey-id", "", "ID of the skykey on the portal used to decrypt the download")

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
