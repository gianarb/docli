package main

import (
    "fmt"
    "github.com/codegangsta/cli"
    "os"
    "gitgub.com/gianarb/digitalocean-go"
)

func main() {
    app := cli.NewApp()
    app.Name = "digitalocean-cli"
    app.Usage = "Digitalocean in your command line"
    var ImgsStruct digitalocean.Images
    var ImgStruct digitalocean.Image
    var DropletStruct digitalocean.Droplet
    var DropletsStruct digitalocean.Droplets
    var RegionsStruct digitalocean.Regions
    var configuration Configuration
    configuration.Parse()
    app.Commands = []cli.Command{
        {
            Name:  "regions",
            Usage: "List of regions",
            Action: func(c *cli.Context) {
                regions := RegionsStruct.List(configuration.Token) 
                for _, img := range regions.Pool {
                    fmt.Printf("%s \t %s \t %t \n", img.Name, img.Slug, img.Available) 
                }
            },
            Flags: []cli.Flag {
                cli.StringFlag {
                    Name: "id",
                    Value: "",
                    Usage: "Resume single droples from id",
                },
            },
        },
        {
            Name:  "images",
            Usage: "List of images",
            Action: func(c *cli.Context) {
                if c.String("id") != "" {
                    imageString := ImgStruct.Single(c.String("id"), configuration.Token)
                    fmt.Printf("%s \n", imageString);
                    return 
                }
                images := ImgsStruct.List(configuration.Token) 
                for _, img := range images.Pool {
                    fmt.Printf("%d \t %s \t %s \t %t \t %s \n", img.Id, img.Name, img.Slug, img.Public, img.Created) 
                }
            },
            Flags: []cli.Flag {
                cli.StringFlag {
                    Name: "id",
                    Value: "",
                    Usage: "Resume single droples from id",
                },
            },
        },
        {
            Name:  "droplets",
            Usage: "List of droplets",
            Action: func(c *cli.Context) {
                if c.String("action") != "" {
                    response := DropletStruct.DoAction(c.String("action"), c.String("id"), configuration.Token)
                    fmt.Printf("%s \n", response);
                    return
                }
                if c.String("delete") != "" {
                    DropletStruct.Delete(c.String("delete"), configuration.Token)
                    println("Droplet deleted")
                    return
                }
                if c.String("id") != "" {
                    dropletString := DropletStruct.Single(c.String("id"), configuration.Token)
                    fmt.Printf("%s \n", dropletString);
                    return
                }
                if c.String("name") != "" {
                    var request digitalocean.DropletRequest
                    request.Name = c.String("name")
                    request.Region = c.String("region")
                    request.Size = c.String("size")
                    request.Image = c.String("image")
                    
                    data := DropletStruct.Create(request, configuration.Token);
                    fmt.Printf("%s \n", data)
                    return
                }
                droplets := DropletsStruct.List(configuration.Token)      
                if len(droplets.Pool) == 0 {
                        fmt.Printf("Zero droplets found\n") 
                } else {
                    fmt.Printf("ID \t NAME \t IP \t Mamory \t CPU \t DISK \t IMAGE NAME \n\n") 
                    for _, dp := range droplets.Pool {
                        fmt.Printf("%d \t %s \t %s \t %d \t %d \t %d \t %s \n", dp.Id, dp.Name, dp.Ip, dp.Memory, dp.Vcpus, dp.Disk, dp.Image.Name) 
                    }
                }
            },
            Flags: []cli.Flag {
                cli.StringFlag {
                    Name: "delete, d",
                    Value: "",
                    Usage: "Delete single droplet --delete=idDroplet",
                },
                cli.StringFlag {
                    Name: "id",
                    Value: "",
                    Usage: "Resume single droples from id",
                },
                cli.StringFlag {
                    Name: "name, n",
                    Value: "",
                    Usage: "Name of new droplet",
                },
                cli.StringFlag {
                    Name: "image, m",
                    Value: "1601",
                    Usage: "Id of image",
                },
                cli.StringFlag {
                    Name: "size, s",
                    Value: "512Mb",
                    Usage: "Mb of VM",
                },
                cli.StringFlag {
                    Name: "region, r",
                    Value: "lon1",
                    Usage: "Region of destion",
                },
                cli.StringFlag {
                    Name: "action",
                    Value: "",
                    Usage: "Region of destion",
                },
            },
        },
    }

    app.Run(os.Args)
}
