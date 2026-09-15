package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "\ufeff") // Xóa BOM của Windows Notepad nếu có
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func main() {
	timeoutSec := 10
	numWorkers := 30 // Đề xuất 20-30 luồng khi dùng server Anisette công cộng để tránh quá tải

	idFile := filepath.Join("folder data", "dataid.txt")
	proxyFile := filepath.Join("folder data", "dataproxy.txt")
	outputDir := "folder save"

	accounts, err := readLines(idFile)
	if err != nil || len(accounts) == 0 {
		log.Fatalf("Không thể đọc danh sách tài khoản tại: %s", idFile)
	}
	proxies, _ := readLines(proxyFile)

	fmt.Printf("%s=== HỆ THỐNG GOLANG GSA STANDALONE (MODULAR ARCHITECTURE) ===%s\n", ColorBold, ColorReset)
	fmt.Printf("-> Đang kiểm tra máy chủ Anisette...\n")
	aniCount := RefreshAnisette()
	fmt.Printf("-> Anisette sẵn sàng: %d máy chủ hoạt động\n", aniCount)

	if aniCount == 0 {
		fmt.Printf("%s[CẢNH BÁO] Không có máy chủ Anisette nào phản hồi. Hãy kiểm tra kết nối mạng hoặc bật Docker local (port 6969/6970).%s\n", ColorRed, ColorReset)
	} else {
		// Bật luồng tự làm mới OTP để tránh hết hạn mã trong suốt quá trình chạy
		StartAnisetteAutoRefresher()
	}

	fmt.Printf("-> Đã nạp: %d Tài khoản (%s) | %d Proxy (%s)\n", len(accounts), idFile, len(proxies), proxyFile)
	fmt.Printf("-> Thư mục lưu kết quả phân loại: [%s/]\n", outputDir)
	fmt.Println(strings.Repeat("=", 75))

	writer, err := NewResultWriter(filepath.Join(outputDir, "dummy.txt"))
	if err != nil {
		log.Fatalf("Lỗi tạo thư mục lưu kết quả: %v", err)
	}
	defer writer.Close()

	proxyMgr := NewProxyManager(proxies)
	engine := NewEngine(numWorkers, timeoutSec, proxyMgr, writer)

	startTime := time.Now()
	engine.Run(accounts)

	valid, wrong, errCnt := engine.Stats()
	fmt.Printf("\n%s=== HOÀN TẤT KIỂM TRA TRONG %s ===%s\n", ColorBold, time.Since(startTime).Round(time.Second), ColorReset)
	fmt.Printf("Hợp lệ (3Q / 2FA / Verify): %d  |  Lỗi / Khóa / Sai pass: %d  |  Lỗi mạng: %d\n", valid, wrong, errCnt)
	fmt.Printf("Kết quả đã được phân loại tự động tại thư mục: %s/\n", outputDir)
}
