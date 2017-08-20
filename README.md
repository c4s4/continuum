CONTINUUM
=========

[![Build Status](https://travis-ci.org/c4s4/continuum.svg?branch=master)](https://travis-ci.org/c4s4/continuum)
[![Code Quality](https://goreportcard.com/badge/github.com/c4s4/continuum)](https://goreportcard.com/report/github.com/c4s4/continuum)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
<!--
[![Coverage Report](https://coveralls.io/repos/github/c4s4/continuum/badge.svg?branch=master)](https://coveralls.io/github/c4s4/continuum?branch=master)
-->

Continuum is a lightweight continuous integration tool: no web interface, no
scheduler. It runs on command line and is triggered by cron.

Installation
------------

Drop your platform executable *continuum_os_arch*, in the *bin* directory of
the distribution archive, somewhere in you *PATH* and rename it *continuum*.
For instance, on 64 bits Linux, you would  copy *continuum_linux_amd64* to 
*/usr/local/bin/continuum*.

Configuration
-------------

Configuration is in YAML format:

```yaml
directory: /tmp
status:    /tmp/continuum-status.yml
port:      6666
email:
  smtp-host: smtp.nowhere.com:25
  recipient: nobody@nowhere.com
  sender:    nobody@nowhere.com
  success:   true
  once:      true
modules:
  - name:    Continuum
    url:     https://github.com/c4s4/continuum.git
    branch:  develop
    command: |
      set -e
      make test
```

The first part indicates:

- **directory**: the directory where modules will be checked out. Currently only
  GIT projects are supported.
- **status**: this is the name of the file were are stored modules status (to 
  determine if their repository changed since last run and if last build was a
  success or a failure).
- **port**: the port that continuum listens to ensure that only one instance is
  running at a time. This port should be free on the host machine.
- **email**: put *~* if you don't want any email.

If you wait to receive email reports, provide following fields:

- **smtp-host**: the hostname and port of your SMTP server.
- **recipient**:  the email of the recipient of the build report.
- **sender**: the email address if the sender of the report.
- **success**: tells if continuum should send an email on success. If *false*,
  it will only send an email on build error.
- **once**: if you want to send a single mail while the status of a module
  changes.

The second part is made of the list of modules, with, for each module:

- **name**: the name of the module.
- **url**: the URL of the module that GIT will use to get the sources.
- **branch**: the branch to build (such as *master* or *develop*).
- **command**: the bash script to run tests, must return 0 on success and a 
  different value on error (as any Unix script should).

You can pass the configuration file to use on command line. If you pass no 
configuration file on command line, continuum will look for following files to
use:

- *~/.continuum.yml*
- */etc/gontinumm.yml*

Crontab
-------

This script is triggered using cron, with a configuration as follows (in file
*/etc/cron.d/continuum*):

```bash
# /etc/cron.d/continuum
# cron configuration to run gontinuum

SHELL=/bin/sh
PATH=/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin
MAILTO=""

# run continuum every 15 minutes
*/15 * * * *    user    continuum >> /tmp/continuum.log
```

This will run continuum every 15 minutes. When continuum starts, it checks if 
repository has changed for all modules, comparing its hash with the one stored
in *status* file.

If repository has changed, continuum clones it and runs command for tests. If 
script returns 0 (which is the Unix standard to tell that a command was
successful), the test is OK, else it is a failure.

Continuum prints a summary of the tests results and sends an email (or not
depending on email settings) for each test. Recommanded email configuration is
to set *success* and *once* to *true*. This will send an email when status of
a module changes (that is on test success when module was broken and test
failure when it was OK).

TODO
----

- Fix spam when failing to clone repository.

*Enjoy!*
