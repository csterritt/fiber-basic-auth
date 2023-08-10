### Fiber Basic Authentication

#### What this is

The idea is to build a web server using the [Fiber library](https://gofiber.io) which
does fast web serving for the Go programming language. It will use a very simple
password-less magic code authentication flow.

#### What you will need to use this

First, set up a file called `alg.info` which contains lines like the following:

    export SERVER_URL='http://localhost:3000'
    export SERVER_PORT=3000

You can, of course, set any `SERVER_PORT` you wish, just make sure that the port
number matches in the `SERVER_URL`.

Next, run:

    ./go

Now you can visit `http://localhost:3000` to see the basic page and do the whole sign in
or sign up flow (but you'll have to read the magic code off of the log. In production,
don't log that value, obviously!).
