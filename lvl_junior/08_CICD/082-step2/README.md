# Gitlab self-hosted

Конфигурация ВМ: 
4 CPU 
8Gb RAM 
100 Gb HDD
OS: Ubuntu 

Далее по гайду: 
https://about.gitlab.com/install/#ubuntu
https://jino.ru/spravka/articles/gitlab_ce.html#%D0%BF%D0%B0%D0%BA%D0%B5%D1%82%D1%8B-%D0%BF%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B9-%D0%B4%D0%B6%D0%B8%D0%BD%D0%BE

1. Устанавливаем зависимости: 
```sh
sudo apt-get update
sudo apt-get install -y curl openssh-server ca-certificates tzdata perl
```
2. Устанавливаем пакет с Гитлабом и сам Гитлаб
```sh
curl -s https://packages.gitlab.com/install/repositories/gitlab/gitlab-ce/script.deb.sh | sudo bash
```

```sh
sudo apt-get install gitlab-ce
```

```sh
sudo gitlab-ctl reconfigure
```

Пароль: 
```sh
cat /etc/gitlab/initial_root_password
```

3. Настройка SSL (https://docs.gitlab.com/ee/topics/offline/quick_start_guide.html#enabling-ssl): 
3.1 Выпускаем сертификат и ключ
```sh
vi openssl.conf
```

```
[req]
default_bits        = 2048
default_md          = sha256
default_keyfile     = gitlab.example.com.key
prompt              = no
distinguished_name  = dn

[dn]
C                   = RU
ST                  = Moscow
L                   = Moscow
O                   = Your Company Name
OU                  = DevOps
CN                  = gitlab.example.com
emailAddress        = admin@gitlab.example.com

[ext]
subjectAltName      = @alt_names

[alt_names]
DNS.1               = gitlab.example.com
DNS.2               = www.gitlab.example.com
```

```sh
sudo mkdir -p /etc/gitlab/ssl
sudo chmod 755 /etc/gitlab/ssl
```

```sh
sudo openssl req -x509 -nodes -days 365 \
    -newkey rsa:2048 \
    -keyout /etc/gitlab/ssl/gitlab.example.com.key \
    -out /etc/gitlab/ssl/gitlab.example.com.crt \
    -config openssl.conf \
    -extensions ext
```
3.2 Обновляем конфиг
```sh
vi /etc/gitlab/gitlab.rb

# external_url 'https://gitlab.example.com'
```

```sh
sudo gitlab-ctl reconfigure
sudo gitlab-ctl status
```