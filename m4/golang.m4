
## ----------------------- ##
## 1. Language selection.  ##
## ----------------------- ##

# AC_LANG(Golang)
# -----------
AC_LANG_DEFINE([Golang], [go], [GOLANG], [GOCMD], [],
[
ac_ext=go
ac_compile='$GOCMD build conftest.$ac_ext >&AS_MESSAGE_LOG_FD'
ac_link='$GOCMD build conftest.$ac_ext >&AS_MESSAGE_LOG_FD'
])

# AC_LANG_GOLANG
# ----------
AU_DEFUN([AC_LANG_GOLANG], [AC_LANG(Golang)])

# ------------------- #
# Producing programs.
# ------------------- #

# AC_LANG_PROGRAM(Golang)([PROLOGUE], [BODY])
# ---------------------------------------
m4_define([AC_LANG_PROGRAM(Golang)],
[package main
$1
func main() {
$2
}])

# _AC_LANG_NULL_PROGRAM(Golang)
# -----------------------------
# Produce source that does nothing.
m4_define([_AC_LANG_NULL_PROGRAM(Golang)],
[AC_LANG_PROGRAM([],
[])])

# _AC_LANG_IO_PROGRAM(Golang)
# -----------------------
# Produce source that performs I/O.
m4_define([_AC_LANG_IO_PROGRAM(Golang)],
[AC_LANG_PROGRAM([import ( "fmt"; "os" )],
[f, err := os.Create("conftest.out")
 if err != nil {
     fmt.Println(err)
     os.Exit(1)
 }
 defer f.Close()
 os.Exit(0)
])])

## -------------------------------------------- ##
## 3. Looking for Compilers and Preprocessors.  ##
## -------------------------------------------- ##


AC_DEFUN([AC_PROG_GOLANG],
[AC_LANG_PUSH(Golang)dnl

AC_PATH_PROG(GOCMD, go, false, [${PATH}${PATH_SEPARATOR}/usr/local/go/bin], false)
_AS_ECHO_LOG([checking for _AC_LANG compiler version])
set X $ac_compile
ac_compiler=$[2]
_AC_DO_LIMIT([$ac_compiler version >&AS_MESSAGE_LOG_FD])


m4_expand_once([_AC_COMPILER_EXEEXT])[]dnl

AC_LANG_POP(Golang)dnl
])# AC_PROG_GOLANG
