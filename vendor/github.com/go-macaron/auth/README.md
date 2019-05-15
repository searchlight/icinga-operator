# auth 
Macaron middleware/handler for http basic authentication. Modified from <https://github.com/martini-contrib/auth>

[API Reference](http://godoc.org/github.com/go-macaron/auth)

## Simple Usage

Use `auth.Basic` to authenticate against a pre-defined username and password:

~~~ go
import (
  "gopkg.in/macaron.v1"
  "github.com/go-macaron/auth"
)

func main() {
  m := macaron.Classic()
  // authenticate every request
  m.Use(auth.Basic("username", "secretpassword"))
  m.Run()
}
~~~

## Advanced Usage

Using `auth.BasicFunc` lets you authenticate on a per-user level, by checking
the username and password in the callback function:

~~~ go
import (
  "gopkg.in/macaron.v1"
  "github.com/go-macaron/auth"
)

func main() {
  m := macaron.Classic()
  // authenticate every request
  m.Use(auth.BasicFunc(func(username, password string) bool {
    return username == "admin" && password == "guessme"
  }))
  m.Run()
}
~~~

Note that checking usernames and passwords with string comparison might be
susceptible to timing attacks. To avoid that, use `auth.SecureCompare` instead:

~~~ go
  m.Use(auth.BasicFunc(func(username, password string) bool {
    return auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "guessme")
  }))
}
~~~

Upon successful authentication, the username is available to all subsequent
handlers via the `auth.User` type:

~~~ go
  m.Get("/", func(user auth.User) string {
    return "Welcome, " + string(user)
  })
}
~~~

## Authors
* [Jeremy Saenz](https://github.com/codegangsta)
* [Brendon Murphy](https://github.com/bemurphy)
* [codeskyblue](https://github.com/codeskyblue)
