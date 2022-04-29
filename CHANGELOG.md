Version Scheme
--------------
Skynet CLI uses the following versioning scheme, vX.X.X
 - First Digit signifies a major (compatibility breaking) release
 - Second Digit signifies a major (non compatibility breaking) release
 - Third Digit signifies a minor or patch release

Version History
---------------

Latest:

## Apr 29, 2022:
### v2.1.0
**Key Updates**
- Add --skynet-api-key option.

## Oct 22, 2020:
### v2.0.1
**Key Updates**
- Fixed single-file directory uploads being uploaded as files instead of
  directories.
- Fixed commands not working due to default options being ignored
- Added Homebrew installation method for OSX

## Sep 4, 2020:
### v2.0.0
**Key Updates**
- Updated to use v2 of the API.
- Added `version` command.

## Aug 13, 2020:
### v1.1.0
**Key Updates**
- The encryption API and support for encryption in uploads and decryption in downloads were added.
- Connection options such as --api-key were added.
- Some bugs were fixed, including a bug that prevented `skynet upload` from working.

## Jun 16, 2020:
### v1.0.0
**Key Updates**
- Created skynet cli with upload and download.
