## 0.3.0 (soon)
***
#### Русский
* Добавлены настройки для уведомлений по электронной почте.
  * Добавлен файл настроек notifications.toml.
* Реализовано уведомление о входах по SSH.
  * Добавлен файл настроек analyzer.toml.
* Служба systemd
  * Изменено WantedBy с multi-user.target на multi-user.target
  * Убрано ExecStop. По факту это не работало. Чтобы остановить сервис с очисткой правил nftables выпоните команду: kor-elf-shield stop
  * Добавлено Restart=on-failure. Нужно для того, чтобы программа перезапустилась после критической ошибки.
***
#### English
* Added settings for email notifications.
  * Added notifications.toml settings file.
* Implemented notification of SSH logins.
  * Added analyzer.toml settings file.
* Systemd service
  * Changed WantedBy from multi-user.target to multi-user.target
  * Removed ExecStop. It didn't actually work. To stop the service and clear the nftables rules, run the command: kor-elf-shield stop
  * Added Restart=on-failure. This is necessary to ensure the program restarts after a critical error.
## 0.2.0 (29.11.2025)
***
#### Русский
* Добавлен параметр clear_mode в firewall.toml. Он позволяет переключать режим очистки всех правил в nftables или только таблицу относящие к программе.
* Добавлен параметр input_priority в firewall.toml. Можно указать приоритет от -50 по 50 к chain input.
* Добавлен параметр output_priority в firewall.toml. Можно указать приоритет от -50 по 50 к chain output.
* Добавлен параметр forward_priority в firewall.toml. Можно указать приоритет от -50 по 50 к chain forward.
***
#### English
* Added the clear_mode parameter to firewall.toml. It allows you to toggle clearing of all rules in nftables or only the program-specific table.
* Added the input_priority parameter to firewall.toml. You can specify a priority from -50 to 50 for chain input.
* Added the output_priority parameter to firewall.toml. You can specify a priority from -50 to 50 for chain output.
* Added the forward_priority parameter to firewall.toml. You can specify a priority from -50 to 50 for chain forward.
***
## 0.1.0 (8.11.2025)
***
#### Русский
* Реализована возможность настраивать nftables:
  * По умолчанию разрешить или блокировать входящий трафик.
  * По умолчанию разрешить или блокировать исходящий трафик.
  * Настройка icmp.
  * Настройка портов.
  * Настройка белых и чёрных списков IP адресов.
* Настройка логирование.
***
#### English
* Implemented the ability to configure nftables:
  * Allow or block incoming traffic by default.
  * Allow or block outgoing traffic by default.
  * ICMP configuration.
  * Port configuration.
  * IP address whitelisting and blacklisting.
* Logging configuration.
***