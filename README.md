## SeeMS

This tool is designed to be a high level scanner for identifying the CMS type and version at any given URL. 

Currently Supports
- Drupal
- Wordpress
- Sharepoint
- Joomla
- Moodle

Still under active development.


### Installation
Run the command `go build` from the cloned directory. 


### Command Line
```
Usage of SeeMS.exe:
  -filename string
        File name of a list of targets. One per line.
  -target string
        URL of an individual target you wish to scan.
  -threads int
        Number of threads to use. (default 10)
```

### Target / File Format
SeeMS expects targets to be in a format of one per line with protocol specified. 

**Accepted**
```
https://foobar.com.au
http://foobar.com:8080
```

**Not Accepted**
```
foobar.com.au
```
