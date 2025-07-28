# PDF Outline Extractor

This project is a Go application designed to process PDF files and extract structured data, including the title and outline (e.g., headings), saving the results as JSON files. It uses poppler-utils for text extraction and processes multiple PDFs concurrently.

## Features

- *PDF Text Extraction*: Uses `pdftotext` from poppler-utils to extract text while preserving layout.
- *Title Detection*: Dynamically identifies the document title using heuristics like prominence and repetition.
- *Outline Extraction*: Detects headings based on patterns (e.g., "Section 1", all caps) and indentation levels.
- *Concurrent Processing*: Processes multiple PDF files in parallel using Go's concurrency features.
- *JSON Output*: Saves extracted data (title and outline with page numbers) as JSON files.

## Prerequisites

- *Go*: Version 1.24.4 or later.
- *poppler-utils*: Required for `pdftotext` to extract text from PDFs.
- *Docker*: Optional, for running the application in a container.

## Installation

1.  **Clone the Repository**:
    ```bash
    git clone <repository-url>
    cd pdf-outline-extractor
    ```

2.  **Install Dependencies**:
    Ensure `poppler-utils` is installed. On Debian-based systems:
    ```bash
    sudo apt-get update
    sudo apt-get install -y poppler-utils
    ```

3.  **Build the Application**:
    ```bash
    go build -o main .
    ```

4.  **Run with Docker** (alternative):
    ```bash
    docker build -t pdf-outline-extractor .
    docker run -v $(pwd)/input:/app/input -v $(pwd)/output:/app/output pdf-outline-extractor
    ```

## Usage

1.  **Prepare Input PDFs**:
    Place your PDF files in the `/app/input` directory (or a directory named `input` in the project root).

2.  **Run the Application**:
    ```bash
    ./main
    ```
    The application processes all `.pdf` files in the input directory and generates corresponding `.json` files in the output directory.

3.  **Output**:
    Each PDF is converted to a JSON file containing:
    - `title`: The extracted document title.
    - `outline`: A list of headings with their level (e.g., H1, H2, H3) and page number.

    **Example JSON output:**
    ```json
    {
        "title": "Sample Document",
        "outline": [
            {
                "level": "H1",
                "text": "Introduction",
                "page": 1
            },
            {
                "level": "H2",
                "text": "Section 1: Overview",
                "page": 2
            }
        ]
    }
    ```

## Project Structure

-   `main.go`: Core application logic for processing PDFs and generating JSON output.
-   `Dockerfile`: Defines the Docker image with Go and `poppler-utils`.
-   `input/`: Directory for input PDF files.
-   `output/`: Directory for generated JSON files.

## How It Works

1.  **PDF Processing**:
    -   Uses `pdftotext` to extract text from PDFs, splitting content by page.
    -   Extracts the title using heuristics (e.g., prominent text on the first page or repeated headers).

2.  **Outline Extraction**:
    -   Identifies headings using regex patterns (e.g., "Section 1", all caps) and indentation.
    -   Assigns heading levels (H1, H2, H3) based on indentation and text clues.
    -   Tracks page numbers for each heading.

3.  **Concurrency**:
    -   Processes multiple PDFs concurrently using Go's `sync.WaitGroup` and goroutines.
    -   Ensures thread-safe error handling and file writing.

4.  **Output**:
    -   Saves structured data as JSON files in the output directory.

## Limitations

-   *Title Detection*: Relies on heuristics, which may not always identify the correct title.
-   *Heading Detection*: Depends on regex and indentation, which may miss complex or non-standard headings.
-   *PDF Compatibility*: Assumes PDFs are text-based; scanned PDFs require OCR preprocessing.
-   *Error Handling*: Errors during processing (e.g., corrupted PDFs) are logged but may halt processing for that file.

## Contributing

Contributions are welcome! To contribute:
1.  Fork the repository.
2.  Create a feature branch (`git checkout -b feature/your-feature`).
3.  Commit your changes (`git commit -m "Add your feature"`).
4.  Push to the branch (`git push origin feature/your-feature`).
5.  Open a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
