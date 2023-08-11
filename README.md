### Fiber Basic Authentication

# DO NOT USE THIS AS IS. You've been warned. See below for more on Production use.

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

#### Production use

First of all, you'll have to figure out some way to get the magic code to your
users! There's nothing here (yet) to support that.

Second, the codes don't expire in twenty minutes or so. Gotta fix that.

Third, there is no check that someone is just trying *every possible code*. Really,
it should fail hard (i.e., get rid of the current code and make you try again) if
you fail to enter the code properly some small number of times (like, three). No,
really, someone will write a bot to try all 900,000 different codes. Right now, that
would work (eventually).

Fourth, search the source for the phrase "PRODUCTION:". You'll find a number of places
where it's *extremely highly recommended* that you take some particular action. Do what
the code says.

Fifth, using the "sqlite3.Storage" engine is probably not a great idea, particularly
if you're running on a serverless-style hosting environment where your sqlite3 database
could suddenly disappear, since the service started a new host. Databases are supposed
to be durable... let them do their job.
