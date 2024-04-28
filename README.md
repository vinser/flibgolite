FLibGoLite - Just enough for free OPDS 
===

### CURRENT STABLE RELEASE v2.0.0

### *Attention users of version 1 !!!*   
The version 2 update has some improvements to the index structure, configuration and localization. Therefore, to switch to version 2, you should uninstall the version 1 service and install the version 2 service from scratch. Rendexing of the book stock in version 2 is an order of magnitude faster than in version 1, especially for FB2 archives, so the transition to version 2 can be completed with acceptable service downtime.

---
__FLibGoLite__ is easy to use home library OPDS server you can install on your PC.

>The Open Publication Distribution System (OPDS) catalog format is a syndication format for electronic publications based on Atom and HTTP. OPDS catalogs enable the aggregation, distribution, discovery, and acquisition of electronic publications. [(Wikipedia)](https://en.wikipedia.org/wiki/Open_Publication_Distribution_System)

__FLibGoLite__ is multiplatform lightweight OPDS server with SQLite database book search index.

Current __FLibGoLite__ release supports [EPUB](https://en.wikipedia.org/wiki/EPUB) and [FB2 (single files and zip archives)](./pkg/fb2/LICENSE) publications format.

__FLibGoLite__ OPDS catalog has been tested and works with mobile book reader applications PocketBook Reader, FBReader, Librera Reader, Cool Reader, as well as desktop applications Foliate and Thorium Reader. You can use any other applications or e-ink devices that can read the listed book formats and work with OPDS catalogs.

__FLibGoLite__ program is written in GO as a single executable and doesn't require any prereqiusites.  
__All you have to do is to download, install and start it.__

##  Download
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

## Install and start
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
At the first run program will create the set of subfolders in the folder where program is located

 	flibgolite
	├─┬─ books  
	| ├─── stock - library book files and archives are stored here
	| └─── trash - files with processing errors will go here
	├─┬─ config - contains main configiration file config.yml and genre tree file
	| └─── locales - subfolder for localization files 
	├─── dbdata - database with book index resides here
	└─── logs - scan and opds rotating logs are here

Put your book files or book file zip archives in `books/stock` folder and start to setup bookreader. Meanwhile book descriptions will be added to book index of OPDS-catalog.

Set bookreader opds-catalog to `http://<PC_name or PC_IP_address>:8085/opds` to choose and download books on your device to read. See bookreader manual/help.

`Tip:` While searching book in bookreader use native keyboard layout for choosed language to fill search pattern. For example, don't use Latin English "i" instead of Cyrillic Ukrainian "i", because it's not the same Unicode symbol. 

## Advanced usage
To understand the features of fine-tuning __FLibGoLite__ application see the [Advanced User Guide](https://vinser.github.io/flibgolite/)

---
___*Suggestions, bug reports and comments are welcome [here](https://github.com/vinser/flibgolite/issues)*___

   

