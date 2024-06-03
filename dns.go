package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Altera para milissegundos
func microsToMillis(micros int64) float64 {
	return float64(micros) / 1000
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dnsServers := map[string]string{
			"1.1.1.1":        "Cloudflare DNS",
			"8.8.8.8":        "Google DNS",
			"208.67.222.222": "OpenDNS",
			"189.45.192.3":   "Unifique",
		}

		var bestProvider string
		var bestLatency time.Duration

		fmt.Fprintf(w, `
            <!DOCTYPE html>
            <html>
            <head>
                <title>Análise de Provedores de DNS</title>
                <style>
                    body {
                        font-family: Arial, sans-serif;
                        margin: 20px;
                    }
                    .dns-server {
                        display: flex;
                        align-items: center;
                        margin-bottom: 10px;
                    }
                    .dns-server img {
                        max-width: 50px;
                        margin-right: 10px;
                    }
                    .dns-server p {
                        margin: 0;
                    }
                </style>
				<meta http-equiv="refresh" content="60">
            </head>
            <body>
                <h1>Análise de Provedores de DNS</h1>
                <div>
        `)

		for ip, name := range dnsServers {
			start := time.Now()
			_, err := net.LookupHost("microsoft.com")
			elapsed := time.Since(start)

			if err == nil && (bestProvider == "" || elapsed < bestLatency) {
				bestProvider = name
				bestLatency = elapsed
			}

			var logoURL string
			switch ip {
			case "1.1.1.1":
				logoURL = "https://boostmypresta.com/109-large_default/configuration-cloudflare.jpg"
			case "8.8.8.8":
				logoURL = "https://upload.wikimedia.org/wikipedia/commons/thumb/2/2f/Google_2015_logo.svg/1200px-Google_2015_logo.svg.png"
			case "208.67.222.222":
				logoURL = "https://itcurated.com/infosecindex/wp-content/uploads/sites/35/2016/12/opendns-cybersecurity.png"
			case "189.45.192.3":
				logoURL = "https://servicos.unifique.com.br/images/logo2.png"
			}

			fmt.Fprintf(w, `
                <div class="dns-server">
                    <img src="%s" alt="%s Logo">
                    <p>%s (%s): %.3f ms</p>
                </div>
            `, logoURL, name, name, ip, microsToMillis(elapsed.Microseconds()))
		}

		fmt.Fprintf(w, `
                </div>
                <h2>Melhor Provedor de DNS</h2>
                <p>Melhor provedor neste momento: %s: %.3f ms</p>
            </body>
            </html>
        `, bestProvider, microsToMillis(bestLatency.Microseconds()))
	})

	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
