package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	uploadDir     = "uploads"
	templateDir   = "static"
	maxUploadSize = 10 << 20
)

type ImageInfo struct {
	Path     string
	Filename string
}

type PageData struct {
	Message          string
	ShowUploadForm   bool
	Images           []ImageInfo
	CurrentStep      string
	OriginalFilePath string
	SelectedMode     int
	SelectedN        int
}

var primitiveModes = []struct {
	ID    int
	Label string
}{
	{ID: 0, Label: "Combo"},
	{ID: 1, Label: "Triangles"},
	{ID: 2, Label: "Rectangles"},
	{ID: 3, Label: "Ellipses"},
}

var nShapeOptions = []int{50, 100, 150, 200}

var funcMap = template.FuncMap{
	"split": func(s, sep string) []string {
		return strings.Split(s, sep)
	},
}

func main() {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/transform", transformHandler)
	http.Handle("/"+uploadDir+"/", http.StripPrefix("/"+uploadDir+"/", http.FileServer(http.Dir(uploadDir))))

	log.Printf("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, PageData{Message: "Upload an image to transform it!", ShowUploadForm: true})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		renderTemplate(w, PageData{Message: "Error: File too large or bad request: " + err.Error(), ShowUploadForm: true})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		renderTemplate(w, PageData{Message: "Error retrieving file: " + err.Error(), ShowUploadForm: true})
		return
	}
	defer file.Close()

	originalFilename := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	uniqueFilename := fmt.Sprintf("%s-%d%s", originalFilename, time.Now().UnixNano(), filepath.Ext(header.Filename))
	uploadedFilePath := filepath.Join(uploadDir, uniqueFilename)

	dst, err := os.Create(uploadedFilePath)
	if err != nil {
		renderTemplate(w, PageData{Message: "Error creating file on server: " + err.Error(), ShowUploadForm: true})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		renderTemplate(w, PageData{Message: "Error saving file: " + err.Error(), ShowUploadForm: true})
		return
	}

	displayTransformationStep(w, r, uploadedFilePath, "selectMode", -1, -1)
}

func transformHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	currentStep := r.FormValue("current_step")
	originalFilePath := r.FormValue("original_file_path")

	selectedModeStr := r.FormValue("selected_mode")
	selectedNStr := r.FormValue("selected_n")

	selectedMode, _ := strconv.Atoi(selectedModeStr)
	selectedN, _ := strconv.Atoi(selectedNStr)

	switch currentStep {
	case "selectMode":
		modeIDStr := r.FormValue("mode_id")
		if modeIDStr == "" {
			displayTransformationStep(w, r, originalFilePath, "selectMode", -1, -1, "Error: No mode selected. Please try again.")
			return
		}
		selectedMode, _ = strconv.Atoi(modeIDStr)
		displayTransformationStep(w, r, originalFilePath, "selectN", selectedMode, -1)

	case "selectN":
		nValStr := r.FormValue("n_value")
		if nValStr == "" {
			displayTransformationStep(w, r, originalFilePath, "selectN", selectedMode, -1, "Error: No number of shapes selected. Please try again.")
			return
		}
		selectedN, _ = strconv.Atoi(nValStr)
		displayTransformationStep(w, r, originalFilePath, "final", selectedMode, selectedN)

	case "downloadFinal":
		finalImagePath := r.FormValue("final_image_path")
		if finalImagePath == "" {
			http.Error(w, "No final image specified for download.", http.StatusBadRequest)
			return
		}

		filename := filepath.Base(finalImagePath)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, finalImagePath)

		log.Println("User downloaded image. Initiating immediate cleanup of 'uploads' directory.")
		if err := emptyUploadsDirectory(uploadDir); err != nil {
			log.Printf("Error emptying uploads directory: %v", err)
		} else {
			log.Println("Uploads directory emptied successfully.")
		}	

	default:
		renderTemplate(w, PageData{Message: "Error: Invalid transformation step.", ShowUploadForm: true})
	}
}

func displayTransformationStep(w http.ResponseWriter, r *http.Request, originalPath string, step string, mode, n int, errMsg ...string) {
	var message string
	if len(errMsg) > 0 && errMsg[0] != "" {
		message = errMsg[0]
	} else {
		switch step {
		case "selectMode":
			message = "Step 1: Choose a transformation mode:"
		case "selectN":
			message = fmt.Sprintf("Step 2: Mode '%s' selected. Now choose the number of shapes (n):", getModeLabel(mode))
		case "final":
			message = "Your transformed image is ready! Click to download."
		}
	}

	data := PageData{
		Message:          message,
		ShowUploadForm:   false,
		CurrentStep:      step,
		OriginalFilePath: originalPath,
		SelectedMode:     mode,
		SelectedN:        n,
	}

	var generatedImages []ImageInfo

	switch step {
	case "selectMode":
		for _, m := range primitiveModes {
			outputFilename := fmt.Sprintf("%s_mode%d-%d.png", strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(filepath.Base(originalPath))), m.ID, time.Now().UnixNano())
			outputPath := filepath.Join(uploadDir, outputFilename)
			err := runPrimitive(originalPath, outputPath, m.ID, 50)
			if err != nil {
				log.Printf("Error running primitive (mode %d): %v", m.ID, err)
				continue
			}
			generatedImages = append(generatedImages, registerImageForDisplay(outputPath)) // Use registerImageForDisplay
		}
	case "selectN":
		for _, nVal := range nShapeOptions {
			outputFilename := fmt.Sprintf("%s_mode%d_n%d-%d.png", strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(filepath.Base(originalPath))), mode, nVal, time.Now().UnixNano())
			outputPath := filepath.Join(uploadDir, outputFilename)
			err := runPrimitive(originalPath, outputPath, mode, nVal)
			if err != nil {
				log.Printf("Error running primitive (mode %d, n %d): %v", mode, nVal, err)
				continue
			}
			generatedImages = append(generatedImages, registerImageForDisplay(outputPath)) // Use registerImageForDisplay
		}
	case "final":
		outputFilename := fmt.Sprintf("%s_final_mode%d_n%d-%d.png", strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(filepath.Base(originalPath))), mode, n, time.Now().UnixNano())
		outputPath := filepath.Join(uploadDir, outputFilename)
		err := runPrimitive(originalPath, outputPath, mode, n)
		if err != nil {
			log.Printf("Error running primitive (final): %v", err)
			http.Error(w, "Error generating final image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		generatedImages = append(generatedImages, registerImageForDisplay(outputPath)) // Use registerImageForDisplay
	}

	data.Images = generatedImages
	renderTemplate(w, data)
}

func renderTemplate(w http.ResponseWriter, data PageData) {
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles(filepath.Join(templateDir, "index.html"))
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func runPrimitive(inputPath, outputPath string, mode, n int) error {
	cmd := exec.Command("primitive",
		"-i", inputPath,
		"-o", outputPath,
		"-m", strconv.Itoa(mode),
		"-n", strconv.Itoa(n),
	)

	log.Printf("Executing command: %s", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("primitive execution failed: %v\nOutput: %s", err, output)
	}
	return nil
}

func getModeLabel(modeID int) string {
	for _, m := range primitiveModes {
		if m.ID == modeID {
			return m.Label
		}
	}
	return "Unknown"
}

func registerImageForDisplay(path string) ImageInfo {
	return ImageInfo{
		Path:     path,
		Filename: filepath.Base(path),
	}
}

func emptyUploadsDirectory(dirPath string) error {
	d, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		filePath := filepath.Join(dirPath, name)
		info, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Error stating file %s during cleanup: %v", filePath, err)
			continue
		}
		if info.Mode().IsRegular() {
			if err = os.Remove(filePath); err != nil {
				log.Printf("Error deleting file %s: %v", filePath, err)
			}
		}
	}
	return nil
}