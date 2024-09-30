package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FormData struct {
	Name           string `form:"name"`
	Phone          string `form:"phone"`
	Email          string `form:"email"`
	Course         string `form:"course"`
	College        string `form:"college"`
	School         string `form:"school"`
	Project1       string `form:"project1"`
	ProjectDetails string `form:"projectDetails"`
}

func main() {
	r := chi.NewRouter()
	http.HandleFunc("/submit-form", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			data := FormData{
				Name:           r.FormValue("name"),
				Phone:          r.FormValue("phone"),
				Email:          r.FormValue("email"),
				Course:         r.FormValue("course"),
				College:        r.FormValue("college"),
				School:         r.FormValue("school"),
				Project1:       r.FormValue("project1"),
				ProjectDetails: r.FormValue("projectDetails"),
			}

			// Create a text file with the data
			fileName := "formdata.txt"
			content := fmt.Sprintf("Name: %s\nPhone: %s\nEmail: %s\nCourse: %s\nCollege: %s\nSchool: %s\nProject 1: %s\nProject Details: %s\n",
				data.Name, data.Phone, data.Email, data.Course, data.College, data.School, data.Project1, data.ProjectDetails)
			ioutil.WriteFile(fileName, []byte(content), 0644)

			// Send the file to the user
			w.Header().Set("Content-Disposition", "attachment; filename=formdata.txt")
			w.Header().Set("Content-Type", "text/plain")
			http.ServeFile(w, r, fileName)

			// Delete the file after sending
			os.Remove(fileName)
			return
		}
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	})

	// Middleware
	r.Use(middleware.Logger)

	// Serve static files from the "build" directory
	r.Handle("/build/*", http.StripPrefix("/build/", http.FileServer(http.Dir("build"))))

	// Serve static files from the "jsm" directory
	r.Handle("/jsm/*", http.StripPrefix("/jsm/", http.FileServer(http.Dir("jsm"))))

	// Serve static files from the "textures" directory
	r.Handle("/textures/*", http.StripPrefix("/textures/", http.FileServer(http.Dir("textures"))))

	// Serve the main HTML file
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Hello().Render(r.Context(), w)
	})

	// Handle file download
	r.Get("/download", func(w http.ResponseWriter, r *http.Request) {
		texFilePath := "C:\\Users\\om\\Desktop\\comp\\abhi\\om.tex"   // Path to your .tex file
		pdfFilePath := "C:\\Users\\om\\Desktop\\comp\\abhi\\temp.pdf" // Path for the generated PDF

		// Set the command and its working directory
		cmd := exec.Command("pdflatex", texFilePath)
		cmd.Dir = "C:\\Users\\om\\Desktop\\comp\\abhi" // Directory where the .tex file is located

		// Capture the output and error
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert .tex to .pdf: %v\nOutput: %s", err, output), http.StatusInternalServerError)
			return
		}

		// Check if the PDF was created
		if _, err := os.Stat(pdfFilePath); os.IsNotExist(err) {
			http.Error(w, "PDF file was not created", http.StatusInternalServerError)
			return
		}

		// Set headers and serve the PDF file
		w.Header().Set("Content-Disposition", "attachment; filename=file.pdf")
		http.ServeFile(w, r, pdfFilePath)

		// Clean up the PDF file after serving
		defer os.Remove(pdfFilePath)
	})

	r.Get("/clicked", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Button clicked")
	})

	// Start the server
	http.ListenAndServe(":3000", r)
}
