package main

import "bufio"
import "github.com/beberlei/hhvm-serve/fcgiclient"
import "errors"
import "net/http"
import "fmt"
import "io/ioutil"
import "flag"
import "os"
import "strings"
import "strconv"

var documentRoot string
var index string
var listen string
var staticHandler *http.ServeMux
var serverPort int
var serverIp string
var serverEnvironment map[string]string

func respond(w http.ResponseWriter, body string, statusCode int, headers map[string]string) {
	w.WriteHeader(statusCode)
	for header, value := range headers {
		w.Header().Set(header, value)
	}
	fmt.Fprintf(w, "%s", body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqParams := ""
	var filename string
	var scriptName string

	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		reqParams = string(body)
	}

	if r.URL.Path == "/.env" {
		respond(w, "Not allowed", 403, map[string]string{})
		return
	} else if r.URL.Path == "/" || r.URL.Path == "" {
		scriptName = "/" + index
		filename = documentRoot + "/" + index
	} else {
		scriptName = r.URL.Path
		filename = documentRoot + r.URL.Path
	}

	// static file exists
	_, err := os.Stat(filename)
	if !strings.HasSuffix(filename, ".php") && err == nil {
		staticHandler.ServeHTTP(w, r)
		return
	}

	if os.IsNotExist(err) {
		scriptName = "/" + index
		filename = documentRoot + "/" + index
	}

	env := make(map[string]string)

	for name,value := range serverEnvironment {
		env[name] = value
	}

	env["REQUEST_METHOD"] = r.Method
	env["SCRIPT_FILENAME"] = filename
	env["SCRIPT_NAME"] = scriptName
	env["SERVER_SOFTWARE"] = "go / fcgiclient "
	env["REMOTE_ADDR"] = r.RemoteAddr
	env["SERVER_PROTOCOL"] = "HTTP/1.1"
	env["PATH_INFO"] = r.URL.Path
	env["DOCUMENT_ROOT"] = documentRoot
	env["QUERY_STRING"] = r.URL.RawQuery
	env["REQUEST_URI"] = r.URL.Path + "?" + r.URL.RawQuery
	//env["HTTP_HOST"] = r.Host
	//env["SERVER_ADDR"] = listen

	for header, values := range r.Header {
		env["HTTP_" + strings.Replace(strings.ToUpper(header), "-", "_", -1)] = values[0]
	}

	fcgi, err := fcgiclient.New(serverIp, serverPort)
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	content, _, err := fcgi.Request(env, reqParams)

	if err != nil {
		fmt.Printf("ERROR: %s - %v", r.URL.Path, err)
	}

	statusCode, headers, body, err := ParseFastCgiResponse(fmt.Sprintf("%s", content))

	respond(w, body, statusCode, headers)

	fmt.Printf("%s \"%s %s %s\" %d %d \"%s\"\n", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, statusCode, len(content), r.UserAgent())
}

func ParseFastCgiResponse(content string) (int, map[string]string, string, error) {
	var headers map[string]string

	parts := strings.SplitN(content, "\r\n\r\n", 2)

	if len(parts) < 2 {
		return 502, headers, "", errors.New("Cannot parse FastCGI Response")
	}

	headerParts := strings.Split(parts[0], ":")
	body := parts[1]
	status := 200

	if strings.HasPrefix(headerParts[0], "Status:") {
		lineParts := strings.SplitN(headerParts[0], " ", 3)
		status, _ = strconv.Atoi(lineParts[1])
	}

	for _, line := range headerParts {
		lineParts := strings.SplitN(line, ":", 2)

		if len(lineParts) < 2 {
			continue
		}

		lineParts[1] = strings.TrimSpace(lineParts[1])

		if lineParts[0] == "Status" {
			continue
		}

		headers[lineParts[0]] = lineParts[1]
	}

	return status, headers, body, nil
}

func ReadEnvironmentFile(path string) {
	file, err := os.Open(path + "/.env")

	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	serverEnvironment = make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if (strings.Contains(line, "=")) {
			parts := strings.Split(line, "=")
			serverEnvironment[parts[0]] = parts[1]
		}
	}
}

func main() {

	cwd, _ := os.Getwd()
	flag.StringVar(&documentRoot, "document-root", cwd, "The document root to serve files from")
	flag.StringVar(&listen, "listen", "localhost:8080", "The webserver bind address to listen to.")
	flag.StringVar(&serverIp, "server", "127.0.0.1", "The FastCGI Server to listen to")
	flag.IntVar(&serverPort, "server-port", 9000, "The FastCGI Port to listen to")
	flag.StringVar(&index, "index", "index.php", "The default script to call when path cannot be served by existing file.")

	flag.Parse()

	ReadEnvironmentFile(cwd)

	staticHandler = http.NewServeMux()
	staticHandler.Handle("/", http.FileServer(http.Dir(documentRoot)))

	fmt.Printf("Listening on http://%s\n", listen)
	fmt.Printf("Document root is %s\n", documentRoot)
	fmt.Printf("Press Ctrl-C to quit.\n")

	http.HandleFunc("/", handler)
	http.ListenAndServe(listen, nil)
}
