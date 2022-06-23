package params

import (
	"github.com/urfave/cli/v2"
)

type Params struct {
	Debug       bool
	MaskName    string
	MaskEpisode string
	MaskSeason  string
	Name        string
	Preview     bool
}

func (p *Params) GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Required:    false,
			HasBeenSet:  true,
			Value:       false,
			Aliases:     []string{"d"},
			Destination: &p.Debug,
		},
		&cli.BoolFlag{
			Name:        "preview",
			Required:    false,
			HasBeenSet:  true,
			Value:       false,
			Aliases:     []string{"p"},
			Destination: &p.Preview,
		},
		&cli.StringFlag{
			Name:        "name",
			Required:    false,
			HasBeenSet:  false,
			Aliases:     []string{"n"},
			Usage:       "Set name serial (mask name ignore)",
			Destination: &p.Name,
		},
		&cli.StringFlag{
			Name:        "mask-name",
			Required:    false,
			HasBeenSet:  true,
			Value:       `(.*)\.[s|S]`,
			Aliases:     []string{"mn"},
			Usage:       "Set regexp mask for name",
			Destination: &p.MaskName,
		},
		&cli.StringFlag{
			Name:        "mask-season",
			Required:    false,
			HasBeenSet:  true,
			Value:       `[s|S]([0-9]+)`,
			Aliases:     []string{"ms"},
			Usage:       "Set regexp mask for number season",
			Destination: &p.MaskSeason,
		},
		&cli.StringFlag{
			Name:        "mask-episode",
			Required:    false,
			HasBeenSet:  true,
			Value:       `[e|E]([0-9]+)`,
			Aliases:     []string{"me"},
			Usage:       "Set regexp mask for number episode",
			Destination: &p.MaskEpisode,
		},
	}
}

func NewParams() Params {
	return Params{}
}
