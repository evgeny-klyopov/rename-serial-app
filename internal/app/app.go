package app

import (
	"errors"
	"fmt"
	"github.com/evgeny-klyopov/bashColor"
	"github.com/evgeny-klyopov/rename-serial-app/internal/models"
	"github.com/evgeny-klyopov/rename-serial-app/internal/params"
	"github.com/urfave/cli/v2"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type maskProperties struct {
	seasonMax   uint8
	episodesMax uint8
	mask        string
	prefix      string
}

type app struct {
	color          bashColor.Colorer
	params         params.Params
	context        *cli.Context
	inputDir       string
	files          []fs.FileInfo
	countFiles     int
	serials        []models.Serial
	maskProperties maskProperties
	mask           string
}

type Apper interface {
	Run() error
}

func NewApp(c *cli.Context, p params.Params, color bashColor.Colorer) Apper {
	return &app{
		params:  p,
		color:   color,
		context: c,
	}
}

func (a *app) Run() error {
	if err := a.setInputDir(); err != nil {
		return err
	}

	a.print("Input directory", a.inputDir)

	if err := a.setFiles(); err != nil {
		return err
	}

	a.print("Found files", fmt.Sprintf("%d", a.countFiles))
	if a.params.Debug == true && a.countFiles > 0 {
		a.print("Files", "")
		for _, f := range a.files {
			a.debug(f.Name())
		}
		fmt.Println()
	}

	if err := a.parseFiles(); err != nil {
		return err
	}

	a.setMask()
	a.print("Mask", a.mask)

	return a.renameSerials()
}

func (a *app) debug(message string) {
	fmt.Println(a.color.GetColor(bashColor.Magenta) + message)
}

func (a *app) print(header string, message string) {
	fmt.Print(a.color.GetColor(bashColor.Yellow) + header + ": ")
	fmt.Println(a.color.GetColor(bashColor.Teal) + message)
}

func (a *app) setInputDir() error {
	var err error
	var base string

	if a.context.Args().First() != "" {
		a.inputDir, err = filepath.Abs(a.context.Args().First())
	}

	if err == nil && a.inputDir == "" {
		base, err = os.Getwd()
		if err != nil {
			return err
		}
		a.inputDir = base
	}

	return err
}

func (a *app) setFiles() error {
	var err error
	a.files, err = ioutil.ReadDir(a.inputDir)
	a.countFiles = len(a.files)

	return err
}

func (a *app) parseByMask(name []byte, mask string) string {
	var result string
	re := regexp.MustCompile(mask)
	search := re.FindAllSubmatch(name, -1)

	if len(search) > 0 && len(search[0]) > 1 {
		result = string(search[0][1])
	}

	return result
}

func (a *app) parseNumberByMask(name []byte, mask string) uint8 {
	e, err := strconv.ParseUint(a.parseByMask(name, mask), 10, 8)
	if err != nil && a.params.Debug {
		a.print("Warning mask", mask)
	}

	return uint8(e)
}

func (a *app) parseFile(f fs.FileInfo) *models.Serial {
	var serial *models.Serial

	if a.params.Debug {
		a.print("File name", f.Name())
	}

	name := []byte(f.Name())

	prefix := a.params.Name
	if prefix == "" {
		prefix = a.parseByMask(name, a.params.MaskName)
	}

	season := a.parseNumberByMask(name, a.params.MaskSeason)
	episode := a.parseNumberByMask(name, a.params.MaskEpisode)

	var isFound bool
	if season > 0 && episode > 0 && prefix != "" {
		isFound = true
		s := models.Serial{
			Name: f.Name(),
			Ext:  filepath.Ext(f.Name()),
			Info: models.SerialInfo{
				Prefix:  prefix,
				Season:  season,
				Episode: episode,
			},
		}
		serial = &s
	}

	if a.params.Debug {
		if !isFound {
			a.print("Warning", "Not found serial by mask")
		}
		a.print("Prefix", prefix)
		a.print("Season", fmt.Sprintf("%d", season))
		a.print("Episode", fmt.Sprintf("%d", episode))
		fmt.Println()
	}

	return serial
}

func (a *app) parseFiles() error {
	var err error
	if a.params.Debug == true {
		a.print("Mask mame", a.params.MaskName)
		a.print("Mask season", a.params.MaskSeason)
		a.print("Mask episode", a.params.MaskEpisode)
		fmt.Println()
	}

	a.maskProperties = maskProperties{}
	a.serials = make([]models.Serial, 0, a.countFiles)
	exceptFiles := make([]string, 0, a.countFiles)
	for _, f := range a.files {
		s := a.parseFile(f)
		if s != nil {
			a.serials = append(a.serials, *s)
			a.maskProperties.prefix = s.Info.Prefix
			a.maskProperties.seasonMax = s.Info.Season
			if s.Info.Episode > a.maskProperties.episodesMax {
				a.maskProperties.episodesMax = s.Info.Episode
			}
		} else {
			exceptFiles = append(exceptFiles, f.Name())
		}
	}

	if a.params.Debug {
		a.print("Count except files", fmt.Sprintf("%d", len(exceptFiles)))
		a.print("Except files", "")
		for _, n := range exceptFiles {
			a.debug(n)
		}
		fmt.Println()
	}

	countSerials := len(a.serials)
	a.print("Found serials", fmt.Sprintf("%d", countSerials))

	if countSerials == 0 {
		err = errors.New("serials not found")
	}

	return err
}

func (a *app) rename(path string, newPath string) error {
	var err error
	fmt.Print(a.color.GetColor(bashColor.Purple) + path)
	fmt.Print(a.color.GetColor(bashColor.Green) + " -> ")
	fmt.Println(a.color.GetColor(bashColor.Magenta) + newPath)

	if !a.params.Preview {
		err = os.Rename(path, newPath)
	}
	return err
}

func (a *app) renameSerials() error {
	a.color.Print(a.color.Green, "\n-------------------------------")
	a.color.Print(a.color.Green, "Rename videos:")

	for _, v := range a.serials {
		name := fmt.Sprintf(a.mask, v.Info.Season, v.Info.Episode) + v.Ext

		if err := a.rename(a.inputDir+"/"+v.Name, a.inputDir+"/"+name); err != nil {
			return err
		}
	}

	a.color.Print(a.color.Green, "\nRename input directory:")

	return a.rename(a.inputDir, "Season "+fmt.Sprintf("%d", a.maskProperties.seasonMax))
}

func (a *app) getCountLetterByNumber(number uint8) uint8 {
	countLetter := len(fmt.Sprintf("%d", number))
	if countLetter < 2 {
		countLetter = 2
	}
	return uint8(countLetter)
}

func (a *app) setMask() {
	a.mask = fmt.Sprintf(
		a.maskProperties.prefix+".s%%0%dde%%0%dd",
		a.getCountLetterByNumber(a.maskProperties.seasonMax),
		a.getCountLetterByNumber(a.maskProperties.episodesMax),
	)
}
