package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"transitrhythm/gtfs/realtime/server/process"
)

var (
	p = fmt.Println
	l = log.Println
)

const (
	ntpSite                 = "time.google.com" //"0.beevik-ntp.pool.ntp.org"
	referenceTimeZone       = "GMT"
	defaultDownloadInterval = time.Duration(time.Second * 1)
)

// Convert reference timestamp into local time
func localTimeFromReferenceTimestamp(timestamp string) (time.Time, error) {
	var localTime time.Time
	referenceTime, err := time.Parse(time.RFC1123, timestamp)
	if err == nil {
		localTime = referenceTime.Local()
	}
	return localTime, err
}

func fileUpdate(url string, filespec string, sinceTime string) (response *http.Response, written int64, lastModified string, err error) {
	// Check if source file has been modified since the file last modified time.
	response, err = getResponse(url, sinceTime)
	if err == nil && response != nil {
		lastModified = response.Header.Get("Last-Modified")
		written, err = SaveFile(response, filespec, lastModified)
	}
	return response, written, lastModified, err
}

func getResponse(url string, sinceTime string) (*http.Response, error) {
	// Determine Last-Modified-Time for source file
	var response *http.Response
	response = nil
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err == nil {
		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}
		//defer response.Body.Close()
		if response.StatusCode == http.StatusOK {
			lastCacheModifiedTime := response.Header.Get("Last-Modified")
			cacheTime, _ := time.Parse(time.RFC1123, lastCacheModifiedTime)
			fileTime, _ := time.Parse(time.RFC1123, sinceTime)
			if fileTime.Sub(cacheTime) < 0 {
				request, err = http.NewRequest(http.MethodGet, url, nil)
				if err == nil {
					// Add date flag to check for file changes since the last download
					request.Header.Set("If-Modified-Since", sinceTime)
					response, err = client.Do(request)
					if err != nil {
						return nil, err
					}
					if response.StatusCode == http.StatusOK {
						return response, err
					}
					//					defer response.Body.Close()
				}
			}
		}
	}
	return response, err
}

// SaveFile copies the file in the body of the reponse into a local file repository
func SaveFile(response *http.Response, filespec string, timestamp string) (written int64, err error) {
	// Create the file
	out, err := os.Create(filespec)
	if err != nil {
		l("File open error: ", err)
		return written, err
	}
	defer out.Close()

	// Write the body to file
	written, err = io.Copy(out, response.Body)
	if err != nil {
		l("File write error: ", err)
		return written, err
	}
	//defer response.Body.Close()

	ftime, err := time.Parse(time.RFC1123, timestamp)
	// change both atime and mtime to lastModifiedTime
	err = os.Chtimes(filespec, ftime, ftime)
	if err != nil {
		l("File timestamp error: ", err)
	}
	return written, err
}

// DownloadFile will download a url-specified cached file to a local file.
func DownloadFile(dst, src string) {
	var response *http.Response
	for {
		// Does destination data file exist?
		var updatedFileModifiedTime string
		var written int64
		// If file exists, then convert local file timestamp into GMT string
		file, err := os.Stat(dst)
		if err == nil {
			lastFileModifiedTime, err := fileModifiedTime(file, referenceTimeZone)
			if err == nil {
				response, written, updatedFileModifiedTime, err = fileUpdate(src, dst, lastFileModifiedTime)
				if err == nil {
					cacheTime, err := localTimeFromReferenceTimestamp(updatedFileModifiedTime)
					if err == nil {
						latency := time.Since(cacheTime)
						p("a. Now =", time.Now(), "; Filepath =", dst, "; Written", written, "; Last modified =", cacheTime, "; Latency:", latency)
					}
				}
			}
		} else {
			// Otherwise, get the data & create a new file
			response, err = http.Get(src)
			if err == nil {
				updatedFileModifiedTime = response.Header.Get("Last-Modified")
				written, _ = SaveFile(response, dst, updatedFileModifiedTime)
				cacheTime, err := localTimeFromReferenceTimestamp(updatedFileModifiedTime)
				if err == nil {
					latency := time.Since(cacheTime)
					p("b. Now =", time.Now(), "; Filepath =", dst, "; Written", written, "; Last modified =", cacheTime, "; Latency:", latency)
				}
			}
		}
		body, err := ioutil.ReadAll(response.Body)
		go process.Process(body, len(body))
		//	response.Body.Close()

		//waitDuration, err := loopDuration(updatedFileModifiedTime, defaultDownloadInterval)
		//p("WaitDuration:", waitDuration)
		time.Sleep(defaultDownloadInterval)
	}
}

// Convert local file timestamp into GMT string
func fileModifiedTime(file os.FileInfo, timeZone string) (string, error) {
	lastModifiedFileTime := file.ModTime()
	location, err := time.LoadLocation(timeZone)
	return lastModifiedFileTime.In(location).Format(time.RFC1123), err
}
