# GO MobSf

GO MobSf is a Golang module designed to automate MobSF APK scans. This module allows you to run MobSF, configure the API key, set input and output targets, and start the scanning process with ease.

![alt text](https://raw.githubusercontent.com/edenbwt/Golang_Mobsf_Processor/main/img/Capture%20d%E2%80%99%C3%A9cran%202024-06-24%20151042.png)

## Configuration

Before running the module, ensure MobSF is up and running. Obtain the API key from MobSF and update the configuration file accordingly.

### Config File (`config.json`)

The `config.json` file in the project root directory contains the following fields:

```json
{
    "apiKey": "API KEY",
    "domain": "DOMAIN TO MOBSF",
    "InFolder": "input",
    "dump": "dump",
    "ext": [".apk", ".ipa", ".xapp"],
    "Version": "1.0.1.3 MULTI UPLOAD"
}
```


1. apiKey: Replace "API KEY" with your MobSF API key.
2. domain: Replace "DOMAIN TO MOBSF" with the domain where MobSF is hosted.
3. InFolder: This is the directory where APK files (or other specified extensions) are located for input to MobSF.
4 .dump: This directory is where MobSF will store the scan reports.
5. ext: An array of file extensions (.apk, .ipa, .xapp) that MobSF will scan.
6 .Version: Optional field indicating the current version of the configuration.
   
## Usage

Setup: Ensure MobSF is running and the API key is correctly configured in config.json.

Configuration: Verify that the InFolder and dump directories are set up correctly for input and output.

Run: Execute the module with the following command:

`go run main.go`

Scan Process: Enjoy the automated scanning process! MobSF will use the configured settings to scan the files in InFolder, generate reports, and save them in the dump directory.

if you enjoy ;)
![alt text](https://raw.githubusercontent.com/edenbwt/Golang_Mobsf_Processor/main/img/mybtc.png)
