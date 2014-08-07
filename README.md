## digitalocean-cli
Manage your digitalocean world in console

### Install
Create in your homedir ``` .digitalocean-cli  ``` file with this configuration
```json
{
    "token": "your-digitalocean-token"
}
```

Clone this repo and build it
```bash
git clone https://github.com/gianarb/digitalocean-cli
cd digitalocean-cli
go build
```
Read help after build
```
./digitalocean-cli
```

* List of images: ```./digitalocean-cli images```
* Single image: ```./digitalocean-cli images --id=idimage```
* Lost of droplets: ```./digitalocean-cli dropltes```
* Single droplet: ```./digitalocean-cli dropltes --id=id```
* Create droplet: ```./digitalocean-cli dropltes --name="my-drop" [--size="512Mb" --region="lon1" --images=idImages ...]```
* Delete droplet: ```./digitalocean-cli dropltes --delete=id```
* Shutdown droplet: ```./digitalocean-cli dropltes --id=idDroplet --action=shutdown```
* Power On droplet: ```./digitalocean-cli dropltes --id=idDroplet --action=power_on```
* Power Off droplet: ```./digitalocean-cli dropltes --id=idDroplet --action=power_off```

Contribute please!
[digitalocean-go](http://github.com/gianarb/digitalbees-go) is dependency for manage digitalocean api integration
 
