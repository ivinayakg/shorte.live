package migrations

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/go-bun/bun-starter-kit/cmd/bun/migrations"
// 	"github.com/uptrace/bun/migrate"
// 	"github.com/urfave/cli/v2"
// )

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("No command provided")
// 		return
// 	}

// 	switch os.Args[1] {
// 	case "create-migrations":
// 		fmt.Println("Hello World")
// 	default:
// 		fmt.Println("Unknown command:", os.Args[1])
// 	}
// }

// func createMiigrations(c *cli.Context) error {
// 	migrator := migrate.NewMigrator(app.DB(), migrations.Migrations)

// 	name := strings.Join(c.Args().Slice(), "_")
// 	mf, err := migrator.CreateGoMigration(ctx, name)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

// 	return nil
// }
