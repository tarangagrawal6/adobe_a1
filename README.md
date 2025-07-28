# PDF Content and Structure Extractor

This project is a Go-based application designed to extract structured information, such as titles and hierarchical outlines, from PDF documents. It reads PDF files from an `input` directory, processes them concurrently, and outputs the extracted data into structured JSON files in the `output` directory.

## Features

* **Document Title Extraction:** Automatically detects and extracts the main title from a PDF.
* **Outline Generation:** Creates a hierarchical outline (e.g., chapters, sections, subsections) based on the document's structure.
* **Concurrent Processing:** Processes multiple PDF files in parallel for improved performance.
* **JSON Output:** Saves the extracted data in a clean, easy-to-parse JSON format.
* **Docker Support:** Includes a `Dockerfile` for easy containerization and deployment, ensuring all dependencies are handled.

## How It Works

The application leverages the `pdftotext` command-line utility (from the `poppler-utils` package) to convert the content of each PDF into text while preserving the layout. The Go program then performs the following steps:

1.  **Text Extraction:** Executes `pdftotext` for each PDF file.
2.  **Page Splitting:** Divides the extracted text into individual pages.
3.  **Title Identification:** Applies heuristics to the first page to find a prominent line of text to use as the title.
4.  **Outline Extraction:** Uses regular expressions and indentation analysis to identify potential headings and determine their hierarchical level (e.g., H1, H2, H3).
5.  **JSON Serialization:** Marshals the extracted title and outline into a JSON object.
6.  **File Output:** Writes the resulting JSON data to a corresponding file in the `output` directory.

## Prerequisites

Before you begin, ensure you have the following installed:

* **Go:** Version 1.24.4 or later.
* **Poppler-utils:** A package that provides PDF rendering and utility tools, including `pdftotext`.
    * On Debian/Ubuntu: `sudo apt-get update && sudo apt-get install poppler-utils`
    * On macOS (using Homebrew): `brew install poppler`
    * On Windows: You can download it via [Chocolatey](https://chocolatey.org/packages/poppler) or install it as part of a WSL environment.

## Setup and Installation

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd <your-repository-directory>
    ```

2.  **Create input and output directories:**
    ```bash
    mkdir input output
    ```

3.  **Place your PDF files** into the `./input/` directory.

## Usage

There are two ways to run this application: directly with Go or using Docker.

### Option 1: Running with Go

Execute the main program from the root of the project directory:

```bash
go run main.go
```

The application will process all PDF files in the `input` directory and save the JSON results in the `output` directory.

### Option 2: Running with Docker

This method is recommended for a consistent and isolated environment, as it handles all dependencies for you.

1.  **Build the Docker image:**
    ```bash
    docker build -t pdf-extractor .
    ```

2.  **Run the Docker container:**
    This command mounts the local `input` and `output` directories into the container.

    ```bash
    docker run --rm -v "$(pwd)/input":/app/input -v "$(pwd)/output":/app/output pdf-extractor
    ```

After the container finishes running, the extracted JSON files will be available in your local `output` directory.

## Output JSON Format

For each input PDF file named `example.pdf`, a corresponding `example.json` file will be created in the `output` directory. The structure of the JSON is as follows:

```json
{
    "title": "The Title of the Document",
    "outline": [
        {
            "level": "H1",
            "text": "Chapter 1: Introduction",
            "page": 1
        },
        {
            "level": "H2",
            "text": "Section 1.1: Background",
            "page": 2
        },
        {
            "level": "H1",
            "text": "Chapter 2: Methodology",
            "page": 5
        }
    ]
}
```

## Code Overview

* `main.go`: The main application file. It handles file I/O, concurrent processing, and orchestrates the PDF-to-JSON conversion.
* `processPDF()`: The core function responsible for processing a single PDF. It calls `pdftotext` and invokes the extraction logic.
* `extractTitle()`: Contains the logic to heuristically determine the document's title.
* `extractOutline()`: Contains the logic to scan pages for headings and build a structured outline.
* `Dockerfile`: Defines the container image, including the Go environment and the `poppler-utils` dependency.
* `go.mod`: Defines the project's Go module and dependencies.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
