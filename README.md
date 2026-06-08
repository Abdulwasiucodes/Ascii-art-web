# ASCII Art Web

## Description

ASCII Art Web is a web-based GUI application that converts user-provided text into ASCII art representations using one of three banner styles: **standard**, **shadow**, and **thinkertoy**.

The server is written entirely in Go using only the standard library. Users interact with a clean web interface to type their text, choose a banner, and instantly see the ASCII art output rendered on the page.

---

## Authors

- Abdulraufu Wasiu Olamilekan

---

## Usage

### Requirements
- Go 1.18 or higher

### Setup

1. Clone the repository:
```bash
git clone https://github.com/Abdulwasiucodes/Ascii-art-web.git
cd Ascii-art-web
```

2. Make sure the `banners/` folder contains:
   - `standard.txt`
   - `shadow.txt`
   - `thinkertoy.txt`

3. Run the server:
```bash
go run .
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

### Running Tests
```bash
go test ./...
```

---

## Implementation Details

### Algorithm

1. **Banner Loading**: Each banner `.txt` file contains ASCII art representations of all printable characters (ASCII 32–126). Characters are arranged sequentially, each occupying exactly 8 lines, separated by a blank line (9 lines per character block).

2. **Character Mapping**: For a given character `c`, its position in the banner file is calculated as:
   ```
   line_index = (c - 32) * 9 + 1 + row
   ```
   where `row` is 0–7 (the 8 lines of the character block), and the `+1` accounts for the leading blank line at the top of the file.

3. **Multi-line Input**: Input text is split on newline characters (`\n`). Each line of input is rendered independently, and empty lines in the input produce empty lines in the output.

4. **Character Validation**: Only printable ASCII characters (32–126) and newline characters are accepted. Any other character returns a 400 Bad Request error.

5. **HTTP Endpoints**:
   - `GET /` — Serves the main HTML page
   - `POST /ascii-art` — Accepts `text` and `banner` form fields, returns the ASCII art result rendered into the main page

6. **HTTP Status Codes**:
   - `200 OK` — Successful generation
   - `400 Bad Request` — Empty input, invalid banner, or non-ASCII characters
   - `404 Not Found` — Unknown route or missing banner file
   - `500 Internal Server Error` — Unhandled server-side errors