package main

import (
    "encoding/json"
    "io/ioutil"
    "fmt"
    "os/user"
)

type Configuration struct {
    Token string `json:"token"`
}

func (c *Configuration) Parse() error {
    usr, err := user.Current()
    if err != nil {
        return err
    }

    file, err := ioutil.ReadFile(fmt.Sprintf("%s/.digitalocean-cli", usr.HomeDir))
    
    if err != nil {
        return err
    }
    
    err = json.Unmarshal(file, &c)

    if err != nil {
        return err
    }
    return nil
}
