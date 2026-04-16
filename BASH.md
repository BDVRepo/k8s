# Bash — базовый справочник для работы в терминале

> Цель: быстро и уверенно работать в Linux shell без лишней теории.
> Формат: только важное + короткие примеры.

---

## 1. Что такое Bash

- Bash — это командная оболочка (shell): принимает команды и выполняет их.
- Терминал — окно, где вы вводите команды Bash.
- Рабочая директория (current directory) — папка, в которой вы сейчас находитесь.

Проверить:
```bash
echo $SHELL
pwd
```

---

## 2. Навигация и файлы

```bash
pwd                 # где я сейчас
ls                  # список файлов
ls -la              # включая скрытые файлы и права
cd /path/to/dir     # перейти в папку
cd ..               # на уровень выше
cd ~                # в домашнюю папку
mkdir mydir         # создать папку
touch file.txt      # создать пустой файл
cp a.txt b.txt      # копия файла
mv old.txt new.txt  # переместить / переименовать
rm file.txt         # удалить файл
rm -r mydir         # удалить папку рекурсивно
```

---

## 3. Просмотр и поиск

```bash
cat file.txt                    # показать весь файл
less file.txt                   # просмотр с прокруткой
head -n 20 file.txt             # первые 20 строк
tail -n 20 file.txt             # последние 20 строк
tail -f app.log                 # смотреть лог в реальном времени
grep "error" app.log            # поиск текста
grep -R "TODO" .                # рекурсивный поиск по папке
find . -name "*.md"             # поиск файлов по маске
```

---

## 4. Перенаправления и пайпы

- `>` перезаписывает файл.
- `>>` дописывает в конец.
- `|` передает вывод одной команды во вход следующей.

```bash
echo "hello" > out.txt
echo "world" >> out.txt
cat out.txt | grep "hello"
ps aux | grep nginx
kubectl get pods -n demo | grep Running
```

---

## 5. Переменные и окружение

```bash
name="Alex"
echo $name

export APP_ENV=dev      # переменная окружения для дочерних процессов
echo $APP_ENV
env | grep APP_ENV
```

Постоянно добавить переменную:
```bash
echo 'export APP_ENV=dev' >> ~/.bashrc
source ~/.bashrc
```

---

## 6. Права доступа и владельцы (chmod/chown/chgrp)

Как читать строку прав:
```bash
-rwxr-x--- 1 alice devops 1234 Apr 16 12:00 deploy.sh
```
- `-` — файл (`d` = директория).
- `rwx` / `r-x` / `---` — права для owner / group / others.
- `alice` — владелец, `devops` — группа файла.

Для файла:
- `r` читать, `w` изменять, `x` запускать.

Для директории:
- `r` видеть список имен,
- `w` создавать/удалять внутри,
- `x` входить (`cd`) и обращаться к файлам по имени.

Базовые команды:
```bash
ls -l file.sh
chmod +x file.sh                # добавить execute
chmod 644 file.txt              # rw-r--r--
chmod 755 script.sh             # rwxr-xr-x
chmod u+x script.sh             # x только владельцу
chmod o-rwx secret.txt          # убрать все права у others
chown user:user file.txt        # сменить владельца и группу
chown :devops file.txt          # сменить только группу
chown -R user:devops app/       # рекурсивно для папки
chgrp devops file.txt           # сменить только группу
```

Шпаргалка по цифрам:
- `r=4`, `w=2`, `x=1`
- `7=rwx`, `6=rw-`, `5=r-x`, `4=r--`, `0=---`
- `644` — обычный файл, `755` — скрипт/директория, `600` — приватный файл.

Группы и пользователи:
```bash
id                              # uid, gid, список групп
groups                          # группы текущего пользователя
sudo groupadd devops            # создать группу
sudo usermod -aG devops alice   # добавить пользователя в группу
newgrp devops                   # применить в текущей сессии
```

Важно:
- `usermod -aG` использовать именно с `-a`, иначе можно потерять старые группы.
- `chown` обычно требует `sudo`, если файл не ваш.

---

## 7. Процессы и управление задачами

```bash
ps aux | grep python
top
htop
kill <pid>
kill -9 <pid>             # жестко завершить (осторожно)
```

Job control:
```bash
sleep 1000                # запустили процесс
Ctrl+Z                    # поставить на паузу
bg                        # продолжить в фоне
jobs                      # список фоновых jobs
fg %1                     # вернуть job 1 на передний план
```

Запуск в фоне:
```bash
command &
nohup command > app.log 2>&1 &
```

---

## 8. История и автодополнение

```bash
history                 # история команд
!123                    # выполнить команду из history по номеру
!!                      # повторить предыдущую команду
Ctrl+R                  # поиск по истории
Tab                      # автодополнение
```

---

## 9. Горячие клавиши Bash (очень полезно)

- `Ctrl + C` — прервать текущую команду.
- `Ctrl + D` — выйти из shell (EOF).
- `Ctrl + Z` — приостановить процесс.
- `Ctrl + L` — очистить экран.
- `Ctrl + A` — в начало строки.
- `Ctrl + E` — в конец строки.
- `Ctrl + U` — удалить от курсора до начала строки.
- `Ctrl + K` — удалить от курсора до конца строки.
- `Ctrl + W` — удалить предыдущее слово.

---

## 10. Как закрывать/сохранять (самое частое)

### `nano`
- Сохранить: `Ctrl + O`, затем `Enter`.
- Выйти: `Ctrl + X`.
- Если спросит сохранить: `Y` (да) / `N` (нет).

### `vim`
- Войти в режим команд: `Esc`.
- Сохранить: `:w` + Enter.
- Выйти: `:q` + Enter.
- Сохранить и выйти: `:wq` + Enter.
- Выйти без сохранения: `:q!` + Enter.

### `less`
- Выйти: `q`.
- Прокрутка: `Space`/`PgDn` вниз, `b`/`PgUp` вверх.
- Поиск: `/text`, затем `n` для следующего совпадения.

### `cat > file`
```bash
cat > notes.txt
hello
world
```
- Завершить ввод и сохранить: `Ctrl + D`.
- Прервать без завершения команды: `Ctrl + C`.

### `tail -f logs`
- Остановить просмотр: `Ctrl + C`.

### `tmux` (если используете)
- Отсоединиться (не убивая сессию): `Ctrl + B`, затем `D`.
- Закрыть окно: `exit` внутри shell.
- Список сессий: `tmux ls`.
- Подключиться: `tmux attach -t <session>`.

---

## 11. Базовые конструкции bash-скриптов

```bash
#!/usr/bin/env bash
set -euo pipefail

name="${1:-world}"

if [[ "$name" == "admin" ]]; then
  echo "Hello, admin"
else
  echo "Hello, $name"
fi

for i in 1 2 3; do
  echo "step $i"
done
```

Сделать исполняемым и запустить:
```bash
chmod +x script.sh
./script.sh Alex
```

---

## 12. Полезный минимум для ежедневной работы DevOps

```bash
# файлы/поиск
ls -la
grep -R "text" .
find . -name "*.yaml"

# сеть
curl -I https://example.com
ss -tulpen

# процессы
ps aux
top
kill <pid>

# архивы
tar -czf archive.tar.gz folder/
tar -xzf archive.tar.gz

# права
chmod +x script.sh
chown -R user:user dir/
```

---

## 13. Типичные ошибки новичков

- Использовать `rm -rf` без `pwd` и проверки пути.
- Запускать команды с `sudo`, не понимая что они делают.
- Путать `>` (перезапись) и `>>` (дозапись).
- Забывать про кавычки в путях с пробелами.
- Править файл в `vim`, не зная как выйти (`:q!` спасает).

---

## 14. Быстрая памятка

- Прервать команду: `Ctrl + C`
- Выйти из shell: `Ctrl + D` или `exit`
- Nano: `Ctrl+O` (save), `Ctrl+X` (exit)
- Vim: `Esc :wq` (save+exit), `Esc :q!` (exit without save)
- Проверить где вы: `pwd`
- Перед удалением: `ls` и еще раз `pwd`
