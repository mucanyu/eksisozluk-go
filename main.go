package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mucanyu/eksisozluk-go/model"
	"github.com/mucanyu/eksisozluk-go/scraper"
	"github.com/urfave/cli"
)

var (
	limitVal, pageVal int
	sukelaVal         bool
	kategoriVal       string

	version = "1.0.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "eksisozluk-go"
	app.Version = version
	app.CustomAppHelpTemplate = `VERSIYON:
			` + app.Version + `

KULLANIM:
			eksisozluk-go [gundem|baslik] [argumanlar...]

KOMUTLAR:
			gundem,  g	 Ekşisözlük'teki gündemleri listeler
			baslik,  b	 Başlık içerisindeki entryleri listeler
			version, v 	 Versiyon numarasını gösterir
			help,    h	 Kullanılabilen komutları listeler ya da bir komut için yardım yazısını gösterir
`
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "gundem",
			Aliases: []string{"g"},
			Usage:   "Eksisozluk'teki gundemleri listeler",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "limit, l",
					Value:       20,
					Usage:       "Listelenecek maksimum gundem sayısı.",
					Destination: &limitVal,
				},
				cli.IntFlag{
					Name:        "sayfa, s",
					Value:       1,
					Usage:       "Istediginiz sayfadaki populer basliklari getirir.",
					Destination: &pageVal,
				},
				cli.StringFlag{
					Name:        "kategori, k",
					Usage:       "Belirlediğiniz kategoride gündemdeki başlıkları getirir.",
					Destination: &kategoriVal,
				},
			},
			Action: func(c *cli.Context) error {
				if limitVal < 1 || pageVal < 1 {
					cli.ShowCommandHelp(c, "gundem")
					return errors.New("\nLimit veya belirtilen sayfanın değeri `1` değerinden az olmamalıdır")
				}
				params := model.GundemParams{Limit: limitVal, Page: pageVal, Kategori: kategoriVal}

				err := scraper.PrintGundem(&params)
				if err != nil {
					return err
				}
				return nil
			},
			CustomHelpTemplate: `GUNDEM:
		Eksisozluk'teki gundemleri listeler

KULLANIM:
		eksisozluk-go gundem [--limit=BASLIK_SAYISI] [--sayfa=GUNDEM_SAYFA_NO]

SECENEKLER:
		--kategori, -k	  Seçtiğiniz kategoriden popüler başlıkları listeler. (ilişkiler, spor, siyaset...)
		--limit,    -l    Listelenecek maksimum başlık sayısı. (varsayılan: 20)` + `
		--sayfa,    -s    Seçtiğiniz sayfadaki popüler başlıkları getirir. (varsayılan: 1)` + "\n",
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				cli.ShowCommandHelp(c, "gundem")
				return nil
			},
		},
		cli.Command{
			Name:  "baslik",
			Usage: "Basliktaki entryleri listeler",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "sayfa, s",
					Usage:       "Başlığın başlangıç noktasını belirler",
					Value:       1,
					Destination: &pageVal,
				},
				cli.IntFlag{
					Name:        "limit, l",
					Usage:       "Bir sayfada listelenecek maksimum entry sayısı",
					Value:       10,
					Destination: &limitVal,
				},
				cli.BoolFlag{
					Name:        "sukela",
					Usage:       "Başlık içerisinde en beğenilen entryleri sıralar",
					Destination: &sukelaVal,
				},
			},
			CustomHelpTemplate: `BASLIK:
		Seçtiğiniz başlığın entrylerini listeler.

KULLANIM:
		eksisozluk-go baslik BASLIK_ISMI [--sukela] [--sayfa=BASLIK_SAYFA_NO] [--limit=ENTRY_LIMITI]

SECENEKLER:
	--sukela			 Başlığın en çok olumlu oy olan entrylerini sıralar (varsayılan: false)
	--limit, -l    Listelenecek maksimum entry sayısı (varsayılan: 10)
	--sayfa, -s    Başlık içerisindeki sayfa seçimi (varsayılan: 1)` + "\n",

			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					cli.ShowCommandHelp(c, "baslik")
					return errors.New("\nERROR: Başlık ismini giriniz")

				} else if limitVal < 1 || pageVal < 1 {
					cli.ShowCommandHelp(c, "baslik")
					return errors.New("\nLimit veya belirtilen sayfanın değeri `1` değerinden az olmamalıdır")

				} else {
					args := strings.Join(c.Args(), " ")
					params := model.BaslikParams{Topic: args, Page: pageVal, Limit: limitVal, Sukela: sukelaVal}

					err := scraper.PrintTopic(&params)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					return nil
				}
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				cli.ShowCommandHelp(c, "baslik")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
