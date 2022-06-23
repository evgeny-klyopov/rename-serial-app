package main

import (
	"github.com/evgeny-klyopov/bashColor"
	"github.com/evgeny-klyopov/rename-serial-app/internal/app"
	"github.com/evgeny-klyopov/rename-serial-app/internal/params"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	color := bashColor.NewColor()

	appHelp, commandHelp := helpTemplate(color)

	cli.AppHelpTemplate = appHelp
	cli.CommandHelpTemplate = commandHelp

	p := params.NewParams()

	app := &cli.App{
		Name: "RSD" +
			"",
		Version:  "v0.0.3",
		Flags:    p.GetFlags(),
		HelpName: "rsd",
		Usage:    "Rename serial dir",
		Action: func(c *cli.Context) error {
			return app.NewApp(c, p, color).Run()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		color.Print(color.Fatal, "\n\nErrors:")
		color.Print(color.Fatal, err.Error())
	}
}

func helpTemplate(c bashColor.Colorer) (appHelp string, commandHelp string) {
	info := c.White(`{{.Name}} - {{.Usage}}`)
	info += c.Green(`{{if .Version}} {{.Version}}{{end}}`)

	appHelp = info + `
` + c.Yellow("Usage:") + `
	{{.HelpName}} {{if .VisibleFlags}}[options]{{end}}
 {{if .Commands}}
` + c.Yellow("Commands:") + `
{{range .Commands}}{{if not .HideHelp}}` + "	" + c.GetColor(bashColor.Green) + `{{join .Names ", "}}` + c.GetColor(bashColor.Default) + `{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
` + c.Yellow("Options:") + `
{{range .VisibleFlags}}  {{.}}
{{end}}{{end}}`

	commandHelp = c.Yellow("Description:") + ` 
   {{.Usage}}
` + c.Yellow("Usage:") + `
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
{{if .VisibleFlags}}
` + c.Yellow("Arguments:") + `
	` + c.GetColor(bashColor.Green) + `stage` + c.GetColor(bashColor.Default) + `{{ "\t"}}{{ "\t"}}{{ "\t"}}{{ "\t"}} Stage or hostname
` + c.Yellow("Options:") + `
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
	return appHelp, commandHelp
}
