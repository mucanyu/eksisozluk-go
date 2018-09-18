package main

import (
	"fmt"
	"strings"
	"time"
)

// Gundem is struct that satisfies "cli.Command" interface.
type Gundem struct {
}

// Help provides detailed information about all of the commands.
func (g *Gundem) Help() string {
	return strings.TrimSpace(`
	Kullanım:
	--version
		Eksisozluk-Go'nun versiyonunu gosterir.
	--help
		Kullanilabilen argumanlari listeler.
	--gundem
		Eksisozluk gundemini listeler.
	`)
}

// Run do the tasks that the command should do.
func (g *Gundem) Run(args []string) int {
	for i := 0; i < 10; i++ {
		fmt.Print(".")
		time.Sleep(time.Second / 4)
	}
	fmt.Println()
	return 1
}

// Synopsis provides short synopsis of the command.
func (g *Gundem) Synopsis() string {
	return "Eksisozluk gundemini listeler."
}
