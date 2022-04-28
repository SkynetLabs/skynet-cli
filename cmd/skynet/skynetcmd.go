package main

import (
	"fmt"
	"os"

	"github.com/SkynetLabs/go-skynet/v2"
	"github.com/spf13/cobra"
	"gitlab.com/NebulousLabs/errors"
)

var (
	skynetDownloadCmd = &cobra.Command{
		Use:   "download [skylink] [destination]",
		Short: "Download a skylink from skynet.",
		Long: `Download a file from skynet using a skylink. Use the --portal flag to
fetch a skylink file from a chosen skynet portal.`,
		Run: skynetdownloadcmd,
	}

	skynetSkykeyCmd = &cobra.Command{
		Use:   "skykey",
		Short: "Perform skykey operations.",
		// A subcommand must be provided.
	}

	skynetSkykeyAddCmd = &cobra.Command{
		Use:   "add [skykey]",
		Short: "Add a skykey to Skynet.",
		Long:  "Store the given base-64 encoded skykey with the skykey manager.",
		Run:   wrap(skynetaddskykeycmd),
	}

	skynetSkykeyCreateCmd = &cobra.Command{
		Use:   "create [skykey name]",
		Short: "Create a skykey on Skynet.",
		Long:  "Returns a new skykey created and stored under that name.",
		Run:   wrap(skynetcreateskykeycmd),
	}

	skynetSkykeyGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Perform skykey get operations.",
		// A subcommand must be provided.
	}

	skynetSkykeyGetIDCmd = &cobra.Command{
		Use:   "id [skykey id]",
		Short: "Get a skykey given its ID.",
		Long:  "Get a skykey given its ID.",
		Run:   wrap(skynetgetskykeyidcmd),
	}

	skynetSkykeyGetNameCmd = &cobra.Command{
		Use:   "name [skykey name]",
		Short: "Get a skykey given its name.",
		Long:  "Get a skykey given its name.",
		Run:   wrap(skynetgetskykeynamecmd),
	}

	skynetSkykeyGetSkykeysCmd = &cobra.Command{
		Use:   "list",
		Short: "Get a list of all skykeys on Skynet.",
		Long:  "Get a list of all skykeys on Skynet.",
		Run:   wrap(skynetgetskykeyscmd),
	}

	skynetUploadCmd = &cobra.Command{
		Use:   "upload [source path]",
		Short: "Upload a file or a directory to Skynet.",
		Long: `Upload a file or a directory to Skynet. A skylink will be produced
which can be shared and used to retrieve the file. If the given path is
a directory all files under that directory will be uploaded individually and
an individual skylink will be produced for each. All files that get uploaded
will be pinned to Skynet Portal used for the upload, meaning that the portal
will pay for storage and repairs until the files are manually deleted.`,
		Run: wrap(skynetuploadcmd),
	}
)

// skynetcmd displays general info about the skynet cli.
func skynetcmd() {
	// Check if the user wants to generate the documentation
	if generateDocs {
		generateSkyetDocs()
	}

	// Print General Info
	fmt.Printf("%s\n\n", binDescription)

	// Get Default Portal
	fmt.Printf("Default Skynet Portal: %v\n", skynet.DefaultPortalURL())
}

// skynetaddskykey stores the given base-64 encoded skykey with the skykey
// manager.
func skynetaddskykeycmd(skykey string) {
	// Get the addskykey options.
	opts := skynet.DefaultAddSkykeyOptions
	client := initClientAndOptions(&opts.Options)

	err := client.AddSkykey(skykey, opts)
	if err != nil {
		err = errors.AddContext(err, fmt.Sprintf("AddSkykey Options: %+v\n", opts))
		die("Unable to add skykey:", err)
	}
	fmt.Println("Successfully added skykey!")
}

// skynetcreateskykey returns a new skykey created and stored under that name.
func skynetcreateskykeycmd(name, skykeyType string) {
	// Get the createskykey options.
	opts := skynet.DefaultCreateSkykeyOptions
	client := initClientAndOptions(&opts.Options)

	skykey, err := client.CreateSkykey(name, skykeyType, opts)
	if err != nil {
		err = errors.AddContext(err, fmt.Sprintf("CreateSkykey Options: %+v\n", opts))
		die("Unable to create skykey:", err)
	}
	fmt.Println("Successfully created skykey! Skykey:", skykey)
}

func skynetgetskykeyidcmd(id string) {
	// Get the getskykeyid options.
	opts := skynet.DefaultGetSkykeyOptions
	client := initClientAndOptions(&opts.Options)

	skykey, err := client.GetSkykeyByID(id, opts)
	if err != nil {
		err = errors.AddContext(err, fmt.Sprintf("GetSkykey Options: %+v\n", opts))
		die("Unable to get skykey by id:", err)
	}
	fmt.Println("Successfully got skykey! Skykey:", skykey)
}

func skynetgetskykeynamecmd(name string) {
	// Get the getskykeyname options.
	opts := skynet.DefaultGetSkykeyOptions
	client := initClientAndOptions(&opts.Options)

	skykey, err := client.GetSkykeyByName(name, opts)
	if err != nil {
		err = errors.AddContext(err, fmt.Sprintf("GetSkykey Options: %+v\n", opts))
		die("Unable to get skykey by name:", err)
	}
	fmt.Println("Successfully got skykey! Skykey:", skykey)
}

// skynetgetskykeyscmd gets a list of all skykeys.
func skynetgetskykeyscmd() {
	// Get the getskykeys options.
	opts := skynet.DefaultGetSkykeysOptions
	client := initClientAndOptions(&opts.Options)

	skykeys, err := client.GetSkykeys(opts)
	if err != nil {
		err = errors.AddContext(err, fmt.Sprintf("GetSkykeys Options: %+v\n", opts))
		die("Unable to get skykeys:", err)
	}
	fmt.Println("Successfully got skykeys! Skykeys:", skykeys)
}

// skynetdownloadcmd will perform the download of a skylink.
func skynetdownloadcmd(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		_ = cmd.UsageFunc()(cmd)
		os.Exit(exitCodeUsage)
	}

	// Get inputs
	skylink := args[0]
	filename := args[1]

	// Get the download options.
	opts := skynet.DefaultDownloadOptions
	client := initClientAndOptions(&opts.Options)
	if downloadSkykeyName != "" {
		opts.SkykeyName = downloadSkykeyName
	}
	if downloadSkykeyID != "" {
		opts.SkykeyID = downloadSkykeyID
	}

	// Download Skylink
	err := client.DownloadFile(filename, skylink, opts)
	if err != nil {
		err = errors.AddContext(err, fmt.Sprintf("Download Options: %+v\n", opts))
		die("Unable to download skylink:", err)
	}
	fmt.Println("Successfully downloaded skylink!")
}

// skynetuploadcmd will upload a file or directory to Skynet. If --dry-run is
// passed, it will fetch the skylinks without uploading.
func skynetuploadcmd(sourcePath string) {
	// Get the upload options.
	opts := skynet.DefaultUploadOptions
	client := initClientAndOptions(&opts.Options)
	if portalFileFieldName != "" {
		opts.PortalFileFieldName = portalFileFieldName
	}
	if portalDirectoryFileFieldName != "" {
		opts.PortalDirectoryFileFieldName = portalDirectoryFileFieldName
	}
	if customFilename != "" {
		opts.CustomFilename = customFilename
	}
	if customDirname != "" {
		opts.CustomDirname = customDirname
	}
	if uploadSkykeyName != "" {
		opts.SkykeyName = uploadSkykeyName
	}
	if uploadSkykeyID != "" {
		opts.SkykeyID = uploadSkykeyID
	}

	skylink, uploadType, err := upload(sourcePath, client, opts)
	if err != nil {
		err = errors.AddContext(err, fmt.Sprintf("Upload Options: %+v\n", opts))
		die(fmt.Sprintf("Unable to upload %v: %v\n", uploadType, err))
	}

	fmt.Printf("Successfully uploaded %v! Skylink: %v\n", uploadType, skylink)
}

// upload uploads the given path.
func upload(sourcePath string, client skynet.SkynetClient, opts skynet.UploadOptions) (skylink string, uploadType string, err error) {
	// Open the source file.
	file, err := os.Open(sourcePath)
	if err != nil {
		return "", "path", errors.AddContext(err, "Unable to open source path")
	}
	fi, err := file.Stat()
	if err != nil {
		err = errors.Extend(err, file.Close())
		return "", "path", errors.AddContext(err, "Unable to fetch source fileinfo")
	}
	err = file.Close()
	if err != nil {
		return "", "path", errors.AddContext(err, "Unable to close file")
	}

	// Upload File
	if !fi.IsDir() {
		skylink, err = client.UploadFile(sourcePath, opts)
		if err != nil {
			return "", "file", errors.AddContext(err, "Unable to upload file")
		}
		return skylink, "file", nil
	}

	// Upload directory
	skylink, err = client.UploadDirectory(sourcePath, opts)
	if err != nil {
		return "", "directory", errors.AddContext(err, "Unable to upload directory")
	}
	return skylink, "directory", nil
}

// initClientAndOptions initializes a client and common options from the
// persistent root flags that are common to all commands. Any available options
// in `opts` will be used if the option is not overridden with a root flag.
func initClientAndOptions(opts *skynet.Options) skynet.SkynetClient {
	if endpointPath != "" {
		opts.EndpointPath = endpointPath
	}
	if apiKey != "" {
		opts.APIKey = apiKey
	}
	if customUserAgent != "" {
		opts.CustomUserAgent = customUserAgent
	}
	// Create a client with specified portal (or "" if not specified) default
	// options. Custom options will be passed into the API call itself.
	client := skynet.NewCustom(skynetPortal, skynet.Options{})
	return client
}
