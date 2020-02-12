## SeeMS

This tool is designed to be a high level scanner for identifying the CMS type and version at any given URL. 

Currently Supports
- Drupal (w/ Versions)
- Wordpress (w/ Versions)
- Sharepoint (w/ Versions)
- Joomla (w/ Versions)
- Moodle (No Version Information)

Still under active development.


### Installation
Run the commands `go get` followed by `go build` from the cloned directory. 


### Command Line
```
Usage of SeeMS.exe:
  -filename string
        File name of a list of targets. One per line.
  -quiet
        Option to suppress status messages.
  -target string
        URL of an individual target you wish to scan.
  -threads int
        Number of threads to use. (default 10)
```