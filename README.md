## digitalocean-cli
Manage your digitalocean world in console

### Install
```bash
go get github.com/gianarb/docli
```
Create in your homedir ``` .digitalocean-cli  ``` file with this configuration
```json
{
    "token": "your-digitalocean-token"
}
```

### Contrib
Clone this repo and build it
```bash
git clone https://github.com/gianarb/docli
cd docli
go build
```

Read help after build
* List of regions: ```docli regions```
* List of images: ```docli images```
* Single image: ```docli images --id=idimage```
* Lost of droplets: ```docli dropltes```
* Single droplet: ```docli dropltes --id=id```
* Create droplet: ```docli dropltes --name="my-drop" [--size="512Mb" --region="lon1" --images=idImages ...]```
* Delete droplet: ```docli dropltes --delete=id```
* Shutdown droplet: ```docli dropltes --id=idDroplet --action=shutdown```
* Power On droplet: ```docli dropltes --id=idDroplet --action=power_on```
* Power Off droplet: ```docli dropltes --id=idDroplet --action=power_off```

