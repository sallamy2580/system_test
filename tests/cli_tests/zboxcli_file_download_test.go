package cli_tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	climodel "github.com/0chain/system_test/internal/cli/model"
	cliutils "github.com/0chain/system_test/internal/cli/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func TestDownload(t *testing.T) {
	t.Parallel()

	// Create a folder to keep all the generated files to be uploaded
	err := os.MkdirAll("tmp", os.ModePerm)
	require.Nil(t, err)

	// Success Scenarios
	t.Run("Download File from Root Directory Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)
		originalFileChecksum := generateChecksum(t, filename)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download File Concurrently Should Work from two Different Directory", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(4096)
		filesize := int64(1024)
		remoteFilePaths := [2]string{"/dir1/", "/dir2/"}

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		fileNameOfFirstDirectory := generateFileAndUpload(t, allocationID, remoteFilePaths[0], filesize)
		fileNameOfSecondDirectory := generateFileAndUpload(t, allocationID, remoteFilePaths[1], filesize)
		originalFirstFileChecksum := generateChecksum(t, fileNameOfFirstDirectory)
		originalSecondFileChecksum := generateChecksum(t, fileNameOfSecondDirectory)

		//deleting uploaded file from /dir1 since we will be downloading it now
		err := os.Remove(fileNameOfFirstDirectory)
		require.Nil(t, err)

		//deleting uploaded file from /dir2 since we will be downloading it now
		err = os.Remove(fileNameOfSecondDirectory)
		require.Nil(t, err)

		var outputList [2][]string
		var errorList [2]error
		var wg sync.WaitGroup

		fileNames := [2]string{fileNameOfFirstDirectory, fileNameOfSecondDirectory}
		for index, fileName := range fileNames {
			wg.Add(1)
			go func(currentFileName string, currentIndex int) {
				defer wg.Done()
				op, err := downloadFile(t, configPath, createParams(map[string]interface{}{
					"allocation": allocationID,
					"remotepath": remoteFilePaths[currentIndex] + filepath.Base(currentFileName),
					"localpath":  "tmp/",
				}), true)
				errorList[currentIndex] = err
				outputList[currentIndex] = op
			}(fileName, index)
		}

		wg.Wait()

		require.Nil(t, errorList[0], strings.Join(outputList[0], "\n"))
		require.Len(t, outputList[0], 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(fileNameOfFirstDirectory),
		)

		require.Equal(t, expected, outputList[0][1])
		downloadedFileFromFirstDirectoryChecksum := generateChecksum(t, "tmp/"+filepath.Base(fileNameOfFirstDirectory))

		require.Equal(t, originalFirstFileChecksum, downloadedFileFromFirstDirectoryChecksum)
		require.Nil(t, errorList[1], strings.Join(outputList[1], "\n"))
		require.Len(t, outputList[1], 2)

		expected = fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(fileNameOfSecondDirectory),
		)

		require.Equal(t, expected, outputList[1][1])
		downloadedFileFromSecondDirectoryChecksum := generateChecksum(t, "tmp/"+filepath.Base(fileNameOfSecondDirectory))
		require.Equal(t, originalSecondFileChecksum, downloadedFileFromSecondDirectoryChecksum)
	})

	t.Run("Download File from a Directory Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/dir/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)
		originalFileChecksum := generateChecksum(t, filename)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download File from Nested Directory Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/nested/dir/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)
		originalFileChecksum := generateChecksum(t, filename)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	//TODO: Directory download seems broken see https://github.com/0chain/blobber/issues/588
	t.Run("Download Entire Directory Should Work but does not see blobber/issues/588", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/nested/dir/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath,
			"localpath":  "tmp/dir",
		}), false)
		require.Error(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)
		require.Equal(t, "Error in file operation: No minimum consensus for file meta data of file", output[0])
	})

	//TODO: Directory share seems broken see https://github.com/0chain/blobber/issues/588
	t.Run("Download File From Shared Folder Should Work but does not see blobber/issues/588", func(t *testing.T) {
		t.Parallel()

		var authTicket, filename string

		filesize := int64(10)
		remotepath := "/"

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share Entire Folder from Another Wallet", func(t *testing.T) {
			allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUpload(t, allocationID, remotepath, filesize)

			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation": allocationID,
				"remotepath": remotepath,
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)

			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)
		})

		// Just register a wallet so that we can work further
		_, err := registerWallet(t, configPath)
		require.Nil(t, err)

		// Download file using auth-ticket: should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  "tmp/dir",
			"remotepath": "/" + filename,
		}), false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)
		require.Equal(t, "Error in file operation: No minimum consensus for file meta data of file", output[0])
	})

	t.Run("Download Entire Shared Folder Should Fail", func(t *testing.T) {
		t.Parallel()

		var authTicket, filename string

		filesize := int64(10)
		remotepath := "/"

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share Entire Folder from Another Wallet", func(t *testing.T) {
			allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUpload(t, allocationID, remotepath, filesize)

			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation": allocationID,
				"remotepath": remotepath,
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)

			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)
		})

		// Just register a wallet so that we can work further
		_, err := registerWallet(t, configPath)
		require.Nil(t, err)

		// Download file using auth-ticket: should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  "tmp/dir",
			"remotepath": "/",
		}), false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)
		require.Equal(t, "Error in file operation: please get files from folder, and download them one by one", output[0])
	})

	t.Run("Download Shared File Should Work", func(t *testing.T) {
		t.Parallel()

		var authTicket, filename, originalFileChecksum string

		filesize := int64(10)
		remotepath := "/"

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share File from Another Wallet", func(t *testing.T) {
			allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUpload(t, allocationID, remotepath, filesize)
			originalFileChecksum = generateChecksum(t, filename)

			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation": allocationID,
				"remotepath": remotepath + filepath.Base(filename),
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)

			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)
		})

		// Just register a wallet so that we can work further
		err := registerWalletAndLockReadTokens(t, configPath)
		require.Nil(t, err)

		// Download file using auth-ticket: should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  "tmp/",
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download Encrypted File Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(10 * MB)
		filesize := int64(10)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateRandomTestFileName(t)
		err := createFileWithSize(filename, filesize)
		require.Nil(t, err)
		originalFileChecksum := generateChecksum(t, filename)

		// Upload parameters
		uploadWithParam(t, configPath, map[string]interface{}{
			"allocation": allocationID,
			"localpath":  filename,
			"remotepath": remotepath + filepath.Base(filename),
			"encrypt":    "",
		})

		// Delete the uploaded file, since we will be downloading it now
		err = os.Remove(filename)
		require.Nil(t, err)

		// Downloading encrypted file should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  os.TempDir(),
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)
		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[len(output)-1])
		downloadedFileChecksum := generateChecksum(t, strings.TrimSuffix(os.TempDir(), "/")+"/"+filepath.Base(filename))
		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download Shared Encrypted File Should Work", func(t *testing.T) {
		t.Parallel()

		var authTicket, filename string

		filesize := int64(10)
		remotepath := "/"
		var allocationID string

		// register viewer wallet
		viewerWalletName := escapedTestName(t) + "_viewer"
		err = registerWalletForNameAndLockReadTokens(t, configPath, viewerWalletName)
		require.Nil(t, err)

		viewerWallet, err := getWalletForName(t, configPath, viewerWalletName)
		require.Nil(t, err)
		require.NotNil(t, viewerWallet)

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share File from Another Wallet", func(t *testing.T) {
			allocationID = setupAllocationAndReadLock(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUploadWithParam(t, allocationID, remotepath, filesize, map[string]interface{}{
				"encrypt": "",
			})
			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation":          allocationID,
				"remotepath":          remotepath + filepath.Base(filename),
				"encryptionpublickey": viewerWallet.EncryptionPublicKey,
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)

			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)
		})

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)

		file := "tmp/" + filepath.Base(filename)

		// Download file using auth-ticket: should work
		output, err := downloadFileForWallet(t, viewerWalletName, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  file,
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		require.Equal(t, expected, output[len(output)-1])

		os.Remove(file) //nolint

		// Download file using auth-ticket and lookuphash: should work
		output, err = downloadFileForWallet(t, viewerWalletName, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"lookuphash": GetReferenceLookup(allocationID, remotepath+filepath.Base(filename)),
			"localpath":  file,
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		require.Equal(t, expected, output[len(output)-1])
	})

	t.Run("Download From Shared Folder by Remotepath Should Work", func(t *testing.T) {
		t.Parallel()

		var authTicket, filename, originalFileChecksum string

		filesize := int64(10)
		remotepath := "/dir/"

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share File from Another Wallet", func(t *testing.T) {
			allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUpload(t, allocationID, remotepath, filesize)
			originalFileChecksum = generateChecksum(t, filename)
			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation": allocationID,
				"remotepath": remotepath,
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)

			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)
		})

		// Just register a wallet so that we can work further
		err := registerWalletAndLockReadTokens(t, configPath)
		require.Nil(t, err)

		// Download file using auth-ticket: should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  "tmp/",
			"remotepath": remotepath + filepath.Base(filename),
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download From Shared Folder by Lookup Hash Should Work", func(t *testing.T) {
		t.Parallel()

		var authTicket, lookuphash, filename, originalFileChecksum string

		filesize := int64(10)
		remotepath := "/dir/"

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share File from Another Wallet", func(t *testing.T) {
			allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUpload(t, allocationID, remotepath, filesize)
			originalFileChecksum = generateChecksum(t, filename)
			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation": allocationID,
				"remotepath": remotepath,
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)

			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)

			h := sha3.Sum256([]byte(fmt.Sprintf("%s:%s%s", allocationID, remotepath, filepath.Base(filename))))
			lookuphash = fmt.Sprintf("%x", h)
			require.NotEqual(t, "", lookuphash, "Lookup Hash: ", lookuphash)
		})

		// Just register a wallet so that we can work further
		err := registerWalletAndLockReadTokens(t, configPath)
		require.Nil(t, err)

		// Download file using auth-ticket: should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  "tmp/",
			"lookuphash": lookuphash,
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download Shared File without Paying Should Not Work", func(t *testing.T) {
		t.Parallel()

		var authTicket, filename string

		filesize := int64(10)
		remotepath := "/"

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share File from Another Wallet", func(t *testing.T) {
			allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUpload(t, allocationID, remotepath, filesize)
			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation": allocationID,
				"remotepath": remotepath + filepath.Base(filename),
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)

			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)
		})

		// Just register a wallet so that we can work further
		_, err := registerWallet(t, configPath)
		require.Nil(t, err)

		// Download file using auth-ticket: should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  "tmp/",
		}), true)
		require.NotNil(t, err)
		require.Len(t, output, 3)
	})

	t.Run("Download Shared File by Paying Should Work", func(t *testing.T) {
		t.Parallel()

		var allocationID, authTicket, filename string

		filesize := int64(10)
		remotepath := "/"

		// This test creates a separate wallet and allocates there, test nesting is required to create another wallet json file
		t.Run("Share File from Another Wallet", func(t *testing.T) {
			allocationID = setupAllocation(t, configPath, map[string]interface{}{
				"size":   10 * 1024,
				"tokens": 1,
			})
			filename = generateFileAndUpload(t, allocationID, remotepath, filesize)
			require.NotEqual(t, "", filename)

			// Delete the uploaded file from tmp folder if it exist,
			// since we will be downloading it now
			err := os.RemoveAll("tmp/" + filepath.Base(filename))
			require.Nil(t, err)

			shareParam := createParams(map[string]interface{}{
				"allocation": allocationID,
				"remotepath": remotepath + filepath.Base(filename),
			})

			output, err := shareFolderInAllocation(t, configPath, shareParam)
			require.Nil(t, err, strings.Join(output, "\n"))
			require.Len(t, output, 1)
			authTicket, err = extractAuthToken(output[0])
			require.Nil(t, err, "extract auth token failed")
			require.NotEqual(t, "", authTicket, "Ticket: ", authTicket)
		})

		err = registerWalletAndLockReadTokens(t, configPath)
		require.Nil(t, err)
		// Download file using auth-ticket: should work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"authticket": authTicket,
			"localpath":  "tmp/",
		}), false)

		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, filepath.Base(filename))
	})

	t.Run("Download File Thumbnail Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		thumbnail := escapedTestName(t) + "thumbnail.png"
		//nolint
		thumbnailSize := generateThumbnail(t, thumbnail)

		defer func() {
			// Delete the downloaded thumbnail file
			err = os.Remove(thumbnail)
			require.Nil(t, err)
		}()

		filename := generateFileAndUploadWithParam(t, allocationID, remotepath, filesize, map[string]interface{}{
			"thumbnailpath": thumbnail,
		})

		// Delete the uploaded file, since we will be downloading it now
		err = os.Remove(filename)
		require.Nil(t, err)

		localPath := filepath.Join(os.TempDir(), filepath.Base(filename))

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  localPath,
			"thumbnail":  nil,
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		stat, err := os.Stat(localPath)
		require.Nil(t, err)
		require.Equal(t, thumbnailSize, int(stat.Size()))
	})

	t.Run("Download to Non-Existent Path Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)
		originalFileChecksum := generateChecksum(t, filename)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/tmp2/" + filepath.Base(filename),
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/tmp2/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download File With Only startblock Should Work", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := getFileStats(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": "/" + filepath.Base(filename),
			"json":       "",
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)

		var stats map[string]climodel.FileStats

		err = json.Unmarshal([]byte(output[0]), &stats)
		require.Nil(t, err)
		var data climodel.FileStats
		for _, data = range stats {
			break
		}

		startBlock := int64(5) // blocks 5 to 10 should be downloaded
		output, err = downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"startblock": startBlock,
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])

		info, err := os.Stat("tmp/" + filepath.Base(filename))
		require.Nil(t, err, "error getting file stats")
		// downloaded file size should equal to ratio of block downloaded by original file size
		require.Equal(t, float64(info.Size()), (float64(data.NumOfBlocks-(startBlock-1))/float64(data.NumOfBlocks))*float64(filesize))
	})

	t.Run("Download File With Only endblock Should Not Work", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		endBlock := int64(5)
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"endblock":   endBlock,
		}), false)

		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 3)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, "invalid parameter: X-Block-Num")
	})

	t.Run("Download File With startblock And endblock Should Work", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536, we upload 20 blocks and download 1 block
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := getFileStats(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": "/" + filepath.Base(filename),
			"json":       "",
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)

		var stats map[string]climodel.FileStats

		err = json.Unmarshal([]byte(output[0]), &stats)
		require.Nil(t, err)
		var data climodel.FileStats
		for _, data = range stats {
			break
		}

		startBlock := 1
		endBlock := 1
		// Minimum Startblock value should be 1 (since gosdk subtracts 1 from start block, so 0 would lead to startblock being -1).
		output, err = downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"startblock": startBlock,
			"endblock":   endBlock,
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])

		info, err := os.Stat("tmp/" + filepath.Base(filename))
		require.Nil(t, err, "error getting file stats")
		// downloaded file size should equal to ratio of block downloaded by original file size
		require.Equal(t, float64(info.Size()), (float64(endBlock-(startBlock-1))/float64(data.NumOfBlocks))*float64(filesize))
	})

	t.Run("Download File With startblock 0 and non-zero endblock should fail", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		startBlock := 0
		endBlock := 5
		// Minimum Startblock value should be 1 (since gosdk subtracts 1 from start block, so 0 would lead to startblock being -1).
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"startblock": startBlock,
			"endblock":   endBlock,
		}), true)
		require.NotNil(t, err)
		require.Len(t, output, 3)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, "invalid parameter: X-Block-Num")
	})

	t.Run("Download File With endblock greater than number of blocks should fail", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		startBlock := 1
		endBlock := 40
		// Minimum Startblock value should be 1 (since gosdk subtracts 1 from start block, so 0 would lead to startblock being -1).
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"startblock": startBlock,
			"endblock":   endBlock,
		}), true)

		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 3)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, "Invalid block number")
	})

	t.Run("Download with endblock less than startblock should fail", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		startBlock := 6
		endBlock := 4
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "/tmp",
			"startblock": startBlock,
			"endblock":   endBlock,
		}), false)

		require.NotNil(t, err)
		require.Len(t, output, 2)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, "start block should be less than end block")
	})

	t.Run("Download with negative startblock should fail", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		startBlock := -6
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"startblock": startBlock,
		}), true)

		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 3)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, "invalid parameter: X-Block-Num")
	})

	t.Run("Download with negative endblock should fail", func(t *testing.T) {
		t.Parallel()

		// 1 block is of size 65536
		allocSize := int64(655360 * 4)
		filesize := int64(655360 * 2)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		endBlock := -6
		startBlock := 1
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"endblock":   endBlock,
			"startblock": startBlock,
		}), false)

		require.NotNil(t, err)
		require.Len(t, output, 2)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, "start block should be less than end block")
	})

	t.Run("Download File With commit Flag Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)
		originalFileChecksum := generateChecksum(t, filename)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
			"commit":     "",
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 3)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])

		match := reCommitResponse.FindStringSubmatch(output[2])
		require.Len(t, match, 2)

		var commitResp climodel.CommitResponse
		err = json.Unmarshal([]byte(match[1]), &commitResp)
		require.Nil(t, err)
		require.NotEmpty(t, commitResp)

		require.Equal(t, "application/octet-stream", commitResp.MetaData.MimeType)
		require.Equal(t, filesize, commitResp.MetaData.Size)
		require.Equal(t, filepath.Base(filename), commitResp.MetaData.Name)
		require.Equal(t, remotepath+filepath.Base(filename), commitResp.MetaData.Path)
		require.Equal(t, "", commitResp.MetaData.EncryptedKey)
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	t.Run("Download File With blockspermarker Flag Should Work", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)
		originalFileChecksum := generateChecksum(t, filename)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation":      allocationID,
			"remotepath":      remotepath + filepath.Base(filename),
			"localpath":       "tmp/",
			"blockspermarker": 1,
		}), true)
		require.Nil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 2)

		expected := fmt.Sprintf(
			"Status completed callback. Type = application/octet-stream. Name = %s",
			filepath.Base(filename),
		)
		require.Equal(t, expected, output[1])
		downloadedFileChecksum := generateChecksum(t, "tmp/"+filepath.Base(filename))

		require.Equal(t, originalFileChecksum, downloadedFileChecksum)
	})

	// Failure Scenarios

	t.Run("Download File from Non-Existent Allocation Should Fail", func(t *testing.T) {
		t.Parallel()

		output, err := registerWallet(t, configPath)
		require.Nil(t, err, strings.Join(output, "\n"))

		output, err = downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": "12334qe",
			"remotepath": "/",
			"localpath":  "tmp/",
		}), false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)

		require.Equal(t, "Error fetching the allocation allocation_fetch_error: "+
			"Error fetching the allocation.internal_error: can't get allocation: error retrieving allocation: 12334qe, error: record not found", output[0])
	})

	t.Run("Download File from Other's Allocation Should Fail", func(t *testing.T) {
		t.Parallel()

		var otherAllocationID, otherFilename string

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		t.Run("Get Other Allocation ID", func(t *testing.T) {
			otherAllocationID = setupAllocation(t, configPath, map[string]interface{}{
				"size": allocSize,
			})
			otherFilename = generateFileAndUpload(t, otherAllocationID, remotepath, filesize)
		})

		// Delete the uploaded file, since we will be downloading it now
		err = os.Remove(otherFilename)
		require.Nil(t, err)

		// Download using otherAllocationID: should not work
		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": otherAllocationID,
			"remotepath": remotepath + filepath.Base(otherFilename),
			"localpath":  "tmp/",
		}), false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.True(t, len(output) > 0)

		require.Equal(t, "Error in file operation: No minimum consensus for file meta data of file", output[len(output)-1])
	})

	t.Run("Download Non-Existent File Should Fail", func(t *testing.T) {
		t.Parallel()

		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   10000,
			"tokens": 1,
		})

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + "hello.txt",
			"localpath":  "tmp/",
		}), false)

		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)
		require.Equal(t, "Error in file operation: No minimum consensus for file meta data of file", output[0])
	})

	t.Run("Download without any Parameter Should Fail", func(t *testing.T) {
		t.Parallel()

		output, err := registerWallet(t, configPath)
		require.Nil(t, err, strings.Join(output, "\n"))

		output, err = downloadFile(t, configPath, "", false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)

		require.Equal(t, "Error: remotepath / authticket flag is missing", output[0])
	})

	t.Run("Download from Allocation without other Parameter Should Fail", func(t *testing.T) {
		t.Parallel()

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   10000,
			"tokens": 1,
		})

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
		}), false)

		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)
		require.Equal(t, "Error: remotepath / authticket flag is missing", output[0])
	})

	t.Run("Download File Without read-lock Should Fail", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocation(t, configPath, map[string]interface{}{
			"size": allocSize,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
		}), false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 3)
		aggregatedOutput := strings.Join(output, " ")
		require.Contains(t, aggregatedOutput, "not enough tokens")
	})

	t.Run("Download File using Expired Allocation Should Fail", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
			"expire": "1h",
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		// Delete the uploaded file, since we will be downloading it now
		err := os.Remove(filename)
		require.Nil(t, err)

		params := createParams(map[string]interface{}{
			"allocation": allocationID,
			"expiry":     "-1h",
		})
		output, err := updateAllocation(t, configPath, params, true)
		require.Nil(t, err, strings.Join(output, "\n"))

		output, err = downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  "tmp/",
		}), false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)
		require.Equal(t, "Error in file operation: No minimum consensus for file meta data of file", output[0])
	})

	t.Run("Download File to Existing File Should Fail", func(t *testing.T) {
		t.Parallel()

		allocSize := int64(2048)
		filesize := int64(256)
		remotepath := "/"

		allocationID := setupAllocationAndReadLock(t, configPath, map[string]interface{}{
			"size":   allocSize,
			"tokens": 1,
		})

		filename := generateFileAndUpload(t, allocationID, remotepath, filesize)

		output, err := downloadFile(t, configPath, createParams(map[string]interface{}{
			"allocation": allocationID,
			"remotepath": remotepath + filepath.Base(filename),
			"localpath":  os.TempDir(),
		}), false)
		require.NotNil(t, err, strings.Join(output, "\n"))
		require.Len(t, output, 1)

		expected := fmt.Sprintf(
			"Download failed. Local file already exists '%s'",
			strings.TrimSuffix(os.TempDir(), "/")+"/"+filepath.Base(filename),
		)
		require.Equal(t, expected, output[0])
	})
}

func setupAllocationAndReadLock(t *testing.T, cliConfigFilename string, extraParam map[string]interface{}) string {
	tokens := float64(1)
	if tok, ok := extraParam["tokens"]; ok {
		token, err := strconv.ParseFloat(fmt.Sprintf("%v", tok), 64)
		require.Nil(t, err)
		tokens = token
	}

	allocationID := setupAllocation(t, cliConfigFilename, extraParam)

	// Lock half the tokens for read pool
	readPoolParams := createParams(map[string]interface{}{
		"tokens": tokens / 2,
	})
	output, err := readPoolLock(t, cliConfigFilename, readPoolParams, true)
	require.Nil(t, err, strings.Join(output, "\n"))
	require.Len(t, output, 1)
	require.Equal(t, "locked", output[0])

	return allocationID
}

func downloadFile(t *testing.T, cliConfigFilename, param string, retry bool) ([]string, error) {
	return downloadFileForWallet(t, escapedTestName(t), cliConfigFilename, param, retry)
}

func downloadFileForWallet(t *testing.T, wallet, cliConfigFilename, param string, retry bool) ([]string, error) {
	cliutils.Wait(t, 15*time.Second) // TODO replace with pollers
	t.Logf("Downloading file...")
	cmd := fmt.Sprintf(
		"./zbox download %s --silent --wallet %s --configDir ./config --config %s",
		param,
		wallet+"_wallet.json",
		cliConfigFilename,
	)

	if retry {
		return cliutils.RunCommand(t, cmd, 3, time.Second*2)
	} else {
		return cliutils.RunCommandWithoutRetry(cmd)
	}
}

func generateChecksum(t *testing.T, filePath string) string {
	t.Logf("Generating checksum for file [%v]...", filePath)

	output, err := cliutils.RunCommandWithoutRetry("shasum -a 256 " + filePath)
	require.Nil(t, err, "Checksum generation for file %v failed", filePath, strings.Join(output, "\n"))
	require.Greater(t, len(output), 0)

	matcher := regexp.MustCompile("(.*) " + filePath + "$")
	require.Regexp(t, matcher, output[0], "Checksum execution output did not match expected", strings.Join(output, "\n"))

	return matcher.FindAllStringSubmatch(output[0], 1)[0][1]
}
