---
layout: page
title: Advanced User Guide
---

__FLibGoLite__ is easy to use home library OPDS server 

>The Open Publication Distribution System (OPDS) catalog format is a syndication format for electronic publications based on Atom and HTTP. OPDS catalogs enable the aggregation, distribution, discovery, and acquisition of electronic publications. [(Wikipedia)](https://en.wikipedia.org/wiki/Open_Publication_Distribution_System)

__FLibGoLite__ is multiplatform lightweight OPDS server with SQLite database book search index.

Current __FLibGoLite__ release supports [EPUB](https://en.wikipedia.org/wiki/EPUB) and [FB2 (single files and zip archives)](./pkg/fb2/LICENSE) publications format.

__FLibGoLite__ OPDS catalog has been tested and works with mobile book reader applications PocketBook Reader, FBReader, Librera Reader, Cool Reader, as well as desktop applications Foliate and Thorium Reader. You can use any other applications or e-ink devices that can read the listed book formats and work with OPDS catalogs.

__FLibGoLite__ program is written in GO as a single executable and doesn't require any prereqiusites.  
__All you have to do is to download, install and start it.__

###  Download
[Download latest release](https://github.com/vinser/flibgolite/releases/tag/v2.0.0) of specific program build for your OS and CPU type  

|OS        |CPU type              |Program executable          |Tested<sup>1</sup> |  
|----------|----------------------|----------------------------|:------:|  
|Windows   | Intel, AMD 64-bit    | flibgolite-linux-amd64.exe |Yes     |  
|OS X (MAC)| Intel, AMD 64-bit    | flibgolite-darwin-64       |No      |  
|OS X (MAC)| ARM 64-bit           | flibgolite-darwin-64       |No      |  
|Linux     | Intel, AMD 64-bit    | flibgolite-linux-amd64     |No      |  
|Linux     | ARM 32-bit (armhf)   | flibgolite-linux-arm-6     |Yes     |  
|Linux     | ARM 64-bit (armv8)   | flibgolite-linux-arm64     |Yes     |  

<sup>1</sup>_Some of executables was only cross-builded and not tested on real desktops, but you can still try them out_  

You may rename downloaded program executable to `flibgolite` or any other name you want.
For convenience, `flibgolite` name will be used below in this README.

### Install and start
Although __FLibGoLite__ program can be run from command line, the preferred setup is program to be installed as a system service running in background that will automaticaly start after power on or reboot.

Service installation and control requires administrator rights.

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

2. On Linux open terminal and run commands using `sudo`:

```bash
  sudo ./flibgolite -service install
  sudo ./flibgolite -service start
  sudo ./flibgolite -service status
```

If status is like "running" you can start to use it.

### Use
At the first run program will create the set of subfolders in the folder where program is located

 	flibgolite
	├─┬─ books  
	| ├─── stock - library book files and archives are stored here
	| └─── trash - files with processing errors will go here
	├─┬─ config - contains main configuration file config.yml and genre tree file
	| └─── locales - subfolder for localization files 
	├─── dbdata - database with book index resides here
	└─── logs - scan and opds rotating logs are here

Put your book files or book file zip archives in `books/stock` folder and start to setup bookreader. Meanwhile book descriptions will be added to book index of OPDS-catalog.

Set bookreader opds-catalog to `http://<PC_name or PC_IP_address>:8085/opds` to choose and download books on your device to read. See bookreader manual/help.

`Tip:` While searching book in bookreader use native keyboard layout for choosed language to fill search pattern. For example, don't use Latin English "i" instead of Cyrillic Ukrainian "i", because it's not the same Unicode symbol. 

### Advanced usage
From command line run `./flibgolite -help` to see run options
```
Usage: flibgolite [OPTION] [data directory]

With no OPTION program will run in console mode (Ctrl+C to exit)
Caution: Only one OPTION can be used at a time

OPTION should be one of:
  -service [action]     control FLibGoLite system service
          where action is one of: install, start, stop, restart, uninstall, status 
  -reindex              empty book stock index and then scan book stock directory to add books to index (database)
  -config               create default config file in ./config folder for customization and exit
  -help                 display this help and exit
  -version              output version information and exit

data directory is optional (current directory by default)
```

Examples:

```bash
./flibgolite                          Run FLibGoLite console mode
sudo ./flibgolite -service install    Install FLibGoLite as a system service
sudo ./flibgolite -service start	
```

### Detalization

#### _1. Main configuration file_

For advanced sutup you can edit `config/config.yml` selfexplanatory configuration file.  
This file by default is located in `config` subfolder of program file location.

#### _2. Locations of folders setup_

To change location of a folder just edit corresponding line in `config.yml`

For example, if you need to setup separate folder for new aquired books uncomment line  

```yml
NEW: "books/new"
``` 

and change `books/new` to the appropriate folder path.

#### _3. OPDS tuning_

You can change OPDS default 8085 http port to yours 
```yml
# OPDS-server port so opds can be found at http://<server name or IP-address>:8085/opds
PORT: 8085
```
Here you can set OPDS-server preferred name
```yml
# OPDS-server title that is displayed in a book reader
TITLE: "FLib Go Go Go!!!"
```
You can change the number of books your bookreader will load at a time when you page (pulldown/update the screen)

```yml
# OPDS feeds entries page size
PAGE_SIZE: 30
```
Do not set this value more than default. With lower values it updates faster.

#### _4. Localization tips_

There are some easy features that may help to tune your language experience

4.1. By default new books processing is limited to English, Russian and Ukrainian books. You can add [others](https://en.wikipedia.org/wiki/IETF_language_tag) like `"de"`, `"fr"`, `"it"` and so on.  

```yml
# Accept only these languages publications. Add others if needed please.
ACCEPTED: "en, ru, uk"
```  

4.2. By default bookreader will show menues and comments in English `"en"` If you are Rusiian or Ukranian you can change this setting to `"ru"` or `"uk`" 

```yml
# Default english locale for opds feeds (bookreaders opds menu tree) can be changed to:
# "uk" for Ukrainian, 
# "ru" for Russian 
DEFAULT: "en"
```

4.3. If your native language is other then three mentioned above for your convinience you can make language file and put it in `config/locales` folder  

```yml
# Locales folder. You can add your own locale file there like en.yml, ru.yml, uk.yml
DIR: "config/locales"
```

For example, for German, copy `en.yml` to `de.yml` and translate the phrases into German to the right of the colon separator. Leave `%d` format symbols untouchced. Something like this:  

```yml
Found authors - %d: Autoren gefunden - %d
```

Don't forget to replace alphabet string `ABC` to German. This will ensure that German names and titles are displayed and sorted correctly.  

4.4. Genres tree selection language adaptation can be done by editing the file `genres.xml` in `config` folder

```yml
  TREE_FILE: "config/genres.xml"
```

This can be done by adding language specific lines in `genres.xml` file

```xml
<genre-descr lang="en" title="Alternative history"/>
<genre-descr lang="ru" title="Альтернативная история"/>
<genre-descr lang="uk" title="Альтернативна історія"/>
<genre-descr lang="de" title="Alternative Geschichte"/>
```

#### _5. Default config.yml_

Default configuration file `config.yml` with folder tree is created at the first programm run. You can edit it and your edits will not be canceled the next time you run the program. Thus, you can distribute the files used by the program into the necessary folders. With reasonable care, you can edit or add any configuration file located by default in the `config` folder and it will not be deleted or overwriten.

```yml
library:
  # Selfexplained folders
  STOCK: "books/stock" # Book stock
  TRASH: "books/trash" # Error and duplicate files and archives wil be moved to this folder 
  # NEW: "books/new" # Uncomment the line to have separate folder for new acquired books

genres:
  TREE_FILE: "config/genres.xml"
  # Alternative genres tree can be used (Russian only, sorry) 
  # TREE_FILE: "config/alt_genres.xml"
  
database:
  DSN: "dbdata/books.db"
  # Delay before start each new acquisitions processing
  POLL_DELAY: 30 
  # Maximum simultaneous new aquisitios processing threads
  MAX_SCAN_THREADS: 30

logs:
  # Logs are here
  # To redirect the log output to console (stdout) just comment out the appropriate line OPDS or SCAN
  OPDS: "logs/opds.log"
  SCAN: "logs/scan.log"
   # Logging levels: D - debug, I - info, W - warnings (default), E - errors
  LEVEL: "W" 

opds:
  # OPDS-server port so opds can be found at http://<server name or IP-address>:8085/opds
  PORT: 8085
  # OPDS-server title that is displayed in a book reader
  TITLE: "FLib Go Go Go!!!"
  # OPDS feeds entries page size
  PAGE_SIZE: 30

locales:
  # Locales folder. You can add your own locale file there like en.yml, ru.yml, uk.yml
  DIR: "config/locales"
  # Default english locale for opds feeds (bookreaders opds menu tree) can be changed to:
  # "uk" for Ukrainian, 
  # "ru" for Russian 
  DEFAULT: "en"
  # Accept only these languages publications. Add others if needed please.
  ACCEPTED: "en, ru, uk"
```

#### _6. Book index database_

Book index is stored in SQLite database file located in `dbdata` folder. It is created at the first program run and __is not intended for manual editing__. 

```yml
DSN: "dbdata/books.db"
```

#### _7. Logging_

While running program writes `opds.log` and `scan.log` located in `logs` folder.

```yml
OPDS: "logs/opds.log"
SCAN: "logs/scan.log"
```
`opds.log` contains records about bookreaders requests.  
`scan.log` contains records about new books and archive indexing.  
To redirect the log output to console (stdout) just comment out the appropriate line OPDS or SCAN.

You don't need to delete logs to free up disk space, as logs are rotated (overwrite) after 7 days.

You can setup logging level (verbosity) to one of: `D` - debug, `I` - info, `W` - warnings (default), `E` - errors
```yml
LEVEL: "W" 
```

#### _8. Run in Docker container_

As an option you may run program in [docker container](README.docker.md)

#### _9. Build from sources_

If you have any security doubts about builded executables or there is no suitable one you may easily build it yourself.    
To build an executable install [Golang](https://go.dev/dl/), [Git](https://git-scm.com/downloads) clone [FLibGoLite repositiry](https://github.com/vinser/flibgolite) and run `go build ./cmd/flibgolite`  
It's better to build it on the host the service will run. You will get executable right for the host OS and hardware.  
For crosscompile install GNU `make` and run it with Makefile


-------------------------------
___*Suggestions, bug reports and comments are welcome [here](https://github.com/vinser/flibgolite/issues)*___

   

