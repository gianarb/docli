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
    var ImgStruct digitalocean.Image
    var DropletStruct digitalocean.Droplet
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
                    opt.Page = page + 1
                }
                Table.Print()
            },
        },
        {
            Name:  "regions",
            Usage: "List of regions",
            Action: func(c *cli.Context) {
                Table := clitable.New([]string{
                    "Name",
                    "Slug",
                    "Available",
                })
                opt := &godo.ListOptions{}
                for {
                    regions, resp, _ := client.Regions.List(opt)
                    for _, r := range regions {
                        Table.AddRow(map[string]interface{}{
                            "Name": r.Name,
                            "Slug": r.Slug,
                            "Available": r.Available,
                        })
                    }
                    if resp.Links.IsLastPage() {
                        break
                    }
                    page, _ := resp.Links.CurrentPage()
                    opt.Page = page + 1
                }
                Table.Print()
            },
        },
        {
            Name:  "keys",
            Usage: "List of keys",
            Action: func(c *cli.Context) {
                if c.String("name") != "" && c.String("ssh_key") != "" {
                    request := &godo.KeyCreateRequest{};
                    request.Name = c.String("name")
                    request.PublicKey = c.String("ssh_key")
                    k, _, err := client.Keys.Create(request);
                    if err != nil {
                        fmt.Printf("ErrErr  %s \n", err)
                        return
                    }
                    fmt.Printf("Id: %d \n",  k.ID)
                    fmt.Printf("Name: %s \n", k.Name)
                    fmt.Printf("Fingerprint: %s \n", k.Fingerprint)
                    fmt.Printf("PublicKey: %s \n\n", k.PublicKey)
                    return
                }
                if c.String("delete") != "" {
                    client.Keys.DeleteByID(c.Int("delete"))
                    println("Key deleted")
                    return
                }
                if c.String("id") != "" {
                    k, _, _ := client.Keys.GetByID(c.Int("id"))
                    fmt.Printf("Id: %d \n",  k.ID)
                    fmt.Printf("Name: %s \n", k.Name)
                    fmt.Printf("Fingerprint: %s \n", k.Fingerprint)
                    fmt.Printf("PublicKey: %s \n\n", k.PublicKey)
                    return 
                }
                opt := &godo.ListOptions{}
                for {
                    keys, resp, _ := client.Keys.List(opt)
                    for _, k := range keys {
                        fmt.Printf("Id: %d \n",  k.ID)
                        fmt.Printf("Name: %s \n", k.Name)
                        fmt.Printf("Fingerprint: %s \n", k.Fingerprint)
                        fmt.Printf("PublicKey: %s \n\n", k.PublicKey)
                    }
                    if resp.Links.IsLastPage() {
                        break
                    }

                    page, _ := resp.Links.CurrentPage()

                    opt.Page = page + 1
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
                    Name: "ssh_key, s",
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


                Table := clitable.New([]string{
                    "Id",
                    "Name",
                    "Slug",
                    "Public",
                })
                opt := &godo.ListOptions{}
                for {
                    images, resp, _ := client.Images.List(opt)
                    for _, img := range images {
                        Table.AddRow(map[string]interface{}{
                            "Id": img.ID,
                            "Name": img.Name,
                            "Splug": img.Slug,
                            "Public": img.Public,
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
                    client.Droplets.Delete(c.Int("delete"))
                    println("Droplet deleted")
                    return
                }
                if c.String("id") != "" {
                    droplet, _, _ := client.Droplets.Get(c.Int("id"))
                    fmt.Printf("ID: %d \n", droplet.Droplet.ID);
                    fmt.Printf("Name: %s \n", droplet.Droplet.Name);
                    fmt.Printf("Memory: %d mb \n", droplet.Droplet.Memory);
                    fmt.Printf("Vcpu: %d \n", droplet.Droplet.Vcpus);
                    fmt.Printf("Region: %s \n", droplet.Droplet.Region.Name);
                    fmt.Printf("Status: %s \n", droplet.Droplet.Status);
                    fmt.Printf("IP: %s \n", droplet.Droplet.Networks.V4[0].IPAddress);
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
                Table := clitable.New([]string{
                    "Id",
                    "Name",
                    "Memory",
                    "Vcpus",
                    "Ip",
                    "Disk",
                    "Image",
                })
                opt := &godo.ListOptions{}
                for {
                    droplets, resp, _ := client.Droplets.List(opt)
                    if len(droplets) == 0 {
                        fmt.Printf("Zero droplets found\n") 
                        return
                    }
                    for _, dp := range droplets {
                        Table.AddRow(map[string]interface{}{
                            "Id": dp.ID,
                            "Name": dp.Name,
                            "Memory": dp.Memory,
                            "Vcpus": dp.Vcpus,
                            "Ip": dp.Networks.V4[0].IPAddress,
                            "Disk": dp.Disk,
                            "Image": dp.Image.Name,
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
