package main

import (
	"fmt"
	"os"
	"time"

	"github.com/packethost/packngo"
	"github.com/urfave/cli"
)

func waitDeviceActive(id string, c *packngo.Client, niter int, interval time.Duration) (*packngo.Device, error) {
	for i := 0; i < niter; i++ {
		<-time.After(interval)
		d, _, err := c.Devices.Get(id, nil)
		if err != nil {
			return nil, err
		}
		if d.State == "active" {
			return d, nil
		}
	}
	return nil, fmt.Errorf("device %s is still not active after %d * %v", id, niter, interval)
}

func reinstallDevice(c *cli.Context) error {
	id := c.String("id")
	if len(id) == 0 {
		return fmt.Errorf("You must set the ID in the --id flag")
	}
	df := c.Bool("deprovision-fast")
	pd := c.Bool("preserve-data")
	wa := c.Bool("wait")
	drf := packngo.DeviceReinstallFields{
		DeprovisionFast: df,
		PreserveData:    pd,
	}
	cl, err := packngo.NewClient()
	if err != nil {
		return err
	}
	_, err = cl.Devices.Reinstall(id, &drf)
	if err != nil {
		return err
	}
	if wa {
		_, err = waitDeviceActive(id, cl, 10, 10*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	app := &cli.App{
		Action: reinstallDevice,
		Name:   "reinstall",
		Usage:  "issues Packet API action to reisntall a device",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Usage: "ID of the device to reinstall",
			},
			&cli.BoolFlag{
				Name:  "deprovision-fast, d",
				Usage: "disk wipes will be skipped during a reinstall",
			},
			&cli.BoolFlag{
				Name:  "preserve-data, p",
				Usage: "no non-root disks will be touched during a reinstall",
			},
			&cli.BoolFlag{
				Name:  "wait, w",
				Usage: "wait until the device gets back the the active state",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
