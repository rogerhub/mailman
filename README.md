SMTP Mailman
============

What is this?
-------------

An implementation of rfc5321, the latest revision of the Simple Mail Transfer Protocol, in Go. I might set up a honeypot with this some time in the future.

Why are you making this??
-------------------------

You clearly don't understand how bored I am right now.

There is a smtp package in the language already..
-------------------------------------------------

Yeah, so what? I thought it would be interesting to try to implement as many features of the specification as I can with just a few basic packages.

Proof of concept
----------------

Here's a sample exchange. Commands were taken from the official spec.

    roger [21:52:35] ~/Development/smtpserver/src/github.com/rogerhub
    $ nc localhost 2225
    220 localhost Greetings
    EHLO foo.com
    250-localhost
    250-PIPELINING
    250 SIZE 10000000
    MAIL FROM:<JQP@bar.com>
    250 OK
    RCPT TO:<Jones@XYZ.COM>
    250 OK
    DATA
    354 Start mail input; end with <CRLF>.<CRLF>
    Received: from bar.com by foo.com ; Thu, 21 May 1998
    	05:33:29 -0700
    Date: Thu, 21 May 1998 05:33:22 -0700
    From: John Q. Public <JQP@bar.com>
    Subject:  The Next Meeting of the Board
    To: Jones@xyz.com

    Bill:
    The next meeting of the board of directors will be
    on Tuesday.
    						John.
    .
    250 OK
    QUIT
    221 Service closing transmission channel
