# Kor-Elf Shield
### <strong>English</strong> &nbsp;&nbsp;&nbsp; <a href="README.md">Русский</a>
<p style="color: red; font-weight: bold">ATTENTION: The program is under active development and is NOT CURRENTLY PRODUCTION READY</p>
<p>I have been using ConfigServer Security and Firewall (csf) for over 10 years to protect my server. But unfortunately, in September I found out that the company that supported this great product closed on August 31, 2025. CSF is written in PERL. And the company uploaded all the source codes to its repository under the GPLv3 license. But I do not know the PERL language. And it's hard for me to read it. :)</p>
<p>I decided to implement my solution in the Go Lang language. It will not be a complete copy of CSF. CSF just inspired me to do something similar to protect my server.</p>

***

<p style="color: red; font-weight: bold; font-size: 20px;">Requirements:</p>

* Run as root
* Linux 5.2+
* nftables
* Systemd
* journalctl

***

### Done:
* The ability to configure nftables has been implemented:
    * Allow or block incoming traffic by default.
    * Allow or block outgoing traffic by default.
    * Setting up icmp.
    * Port configuration.
    * Setting up white and black lists of IP addresses.
* Setting up logging.
* Make friends with docker (partially).
* Implement notification settings (for now only by e-mail).
* Send notifications during ssh authorization.
* Password brute-force protection.

### The plans include:
* Notify if a new user appears in the system.
* Notify if system files have changed.
***

## Compiling from source code:
1. git clone https://git.kor-elf.net/kor-elf-shield/kor-elf-shield.git
2. cd kor-elf-shield
3. go build -o kor-elf-shield
4. sudo cp kor-elf-shield /usr/local/bin/kor-elf-shield
5. sudo chmod +x /usr/local/bin/kor-elf-shield
6. sudo cp assets/configs /etc/kor-elf-shield
7. sudo cp assets/kor-elf-shield.service /etc/systemd/system/kor-elf-shield.service
8. sudo cp assets/kor-elf-shield.logrotate /etc/logrotate.d/kor-elf-shield
9. Edit the required parameters in:<br> /etc/kor-elf-shield/kor-elf-shield.toml<br> /etc/kor-elf-shield/firewall.toml
10. sudo systemctl daemon-reload
11. sudo systemctl enable kor-elf-shield <br><strong>Incompatible with ufw and firewalld. These must be disabled, otherwise there will be a conflict with nftables settings. Any programs that work with nftables may conflict with kor-elf-shield.</strong>
12. sudo systemctl start kor-elf-shield

<p>You can download compiled ready-made versions here: <a href="https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/releases">https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/releases</a>.</p>

<p><strong>Attention:</strong> By default, the settings are set to Test mode. After a certain period (5 minutes by default), the program shuts down and all nftables rules are cleared. This is necessary so you can first debug all the rules and then enable them in production mode.</p>

<p><strong>Before launching the program, we recommend that you read the <a href="https://shield.kor-elf.net/docs/0.x/first_launch/language/en" target="_blank">first launch instructions</a>.</strong></p>

***

## Settings:

<p><strong>/etc/kor-elf-shield/kor-elf-shield.toml</strong> - General settings are located here. Information can be found here: <a href="https://shield.kor-elf.net/docs/0.x/kor-elf-shield.toml/language/en" target="_blank">https://shield.kor-elf.net/docs/0.x/kor-elf-shield.toml/language/en</a></p>

<p><strong>/etc/kor-elf-shield/firewall.toml</strong> - Here are the settings related to nftables. Information can be found here: <a href="https://shield.kor-elf.net/docs/0.x/firewall.toml/language/en" target="_blank">https://shield.kor-elf.net/docs/0.x/firewall.toml/language/en</a></p>

***

<p>The software is MIT (see <a href="https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/branch/main/LICENSE">LICENSE</a>) and uses third-party libraries that are distributed on their own terms (see <a href="https://git.kor-elf.net/kor-elf-shield/kor-elf-shield/src/branch/main/LICENSE-3RD-PARTY.txt">LICENSE-3RD-PARTY.txt</a>).</p>