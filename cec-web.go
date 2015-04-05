package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/robbiet480/cec"
	"strings"
	"log"
	"os"
)

type Options struct {
	Host string `short:"i" long:"ip" description:"ip to listen on" default:"127.0.0.1"`
	Port string `short:"p" long:"port" description:"tcp port to listen on" default:"8080"`
	Adapter string `short:"a" long:"adapter" description:"cec adapter to connect to [RPI, usb, ...]"`
	Name string `short:"n" long:"name" description:"OSD name to announce on the cec bus" default:"REST Gateway"`
	Type string `short:"t" long:"type" description:"The device type to register as" default:"tuner"`
}

var options Options 
var parser = flags.NewParser(&options, flags.Default) 

func main() {
	if _, err := parser.Parse(); err != nil { 
		os.Exit(1) 
	} 

	cec.Open(options.Adapter, options.Name, options.Type)
	
	r := gin.Default()
	r.GET("/info", info)
	r.GET("/sourcestatus", source_status)
	r.GET("/power/:device", power_status)
	r.PUT("/power/:device", power_on)
	r.DELETE("/power/:device", power_off)
	r.PUT("/volume/up", vol_up)
	r.PUT("/volume/down", vol_down)
	r.PUT("/volume/mute", vol_mute)
	//r.POST("/key/:device", key_post)
	r.PUT("/key/:device/:key", key)
	r.PUT("/channel/:device/:channel", change_channel)
	r.POST("/transmit", transmit)

	r.Run(options.Host + ":" + options.Port)
}

func info(c *gin.Context) {
	c.JSON(200, cec.List())
}

func power_on(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))

	cec.PowerOn(addr)
	c.String(204, "")
}

func power_off(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))

	cec.Standby(addr)
	c.String(204, "")
}

func power_status(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))

	status := cec.GetDevicePowerStatus(addr)
	if status == "on" {
		c.String(204, "")
	} else if status == "standby" {
		c.String(404, "")
	} else {
		c.String(500, "invalid power state")
	}
}

func source_status(c *gin.Context) {
	active_devices := cec.GetActiveDevices()

	for address, active := range active_devices {
		if (active) && (cec.IsActiveSource(address)) {
			c.String(200, "INPUT HDMI "+strings.Split(cec.GetDevicePhysicalAddress(address), ".")[0]);
		}
	}
}

func change_channel(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))
	channel := c.Params.ByName("channel")

	for _, number := range channel {
		cec.Key(addr, "0x2"+string(number))
	}

	c.String(200, channel)
}

func transmit(c *gin.Context) {
	var commands []string
	c.Bind(&commands)

	for _, val := range commands {
		cec.Transmit(val)
	}
	c.String(204, "")
}	

func vol_up(c *gin.Context) {
	cec.VolumeUp()
	c.String(204, "")
}	

func vol_down(c *gin.Context) {
	cec.VolumeDown()
	c.String(204, "")
}	

func vol_mute(c *gin.Context) {
	cec.Mute()
	c.String(204, "")
}	

func key(c *gin.Context) {
	addr := cec.GetLogicalAddressByName(c.Params.ByName("device"))
	key := c.Params.ByName("key")

	cec.Key(addr, key)
	c.String(204, "")
}
