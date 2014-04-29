=========
GONTINUUM
=========

Gontinuum is a lightweight continuous integration tool: no web interface, no scheduler. It runs on command line and is triggered by cron.

Installation
============

TBD.

Configuration
=============

Configuration is in YAML format::

  directory:   /tmp
  repo_hash:   /tmp/repo-hash.yml
  port:        6666
  email:
    smtp_host: smtp.foo.com:25
    recipient: foo@bar.com
    sender:    foo@bar.com
    success:   false
  modules:
    - name:    module1
      url:     https://github.com/user/module1.git
      command: |
        set -e
        commands to run tests
    - name:    module2
      url:     https://github.com/user/module2.git
      command: |
        set -e
        commands to run tests

The first part indicates:

- directory: the directory where modules will be checked out. Currently only GIT projects are supported.
- repo_hash: this is the name of the file were are stored repositories hash (to determine if they changed since last run).
- port: the port that gontinuum listens to ensure that only one instance is running at a time. This port should be free on the host machine.
- email: put *~* if you don't want any email.

If you wait to receive email reports, provide following fields:

- smtp_host: the hostname and port of your SMTP server.
- recipient:  the email of the recipient of the build report.
- sender: the email address if the sender of the report.
- success: tells if continuum should send an email on success. If *false*, it will only send an email on build error.

The second one is the list of modules, with, for each module:

- name: the name of the module.
- url: the URL of the module that GIT will use to get the sources.
- command: the command to run tests, must return 0 on success and a different value on error (as any Unix script should).

You can pass the configuration file to use on command line. If you pass no configuration file on command line, gontinuum will look for following files to use:

- ~/.gontinuum.yml
- ~/etc/gontinuum.yml
- /etc/gontinumm.yml

Crontab
=======

This script is triggered using cron, with as configuration as follows (in file */etc/crontab*)::

  0   * * * *  me    gontinuum

This will run gontinuum every hour. When gontinuum starts, it checks if repository has changed for all modules, comparing its hash with the one stored in *repo_hash* file.

If repository has changed, gontinuum clones it and runs command for tests. If commands return 0 (which is the Unix standard to tell that a command was successful), the test is OK, else it is a failure.

Gontinuum prints a summary of the tests results and sends an email if one test failed. It also sends a report if no test failed and *success* configuration field was set to *true*.

Todo
====

Here is a list of what is planned in the future:

- Manage other SCM (such as SVN and CVS).

Releases
========

- **0.1.0** (*2014-04-??*): First public release.

Enjoy!

