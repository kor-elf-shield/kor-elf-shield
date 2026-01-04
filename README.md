# Kor-Elf Shield
### <strong>Русский</strong> &nbsp;&nbsp;&nbsp; <a href="README_ENG.md">English</a>
<p style="color: red; font-weight: bold">ВНИМАНИЕ: Программа на стадии активной разработки и на данный момент NOT PRODUCTION READY</p>
<p>Я больше 10-ти лет использовал ConfigServer Security and Firewall (csf) для защиты своего сервера. Но, к сожалению, в сентябре узнал, что компания, которая поддерживала этот великолепный продукт, закрылась 31 августа 2025 года. CSF написан на языке PERL. И компания выложила все исходные коды в свой репозиторий под лицензией GPLv3. Но я не знаю язык PERL. И мне его тяжело читать. :)</p>
<p>Я решил реализовать своё решение на языке Go Lang. Это не будет полная копия CSF. Просто CSF вдохновило меня сделать, что-то похожее для защиты своего сервера.</p>

***

<p style="color: red; font-weight: bold; font-size: 20px;">Требования:</p>

* Запуск от имени root
* Linux 5.2+ 
* nftables
* Systemd

***

### Сделанно:
* Реализована возможность настраивать nftables:
    * По умолчанию разрешить или блокировать входящий трафик.
    * По умолчанию разрешить или блокировать исходящий трафик.
    * Настройка icmp.
    * Настройка портов.
    * Настройка белых и чёрных списков IP адресов.
* Настройка логирование.
* Подружить с docker (частично).
* Внедрить настройку уведомлений (пока только e-mail).
* Отправлять уведомления при авторизации ssh.

### В планах:
* Защита от перебора паролей (brute-force).
* Уведомлять, если появится новый пользователь в системе.
* Уведомлять, если изменились системные файлы.
***

## Компиляция из исходного кода:
1. git clone https://git.kor-elf.net/kor-elf-shield/kor-elf-shield.git
2. cd kor-elf-shield
3. go build -o kor-elf-shield
4. sudo cp kor-elf-shield /usr/local/bin/kor-elf-shield
5. sudo chmod +x /usr/local/bin/kor-elf-shield
6. sudo cp -r assets/configs /etc/kor-elf-shield
7. sudo cp assets/kor-elf-shield.service /etc/systemd/system/kor-elf-shield.service
8. sudo cp assets/kor-elf-shield.logrotate /etc/logrotate.d/kor-elf-shield
9. Отредактировать нужные параметры в:<br> /etc/kor-elf-shield/kor-elf-shield.toml<br> /etc/kor-elf-shield/firewall.toml
10. sudo systemctl daemon-reload
11. sudo systemctl enable kor-elf-shield <br><strong>Не совместим с ufw, firewalld. Их надо отключить иначе будет конфликт в настройках nftables. Все программы, которые работают с nftables могут конфликтовать с kor-elf-shield.</strong>
12. sudo systemctl start kor-elf-shield

<p>Скачать скомпилированные готовые версии можно тут: <a href="https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/releases">https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/releases</a>.</p>

<p><strong>Внимание:</strong> по умолчанию в настройках стоит Тестовый режим, после определённого периода (по умолчанию 5 минут) программа выключается и все правила nftables очищаются. Это нужно для того, чтобы вы могли вначале отладить все правила, а после включить в боевом режиме.</p>

<p><strong>Перед запуском программы рекомендуем ознакомиться с <a href="https://shield.kor-elf.net/docs/0.x/first_launch" target="_blank">инструкцией по первому запуску</a>.</strong></p>

***

## Настройки:

<p><strong>/etc/kor-elf-shield/kor-elf-shield.toml</strong> - тут находятся общие настройки. Информацию можно посмотреть тут: <a href="https://shield.kor-elf.net/docs/0.x/kor-elf-shield.toml" target="_blank">https://shield.kor-elf.net/docs/0.x/kor-elf-shield.toml</a></p>

<p><strong>/etc/kor-elf-shield/firewall.toml</strong> - тут находятся настройки, связанные с nftables. Информацию можно посмотреть тут: <a href="https://shield.kor-elf.net/docs/0.x/firewall.toml" target="_blank">https://shield.kor-elf.net/docs/0.x/firewall.toml</a></p>

***

<p>Программное обеспечение является MIT (см. <a href="https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/branch/main/LICENSE">LICENSE</a>) и использует сторонние библиотеки, которые распространяются на их собственных условиях (см. <a href="https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/branch/main/LICENSE-3RD-PARTY.txt">LICENSE-3RD-PARTY.txt</a>).</p>