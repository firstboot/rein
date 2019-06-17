# rein

![rein](https://raw.githubusercontent.com/firstboot/rein/master/rein-logo.png)

* This program is used to efficiently redirect connections from one IP address/port combination to another. 



It is useful when operating virtual servers, dockers, firewalls and the like.  It was developed by golang.



### Simple deployment 

#### CentOS/RHEL/Ubuntu

```shell
cd ~
wget https://github.com/firstboot/rein/blob/master/release/linux/rein
wget https://github.com/firstboot/rein/blob/master/release/linux/rein.json
chmod +x rein
-- modify your conf
-- nano rein.conf
./rein
```

#### Windows

- download 2 files

  https://github.com/firstboot/rein/blob/master/release/windows/rein.zip

  https://github.com/firstboot/rein/blob/master/release/windows/rein.json

- decompress `rein.zip`

- modifying conf `rein.json` 

  ```json
  {
      "upstream": [
          {"source": "0.0.0.0:8150", "target": "127.0.0.1:22"}
      ]
  }
  ```

- running

  ```powershell
  rein
  ```

  














