// Depends on goyaml:
// 
//   go get gopkg.in/yaml.v1
//
// Sample configuration file:
// 
//   directory:  /home/casa/tmp
//   email:
//     smtp_host: smtp.orange.fr
//     recipient: casa@sweetohm.net
//     sender:    casa@sweetohm.net
//     success:   true
//   
//   modules:
//     continuum:
//       url:     https://github.com/c4s4/continuum.git
//       command: |
//         set -e
//         export PATH=/opt/python/current/bin:$PATH
//         virtualenv env --no-site-packages
//         . env/bin/activate
//         pip install -r etc/requirements.txt
//         bee test

package main

import (
    "os"
    "fmt"
    "path"
    "os/exec"
    "io/ioutil"
    "path/filepath"
    "gopkg.in/yaml.v1"
)

type Config struct {
    Directory string
    Email struct {
        SmtpHost string "smtp_host"
        Recipient string
        Sender string
        Success bool
    }
    Modules map[string] struct {
        Url string
        Command string
    }
}

type Build struct {
    Success bool
    Output string
}

type Builds map[string]Build

func (builds Builds) Success() bool {
    for module := range(builds) {
        if !builds[module].Success {
            return false
        }
    }
    return true
}

func loadConfig(file string) Config {
    config := Config{}
    text, err := ioutil.ReadFile(file)
    if err != nil {
        panic(err.Error())
    }
    err = yaml.Unmarshal(text, &config)
    if err != nil {
        panic(err.Error())
    }
    // make directory absolute path
    absdir, err := filepath.Abs(config.Directory)
    if err != nil {
        panic(err.Error())
    }
    config.Directory = absdir
    return config
}

func buildModule(module string, config Config) Build {
    fmt.Printf("Building '%s'... ", module)
    module_dir := path.Join(config.Directory, module)
    // go in build directory
    err := os.Chdir(config.Directory)
    if err != nil {
        return Build {
            Success: false,
            Output: err.Error(),
        }
    }
    // delete module directory if it already exists
    if _, err := os.Stat(module_dir); err == nil {
        os.RemoveAll(module_dir)
    }
    // git clone the module repository
    cmd := exec.Command("git", "clone", config.Modules[module].Url)
    output, err := cmd.CombinedOutput()
    defer os.RemoveAll(module_dir)
    if err != nil {
        fmt.Println("ERROR")
        return Build{
            Success: false,
            Output: string(output),
        }
    } else {
        os.Chdir(module_dir)
        // run the build command
        cmd := exec.Command("bash", "-c", config.Modules[module].Command)
        output, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Println("ERROR")
            return Build {
                Success: false,
                Output: string(output),
            }
        } else {
            fmt.Println("OK")
            return Build{
                Success: true,
                Output: string(output),
            }
        }
    }
}

func buildModules(config Config) Builds {
    builds := make(Builds)
    for module := range(config.Modules) {
        builds[module] = buildModule(module, config)
    }
    return builds
}

func sendReport(builds Builds) {
    if builds.Success() {
        fmt.Println("OK")
    } else {
        fmt.Println("ERROR")
        for module := range(builds) {
            fmt.Println("===================================")
            fmt.Println(module)
            fmt.Println("-----------------------------------")
            fmt.Println(builds[module].Output)
            fmt.Println("-----------------------------------")
        }
    }
}

func main() {
    for i:=1; i<len(os.Args); i++ {
        config := loadConfig(os.Args[i])
        builds := buildModules(config)
        sendReport(builds)
    }
}

