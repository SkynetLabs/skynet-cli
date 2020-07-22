package main

import (
	"fmt"
	"os"

	"github.com/NebulousLabs/go-skynet"
	"github.com/spf13/cobra"
)

var (
	skynetBlacklistCmd = &cobra.Command{
		Use:   "blacklist [skylink]",
		Short: "Needs SDK Implementation",
		Long:  "Needs SDK Implementation",

		/*
		   		Short: "Blacklist a skylink from skynet.",
		   		Long: `Blacklist a skylink from skynet. Use the --remove flag to
		   remove a skylink from the blacklist.`,
		*/
		Run: skynetblacklistcmd,
	}

	skynetConvertCmd = &cobra.Command{
		Use:   "convert [source siaPath] [destination siaPath]",
		Short: "Needs SDK Implementation",
		Long:  "Needs SDK Implementation",
		/*
				Short: "Convert a siafile to a skyfile with a skylink.",
				Long: `Convert a siafile to a skyfile and then generate its skylink. A new skylink
			will be created in the user's skyfile directory. The skyfile and the original
			siafile are both necessary to pin the file and keep the skylink active. The
			skyfile will consume an additional 40 MiB of storage.`,
		*/
		Run: wrap(skynetconvertcmd),
	}

	skynetDownloadCmd = &cobra.Command{
		Use:   "download [skylink] [destination]",
		Short: "Download a skylink from skynet.",
		Long: `Download a file from skynet using a skylink. Use the --portal flag to
fetch a skylink file from a chosen skynet portal.`,
		Run: skynetdownloadcmd,
	}

	skynetLsCmd = &cobra.Command{
		Use:   "ls",
		Short: "Needs SDK Implementation",
		Long:  "Needs SDK Implementation",
		/*
		   Short: "List all skyfiles that the user has pinned.",
		   		Long: `List all skyfiles that the user has pinned along with the corresponding
		   skylinks. By default, only files in var/skynet/ will be displayed. The --root
		   flag can be used to view skyfiles pinned in other folders.`,
		*/
		Run: skynetlscmd,
	}

	skynetPinCmd = &cobra.Command{
		Use:   "pin [skylink] [destination siapath]",
		Short: "Needs SDK Implementation",
		Long:  "Needs SDK Implementation",
		/*
					Short: "Pin a skylink from skynet by re-uploading it yourself.",
					Long: `Pin the file associated with this skylink by re-uploading an exact copy. This
			ensures that the file will still be available on skynet as long as you continue
			maintaining the file in your renter.`,
		*/
		Run: wrap(skynetpincmd),
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
		Use:   "skykeys",
		Short: "Get a list of all skykeys on Skynet.",
		Long:  "Get a list of all skykeys on Skynet.",
		Run:   wrap(skynetgetskykeyscmd),
	}

	skynetUnpinCmd = &cobra.Command{
		Use:   "unpin [siapath]",
		Short: "Needs SDK Implementation",
		Long:  "Needs SDK Implementation",
		/*
					Short: "Unpin pinned skyfiles or directories.",
					Long: `Unpin one or more pinned skyfiles or directories at the given siapaths. The
			files and directories will continue to be available on Skynet if other nodes have pinned them.`,
		*/
		Run: skynetunpincmd,
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
	fmt.Printf("Use Skynet to upload and share content\n\n")

	// Get Default Portal
	downloadOptions := skynet.DefaultDownloadOptions
	fmt.Printf("Default Skynet Portal: %v \n\n", downloadOptions.PortalURL)
}

// skynetaddskykey stores the given base-64 encoded skykey with the skykey
// manager.
func skynetaddskykeycmd(skykey string) {
	// Get the addskykey options.
	opts := skynet.DefaultAddSkykeyOptions
	opts.Options = getCommonOptions(opts.Options)
	fmt.Printf("AddSkykey Options: %+v\n", opts)

	err := skynet.AddSkykey(skykey, opts)
	if err != nil {
		die("Unable to add skykey:", err)
	}
	fmt.Println("Successfully added skykey!")
}

// skynetcreateskykey returns a new skykey created and stored under that name.
func skynetcreateskykeycmd(name, skykeyType string) {
	// Get the createskykey options.
	opts := skynet.DefaultCreateSkykeyOptions
	opts.Options = getCommonOptions(opts.Options)
	fmt.Printf("CreateSkykey Options: %+v\n", opts)

	skykey, err := skynet.CreateSkykey(name, skykeyType, opts)
	if err != nil {
		die("Unable to create skykey:", err)
	}
	fmt.Println("Successfully created skykey! Skykey:", skykey)
}

func skynetgetskykeyidcmd(id string) {
	// Get the getskykeyid options.
	opts := skynet.DefaultGetSkykeyOptions
	opts.Options = getCommonOptions(opts.Options)
	fmt.Printf("GetSkykey Options: %+v\n", opts)

	skykey, err := skynet.GetSkykeyByID(id, opts)
	if err != nil {
		die("Unable to get skykey by id:", err)
	}
	fmt.Println("Successfully got skykey! Skykey:", skykey)
}

func skynetgetskykeynamecmd(name string) {
	// Get the getskykeyname options.
	opts := skynet.DefaultGetSkykeyOptions
	opts.Options = getCommonOptions(opts.Options)
	fmt.Printf("GetSkykey Options: %+v\n", opts)

	skykey, err := skynet.GetSkykeyByName(name, opts)
	if err != nil {
		die("Unable to get skykey by name:", err)
	}
	fmt.Println("Successfully got skykey! Skykey:", skykey)
}

// skynetgetskykeyscmd gets a list of all skykeys.
func skynetgetskykeyscmd() {
	// Get the getskykeys options.
	opts := skynet.DefaultGetSkykeysOptions
	opts.Options = getCommonOptions(opts.Options)
	fmt.Printf("GetSkykeys Options: %+v\n", opts)

	skykeys, err := skynet.GetSkykeys(opts)
	if err != nil {
		die("Unable to get skykeys:", err)
	}
	fmt.Println("Successfully got skykeys! Skykeys:", skykeys)
}

// skynetblacklistcmd handles adding and removing a skylink from the Skynet
// Blacklist
func skynetblacklistcmd(cmd *cobra.Command, args []string) {
	fmt.Println("Skynet Blacklist not support in SDK")
}

// skynetconvertcmd will convert an existing siafile to a skyfile and skylink on
// the Sia network.
func skynetconvertcmd(sourceSiaPathStr, destSiaPathStr string) {
	fmt.Println("Skynet convert not support in SDK")
}

// skynetdownloadcmd will perform the download of a skylink.
func skynetdownloadcmd(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		cmd.UsageFunc()(cmd)
		os.Exit(exitCodeUsage)
	}

	// Get inputs
	skylink := args[0]
	filename := args[1]

	// Get the download options.
	opts := skynet.DefaultDownloadOptions
	opts.Options = getCommonOptions(opts.Options)
	if downloadSkykeyName != "" {
		opts.SkykeyName = downloadSkykeyName
	}
	if downloadSkykeyID != "" {
		opts.SkykeyID = downloadSkykeyID
	}
	fmt.Printf("Download Options: %+v\n", opts)

	// Check whether the portal flag is set, if so update the portal to download
	// from.
	if skynetPortal != "" {
		opts.PortalURL = skynetPortal
	}

	// Download Skylink
	err := skynet.DownloadFile(filename, skylink, opts)
	if err != nil {
		die("Unable to download skylink:", err)
	}
	fmt.Println("Successfully downloaded skylink!")
}

// skynetlscmd is the handler for the command `siac skynet ls`. Works very
// similarly to 'siac renter ls' but defaults to the SkynetFolder and only
// displays files that are pinning skylinks.
func skynetlscmd(cmd *cobra.Command, args []string) {
	fmt.Println("Skynet ls not implemented in SDK")
}

// skynetpincmd will pin the file from this skylink.
func skynetpincmd(sourceSkylink, destSiaPath string) {
	fmt.Println("Skynet pin not implemented in SDK")
}

// skynetunpincmd will unpin and delete either a single or multiple files or
// directories from the Renter.
func skynetunpincmd(cmd *cobra.Command, skyPathStrs []string) {
	fmt.Println("Skynet pin not implemented in SDK")
}

// skynetuploadcmd will upload a file or directory to Skynet. If --dry-run is
// passed, it will fetch the skylinks without uploading.
func skynetuploadcmd(sourcePath string) {
	// Get the upload options.
	opts := skynet.DefaultUploadOptions
	opts.Options = getCommonOptions(opts.Options)
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
	fmt.Printf("Upload Options: %+v\n", opts)

	// Open the source file.
	file, err := os.Open(sourcePath)
	if err != nil {
		die("Unable to open source path:", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			die("Unable to close file:", err)
		}
	}()
	fi, err := file.Stat()
	if err != nil {
		die("Unable to fetch source fileinfo:", err)
	}

	// Upload File
	if !fi.IsDir() {
		skylink, err := skynet.UploadFile(sourcePath, opts)
		if err != nil {
			die("Unable to upload file:", err)
		}
		fmt.Println("Successfully uploaded skyfile! Skylink:", skylink)
		return
	}

	// Upload directory
	skylink, err := skynet.UploadDirectory(sourcePath, opts)
	if err != nil {
		die("Unable to upload directory:", err)
	}
	fmt.Println("Successfully uploaded directory! Skylink:", skylink)
}

func getCommonOptions(opts skynet.Options) skynet.Options {
	if skynetPortal != "" {
		opts.PortalURL = skynetPortal
	}
	if endpointPath != "" {
		opts.EndpointPath = endpointPath
	}
	if customUserAgent != "" {
		opts.CustomUserAgent = customUserAgent
	}
	return opts
}
