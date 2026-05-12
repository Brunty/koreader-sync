package main

import (
	"context"
	cryptomd5 "crypto/md5"
	"encoding/hex"
	"errors"
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

		fmt.Print("Password: ")
		bytePw, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println("Error reading password")
			os.Exit(1)
		}

		fmt.Println("")
		fmt.Println("")

		err = changePassword(ctx, userRepo, *username, string(bytePw))
		if err != nil {
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

func changePassword(ctx context.Context, userRepo userpackage.UserRepository, username, plainPw string) error {
	user, _ := userRepo.SelectByUsername(ctx, username)

	if user == nil {
		return errors.New("user not found")
	}

	// We MD5 the plain password before Bcrypt hashing it because the KOReader Sync spec specifies that client
	// devices will md5 the password before providing it in the x-auth-key header.
	// Because of this we don't need to handle the error as md5 hashes are 32 chars long
	// Plus if it did error in this case, it's an internal tool, not the end of the world
	hashedPw, _ := crypto.BcryptHashPassword(md5(plainPw))

	user.Password = hashedPw

	_, err := userRepo.Update(ctx, *user)
	if err != nil {
		return fmt.Errorf("error storing user\n%w", err)
	}

	return nil
}
