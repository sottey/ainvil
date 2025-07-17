package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type ExportFile struct {
	Pendant  string
	Date     string
	FileName string
	FullPath string
	WebPath  string
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local web server to browse Ainvil export files",
	Run: func(cmd *cobra.Command, args []string) {
		outDir, _ := cmd.Flags().GetString("out")
		port, _ := cmd.Flags().GetInt("port")

		files := []ExportFile{}
		_ = filepath.Walk(outDir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
				relPath, _ := filepath.Rel(outDir, path)
				parts := strings.Split(relPath, string(os.PathSeparator))
				if len(parts) >= 4 {
					date := fmt.Sprintf("%s-%s-%s", parts[0], parts[1], parts[2])
					name := parts[3]
					pendant := strings.SplitN(name, "_", 2)[0]
					files = append(files, ExportFile{
						Pendant:  pendant,
						Date:     date,
						FileName: name,
						FullPath: path,
						WebPath:  "/view?file=" + relPath,
					})
				}
			}
			return nil
		})

		sort.Slice(files, func(i, j int) bool {
			return files[i].Date > files[j].Date
		})

		tmpl := template.Must(template.New("index").Parse(`
        <!DOCTYPE html>
        <html>
        <head><title>Ainvil Export Viewer</title></head>
        <body style="font-family: sans-serif;">
        <h1>Ainvil Export Files</h1>
        <ul>
        {{ range . }}
            <li><a href="{{ .WebPath }}">{{ .Date }} - {{ .Pendant }} - {{ .FileName }}</a></li>
        {{ end }}
        </ul>
        </body>
        </html>
        `))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl.Execute(w, files)
		})

		http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
			file := r.URL.Query().Get("file")
			if file == "" {
				http.Error(w, "Missing file param", 400)
				return
			}
			full := filepath.Join(outDir, file)
			data, err := os.ReadFile(full)
			if err != nil {
				http.Error(w, "Error reading file", 500)
				return
			}
			var pretty map[string]interface{}
			if err := json.Unmarshal(data, &pretty); err != nil {
				http.Error(w, "Invalid JSON", 500)
				return
			}
			out, _ := json.MarshalIndent(pretty, "", "  ")
			w.Header().Set("Content-Type", "application/json")
			w.Write(out)
		})

		fmt.Printf("Serving at http://localhost:%d ...\n", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	},
}

func init() {
	flagSet := serveCmd.Flags()
	flagSet.String("out", "./out", "Output directory to serve files from")
	flagSet.Int("port", 8080, "Port to serve on")
	rootCmd.AddCommand(serveCmd)
}
