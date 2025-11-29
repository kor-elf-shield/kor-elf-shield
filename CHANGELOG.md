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