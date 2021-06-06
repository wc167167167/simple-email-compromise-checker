# simple-email-compromise-checker

This is an simple sample of a checker API that listens at port.

To build and run it, please use a Mac or Windows machine just cd to the project root
and do a `sh ./buildAndStart.sh`

It builds the project as a docker image and also pulls another database image and run them.
They will be listening to port `80`, `8987`, `8988`, `8989`, and `8990` so make sure the
ports are free.

As a simple example project it does not really maintain the data, which means that the data are
'mocked' and even attempt to add more data into the db, once the images are down they will be
cleaned up.

The 'mocked' compromised emails are in forms of `jx${number}compromised@fake_email.com`, where
the number shall be from 0 to 99, for instance `jx10compromised@fake_email.com`

To test the API you can simply do a `curl localhost/check?email=test@some.com`.
