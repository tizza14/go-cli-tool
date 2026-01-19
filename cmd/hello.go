package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/go-cli-tool/internal/greeter"
)

var (
	name  string
	upper bool
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Print a greeting message",
	Long: `The hello command prints a friendly greeting message.
You can customize the greeting by providing a name.`,
	Example: `  go-cli-tool hello
  go-cli-tool hello --name John
  go-cli-tool hello --name John --upper`,
	RunE: func(cmd *cobra.Command, args []string) error {
		g := greeter.New()
		message := g.Greet(name, upper)
		fmt.Println(message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)

	helloCmd.Flags().StringVarP(&name, "name", "n", "World", "name to greet")
	helloCmd.Flags().BoolVarP(&upper, "upper", "u", false, "convert message to uppercase")
}
