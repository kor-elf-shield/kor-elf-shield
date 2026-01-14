## 0.5.0 (soon)
***
#### Русский
* В настройках analyzer.toml добавил параметры local_enable и local_notify.
  * local_enable = Включает отслеживание локальных авторизаций (TTY, физический доступ). По умолчанию включён.
  * local_notify = Включает уведомления о локальных авторизациях. По умолчанию включён.
***
#### English
* Added local_enable and local_notify parameters to analyzer.toml settings.
  * local_enable = Enables tracking of local logins (TTY, physical access). Enabled by default.
  * local_notify = Enables notifications about local logins. Enabled by default.
***
## 0.4.0 (11.1.2026)
***
#### Русский
* Удалён параметр options.docker_support из файла firewall.toml. Настройки от Docker перенесены в файл docker.toml.
* В настройках docker.toml добавил возможность переключать режим работы с Docker через параметр rule_strategy.
  * incremental = добавляются или удаляются только правила конкретного контейнера (сейчас по умолчанию)
  * rebuild = при любом изменении все цепочки Docker пересоздаются целиком (старый режим)
* Исправлена ошибка:
  * Настройка binaryLocations.docker не работала.
  * Программа аварийно завершалась после остановки Docker'а.
  * Указанные в настройках IP-адреса не блокировались во время перенаправления в контейнер Docker.
***
#### English
* Removed the options.docker_support parameter from firewall.toml. Docker settings have been moved to the docker.toml file.
* Added the ability to switch Docker operation mode via the rule_strategy parameter to the docker.toml settings.
  * incremental = only rules for a specific container are added or removed (currently the default)
  * rebuild = any change rebuilds all Docker chains (old mode)
* Fixed error:
  * The binaryLocations.docker setting did not work.
  * The program crashed after Docker was stopped.
  * The IP addresses specified in the settings were not blocked during redirection to the Docker container.
***
## 0.3.0 (4.1.2026)
***
#### Русский
* Добавлена частичная поддержка Docker.
  * Добавлен параметр options.docker_support в firewall.toml. Это включает поддержку Docker.
  * Каждый запуск контейнера будет полностью пересчитываться правила у chain, которые относятся к Docker. (в будущем планирую это переработать)
* Добавлены настройки для уведомлений по электронной почте.
  * Добавлен файл настроек notifications.toml.
* Реализовано уведомление о входах по SSH.
  * Добавлен файл настроек analyzer.toml.
* Служба systemd
  * Изменено WantedBy с sysinit.target на multi-user.target
  * Убрано ExecStop. По факту это не работало. Чтобы остановить сервис с очисткой правил nftables выпоните команду: kor-elf-shield stop
  * Добавлено Restart=on-failure. Нужно для того, чтобы программа перезапустилась после критической ошибки.
***
#### English
* Added partial Docker support.
  * Added the options.docker_support parameter to firewall.toml. This enables Docker support.
  * Each container launch will completely recalculate the Docker-specific rules in chain. (I plan to rework this in the future)
* Added settings for email notifications.
  * Added notifications.toml settings file.
* Implemented notification of SSH logins.
  * Added analyzer.toml settings file.
* Systemd service
  * Changed WantedBy from sysinit.target to multi-user.target
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