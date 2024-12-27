
## main.go

The `main.go` file is the main entry point of the application. It sets up the web server, handles HTTP requests, and streams data to the client.

### REST API

The project includes REST API endpoints for updating scores and winners:

- `POST /updateScore`: Updates the score for a player. The request body should be a JSON object with `playerColor` (either "red" or "blue") and `points` (the number of points to add).
- `POST /updateWinner`: Updates the winner status for a player. The request body should be a JSON object with `playerColor` (either "red" or "blue") and `isWinner` (a boolean indicating if the player is the winner).

### Server-Sent Events (SSE)

The project uses Server-Sent Events (SSE) to stream score updates to the client in real-time:

- `GET /streamScore`: Streams the current scores and winner status to the client. The client receives updates every 2 seconds.

### Stream Data

The `streamScoreHandler` function in `main.go` handles the SSE connection and streams data to the client. It sends the current scores and winner status in JSON format.

```go
func streamScoreHandler(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    for {
        if redPlayer.IsWinner || bluePlayer.IsWinner {
            winner := fmt.Sprintf("data: {\"red\": %d, \"blue\": %d, \"winner\": \"%s\"}\n\n", redPlayer.Score, bluePlayer.Score, func() string {
                if redPlayer.IsWinner {
                    return "red"
                } else {
                    return "blue"
                }
            }())
            w.Write([]byte(winner))
            flusher.Flush()
            return
        }

        content := fmt.Sprintf("data: {\"red\": %d, \"blue\": %d}\n\n", redPlayer.Score, bluePlayer.Score)
        w.Write([]byte(content))
        flusher.Flush()
        time.Sleep(2 * time.Second)
    }
}
```

## static/index.html

The `static/index.html` file is the client-side part of the application. It uses HTMX to send requests to the server and update the scores in real-time. It also sets up an SSE connection to receive score updates from the server.

### JavaScript

```javascript
<script>
    // Set up SSE connection
    const eventSource = new EventSource('/streamScore');
    
    eventSource.onmessage = function(event) {
        const data = JSON.parse(event.data);
        
        // Update scores
        document.getElementById('red-score').textContent = data.red;
        document.getElementById('blue-score').textContent = data.blue;
        
        // Check for winner
        if (data.winner) {
            const winnerBanner = document.createElement('div');
            winnerBanner.className = `winner-banner ${data.winner}-winner`;
            winnerBanner.textContent = `${data.winner.toUpperCase()} FIGHTER WINS! 🏆`;
            document.body.appendChild(winnerBanner);
        }
    };
</script>
```

### HTMX

HTMX is used to set up endpoints for updating scores and winners. It allows you to send requests to the server and update parts of the web page without reloading the entire page.

Example usage of HTMX in `index.html`:

```html
<button hx-post="/updateScore" hx-vals='{"playerColor": "red", "points": 1}'>Add Point to Red</button>
<button hx-post="/updateScore" hx-vals='{"playerColor": "blue", "points": 1}'>Add Point to Blue</button>
<button hx-post="/updateWinner" hx-vals='{"playerColor": "red", "isWinner": true}'>Red Wins</button>
<button hx-post="/updateWinner" hx-vals='{"playerColor": "blue", "isWinner": true}'>Blue Wins</button>
```

In this project, you learned how to create a simple web application using Go, REST APIs, Server-Sent Events (SSE), and streaming data. This project demonstrates how to build a real-time score tracking system for a Combat Sport.

