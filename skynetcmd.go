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
	// Print General Info
	fmt.Printf("Use Skynet to upload and share content\n\n")

	// Get Default Portal
	downloadOptions := skynet.DefaultDownloadOptions
	fmt.Printf("Default Skynet Portal: %v \n\n", downloadOptions.PortalURL)
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

	// Define Options
	opts := skynet.DefaultDownloadOptions

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
}

// skynetlscmd is the handler for the command `siac skynet ls`. Works very
// similar to 'siac renter ls' but defaults to the SkynetFolder and only
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
	// Get the upload options
	opts := skynet.DefaultUploadOptions
	if skynetPortal != "" {
		opts.PortalURL = skynetPortal
	}
	if portalUploadPath != "" {
		opts.PortalUploadPath = portalUploadPath
	}
	if portalFileFieldName != "" {
		opts.PortalFileFieldName = portalFileFieldName
	}
	if portalDirectoryFileFieldName != "" {
		opts.PortalDirectoryFileFieldName = portalDirectoryFileFieldName
	}
	if customFilename != "" {
		opts.CustomFilename = customFilename
	}

	// Open the source file.
	file, err := os.Open(sourcePath)
	if err != nil {
		die("Unable to open source path:", err)
	}
	defer file.Close()
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
	return
}
