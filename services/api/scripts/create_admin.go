package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"medic-api/config"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username admin [admin]: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		username = "admin"
	}

	fmt.Print("Nama lengkap admin: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName)

	if fullName == "" {
		fmt.Println("Nama lengkap wajib diisi.")
		os.Exit(1)
	}

	fmt.Print("Password admin: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Gagal membaca password:", err)
		os.Exit(1)
	}

	password := string(passwordBytes)
	if len(password) < 8{
		fmt.Println("Password minimal 8 karakter.")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		fmt.Println("Gagal membuat hash password:", err)
		os.Exit(1)
	}

	db, err := config.ConnectDB()
	if err != nil {
		fmt.Println("Gagal koneksi database:", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx := context.Background()

	result, err := db.ExecContext(ctx, `
		IF NOT EXISTS (
			SELECT 1 FROM dbo.users WHERE username = @p1
		)
		BEGIN
			INSERT INTO dbo.users
				(username, password_hash, full_name, role, is_active)
			VALUES
				(@p1, @p2, @p3, 'admin', 1);
		END;
	`, username, string(hash), fullName)

	if err != nil {
		fmt.Println("Gagal membuat akun:", err)
		os.Exit(1)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Akun diproses, tetapi hasil tidak dapat diperiksa:", err)
		os.Exit(1)
	}

	if rows == 0 {
		fmt.Println("Username sudah ada; tidak ada akun baru dibuat.")
		return
	}

	fmt.Println("Akun admin berhasil dibuat.")
	fmt.Println("Username:", username)
	fmt.Println("Role: admin")
}
