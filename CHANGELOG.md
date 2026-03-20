## 0.9.0 (Soon)
#### Русский
* Добавилась поддержка Port knocking.
  * В firewall.toml добавился раздел Port knocking.
* Теперь вы можете получить список IP-адресов от различных сервисов для блокировки доступа:
  * Spamhaus Don't Route Or Peer Lists
  * DShield.org Recommended Block List
  * TOR Exit Nodes List
  * Project Honey Pot Directory of Dictionary Attacker IPs
  * C.I. Army Malicious IP List
  * BruteForceBlocker IP List
  * Blocklist.de
  * Stop Forum Spam
  * GreenSnow Hack List
* В настройки файла kor-elf-shield.toml добавлен параметр otherSettingsPath.blocklists.
* Добавлен новый файл настроек, blocklists.toml. Он содержит параметры для получения списка IP-адресов для блокировки.
***
#### English
* Added support for port knocking.
  * A Port knocking section has been added to firewall.toml.
* Now you can get a list of IP addresses from various services to block access:
  * Spamhaus Don't Route Or Peer Lists
  * DShield.org Recommended Block List
  * TOR Exit Nodes List
  * Project Honey Pot Directory of Dictionary Attacker IPs
  * C.I. Army Malicious IP List
  * BruteForceBlocker IP List
  * Blocklist.de
  * Stop Forum Spam
  * GreenSnow Hack List
* Added the otherSettingsPath.blocklists parameter to the kor-elf-shield.toml settings.
* Added a new settings file, blocklists.toml. It contains settings for obtaining a list of IP addresses to block.
## 0.8.0 (9.3.2026)
***
#### Русский
* Теперь можно тонко настроить блокировку портов для IP адреса, который пытается подобрать пароль.
  * В файл настроек analyzer.toml в [[bruteForceProtection.groups]] добавлен новый параметр "block_type". 
  * В файл настроек analyzer.toml в [[bruteForceProtection.groups]] добавлен новый параметр "ports".
  * В файл настроек analyzer.toml в [[bruteForceProtection.groups.rate_limits]] добавлен новый параметр "block_type".
  * В файл настроек analyzer.toml в [[bruteForceProtection.groups.rate_limits]] добавлен новый параметр "ports".
  * Смотрите полный список по ссылке: https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/commit/d0a358a445b1dec850d8b84c06e86bd6872796cf/assets/configs/analyzer.toml
* Команда `kor-elf-shield ban clear` была переименованна в `kor-elf-shield block clear`.
* Добавлена команда `kor-elf-shield block add`. Через эту команду можно заблокировать IP адрес. Смотрите подробно в `kor-elf-shield block add --help`.
* Добавлена команда `kor-elf-shield block delete`. Через эту команду можно удалить заблокированный IP адрес. Смотрите подробно в `kor-elf-shield block delete --help`.
***
#### English
* You can now fine-tune port blocking for the IP address attempting to brute-force a password.
  * A new "block_type" parameter has been added to the analyzer.toml settings file in [[bruteForceProtection.groups]].
  * A new "ports" parameter has been added to the analyzer.toml settings file in [[bruteForceProtection.groups]].
  * A new "block_type" parameter has been added to the analyzer.toml settings file in [[bruteForceProtection.groups.rate_limits]].
  * A new "ports" parameter has been added to the analyzer.toml settings file in [[bruteForceProtection.groups.rate_limits]]. 
  * See the full list at: https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/commit/d0a358a445b1dec850d8b84c06e86bd6872796cf/assets/configs/analyzer.toml
* The `kor-elf-shield ban clear` command has been renamed to `kor-elf-shield block clear`.
* The `kor-elf-shield block add` command has been added. This command can be used to block an IP address. See `kor-elf-shield block add --help` for details.
* The `kor-elf-shield block delete` command has been added. This command can be used to delete a blocked IP address. See `kor-elf-shield block delete --help` for details.
***
## 0.7.0 (28.2.2026)
***
#### Русский
* Добавлена возможность настройки отслеживания событий в журналах.
* Добавлены настройки для защиты от перебора паролей.
* В файл настроек analyzer.toml добавлены новые параметры. Смотрите полный список по ссылке: https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/commit/187c447301b9c0bfa41ec2b2c9435ab0ce44bed6/assets/configs/analyzer.toml
* Добавлена команда `kor-elf-shield ban clear`, которая разблокирует все IP адреса. Которые были забанены.
***
#### English
* Added the ability to customize event tracking in logs.
* Added settings to protect against password guessing.
* New parameters have been added to the analyzer.toml settings file. See the full list at: https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/commit/187c447301b9c0bfa41ec2b2c9435ab0ce44bed6/assets/configs/analyzer.toml
* Added the `kor-elf-shield ban clear` command, which unbans all banned IP addresses.
***
## 0.6.0 (8.2.2026)
***
#### Русский
* Добавлена возможность повторной отправки уведомления, если в прошлый раз произошла ошибка.
* Добавлена команда `kor-elf-shield notifications queue count`, которая возвращает количество уведомлений в очереди в базе данных.
* Добавлена команда `kor-elf-shield notifications queue clear`, которая удаляет все уведомления из очереди в базе данных.
* В файл настроек kor-elf-shield.toml добавлены новые параметры:
  * data_dir = Каталог для постоянных данных приложения (state): локальная база данных, кэш/индексы, файлы состояния и другие служебные файлы. Должен быть доступен на запись пользователю, от имени которого запущен демон. Если каталог не существует — будет создан. По умолчанию: "/var/lib/kor-elf-shield/"
* В файл настроек notifications.toml добавлены новые параметры:
  * enable_retries = Включает повторные попытки отправить уведомление, если сразу не получилось. По умолчанию: true
  * retry_interval = Интервал времени в секундах между попытками. По умолчанию: 600
***
#### English
* Added the ability to retry sending a notification if an error occurred the previous time.
* Added the `kor-elf-shield notifications queue count` command, which returns the number of notifications in the queue in the database.
* Added the `kor-elf-shield notifications queue clear` command, which removes all notifications from the queue in the database.
* New parameters have been added to the kor-elf-shield.toml settings file:
  * data_dir = Directory for persistent application data (state): local database, cache/indexes, state files, and other internal data. Must be writable by the daemon user. If the directory does not exist, it will be created. Default: "/var/lib/kor-elf-shield/" 
* New parameters have been added to the notifications.toml settings file:
  * enable_retries = Enables repeated attempts to send a notification if the first attempt fails. Default: true
  * retry_interval = The time interval in seconds between attempts. Default: 600
***
## 0.5.0 (17.1.2026)
***
#### Русский
* В настройках analyzer.toml добавил параметры local_enable и local_notify.
  * local_enable = Включает отслеживание локальных авторизаций (TTY, физический доступ). По умолчанию включён.
  * local_notify = Включает уведомления о локальных авторизациях. По умолчанию включён.
* В настройках analyzer.toml добавил параметры su_enable и su_notify.
  * su_enable = Включает отслеживание авторизаций через su. По умолчанию включён.
  * su_notify = Включает уведомления об авторизациях через su. По умолчанию включён.
* В настройках analyzer.toml добавил параметры sudo_enable и sudo_notify.
  * sudo_enable = Включает отслеживание авторизаций через sudo. По умолчанию выключен.
  * sudo_notify = Включает уведомления об авторизациях через sudo. По умолчанию включён.
***
#### English
* Added local_enable and local_notify parameters to analyzer.toml settings.
  * local_enable = Enables tracking of local logins (TTY, physical access). Enabled by default.
  * local_notify = Enables notifications about local logins. Enabled by default.
* Added su_enable and su_notify parameters to analyzer.toml settings.
  * su_enable = Enables tracking of logins via su. Enabled by default.
  * su_notify = Enables notifications about logins via su. Enabled by default.
* Added sudo_enable and sudo_notify parameters to analyzer.toml settings.
  * sudo_enable = Enables tracking of logins via sudo. Off by default.
  * sudo_notify = Enables notifications about logins via sudo. Enabled by default.
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