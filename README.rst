=========
GONTINUUM
=========

Gontinuum is a lightweight continuous integration tool: no web interface, no scheduler. It runs on command line and is triggered by cron.

Installation
============



Configuration
=============

Configuration is in YAML format::

  directory:  /tmp
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
      command:
        set -e
        commands to run tests

The first part indicates:

- directory: the directory where modules will be checked out. Currently only GIT projects are supported.
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

Crontab
=======

This script is triggered using cron, with as configuration as follows (in file */etc/crontab*)::

  # run gontinuum at 4 every night
  0   4 * * *  me    gontinuum ~/etc/gontinuum.yml

Todo
====

Here is a list of what is planned in the future:

- Manage other SCM (such as SVN and CVS).
- Implement conditional build (when project was modified).

Releases
========

- **0.1.0** (*2014-04-??*): First public release.

Enjoy!

