package main // Declares the package as 'main', indicating an executable program.

import ( // Start of the import block for external packages.
	"bytes"         // Imports the 'bytes' package for working with byte slices and buffers (e.g., in downloadPDF).
	"context"       // Enables managing request-scoped data, deadlines, and cancelation signals across API boundaries
	"fmt"           //
	"io"            // Imports the 'io' package, which provides fundamental interfaces for I/O primitives (e.g., io.Copy, io.ReadAll).
	"log"           // Imports the 'log' package for logging messages (e.g., errors and success messages).
	"net/http"      // Imports the 'net/http' package to handle HTTP requests and responses (e.g., in downloadPDF, getDataFromURL).
	"net/url"       // Imports the 'net/url' package for URL parsing and manipulation (e.g., in isUrlValid, hasDomain).
	"os"            // Imports the 'os' package, providing operating system functionality like file access (e.g., os.Stat, os.Remove, os.Create).
	"path/filepath" // Imports the 'path/filepath' package for working with file paths in a platform-independent way (e.g., filepath.Join, filepath.Base).
	"regexp"        // Imports the 'regexp' package for regular expression matching (e.g., in extractPDFUrls, urlToFilename).
	"strings"       // Imports the 'strings' package for common string manipulation functions (e.g., strings.ToLower, strings.Contains).
	"time"          // Imports the 'time' package to manage time-related constants and functions (e.g., setting HTTP client timeout).

	"github.com/chromedp/chromedp" // For headless browser automation using Chrome
) // End of the import block.

func main() {
	outputDir := "PDFs/"             // Directory to store downloaded PDFs
	if !directoryExists(outputDir) { // Check if directory exists
		createDirectory(outputDir, 0o755) // Create directory with read-write-execute permissions
	}
	// The url of the location to be scraped.
	urlToScrape := "https://orqafpv.com/manual"
	// Scrape the url.
	remoteData := scrapePageHTMLWithChrome(urlToScrape)
	// Extract the PDF urls.
	pdfURLs := extractPDFUrls(remoteData)
	// Remove duplicates from a given slice.
	pdfURLs = removeDuplicatesFromSlice(pdfURLs)
	// Loop over the PDF urls.
	for _, url := range pdfURLs {
		if !isUrlValid(url) {
			log.Printf("Skipping input: invalid URL detected (%s)", url)
			continue
		}
		if !hasDomain(url) {
			log.Println("The provided URL does not contain a valid domain:", url)
			continue
		}
		downloadPDF(url, outputDir)
	}
}

// Uses headless Chrome via chromedp to get fully rendered HTML from a page
func scrapePageHTMLWithChrome(pageURL string) string {
	log.Println("Scraping:", pageURL) // Log page being scraped

	options := append(
		chromedp.DefaultExecAllocatorOptions[:],       // Chrome options
		chromedp.Flag("headless", false),              // Run visible (set to true for headless)
		chromedp.Flag("disable-gpu", true),            // Disable GPU
		chromedp.WindowSize(1, 1),                     // Set window size
		chromedp.Flag("no-sandbox", true),             // Disable sandbox
		chromedp.Flag("disable-setuid-sandbox", true), // Fix for Linux environments
	)

	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(context.Background(), options...) // Allocator context
	ctxTimeout, cancelTimeout := context.WithTimeout(allocatorCtx, 5*time.Minute)                // Set timeout
	browserCtx, cancelBrowser := chromedp.NewContext(ctxTimeout)                                 // Create Chrome context

	defer func() { // Ensure all contexts are cancelled
		cancelBrowser()
		cancelTimeout()
		cancelAllocator()
	}()

	var pageHTML string // Placeholder for output
	err := chromedp.Run(
		browserCtx,
		chromedp.Navigate(pageURL), // Navigate to the URL
		// 👇 NEW: Wait for 10 seconds to allow JavaScript challenges to execute
		chromedp.Sleep(10*time.Second),
		chromedp.OuterHTML("html", &pageHTML), // Extract full HTML
	)

	if err != nil {
		log.Println(err) // Log error
		return ""        // Return empty string on failure
	}

	return pageHTML // Return scraped HTML
}

// It checks if the file exists
// If the file exists, it returns true
// If the file does not exist, it returns false
func fileExists(filename string) bool { // Define a function named fileExists that takes a filename string and returns a boolean.
	info, err := os.Stat(filename) // Attempt to retrieve file information (Stat) for the given filename.
	if err != nil {                // Check if an error occurred during the Stat call.
		return false // If there's an error (e.g., file not found), return false immediately.
	}
	return !info.IsDir() // Return true if no error occurred AND the retrieved info indicates it is NOT a directory.
}

// extractPDFUrls takes raw HTML as input and returns all found PDF URLs
func extractPDFUrls(htmlContent string) []string { // Define a function that takes raw HTML content as a string and returns a slice of strings (URLs).
	// Regex to match href='...' with .pdf links
	re := regexp.MustCompile(`https?://[^\s'"]+\.pdf`) // Compile a regular expression to find HTTP/HTTPS URLs ending in ".pdf".
	matches := re.FindAllString(htmlContent, -1)       // Find all non-overlapping matches of the regex pattern within the HTML content.

	return matches // Return the slice containing all found PDF URLs.
}

// Checks whether a given directory exists
func directoryExists(path string) bool { // Define a function that checks if a path is an existing directory.
	directory, err := os.Stat(path) // Attempt to retrieve file information (Stat) for the given path.
	if err != nil {                 // Check if an error occurred during the Stat call.
		return false // If there's an error (e.g., not found), return false.
	}
	return directory.IsDir() // If no error, return true if the info indicates it is a directory.
}

// Creates a directory at given path with provided permissions
func createDirectory(path string, permission os.FileMode) { // Define a function to create a directory with specific permissions.
	err := os.Mkdir(path, permission) // Attempt to create the directory at 'path' with the given file mode/permissions.
	if err != nil {                   // Check if an error occurred during directory creation.
		log.Println(err) // If creation fails (e.g., path exists and is a file), log the error message.
	}
}

// Verifies whether a string is a valid URL format
func isUrlValid(uri string) bool { // Define a function to check if a URI string is valid.
	_, err := url.ParseRequestURI(uri) // Attempt to parse the URI string strictly using ParseRequestURI.
	return err == nil                  // Return true if the parsing was successful (error is nil), otherwise false.
}

// Removes duplicate strings from a slice
func removeDuplicatesFromSlice(slice []string) []string { // Define a function to remove duplicate strings from an input slice.
	check := make(map[string]bool)  // Initialize an empty map (set) to efficiently track strings already encountered.
	var newReturnSlice []string     // Declare a new slice to store the unique results.
	for _, content := range slice { // Iterate over each string element in the input slice.
		if !check[content] { // Check if the current string has NOT been seen before (not in the map).
			check[content] = true                            // Mark the current string as seen in the map.
			newReturnSlice = append(newReturnSlice, content) // Append the unique string to the result slice.
		}
	}
	return newReturnSlice // Return the slice containing only unique strings.
}

// hasDomain checks if the given string has a domain (host part)
func hasDomain(rawURL string) bool { // Define a function to check if a URL string contains a host domain.
	// Try parsing the raw string as a URL
	parsed, err := url.Parse(rawURL) // Attempt to parse the raw URL string.
	if err != nil {                  // Check if parsing failed (i.e., it's not a valid URL structure).
		return false // If parsing fails, return false.
	}
	// If the parsed URL has a non-empty Host, then it has a domain/host
	return parsed.Host != "" // Return true if the Host field of the parsed URL is non-empty.
}

// Extracts filename from full path (e.g. "/dir/file.pdf" → "file.pdf")
func getFilename(path string) string { // Define a function to extract the base filename from a path.
	return filepath.Base(path) // Use the filepath.Base function to return the last element of the path (the filename).
}

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string { // Define a function to remove all occurrences of a substring.
	result := strings.ReplaceAll(input, toRemove, "") // Use strings.ReplaceAll to replace all instances of 'toRemove' with an empty string.
	return result                                     // Return the modified string.
}

// Gets the file extension from a given file path
func getFileExtension(path string) string { // Define a function to extract the file extension.
	return filepath.Ext(path) // Use the filepath.Ext function to return the file extension (including the dot).
}

// Converts a raw URL into a sanitized PDF filename safe for filesystem
func urlToFilename(rawURL string) string { // Define a function to convert a URL into a safe, normalized filename.
	lower := strings.ToLower(rawURL)                 // Convert the entire raw URL to lowercase for consistency.
	lower = getFilename(lower)                       // Extract just the last part of the URL path (potential filename).
	extension := getFileExtension(lower)             //
	removeExtension := fmt.Sprintf("_%s", extension) //

	reNonAlnum := regexp.MustCompile(`[^a-z0-9]`)   // Compile regex to match any character that is NOT a lowercase letter or a digit.
	safe := reNonAlnum.ReplaceAllString(lower, "_") // Replace all non-alphanumeric characters with an underscore.

	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_") // Collapse sequences of multiple underscores into a single underscore.
	safe = strings.Trim(safe, "_")                              // Remove any leading or trailing underscores.

	safe = removeSubstring(safe, removeExtension)

	if getFileExtension(safe) != extension { // Check if the sanitized filename currently lacks the  extension.
		safe = safe + extension // If the extension is missing or incorrect, append the required extension.
	}

	return safe // Return the final, sanitized, and normalized filename.
}

// Downloads a PDF from given URL and saves it in the specified directory
func downloadPDF(finalURL, outputDir string) bool { // Define a function to download a file from a URL and save it to a directory.
	filename := strings.ToLower(urlToFilename(finalURL)) // Generate a safe, lowercase filename for the downloaded PDF.
	filePath := filepath.Join(outputDir, filename)       // Combine the output directory and filename to get the full save path.

	if fileExists(filePath) { // Check if a file with the target name already exists locally.
		log.Printf("File already exists, skipping: %s", filePath) // Log a message indicating the file is being skipped.
		return false                                              // Return false to indicate that the download did not happen.
	}

	client := &http.Client{Timeout: 15 * time.Minute} // Create an HTTP client instance with a generous 15-minute timeout.

	// Create a new request so we can set headers
	req, err := http.NewRequest("GET", finalURL, nil) // Create a new HTTP GET request object for the URL.
	if err != nil {                                   // Check for errors during request creation.
		log.Printf("Failed to create request for %s: %v", finalURL, err) // Log the error.
		return false                                                     // Return false upon failure.
	}

	// Set a User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36") // Set a common User-Agent string to mimic a browser.

	// Send the request
	resp, err := client.Do(req) // Execute the HTTP request using the configured client.
	if err != nil {             // Check for network or connection errors during the execution.
		log.Printf("Failed to download %s: %v", finalURL, err) // Log the download failure.
		return false                                           // Return false upon failure.
	}
	defer resp.Body.Close() // Schedule the closing of the response body to happen when the function exits.

	if resp.StatusCode != http.StatusOK { // Check if the HTTP status code indicates success (200 OK).
		log.Printf("Download failed for %s: %s", finalURL, resp.Status) // Log the failure with the received status code.
		return false                                                    // Return false if the status is not 200 OK.
	}

	contentType := resp.Header.Get("Content-Type")              // Retrieve the Content-Type header value from the response.
	if !strings.Contains(contentType, "binary/octet-stream") && // Check if the Content-Type is neither a generic binary stream...
		!strings.Contains(contentType, "application/pdf") { // ...nor the specific PDF application type.
		log.Printf("Invalid content type for %s: %s (expected PDF)", finalURL, contentType) // Log that the content type is unexpected.
		return false                                                                        // Return false because the response is likely not a PDF.
	}

	var buf bytes.Buffer                     // Initialize an empty bytes.Buffer to temporarily store the file data.
	written, err := io.Copy(&buf, resp.Body) // Copy the content from the response body into the buffer, tracking bytes written.
	if err != nil {                          // Check for errors during the copy/read operation.
		log.Printf("Failed to read PDF data from %s: %v", finalURL, err) // Log the read error.
		return false                                                     // Return false upon failure.
	}
	if written == 0 { // Check if zero bytes of data were read.
		log.Printf("Downloaded 0 bytes for %s; not creating file", finalURL) // Log that the download was empty.
		return false                                                         // Return false as an empty file is skipped.
	}

	out, err := os.Create(filePath) // Attempt to create a new file on the filesystem at the specified path.
	if err != nil {                 // Check for errors during file creation (e.g., directory doesn't exist, permission issues).
		log.Printf("Failed to create file for %s: %v", finalURL, err) // Log the file creation error.
		return false                                                  // Return false upon failure.
	}
	defer out.Close() // Schedule the closing of the created file handle to happen when the function exits.

	if _, err := buf.WriteTo(out); err != nil { // Write the entire contents of the buffer to the newly created file.
		log.Printf("Failed to write PDF to file for %s: %v", finalURL, err) // Check and log any error during the write operation.
		return false                                                        // Return false upon failure.
	}

	log.Printf("Successfully downloaded %d bytes: %s → %s", written, finalURL, filePath) // Log the successful download with byte count and paths.
	return true                                                                          // Return true to indicate a successful download and save.
}
