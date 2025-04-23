# GoLang Guestbook

This project is a simple guestbook web application written in Go. It was originally created as a class project and has been updated in 2025 with improved knowledge and practices.

## Features

- Submit messages to the guestbook with an optional name.
- View all submitted messages in reverse chronological order.
- Simple and clean HTML interface styled with CSS.
- Messages are stored in the browser's local storage for offline access.

## Updated Project Structure

- `guest.go`: The main Go application file containing the server logic.
- `templates/template.html`: The HTML template for the guestbook interface.
- `static/style.css`: The CSS file for styling the guestbook.
- `static/guestbook.js`: The JavaScript file for handling local storage and form submission.

## How to Run

1. Make sure you have [Go](https://golang.org/) installed on your system.
2. Clone this repository or download the source code.
3. Navigate to the project directory.
4. Run the following command to start the server:
   ```bash
   go run guest.go
   ```
5. Open your browser and go to [http://localhost:8080](http://localhost:8080).

## How to Use

1. Fill in your name (optional) and message in the form.
2. Click "Sign Guestbook" to submit your message.
3. View your message along with others on the main page.
4. Messages are saved in your browser's local storage and will persist even if the page is refreshed.

## Recent Updates

- Modularized the HTML template by moving it to a separate file (`templates/template.html`).
- Improved code maintainability and separation of concerns.

## License

This project is open-source and available under the MIT License.
