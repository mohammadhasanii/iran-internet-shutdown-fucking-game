package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>iran-internet-shutdown-fucking-game</title>
<style>
  * { margin:0; padding:0; box-sizing:border-box; }
  body {
    background: #0a0a0a;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    font-family: 'Courier New', monospace;
    user-select: none;
    overflow: hidden;
  }
  .game-container {
    text-align: center;
  }
  .title {
    color: #777;
    font-size: 0.9rem;
    margin-bottom: 1.5rem;
    letter-spacing: 1px;
    text-transform: uppercase;
    opacity: 0.8;
  }
  .counter {
    font-size: 20vw;
    color: #ff3333;
    text-shadow: 0 0 30px #ff0000, 0 0 60px #ff0000;
    animation: pulse 0.8s infinite;
  }
  .counter.paused {
    color: #cccccc;
    text-shadow: 0 0 20px #888;
    animation: none;
  }
  @keyframes pulse {
    0% { text-shadow: 0 0 20px #ff0000, 0 0 40px #ff0000; }
    50% { text-shadow: 0 0 40px #ff4444, 0 0 80px #ff4444; }
    100% { text-shadow: 0 0 20px #ff0000, 0 0 40px #ff0000; }
  }
  .goal {
    font-size: 1.5rem;
    color: #aaa;
    margin-bottom: 1rem;
    letter-spacing: 2px;
  }
  .message {
    font-size: 2rem;
    height: 3rem;
    color: #ffcc00;
    margin-top: 1rem;
    opacity: 0;
    transition: opacity 0.3s;
  }
  .message.show {
    opacity: 1;
  }
  .score {
    position: absolute;
    top: 20px;
    right: 20px;
    color: #fff;
    font-size: 1.2rem;
  }
  .instruction {
    position: absolute;
    bottom: 20px;
    width: 100%;
    text-align: center;
    color: #555;
    font-size: 1rem;
  }
</style>
</head>
<body>
  <div class="score">Perfect Stops: <span id="score">0</span></div>
  <div class="game-container">
    <div class="title">iran-internet-shutdown-fucking-game</div>
    <div class="goal">Stop at <strong>10:00</strong></div>
    <div id="number" class="counter">00:00</div>
    <div id="message" class="message"></div>
  </div>
  <div class="instruction">Press ANY key to <strong id="action-text">PAUSE</strong></div>

  <script>
    (function() {
      const numberEl = document.getElementById('number');
      const messageEl = document.getElementById('message');
      const scoreEl = document.getElementById('score');
      const actionText = document.getElementById('action-text');
      const counterEl = document.querySelector('.counter');

      let total = 0;            // 0 to 9999 (MMSS)
      let running = true;
      let intervalId = null;
      let score = 0;
      const speed = 10;         // ms per tick — 100 ticks per second, seconds cycle fast!
      const targetTotal = 1000; // 10:00

      function formatTime(val) {
        const mins = Math.floor(val / 100);
        const secs = val % 100;
        return String(mins).padStart(2, '0') + ':' + String(secs).padStart(2, '0');
      }

      function increment() {
        total = (total + 1) % 10000; // 0000 to 9999
        numberEl.textContent = formatTime(total);
      }

      function startCounting() {
        running = true;
        counterEl.classList.remove('paused');
        actionText.textContent = 'PAUSE';
        messageEl.classList.remove('show');
        if (intervalId) clearInterval(intervalId);
        intervalId = setInterval(increment, speed);
      }

      function stopCounting() {
        running = false;
        counterEl.classList.add('paused');
        actionText.textContent = 'RESTART';
        if (intervalId) {
          clearInterval(intervalId);
          intervalId = null;
        }
        if (total === targetTotal) {
          messageEl.textContent = '🎯 PERFECT! +1';
          messageEl.style.color = '#00ff88';
          score++;
          scoreEl.textContent = score;
        } else {
          messageEl.textContent = '❌ Missed! Try again.';
          messageEl.style.color = '#ff5555';
        }
        messageEl.classList.add('show');
        setTimeout(function() {
          messageEl.classList.remove('show');
        }, 800);
      }

      function resetAndStart() {
        total = 0;
        numberEl.textContent = '00:00';
        startCounting();
      }

      function handleKeyPress(e) {
        e.preventDefault();
        if (running) {
          stopCounting();
        } else {
          resetAndStart();
        }
      }

      startCounting();
      document.addEventListener('keydown', handleKeyPress);
      document.addEventListener('keyup', function(e) { e.preventDefault(); });
    })();
  </script>
</body>
</html>`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(indexHTML))
	})

	srv := &http.Server{Addr: ":9999"}
	go func() {
		fmt.Println("💀 iran-internet-shutdown-fucking-game -> http://localhost::9999")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server crashed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n🔌 Shutdown complete. Even the game is gone — like our internet.")
}
