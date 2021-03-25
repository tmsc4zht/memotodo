@rem this code based on https://github.com/zetamatta/nyagos/blob/master/make.cmd
@powershell "iex((@('')*3+(cat '%~f0'|select -skip 3))-join[char]10)"
@exit /b %ERRORLEVEL%

$VerbosePreference = "Continue"

$version = (git describe --tags --abbrev=0)
$revision = (git rev-parse --short HEAD)

Write-Verbose "Build as version='$version' revision='$revision'"

Write-Verbose "$ go build"
go build -o todo.exe -trimpath -ldflags "-s -w -X main.Version=$version -X main.Revision=$revision"