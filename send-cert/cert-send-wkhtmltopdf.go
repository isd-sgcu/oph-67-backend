package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// config
const (
	// SendGrid email sender
	senderEmail = "isd.team.sgcu-no-reply@cuopenhouse2025.com"
	// CSV file path
	csvFilePath = "students1.csv"
	// HTML template file path
	htmlTemplatePath = "cert_template.html"
	// Font file path
	fontFilePath = "THSarabunNew.ttf"
	// Image file path
	imgFilePath = "template.png"
	// SendGrid API key
	sendGridAPIKEY = "<SendGrid API KEY>"
	// CertPrivateKey
	certPrivateKey = "<CertPrivateKey>" // base64 encoded private key
	// config
	isWkhtmltopdf = false
)

type CertData struct {
	UserName string
	Token    string
}

type User struct {
	Name  string
	Email string
}

func main() {
	// Open log file
	f := filepath.Join(Getwd(), "log-temp.txt")
	logFile, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Multi-writer to log both to file and console (optional)
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	if _, err := os.Stat("progress"); os.IsNotExist(err) {
		if err := os.Mkdir("progress", fs.ModePerm); err != nil {
			log.Fatalf("Failed to create progress folder: %v", err)
		}
	}

	f = filepath.Join(Getwd(), csvFilePath)
	users, err := loadCSV(f)
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	for i, user := range users {
		progressPath := getProgressFilePath(user.Name)
		email_progressPath := getProgressFilePath(user.Email)
		// Skip if already processed
		if _, err := os.Stat(progressPath); err == nil {
			log.Printf("⏩ Already name processed: %s path: %s", user.Name, progressPath)
			continue
		}
		if _, err := os.Stat(email_progressPath); err == nil {
			log.Printf("⏩ Already email processed: %s path: %s", user.Email, email_progressPath)
			continue
		}

		log.Printf("Progression: %d/%d - Processing %s", i+1, len(users), user.Name)
		if err := processUser(user); err != nil {
			log.Printf("❌ Error processing %s: %v", user.Email, err)
			continue
		}

		// Write checkpoint
		if err := os.WriteFile(progressPath, []byte("done"), 0644); err != nil {
			log.Printf("⚠️ Failed to write progress file for %s: %v", user.Name, err)
		}
		if err := os.WriteFile(email_progressPath, []byte("done"), 0644); err != nil {
			log.Printf("⚠️ Failed to write progress file for %s: %v", user.Email, err)
		}
	}
}

func loadCSV(filePath string) ([]User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var users []User
	for i, row := range records {
		if i == 0 { // skip header
			continue
		}
		if len(row) < 9 {
			// if len(row) < 6 {
			continue
		}
		firstName := strings.Trim(row[5], " ")
		lastName := strings.Trim(row[6], " ")
		// firstName := strings.Trim(row[2], " ")
		// lastName := strings.Trim(row[3], " ")

		users = append(users, User{
			Name:  firstName + " " + lastName,
			Email: strings.Trim(row[8], " "),
			// Email: strings.Trim(row[5], " "),
		})
	}
	return users, nil
}

func processUser(user User) error {
	privateKey, err := base64.RawURLEncoding.DecodeString(certPrivateKey)
	if err != nil {
		return fmt.Errorf("decode private key: %w", err)
	}

	signature := ed25519.Sign(privateKey, []byte(user.Name))
	token := base64.RawURLEncoding.EncodeToString(signature)

	// embed content for future verification
	userName := base64.RawURLEncoding.EncodeToString([]byte(user.Name))

	token = userName + "." + token

	var fontPath, imgPath, outputPDF string

	outputPDF = "cert_" + user.Email + ".pdf"
	outputPDF = filepath.Join(Getwd(), outputPDF)

	if isWkhtmltopdf {
		data := CertData{
			UserName: user.Name,
			Token:    token,
		}

		// Generate HTML
		f := filepath.Join(Getwd(), htmlTemplatePath)
		tmpl, err := template.ParseFiles(f)
		if err != nil {
			return fmt.Errorf("parse HTML template: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("execute template: %w", err)
		}

		// Write HTML file
		htmlFile := "temp_" + user.Email + ".html"
		htmlFile = filepath.Join(Getwd(), htmlFile)
		if err := os.WriteFile(htmlFile, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write HTML: %w", err)
		}

		f = filepath.Join(Getwd(), fontFilePath)
		fontPath, _ = filepath.Abs(f)

		f = filepath.Join(Getwd(), imgFilePath)
		imgPath, _ = filepath.Abs(f)

		cmd := exec.Command("wkhtmltopdf.exe",
			"--no-pdf-compression",
			"--margin-bottom", "0",
			"--margin-left", "0",
			"--margin-right", "0",
			"--margin-top", "0",
			"--orientation", "Landscape",
			"--page-size", "A4",
			"--allow", fontPath,
			"--allow", imgPath,
			// silence the console output
			"--quiet",
			htmlFile,
			outputPDF,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("generate PDF: %w", err)
		}

		// Cleanup
		defer func() {
			_ = os.Remove(htmlFile)
			_ = os.Remove(outputPDF)
		}()
	} else {
		err = fetchCertificatePDF(user.Name, token, outputPDF)
		if err != nil {
			return fmt.Errorf("fetch PDF: %w", err)
		}

		// Cleanup
		defer func() {
			_ = os.Remove(outputPDF)
		}()
	}
	log.Println("✅ Certificate created:", outputPDF)

	// Send email
	if err := sendEmail(user.Email, user.Name, outputPDF); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func sendEmail(toEmail, userName, pdfPath string) error {
	from := mail.NewEmail("CUOPENHOUSE2025 Team", senderEmail)
	subject := "CUOPENHOUSE2025 - 🎉 Congratulations! Here's Your Certificate"

	to := mail.NewEmail(userName, toEmail)

	var plainTextContent, htmlContent string

	// {
	// 	plainTextContent = fmt.Sprintf(`เรียนคุณ %s ขอบคุณน้อง ๆ ทุกคนที่ให้ความสนใจและเข้าร่วมกิจกรรมของเราในงาน CU Open House 2025 ทางทีมงานขอส่งมอบ เกียรติบัตรการเข้าร่วมกิจกรรม มาให้ในรูปแบบไฟล์แนบในอีเมลฉบับนี้
	// 		📎 กรุณาตรวจสอบชื่อ-นามสกุลในเกียรติบัตรให้เรียบร้อย

	// 		หมายเหตุสำคัญ:
	// 		ชื่อ-นามสกุลในเกียรติบัตรไม่สามารถแก้ไขได้ เนื่องจากข้อมูลดังกล่าวมาจากแบบฟอร์ม Google Form ที่น้อง ๆ ได้กรอกไว้ด้วยตนเอง
	// 	เพื่อความปลอดภัย เกียรติบัตรแต่ละฉบับมี รหัสยืนยันเฉพาะตัว (Unique ID) อยู่ที่มุมขวาล่างของเกียรติบัตร เพื่อป้องกันการปลอมแปลง โปรดอย่าแก้ไขหรือส่งต่อไฟล์โดยมิชอบ

	// 	หากมีคำถามหรือพบปัญหาในการเปิดไฟล์แนบ สามารถติดต่อทีมงานได้ทางอีเมลนี้ isd.team.sgcu@gmail.com
	// 	ขอขอบคุณอีกครั้งที่เข้าร่วมกิจกรรมกับเรา หวังว่าจะได้พบกันอีกในโอกาสหน้า 😊

	// 	ขอแสดงความนับถือ
	// 	ฝ่ายพัฒนาระบบ CU Open House 2025`, userName)

	// 	htmlContent = fmt.Sprintf(`<p>เรียนคุณ <strong>%s</strong></p>

	// 	<p>ขอบคุณน้อง ๆ ทุกคนที่ให้ความสนใจและเข้าร่วมกิจกรรมของเราในงาน <strong>CU Open House 2025</strong> ทางทีมงานขอส่งมอบ <strong>เกียรติบัตรการเข้าร่วมกิจกรรม</strong> มาให้ในรูปแบบไฟล์แนบในอีเมลฉบับนี้</p>

	// 	<p>📎 <strong>กรุณาตรวจสอบชื่อ-นามสกุลในเกียรติบัตรให้เรียบร้อย</strong></p>

	// 	<p><strong>หมายเหตุสำคัญ:</strong><br>
	// 	ชื่อ-นามสกุลในเกียรติบัตรไม่สามารถแก้ไขได้ เนื่องจากข้อมูลดังกล่าวมาจากแบบฟอร์ม Google Form ที่น้อง ๆ ได้กรอกไว้ด้วยตนเอง<br>
	// 	เพื่อความปลอดภัย เกียรติบัตรแต่ละฉบับมี <strong>รหัสยืนยันเฉพาะตัว (Unique ID)</strong> อยู่ที่ <strong>มุมขวาล่างของเกียรติบัตร</strong> เพื่อป้องกันการปลอมแปลง โปรดอย่าแก้ไขหรือส่งต่อไฟล์โดยมิชอบ</p>

	// 	<p>หากมีคำถามหรือพบปัญหาในการเปิดไฟล์แนบ สามารถติดต่อทีมงานได้ทางอีเมล <strong>isd.team.sgcu@gmail.com</strong></p>

	// 	<p>ขอขอบคุณอีกครั้งที่เข้าร่วมกิจกรรมกับเรา<br>
	// 	หวังว่าจะได้พบกันอีกในโอกาสหน้า 😊</p>

	// 	<p>ขอแสดงความนับถือ<br>
	// 	ฝ่ายพัฒนาระบบ CU Open House 2025</p>`, userName)
	// }

	{
		plainTextContent = fmt.Sprintf(`เรียนคุณ %s

	ขอบคุณน้อง ๆ ทุกคนที่ให้ความสนใจและเข้าร่วมกิจกรรมของเราในงาน CU Open House 2025 ทางทีมงานขอส่งมอบ เกียรติบัตรการเข้าร่วมกิจกรรม มาให้ในรูปแบบไฟล์แนบในอีเมลฉบับนี้
	📎 กรุณาตรวจสอบชื่อ-นามสกุลในเกียรติบัตรให้เรียบร้อย
	
	หมายเหตุสำคัญ:
	
	1. ชื่อ-นามสกุลในเกียรติบัตรไม่สามารถแก้ไขได้ เนื่องจากข้อมูลดังกล่าวมาจากแบบฟอร์ม Google Form ที่น้อง ๆ ได้กรอกไว้ด้วยตนเอง
	
	2. เพื่อความปลอดภัย เกียรติบัตรแต่ละฉบับมี รหัสยืนยันเฉพาะตัว (Unique ID) อยู่ที่มุมขวาล่างของเกียรติบัตร เพื่อป้องกันการปลอมแปลง โปรดอย่าแก้ไขหรือส่งต่อไฟล์โดยมิชอบ
		
	หากมีคำถามหรือพบปัญหาในการเปิดไฟล์แนบ สามารถติดต่อทีมงานได้ทางอีเมล isd.team.sgcu@gmail.com
	
	ขอขอบคุณอีกครั้งที่เข้าร่วมกิจกรรมกับเรา
	หวังว่าจะได้พบกันอีกในโอกาสหน้า 
	
	ขอแสดงความนับถือ
	ฝ่ายพัฒนาระบบ CU Open House 2025`, userName)

		htmlContent = fmt.Sprintf(`<p>เรียนคุณ <strong>%s</strong></p>
	
	<p>ขอบคุณน้อง ๆ ทุกคนที่ให้ความสนใจและเข้าร่วมกิจกรรมของเราในงาน <strong>CU Open House 2025</strong> ทางทีมงานขอส่งมอบ <strong>เกียรติบัตรการเข้าร่วมกิจกรรม</strong> มาให้ในรูปแบบไฟล์แนบในอีเมลฉบับนี้<br>
	📎 <strong>กรุณาตรวจสอบชื่อ-นามสกุลในเกียรติบัตรให้เรียบร้อย</strong></p>
	
	<p><strong>หมายเหตุสำคัญ:</strong></p>
	<ol>
	  <li>ชื่อ-นามสกุลในเกียรติบัตรไม่สามารถแก้ไขได้ เนื่องจากข้อมูลดังกล่าวมาจากแบบฟอร์ม Google Form ที่น้อง ๆ ได้กรอกไว้ด้วยตนเอง</li>
	  <li>เพื่อความปลอดภัย เกียรติบัตรแต่ละฉบับมี <strong>รหัสยืนยันเฉพาะตัว (Unique ID)</strong> อยู่ที่มุมขวาล่างของเกียรติบัตร เพื่อป้องกันการปลอมแปลง โปรดอย่าแก้ไขหรือส่งต่อไฟล์โดยมิชอบ</li>
	</ol>
	
	<p>หากมีคำถามหรือพบปัญหาในการเปิดไฟล์แนบ สามารถติดต่อทีมงานได้ทางอีเมล <strong>isd.team.sgcu@gmail.com</strong></p>
	
	<p>ขอขอบคุณอีกครั้งที่เข้าร่วมกิจกรรมกับเรา<br>
	หวังว่าจะได้พบกันอีกในโอกาสหน้า</p>
	
	<p>ขอแสดงความนับถือ<br>
	ฝ่ายพัฒนาระบบ CU Open House 2025</p>`, userName)
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Attach PDF
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return fmt.Errorf("read pdf: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(pdfBytes)
	attachment := mail.NewAttachment()
	attachment.SetContent(encoded)
	attachment.SetType("application/pdf")
	attachment.SetFilename("certificate.pdf")
	attachment.SetDisposition("attachment")
	message.AddAttachment(attachment)

	client := sendgrid.NewSendClient(sendGridAPIKEY)
	response, err := client.SendWithContext(context.Background(), message)
	if err != nil {
		return fmt.Errorf("sendgrid: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid response: %s", response.Body)
	}

	log.Printf("📬 Email sent to %s", toEmail)
	return nil
}

func getProgressFilePath(userName string) string {
	progressDir := filepath.Join(Getwd(), "progress")
	if err := os.MkdirAll(progressDir, fs.ModePerm); err != nil {
		panic(err)
	}

	fileName := userName
	fileName = strings.Replace(fileName, " ", "_", -1)
	return filepath.Join(Getwd(), "progress", fileName)
}

func Getwd() string {
	// Get the current working directory
	dir, err := os.Getwd()
	dir = filepath.Join(dir, "send-cert")
	if err != nil {
		panic(err)
	}
	return dir
}

func fetchCertificatePDF(userName, token, outputPDF string) error {
	url := "http://localhost:3000/api/generate-pdf" // works on local but not test on prod yet. so I'll hardcode it for now

	// format url to URL-encoded
	userName = ToURLSafe(userName)
	token = ToURLSafe(token)

	resp, err := http.Get(fmt.Sprintf("%s?userName=%s&token=%s", url, userName, token))
	if err != nil {
		return fmt.Errorf("request to PDF service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %s", resBody)
	}

	pdfBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read PDF response: %w", err)
	}

	if err := os.WriteFile(outputPDF, pdfBytes, 0644); err != nil {
		return fmt.Errorf("save PDF file: %w", err)
	}

	log.Println("✅ Certificate fetched and saved:", outputPDF)
	return nil
}

func ToURLSafe(s string) string {
	return url.QueryEscape(s)
}
