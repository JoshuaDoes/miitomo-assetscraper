package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-unarr"
)

type Manifest struct {
	PackageURL        string           `json:"packageUrl"`
	RemoteVersionURL  string           `json:"remoteVersionUrl"`
	RemoteManifestURL string           `json:"remoteManifestUrl"`
	Version           string           `json:"version"`
	EngineVersion     string           `json:"engineVersion"`
	Assets            map[string]Asset `json:"assets"`
}
type Asset struct {
	MD5        string `json:"md5"`
	Path       string `json:"path"`
	Compressed bool   `json:"compressed"`
	Version    string `json:"version"`
	MiiStudio  bool   `json:"miistudio"`
	Later      bool   `json:"later"`
}
type CompletedDownloads struct {
	Assets []string `json:"assets"`
}

var (
	manifestURL string

	completedDownloads CompletedDownloads
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: You must specify the Miitomo Manifest URL.")
		return
	}

	manifestURL = os.Args[1]

	_, err := url.ParseRequestURI(manifestURL)
	if err != nil {
		fmt.Println("Error: The URL you specified is not a valid URL.")
		fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
		return
	}

	fmt.Println("> Fetching asset manifest")
	manifestGET, err := http.Get(manifestURL)
	if err != nil {
		fmt.Println("Error: Unable to fetch a response from the specified URL.")
		fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
		return
	}

	fmt.Println("> Parsing asset manifest")
	manifest := Manifest{}
	err = unmarshal(manifestGET, &manifest)
	if err != nil {
		fmt.Println("Error: Unable to parse the response from the specified URL.")
		fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
		return
	}

	fmt.Println("> Creating temporary directory")
	mkdir("temp")
	fmt.Println("> Deferring deletion of temporary directory to end of program")

	fmt.Println("> Getting list of assets to download")
	var assetList []string
	longestAssetName := 0
	for assetKey := range manifest.Assets {
		assetList = append(assetList, assetKey)
		if len(assetKey) > longestAssetName {
			longestAssetName = len(assetKey)
		}
	}

	fmt.Println("> Getting list of previously downloaded assets")
	loadCompletedDownloads()

	fmt.Println("> Beginning download of assets")
	longestPrintString := 0
	for index, assetKey := range assetList {
		percentageDone := percentOf((index + 1), len(assetList))

		asset := manifest.Assets[assetKey]
		assetFilePath := asset.Path
		assetDirectory := filepath.Dir(assetFilePath)

		if isCompletedDownload(assetKey) {
			printString := fmt.Sprintf("(%d/%d) Skipping [%s] as it already exists (%.2f%%)", (index + 1), len(assetList), assetFilePath, percentageDone)
			if len(printString) > longestPrintString {
				longestPrintString = len(printString)
			} else if len(printString) < longestPrintString {
				printString = printString + strings.Repeat(" ", longestPrintString-len(printString))
			}
			fmt.Print(printString + "\n")
		} else {
			printString := fmt.Sprintf("(%d/%d) Downloading [%s]...%s (%.2f%%)", (index + 1), len(assetList), assetFilePath, getDots(assetFilePath, longestAssetName), percentageDone)
			if len(printString) > longestPrintString {
				longestPrintString = len(printString)
			} else if len(printString) < longestPrintString {
				printString = printString + strings.Repeat(" ", longestPrintString-len(printString))
			}
			fmt.Print(printString + "\r")

			fileLocation, md5hash, err := fileDownload(manifest.PackageURL+"/"+assetFilePath, "temp/"+assetKey)
			if err != nil || fileLocation == "" {
				fmt.Println("Error: Unable to download the asset.")
				fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
				return
			}
			if md5hash != asset.MD5 {
				fmt.Println("Error: MD5 checksum does not match the downloaded asset.")
				fmt.Println("- manifest MD5 [" + asset.MD5 + "]")
				fmt.Println("- download MD5 [" + md5hash + "]")
				return
			}

			if asset.Compressed {
				zip, err := unarr.NewArchive(fileLocation)
				if err != nil {
					fmt.Println("Error: Unable to open the archive for extraction.")
					fmt.Println("- archive [" + fileLocation + "]")
					fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
					return
				}

				err = zip.Extract(assetDirectory)
				if err != nil {
					fmt.Println("Error: Unable to extract the archive.")
					fmt.Println("- archive [" + fileLocation + "]")
					fmt.Println("- directory [" + assetDirectory + "]")
					fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
					return
				}

				zip.Close()
				os.Remove("temp/" + assetKey)
			} else {
				err = mkdir(assetDirectory)
				if err != nil {
					fmt.Println("Error: Unable to create directory for asset.")
					fmt.Println("- directory [" + assetDirectory + "]")
					fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
					return
				}

				err = os.Rename("temp/assetKey", assetFilePath)
				if err != nil {
					fmt.Println("Error: Unable to move asset from temp directory.")
					fmt.Println("Additional details: " + fmt.Sprintf("%v", err))
					return
				}
			}

			addCompletedDownload(assetKey)
		}
	}

	os.Remove("temp")
	fmt.Println("\n> Done!")
}

func unmarshal(body *http.Response, target interface{}) error {
	defer body.Body.Close()
	return json.NewDecoder(body.Body).Decode(target)
}

func loadCompletedDownloads() {
	completedDownloadsJSON, err := ioutil.ReadFile("state.dat")
	if err == nil {
		_ = json.Unmarshal(completedDownloadsJSON, &completedDownloads)
	}
}
func isCompletedDownload(asset string) bool {
	for assetN := range completedDownloads.Assets {
		if completedDownloads.Assets[assetN] == asset {
			return true
		}
	}
	return false
}
func addCompletedDownload(asset string) {
	completedDownloads.Assets = append(completedDownloads.Assets, asset)
	completedDownloadsJSON, _ := json.MarshalIndent(completedDownloads, "", "\t")
	_ = ioutil.WriteFile("state.dat", completedDownloadsJSON, 644)
}

func getDots(assetName string, longestAssetName int) string {
	if longestAssetName-len(assetName) < 0 {
		return strings.Repeat(".", int(math.Abs(float64(longestAssetName-len(assetName)))))
	} else if longestAssetName-len(assetName) == 0 {
		return ""
	}
	return strings.Repeat(".", longestAssetName-len(assetName))
}
func percentOf(current int, all int) float64 {
	percent := (float64(current) * float64(100)) / float64(all)
	return percent
}

func fileDownload(url, location string) (string, string, error) {
	err := mkdir(filepath.Dir(location))
	if err != nil {
		return "", "", err
	}

	file, err := os.Create(location)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", "", err
	}
	file.Close()

	file, err = os.Open(location)
	if err != nil {
		return "", "", err
	}
	md5HashInterface := md5.New()
	_, err = io.Copy(md5HashInterface, file)
	if err != nil {
		return "", "", err
	}
	md5HashBytes := md5HashInterface.Sum(nil)[:16]
	file.Close()

	return location, fmt.Sprintf("%x", md5HashBytes), nil
}

func mkdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 777)
		if err != nil {
			return err
		}
	}
	return nil
}
