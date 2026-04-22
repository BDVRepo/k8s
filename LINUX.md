# Linux — справочник с задачами

> Структура: теория → команды → задача → решение. Без временных рамок.

Связка: [BASH.md](BASH.md)

---

## Содержание

1. [Навигация и файлы](#1-навигация-и-файлы)
2. [Права доступа](#2-права-доступа)
3. [Процессы](#3-процессы)
4. [Сигналы и фоновые задачи](#4-сигналы-и-фоновые-задачи)
5. [Диск и память](#5-диск-и-память)
6. [Архивы](#6-архивы)
7. [Сеть](#7-сеть)
8. [Пакеты](#8-пакеты)
9. [systemd и сервисы](#9-systemd-и-сервисы)
10. [Мониторинг системы](#10-мониторинг-системы)
11. [Диагностика: типичные сценарии](#11-диагностика-типичные-сценарии)
12. [Скрипты: простые](#12-скрипты-простые)
13. [Скрипты: сложные и самописные утилиты](#13-скрипты-сложные-и-самописные-утилиты)
14. [Cron](#14-cron)
15. [Шпаргалка](#15-шпаргалка)

---

## 1. Навигация и файлы

```bash
pwd                      # где я сейчас
ls                       # содержимое папки
ls -la                   # подробно + скрытые файлы
ls -lh                   # размеры в читаемом виде
ls -lt                   # сортировка по времени изменения
cd /var/log              # перейти в папку
cd ~                     # домой
cd -                     # в предыдущую папку

mkdir mydir              # создать папку
mkdir -p a/b/c           # создать все вложенные сразу
touch file.txt           # создать пустой файл
cp a.txt b.txt           # копия файла
cp -r dir1/ dir2/        # копия папки рекурсивно
mv old.txt new.txt       # переместить / переименовать
rm file.txt              # удалить файл
rm -rf mydir             # удалить папку со всем содержимым (осторожно!)
ln -s /path/to/file link # символическая ссылка
ln /path/to/file link    # жёсткая ссылка

cat file.txt             # показать файл
less file.txt            # просмотр с прокруткой (q — выйти)
head -n 20 file.txt      # первые 20 строк
tail -n 20 file.txt      # последние 20 строк
tail -f app.log          # следить за логом в реальном времени
wc -l file.txt           # количество строк в файле
```

**Чем hard link отличается от symbolic link:**
- **hard link**: ещё одно имя для того же файла (одинаковый inode). Удалишь оригинал — данные не пропадут, пока есть хотя бы одно имя.
- **symbolic link**: просто ярлык на путь. Если оригинал удалить — ссылка станет битой.

```bash
ls -li file1 hardlink    # одинаковый inode-номер = hard link
ls -la symlink           # видно "-> /path/to/original"
```

**Задача:** создай файл `test.txt`, сделай к нему hard link и символическую ссылку. Удали оригинал. Проверь что hard link жив, а symbolic ссылка — битая.

<details>
<summary>Решение</summary>

```bash
echo "data" > test.txt
ln test.txt hardlink.txt
ln -s test.txt symlink.txt

ls -li test.txt hardlink.txt symlink.txt   # inode у первых двух одинаковый

rm test.txt

cat hardlink.txt    # работает — данные живы
cat symlink.txt     # ошибка — битая ссылка
```
</details>

---

## 2. Права доступа

Строка `ls -l` выглядит так:
```
-rwxr-x--- 1 alice devops 1234 Apr 21 deploy.sh
```

- Первый символ: `-` файл, `d` директория, `l` ссылка
- Следующие 9: `rwx` владелец / `r-x` группа / `---` остальные
- `r`=4, `w`=2, `x`=1 — складываются для цифрового вида

**Особенности на каталоге:**
- `r` — можно видеть список файлов
- `x` — можно войти (`cd`) и обращаться к файлам
- `w` — можно создавать и удалять файлы внутри

```bash
chmod +x script.sh              # добавить исполнение
chmod 644 file.txt              # rw-r--r--
chmod 755 script.sh             # rwxr-xr-x
chmod -R 755 dir/               # рекурсивно для папки
chmod u+x,g-w file.sh           # владельцу +x, группе -w

chown user:group file.txt       # сменить владельца и группу
chown -R user:group dir/        # рекурсивно

umask                           # посмотреть текущую маску (022 = файлы 644, папки 755)
umask 027                       # установить маску

getfacl file.txt                # посмотреть ACL
setfacl -m u:bob:rx file.txt    # дать bob права rx через ACL
```

**Задача:** создай скрипт `deploy.sh`. Дай права: владелец — читать/писать/запускать; группа — только читать и запускать; остальные — ничего. Проверь командой `ls -l`.

<details>
<summary>Решение</summary>

```bash
touch deploy.sh
chmod 750 deploy.sh   # rwxr-x---
ls -l deploy.sh
```
</details>

---

## 3. Процессы

**Как запускается процесс:**
1. Shell вызывает `fork()` — ядро клонирует текущий процесс, создаётся новый PID
2. В дочернем процессе вызывается `exec()` — загружается программа (PID не меняется)

**Зомби** — процесс завершился, но родитель ещё не «прочитал» его статус. Запись в таблице процессов висит. Лечится правкой родителя или перезапуском.

```bash
ps aux                          # все процессы
ps aux | grep nginx             # найти nginx
ps -ef --forest                 # дерево процессов (кто кого запустил)
pgrep -l nginx                  # PID + имя по маске
pidof nginx                     # только PID

top                             # живая статистика (q для выхода)
htop                            # то же но красивее (если установлен)

ps aux --sort=-%cpu | head -5   # топ по CPU
ps aux --sort=-%mem | head -5   # топ по памяти

cat /proc/<PID>/cmdline | tr '\0' ' '   # с какими аргументами запущен
ls -la /proc/<PID>/exe                  # что за программа
cat /proc/<PID>/status | grep -E 'Pid|PPid|State'
```

**Задача:** найди все процессы пользователя root, отсортируй по памяти, выведи первые 5.

<details>
<summary>Решение</summary>

```bash
ps aux | awk '$1=="root"' | sort -k4 -rn | head -5
```
</details>

---

## 4. Сигналы и фоновые задачи

| Сигнал | Номер | Смысл |
|--------|-------|--------|
| `SIGTERM` | 15 | вежливо попросить завершиться (можно перехватить) |
| `SIGKILL` | 9 | жёсткое убийство (нельзя перехватить) |
| `SIGHUP` | 1 | перезагрузить конфиг (или завершить если нет обработчика) |
| `SIGINT` | 2 | Ctrl+C |
| `SIGSTOP` | 19 | приостановить процесс |

```bash
kill -TERM <pid>           # вежливо (то же что kill без флага)
kill -9 <pid>              # только если TERM не помог
kill -HUP <pid>            # попросить перечитать конфиг (nginx, sshd)
killall nginx              # убить все процессы с именем nginx
pkill -f "python script"   # убить по части командной строки
kill -l                    # список всех сигналов

# Фоновые задачи
sleep 300 &                # запустить в фоне
jobs -l                    # список фоновых задач текущего shell
fg %1                      # вернуть задачу 1 на передний план
bg %1                      # продолжить задачу 1 в фоне
Ctrl+Z                     # приостановить текущий процесс

nohup ./longscript.sh > out.log 2>&1 &   # запустить, не умирать при закрытии терминала
disown -h %1               # отвязать от shell (альтернатива nohup после запуска)
```

**Задача:** запусти `sleep 500` в фоне. Найди его PID через `jobs -l` и через `ps`. Отправь SIGTERM. Проверь код выхода `$?`.

<details>
<summary>Решение</summary>

```bash
sleep 500 &
jobs -l
ps aux | grep sleep
kill -TERM <PID>
echo $?    # 0 если kill сработал
```
</details>

---

## 5. Диск и память

**Важно:** `df` и `du` могут «не совпадать» — `df` смотрит ФС целиком, `du` обходит дерево. Если файл удалили, но процесс ещё держит `fd` — блоки считаются занятыми в `df`.

```bash
df -h                           # свободное место (человекочитаемо)
df -i                           # inode — другая причина «диск полон»
du -sh /var/log                 # сколько занимает папка
du -sh /var/log/* | sort -rh | head -10   # самые тяжёлые папки
du -h --max-depth=1 /           # по верхнеуровневым каталогам

free -h                         # RAM и Swap
cat /proc/meminfo | head -5     # детально
vmstat 1                        # CPU, память, swap (каждую секунду)
iostat -xz 1                    # нагрузка на диск

lsof +L1                        # файлы удалённые, но ещё открытые процессами
```

**Задача:** найди 5 самых больших файлов в `/var/log` (или `/tmp`).

<details>
<summary>Решение</summary>

```bash
find /var/log -type f -printf '%s %p\n' 2>/dev/null | sort -rn | head -5
```
</details>

---

## 6. Архивы

`tar` — основной инструмент. Флаги:
- `c` — создать
- `x` — распаковать
- `z` — gzip (`.tar.gz`)
- `j` — bzip2 (`.tar.bz2`)
- `f` — имя файла
- `v` — verbose (показывать что делает)

```bash
tar -czf archive.tar.gz folder/          # создать .tar.gz
tar -cjf archive.tar.bz2 folder/        # создать .tar.bz2 (лучше сжатие, медленнее)
tar -xzf archive.tar.gz                  # распаковать .tar.gz
tar -xzf archive.tar.gz -C /target/dir  # распаковать в конкретную папку
tar -tzf archive.tar.gz                  # посмотреть содержимое без распаковки
tar -czf archive.tar.gz -C /parent dir  # архив с относительными путями

gzip file.txt        # сжать → file.txt.gz (оригинал удаляется)
gunzip file.txt.gz   # распаковать
zip -r arch.zip dir/ # zip архив
unzip arch.zip       # распаковать zip
```

**Задача:** заархивируй папку `/tmp` в файл `/tmp/backup.tar.gz`, не включая сам файл архива. Проверь содержимое архива.

<details>
<summary>Решение</summary>

```bash
tar -czf /tmp/backup.tar.gz -C / tmp --exclude=tmp/backup.tar.gz
tar -tzf /tmp/backup.tar.gz | head -10
```
</details>

---

## 7. Сеть

**Состояния TCP соединений:**
- `LISTEN` — ждёт входящих
- `ESTABLISHED` — соединение установлено
- `TIME_WAIT` — соединение закрыто, но ждёт возможных запоздавших пакетов (нормально!)
- `CLOSE_WAIT` — наша сторона ждёт закрытия (может означать баг в приложении)

```bash
ip addr show                    # IP адреса всех интерфейсов
ip route                        # таблица маршрутов
ip link show                    # состояние сетевых интерфейсов

ss -tulpen                      # все порты: tcp+udp+listening+pid
ss -tlnp                        # только listening tcp, без имён
ss -tlnp 'sport = :80'          # кто слушает порт 80
ss -tan state established       # только установленные соединения
ss -tan | awk '{print $1}' | sort | uniq -c   # количество по состояниям

curl -I https://example.com                         # HTTP заголовки
curl -o /dev/null -w "%{http_code}\n" http://...   # только код ответа
curl -v http://example.com                          # подробный вывод
curl -s http://api/endpoint | jq .                  # JSON ответ
wget -q -O /dev/null http://example.com             # проверить доступность

ping -c 4 8.8.8.8               # 4 пинга
ping -W 1 -c 1 host             # один пинг с таймаутом 1 сек
host google.com                 # DNS-разрешение
dig google.com                  # подробно DNS
nslookup google.com             # альтернатива dig
```

**Задача:** найди все процессы которые слушают TCP порты. Выведи номер порта и имя процесса.

<details>
<summary>Решение</summary>

```bash
ss -tlnp | awk 'NR>1 {print $4, $6}' | sed 's/.*://; s/pid=//'
# или проще:
ss -tlnp
```
</details>

---

## 8. Пакеты

```bash
# Debian / Ubuntu (apt)
apt update                       # обновить список пакетов
apt upgrade                      # обновить установленные
apt install nginx                # установить пакет
apt remove nginx                 # удалить (конфиги остаются)
apt purge nginx                  # удалить с конфигами
apt search nginx                 # поиск пакета
apt show nginx                   # информация о пакете
dpkg -l | grep nginx             # проверить установлен ли

# RHEL / CentOS / Fedora (dnf/yum)
dnf install nginx
dnf remove nginx
dnf update
rpm -qa | grep nginx             # список rpm-пакетов

# Проверить что установлено и откуда
which curl                       # путь к команде
command -v curl                  # то же, работает в скриптах
type curl                        # тип (alias/function/file)
```

---

## 9. systemd и сервисы

**Идея:** systemd управляет сервисами через unit-файлы (`/etc/systemd/system/`).

```bash
systemctl status nginx           # статус сервиса
systemctl start nginx            # запустить
systemctl stop nginx             # остановить
systemctl restart nginx          # перезапустить
systemctl reload nginx           # перезагрузить конфиг без остановки
systemctl enable nginx           # автозапуск при старте системы
systemctl disable nginx          # отключить автозапуск
systemctl is-active nginx        # active или inactive
systemctl cat nginx              # посмотреть unit-файл

journalctl -u nginx              # логи сервиса
journalctl -u nginx -f           # следить за логами в реальном времени
journalctl -u nginx --since "10 min ago"
journalctl -b                    # логи с последней загрузки
journalctl -b -p err             # только ошибки
journalctl --disk-usage          # сколько места занимают логи
```

Пример простого unit-файла (`/etc/systemd/system/myapp.service`):

```ini
[Unit]
Description=My Application
After=network.target

[Service]
ExecStart=/usr/bin/myapp --config /etc/myapp/config.yaml
Restart=always
User=myapp
WorkingDirectory=/var/lib/myapp
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

После создания: `systemctl daemon-reload && systemctl enable --now myapp`

**Задача:** посмотри логи любого установленного сервиса за последние 10 минут. Найди ошибки.

<details>
<summary>Решение</summary>

```bash
systemctl list-units --type=service --state=running   # список активных сервисов
journalctl -u ssh --since "10 min ago"               # логи ssh
journalctl -u ssh --since "10 min ago" -p err        # только ошибки
```
</details>

---

## 10. Мониторинг системы

```bash
# Общая картина
uptime                          # время работы + load average
uname -a                        # версия ядра
lsb_release -a                  # версия дистрибутива
hostname                        # имя хоста

# CPU
nproc                           # количество ядер
cat /proc/cpuinfo | grep "model name" | head -1
mpstat 1 3                      # статистика CPU (нужен sysstat)
top -bn1 | grep "Cpu(s)"        # снимок загрузки CPU

# Память
free -h
cat /proc/meminfo | grep -E 'MemTotal|MemAvailable|SwapTotal'
vmstat 1 5                      # 5 снимков каждую секунду

# Диск
iostat -xz 1 3                  # нагрузка на диски (нужен sysstat)
iotop                           # кто читает/пишет (нужен iotop)
lsblk                           # блочные устройства

# Открытые файлы и дескрипторы
lsof -p <PID>                   # файлы одного процесса
lsof -i :80                     # кто держит порт 80
lsof -i tcp                     # все TCP-соединения
cat /proc/sys/fs/file-nr        # открытых fd: текущих / макс

# Отладка системных вызовов
strace ls                       # syscalls команды ls
strace -p <PID>                 # syscalls живого процесса
strace -e trace=openat,read ls  # только конкретные syscalls
strace -c ls                    # статистика по syscalls
```

**Задача:** найди топ-5 процессов по потреблению памяти и выведи: пользователь, PID, % памяти, команда.

<details>
<summary>Решение</summary>

```bash
ps aux --sort=-%mem | awk 'NR==1 || NR<=6 {printf "%-10s %-8s %-6s %s\n", $1, $2, $4, $11}'
```
</details>

---

## 11. Диагностика: типичные сценарии

### Порт занят — кто его держит

```bash
ss -tlnp | grep :8080
lsof -i :8080
# получили PID, смотрим подробнее:
ps -fp <PID>
ls -la /proc/<PID>/exe
cat /proc/<PID>/cmdline | tr '\0' ' '
```

### Место на диске кончилось

```bash
df -h                                       # 1. какой раздел полный
df -i                                       # 1b. может inode кончились
du -sh /var/log/* | sort -rh | head         # 2. тяжёлые папки
find /var/log -size +100M                   # 3. большие файлы
lsof +L1                                    # 4. удалённые но открытые файлы
# если lsof показал процесс:
kill -HUP <PID>   # или restart — освободит файл
```

### Сервис не поднимается

```bash
systemctl status myservice
journalctl -u myservice -n 50 --no-pager
journalctl -u myservice --since "5 min ago"
# если нужно ещё глубже:
strace -f -e trace=file /usr/bin/myservice   # что за файлы ищет
```

---

### Сервис запущен, но не отвечает — пошаговая отладка

Применимо к любому сервису (nginx, API, Postgres, Redis и т.д.).

**Шаг 1 — убедиться что сервис вообще живой**
```bash
systemctl status <service> --no-pager -l
```
Смотришь:
- `Active: active (running)` — живой, идём дальше.
- `Active: failed` — сервис упал, смотри шаг 2.
- `Active: inactive (dead)` — сервис не запущен, `systemctl start <service>`.

---

**Шаг 2 — читать логи сервиса**
```bash
journalctl -u <service> -n 100 --no-pager
journalctl -u <service> --since "10 min ago"
```
Что искать:
- строки с `error`, `failed`, `permission denied`, `address already in use`.
- Стек ошибки или конкретную причину отказа.

Если это приложение (не systemd-сервис) — смотри в его лог-файл:
```bash
tail -f /var/log/<appname>/<appname>.log
grep -i "error\|fatal" /var/log/<appname>/<appname>.log | tail -30
```

---

**Шаг 3 — проверить конфиг (если есть команда проверки)**
```bash
# nginx:
sudo nginx -t
# --no-pager: вывод сразу в терминал без less
# -l: не обрезать длинные строки
```
Аналоги для других сервисов:
```bash
apache2 -t              # Apache
sshd -t                 # SSH
haproxy -c -f /etc/haproxy/haproxy.cfg
```

---

**Шаг 4 — проверить слушает ли порт**
```bash
ss -tlnp | grep :<port>
```
Примеры:
```bash
ss -tlnp | grep :80     # HTTP
ss -tlnp | grep :443    # HTTPS
ss -tlnp | grep :5432   # Postgres
ss -tlnp | grep :6379   # Redis
```
Что означает вывод:
- `LISTEN 0 ... 0.0.0.0:80` — слушает на всех IPv4-интерфейсах. Хорошо.
- `LISTEN 0 ... [::]:80` — слушает IPv6. Тоже нормально.
- Строки нет совсем — сервис НЕ слушает этот порт: либо не запустился, либо конфиг слушает другой порт.

---

**Шаг 5 — проверить ответ локально**
```bash
curl -I http://127.0.0.1           # HTTP
curl -kI https://127.0.0.1         # HTTPS (без проверки сертификата)
curl -I http://127.0.0.1:8080      # нестандартный порт
```
- `200 OK` — сервис отвечает, проблема в чём-то другом (DNS, firewall, upstream).
- `Connection refused` — никто не слушает порт.
- `502/503` — сервис слушает, но backend/upstream не отвечает.
- Timeout — порт может быть заблокирован firewall.

---

**Шаг 6 — проверить ресурсы**
```bash
free -h                         # не кончилась ли память
df -h                           # не кончился ли диск
ps aux --sort=-%cpu | head -5   # не перегружен ли CPU
lsof -p <PID> | wc -l          # сколько открытых файлов (limit?)
```
Частые причины:
- OOM (Out of Memory) — процесс убивает ядро. Проверь: `dmesg -T | grep -i "killed process"`.
- Диск кончился → `df -h`.
- Лимит на открытые файлы → `ulimit -n`, `cat /proc/<PID>/limits | grep "open files"`.

---

**Шаг 7 — если непонятно что делает процесс**
```bash
strace -p <PID> -f -e trace=network,file 2>&1 | head -30
```
Покажет в реальном времени какие системные вызовы делает процесс: куда подключается, какие файлы открывает.

---

**Итоговый чеклист «сервис не отвечает»:**

| Шаг | Команда | Что ищем |
|-----|---------|----------|
| 1 | `systemctl status <svc> --no-pager -l` | `active/failed/inactive` |
| 2 | `journalctl -u <svc> -n 100 --no-pager` | ошибки, причина |
| 3 | `<app> -t` (если есть) | ошибки конфига |
| 4 | `ss -tlnp \| grep :<port>` | слушает ли порт |
| 5 | `curl -I http://127.0.0.1:<port>` | HTTP-ответ |
| 6 | `free -h`, `df -h`, `dmesg \| grep killed` | ресурсы/OOM |
| 7 | `strace -p <PID>` | syscalls если совсем непонятно |

### Высокий CPU

```bash
top -bn1 | head -20
ps aux --sort=-%cpu | head -10
strace -p <PID> -c -f       # что делает процесс
```

### Сеть не работает

```bash
ip addr                    # есть ли IP?
ping 8.8.8.8               # IP связность
ping google.com            # DNS работает?
curl -v http://example.com # HTTP?
ss -tulpen                 # что слушает локально?
```

---

## 12. Скрипты: простые

Создавай файлы в `/home/bdv/projects/k8s/bash/linux-scripts/`

---

### `check-port.sh` — кто слушает порт

```bash
#!/bin/bash
port="${1:?Использование: $0 <port>}"
echo "Порт $port:"
result=$(ss -tlnp "sport = :$port")
if [ -z "$result" ]; then
    echo "Свободен"
else
    echo "$result"
fi
```

---

### `diskusage.sh` — топ папок по размеру

```bash
#!/bin/bash
dir="${1:-.}"
n="${2:-10}"
echo "Топ-$n в $dir:"
du -h "$dir" --max-depth=1 2>/dev/null | sort -rh | head -"$n"
```

Запуск: `./diskusage.sh /var/log 5`

---

### `proc-info.sh` — информация о процессе

```bash
#!/bin/bash
pid="${1:?Использование: $0 <PID>}"

if [ ! -d "/proc/$pid" ]; then
    echo "Процесс $pid не найден"
    exit 1
fi

echo "PID:     $pid"
echo "Команда: $(cat /proc/$pid/cmdline | tr '\0' ' ')"
echo "Статус:  $(grep State /proc/$pid/status | awk '{print $2, $3}')"
echo "Память:  $(grep VmRSS /proc/$pid/status | awk '{print $2, $3}')"
echo "FD:      $(ls /proc/$pid/fd 2>/dev/null | wc -l) открытых файлов"
echo "Exe:     $(readlink /proc/$pid/exe 2>/dev/null)"
```

---

### `netcheck.sh` — проверить список хостов

```bash
#!/bin/bash
hosts=("8.8.8.8" "1.1.1.1" "google.com" "github.com")

for host in "${hosts[@]}"; do
    if ping -c 1 -W 1 "$host" > /dev/null 2>&1; then
        echo "  OK   $host"
    else
        echo "  FAIL $host"
    fi
done
```

---

### `user-info.sh` — пользователи в системе

```bash
#!/bin/bash
echo "=== Пользователи с оболочкой ==="
grep -v '/nologin\|/false' /etc/passwd | awk -F: '{print $1, $3, $6}' | column -t

echo ""
echo "=== Залогиненные сейчас ==="
who
```

---

### `find-big-files.sh` — большие файлы

```bash
#!/bin/bash
dir="${1:-/}"
size="${2:-100M}"
echo "Файлы больше $size в $dir:"
find "$dir" -type f -size +"$size" 2>/dev/null -printf '%s %p\n' \
    | sort -rn \
    | head -20 \
    | awk '{printf "%-10s %s\n", $1, $2}'
```

Запуск: `./find-big-files.sh /var/log 10M`

---

**Задача:** напиши скрипт `free-space.sh` который проверяет все разделы. Если занятость > 80% — выводит `ВНИМАНИЕ: /mountpoint занят на N%`, иначе — `OK: /mountpoint N%`.

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
df -h | tail -n +2 | while read -r line; do
    usage=$(echo "$line" | awk '{print $5}' | tr -d '%')
    mount=$(echo "$line" | awk '{print $6}')
    if ! echo "$usage" | grep -qE '^[0-9]+$'; then continue; fi
    if [ "$usage" -gt 80 ]; then
        echo "ВНИМАНИЕ: $mount занят на $usage%"
    else
        echo "OK:       $mount $usage%"
    fi
done
```
</details>

---

## 13. Скрипты: сложные и самописные утилиты

---

### `sysmon.sh` — мониторинг системы в реальном времени

```bash
#!/bin/bash
interval="${1:-5}"

while true; do
    clear
    echo "=== $(date '+%Y-%m-%d %H:%M:%S') | обновление каждые ${interval}с | Ctrl+C выход ==="
    echo ""

    # Uptime и load
    echo "Uptime: $(uptime -p)  Load: $(uptime | awk -F'load average:' '{print $2}')"
    echo ""

    # Память
    echo "--- Память ---"
    free -h | awk 'NR==1 || NR==2'
    echo ""

    # Диск
    echo "--- Диск ---"
    df -h | grep -v 'tmpfs\|udev'
    echo ""

    # Топ CPU
    echo "--- Топ 5 по CPU ---"
    ps aux --sort=-%cpu | awk 'NR>1 && NR<=6 {printf "%-10s %-6s %-6s %s\n", $1, $2, $3"%", $11}'
    echo ""

    # Топ памяти
    echo "--- Топ 5 по памяти ---"
    ps aux --sort=-%mem | awk 'NR>1 && NR<=6 {printf "%-10s %-6s %-6s %s\n", $1, $2, $4"%", $11}'

    sleep "$interval"
done
```

---

### `logwatch.sh` — следить за паттерном в логе

```bash
#!/bin/bash
LOG="${1:?Использование: $0 <log_file> [pattern]}"
PATTERN="${2:-ERROR}"

[ -f "$LOG" ] || { echo "Нет файла: $LOG"; exit 1; }

echo "Слежу за '$PATTERN' в $LOG (Ctrl+C для выхода)..."
tail -f "$LOG" | grep --line-buffered -i "$PATTERN" | while read -r line; do
    echo "$(date '+%H:%M:%S') | $line"
done
```

---

### `log-rotate.sh` — ротация логов вручную

```bash
#!/bin/bash
LOG="${1:?Использование: $0 <log_file> [max_mb]}"
MAX_MB="${2:-100}"

[ -f "$LOG" ] || { echo "Нет файла: $LOG"; exit 1; }

size_mb=$(du -m "$LOG" | awk '{print $1}')

if [ "$size_mb" -ge "$MAX_MB" ]; then
    ts=$(date +%Y%m%d_%H%M%S)
    backup="${LOG}.${ts}"
    mv "$LOG" "$backup"
    gzip "$backup"
    touch "$LOG"
    echo "Ротация: ${size_mb}MB → $backup.gz"
else
    echo "Ротация не нужна: ${size_mb}MB < ${MAX_MB}MB"
fi
```

---

### `ssl-check.sh` — когда истекает SSL-сертификат

```bash
#!/bin/bash
host="${1:?Использование: $0 <host> [port]}"
port="${2:-443}"

expiry=$(echo | openssl s_client -connect "$host:$port" -servername "$host" 2>/dev/null \
    | openssl x509 -noout -enddate 2>/dev/null \
    | cut -d= -f2)

[ -z "$expiry" ] && { echo "Не удалось получить сертификат $host:$port"; exit 1; }

expiry_epoch=$(date -d "$expiry" +%s 2>/dev/null)
now_epoch=$(date +%s)
days=$(( (expiry_epoch - now_epoch) / 86400 ))

if [ "$days" -lt 14 ]; then
    echo "ВНИМАНИЕ $host: осталось $days дней! (истекает $expiry)"
else
    echo "OK $host: $days дней (истекает $expiry)"
fi
```

Запуск: `./ssl-check.sh github.com`

---

### `backup.sh` — умный бэкап с ротацией

```bash
#!/bin/bash
SOURCE="${1:?Использование: $0 <source_dir> <backup_dir> [keep_days]}"
DEST="${2:?Использование: $0 <source_dir> <backup_dir> [keep_days]}"
KEEP="${3:-7}"

[ -d "$SOURCE" ] || { echo "Нет папки: $SOURCE"; exit 1; }
mkdir -p "$DEST"

ARCHIVE="$DEST/backup_$(basename "$SOURCE")_$(date +%Y%m%d_%H%M%S).tar.gz"

echo "Создаю архив: $ARCHIVE"
tar -czf "$ARCHIVE" -C "$(dirname "$SOURCE")" "$(basename "$SOURCE")" \
    && echo "Готово: $(du -sh "$ARCHIVE" | cut -f1)" \
    || { echo "Ошибка создания архива"; exit 1; }

echo "Удаляю архивы старше $KEEP дней..."
find "$DEST" -name "*.tar.gz" -mtime +"$KEEP" -print -delete

echo "Архивов в $DEST: $(ls "$DEST"/*.tar.gz 2>/dev/null | wc -l)"
```

Запуск: `./backup.sh /home/bdv/projects /tmp/backups 7`

---

### `health-check.sh` — проверка здоровья сервисов

```bash
#!/bin/bash
# Список: "имя URL или порт"
services=(
    "nginx:http://localhost:80"
    "api:http://localhost:8080/health"
    "postgres:tcp://localhost:5432"
)

check_http() {
    code=$(curl -s -o /dev/null -w "%{http_code}" --max-time 3 "$1")
    [ "$code" = "200" ] || [ "$code" = "204" ]
}

check_tcp() {
    host=$(echo "$1" | sed 's|tcp://||' | cut -d: -f1)
    port=$(echo "$1" | cut -d: -f3)
    timeout 2 bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null
}

ok=0
fail=0
for entry in "${services[@]}"; do
    name="${entry%%:*}"
    url="${entry#*:}"

    if [[ "$url" == tcp://* ]]; then
        check_tcp "$url" && status="OK" || status="FAIL"
    else
        check_http "$url" && status="OK" || status="FAIL"
    fi

    if [ "$status" = "OK" ]; then
        echo "  OK   $name ($url)"
        ok=$(( ok + 1 ))
    else
        echo "  FAIL $name ($url)"
        fail=$(( fail + 1 ))
    fi
done

echo ""
echo "Итого: $ok OK, $fail FAIL"
[ "$fail" -gt 0 ] && exit 1 || exit 0
```

---

### `find-zombies.sh` — найти зомби-процессы

```bash
#!/bin/bash
zombies=$(ps aux | awk '$8=="Z" {print $2, $11}')

if [ -z "$zombies" ]; then
    echo "Зомби не найдено"
    exit 0
fi

echo "Зомби-процессы:"
echo "$zombies"
echo ""
echo "Их родители:"
ps aux | awk '$8=="Z" {print $2}' | while read -r zpid; do
    ppid=$(grep PPid /proc/"$zpid"/status 2>/dev/null | awk '{print $2}')
    [ -n "$ppid" ] && echo "  PID=$zpid → Родитель PID=$ppid: $(ps -p "$ppid" -o comm= 2>/dev/null)"
done
```

---

## 14. Cron

```bash
crontab -e          # редактировать
crontab -l          # посмотреть
crontab -r          # удалить всё (осторожно!)
```

Формат строки:
```
* * * * * команда
│ │ │ │ │
│ │ │ │ └── день недели (0-7, 0 и 7 = воскресенье)
│ │ │ └──── месяц (1-12)
│ │ └────── день месяца (1-31)
│ └──────── час (0-23)
└────────── минута (0-59)
```

```bash
*/5 * * * *   /scripts/monitor.sh >> /var/log/monitor.log 2>&1   # каждые 5 мин
0 3 * * *     /scripts/backup.sh                                  # каждый день в 3:00
0 9 * * 1     /scripts/report.sh                                  # каждый понедельник в 9:00
0 0 1 * *     /scripts/cleanup.sh                                 # первый день месяца
```

**Важно для cron:**
- Используй **полные пути** к командам: `/usr/bin/curl`, не `curl`
- Пиши `>> /path/to/log 2>&1` — иначе не увидишь ошибок
- `PATH` в cron минимальный — в начале crontab добавь: `PATH=/usr/local/bin:/usr/bin:/bin`
- Проверяй работу: `systemctl status cron` или `grep CRON /var/log/syslog`

---

## 15. Шпаргалка

```bash
# Файлы
ls -lah              cp -r  mv  rm -rf  ln -s  touch  mkdir -p

# Права
chmod 755 f          chmod -R 644 dir/    chown user:group f

# Процессы
ps aux | grep X      pgrep X   pidof X   kill -9 PID   killall X

# Сигналы
kill -TERM   kill -KILL   kill -HUP   kill -l

# Диск
df -h   df -i   du -sh dir/   du --max-depth=1

# Память
free -h   vmstat 1   cat /proc/meminfo

# Сеть
ip addr   ss -tulpen   curl -I url   ping -c4 host

# Логи
journalctl -u service -f
journalctl -b -p err
tail -f /var/log/syslog | grep ERROR

# Архивы
tar -czf out.tar.gz dir/    tar -xzf in.tar.gz
tar -tzf arch.tar.gz        # просмотр без распаковки

# systemd
systemctl {start|stop|restart|status|enable} service
journalctl -u service -f

# Поиск
find . -name "*.log" -mtime +7
find . -size +100M -type f
lsof -i :8080                # кто держит порт
lsof +L1                     # удалённые открытые файлы
```

| Утилита | Зачем |
|---------|--------|
| `strace -p PID` | что делает процесс (системные вызовы) |
| `lsof -p PID` | какие файлы/сокеты открыл процесс |
| `ss -tulpen` | все порты с процессами |
| `vmstat 1` | CPU + память + swap в реальном времени |
| `iostat -xz 1` | нагрузка на диск |
| `dmesg -T \| tail` | сообщения ядра с временем |
| `who` / `w` | кто сейчас залогинен |
| `last` | история логинов |
