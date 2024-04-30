<h1 align="center"> (っ◔◡◔)っ ♥ Flex ♥ </h1>

<p align="center">

<p align="center">
    <img src="https://img.shields.io/github/license/hikitani/flex" alt="License">
    <img alt="Static Badge" src="https://img.shields.io/badge/lang-golang-blue">
    <a href='https://coveralls.io/github/hikitani/flex?branch=main'><img src='https://coveralls.io/repos/github/hikitani/flex/badge.svg?branch=main' alt='Coverage Status' /></a>
    <a href="https://goreportcard.com/report/github.com/hikitani/flex"><img src="https://goreportcard.com/badge/github.com/hikitani/flex" alt="Go Reference"></a>
    <a href="https://pkg.go.dev/github.com/hikitani/flex"><img src="https://pkg.go.dev/badge/github.com/hikitani/flex.svg" alt="Go Reference"></a>
</p>

</p>

<p align="center"> Library for introspection of structures.</p>
<p align="center">Allows you to examine the structure of objects (including private fields) at runtime.</p>

<h2 align="center"> :warning: Important </h2>

<p align="center">
<a href="https://git.io/typing-svg"><img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&pause=1000&color=F75656&center=true&random=false&width=435&lines=Not+recommended+for+production+use;Use+the+library+with+caution" alt="Typing SVG" /></a>
<p>

<p align="center">This library is written for fun, and the original purpose of writing it is to better understand the go language.</p>


<h2 align="center"> :sparkles: Functions </h2>

* `StructToMap[T any]` - Getting all fields (as well as nested fields) as a map.
* `ValuesOf[Target, From any]` - Search for all values matching the specified Target type.
* `FieldValue[T any]` - Returning a value in a structure along a specified path.

<h2 align="center"> :arrow_down: Install </h2>

<p align="center"><code>go get github.com/hikitani/flex</code></p>
