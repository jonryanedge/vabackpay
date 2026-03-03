package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func formatCurrency(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)
	return formatWithCommas(s)
}

func formatWithCommas(s string) string {
	parts := strings.Split(s, ".")
	intPart := parts[0]
	result := ""
	for i, r := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result += ","
		}
		result += string(r)
	}
	if len(parts) > 1 {
		result += "." + parts[1]
	}
	return result
}

var monthNames = []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

type Config struct {
	Debug     bool
	SMTPHost  string
	SMTPPort  int
	SMTPUser  string
	SMTPPass  string
	FromEmail string
	FromName  string
	AppURL    string
}

var config Config

type YearlyData struct {
	Year        int     `json:"year"`
	MonthlyRate float64 `json:"monthly_rate"`
	YearTotal   float64 `json:"year_total"`
}

type CalcResult struct {
	TotalBackpay float64      `json:"total_backpay"`
	YearlyData   []YearlyData `json:"yearly_data"`
	StartMonth   int          `json:"start_month"`
	StartYear    int          `json:"start_year"`
	EndMonth     int          `json:"end_month"`
	EndYear      int          `json:"end_year"`
	Rating       int          `json:"rating"`
	Dependents   int          `json:"dependents"`
}

type EmailSubmission struct {
	ID             string       `json:"id"`
	Email          string       `json:"email"`
	Timestamp      time.Time    `json:"timestamp"`
	TotalBackpay   float64      `json:"total_backpay"`
	StartMonth     int          `json:"start_month"`
	StartYear      int          `json:"start_year"`
	EndMonth       int          `json:"end_month"`
	EndYear        int          `json:"end_year"`
	Rating         int          `json:"rating"`
	Dependents     int          `json:"dependents"`
	YearlyData     []YearlyData `json:"yearly_data"`
	EmailSent      bool         `json:"email_sent"`
	EmailSentError string       `json:"email_sent_error,omitempty"`
}

var rates = map[int]map[int]float64{
	2026: {10: 180.42, 20: 356.66, 30: 552.47, 40: 795.84, 50: 1132.90, 60: 1435.02, 70: 1808.45, 80: 2102.15, 90: 2362.30, 100: 3938.58},
	2025: {10: 175.51, 20: 346.95, 30: 537.42, 40: 774.16, 50: 1102.04, 60: 1395.93, 70: 1759.19, 80: 2044.89, 90: 2297.96, 100: 3831.30},
	2024: {10: 171.25, 20: 337.88, 30: 523.43, 40: 754.08, 50: 1074.19, 60: 1361.50, 70: 1715.71, 80: 1994.83, 90: 2241.39, 100: 3737.85},
	2023: {10: 165.92, 20: 327.42, 30: 507.40, 40: 731.00, 50: 1040.84, 60: 1319.91, 70: 1662.75, 80: 1933.06, 90: 2172.39, 100: 3621.95},
	2022: {10: 152.64, 20: 301.74, 30: 467.35, 40: 673.28, 50: 958.38, 60: 1215.28, 70: 1531.30, 80: 1780.67, 90: 2000.97, 100: 3332.06},
	2021: {10: 144.14, 20: 284.93, 30: 441.35, 40: 635.77, 50: 905.04, 60: 1146.39, 70: 1444.71, 80: 1679.35, 90: 1887.18, 100: 3146.42},
	2020: {10: 142.29, 20: 281.27, 30: 435.69, 40: 627.61, 50: 893.43, 60: 1131.68, 70: 1426.17, 80: 1657.80, 90: 1862.96, 100: 3106.04},
	2019: {10: 140.05, 20: 276.84, 30: 428.83, 40: 617.73, 50: 879.36, 60: 1113.86, 70: 1403.71, 80: 1631.69, 90: 1833.62, 100: 3057.13},
	2018: {10: 136.24, 20: 269.30, 30: 417.15, 40: 600.90, 50: 855.41, 60: 1083.52, 70: 1365.48, 80: 1587.25, 90: 1783.68, 100: 2973.86},
	2017: {10: 133.57, 20: 264.22, 30: 408.97, 40: 589.12, 50: 838.64, 60: 1062.97, 70: 1338.70, 80: 1556.04, 90: 1748.71, 100: 2915.55},
	2016: {10: 133.17, 20: 263.23, 30: 407.75, 40: 587.36, 50: 836.13, 60: 1059.09, 70: 1334.71, 80: 1551.48, 90: 1743.48, 100: 2906.83},
	2015: {10: 133.17, 20: 263.23, 30: 407.75, 40: 587.36, 50: 836.13, 60: 1059.09, 70: 1334.71, 80: 1551.48, 90: 1743.48, 100: 2858.75},
	2014: {10: 130.94, 20: 258.83, 30: 400.93, 40: 577.54, 50: 822.15, 60: 1041.39, 70: 1312.40, 80: 1525.55, 90: 1714.34, 100: 2816.00},
	2013: {10: 128.60, 20: 254.24, 30: 393.87, 40: 567.53, 50: 808.04, 60: 1023.40, 70: 1289.60, 80: 1498.89, 90: 1684.40, 100: 2769.00},
	2012: {10: 123.84, 20: 244.95, 30: 379.28, 40: 546.33, 50: 777.63, 60: 985.22, 70: 1241.87, 80: 1442.87, 90: 1621.21, 100: 2673.75},
	2011: {10: 123.84, 20: 244.95, 30: 379.28, 40: 546.33, 50: 777.63, 60: 985.22, 70: 1241.87, 80: 1442.87, 90: 1621.21, 100: 2673.75},
	2010: {10: 117.00, 20: 231.33, 30: 358.13, 40: 515.79, 50: 734.18, 60: 930.08, 70: 1172.32, 80: 1362.45, 90: 1531.20, 100: 2527.25},
	2009: {10: 110.61, 20: 218.85, 30: 338.76, 40: 487.95, 50: 694.41, 60: 879.73, 70: 1108.77, 80: 1288.22, 90: 1447.72, 100: 2389.00},
	2008: {10: 117.00, 20: 231.33, 30: 358.13, 40: 515.79, 50: 734.18, 60: 930.08, 70: 1172.32, 80: 1362.45, 90: 1531.20, 100: 2527.00},
	2007: {10: 115.00, 20: 227.40, 30: 352.00, 40: 507.00, 50: 721.00, 60: 913.00, 70: 1150.00, 80: 1336.00, 90: 1501.00, 100: 2471.00},
	2006: {10: 110.00, 20: 217.48, 30: 336.61, 40: 484.94, 50: 690.17, 60: 874.36, 70: 1101.87, 80: 1280.22, 90: 1438.46, 100: 2389.00},
	2005: {10: 107.00, 20: 211.54, 30: 327.38, 40: 471.60, 50: 671.00, 60: 849.86, 70: 1070.86, 80: 1243.86, 90: 1397.54, 100: 2316.00},
	2004: {10: 108.00, 20: 210.00, 30: 324.00, 40: 466.00, 50: 663.00, 60: 839.00, 70: 1056.00, 80: 1227.00, 90: 1380.00, 100: 2299.00},
	2003: {10: 106.00, 20: 205.00, 30: 316.00, 40: 454.00, 50: 646.00, 60: 817.00, 70: 1029.00, 80: 1195.00, 90: 1344.00, 100: 2239.00},
	2002: {10: 104.00, 20: 201.00, 30: 310.00, 40: 445.00, 50: 633.00, 60: 801.00, 70: 1008.00, 80: 1171.00, 90: 1317.00, 100: 2193.00},
	2001: {10: 103.00, 20: 199.00, 30: 306.00, 40: 439.00, 50: 625.00, 60: 790.00, 70: 995.00, 80: 1155.00, 90: 1299.00, 100: 2163.00},
	2000: {10: 101.00, 20: 194.00, 30: 298.00, 40: 427.00, 50: 609.00, 60: 769.00, 70: 969.00, 80: 1125.00, 90: 1266.00, 100: 2107.00},
}

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: .env file not found, using defaults")
	}

	config.Debug = os.Getenv("DEBUG") == "true"
	config.SMTPHost = os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}
	config.SMTPPort = port
	config.SMTPUser = os.Getenv("SMTP_USERNAME")
	config.SMTPPass = os.Getenv("SMTP_PASSWORD")
	config.FromEmail = os.Getenv("FROM_EMAIL")
	config.FromName = os.Getenv("FROM_NAME")
	config.AppURL = os.Getenv("APP_URL")

	if config.AppURL == "" {
		config.AppURL = "http://localhost:8080"
	}
}

func getDependentAddOn(rating int, dependents int) float64 {
	if dependents == 0 || rating < 30 {
		return 0
	}
	if rating >= 70 {
		return float64(dependents) * 153.00
	}
	return float64(dependents) * 109.00
}

func calculateBackpay(startMonth, startYear, endMonth, endYear, rating int, dependents int) CalcResult {
	var totalBackpay float64
	var yearlyData []YearlyData

	for year := startYear; year <= endYear; year++ {
		yearRates, ok := rates[year]
		if !ok {
			continue
		}

		baseRate := yearRates[rating]
		dependentAddOn := getDependentAddOn(rating, dependents)
		monthlyRate := baseRate + dependentAddOn

		monthsInYear := 12
		if year == endYear {
			monthsInYear = endMonth
		}
		if year == startYear {
			monthsInYear = monthsInYear - startMonth + 1
		}

		yearTotal := monthlyRate * float64(monthsInYear)
		totalBackpay += yearTotal

		yearlyData = append(yearlyData, YearlyData{
			Year:        year,
			MonthlyRate: monthlyRate,
			YearTotal:   yearTotal,
		})
	}

	return CalcResult{
		TotalBackpay: totalBackpay,
		YearlyData:   yearlyData,
		StartMonth:   startMonth,
		StartYear:    startYear,
		EndMonth:     endMonth,
		EndYear:      endYear,
		Rating:       rating,
		Dependents:   dependents,
	}
}

func saveEmailSubmission(submission EmailSubmission) error {
	filename := "/Users/jre/code/go/backpaycalc/email_submissions.json"

	var submissions []EmailSubmission

	// Read existing file
	data, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(data, &submissions)
	}

	// Add new submission
	submissions = append(submissions, submission)

	// Write back to file
	newData, err := json.MarshalIndent(submissions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, newData, 0644)
}

func sendEmail(to, subject, body string) error {
	if config.Debug {
		fmt.Printf("[DEBUG] Email would be sent:\n")
		fmt.Printf("  To: %s\n", to)
		fmt.Printf("  Subject: %s\n", subject)
		fmt.Printf("  Body: %s\n", body)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}

func generateEmailBody(result CalcResult) string {
	html := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h1 style="color: #333;">VA Disability Backpay Calculation</h1>
			
			<div style="background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0;">
				<h2 style="margin: 0 0 10px 0;">Total Backpay: $%s</h2>
			</div>
			
			<h3>Calculation Details</h3>
			<ul>
				<li><strong>Effective Date:</strong> %s %d</li>
				<li><strong>End Date:</strong> %s %d</li>
				<li><strong>Disability Rating:</strong> %d%%</li>
				<li><strong>Dependents:</strong> %d</li>
			</ul>
			
			<h3>Yearly Breakdown</h3>
			<table style="width: 100%%; border-collapse: collapse;">
				<tr style="background: #f5f5f5;">
					<th style="padding: 10px; text-align: left;">Year</th>
					<th style="padding: 10px; text-align: right;">Monthly</th>
					<th style="padding: 10px; text-align: right;">Yearly Total</th>
				</tr>
	`, formatCurrency(result.TotalBackpay),
		monthNames[result.StartMonth], result.StartYear,
		monthNames[result.EndMonth], result.EndYear,
		result.Rating, result.Dependents)

	for _, yd := range result.YearlyData {
		html += fmt.Sprintf(`
				<tr>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">%d</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee; text-align: right;">$%s</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee; text-align: right;">$%s</td>
				</tr>
		`, yd.Year, formatCurrency(yd.MonthlyRate), formatCurrency(yd.YearTotal))
	}

	html += `
			</table>
			
			<p style="margin-top: 30px; color: #666; font-size: 14px;">
				This email was sent from the VA Disability Backpay Calculator.<br>
				Calculate your backpay anytime at: ` + config.AppURL + `
			</p>
		</body>
		</html>
	`

	return html
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.Execute(w, nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	startMonth, _ := strconv.Atoi(r.Form.Get("start_month"))
	startYear, _ := strconv.Atoi(r.Form.Get("start_year"))
	endMonth, _ := strconv.Atoi(r.Form.Get("end_month"))
	endYear, _ := strconv.Atoi(r.Form.Get("end_year"))
	rating, _ := strconv.Atoi(r.Form.Get("rating"))
	hasDependents := r.Form.Get("has_dependents") == "on"
	dependents := 0
	if hasDependents {
		dependents, _ = strconv.Atoi(r.Form.Get("dependents"))
	}

	result := calculateBackpay(startMonth, startYear, endMonth, endYear, rating, dependents)

	tmpl, err := template.New("results.html").Funcs(template.FuncMap{
		"formatCurrency": formatCurrency,
	}).ParseFiles("templates/results.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, result)
}

func emailHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.Form.Get("email")
	startMonth, _ := strconv.Atoi(r.Form.Get("start_month"))
	startYear, _ := strconv.Atoi(r.Form.Get("start_year"))
	endMonth, _ := strconv.Atoi(r.Form.Get("end_month"))
	endYear, _ := strconv.Atoi(r.Form.Get("end_year"))
	rating, _ := strconv.Atoi(r.Form.Get("rating"))
	dependents, _ := strconv.Atoi(r.Form.Get("dependents"))

	if email == "" {
		w.Write([]byte("Error: email is required"))
		return
	}

	// Calculate the results
	result := calculateBackpay(startMonth, startYear, endMonth, endYear, rating, dependents)

	// Generate email body
	subject := fmt.Sprintf("Your VA Disability Backpay Calculation: $%s", formatCurrency(result.TotalBackpay))
	body := generateEmailBody(result)

	// Try to send email
	var emailError string
	err := sendEmail(email, subject, body)
	if err != nil {
		emailError = err.Error()
		fmt.Printf("Error sending email: %v\n", err)
	}

	// Save submission to JSON
	submission := EmailSubmission{
		ID:             fmt.Sprintf("%d", time.Now().UnixNano()),
		Email:          email,
		Timestamp:      time.Now(),
		TotalBackpay:   result.TotalBackpay,
		StartMonth:     startMonth,
		StartYear:      startYear,
		EndMonth:       endMonth,
		EndYear:        endYear,
		Rating:         rating,
		Dependents:     dependents,
		YearlyData:     result.YearlyData,
		EmailSent:      err == nil,
		EmailSentError: emailError,
	}

	saveEmailSubmission(submission)

	if config.Debug {
		w.Header().Set("Content-Type", "text/html")
		if err != nil {
			w.Write([]byte("OK (debug mode - email not sent)"))
		} else {
			w.Write([]byte("OK (debug mode - email not sent)"))
		}
	} else {
		w.Header().Set("Content-Type", "text/html")
		if err != nil {
			w.Write([]byte("Error sending email"))
		} else {
			w.Write([]byte("OK"))
		}
	}
}

func main() {
	loadConfig()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calculate", calculateHandler)
	http.HandleFunc("/email", emailHandler)

	fmt.Println("Server starting on http://localhost:8080")
	if config.Debug {
		fmt.Println("DEBUG MODE: Emails will not be sent")
	}
	http.ListenAndServe(":8080", nil)
}
