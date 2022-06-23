## Install
* use go
    ```shell
    go get github.com/evgeny-klyopov/rename-serial-app
    ```
* source
    * Mac OS (m1)
      ```shell
      curl -s https://api.github.com/repos/evgeny-klyopov/rename-serial-app/releases/latest \
      | grep browser_download_url \
      | grep rsd.macos-arm64.tar.gz \
      | cut -d '"' -f 4 \
      | wget -qi - 
      tar -xvf rsd.macos-arm64.tar.gz && mv rsd /usr/bin/rsd 
      ```
    * Mac OS
      ```shell
      curl -s https://api.github.com/repos/evgeny-klyopov/rename-serial-app/releases/latest \
      | grep browser_download_url \
      | grep rsd.macos-amd64.tar.gz \
      | cut -d '"' -f 4 \
      | wget -qi - 
      tar -xvf rsd.macos-amd64.tar.gz && mv rsd /usr/bin/rsd 
      ```
      ```
    * Linux
      ```shell
      curl -s https://api.github.com/repos/evgeny-klyopov/rename-serial-app/releases/latest \
      | grep browser_download_url \
      | grep rsd.linux-amd64.tar.gz \
      | cut -d '"' -f 4 \
      | wget -qi - 
      tar -xvf rsd.linux-amd64.tar.gz && mv rsd /usr/bin/rsd 
      ```
    * Windows
      ```shell
      curl -s https://api.github.com/repos/evgeny-klyopov/rename-serial-app/releases/latest \
      | grep browser_download_url \
      | grep rsd.windows-amd64.tar.gz \
      | cut -d '"' -f 4 \
      | wget -qi -
      tar -xvf rsd.windows-amd64.tar.gz && mv rsd.exe /usr/bin/rsd
      ```