package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func version(*cobra.Command, []string) {
	fmt.Println(appVersion)
}

func test(*cobra.Command, []string) {

}

func update(*cobra.Command, []string) {

}
