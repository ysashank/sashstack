package handlers

import (
	"html"
	"log"
	"net"
	"net/http"
)

func Waitlist(w http.ResponseWriter, r *http.Request) {

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if !rateLimiter.Allow(ip) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	type waitlistForm struct {
		Email     string
		Country   string
		Frequency string
		Referral  string
		Interest  string
	}
	var data waitlistForm
	data.Email = html.EscapeString(r.FormValue("email"))
	data.Country = html.EscapeString(r.FormValue("country"))
	data.Frequency = html.EscapeString(r.FormValue("frequency"))
	data.Referral = html.EscapeString(r.FormValue("referral"))
	data.Interest = html.EscapeString(r.FormValue("interest"))
	if data.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	// Validate email format
	if !isValidEmail(data.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if data.Country == "" || data.Frequency == "" || data.Referral == "" || data.Interest == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if len(data.Email) > 254 || len(data.Country) > 64 || len(data.Referral) > 128 || len(data.Interest) > 64 {
		http.Error(w, "Input too long", http.StatusBadRequest)
		return
	}

	if r.FormValue("website") != "" {
		log.Printf("⚠️ Bot blocked: IP=%s UA=%s", r.RemoteAddr, r.UserAgent())
		http.Error(w, "Bot detected", http.StatusBadRequest)
		return
	}

	if err := Email(cfg.WaitlistReceiverEmail, "templates/email/waitlistSubmitted.html", map[string]string{
		"{{Email}}":     data.Email,
		"{{Country}}":   data.Country,
		"{{Frequency}}": data.Frequency,
		"{{Referral}}":  data.Referral,
		"{{Interest}}":  data.Interest,
	}); err != nil {
		log.Printf("Error sending waitlist email: %v", err)
		http.Error(w, "Internal error. Please try again later.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/?submitted=true", http.StatusSeeOther)
}

func isValidEmail(email string) bool {
	// Simple email validation logic
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := 0
	for i, char := range email {
		if char == '@' {
			at++
			if at > 1 || i == 0 || i == len(email)-1 {
				return false
			}
		} else if char == '.' {
			if i == 0 || i == len(email)-1 || email[i-1] == '.' {
				return false
			}
		}
	}
	return at == 1
}
