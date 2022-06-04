:us: FLibGoLite :gb:
===
[:ru: *русский* ](README_RU.md)  
[:ukraine: *українською* ](README_UK.md)

### BETA RELEASE v0.2.x * 
_*This software release has not been tested thoroughly yet but based on __[flibgo](https://github.com/vinser/flibgo.git)__ it does the job_

---

__FLibGoLite__ is easy to use home library OPDS server 

>The Open Publication Distribution System (OPDS) catalog format is a syndication format for electronic publications based on Atom and HTTP. OPDS catalogs enable the aggregation, distribution, discovery, and acquisition of electronic publications. [(Wikipedia)](https://en.wikipedia.org/wiki/Open_Publication_Distribution_System)

__FLibGoLite__is multiplatform lightweight OPDS server with SQLite database book search index
Current __FLibGoLite__ release only supports FB2 publications, both individual files and zip archives.

OPDS-catalog is checked and works with mobile bookreaders applications FBReader and PocketBook Reader


__FLibGoLite__ program is written in GO as a single executable and doesn't require any prerequsites.  
All you have to do is to download, install and start it.

##  Download
---
Download your OS and Hardware specific program build

|OS        |Hardware              |Program executable          |Tested  |  
|----------|----------------------|----------------------------|:------:|  
|Windows   | Intel, AMD 32-bit    | flibgolite-linux-386.exe   |Yes     |  
|Windows   | Intel, AMD 64-bit    | flibgolite-linux-amd64.exe |Yes     |  
|OS X (MAC)| 64-bit               | flibgolite-darwin-64       |No      |  
|Linux     | Intel, AMD 32-bit    | flibgolite-linux-386       |No      |  
|Linux     | Intel, AMD 64-bit    | flibgolite-linux-amd64     |No      |  
|Linux     | ARM 32-bit (armhf)   | flibgolite-linux-arm-6     |Yes     |  
|Linux     | ARM 64-bit (armv8)   | flibgolite-linux-arm64     |Yes     |  
  
You may rename downloaded program executable to `flibgolite` or any other name you want.
For convenience, `flibgolite` name will be used below in this README.

## Install and start
---

Although __FLibGoLite__ program can be run from command line, the preferred setup is program to be installed as a system service running in background that will automaticaly start after power on or reboot.

Service installation and control requires administrator rights. On Linux you may use `sudo`.

On Windows open Powershell as Administrator and run commands to install, start and check service status

1. In Windows Powershell terminal run command

Install service:
```sh
  ./flibgolite -service install
```
Start service
```sh
  ./flibgolite -service start
```
And check that service is running
```sh
  ./flibgolite -service status
```

2. On Linux open terminal and run commands:

```bash
  sudo ./flibgolite -service install
  sudo ./flibgolite -service start
  sudo ./flibgolite -service status
```

If status is like "running" you can start to use it.

## Use
---

At the first run program will create the set of subfolders in the folder where program is located

 	flibgolite
	├─┬─ books  
	| ├─── new   - new book files and/or zip archives should be placed here to be added to the index
	| ├─── stock - indexed library book files and archives are stored here
	| └─── trash - files with processing errors will go here
	├─┬─ config - contains main configiration file config.yml and genre tree file
	| └─── locales - subfolder for localization files 
	├─── dbdata - database with book index resides here
	└─── logs - scan and opds rotating logs are here

Put your book files or book file zip archives in `books/new` folder and start to setup bookreader. Meanwhile  book decriptions will be added to book index of OPDS-catalog/

Set bookreader opds-catalog to `http://<server_name or IP-address>:8085/opds` to choose and download books on your device to read. See bookreader manual/help.

Tip: While searching book in bookreader use native keyboard layout to fill dearch pattern. For example, don't use Latin English "i" instead of Cyrillic Ukrainian "i", because it's not the same Unicode symbol. 

---
## Advanced usage

From command line run `./flibgolite -help` to see run options

	Usage:
	./flibgolite [OPTION]

	With no OPTION program will run in console mode (Ctrl+C to exit)  
	Caution: Only one OPTION can be used at a time

	OPTION should be one of:

	-service [action]     control FLibGoLite service
		where action is one of: install, start, stop, restart, uninstall, status  
	-reindex              empty book stock index and then scan book stock directory to add books to index
	-config               create default config file in ./config folder for customization and exit
	-help                 display this help and exit
	-version              output version information and exit


Examples:

```bash
./flibgolite                      	  Run FLibGoLite console mode
sudo ./flibgolite -service install    Install FLibGoLite as a system service
sudo ./flibgolite -service start	
```

Detalization

<details><summary><i><b>1. Main configuration file</i></b></summary>
<p>

For advanced sutup you can edit `config/config.yml` selfexplanatory configuration file.  
This file by default is located in `config` subfolder of program file location.  

</p>
</details>

<details><summary><i><b>2. Locations of folders setup</i></b></summary>
<p>

To change location of a folder just edit corresponding line in `config.yml`

For example, if you need to change the folder for new aquired books
```yml
NEW: "books/new"
``` 
just change `books/new` to the appropriate folder path.

</p>
</details>

<details><summary><i><b>3. OPDS tuning</i></b></summary>
<p>

You can change OPDS default 8085 http port to yours 
```yml
# OPDS-server port so opds can be found at http://<server name or IP-address>:8085/opds
PORT: 8085
```
You can change the number of books your bookreader will load when you page (pulldown the screen)

```yml
# OPDS feeds entries page size
PAGE_SIZE: 30
```
Do not set this value more than default. With lower values it updates faster.
</p>
</details>

<details><summary><i><b>4. Localization tips</i></b></summary>
<p>

There are some easy features that may help to tune your language experience

1. By default new books processing is limited to English, Russian and Ukrainian books. You can add [others](https://en.wikipedia.org/wiki/IETF_language_tag) like `"de"`, `"fr"`, `"it"` and so on.

```yml
# Accept only these languages publications. Add others if needed please.
ACCEPTED: "en, ru, uk"
```

2. By default bookreader will show menues and comments in English `"en"` If you are Rusiian or Ukranian you can change this setting to `"ru"` or `"uk`" 

```yml
# Default english locale for opds feeds (bookreaders opds menu tree) can be changed to:
# "uk" for Ukrainian, 
# "ru" for Russian 
DEFAULT: "uk"
```
3. If your native language is other then tree mentioned above for your convinience you can make language file and put it in `config/locales` folder

```yml
# Locales folder. You can add your own locale file there like en.yml, ru.yml, uk.yml
DIR: "config/locales"
```

For example, for German, copy `en.yml` to `de.yml` and translate the phrases into German to the right of the colon separator. Leave `%d` format symbols untouchced. Something like this:

```yml
Found authors - %d: Found Autoren gefunden - %d
```

Don't forget to replace alphabet string `ABC` to German. This ensures that the selections are in the correct alphabetical order.

4. Genres tree selection language adaptation can be done by editing the file `genres.xml` in `config` folder

```yml
  TREE_FILE: "config/genres.xml"
  # Alternative genres tree can be used (Russian only, sorry) 
  # TREE_FILE: "config/alt_genres.xml"
```

This can be done by adding language specific lines in `genres.xml` file

```xml
<genre-descr lang="en" title="Alternative history"/>
<genre-descr lang="ru" title="Альтернативная история"/>
<genre-descr lang="uk" title="Альтернативна історія"/>
<genre-descr lang="de" title="Alternative Geschichte"/>
```
</p>
</details>

<details><summary><i><b>5. Default config.yml</i></b></summary>
<p>

Default configuration file `config.yml` with folder tree is created at the first programm run. You can edit it and your edits will not be canceled the next time you run the program. Thus, you can distribute the files used by the program into the necessary folders. With reasonable care, you can edit or add any configuration file located by default in the `config` folder and it will not be deleted or overwriten.

```yml
library:
  # Selfexplained folders
  STOCK: "books/stock"
  NEW: "books/new"
  TRASH: "books/trash"

genres:
  TREE_FILE: "config/genres.xml"
  # Alternative genres tree can be used (Russian only, sorry) 
  # TREE_FILE: "config/alt_genres.xml"
  
database:
  DSN: "dbdata/flibgolite.db"
  # Delay before start each new acquisitions processing
  POLL_DELAY: 30 
  # Maximum simultaneous new aquisitios processing threads
  MAX_SCAN_THREADS: 3

logs:
  # Logs are here
  OPDS: "logs/opds.log"
  SCAN: "logs/scan.log"
  DEBUG: false 

opds:
  # OPDS-server port so opds can be found at http://<server name or IP-address>:8085/opds
  PORT: 8085
  # OPDS feeds entries page size
  PAGE_SIZE: 30

locales:
  # Locales folder. You can add your own locale file there like en.yml, ru.yml, uk.yml
  DIR: "config/locales"
  # Default english locale for opds feeds (bookreaders opds menu tree) can be changed to:
  # "uk" for Ukrainian, 
  # "ru" for Russian 
  DEFAULT: "uk"
  # Accept only these languages publications. Add others if needed please.
  ACCEPTED: "en, ru, uk"
```
</p>
</details>

<details><summary><i><b>6. Book index database</i></b></summary>
<p>

Book index is stored in SQLite database file located in `dbdata` folder. It is created at the first program run and __is not intended for manual editing__. 

```yml
  DSN: "dbdata/flibgolite.db"
```

</p>
</details>

<details><summary><i><b>7. Logging</i></b></summary>
<p>

While running program writes `opds.log` and `scan.log` located in `logs` folder.

```yml
OPDS: "logs/opds.log"
SCAN: "logs/scan.log"
```
`opds.log` contains records about bookreaders requests.  
`scan.log` contains records about new books and archive indexing.

You don't need to delete logs to free up disk space, as logs are rotated (overwrite) after 7 days.

</p>
</details>

<details><summary><i><b>8. Build from sources</i></b></summary>
<p>

If you have any security doubts about builded executables or there is no suitable one you may easily build it yourself.    
To build an executable install [Golang](https://go.dev/dl/), clone [FLibGoLite repositiry](https://github.com/vinser/flibgolite) and run `go build ./cmd/flibgolite`  
It's better to build it on the host the service will run. You will get executable right for the host OS and hardware.  
For crosscompile install and run GNU make with Makefile

</p>
</details>

---

*Comments and suggestions are welcome*

ANY CONCEPT CAN BE RETHINKED :)
   

