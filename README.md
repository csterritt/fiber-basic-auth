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
    export COOKIE_KEY='A 32-char, base-64-encoded string'

For the `COOKIE_KEY`, you can use:

    openssl rand -base64 32

to generate a usable value. This should be different for development and production.

You can, of course, set any `SERVER_PORT` you wish, just make sure that the port
number matches in the `SERVER_URL`.

Next, run:

    ./go

Now you can visit `http://localhost:3000` to see the basic page and do the whole sign in
or sign up flow (but you'll have to read the magic code off of the log. In production,
don't log that value, obviously!).

## Testing

If you want to run the end-to-end tests, make sure you have Node version 16 or later
installed, and run:

    npm install

Then you can run all the end-to-end tests via the Playwright testing environment with
the command:

    ./run-e2e-tests

You can run one specific test by naming it on the command line, e.g.

    ./run-e2e-tests e2e-tests/visit-sign-up-return.spec.ts

You can run the Playwright `codegen` tool with the `codegen` script, or bring up the
Playwright UI (which allows inspection of tests at each step) with the `ui-tests`
script.

## Production use

**NOTE**: There is a script here named `prod_deploy`. Currently, it's set up to do a
Linux build (via the `linux_build` script), and then fail. The `linux_build` script
removes any lines with a `PRODUCTION:REMOVE` comment, then compiles for a linux host.
This type of host may not be what you want, but you **definitely** want the
`PRODUCTION:REMOVE` lines to be removed in any code you deploy to production.
Before doing this removal, the script verifies that there is no code currently checked
out. Then it removes those lines, compiles the code, and does a `git reset` on any
modified files. Finally, if that succeeds, `prod_deploy` deploys your code to your
production server(s) or service(s), or at least it will once you write that code!
For now, it just `echo`s success and exits.

First, you'll have to figure out some way to get the magic code to your
users! There's nothing here (yet) to support that. All the text mentions emails;
if you want to use something else, like text messages, change that too.

Second, the codes expire in twenty minutes. You can change this duration by
setting the `CodeExpireTimeInSeconds` value in `constants.constants.go`.

Third, to prevent someone trying *every possible code*, the user will get a failure
if they enter a wrong code more than three times. You can change the failure  count
by setting the `WrongCodeFailureCount` value in `constants.constants.go`.

Fourth, search the source for the phrase `PRODUCTION:`. You'll find a number of places
where it's *extremely highly recommended* that you take some particular action. Do what
it says. Note that `PRODUCTION:REMOVE` lines will be automatically removed by the
`linux_build` script, as mentioned above.

Fifth, using the "sqlite3.Storage" engine is probably not a great idea, particularly
if you're running on a serverless-style hosting environment where your sqlite3 database
could suddenly disappear, since the service started a new host. Databases are supposed
to be durable... let them do their job.

Sixth, sign-up should direct our new user to a welcome page with some basic info. Or
maybe it's just a pop-up/toast/other notification.

Seventh, allow resubmitting the code (but only after a suitable interval).

Eighth, when the session fails to save, redirect to a 500 error page, notify admins.

## License

See the LICENSE.txt file.
