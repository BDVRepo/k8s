# Bash — справочник с задачами

> Объяснения простые. Для каждой темы: что это, пример, задача, решение.

Связка с ОС: [LINUX.md](LINUX.md)

---

## Содержание

1. [Переменные](#1-переменные)
2. [Математика](#2-математика)
3. [Условия if/elif/else](#3-условия-ifelifelse)
4. [Проверки файлов](#4-проверки-файлов)
5. [Сравнение чисел и строк](#5-сравнение-чисел-и-строк)
6. [Цикл for](#6-цикл-for)
7. [Цикл while и until](#7-цикл-while-и-until)
8. [case](#8-case)
9. [Функции](#9-функции)
10. [Аргументы скрипта](#10-аргументы-скрипта)
11. [Массивы](#11-массивы)
12. [Перенаправления и пайпы](#12-перенаправления-и-пайпы)
13. [sed](#13-sed)
14. [awk](#14-awk)
15. [grep](#15-grep)
16. [sort и uniq](#16-sort-и-uniq)
17. [find](#17-find)
18. [Коннекторы && || ;](#18-коннекторы----)
19. [Самописные утилиты](#19-самописные-утилиты)
20. [Шпаргалка](#20-шпаргалка)

---

## 1. Переменные

```bash
name="Alex"         # присвоить (без пробелов вокруг =!)
echo "$name"        # использовать (всегда в кавычках)
echo "${name}!"     # если рядом другие символы

# Переменная из команды
today=$(date +%F)
echo "$today"

# Значение по умолчанию: если $1 не передан — берём "guest"
user="${1:-guest}"
```

**Задача:** напиши скрипт, который выводит `Привет, <имя>!` где имя берётся из `$1`, а если не передано — выводит `Привет, Аноним!`

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
name="${1:-Аноним}"
echo "Привет, $name!"
```
</details>

---

## 2. Математика

```bash
result=$(( 5 + 3 ))     # 8
echo $(( 10 / 3 ))      # 3 (целочисленное)
echo $(( 10 % 3 ))      # 1 (остаток)
x=5
(( x++ ))               # x теперь 6
echo $x
```

**Задача:** напиши скрипт, который принимает два числа `$1` и `$2` и выводит их сумму, разность, произведение.

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
a=$1
b=$2
echo "Сумма:      $(( a + b ))"
echo "Разность:   $(( a - b ))"
echo "Произведение: $(( a * b ))"
```
</details>

---

## 3. Условия if/elif/else

```bash
#!/bin/bash
x=10

if [ $x -gt 5 ]; then
    echo "больше 5"
elif [ $x -eq 5 ]; then
    echo "равно 5"
else
    echo "меньше 5"
fi
```

Скобки `[ ]` — это команда `test`. Внутри обязательны пробелы.

**Задача:** скрипт принимает число `$1`. Если > 100 — «большое», если > 10 — «среднее», иначе — «маленькое».

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
n=$1
if [ $n -gt 100 ]; then
    echo "большое"
elif [ $n -gt 10 ]; then
    echo "среднее"
else
    echo "маленькое"
fi
```
</details>

---

## 4. Проверки файлов

| Флаг | Смысл |
|------|--------|
| `-f file` | файл существует и является обычным файлом |
| `-d file` | существует и это директория |
| `-e file` | существует (любой тип) |
| `-r file` | можно читать |
| `-w file` | можно писать |
| `-x file` | можно исполнять |
| `-s file` | существует и не пустой |
| `-z "$var"` | строка пустая |
| `-n "$var"` | строка НЕ пустая |

```bash
#!/bin/bash
file="/etc/passwd"

if [ -f "$file" ]; then
    echo "Файл существует"
fi

if [ ! -d "/tmp/mydir" ]; then
    echo "Папки нет, создаю"
    mkdir /tmp/mydir
fi
```

**Задача:** скрипт принимает путь `$1`. Проверить: если это файл — вывести его размер (`du -sh`), если директория — вывести количество файлов внутри, если не существует — вывести ошибку.

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
path=$1
if [ -f "$path" ]; then
    echo "Файл, размер: $(du -sh "$path" | cut -f1)"
elif [ -d "$path" ]; then
    count=$(find "$path" -type f | wc -l)
    echo "Директория, файлов: $count"
else
    echo "Не существует: $path"
    exit 1
fi
```
</details>

---

## 5. Сравнение чисел и строк

**Числа** (только целые):

| Оператор | Смысл |
|----------|--------|
| `-eq` | равно |
| `-ne` | не равно |
| `-gt` | больше |
| `-ge` | больше или равно |
| `-lt` | меньше |
| `-le` | меньше или равно |

**Строки:**

```bash
[ "$a" = "$b" ]    # равны
[ "$a" != "$b" ]   # не равны
[ -z "$a" ]        # строка пустая
[ -n "$a" ]        # строка не пустая
```

**Задача:** скрипт принимает `$1`. Если это слово «admin» — вывести «Доступ разрешён», иначе «Доступ запрещён».

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
user="$1"
if [ "$user" = "admin" ]; then
    echo "Доступ разрешён"
else
    echo "Доступ запрещён"
fi
```
</details>

---

## 6. Цикл for

```bash
# По списку значений
for fruit in apple orange banana; do
    echo "Фрукт: $fruit"
done

# По числам (seq)
for i in $(seq 1 5); do
    echo "Шаг $i"
done

# По файлам
for file in *.sh; do
    echo "Скрипт: $file"
done
```

**Задача:** напиши скрипт, который получает список имён через `$@` и для каждого выводит `Привет, <имя>!`

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
for name in "$@"; do
    echo "Привет, $name!"
done
```

Запуск: `./script.sh Alice Bob Charlie`
</details>

---

## 7. Цикл while и until

**while** — выполнять ПОКА условие истинно:

```bash
count=1
while [ $count -le 5 ]; do
    echo "Шаг $count"
    count=$(( count + 1 ))
done
```

**until** — выполнять ПОКА условие ЛОЖНО (противоположность while):

```bash
count=1
until [ $count -gt 5 ]; do
    echo "Шаг $count"
    count=$(( count + 1 ))
done
```

**Читать файл построчно:**

```bash
while IFS= read -r line; do
    echo "Строка: $line"
done < /etc/passwd
```

**Задача:** напиши скрипт, который каждые 2 секунды проверяет существует ли файл `/tmp/ready`. Как только файл появился — вывести «Файл найден!» и выйти. Максимум 10 попыток.

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
attempts=0
until [ -f /tmp/ready ]; do
    attempts=$(( attempts + 1 ))
    echo "Попытка $attempts — файла нет, жду..."
    if [ $attempts -ge 10 ]; then
        echo "Timeout!"
        exit 1
    fi
    sleep 2
done
echo "Файл найден!"
```
</details>

---

## 8. case

Когда нужно проверить одно значение на много вариантов — красивее чем цепочка `elif`.

```bash
#!/bin/bash
fruit=$1

case $fruit in
    apple)
        echo "Яблоко" ;;
    orange|lemon)
        echo "Цитрус" ;;
    *)
        echo "Неизвестный фрукт" ;;
esac
```

- `*)` — «всё остальное» (как `else`)
- `|` — ИЛИ (несколько значений в одной ветке)
- `;;` — конец ветки (обязательно)

**Задача:** скрипт принимает `$1` — расширение файла (txt, csv, json, xml). Для каждого — вывести тип: «текст», «таблица», «данные», «разметка». Для неизвестного — «неизвестный тип».

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
ext=$1
case $ext in
    txt|md)    echo "текст" ;;
    csv|xlsx)  echo "таблица" ;;
    json|yaml) echo "данные" ;;
    xml|html)  echo "разметка" ;;
    *)         echo "неизвестный тип" ;;
esac
```
</details>

---

## 9. Функции

```bash
#!/bin/bash

# Объявление
greet() {
    echo "Привет, $1!"
}

# Вызов
greet "Alice"
greet "Bob"

# Функция возвращает значение через echo
add() {
    echo $(( $1 + $2 ))
}

result=$(add 3 4)
echo "Сумма: $result"
```

Переменные внутри функции — глобальные по умолчанию. Чтобы сделать локальными: `local x=5`

**Задача:** напиши функцию `is_even`, которая принимает число и выводит «чётное» или «нечётное». Вызови её для чисел 2, 7, 10.

<details>
<summary>Решение</summary>

```bash
#!/bin/bash

is_even() {
    if [ $(( $1 % 2 )) -eq 0 ]; then
        echo "$1 — чётное"
    else
        echo "$1 — нечётное"
    fi
}

is_even 2
is_even 7
is_even 10
```
</details>

---

## 10. Аргументы скрипта

| Переменная | Смысл |
|------------|--------|
| `$0` | имя скрипта |
| `$1`, `$2` … | аргументы по позиции |
| `$#` | количество аргументов |
| `$@` | все аргументы (каждый отдельно) |
| `$*` | все аргументы (одна строка) |
| `$?` | код выхода последней команды |
| `$$` | PID текущего shell |

```bash
#!/bin/bash
echo "Скрипт: $0"
echo "Аргументов: $#"
echo "Первый: $1"

for arg in "$@"; do
    echo "  - $arg"
done
```

**Задача:** напиши скрипт, который проверяет что передан ровно 2 аргумента. Если нет — выводит `Использование: $0 <arg1> <arg2>` и выходит с кодом 1.

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
if [ $# -ne 2 ]; then
    echo "Использование: $0 <arg1> <arg2>"
    exit 1
fi
echo "Arg1: $1"
echo "Arg2: $2"
```
</details>

---

## 11. Массивы

```bash
fruits=("apple" "orange" "banana")

echo "${fruits[0]}"          # первый элемент
echo "${fruits[@]}"          # все элементы
echo "${#fruits[@]}"         # количество элементов

fruits+=("grape")            # добавить элемент

for fruit in "${fruits[@]}"; do
    echo "$fruit"
done
```

**Задача:** создай массив из имён пользователей. Для каждого выведи — есть ли такой пользователь в системе (`grep "$user" /etc/passwd`).

<details>
<summary>Решение</summary>

```bash
#!/bin/bash
users=("root" "bdv" "nobody" "ghost")

for user in "${users[@]}"; do
    if grep -q "^$user:" /etc/passwd; then
        echo "ЕСТЬ:    $user"
    else
        echo "НЕТ:     $user"
    fi
done
```
</details>

---

## 12. Перенаправления и пайпы

```bash
cmd > file.txt       # stdout в файл (перезапись)
cmd >> file.txt      # stdout в файл (дозапись)
cmd 2> err.txt       # stderr в файл
cmd > file.txt 2>&1  # и stdout и stderr в файл
cmd > /dev/null      # выбросить stdout (тихий режим)

cmd1 | cmd2          # stdout cmd1 → stdin cmd2
```

Важно: порядок имеет значение!
```bash
cmd > file 2>&1   # ОБА в file (правильно)
cmd 2>&1 > file   # stderr в терминал, stdout в file (не то!)
```

**Задача:** запусти `ls /tmp /nonexistent` и сохрани: обычный вывод в `out.txt`, ошибки в `err.txt`.

<details>
<summary>Решение</summary>

```bash
ls /tmp /nonexistent > out.txt 2> err.txt
cat out.txt
cat err.txt
```
</details>

---

## 13. sed

Потоковый редактор. Меняет текст в файле или потоке.

```bash
# Заменить первое вхождение в каждой строке
echo "hello hello" | sed 's/hello/hi/'        # hi hello

# Заменить все вхождения (флаг g)
echo "hello hello" | sed 's/hello/hi/g'       # hi hi

# Правка файла на месте
sed -i 's/old/new/g' config.txt

# Удалить строки содержащие слово
sed '/DEBUG/d' app.log

# Вывести только строки 5-10
sed -n '5,10p' file.txt
```

**Задача:** есть файл `config.txt` с содержимым `host=localhost`. Замени `localhost` на `192.168.1.1` и выведи результат (без изменения файла).

<details>
<summary>Решение</summary>

```bash
sed 's/localhost/192.168.1.1/' config.txt
```
</details>

---

## 14. awk

Обрабатывает текст по колонкам. `$1`, `$2` — первая, вторая колонка.

```bash
echo "alice 30 dev" | awk '{print $1}'              # alice
echo "alice 30 dev" | awk '{print $1, $3}'          # alice dev

# Разделитель : (для /etc/passwd)
awk -F: '{print $1}' /etc/passwd                    # имена пользователей

# Фильтр: вывести строки где 3-е поле > 1000
awk -F: '$3 > 1000 {print $1, $3}' /etc/passwd

# Сумма значений в колонке
echo -e "10\n20\n30" | awk '{sum += $1} END {print sum}'   # 60
```

**Задача:** из вывода `ps aux` вывести только имя пользователя и имя процесса (первая и одиннадцатая колонки), только первые 5 строк.

<details>
<summary>Решение</summary>

```bash
ps aux | awk 'NR>1 {print $1, $11}' | head -5
```
</details>

---

## 15. grep

```bash
grep "error" file.txt          # строки содержащие "error"
grep -i "error" file.txt       # без учёта регистра
grep -v "debug" file.txt       # строки БЕЗ "debug"
grep -n "error" file.txt       # с номерами строк
grep -c "error" file.txt       # посчитать совпадения
grep -r "TODO" .               # рекурсивно по папке
grep -l "error" *.log          # только имена файлов
grep -E "error|warn" file.txt  # регулярка: error ИЛИ warn
```

**Задача:** найти все строки в `/etc/passwd` где оболочка `/bin/bash`, вывести только имена пользователей.

<details>
<summary>Решение</summary>

```bash
grep "/bin/bash" /etc/passwd | awk -F: '{print $1}'
```
</details>

---

## 16. sort и uniq

```bash
sort file.txt              # алфавитно
sort -r file.txt           # в обратном порядке
sort -n numbers.txt        # числовая сортировка (10 > 2, а не "1" > "2")
sort -rn numbers.txt       # числовая обратная
sort -u file.txt           # убрать дубли

uniq file.txt              # убрать соседние дубли (нужен sort перед ним)
uniq -c file.txt           # посчитать повторения
sort file.txt | uniq -c | sort -rn   # частота слов
```

**Задача:** есть список слов (каждое на новой строке). Вывести уникальные слова в алфавитном порядке с подсчётом сколько раз каждое встречалось.

<details>
<summary>Решение</summary>

```bash
sort words.txt | uniq -c | sort -rn
```
</details>

---

## 17. find

```bash
find . -name "*.log"                  # по маске имени
find . -type f                        # только файлы
find . -type d                        # только директории
find . -mtime +7                      # изменялся > 7 дней назад
find . -size +100M                    # больше 100 MB
find . -name "*.sh" -exec chmod +x {} \;   # выполнить команду для каждого
find . -name "*.log" -delete          # найти и удалить
find . -name "*.log" -print0 | xargs -0 grep "ERROR"   # безопасный поиск с пробелами
```

**Задача:** найти все `.sh` файлы в `/home` которые не имеют права на исполнение (`! -perm -u+x`) и вывести их имена.

<details>
<summary>Решение</summary>

```bash
find /home -name "*.sh" ! -perm -u+x
```
</details>

---

## 18. Коннекторы `&&` `||` `;`

```bash
mkdir newdir && cd newdir          # второе — только если первое OK
cd /bad || echo "Нет папки"        # второе — только если первое УПАЛО
echo "Раз" ; echo "Два"           # просто по очереди (игнорирует ошибки)
```

**Задача:** напиши однострочник: создать папку `deploy`, войти в неё и вывести текущий путь — всё через `&&`.

<details>
<summary>Решение</summary>

```bash
mkdir deploy && cd deploy && pwd
```
</details>

---

## 19. Самописные утилиты

Все файлы создавай в `purpleschool/tools/`.

### `colorlog` — вывод с цветами

```bash
#!/bin/bash
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${BLUE}[INFO]${NC}  $1"; }
log_ok()    { echo -e "${GREEN}[OK]${NC}    $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }

log_info "Запускаем"
log_ok "Подключились"
log_warn "Диск 80%"
log_error "Сервис упал"
```

---

### `sysinfo` — сводка о системе

```bash
#!/bin/bash
echo "===== Система ====="
echo "Hostname:  $(hostname)"
echo "OS:        $(uname -sr)"
echo "Uptime:    $(uptime -p)"
echo ""
echo "===== CPU ====="
echo "Ядер:      $(nproc)"
echo "Загрузка:  $(uptime | awk -F'load average:' '{print $2}')"
echo ""
echo "===== Память ====="
free -h | awk 'NR==2 {printf "Всего: %s  Занято: %s  Свободно: %s\n", $2, $3, $4}'
echo ""
echo "===== Диск ====="
df -h | grep -v tmpfs
```

---

### `cleanup` — удалить файлы старше N дней

```bash
#!/bin/bash
DIR="${1:?Использование: $0 <dir> <days>}"
DAYS="${2:?Использование: $0 <dir> <days>}"

[ -d "$DIR" ] || { echo "Нет папки: $DIR"; exit 1; }

echo "Удаляю в $DIR файлы старше $DAYS дней..."
find "$DIR" -type f -mtime +"$DAYS" -print -delete
echo "Готово"
```

Запуск: `./cleanup /tmp/logs 7`

---

### `envcheck` — проверить нужные переменные окружения

```bash
#!/bin/bash
required=("DB_HOST" "DB_PORT" "APP_SECRET")
missing=0

for var in "${required[@]}"; do
    if [ -z "${!var}" ]; then
        echo "НЕТ: $var"
        missing=$(( missing + 1 ))
    else
        echo "OK:  $var = ${!var}"
    fi
done

[ "$missing" -gt 0 ] && { echo "Задай $missing переменных"; exit 1; }
echo "Всё OK"
```

Запуск: `DB_HOST=localhost DB_PORT=5432 APP_SECRET=abc ./envcheck`

---

### `repeat` — повторить команду N раз

```bash
#!/bin/bash
n="${1:?Использование: $0 <N> <команда>}"
shift

for i in $(seq 1 "$n"); do
    echo "--- Попытка $i/$n ---"
    "$@"
    sleep 1
done
```

Запуск: `./repeat 3 curl -s http://localhost:8080/health`

---

### `portcheck` — проверить список портов

```bash
#!/bin/bash
hosts=("127.0.0.1:22" "127.0.0.1:80" "127.0.0.1:5432")

for entry in "${hosts[@]}"; do
    host="${entry%%:*}"
    port="${entry##*:}"
    if timeout 1 bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null; then
        echo "OPEN   $host:$port"
    else
        echo "CLOSED $host:$port"
    fi
done
```

---

### `log_errors` — найти ошибки в логе

```bash
#!/bin/bash
LOG="${1:?Использование: $0 <log_file>}"
PATTERN="${2:-ERROR}"

[ -f "$LOG" ] || { echo "Нет файла: $LOG"; exit 1; }

count=$(grep -c "$PATTERN" "$LOG" || true)
echo "Найдено '$PATTERN': $count"
echo ""
grep -n "$PATTERN" "$LOG" | tail -20
```

Запуск: `./log_errors /var/log/syslog error`

---

### `disk_alert` — предупреждение если диск > 80%

```bash
#!/bin/bash
THRESHOLD="${1:-80}"

df -h | tail -n +2 | while read -r line; do
    usage=$(echo "$line" | awk '{print $5}' | tr -d '%')
    mount=$(echo "$line" | awk '{print $6}')
    if echo "$usage" | grep -qE '^[0-9]+$' && [ "$usage" -gt "$THRESHOLD" ]; then
        echo "ВНИМАНИЕ: $mount занят на $usage%"
    fi
done
```

---

## 20. Шпаргалка

```bash
# Переменная
name="Alex"
echo "$name"
today=$(date +%F)
default="${1:-значение_по_умолчанию}"

# Математика
$(( a + b ))   $(( a * b ))   $(( a % b ))

# Условие
if [ "$a" = "$b" ]; then ... fi
if [ $n -gt 10 ]; then ... fi
if [ -f "$file" ]; then ... fi
if [ -z "$var" ]; then ... fi      # пустая строка

# Циклы
for i in $(seq 1 5); do ... done
for item in "${arr[@]}"; do ... done
while [ $n -lt 10 ]; do ... done
until [ -f /tmp/ready ]; do ... done

# case
case $var in
    val1) ... ;;
    val2|val3) ... ;;
    *) ... ;;
esac

# Функция
myfunc() { echo "$1"; }
result=$(myfunc "arg")

# Массив
arr=("a" "b")
arr+=("c")
echo "${arr[0]}"
echo "${#arr[@]}"

# Аргументы
$0 $1 $2   $# $@ $?

# Перенаправления
cmd > file        # stdout в файл
cmd >> file       # дозапись
cmd 2> err        # stderr в файл
cmd > file 2>&1   # всё в файл
cmd > /dev/null   # выбросить

# set — безопасный режим
set -euo pipefail
```

| Спецпеременная | Смысл |
|----------------|--------|
| `$?` | код выхода последней команды |
| `$$` | PID текущего shell |
| `$!` | PID последнего фонового процесса |
| `$0` | имя скрипта |
| `$#` | количество аргументов |
| `$@` | все аргументы |

---

## Частые ошибки

| Ошибка | Почему плохо | Как правильно |
|--------|-------------|---------------|
| `name = "Alex"` | пробелы = вызов команды | `name="Alex"` |
| `[ $var = "x" ]` | сломается если `$var` пуст | `[ "$var" = "x" ]` |
| `dir = pwd` | присвоение не работает так | `dir=$(pwd)` |
| `for f in $(ls)` | ломается на пробелах в именах | `for f in *` или `find -print0` |
| Скрипт без `#!/bin/bash` | непредсказуемый shell | добавь shebang |
| Скрипт без `chmod +x` | `Permission denied` | `chmod +x script.sh` |
