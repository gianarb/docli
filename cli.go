package main

import (
    "fmt"
    "github.com/codegangsta/cli"
    "os"
    "github.com/crackcomm/go-clitable"
    "github.com/gianarb/digitalocean-go"
    "code.google.com/p/goauth2/oauth"
    "github.com/digitalocean/godo"
)

func main() {
    app := cli.NewApp()
    app.Name = "docli"
    app.Usage = "Digitalocean in your command line"
    var ImgsStruct digitalocean.Images
    var ImgStruct digitalocean.Image
    var KeyStruct digitalocean.Key
    var KeysStruct digitalocean.Keys
    var DropletStruct digitalocean.Droplet
    var DropletsStruct digitalocean.Droplets
    var RegionsStruct digitalocean.Regions
    var configuration Configuration
    configuration.Parse()

    t := &oauth.Transport{
        Token: &oauth.Token{AccessToken: configuration.Token},
    }
    client := godo.NewClient(t.Client())

    app.Commands = []cli.Command{
        {
            Name:  "sizes",
            Usage: "List of sizes",
            Action: func(c *cli.Context) {
                Table := clitable.New([]string{
                    "Slug",
                    "Memory",
                    "Vcpus",
                    "Disks",
                    "Price monthly",
                    "Price hourly",
                })
                opt := &godo.ListOptions{}
                for {
                    sizes, resp, _ := client.Sizes.List(opt)
                    for _, size := range sizes {
                        Table.AddRow(map[string]interface{}{
                            "Slug": size.Slug,
                            "Memory": size.Memory,
                            "Vcpus": size.Vcpus,
                            "Disks": size.Disk,
                            "Price monthly": size.PriceMonthly,
                            "Price hourly": size.PriceHourly,
                        })
                    }
                    if resp.Links.IsLastPage() {
                        break
                    }

                    page, _ := resp.Links.CurrentPage()

                    // set the page we want for the next request
                    opt.Page = page + 1
                }
                Table.Print()
            },
        },
        {
            Name:  "regions",
            Usage: "List of regions",
            Action: func(c *cli.Context) {
                regions := RegionsStruct.List(configuration.Token)
                Table := clitable.New([]string{
                    "Name",
                    "Slug",
                    "Available",
                })
                for _, img := range regions.Pool {
                    Table.AddRow(map[string]interface{}{
                        "Name": img.Name,
                        "Slug": img.Slug,
                        "Available": img.Available,
                    })
                }
                Table.Print()
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
            Name:  "keys",
            Usage: "List of keys",
            Action: func(c *cli.Context) {
                if c.String("name") != "" && c.String("ssh_key") != "" {
                    var request digitalocean.KeyRequest
                    request.Name = c.String("name")
                    request.Public_key = c.String("ssh_key")
                    
                    data := KeyStruct.Create(request, configuration.Token);
                    fmt.Printf("%s \n", data)
                    return
                }
                if c.String("delete") != "" {
                    KeyStruct.Delete(c.String("delete"), configuration.Token)
                    println("Key deleted")
                    return 
                }
                if c.String("id") != "" {
                    keyString := KeyStruct.Single(c.String("id"), configuration.Token)
                    fmt.Printf("%s \n", keyString);
                    return 
                }
                keys := KeysStruct.List(configuration.Token) 
                for _, k := range keys.Pool {
                        fmt.Printf("Id: %d \n",  k.Id)
                        fmt.Printf("Name: %s \n", k.Name)
                        fmt.Printf("Fingerprint: %s \n", k.Fingerprint)
                        fmt.Printf("PublicKey: %s \n\n", k.PublicKey)
                }
            },
            Flags: []cli.Flag {
                cli.StringFlag {
                    Name: "delete, d",
                    Value: "",
                    Usage: "Delete single image --delete=idImage",
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
                    Name: "ssh_key, m",
                    Value: "1601",
                    Usage: "Id of image",
                },
            },
        },
        {
            Name:  "images",
            Usage: "List of images",
            Action: func(c *cli.Context) {
                if c.String("delete") != "" {
                    ImgStruct.Delete(c.String("delete"), configuration.Token)
                    println("Image deleted")
                    return 
                }
                if c.String("id") != "" {
                    imageString := ImgStruct.Single(c.String("id"), configuration.Token)
                    fmt.Printf("%s \n", imageString);
                    return 
                }
                images := ImgsStruct.List(configuration.Token) 
                Table := clitable.New([]string{
                    "Id",
                    "Name",
                    "Slug",
                    "Public",
                    "Created",
                })
                for _, img := range images.Pool {
                    Table.AddRow(map[string]interface{}{
                        "Id": img.Id,
                        "Name": img.Name,
                        "Splug": img.Slug,
                        "Public": img.Public,
                        "Created": img.Created,
                    })
                }
                Table.Print()
            },
            Flags: []cli.Flag {
                cli.StringFlag {
                    Name: "delete, d",
                    Value: "",
                    Usage: "Delete single image --delete=idImage",
                },
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
                    Table := clitable.New([]string{
                        "Id",
                        "Name",
                        "Memory",
                        "Vcpus",
                        "Ip",
                        "Disk",
                        "Image",
                    })
                    for _, dp := range droplets.Pool {
                        Table.AddRow(map[string]interface{}{
                            "Id": dp.Id,
                            "Name": dp.Name,
                            "Memory": dp.Memory,
                            "Vcpus": dp.Vcpus,
                            "Ip": dp.Networks.V4[0].IpAddress,
                            "Disk": dp.Disk,
                            "Image": dp.Image.Name,
                        })
                    }
                    Table.Print()
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
