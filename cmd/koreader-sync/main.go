package main

import (
	"context"
	cryptomd5 "crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"syscall"

	"github.com/brunty/koreader-sync-server/crypto"
	database "github.com/brunty/koreader-sync-server/db"
	userpackage "github.com/brunty/koreader-sync-server/user"

	"golang.org/x/term"
)

func main() {
	ctx := context.Background()

	err := database.Init("./data/data.db.sqlite3")
	if err != nil {
		slog.Error("database init error", slog.Any("error", err))
		return
	}

	database.CreateTables()
	defer database.DBCon.Close()

	if len(os.Args) < 2 {
		fmt.Println("expected subcommand")
		os.Exit(1)
	}

	userRepo := userpackage.NewUserRepository(database.DBCon)

	switch os.Args[1] {
	case "change-password":
		changePwCmd := flag.NewFlagSet("changePw", flag.ExitOnError)
		username := changePwCmd.String("username", "", "Username to change the password for")
		changePwCmd.Parse(os.Args[2:])

		fmt.Println("Looking for user: " + *username)
		user, _ := userRepo.SelectByUsername(ctx, *username)

		if user == nil {
			fmt.Println("User not found")
			os.Exit(1)
		}

		fmt.Print("Password: ")
		bytePw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error reading password")
			os.Exit(1)
		}

		fmt.Println("")
		fmt.Println("")

		plainPw := string(bytePw)

		// We MD5 the plain password before Bcrypt hashing it because the KOReader Sync spec specifies that client
		// devices will md5 the password before providing it in the x-auth-key header
		md5Pw := md5(plainPw)

		hashedPw, err := crypto.BcryptHashPassword(md5Pw)
		if err != nil {
			fmt.Println("Error hashing password")
			os.Exit(1)
		}

		user.Password = hashedPw

		_, err = userRepo.Update(ctx, *user)

		if err != nil {
			fmt.Println("Error storing user")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Updated password for " + *username)
	default:
		fmt.Println("Unknown command")
	}
}

func md5(text string) string {
	hash := cryptomd5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
