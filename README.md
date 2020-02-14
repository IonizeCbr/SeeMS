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
Run the command `go build` from the cloned directory. 


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

### Example Output
```
.\SeeMS.exe -filename .\hosts.txt
Sharepoint (16.0.0.19722) @ https://ionize.sharepoint.com/
Moodle (Unknown) @ https://moodle.telt.unsw.edu.au/
No recognised CMS @ https://google.com
Joomla (3.9.15) @ https://joomla.org
Drupal (7.61 - 7.67) @ https://drupal.org
Wordpress (5.3.2) @ https://ionize.com.au
```
